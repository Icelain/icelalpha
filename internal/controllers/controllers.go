package controllers

import (
	"context"
	"encoding/json"
	"icealpha/internal/controllers/jwtauth"
	"icealpha/internal/controllers/oauth"
	"icealpha/internal/controllers/user"
	"icealpha/internal/database"
	"icealpha/internal/router"
	"icealpha/internal/types"
	"log/slog"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

// Start all controllers and configure cookie store
func HandleAll(r *router.Router) {

	// set oauth2 config
	oauth.SetGithubOAuthConfig()

	// complete database syncing
	go func() {

		ticker := time.NewTicker(time.Minute * 5)

		for range ticker.C {

			database.Sync(r.S.DB, r.S.CreditCache)

		}

	}()

	r.R.Get("/api", HandleAPIIndex(r))

	r.R.Post("/api/user/handleimage", user.AuthMiddleware(user.HandleSolveInputImage(r), r))
	r.R.Post("/api/user/handletext", user.AuthMiddleware(user.HandleSolveTextInput(r), r))
	r.R.Post("/api/user/test", user.AuthMiddleware(user.TestController(r), r))
	r.R.Post("/api/user/nauthtest", user.TestController(r))

	r.R.Get("/api/oauth", HandleOAuthFlow(r))
	r.R.Get("/api/oauth/logout", HandleOAuthLogout(r))
	r.R.Get("/api/oauth/{provider}/callback", HandleOAuthCallback(r))
}

// GET :: -> Json(status: string)
func HandleAPIIndex(rtr *router.Router) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if err := json.NewEncoder(w).Encode(map[string]string{

			"icelalpha api status": "up",
		}); err != nil {

			http.Error(w, "error while serving /api", http.StatusInternalServerError)
			rtr.Logger.Error("Failed to serve /api", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})

		}

	}

}

// GET :: -> Redirect(url)
func HandleOAuthFlow(rtr *router.Router) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		jwtToken := r.URL.Query().Get("jwtToken")

		if jwtToken != "" {

			if token, ok := rtr.S.JwtSession.TokenPool.Load(jwtToken); ok {
				if jwtToken, err := jwtauth.VerifyToken(token.(string), rtr.S.JwtSession.SecretKey); err == nil {

					if _, err := jwtToken.Claims.GetSubject(); err == nil {

						http.Redirect(w, r, "/api", http.StatusTemporaryRedirect)
						return

					}

				}
			}

		}

		provider := r.URL.Query().Get("provider")

		switch provider {

		case "github":

			state := oauth.SetNewOAuthStateCookie(w)
			url := oauth.GithubOAuthConfig.AuthCodeURL(state, oauth2.ApprovalForce)
			http.Redirect(w, r, url, http.StatusTemporaryRedirect)

		}

	}

}

// GET :: -> SessionCookie | Redirect
func HandleOAuthCallback(rtr *router.Router) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		provider := r.PathValue("provider")

		switch provider {

		case "github":

			var githubUser oauth.GithubUser
			var err error
			switch provider {
			case "github":
				githubUser, _, err = oauth.HandleGithubOAuthCallback(rtr, oauth.GithubOAuthConfig, w, r)
				if err != nil {
					rtr.Logger.Error("err handling github oauth callback", "err", err)
					//http.Redirect(w, r, "/login", http.StatusSeeOther)
					return
				}

				if !rtr.S.DB.CheckUserExists(context.Background(), githubUser.Email) {

					if err := rtr.S.DB.InsertUser(context.Background(), githubUser.Username, githubUser.Email); err != nil {

						http.Error(w, "error creating user", http.StatusInternalServerError)
						return

					}

					http.Redirect(w, r, "/api", http.StatusTemporaryRedirect)
					return

				}

				_, ok := rtr.S.CreditCache.Load(githubUser.Email)
				if !ok {

					user, err := rtr.S.DB.GetUser(context.Background(), githubUser.Email)
					if err != nil {

						http.Error(w, "internal error occurred", http.StatusInternalServerError)
						return

					}

					rtr.S.CreditCache.Store(githubUser.Email, user.CreditBalance)

				}

				// create a user session jwt
				tokenString, err := jwtauth.CreateJWTToken(githubUser.Email, rtr.S.JwtSession.SecretKey)
				if err != nil {
					http.Error(w, "internal error occurred", http.StatusInternalServerError)
					return

				}

				// http.SetCookie(w, &http.Cookie{

				// 	Name:     "jwtToken",
				// 	Value:    tokenString,
				// 	Expires:  time.Now().Add(time.Hour),
				// 	Path:     "/",
				// 	Secure:   false,
				// 	HttpOnly: true,
				// })

				if err := json.NewEncoder(w).Encode(&types.JWTCreatedResponse{Token: tokenString}); err != nil {

					http.Error(w, "internal error occurred while creating jwttoken", http.StatusInternalServerError)
					return

				}

				rtr.S.JwtSession.TokenPool.Store(tokenString, struct{}{})

			}
		}

	}

}

func HandleOAuthLogout(rtr *router.Router) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		http.SetCookie(w, &http.Cookie{Name: "jwtToken", Value: "loggedOut"})

	}

}

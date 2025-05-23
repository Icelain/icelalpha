package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"icealpha/internal/controllers/jwtauth"
	"icealpha/internal/controllers/oauth"
	"icealpha/internal/controllers/user"
	"icealpha/internal/database"
	"icealpha/internal/router"
	"log/slog"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

// Start all controllers and configure cookie store
func HandleAll(r *router.Router) {

	// set oauth2 config

	var err error
	if err = oauth.SetGithubOAuthConfig(); err != nil {

		r.Logger.Error(err.Error())
		return

	}
	if err = oauth.SetGoogleOAuthConfig(); err != nil {

		r.Logger.Error(err.Error())
		return

	}

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

		case "google":

			state := oauth.SetNewOAuthStateCookie(w)
			url := oauth.GoogleOAuthConfig.AuthCodeURL(state, oauth2.ApprovalForce)
			http.Redirect(w, r, url, http.StatusTemporaryRedirect)

		}

	}

}

// GET :: -> SessionCookie | Redirect
func HandleOAuthCallback(rtr *router.Router) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		provider := r.PathValue("provider")

		var user oauth.AuthUser

		switch provider {

		case "github":

			var githubUser oauth.GithubUser
			var err error
			githubUser, _, err = oauth.HandleGithubOAuthCallback(rtr, w, r)
			if err != nil {
				rtr.Logger.Error("err handling github oauth callback", "err", err)
				http.Redirect(w, r, "localhost:3000", http.StatusSeeOther)
				return
			}

			user = githubUser

		case "google":

			var googleUser oauth.GoogleUser
			var err error

			googleUser, _, err = oauth.HandleGoogleCallback(rtr, w, r)
			if err != nil {
				rtr.Logger.Error("err handling google oauth callback", "err", err)
				http.Redirect(w, r, "localhost:3000", http.StatusSeeOther)
				return
			}

			user = googleUser

		}

		var justCreated bool

		if !rtr.S.DB.CheckUserExists(context.Background(), user.GetEmail()) {

			if err := rtr.S.DB.InsertUser(context.Background(), user.GetUsername(), user.GetEmail()); err != nil {

				http.Error(w, "error creating user", http.StatusInternalServerError)
				return

			}

			justCreated = true

		}

		if justCreated {

			rtr.S.CreditCache.Store(user.GetEmail(), 5)

		} else {
			_, ok := rtr.S.CreditCache.Load(user.GetEmail())
			if !ok {

				dbUser, err := rtr.S.DB.GetUser(context.Background(), user.GetEmail())
				if err != nil {

					http.Error(w, "internal error occurred", http.StatusInternalServerError)
					return

				}

				rtr.S.CreditCache.Store(dbUser.Email, dbUser.CreditBalance)

			}
		}

		// create a user session jwt
		tokenString, err := jwtauth.CreateJWTToken(user.GetEmail(), rtr.S.JwtSession.SecretKey)
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

		http.Redirect(w, r, fmt.Sprintf("/dummyboard?jwtToken=%s", tokenString), http.StatusTemporaryRedirect)

		rtr.S.JwtSession.TokenPool.Store(tokenString, struct{}{})

	}

}

func HandleOAuthLogout(rtr *router.Router) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if token := r.Header.Get("jwttoken"); token != "" {

			rtr.S.JwtSession.TokenPool.Delete(token)

		}

		http.Redirect(w, r, "localhost:3000/", http.StatusTemporaryRedirect)
	}

}

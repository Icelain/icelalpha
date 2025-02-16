package controllers

import (
	"context"
	"encoding/json"
	"icealpha/internal/controllers/auth"
	"icealpha/internal/router"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

func HandleAll(r *router.Router) {

	// set oauth2 config
	auth.SetGithubOAuthConfig()

	// set cookiestore
	r.S.CookieStore = sessions.NewCookieStore([]byte("SESSION_KEY"))

	HandleAPIIndex("/api", r)
	HandleOAuthFlow("/api/oauth", r)
	HandleOAuthCallback("/api/oauth/{provider}/callback", r)

}

func HandleAPIIndex(pattern string, rtr *router.Router) {

	rtr.R.Get(pattern, func(w http.ResponseWriter, r *http.Request) {

		if err := json.NewEncoder(w).Encode(map[string]string{

			"icelalpha api status": "up",
		}); err != nil {

			http.Error(w, "error while serving /api", http.StatusInternalServerError)
			rtr.Logger.Error("Failed to serve /api", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})

		}

	})

}

func HandleOAuthFlow(pattern string, rtr *router.Router) {

	rtr.R.Get(pattern, func(w http.ResponseWriter, r *http.Request) {

		provider := r.URL.Query().Get("provider")

		switch provider {

		case "github":

			if auth.CheckSessionExists(r, rtr.S.CookieStore) {

				http.Redirect(w, r, "/api", http.StatusTemporaryRedirect)
				return

			}

			state := auth.SetNewOAuthStateCookie(w)
			url := auth.GithubOAuthConfig.AuthCodeURL(state, oauth2.ApprovalForce)
			http.Redirect(w, r, url, http.StatusTemporaryRedirect)

		}

	})

}

func HandleOAuthCallback(pattern string, rtr *router.Router) {

	rtr.R.Get(pattern, func(w http.ResponseWriter, r *http.Request) {

		provider := r.PathValue("provider")

		switch provider {

		case "github":

			var githubUser auth.GithubUser
			var redirectPath string
			var err error
			switch provider {
			case "github":
				githubUser, redirectPath, err = auth.HandleGithubOAuthCallback(rtr, auth.GithubOAuthConfig, w, r)
				if err != nil {
					rtr.Logger.Error("err handling github oauth callback", "err", err)
					http.Redirect(w, r, "/login", http.StatusSeeOther)
					return
				}

				if rtr.S.DB.CheckUserExists(context.Background(), githubUser.Email) {

					http.Redirect(w, r, redirectPath, http.StatusTemporaryRedirect)
					return

				}

				if err := rtr.S.DB.InsertUser(context.Background(), githubUser.Username, githubUser.Email); err != nil {

					http.Error(w, "error creating user", http.StatusInternalServerError)
					return

				}

				// create a user session cookie

				session, err := rtr.S.CookieStore.Get(r, "usersession")
				if err != nil {

					http.Error(w, "error creating session cookie", http.StatusInternalServerError)
					return

				}

				session.Options.MaxAge = int(time.Now().Add(time.Hour * 24).Unix())
				session.Values["username"] = githubUser.Username
				session.Values["email"] = githubUser.Email

				if err = session.Save(r, w); err != nil {

					http.Error(w, "error saving session cookie: "+err.Error(), http.StatusInternalServerError)
					return

				}

				http.Redirect(w, r, redirectPath, http.StatusTemporaryRedirect)
				return

			}

		}

	})

}

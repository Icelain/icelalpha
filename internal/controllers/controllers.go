package controllers

import (
	"encoding/json"
	"icealpha/internal/controllers/auth"
	"icealpha/internal/router"
	"log/slog"
	"net/http"
)

func HandleAll(r *router.Router) {

	// set oauth2 config
	auth.SetGithubOAuthConfig()

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

			state := auth.SetNewOAuthStateCookie(w)

			url := auth.GithubOAuthConfig.AuthCodeURL(state)
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
			}
			if err != nil {
				rtr.Logger.Error("err handling github oauth callback", "err", err)
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			// "user" is an instance of User that can be used to
			// create a new user or sign in an existing user

			// create a session cookie to keep the user signed in

			rtr.Logger.Info(githubUser.Email)

			http.Redirect(w, r, redirectPath, http.StatusSeeOther)

		}

	})

}

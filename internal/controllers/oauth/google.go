package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"icealpha/internal/router"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleUser struct {
	Email     string `json:"email"`
	FirstName string `json:"given_name"`
	LastName  string `json:"family_name"`
	AvatarURL string `json:"picture"`
}

var GoogleOAuthConfig *oauth2.Config

const oauthGoogleUserInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

func SetGoogleOAuthConfig() {
	GoogleOAuthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/oauth/google/callback",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

func HandleGoogleCallback(rtr *router.Router, w http.ResponseWriter, r *http.Request) (GoogleUser, string, error) {
	u := GoogleUser{}
	path := "/"
	oauthState, err := r.Cookie("oauthstate")

	if err != nil || r.FormValue("state") != oauthState.Value {
		rtr.Logger.Error("invalid oauth google state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return u, path, errors.New("invalid OAuth Google state")
	}

	user, err := getGoogleUserData(r.FormValue("code"))
	if err != nil {
		return u, path, fmt.Errorf("err getting user data from google: %+v", err)
	}

	u.FirstName = user.FirstName
	u.LastName = user.LastName
	u.Email = user.Email
	u.AvatarURL = user.AvatarURL

	return u, path, nil
}

func getGoogleUserData(code string) (GoogleUser, error) {
	gu := GoogleUser{}
	token, err := GoogleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		return gu, err
	}
	response, err := http.Get(oauthGoogleUserInfoURL + token.AccessToken)
	if err != nil {
		return gu, err
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&gu)
	return gu, err
}

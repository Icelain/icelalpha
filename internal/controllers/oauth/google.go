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

// Google user with default oauth2 permissions
type GoogleUser struct {
	Email     string `json:"email"`
	FirstName string `json:"given_name"`
	LastName  string `json:"family_name"`
	AvatarURL string `json:"picture"`
}

func (g GoogleUser) GetEmail() string {

	return g.Email
}

func (g GoogleUser) GetUsername() string {

	return g.FirstName + " " + g.LastName
}

func (g GoogleUser) GetAvatarURL() string {

	return g.AvatarURL

}

var GoogleOAuthConfig *oauth2.Config

const oauthGoogleUserInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

func SetGoogleOAuthConfig() error {

	clientId := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

	if clientId == "" || clientSecret == "" {

		return errors.New("Google ClientID or/and ClientSecret not found")

	}
	GoogleOAuthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/api/oauth/google/callback",
		ClientID:     clientId,
		ClientSecret: clientSecret,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return nil
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

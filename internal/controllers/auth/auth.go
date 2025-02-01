package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"icealpha/internal/router"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleUser struct {
	Email     string `json:"email"`
	FirstName string `json:"given_name"`
	LastName  string `json:"family_name"`
	AvatarURL string `json:"picture"`
	UserID    string `json:"id"`
	Name      string `json:"name"`
}

const (
	oauthGoogleUserInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
)

var (
	googleOAuthConfig *oauth2.Config
)

func SetNewGoogleOAuthConfig(srv *router.Router) {
	googleOAuthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:3100/auth/google/callback",
		ClientID:     srv.S.GOOGLE_CLIENT_KEY,
		ClientSecret: srv.S.GOOGLE_CLIENT_SECRET,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

func generateOAuthStateCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(1 * time.Hour)
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{
		Name:    "oauthstate",
		Value:   state,
		Expires: expiration,
		Path:    "/",
	}
	http.SetCookie(w, &cookie)

	return state
}

func OAuthFlow(w http.ResponseWriter, r *http.Request) {

	oauthState := generateOAuthStateCookie(w)
	url := googleOAuthConfig.AuthCodeURL(oauthState, oauth2.SetAuthURLParam("prompt", "select_account"))
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)

}

func HandleGoogleCallback(srv *router.Router, w http.ResponseWriter, r *http.Request) (GoogleUser, string, error) {
	u := GoogleUser{}
	path := "/"
	oauthState, err := r.Cookie("oauthstate")

	if err != nil || r.FormValue("state") != oauthState.Value {
		srv.Logger.Error("invalid oauth google state")
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
	u.UserID = user.UserID
	u.Name = user.Name

	return u, path, nil
}

func getGoogleUserData(code string) (GoogleUser, error) {
	gu := GoogleUser{}
	token, err := googleOAuthConfig.Exchange(context.Background(), code)
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

func HandleGoogleLogout(srv *router.Router, w http.ResponseWriter, r *http.Request) {

	http.SetCookie(w, &http.Cookie{

		Name:    "oauthstate",
		Value:   "expired",
		Expires: time.Now(),
		Path:    "/",
	})

}

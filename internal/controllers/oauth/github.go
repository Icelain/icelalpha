package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"icealpha/internal/router"
	"net/http"
	"os"
)

const (
	oauthGithubUserURL       = "https://api.github.com/user"
	oauthGithubUserEmailsURL = "https://api.github.com/user/emails"
)

var (
	GithubOAuthConfig *oauth2.Config
)

type GithubUser struct {
	Email     string `json:"email"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
}

func (g GithubUser) GetEmail() string {

	return g.Email
}

func (g GithubUser) GetUsername() string {

	return g.Username
}

func (g GithubUser) GetAvatarURL() string {

	return g.AvatarURL
}

func SetGithubOAuthConfig() error {
	clientId := os.Getenv("GITHUB_CLIENT_ID")
	clientSecret := os.Getenv("GITHUB_CLIENT_SECRET")

	if clientId == "" || clientSecret == "" {

		return errors.New("Github ClientID or/and ClientSecret not found")

	}

	GithubOAuthConfig = &oauth2.Config{

		ClientID:     clientId,
		ClientSecret: clientSecret,
		Endpoint:     github.Endpoint,
	}

	return nil

}

func HandleGithubOAuthCallback(router *router.Router, w http.ResponseWriter, r *http.Request) (GithubUser, string, error) {

	u := GithubUser{}
	path := "/"
	var err error
	oauthState, err := r.Cookie("oauthstate")

	if err != nil || r.FormValue("state") != oauthState.Value {
		router.Logger.Error("invalid oauth github state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return u, path, errors.New("invalid OAuth Github state")
	}

	code := r.URL.Query().Get("code")
	token, err := GithubOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		return u, path, err
	}

	githubUser, err := getGithubUserData(token.AccessToken)
	if err != nil {

		return u, path, err

	}

	return githubUser, path, nil

}

func getGithubUserData(accessToken string) (GithubUser, error) {
	gu := GithubUser{}

	req, err := http.NewRequest(http.MethodGet, oauthGithubUserURL, nil)
	if err != nil {
		return gu, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("token %s", accessToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return gu, err
	}

	if err := json.NewDecoder(res.Body).Decode(&gu); err != nil {
		return gu, err
	}
	return gu, nil
}

func getUserEmailFromGithub(accessToken string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, oauthGithubUserEmailsURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", fmt.Sprintf("token %s", accessToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	responseEmails := []struct {
		// contains other fields as well, but these
		// are the only one's we're interested in
		Email   string `json:"email"`
		Primary bool   `json:"primary"`
	}{}

	if err := json.NewDecoder(res.Body).Decode(&responseEmails); err != nil {
		return "", err
	}

	for _, re := range responseEmails {
		if re.Primary {
			return re.Email, nil
		}
	}

	return "", errors.New("no primary email found")
}

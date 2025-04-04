package oauth

import (
	"crypto/rand"
	"encoding/base64"

	"net/http"
	"time"
)

type AuthUser interface {
	GetEmail() string
	GetUsername() string
	GetAvatarURL() string
}

func SetNewOAuthStateCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(24 * time.Hour)
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{
		Name:    "oauthstate",
		Value:   state,
		Expires: expiration,
	}
	http.SetCookie(w, &cookie)

	return state
}

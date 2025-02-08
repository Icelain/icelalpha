package state

import "github.com/gorilla/sessions"

type State struct {
	CookieStore *sessions.CookieStore
}

package state

import (
	"icealpha/internal/database"

	"github.com/gorilla/sessions"
)

type State struct {
	CookieStore *sessions.CookieStore
	DB          *database.PostgresDriver
}

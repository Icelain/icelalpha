package state

import (
	"icealpha/internal/database"
	"icealpha/pkg/imglatex"

	"github.com/gorilla/sessions"
)

type State struct {
	CookieStore *sessions.CookieStore
	DB          *database.PostgresDriver
	ImgLatex    *imglatex.ImgLatex
}

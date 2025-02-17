package state

import (
	"icealpha/internal/database"
	"icealpha/pkg/imglatex"
	"icealpha/pkg/inference"

	"github.com/gorilla/sessions"
)

type State struct {
	CookieStore *sessions.CookieStore
	DB          *database.PostgresDriver
	ImgLatex    *imglatex.ImgLatex
	LLMClient   inference.LLMClient
}

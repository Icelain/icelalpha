package state

import (
	"icealpha/internal/database"
	"icealpha/pkg/imglatex"
	"icealpha/pkg/inference"
	"sync"

	"github.com/gorilla/sessions"
)

type State struct {
	CookieStore *sessions.CookieStore
	DB          *database.PostgresDriver
	CreditCache *sync.Map
	ImgLatex    *imglatex.ImgLatex
	LLMClient   inference.LLMClient
}

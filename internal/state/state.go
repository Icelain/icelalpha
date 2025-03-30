package state

import (
	"icealpha/internal/controllers/jwtauth"
	"icealpha/internal/database"
	"icealpha/pkg/imglatex"
	"icealpha/pkg/inference"
	"sync"
)

type State struct {
	JwtSession  *jwtauth.JWTSession
	DB          *database.PostgresDriver
	CreditCache *sync.Map
	ImgLatex    *imglatex.ImgLatex
	LLMClient   inference.LLMClient
}

package state

import (
	"icealpha/internal/database"
	"icealpha/pkg/imglatex"
	"icealpha/pkg/inference"
	"sync"
)

type State struct {
	DB          *database.PostgresDriver
	CreditCache *sync.Map
	ImgLatex    *imglatex.ImgLatex
	LLMClient   inference.LLMClient
}

package router

import (
	"fmt"
	"icealpha/internal/database"
	"icealpha/internal/state"
	"icealpha/pkg/imglatex"
	"icealpha/pkg/inference"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

// Router maintains state, logic and logging in one place
type Router struct {
	R      *chi.Mux
	Logger *slog.Logger
	S      *state.State
	Config *RouterConfig
}

func NewRouter() *Router {

	router := &Router{}

	mux := chi.NewMux()
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{

		Level: slog.LevelInfo,
	}))
	router.R = mux
	router.S = &state.State{}
	router.Logger = logger

	return router

}

// Configuration struct for Router
type RouterConfig struct {
	Port      uint
	DB        *database.PostgresDriver
	ImgLatex  *imglatex.ImgLatex
	LLMClient inference.LLMClient
}

func (r *Router) SetConfig(config *RouterConfig) {

	r.Config = config
	r.S.DB = config.DB
	r.S.ImgLatex = config.ImgLatex
	r.S.LLMClient = config.LLMClient

}

func (r *Router) Serve() error {

	r.Logger.Info(fmt.Sprintf("Server running on port :%d", r.Config.Port))
	return http.ListenAndServe(fmt.Sprintf(":%d", r.Config.Port), r.R)

}

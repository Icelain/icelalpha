package router

import (
	"fmt"
	"icealpha/internal/database"
	"icealpha/internal/state"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

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

type RouterConfig struct {
	Port uint
	DB   *database.PostgresDriver
}

func (r *Router) SetConfig(config *RouterConfig) {

	r.Config = config

}

func (r *Router) Serve() error {

	r.Logger.Info(fmt.Sprintf("Server running on port :%d", r.Config.Port))
	return http.ListenAndServe(fmt.Sprintf(":%d", r.Config.Port), r.R)

}

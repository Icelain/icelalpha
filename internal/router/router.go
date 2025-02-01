package router

import (
	"github.com/go-chi/chi/v5"
)

type Router struct {
	R *chi.Mux
	S *state.State
}

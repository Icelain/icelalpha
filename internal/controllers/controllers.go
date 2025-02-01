package controllers

import (
	"encoding/json"
	"icealpha/internal/router"
	"log/slog"
	"net/http"
)

func HandleAll(r *router.Router) {

	HandleAPIIndex("/api", r)

}

func HandleAPIIndex(pattern string, router *router.Router) {

	router.R.Get(pattern, func(w http.ResponseWriter, r *http.Request) {

		if err := json.NewEncoder(w).Encode(map[string]string{

			"icelalpha api status": "up",
		}); err != nil {

			http.Error(w, "error while serving /api", http.StatusInternalServerError)
			router.Logger.Error("Failed to serve /api", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})

		}

	})

}

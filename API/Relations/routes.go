package main

import (
	"net/http"
	authmiddleware "relations/auth"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	// Specify who is allowed to connect
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-TOKEN"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	//To check if service up or not
	mux.Use(middleware.Heartbeat("/ping"))

	// Add JWT middleware
	mux.Use(authmiddleware.JWTMiddleware)

	//Add route at root level
	mux.Post("/Follow", app.Follow)
	mux.Post("/UnFollow", app.UnFollow)

	return mux
}

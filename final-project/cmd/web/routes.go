package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *Config) routes() http.Handler {
	// mux router
	mux := chi.NewRouter()

	// set up middleware
	mux.Use(middleware.Recoverer)
	mux.Use(app.SessionLoad)

	// create application routes
	mux.Get("/", app.HomePage)

	// log-in/out routes
	mux.Get("/login", app.LoginPage)
	mux.Post("/login", app.PostLoginPage)
	mux.Get("/logout", app.Logout)

	// registration/activation routes
	mux.Get("/register", app.RegisterPage)
	mux.Post("/register", app.PostRegisterPage)
	mux.Get("/activate-account", app.ActivateAccount)

	return mux
}

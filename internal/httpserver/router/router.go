// Package router contains all routes for server. Based on chi router
package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"github.com/ole-larsen/green-api/internal/httpserver/handlers"
	"github.com/ole-larsen/green-api/internal/httpserver/middlewares"
)

type Mux struct {
	Router chi.Router
}

func NewMux() *Mux {
	return &Mux{
		Router: chi.NewRouter(),
	}
}

func (m *Mux) SetMiddlewares() *Mux {
	// A good base middleware stack
	m.Router.Use(middleware.RequestID)
	m.Router.Use(middleware.RealIP)
	m.Router.Use(middleware.Recoverer)
	m.Router.Use(middlewares.RSAMiddleware)
	m.Router.Use(middlewares.HashMiddleware)
	m.Router.Use(middlewares.GzipMiddleware)
	m.Router.Use(middlewares.LoggingMiddleware)

	return m
}

func (m *Mux) SetHandlers() *Mux {
	m.Router.Get("/", handlers.HTMLHandler(make(chan struct{})))

	m.Router.Get("/status", handlers.StatusHandler)

	m.Router.Mount("/debug", middleware.Profiler())
	m.Router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), // The url pointing to API definition
	))

	return m
}

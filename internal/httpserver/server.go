// Package httpserver is Webserver for server.
// Copyright 2024 The Oleg Nazarov. All rights reserved.
package httpserver

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ole-larsen/green-api/internal/httpserver/router"
)

type HTTPServer struct {
	Router *router.Mux
	Host   string
	Port   int
}

func NewHTTPServer() *HTTPServer {
	return &HTTPServer{}
}

func (s *HTTPServer) SetHost(h string) *HTTPServer {
	s.Host = h
	return s
}

func (s *HTTPServer) SetPort(p int) *HTTPServer {
	s.Port = p
	return s
}

func (s *HTTPServer) SetRouter(r *router.Mux) *HTTPServer {
	s.Router = r
	return s
}

func (s *HTTPServer) GetHost() string {
	return s.Host // it can be ""
}

func (s *HTTPServer) GetPort() int {
	return s.Port // it can't be 0
}

func (s *HTTPServer) GetRouter() *router.Mux {
	return s.Router
}

func (s *HTTPServer) ListenAndServe() error {
	const defaultTimeout = 3

	server := &http.Server{
		Addr:              s.Host + ":" + fmt.Sprintf("%d", s.Port),
		Handler:           s.Router.Router,
		ReadHeaderTimeout: defaultTimeout * time.Second,
	}

	return server.ListenAndServe()
}

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
	router *router.Mux
	host   string
	port   int
}

func NewHTTPServer() *HTTPServer {
	return &HTTPServer{}
}

func (s *HTTPServer) SetHost(h string) *HTTPServer {
	s.host = h
	return s
}

func (s *HTTPServer) SetPort(p int) *HTTPServer {
	s.port = p
	return s
}

func (s *HTTPServer) SetRouter(r *router.Mux) *HTTPServer {
	s.router = r
	return s
}

func (s *HTTPServer) GetHost() string {
	return s.host // it can be ""
}

func (s *HTTPServer) GetPort() int {
	return s.port // it can't be 0
}

func (s *HTTPServer) GetRouter() *router.Mux {
	return s.router
}

func (s *HTTPServer) ListenAndServe() error {
	const defaultTimeout = 3

	server := &http.Server{
		Addr:              s.host + ":" + fmt.Sprintf("%d", s.port),
		Handler:           s.router.Router,
		ReadHeaderTimeout: defaultTimeout * time.Second,
	}

	return server.ListenAndServe()
}

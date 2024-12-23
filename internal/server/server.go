package server

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/ole-larsen/green-api/internal/httpserver"
	"github.com/ole-larsen/green-api/internal/httpserver/router"
	"github.com/ole-larsen/green-api/internal/log"
	"github.com/ole-larsen/green-api/internal/server/config"
)

var (
	logger = log.NewLogger("info", log.DefaultBuildLogger)
)

// Server represents the server instance, encapsulating settings,
// logger, signal handling, and storage and gRPC server components.
type Server struct {
	http     *httpserver.HTTPServer
	settings *config.Config
	logger   *log.Logger
	signal   chan os.Signal
	done     chan struct{}
}

// NewServer creates and returns a new Server instance with default logger settings.
func NewServer() *Server {
	return &Server{
		logger: logger,
	}
}

var SetupFunc = Setup

// Setup initializes the server with provided configuration settings and sets up
// storage. Returns an error if initialization fails.
func Setup(_ context.Context, settings *config.Config) (*Server, error) {
	s := NewServer()

	if err := s.Init(settings, make(chan os.Signal, 1), make(chan struct{})); err != nil {
		return nil, err
	}

	return s, nil
}

// Run starts the server and begins listening for shutdown signals. It runs the gRPC server
// and handles shutdown on receiving system interrupt signals like SIGINT or SIGTERM.
func (s *Server) Run(ctx context.Context, cancel context.CancelFunc) {
	defer cancel()

	defer close(s.signal)

	// shutdown workers
	go func() {
		signal.Notify(s.signal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	}()

	go func(signal chan os.Signal, done chan struct{}) {
		<-signal
		close(done)
		s.logger.Infow("...graceful server shutdown")
	}(s.signal, s.done)

	port := s.settings.Port
	host := s.settings.Host

	s.logger.Infow("...starting server",
		"host", host,
		"port", port,
		"goroutines", runtime.NumGoroutine(),
	)

	go func() {
		if err := s.http.ListenAndServe(); err != nil {
			s.logger.Errorln(err)
		}
	}()

	for {
		select {
		case <-s.done:
			s.logger.Infow("...stop server",
				"goroutines", runtime.NumGoroutine(),
			)

			return
		case <-ctx.Done():
			s.logger.Infow("stop server by ctx")
			return
		}
	}
}

// Init initializes the server with the given settings, signal channels,
// storage interface, and gRPC interface. Returns an error if any component is missing.
func (s *Server) Init(
	settings *config.Config,
	sgnl chan os.Signal,
	done chan struct{},
) error {
	s.SetSettings(settings).
		SetSignal(sgnl).
		SetDone(done)

	if s.settings == nil {
		return NewError(errors.New("config is missing"))
	}

	if s.signal == nil {
		return NewError(errors.New("signal is missing"))
	}

	if s.done == nil {
		return NewError(errors.New("done is missing"))
	}

	r := router.NewMux().
		SetMiddlewares().
		SetHandlers()

	s.
		setHTTPServer(httpserver.NewHTTPServer().
			SetHost(s.settings.Host).
			SetPort(s.settings.Port).
			SetRouter(r))

	if s.http == nil {
		return NewError(errors.New("http server is missing"))
	}

	if s.http != nil {
		if s.http.GetPort() == 0 {
			return NewError(errors.New("http server port is missing"))
		}

		if s.http.GetRouter() == nil {
			return NewError(errors.New("http server router is missing"))
		}
	}

	return nil
}

// SetSettings sets the server configuration.
func (s *Server) SetSettings(settings *config.Config) *Server {
	s.settings = settings
	return s
}

// SetSignal sets the signal channel for handling OS signals.
func (s *Server) SetSignal(sgnl chan os.Signal) *Server {
	s.signal = sgnl
	return s
}

// SetDone sets the done channel to signal when the server should stop.
func (s *Server) SetDone(done chan struct{}) *Server {
	s.done = done
	return s
}

func (s *Server) setHTTPServer(hs *httpserver.HTTPServer) *Server {
	s.http = hs
	return s
}

// GetSettings retrieves the server configuration settings.
func (s *Server) GetSettings() *config.Config {
	return s.settings
}

// GetSignal retrieves the signal channel used by the server.
func (s *Server) GetSignal() chan os.Signal {
	return s.signal
}

// GetDone retrieves the done channel used by the server.
func (s *Server) GetDone() chan struct{} {
	return s.done
}

// GetLogger retrieves the logger used by the server.
func (s *Server) GetLogger() *log.Logger {
	return s.logger
}

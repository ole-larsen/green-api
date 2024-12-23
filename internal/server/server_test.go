package server_test

import (
	"context"
	"errors"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/ole-larsen/green-api/internal/server"
	"github.com/ole-larsen/green-api/internal/server/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap/zaptest"
)

// MockConfig for simulating different configuration scenarios.
type MockConfig struct {
	*config.Config
}

func TestNewServer(t *testing.T) {
	srv := server.NewServer()
	assert.NotNil(t, srv, "server instance should not be nil")
	assert.NotNil(t, srv.GetLogger(), "server logger should not be nil")
	assert.Nil(t, srv.GetSettings(), "server settings should be nil initially")
	assert.Nil(t, srv.GetSignal(), "server signal channel should be nil initially")
	assert.Nil(t, srv.GetDone(), "server done channel should be nil initially")
}

func TestServer_Init(t *testing.T) {
	tests := []struct {
		expectedError  error
		settings       *config.Config
		signal         chan os.Signal
		done           chan struct{}
		name           string
		verifyChannels bool
	}{
		{
			name: "successful initialization",
			settings: &config.Config{
				Host: "localhost",
				Port: 8080,
			},
			signal:         make(chan os.Signal, 1),
			done:           make(chan struct{}),
			expectedError:  nil,
			verifyChannels: true,
		},
		{
			name:           "missing configuration",
			settings:       nil,
			signal:         make(chan os.Signal, 1),
			done:           make(chan struct{}),
			expectedError:  server.NewError(errors.New("config is missing")),
			verifyChannels: true,
		},
		{
			name: "missing signal channel",
			settings: &config.Config{
				Host: "localhost",
				Port: 8080,
			},
			signal:         nil,
			done:           make(chan struct{}),
			expectedError:  server.NewError(errors.New("signal is missing")),
			verifyChannels: true,
		},
		{
			name: "missing done channel",
			settings: &config.Config{
				Host: "localhost",
				Port: 8080,
			},
			signal:         make(chan os.Signal, 1),
			done:           nil,
			expectedError:  server.NewError(errors.New("done is missing")),
			verifyChannels: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			srv := server.NewServer()
			err := srv.Init(tt.settings, tt.signal, tt.done)

			// Check if the error message is as expected
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			// Verify if channels are properly set
			if tt.verifyChannels {
				if tt.settings != nil {
					assert.NotNil(t, srv.GetSettings(), "settings should not be nil")
				} else {
					assert.Nil(t, srv.GetSettings(), "settings should be nil")
				}

				if tt.signal != nil {
					assert.NotNil(t, srv.GetSignal(), "signal channel should not be nil")
				} else {
					assert.Nil(t, srv.GetSignal(), "signal channel should be nil")
				}

				if tt.done != nil {
					assert.NotNil(t, srv.GetDone(), "done channel should not be nil")
				} else {
					assert.Nil(t, srv.GetDone(), "done channel should be nil")
				}
			}
		})
	}
}

// TestServerSetSettings tests the SetSettings method.
func TestServerSetSettings(t *testing.T) {
	srv := server.NewServer()
	settings := &config.Config{}
	srv.SetSettings(settings)

	assert.Equal(t, settings, srv.GetSettings())
}

// TestServerSetSignal tests the SetSignal method.
func TestServerSetSignal(t *testing.T) {
	srv := server.NewServer()
	signal := make(chan os.Signal)
	srv.SetSignal(signal)

	assert.Equal(t, signal, srv.GetSignal())
}

// TestServerSetDone tests the SetDone method.
func TestServerSetDone(t *testing.T) {
	srv := server.NewServer()
	done := make(chan struct{})
	srv.SetDone(done)

	assert.Equal(t, done, srv.GetDone())
}

// TestServerInitErrorMissingSettings tests the error case when settings are missing.
func TestServerInitErrorMissingSettings(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	srv := server.NewServer()
	signal := make(chan os.Signal)
	done := make(chan struct{})

	err := srv.Init(nil, signal, done)
	assert.Error(t, err)
	assert.Equal(t, "[server]: config is missing", err.Error())
}

// TestServerInitErrorMissingSignal tests the error case when signal is missing.
func TestServerInitErrorMissingSignal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	srv := server.NewServer()
	settings := &config.Config{}
	done := make(chan struct{})

	err := srv.Init(settings, nil, done)
	assert.Error(t, err)
	assert.Equal(t, "[server]: signal is missing", err.Error())
}

// TestServerInitErrorMissingDone tests the error case when done is missing.
func TestServerInitErrorMissingDone(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	srv := server.NewServer()
	settings := &config.Config{}
	signal := make(chan os.Signal)

	err := srv.Init(settings, signal, nil)
	assert.Error(t, err)
	assert.Equal(t, "[server]: done is missing", err.Error())
}

func TestServer_setSettings(t *testing.T) {
	srv := server.NewServer()
	settings := &config.Config{
		Host: "localhost",
		Port: 8080,
	}

	// Test with valid settings
	srv.SetSettings(settings)
	assert.Equal(t, settings, srv.GetSettings(), "server settings should match the input settings")

	// Test with nil settings
	srv.SetSettings(nil)
	assert.Nil(t, srv.GetSettings(), "server settings should be nil when nil input is provided")
}

func TestServer_setSignal(t *testing.T) {
	srv := server.NewServer()
	signalChan := make(chan os.Signal, 1)

	// Test with valid signal channel
	srv.SetSignal(signalChan)
	assert.Equal(t, signalChan, srv.GetSignal(), "server signal channel should match the input channel")

	// Test with nil signal channel
	srv.SetSignal(nil)
	assert.Nil(t, srv.GetSignal(), "server signal channel should be nil when nil input is provided")
}

func TestServer_setDone(t *testing.T) {
	srv := server.NewServer()
	doneChan := make(chan struct{})

	// Test with valid done channel
	srv.SetDone(doneChan)
	assert.Equal(t, doneChan, srv.GetDone(), "server done channel should match the input channel")

	// Test with nil done channel
	srv.SetDone(nil)
	assert.Nil(t, srv.GetDone(), "server done channel should be nil when nil input is provided")
}

func TestServer_Setup_Success(t *testing.T) {
	ctx := context.Background()
	settings := &config.Config{
		Host: "localhost",
		Port: 8080,
	}

	srv, err := server.Setup(ctx, settings)
	require.NoError(t, err)
	require.NotNil(t, srv)

	assert.Equal(t, settings, srv.GetSettings())
	assert.NotNil(t, srv.GetSignal())
	assert.NotNil(t, srv.GetDone())
}

func TestServer_Run(t *testing.T) {

	// Create a new server instance
	srv := server.NewServer()

	// Prepare the context and cancel function
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Prepare the server settings
	settings := &config.Config{
		Host: "localhost",
		Port: 8080,
	}

	// Initialize the server with valid settings
	err := srv.Init(settings, make(chan os.Signal, 1), make(chan struct{}))
	assert.NoError(t, err, "server initialization should succeed")

	// Run the server in a separate goroutine
	go srv.Run(ctx, cancel)

	// Simulate a server stop by sending an interrupt signal
	srv.GetSignal() <- syscall.SIGINT

	// Give the server some time to process the signal and shut down
	time.Sleep(1 * time.Second)

	// Ensure the server has gracefully shut down
	select {
	case <-srv.GetDone():
		assert.True(t, true, "server should have shut down gracefully on signal")
	default:
		assert.Fail(t, "server did not shut down as expected")
	}

	// Ensure the server stops on context cancellation
	cancel()

	// Give the server some time to handle the context cancellation
	time.Sleep(1 * time.Second)

	select {
	case <-ctx.Done():
		assert.True(t, true, "server should have stopped on context cancellation")
	default:
		assert.Fail(t, "server did not stop on context cancellation")
	}
}

func TestServer_Run_ContextDone(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a logger
	logger := zaptest.NewLogger(t)
	defer func() {
		err := logger.Sync()
		if err != nil {
			t.Fatal("should not fall")
		}
	}()

	// Create a new server instance
	srv := server.NewServer()

	// Prepare the context and cancel function
	ctx, cancel := context.WithCancel(context.Background())

	// Prepare the server settings
	settings := &config.Config{
		Host: "localhost",
		Port: 8080,
	}

	// Initialize the server with valid settings
	err := srv.Init(settings, make(chan os.Signal, 1), make(chan struct{}))
	assert.NoError(t, err, "server initialization should succeed")

	// Run the server in a separate goroutine
	go srv.Run(ctx, cancel)

	// Cancel the context to simulate shutdown by context cancellation
	cancel()

	// Give the server some time to handle the context cancellation
	time.Sleep(1 * time.Second)

	// Ensure the server has gracefully shut down
	select {
	case <-ctx.Done():
		assert.True(t, true, "server should have stopped on context cancellation")
	default:
		assert.Fail(t, "server did not stop on context cancellation")
	}
}

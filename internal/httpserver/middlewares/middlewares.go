// Package middlewares contains all middleware logic for server
package middlewares

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/ole-larsen/green-api/internal/log"
)

type Middleware func(http.Handler) http.Handler

func Conveyor(h http.Handler, ms ...Middleware) http.Handler {
	if h == nil {
		return nil
	}

	for _, middleware := range ms {
		h = middleware(h)
	}

	return h
}

func LoggingMiddleware(next http.Handler) http.Handler {
	logFn := func(rw http.ResponseWriter, r *http.Request) {
		logger := log.NewLogger("info", log.DefaultBuildLogger)
		body, err := io.ReadAll(r.Body)

		defer func() {
			e := r.Body.Close()
			if e != nil {
				logger.Errorln(e)
				return
			}
		}()

		start := time.Now()

		lw := LoggingResponseWriter{
			ResponseWriter: rw,
			ResponseData: &ResponseData{
				RequestURI: r.RequestURI,
				Body:       body,
				Status:     0,
				Size:       0,
			},
		}

		if err == nil {
			r.Body = io.NopCloser(bytes.NewReader(body))
		}

		next.ServeHTTP(&lw, r)

		err = dumpRequest(logger, &lw, r, start)
		if err != nil {
			logger.Errorln(err)
		}
	}

	return http.HandlerFunc(logFn)
}

func LengthContentMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// переменная reader будет равна r.Body или *gzip.Reader
		var reader io.Reader

		if r.Header.Get(`Content-Encoding`) == `gzip` {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}

			reader = gz

			defer func() {
				err := gz.Close()
				if err != nil {
					http.Error(rw, err.Error(), http.StatusInternalServerError)
					return
				}
			}()
		} else {
			reader = r.Body
		}

		body, err := io.ReadAll(reader)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		r.Body = io.NopCloser(bytes.NewReader(body))

		next.ServeHTTP(rw, r)
	})
}

func HashMiddleware(h http.Handler) http.Handler {
	logFn := func(rw http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.ServeHTTP(rw, r)
			return
		}

		encrypted := r.Header.Get("HashSHA256")
		if encrypted == "" {
			h.ServeHTTP(rw, r)
			return
		}

		var body = make([]byte, 0)

		// modify request body
		r.Body = io.NopCloser(bytes.NewReader(body))

		// serve next
		h.ServeHTTP(rw, r)
	}

	return http.HandlerFunc(logFn)
}

func RSAMiddleware(h http.Handler) http.Handler {
	logFn := func(rw http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.ServeHTTP(rw, r)
			return
		}

		var (
			body = make([]byte, 0)
		)

		// modify request body
		r.Body = io.NopCloser(bytes.NewReader(body))

		// serve next
		h.ServeHTTP(rw, r)
	}

	return http.HandlerFunc(logFn)
}

func dumpRequest(logger *log.Logger, lrw *LoggingResponseWriter, r *http.Request, start time.Time) error {
	data, err := httputil.DumpRequest(r, true)

	if err != nil {
		return err
	}

	dump := strings.Split(string(data), "\r\n")

	if len(dump) > 0 {
		logParams := make(map[string]string)

		const reqLen = 2

		for i := 1; i < len(dump); i++ {
			params := strings.Split(dump[i], ": ")

			if len(params) == reqLen {
				logParams[params[0]] = params[1]
			}
		}

		duration := time.Since(start)

		msg := dump[0]

		if lrw.ResponseData.Status == http.StatusInternalServerError {
			logger.Infow(msg,
				"url", lrw.ResponseData.RequestURI,
				"host", logParams["Host"],
				"content-type", logParams["Content-Type"],
				"accept-encoding", logParams["Accept-Encoding"],
				"content-encoding", logParams["Content-Encoding"],
				"content-length", logParams["Content-Length"],
				"user-agent", logParams["User-Agent"],
				"duration", duration,
				"status", lrw.ResponseData.Status,
				"size", lrw.ResponseData.Size,
				"body", string(lrw.ResponseData.Body),
			)
		}

		return nil
	}

	return fmt.Errorf("wrong request")
}

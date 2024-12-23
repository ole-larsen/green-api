package middlewares

import (
	"io"
	"net/http"
	"strings"

	"github.com/ole-larsen/green-api/internal/compressor"
)

func GzipMiddleware(next http.Handler) http.Handler {
	logFn := func(rw http.ResponseWriter, r *http.Request) {
		ow := rw

		acceptGzip := strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")

		if acceptGzip {
			rw.Header().Set("Content-Encoding", "gzip")
			rw.Header().Set("Accept-Encoding", "")

			rw.Header().Set("Content-Type", r.Header.Get("Content-Type"))

			cw := compressor.NewCompressWriter(rw)
			ow = cw

			defer func() {
				if err := cw.Close(); err != nil {
					return
				}
			}()
		}

		contentGzip := strings.Contains(r.Header.Get("Content-Encoding"), "gzip")

		if contentGzip {
			rw.Header().Set("Content-Encoding", "")
			rw.Header().Set("Accept-Encoding", "gzip")

			cr, err := compressor.NewCompressReader(r.Body)
			if err != nil {
				rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
				rw.WriteHeader(http.StatusInternalServerError)

				return
			}

			defer func() {
				if err := cr.Close(); err != nil {
					return
				}
			}()

			r.Body = io.NopCloser(cr)
		}

		next.ServeHTTP(ow, r)
	}

	return http.HandlerFunc(logFn)
}

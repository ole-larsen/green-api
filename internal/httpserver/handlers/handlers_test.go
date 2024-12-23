package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatusHandler(t *testing.T) {
	type want struct {
		response    string
		contentType string
		code        int
	}

	tests := []struct {
		name string
		want want
	}{
		{
			name: "positive test #1",
			want: want{
				code:        200,
				response:    `{"status":"ok"}`,
				contentType: "application/json",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/status", http.NoBody)
			w := httptest.NewRecorder()
			StatusHandler(w, request)

			resp := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, resp.StatusCode)
			// получаем и проверяем тело запроса
			defer func() {
				err := resp.Body.Close()
				require.NoError(t, err)
			}()

			resBody, err := io.ReadAll(resp.Body)

			require.NoError(t, err)
			assert.JSONEq(t, test.want.response, string(resBody))
			assert.Equal(t, test.want.contentType, resp.Header.Get("Content-Type"))
		})
	}
}

func ExampleStatusHandler() {
	request := httptest.NewRequest(http.MethodGet, "/status", http.NoBody)
	w := httptest.NewRecorder()
	StatusHandler(w, request)

	resp := w.Result()

	err := resp.Body.Close()
	if err != nil {
		panic(err)
	}
}

func TestBadRequestHandler(t *testing.T) {
	type want struct {
		response    string
		contentType string
		code        int
	}

	tests := []struct {
		name string
		want want
	}{
		{
			name: "positive test #1",
			want: want{
				code:        http.StatusBadRequest,
				response:    "400 bad request\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/status", http.NoBody)
			w := httptest.NewRecorder()
			BadRequest(w, request)

			resp := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, resp.StatusCode)

			defer func() {
				err := resp.Body.Close()
				require.NoError(t, err)
			}()

			resBody, err := io.ReadAll(resp.Body)

			require.NoError(t, err)
			assert.Equal(t, test.want.response, string(resBody))
			assert.Equal(t, test.want.contentType, resp.Header.Get("Content-Type"))
		})
	}
}

func TestNotFoundHandler(t *testing.T) {
	type want struct {
		response    string
		contentType string
		code        int
	}

	tests := []struct {
		name string
		want want
	}{
		{
			name: "positive test #1",
			want: want{
				code:        http.StatusNotFound,
				response:    "404 page not found\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/status", http.NoBody)
			w := httptest.NewRecorder()
			NotFoundRequest(w, request)

			resp := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, resp.StatusCode)

			defer func() {
				err := resp.Body.Close()
				require.NoError(t, err)
			}()

			resBody, err := io.ReadAll(resp.Body)

			require.NoError(t, err)
			assert.Equal(t, test.want.response, string(resBody))
			assert.Equal(t, test.want.contentType, resp.Header.Get("Content-Type"))
		})
	}
}

func TestNotAllowedRequestHandler(t *testing.T) {
	type want struct {
		response    string
		contentType string
		code        int
	}

	tests := []struct {
		name string
		want want
	}{
		{
			name: "positive test #1",
			want: want{
				code:        http.StatusMethodNotAllowed,
				response:    "405 method not allowed\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/status", http.NoBody)
			w := httptest.NewRecorder()
			NotAllowedRequest(w, request)

			resp := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, resp.StatusCode)

			defer func() {
				err := resp.Body.Close()
				require.NoError(t, err)
			}()

			resBody, err := io.ReadAll(resp.Body)

			require.NoError(t, err)
			assert.Equal(t, test.want.response, string(resBody))
			assert.Equal(t, test.want.contentType, resp.Header.Get("Content-Type"))
		})
	}
}

func TestInternalServerErrorHandler(t *testing.T) {
	type want struct {
		response    string
		contentType string
		code        int
	}

	tests := []struct {
		name string
		want want
	}{
		{
			name: "positive test #1",
			want: want{
				code:        http.StatusInternalServerError,
				response:    "500 internal server error\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/status", http.NoBody)
			w := httptest.NewRecorder()
			InternalServerErrorRequest(w, request)

			resp := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, resp.StatusCode)

			defer func() {
				err := resp.Body.Close()
				require.NoError(t, err)
			}()

			resBody, err := io.ReadAll(resp.Body)

			require.NoError(t, err)
			assert.Equal(t, test.want.response, string(resBody))
			assert.Equal(t, test.want.contentType, resp.Header.Get("Content-Type"))
		})
	}
}

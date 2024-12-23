package router_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/ole-larsen/green-api/internal/httpserver/router"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (resp *http.Response, res string) {
	t.Helper()

	req, err := http.NewRequest(method, ts.URL+path, body)

	require.NoError(t, err)

	resp, err = ts.Client().Do(req)

	require.NoError(t, err)

	defer func() {
		e := resp.Body.Close()
		require.NoError(t, e)
	}()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestNewRouter(t *testing.T) {
	tests := []struct {
		want *router.Mux
		name string
	}{
		{
			name: "create router",
			want: &router.Mux{
				Router: chi.NewRouter(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := router.NewMux()
			require.NotNil(t, r.Router)
		})
	}
}

func TestRouter(t *testing.T) {
	ts := httptest.NewServer(router.NewMux().
		SetMiddlewares().
		SetHandlers().Router)
	defer ts.Close()

	var testTable = []struct {
		url    string
		method string
		want   string
		body   []byte
		status int
	}{
		// {"/", "ok", http.StatusOK},
		{"/status", http.MethodGet, `{"status":"ok"}`, nil, http.StatusOK},
	}

	for _, v := range testTable {
		body := bytes.NewReader(v.body)
		resp, get := testRequest(t, ts, v.method, v.url, body)

		func() {
			defer func() {
				err := resp.Body.Close()
				require.NoError(t, err)
			}()
		}()

		require.Equal(t, v.status, resp.StatusCode)

		if !strings.Contains(v.url, "updates") {
			assert.Equal(t, v.want, get)
		} else if resp.StatusCode == http.StatusOK {
			assert.JSONEq(t, v.want, get)
		}
	}
}

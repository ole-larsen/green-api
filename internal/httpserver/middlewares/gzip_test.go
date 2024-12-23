package middlewares

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/ole-larsen/green-api/internal/log"
)

var logger = log.NewLogger("info", log.DefaultBuildLogger)

func webhookPlain(rw http.ResponseWriter, _ *http.Request) {
	rw.WriteHeader(http.StatusOK)

	status, err := rw.Write([]byte("OK"))
	if err != nil {
		logger.Errorln(status, err)
	}
}

func webhookPlainJSON(rw http.ResponseWriter, _ *http.Request) {
	rw.WriteHeader(http.StatusOK)

	status, err := rw.Write([]byte(`{"status":"ok"}`))
	if err != nil {
		logger.Errorln(status, err)
	}
}

func webhookJSON(rw http.ResponseWriter, _ *http.Request) {
	// content-type must be set
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)

	status, err := rw.Write([]byte(`{"status":"ok"}`))
	if err != nil {
		logger.Errorln(status, err)
	}
}

type Request struct {
	Text     string `json:"text"`
	Timezone string `json:"timezone"`
}

type Response struct {
	Response ResponsePayload `json:"response"`
}

type ResponsePayload struct {
	Text string `json:"text"`
}

func webhookJSONWithBody(rw http.ResponseWriter, r *http.Request) {
	var body Request

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
		rw.WriteHeader(http.StatusInternalServerError)

		return
	}

	tz, err := time.LoadLocation(body.Timezone)
	if err != nil {
		rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
		rw.WriteHeader(http.StatusBadRequest)

		return
	}

	now := time.Now().In(tz)
	hour, minute, _ := now.Clock()

	text := fmt.Sprintf("Точное время %d часов, %d минут. %s", hour, minute, body.Text)

	resp, err := json.Marshal(Response{
		Response: ResponsePayload{
			Text: text,
		},
	})
	if err != nil {
		rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
		rw.WriteHeader(http.StatusInternalServerError)

		return
	}

	// content-type must be set
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	status, err := rw.Write(resp)

	if err != nil {
		logger.Errorln(status, err)
	}
}

func TestGzipMiddleware_SendNotGzip(t *testing.T) {
	var err error

	type args struct {
		handler         http.Handler
		body            io.Reader
		method          string
		contentType     string
		acceptEncoding  string
		contentEncoding string
	}

	type want struct {
		response        string
		contentType     string
		acceptEncoding  string
		contentEncoding string
		status          int
	}

	tz, err := time.LoadLocation("")
	require.NoError(t, err)

	now := time.Now().In(tz)
	hour, minute, _ := now.Clock()

	response := `{"response":{"text":"Точное время ` + fmt.Sprintf("%d", hour) + ` часов, ` + fmt.Sprintf("%d", minute) + ` минут. this is a very long string from client"}}`
	requestBody := `{
	        "text": "this is a very long string from client"
	    }`

	buf := bytes.NewBuffer(nil)
	zb := gzip.NewWriter(buf)
	_, err = zb.Write([]byte(requestBody))
	require.NoError(t, err)
	err = zb.Close()
	require.NoError(t, err)

	const JSONContentType = "application/json"

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "no compressed [GET]",
			args: args{
				method:          http.MethodGet,
				acceptEncoding:  "",
				contentEncoding: "",
				contentType:     "text/plain; charset=utf-8",
				handler:         GzipMiddleware(http.HandlerFunc(webhookPlain)),
				body:            nil,
			},
			want: want{
				response:        "OK",
				status:          http.StatusOK,
				acceptEncoding:  "",
				contentEncoding: "",
				contentType:     "text/plain; charset=utf-8",
			},
		},
		{
			name: "no compressed [POST]",
			args: args{
				method:          http.MethodPost,
				acceptEncoding:  "",
				contentEncoding: "",
				contentType:     "text/plain; charset=utf-8",
				handler:         GzipMiddleware(http.HandlerFunc(webhookPlain)),
				body:            nil,
			},
			want: want{
				response:        "OK",
				status:          http.StatusOK,
				acceptEncoding:  "",
				contentEncoding: "",
				contentType:     "text/plain; charset=utf-8",
			},
		},
		{
			name: "no compressed json [GET]",
			args: args{
				method:          http.MethodGet,
				acceptEncoding:  "",
				contentEncoding: "",
				contentType:     "text/plain; charset=utf-8",
				handler:         GzipMiddleware(http.HandlerFunc(webhookPlainJSON)),
				body:            nil,
			},
			want: want{
				response:        `{"status":"ok"}`,
				status:          http.StatusOK,
				acceptEncoding:  "",
				contentEncoding: "",
				contentType:     "text/plain; charset=utf-8",
			},
		},
		{
			name: "no compressed [POST]",
			args: args{
				method:          http.MethodPost,
				acceptEncoding:  "",
				contentEncoding: "",
				contentType:     "text/plain; charset=utf-8",
				handler:         GzipMiddleware(http.HandlerFunc(webhookPlainJSON)),
				body:            nil,
			},
			want: want{
				response:        `{"status":"ok"}`,
				status:          http.StatusOK,
				acceptEncoding:  "",
				contentEncoding: "",
				contentType:     "text/plain; charset=utf-8",
			},
		},
		{
			name: "no compressed no body [GET]",
			args: args{
				method:          http.MethodGet,
				acceptEncoding:  "",
				contentEncoding: "",
				contentType:     JSONContentType,
				handler:         GzipMiddleware(http.HandlerFunc(webhookJSON)),
				body:            nil,
			},
			want: want{
				response:        `{"status":"ok"}`,
				status:          http.StatusOK,
				acceptEncoding:  "",
				contentEncoding: "",
				contentType:     JSONContentType,
			},
		},
		{
			name: "no comressed no body [POST]",
			args: args{
				method:          http.MethodPost,
				acceptEncoding:  "",
				contentEncoding: "",
				contentType:     JSONContentType,
				handler:         GzipMiddleware(http.HandlerFunc(webhookJSON)),
				body:            nil,
			},
			want: want{
				response:        `{"status":"ok"}`,
				status:          http.StatusOK,
				acceptEncoding:  "",
				contentEncoding: "",
				contentType:     JSONContentType,
			},
		},
		{
			name: "no compressed with body [POST]: positive",
			args: args{
				method:          http.MethodPost,
				acceptEncoding:  "",
				contentEncoding: "",
				contentType:     JSONContentType,
				handler:         GzipMiddleware(http.HandlerFunc(webhookJSONWithBody)),
				body:            bytes.NewBuffer([]byte(requestBody)),
			},
			want: want{
				response:        response,
				status:          http.StatusOK,
				acceptEncoding:  "",
				contentEncoding: "",
				contentType:     JSONContentType,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(tt.args.handler)
			defer srv.Close()

			r := httptest.NewRequest(tt.args.method, srv.URL, tt.args.body)
			r.RequestURI = ""

			require.Equal(t, 0, len(r.Header))
			require.Equal(t, 0, len(r.Header.Values("Accept-Encoding")))
			require.Equal(t, 0, len(r.Header.Values("Content-Encoding")))
			require.Equal(t, 0, len(r.Header.Values("Content-Type")))

			r.Header.Set("Accept-Encoding", tt.args.acceptEncoding)
			r.Header.Set("Content-Encoding", tt.args.contentEncoding)
			r.Header.Set("Content-Type", tt.args.contentType)

			resp, err := http.DefaultClient.Do(r)
			require.NoError(t, err)

			require.Equal(t, tt.want.status, resp.StatusCode)

			require.Equal(t, tt.want.acceptEncoding, resp.Header.Get("Accept-Encoding"))
			require.Equal(t, tt.want.contentEncoding, resp.Header.Get("Content-Encoding"))
			require.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))

			t.Log(resp.Header.Get("Accept-Encoding"))
			t.Log(resp.Header.Get("Content-Encoding"))
			t.Log(resp.Header.Get("Content-Type"))

			defer func() {
				e := resp.Body.Close()
				require.NoError(t, e)
			}()

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			t.Log(string(body))
			require.Equal(t, tt.want.response, string(body))

			if tt.want.contentType == JSONContentType {
				require.JSONEq(t, tt.want.response, string(body))
			}
		})
	}
}

func TestGzipMiddleware_AcceptNotGzip(t *testing.T) {
	var err error

	type args struct {
		handler         http.Handler
		body            io.Reader
		method          string
		contentType     string
		acceptEncoding  string
		contentEncoding string
	}

	type want struct {
		response        string
		contentType     string
		acceptEncoding  string
		contentEncoding string
		status          int
	}

	tz, err := time.LoadLocation("")
	require.NoError(t, err)

	now := time.Now().In(tz)
	hour, minute, _ := now.Clock()

	response := `{"response":{"text":"Точное время ` + fmt.Sprintf("%d", hour) + ` часов, ` + fmt.Sprintf("%d", minute) + ` минут. this is a very long string from client"}}`

	requestBody := `{
	    "text": "this is a very long string from client"
	}`

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Accept-Encoding: gzip [GET]",
			args: args{
				method:          http.MethodGet,
				acceptEncoding:  "gzip",
				contentEncoding: "",
				contentType:     "text/plain; charset=utf-8",
				handler:         GzipMiddleware(http.HandlerFunc(webhookPlain)),
				body:            nil,
			},
			want: want{
				response:        "OK",
				status:          http.StatusOK,
				acceptEncoding:  "",
				contentEncoding: "gzip",
				contentType:     "text/plain; charset=utf-8",
			},
		},
		{
			name: "no compressed [POST]",
			args: args{
				method:          http.MethodPost,
				acceptEncoding:  "gzip",
				contentEncoding: "",
				contentType:     "text/plain; charset=utf-8",
				handler:         GzipMiddleware(http.HandlerFunc(webhookPlain)),
				body:            nil,
			},
			want: want{
				response:        "OK",
				status:          http.StatusOK,
				acceptEncoding:  "",
				contentEncoding: "gzip",
				contentType:     "text/plain; charset=utf-8",
			},
		},
		{
			name: "no compressed json [GET]",
			args: args{
				method:          http.MethodGet,
				acceptEncoding:  "gzip",
				contentEncoding: "",
				contentType:     "text/plain; charset=utf-8",
				handler:         GzipMiddleware(http.HandlerFunc(webhookPlainJSON)),
				body:            nil,
			},
			want: want{
				response:        `{"status":"ok"}`,
				status:          http.StatusOK,
				acceptEncoding:  "",
				contentEncoding: "gzip",
				contentType:     "text/plain; charset=utf-8",
			},
		},
		{
			name: "no compressed [POST]",
			args: args{
				method:          http.MethodPost,
				acceptEncoding:  "gzip",
				contentEncoding: "",
				contentType:     "text/plain; charset=utf-8",
				handler:         GzipMiddleware(http.HandlerFunc(webhookPlainJSON)),
				body:            nil,
			},
			want: want{
				response:        `{"status":"ok"}`,
				status:          http.StatusOK,
				acceptEncoding:  "",
				contentEncoding: "gzip",
				contentType:     "text/plain; charset=utf-8",
			},
		},
		{
			name: "no compressed no body [GET]",
			args: args{
				method:          http.MethodGet,
				acceptEncoding:  "gzip",
				contentEncoding: "",
				contentType:     "application/json",
				handler:         GzipMiddleware(http.HandlerFunc(webhookJSON)),
				body:            nil,
			},
			want: want{
				response:        `{"status":"ok"}`,
				status:          http.StatusOK,
				acceptEncoding:  "",
				contentEncoding: "gzip",
				contentType:     "application/json",
			},
		},
		{
			name: "no comressed no body [POST]",
			args: args{
				method:          http.MethodPost,
				acceptEncoding:  "gzip",
				contentEncoding: "",
				contentType:     "application/json",
				handler:         GzipMiddleware(http.HandlerFunc(webhookJSON)),
				body:            nil,
			},
			want: want{
				response:        `{"status":"ok"}`,
				status:          http.StatusOK,
				acceptEncoding:  "",
				contentEncoding: "gzip",
				contentType:     "application/json",
			},
		},
		{
			name: "no compressed with body [POST]: positive",
			args: args{
				method:          http.MethodPost,
				acceptEncoding:  "gzip",
				contentEncoding: "",
				contentType:     "application/json",
				handler:         GzipMiddleware(http.HandlerFunc(webhookJSONWithBody)),
				body:            bytes.NewBuffer([]byte(requestBody)),
			},
			want: want{
				response:        response,
				status:          http.StatusOK,
				acceptEncoding:  "",
				contentEncoding: "gzip",
				contentType:     "application/json",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(tt.args.handler)
			defer srv.Close()

			r := httptest.NewRequest(tt.args.method, srv.URL, tt.args.body)
			r.RequestURI = ""

			require.Equal(t, 0, len(r.Header))
			require.Equal(t, 0, len(r.Header.Values("Accept-Encoding")))
			require.Equal(t, 0, len(r.Header.Values("Content-Encoding")))
			require.Equal(t, 0, len(r.Header.Values("Content-Type")))

			r.Header.Set("Accept-Encoding", tt.args.acceptEncoding)
			r.Header.Set("Content-Encoding", tt.args.contentEncoding)
			r.Header.Set("Content-Type", tt.args.contentType)

			resp, err := http.DefaultClient.Do(r)
			require.NoError(t, err)

			require.Equal(t, tt.want.status, resp.StatusCode)

			require.Equal(t, tt.want.acceptEncoding, resp.Header.Get("Accept-Encoding"))
			require.Equal(t, tt.want.contentEncoding, resp.Header.Get("Content-Encoding"))
			require.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))

			t.Log(resp.Header.Get("Accept-Encoding"))
			t.Log(resp.Header.Get("Content-Encoding"))
			t.Log(resp.Header.Get("Content-Type"))

			defer func() {
				e := resp.Body.Close()
				require.NoError(t, e)
			}()

			zr, err := gzip.NewReader(resp.Body)
			require.NoError(t, err)

			body, err := io.ReadAll(zr)
			require.NoError(t, err)

			require.Equal(t, tt.want.response, string(body))

			if tt.want.contentType == "application/json" {
				require.JSONEq(t, tt.want.response, string(body))
			}
		})
	}
}

func TestGzipMiddleware_AcceptGZIP(t *testing.T) {
	type args struct {
		handler         http.Handler
		body            io.Reader
		method          string
		contentType     string
		acceptEncoding  string
		contentEncoding string
	}

	type want struct {
		response        string
		contentType     string
		acceptEncoding  string
		contentEncoding string
		status          int
	}

	tz, err := time.LoadLocation("")
	require.NoError(t, err)

	now := time.Now().In(tz)
	hour, minute, _ := now.Clock()

	response := `{"response":{"text":"Точное время ` + fmt.Sprintf("%d", hour) + ` часов, ` + fmt.Sprintf("%d", minute) + ` минут. this is a very long string from client"}}`

	requestBody := `{
        "text": "this is a very long string from client"
    }`

	buf := bytes.NewBufferString(requestBody)

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "accept gzippped body [POST]: positive",
			args: args{
				method:          http.MethodPost,
				acceptEncoding:  "gzip",
				contentEncoding: "",
				contentType:     "application/json",
				handler:         GzipMiddleware(http.HandlerFunc(webhookJSONWithBody)),
				body:            buf,
			},
			want: want{
				response:        response,
				status:          http.StatusOK,
				acceptEncoding:  "",
				contentEncoding: "gzip",
				contentType:     "application/json",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(tt.args.handler)
			defer srv.Close()

			r := httptest.NewRequest(tt.args.method, srv.URL, tt.args.body)
			r.RequestURI = ""
			r.Header.Set("Accept-Encoding", tt.args.acceptEncoding)
			r.Header.Set("Content-Encoding", tt.args.contentEncoding)
			resp, err := http.DefaultClient.Do(r)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, resp.StatusCode)

			defer func() {
				e := resp.Body.Close()
				require.NoError(t, e)
			}()

			zr, err := gzip.NewReader(resp.Body)
			require.NoError(t, err)

			b, err := io.ReadAll(zr)
			require.NoError(t, err)

			t.Log(string(b))
			require.JSONEq(t, response, string(b))
		})
	}
}

func TestGzipMiddleware_SendGzip(t *testing.T) {
	var err error

	type args struct {
		handler         http.Handler
		body            io.Reader
		method          string
		contentType     string
		acceptEncoding  string
		contentEncoding string
	}

	type want struct {
		response        string
		contentType     string
		acceptEncoding  string
		contentEncoding string
		status          int
	}

	tz, err := time.LoadLocation("")
	require.NoError(t, err)

	now := time.Now().In(tz)
	hour, minute, _ := now.Clock()

	response := `{"response":{"text":"Точное время ` + fmt.Sprintf("%d", hour) + ` часов, ` + fmt.Sprintf("%d", minute) + ` минут. this is a very long string from client"}}`
	requestBody := `{
	        "text": "this is a very long string from client"
	    }`

	buf := bytes.NewBuffer(nil)
	zb := gzip.NewWriter(buf)
	_, err = zb.Write([]byte(requestBody))
	require.NoError(t, err)
	err = zb.Close()
	require.NoError(t, err)

	t.Log(buf)

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "no compressed [GET]",
			args: args{
				method:          http.MethodGet,
				acceptEncoding:  "",
				contentEncoding: "gzip",
				contentType:     "text/plain; charset=utf-8",
				handler:         GzipMiddleware(http.HandlerFunc(webhookPlain)),
				body:            nil,
			},
			want: want{
				response:        "",
				status:          http.StatusInternalServerError,
				acceptEncoding:  "gzip",
				contentEncoding: "",
				contentType:     "text/plain; charset=utf-8",
			},
		},
		{
			name: "no compressed [POST]",
			args: args{
				method:          http.MethodPost,
				acceptEncoding:  "",
				contentEncoding: "gzip",
				contentType:     "text/plain; charset=utf-8",
				handler:         GzipMiddleware(http.HandlerFunc(webhookPlain)),
				body:            nil,
			},
			want: want{
				response:        "",
				status:          http.StatusInternalServerError,
				acceptEncoding:  "gzip",
				contentEncoding: "",
				contentType:     "text/plain; charset=utf-8",
			},
		},
		{
			name: "no compressed json [GET]",
			args: args{
				method:          http.MethodGet,
				acceptEncoding:  "",
				contentEncoding: "gzip",
				contentType:     "text/plain; charset=utf-8",
				handler:         GzipMiddleware(http.HandlerFunc(webhookPlainJSON)),
				body:            nil,
			},
			want: want{
				response:        "",
				status:          http.StatusInternalServerError,
				acceptEncoding:  "gzip",
				contentEncoding: "",
				contentType:     "text/plain; charset=utf-8",
			},
		},
		{
			name: "no compressed [POST]",
			args: args{
				method:          http.MethodPost,
				acceptEncoding:  "",
				contentEncoding: "gzip",
				contentType:     "text/plain; charset=utf-8",
				handler:         GzipMiddleware(http.HandlerFunc(webhookPlainJSON)),
				body:            nil,
			},
			want: want{
				response:        "",
				status:          http.StatusInternalServerError,
				acceptEncoding:  "gzip",
				contentEncoding: "",
				contentType:     "text/plain; charset=utf-8",
			},
		},
		{
			name: "no compressed no body [GET]",
			args: args{
				method:          http.MethodGet,
				acceptEncoding:  "",
				contentEncoding: "gzip",
				contentType:     "application/json",
				handler:         GzipMiddleware(http.HandlerFunc(webhookJSON)),
				body:            nil,
			},
			want: want{
				response:        "",
				status:          http.StatusInternalServerError,
				acceptEncoding:  "gzip",
				contentEncoding: "",
				contentType:     "text/plain; charset=utf-8",
			},
		},
		{
			name: "no comressed no body [POST]",
			args: args{
				method:          http.MethodPost,
				acceptEncoding:  "",
				contentEncoding: "gzip",
				contentType:     "application/json",
				handler:         GzipMiddleware(http.HandlerFunc(webhookJSON)),
				body:            nil,
			},
			want: want{
				response:        "",
				status:          http.StatusInternalServerError,
				acceptEncoding:  "gzip",
				contentEncoding: "",
				contentType:     "text/plain; charset=utf-8",
			},
		},
		{
			name: "no compressed with body [POST]: positive",
			args: args{
				method:          http.MethodPost,
				acceptEncoding:  "",
				contentEncoding: "gzip",
				contentType:     "application/json",
				handler:         GzipMiddleware(http.HandlerFunc(webhookJSONWithBody)),
				body:            bytes.NewBuffer([]byte(requestBody)),
			},
			want: want{
				response:        "",
				status:          http.StatusInternalServerError,
				acceptEncoding:  "gzip",
				contentEncoding: "",
				contentType:     "text/plain; charset=utf-8",
			},
		},
		{
			name: "compressed with body [POST]: positive",
			args: args{
				method:          http.MethodPost,
				acceptEncoding:  "",
				contentEncoding: "gzip",
				contentType:     "application/json",
				handler:         GzipMiddleware(http.HandlerFunc(webhookJSONWithBody)),
				body:            buf,
			},
			want: want{
				response:        response,
				status:          http.StatusOK,
				acceptEncoding:  "gzip",
				contentEncoding: "",
				contentType:     "application/json",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(tt.args.handler)
			defer srv.Close()

			r := httptest.NewRequest(tt.args.method, srv.URL, tt.args.body)
			r.RequestURI = ""

			require.Equal(t, 0, len(r.Header))
			require.Equal(t, 0, len(r.Header.Values("Accept-Encoding")))
			require.Equal(t, 0, len(r.Header.Values("Content-Encoding")))
			require.Equal(t, 0, len(r.Header.Values("Content-Type")))

			r.Header.Set("Accept-Encoding", tt.args.acceptEncoding)
			r.Header.Set("Content-Encoding", tt.args.contentEncoding)
			r.Header.Set("Content-Type", tt.args.contentType)

			resp, err := http.DefaultClient.Do(r)
			require.NoError(t, err)

			require.Equal(t, tt.want.status, resp.StatusCode)

			require.Equal(t, tt.want.acceptEncoding, resp.Header.Get("Accept-Encoding"))
			require.Equal(t, tt.want.contentEncoding, resp.Header.Get("Content-Encoding"))
			require.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))

			t.Log(resp.Header.Get("Accept-Encoding"))
			t.Log(resp.Header.Get("Content-Encoding"))
			t.Log(resp.Header.Get("Content-Type"))

			defer func() {
				e := resp.Body.Close()
				require.NoError(t, e)
			}()

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			t.Log(string(body))
			require.Equal(t, tt.want.response, string(body))

			if tt.want.contentType == "application/json" {
				require.JSONEq(t, tt.want.response, string(body))
			}
		})
	}
}

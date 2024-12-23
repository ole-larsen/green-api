package httpserver_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ole-larsen/green-api/internal/httpserver"
	"github.com/ole-larsen/green-api/internal/httpserver/router"
)

func TestNewHttpServer(t *testing.T) {
	type args struct {
		host string
		port int
	}

	tests := []struct {
		name string
		want *httpserver.HTTPServer
		args args
	}{
		{
			name: "test defailt http server",
			want: &httpserver.HTTPServer{
				Host: "",
				Port: 0,
			},
			args: args{
				host: "localhost",
				port: 8080,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := httpserver.NewHTTPServer(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHttpServer() = %v, want %v", got, tt.want)
			}

			s := httpserver.NewHTTPServer()

			require.Equal(t, "", s.Host)
			require.Equal(t, 0, s.Port)
			require.Nil(t, s.Router)

			s.SetHost(tt.args.host)
			s.SetPort(tt.args.port)

			require.Equal(t, tt.args.host, s.Host)
			require.Equal(t, tt.args.port, s.Port)

			require.Nil(t, s.Router)

			r := router.NewMux()

			s.SetRouter(r)

			require.NotNil(t, s.Router)
		})
	}
}

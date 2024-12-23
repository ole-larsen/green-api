package httpserver

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ole-larsen/green-api/internal/httpserver/router"
)

func TestNewHttpServer(t *testing.T) {
	type args struct {
		host string
		port int
	}

	tests := []struct {
		name string
		want *HTTPServer
		args args
	}{
		{
			name: "test defailt http server",
			want: &HTTPServer{
				host: "",
				port: 0,
			},
			args: args{
				host: "localhost",
				port: 8080,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHTTPServer(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHttpServer() = %v, want %v", got, tt.want)
			}

			s := NewHTTPServer()

			require.Equal(t, "", s.host)
			require.Equal(t, 0, s.port)
			require.Nil(t, s.router)

			s.SetHost(tt.args.host)
			s.SetPort(tt.args.port)

			require.Equal(t, tt.args.host, s.host)
			require.Equal(t, tt.args.port, s.port)

			require.Nil(t, s.router)

			r := router.NewMux()

			s.SetRouter(r)

			require.NotNil(t, s.router)
		})
	}
}

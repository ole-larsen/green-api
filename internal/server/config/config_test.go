package config_test

import (
	"flag"
	"fmt"
	"testing"

	"github.com/ole-larsen/green-api/internal/common"
	"github.com/ole-larsen/green-api/internal/server/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	defaultAddress = "localhost:8080"
	defaultSecret  = "supersecret"
)

func Test_NewConfig(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "#1 server config test. check only one config was created",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.GetConfig()

			configPtr := fmt.Sprintf("%p", cfg)

			assert.NotEqual(t, "", cfg.Host)
			assert.NotEqual(t, "", cfg.Port)
			assert.Equal(t, "", cfg.DSN)
			assert.Equal(t, defaultSecret, cfg.Secret)

			for i := 0; i < 10; i++ {
				c := config.GetConfig()
				cPtr := fmt.Sprintf("%p", c)
				assert.Equal(t, configPtr, cPtr)
			}

			cfg = &config.Config{}

			assert.Empty(t, cfg)

			cfg = config.InitConfig()

			assert.Empty(t, cfg)
		})
	}
}

func Test_InitConfig(t *testing.T) {
	type args struct {
		opts []func(*config.Config)
	}

	address := defaultAddress
	dsn := "postgresql://postgres:postgres@172.17.0.1:5432/yandex"
	secret := defaultSecret

	tests := []struct {
		name string
		args args
	}{
		{
			name: "test init config functional options with env variables",
			args: args{
				opts: []func(*config.Config){
					config.WithAddress(address, nil),
					config.WithDSN(dsn, nil),
					config.WithSecret(secret, nil),
				},
			},
		},
		{
			name: "test init config functional options with parsed flags",
			args: args{
				opts: []func(*config.Config){
					config.WithAddress("", &address),
					config.WithDSN("", &dsn),
					config.WithSecret("", &secret),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.InitConfig(tt.args.opts...)
			assert.Equal(t, "localhost", cfg.Host)
			assert.Equal(t, 8080, cfg.Port)
			assert.Equal(t, dsn, cfg.DSN)
			assert.Equal(t, secret, cfg.Secret)
		})
	}
}

func Test_parseFlags(t *testing.T) {
	address := defaultAddress
	dsn := ""
	secret := defaultSecret

	tests := []struct {
		want config.Opts
		name string
	}{
		{
			name: "test parseFlags",
			want: config.Opts{
				APtr: &address,
				DPtr: &dsn,
				SPtr: &secret,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotEmpty(t, flag.Lookup("a"))
			assert.NotEmpty(t, flag.Lookup("d"))
			assert.NotEmpty(t, flag.Lookup("s"))

			// check default values
			assert.Equal(t, *tt.want.APtr, (flag.Lookup("a").Value.(flag.Getter).Get().(string)))
			assert.Equal(t, *tt.want.DPtr, (flag.Lookup("d").Value.(flag.Getter).Get().(string)))
			assert.Equal(t, *tt.want.SPtr, (flag.Lookup("s").Value.(flag.Getter).Get().(string)))
		})
	}
}

func Test_withAddress(t *testing.T) {
	type args struct {
		aPtr *string
		a    string
	}

	address := defaultAddress
	invalidAddressFormat := "localhost"       // Missing port
	invalidAddressPort := "localhost:notPort" // Invalid port

	tests := []struct {
		name      string
		args      args
		wantHost  string
		wantPort  int
		wantPanic bool
	}{
		{
			name: "valid address with environment variable",
			args: args{
				a:    address,
				aPtr: nil,
			},
			wantHost:  "localhost",
			wantPort:  8080,
			wantPanic: false,
		},
		{
			name: "valid address with command line argument",
			args: args{
				a:    "",
				aPtr: &address,
			},
			wantHost:  "localhost",
			wantPort:  8080,
			wantPanic: false,
		},
		{
			name: "invalid address format",
			args: args{
				a:    invalidAddressFormat,
				aPtr: nil,
			},
			wantPanic: true,
		},
		{
			name: "invalid address format with command line argument",
			args: args{
				a:    "",
				aPtr: &invalidAddressFormat,
			},
			wantPanic: true,
		},
		{
			name: "invalid port number",
			args: args{
				a:    invalidAddressPort,
				aPtr: nil,
			},
			wantPanic: true,
		},
		{
			name: "invalid port number with command line argument",
			args: args{
				a:    "",
				aPtr: &invalidAddressPort,
			},
			wantPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				assert.Panics(t, func() {
					config.InitConfig(config.WithAddress(tt.args.a, tt.args.aPtr))
				})
			} else {
				cfg := config.InitConfig(config.WithAddress(tt.args.a, tt.args.aPtr))
				assert.Equal(t, tt.wantHost, cfg.Host)
				assert.Equal(t, tt.wantPort, cfg.Port)
			}
		})
	}
}

func Test_WithSecret(t *testing.T) {
	type args struct {
		sPtr *string
		s    string
	}

	secret := defaultSecret

	tests := []struct {
		args args
		want func(*config.Config)
		name string
	}{
		{
			name: "withSecret test 1",
			args: args{
				s:    secret,
				sPtr: nil,
			},
		},
		{
			name: "withSecret test 2",
			args: args{
				s:    "",
				sPtr: &secret,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.InitConfig(config.WithSecret(tt.args.s, tt.args.sPtr))
			assert.Equal(t, secret, cfg.Secret)
		})
	}
}

func Test_withServeKey(t *testing.T) {
	type args struct {
		k string
	}

	tests := []struct {
		args args
		want func(*config.Config)
		name string
	}{
		{
			name: "withKey test",
			args: args{
				k: common.ServerMockKey,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.GetConfig()

			cfg.Reload(config.WithServerKey([]byte(tt.args.k)))
			require.Equal(t, []byte(tt.args.k), cfg.ServerKey)
		})
	}
}
func Test_withServerCert(t *testing.T) {
	type args struct {
		k string
	}

	tests := []struct {
		args args
		want func(*config.Config)
		name string
	}{
		{
			name: "withCert test",
			args: args{
				k: common.ServerMockKey,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.GetConfig()

			cfg.Reload(config.WithServerCrt([]byte(tt.args.k)))
			require.Equal(t, []byte(tt.args.k), cfg.ServerCrt)
		})
	}
}

func Test_withProtocol(t *testing.T) {
	type args struct {
		p string
	}

	tests := []struct {
		args args
		want func(*config.Config)
		name string
	}{
		{
			name: "withProtocol test",
			args: args{
				p: common.HTTPProtocol,
			},
		},
		{
			name: "withProtocol test",
			args: args{
				p: common.HTTPSProtocol,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.GetConfig()

			cfg.Reload(config.WithProtocol(tt.args.p))
			require.Equal(t, tt.args.p, cfg.Protocol)
		})
	}
}

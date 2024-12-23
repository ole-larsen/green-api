// Package config to configure server. The part of server package.
// Copyright 2024 The Oleg Nazarov. All rights reserved.
package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/ole-larsen/green-api/internal/common"
)

type Config struct {
	Host     string
	Protocol string
	DSN      string

	Secret    string
	ServerKey []byte
	ServerCrt []byte
	Port      int
}

type Opts struct {
	APtr *string
	DPtr *string
	SPtr *string
}

var (
	config = &Config{}
	once   sync.Once
)

// GetConfig rewrites using singleton pattern.
func GetConfig() *Config {
	once.Do(func() {
		// Если указана переменная окружения, то используется она.
		// Если нет переменной окружения, но есть аргумент командной строки (флаг), то используется он.
		// Если нет ни переменной окружения, ни флага, то используется значение по умолчанию.\
		f := parseFlags()
		config = InitConfig(
			WithAddress(os.Getenv("ADDRESS"), f.APtr),
			WithDSN(os.Getenv("DATABASE_DSN"), f.DPtr),
			WithSecret(os.Getenv("SECRET"), f.SPtr),
			WithProtocol(common.HTTPSProtocol),
		)
	})

	return config
}

func (c *Config) Reload(opts ...func(*Config)) {
	for _, opt := range opts {
		opt(c)
	}
}

func InitConfig(opts ...func(*Config)) *Config {
	c := &Config{}
	c.Reload(opts...)

	return c
}

func parseFlags() Opts {
	flags := Opts{
		APtr: flag.String("a", "localhost:8080", "адрес эндпоинта HTTP-сервера (по умолчанию localhost:8080)"),
		DPtr: flag.String("d", "", "строка с адресом подключения к БД"),
		SPtr: flag.String("s", "supersecret", "секрет для соли"),
	}

	flag.Parse()

	return flags
}

func WithDSN(d string, dPtr *string) func(*Config) {
	return func(c *Config) {
		if d == "" && dPtr != nil {
			d = *dPtr
		}

		c.DSN = d
	}
}

func WithAddress(a string, aPtr *string) func(*Config) {
	return func(c *Config) {
		if a == "" && aPtr != nil {
			a = *aPtr
		}

		addr := strings.Split(a, ":")

		const reqLen = 2

		if len(addr) != reqLen {
			panic(fmt.Errorf("wrong a parameters"))
		}

		c.Host = addr[0]

		port, err := strconv.Atoi(addr[1])
		if err != nil {
			panic(fmt.Errorf("wrong a parameters"))
		}

		c.Port = port
	}
}

func WithSecret(s string, sPtr *string) func(*Config) {
	return func(c *Config) {
		if s == "" && sPtr != nil {
			s = *sPtr
		}

		c.Secret = s
	}
}

func WithServerCrt(crt []byte) func(*Config) {
	return func(c *Config) {
		c.ServerCrt = crt
	}
}

func WithServerKey(key []byte) func(*Config) {
	return func(c *Config) {
		c.ServerKey = key
	}
}

func WithProtocol(protocol string) func(*Config) {
	return func(c *Config) {
		c.Protocol = protocol
	}
}

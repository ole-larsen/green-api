package main

import (
	"context"
	"embed"
	"fmt"

	"github.com/ole-larsen/green-api/internal/server"
	"github.com/ole-larsen/green-api/internal/server/config"
)

//go:embed certs/server.crt certs/server.key
var certsFS embed.FS

// main runs the build information and prints it to the provided writer.
func main() {
	certBytes, err := certsFS.ReadFile("certs/server.crt")
	if err != nil {
		fmt.Println("Error reading certificate:", err)
		return
	}

	keyBytes, err := certsFS.ReadFile("certs/server.key")
	if err != nil {
		fmt.Println("Error reading key:", err)
		return
	}

	fmt.Println("Certificate and key read successfully")

	ctx, cancel := context.WithCancel(context.Background())

	settings := config.GetConfig()
	settings.Reload(config.WithServerCrt(certBytes), config.WithServerKey(keyBytes))
	fmt.Printf("protocol: %s\n", settings.Protocol)

	srv, err := server.SetupFunc(ctx, settings)
	if err != nil {
		panic(err)
	}

	srv.Run(ctx, cancel)
}

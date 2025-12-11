package main

import (
	"embed"
	"flag"
	"fmt"
	"gomoco/internal/api"
	"gomoco/internal/server"
	"log"
)

//go:embed web/dist
var staticFiles embed.FS

var (
	port    = flag.Int("port", 8080, "API server port")
	version = flag.Bool("version", false, "Show version information")
)

const (
	appVersion = "1.1.0"
	appName    = "Gomoco"
)

func main() {
	flag.Parse()

	// Show version
	if *version {
		fmt.Printf("%s v%s\n", appName, appVersion)
		fmt.Println("A lightweight mock server written in Go")
		return
	}

	// Validate port
	if *port < 1 || *port > 65535 {
		log.Fatalf("Invalid port number: %d (must be between 1 and 65535)", *port)
	}

	// Initialize mock server manager
	manager := server.NewManager()

	// Start API server
	apiServer := api.NewServer(manager, staticFiles)

	addr := fmt.Sprintf(":%d", *port)
	log.Printf("Starting %s v%s on http://localhost%s", appName, appVersion, addr)
	if err := apiServer.Run(addr); err != nil {
		log.Fatal("Failed to start API server:", err)
	}
}

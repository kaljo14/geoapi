package main

import (
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/server"

	geomcp "github.com/neofyis/geopulse/internal/mcp"
)

func main() {
	apiURL := os.Getenv("GEOPULSE_API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8080"
	}

	s := server.NewMCPServer(
		"geopulse",
		"0.1.0",
		server.WithToolCapabilities(true),
	)

	client := geomcp.NewClient(apiURL)
	geomcp.RegisterTools(s, client)

	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "mcp server error: %v\n", err)
		os.Exit(1)
	}
}

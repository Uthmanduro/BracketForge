package main

import (
	"fmt"

	"github.com/uthmanduro/BracketForge/internal/config"
	"github.com/uthmanduro/BracketForge/internal/database"
	"github.com/uthmanduro/BracketForge/internal/server"
)

func main() {
	// Load configuration
	config, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	// Setup database connection
	db, err := database.SetupDB(config)
	if err != nil {
		fmt.Printf("Error setting up database: %v\n", err)
		return
	}

	// Initialize and start the server
	server := server.NewServer(config, db)
	if err := server.Start(); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
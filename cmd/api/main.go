package main

import (
	"fmt"

	"github.com/uthmanduro/BracketForge/internal/config"
	"github.com/uthmanduro/BracketForge/internal/database"
	"github.com/uthmanduro/BracketForge/internal/handler"
	"github.com/uthmanduro/BracketForge/internal/repository"
	"github.com/uthmanduro/BracketForge/internal/server"
	"github.com/uthmanduro/BracketForge/internal/service"
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

	// run database migrations here if needed
	if err := database.MigrateDB(db); err != nil {
		fmt.Printf("Error running database migrations: %v\n", err)
		return
	}

	// Initialize repositories, services, and server here if needed
	orgRepo := repository.NewOrganizationRepository(db)
	userRepo := repository.NewUserRepository(db)
	tournamentRepo := repository.NewTournamentRepo(db)

	// Initialize services with repositories
	orgService := service.NewOrganizationService(orgRepo)
	authService := service.NewAuthService(config, userRepo)
	tournamentService := service.NewTournamentService(tournamentRepo)

	// Initialize handler with services
	orgHandler := handler.NewOrganizationHandler(orgService)
	userHandler := handler.NewUserHandler(authService)
	tournamentHandler := handler.NewTournamentHandler(tournamentService)

	// Initialize and start the server
	server := server.NewServer(config, db, orgHandler, userHandler, tournamentHandler)
	if err := server.Start(); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
// @title           BracketForge API
// @version         1.0.0
// @description     Racket sports tournament management engine. Supports single elimination, round robin, and group + knockout formats with set-based scoring, three-way tiebreakers, BYE handling, and result correction.
// @termsOfService  http://swagger.io/terms/
 
// @contact.name   BracketForge Support
// @contact.email  support@bracketforge.io
 
// @host      localhost:8085
// @BasePath  /api/v1
 
// @securityDefinitions.apikey BearerAuth
// @in                         header
// @name                       Authorization
// @description                Enter your JWT token as: Bearer <token>

package main

import (
	"fmt"

	_ "github.com/uthmanduro/BracketForge/docs"

	"github.com/uthmanduro/BracketForge/internal/config"
	"github.com/uthmanduro/BracketForge/internal/database"
	"github.com/uthmanduro/BracketForge/internal/engine"
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

	store := repository.NewStore(db)

	// Initialize repositories, services, and server here if needed
	orgRepo := repository.NewOrganizationRepository(db)
	userRepo := repository.NewUserRepository(db)
	tournamentRepo := repository.NewTournamentRepo(db)
	playerRepo := repository.NewPlayerRepository(db)
	stageRepo  := repository.NewStageRepo(db)
	matchRepo  := repository.NewMatchRepo(db)
	standRepo  := repository.NewStandingsRepo(db)

	// ── Engines ───────────────────────────────────────────────────────
	rrEngine       := engine.NewRoundRobinEngine(matchRepo, standRepo)
	seEngine       := engine.NewSingleEliminationEngine(matchRepo)
	groupEngine    := engine.NewGroupEngine(stageRepo, matchRepo, standRepo, rrEngine)
	resultEngine   := engine.NewResultEngine(matchRepo, standRepo, stageRepo, store)
	walkoverEngine := engine.NewWalkoverEngine(matchRepo)
 

	// Initialize services with repositories
	orgService := service.NewOrganizationService(orgRepo)
	userService := service.NewUserService(userRepo, config.JWTSecret)
	tournamentService := service.NewTournamentService(
		tournamentRepo, stageRepo, matchRepo, standRepo,
		playerRepo, seEngine, rrEngine, groupEngine,
		resultEngine, walkoverEngine,
	)
	playerService := service.NewPlayerService(playerRepo)

	// Initialize handler with services
	orgHandler := handler.NewOrganizationHandler(orgService)
	userHandler := handler.NewUserHandler(userService)
	tournamentHandler := handler.NewTournamentHandler(tournamentService)
	playerHandler := handler.NewPlayerHandler(playerService)

	// Initialize and start the server
	server := server.NewServer(config, db, orgHandler, userHandler, tournamentHandler, playerHandler)
	if err := server.Start(); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
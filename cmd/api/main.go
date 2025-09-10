package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"foodflow/config"
	"foodflow/internal/core/repos"
	"foodflow/internal/core/services"
	"foodflow/internal/db"
	"foodflow/internal/http"
	"foodflow/internal/http/middleware"
	"foodflow/internal/lib"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Warn().Err(err).Msg("No .env file found, using system environment variables")
	}

	// Initialize logger
	lib.InitLogger()
	log.Info().Msg("Starting FoodFlow API server")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Initialize database connection
	database, err := db.NewConnection(cfg.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer database.Close()

	// Initialize JWT manager
	jwtManager, err := lib.NewJWTManager(cfg.JWT)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize JWT manager")
	}

	// Initialize repositories
	userRepo := repos.NewUserRepo(database.GetPool())
	profileRepo := repos.NewProfileRepo(database.GetPool())
	orgRepo := repos.NewOrgRepo(database.GetPool())
	collabRepo := repos.NewCollabRepo(database.GetPool())
	offerRepo := repos.NewOfferRepo(database.GetPool())
	claimRepo := repos.NewClaimRepo(database.GetPool())
	creditRepo := repos.NewCreditRepo(database.GetPool())
	tokenRepo := repos.NewTokenRepo(database.GetPool())
	redemptionRepo := repos.NewRedemptionRepo(database.GetPool())
	remoteRepo := repos.NewRemoteRepo(database.GetPool())

	// Initialize services
	authService := services.NewAuthService(userRepo, profileRepo, jwtManager)
	profileService := services.NewProfileService(profileRepo)
	orgService := services.NewOrgService(orgRepo)
	collabService := services.NewCollabService(collabRepo)
	offerService := services.NewOfferService(offerRepo, collabRepo, profileRepo)
	claimService := services.NewClaimService(claimRepo, offerRepo, orgRepo, creditRepo, profileRepo)
	creditService := services.NewCreditService(creditRepo, orgRepo)
	tokenService := services.NewTokenService(tokenRepo, collabRepo)
	redemptionService := services.NewRedemptionService(database, offerRepo, claimRepo, creditRepo, tokenRepo, redemptionRepo, cfg.App)
	remoteService := services.NewRemoteService(remoteRepo)
	matchingService := services.NewMatchingService(offerRepo, claimRepo, creditRepo, orgRepo, profileRepo, remoteRepo)
	adminService := services.NewAdminService(userRepo, orgRepo, collabRepo, offerRepo, claimRepo, redemptionRepo)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtManager)

	// Initialize router
	router := http.NewRouter(
		cfg,
		authService,
		profileService,
		orgService,
		collabService,
		offerService,
		claimService,
		redemptionService,
		creditService,
		tokenService,
		remoteService,
		adminService,
		matchingService,
		authMiddleware,
	)

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Info().Str("signal", sig.String()).Msg("Received shutdown signal")
		cancel()
	}()

	// Start HTTP server
	if err := router.Start(ctx); err != nil {
		log.Fatal().Err(err).Msg("Failed to start HTTP server")
	}

	log.Info().Msg("FoodFlow API server stopped")
}

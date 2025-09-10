package http

import (
	"context"
	"net/http"
	"time"

	"foodflow/config"
	"foodflow/internal/core"
	"foodflow/internal/http/handlers"
	"foodflow/internal/http/middleware"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Router struct {
	config            *config.Config
	authHandler       *handlers.AuthHandler
	profileHandler    *handlers.ProfileHandler
	orgHandler        *handlers.OrgHandler
	collabHandler     *handlers.CollabHandler
	offersHandler     *handlers.OffersHandler
	claimsHandler     *handlers.ClaimsHandler
	redemptionHandler *handlers.RedemptionHandler
	creditsHandler    *handlers.CreditsHandler
	tokensHandler     *handlers.TokensHandler
	remoteHandler     *handlers.RemoteHandler
	adminHandler      *handlers.AdminHandler
	authMiddleware    *middleware.AuthMiddleware
}

func NewRouter(
	config *config.Config,
	authService core.AuthService,
	profileService core.ProfileService,
	orgService core.OrgService,
	collabService core.CollabService,
	offerService core.OfferService,
	claimService core.ClaimService,
	redemptionService core.RedemptionService,
	creditService core.CreditService,
	tokenService core.TokenService,
	remoteService core.RemoteService,
	adminService core.AdminService,
	matchingService core.MatchingService,
	authMiddleware *middleware.AuthMiddleware,
) *Router {
	return &Router{
		config:            config,
		authHandler:       handlers.NewAuthHandler(authService),
		profileHandler:    handlers.NewProfileHandler(profileService),
		orgHandler:        handlers.NewOrgHandler(orgService),
		collabHandler:     handlers.NewCollabHandler(collabService),
		offersHandler:     handlers.NewOffersHandler(offerService, matchingService),
		claimsHandler:     handlers.NewClaimsHandler(claimService),
		redemptionHandler: handlers.NewRedemptionHandler(redemptionService),
		creditsHandler:    handlers.NewCreditsHandler(creditService),
		tokensHandler:     handlers.NewTokensHandler(tokenService),
		remoteHandler:     handlers.NewRemoteHandler(remoteService),
		adminHandler:      handlers.NewAdminHandler(adminService),
		authMiddleware:    authMiddleware,
	}
}

func (r *Router) SetupRoutes() *gin.Engine {
	// Set Gin mode
	if r.config.Server.Port == "8080" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"timestamp": time.Now().UTC(),
			"service":   "foodflow-api",
		})
	})

	// API v1 routes
	v1 := router.Group("/v1")
	{
		// Auth routes (no authentication required)
		auth := v1.Group("/auth")
		{
			auth.POST("/signup", r.authHandler.Signup)
			auth.POST("/login", r.authHandler.Login)
		}

		// Protected routes
		protected := v1.Group("")
		protected.Use(r.authMiddleware.RequireAuth())
		{
			// Auth
			protected.GET("/auth/me", r.authHandler.GetMe)

			// Profile
			profile := protected.Group("/profile")
			{
				profile.GET("", r.profileHandler.GetProfile)
				profile.PUT("", r.profileHandler.UpdateProfile)
			}

			// Organization routes
			org := protected.Group("/organizations")
			{
				org.POST("", r.orgHandler.CreateOrganization)
				org.GET("/:id", r.orgHandler.GetOrganization)
				org.PUT("/:id", r.orgHandler.UpdateOrganization)
				org.GET("", r.orgHandler.ListOrganizations)
			}

			// Collaborator routes
			collab := protected.Group("/collaborators")
			{
				collab.POST("", r.collabHandler.CreateCollaborator)
				collab.GET("/:id", r.collabHandler.GetCollaborator)
				collab.PUT("/:id", r.collabHandler.UpdateCollaborator)
				collab.GET("", r.collabHandler.ListCollaborators)
			}

			// Offers routes
			offers := protected.Group("/offers")
			{
				offers.POST("", r.offersHandler.CreateOffer)
				offers.GET("/:id", r.offersHandler.GetOffer)
				offers.PUT("/:id", r.offersHandler.UpdateOffer)
				offers.GET("", r.offersHandler.GetOffers)
				offers.GET("/nearby", r.offersHandler.GetNearbyOffers)
			}

			// Claims routes
			claims := protected.Group("/claims")
			{
				claims.POST("", r.claimsHandler.CreateClaim)
				claims.GET("/:id", r.claimsHandler.GetClaim)
				claims.PUT("/:id/status", r.claimsHandler.UpdateClaimStatus)
				claims.GET("", r.claimsHandler.ListClaims)
			}

			// Redemption routes
			redemption := protected.Group("/redemptions")
			{
				redemption.POST("", r.redemptionHandler.CreateRedemption)
				redemption.GET("/:id", r.redemptionHandler.GetRedemption)
				redemption.PUT("/:id/status", r.redemptionHandler.UpdateRedemptionStatus)
				redemption.GET("", r.redemptionHandler.ListRedemptions)
			}

			// Credits routes
			credits := protected.Group("/credits")
			{
				credits.GET("", r.creditsHandler.GetCredits)
				credits.GET("/history", r.creditsHandler.GetCreditHistory)
				credits.POST("/spend", r.creditsHandler.SpendCredits)
			}

			// Tokens routes
			tokens := protected.Group("/tokens")
			{
				tokens.GET("", r.tokensHandler.GetTokens)
				tokens.GET("/history", r.tokensHandler.GetTokenHistory)
				tokens.POST("/redeem", r.tokensHandler.RedeemTokens)
			}

			// Remote organization routes
			remote := protected.Group("/remote-orgs")
			{
				remote.POST("", r.remoteHandler.CreateRemoteOrg)
				remote.GET("/:id", r.remoteHandler.GetRemoteOrg)
				remote.PUT("/:id", r.remoteHandler.UpdateRemoteOrg)
				remote.GET("", r.remoteHandler.ListRemoteOrgs)
				remote.POST("/assign", r.remoteHandler.AssignRemoteOrg)
			}

			// Admin routes
			admin := protected.Group("/admin")
			{
				admin.GET("/stats", r.adminHandler.GetSystemStats)
				admin.GET("/users/stats", r.adminHandler.GetUserStats)
				admin.GET("/donations/stats", r.adminHandler.GetDonationStats)
				admin.GET("/users", r.adminHandler.ListUsers)
				admin.PUT("/users/:id/status", r.adminHandler.UpdateUserStatus)
				admin.GET("/audit-logs", r.adminHandler.GetAuditLogs)
			}
		}
	}

	return router
}

func (r *Router) Start(ctx context.Context) error {
	router := r.SetupRoutes()

	server := &http.Server{
		Addr:         ":" + r.config.Server.Port,
		Handler:      router,
		ReadTimeout:  r.config.Server.ReadTimeout,
		WriteTimeout: r.config.Server.WriteTimeout,
		IdleTimeout:  r.config.Server.IdleTimeout,
	}

	// Start server in goroutine
	go func() {
		log.Info().Str("port", r.config.Server.Port).Msg("Starting HTTP server")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Failed to start HTTP server")
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()

	// Graceful shutdown
	log.Info().Msg("Shutting down HTTP server")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return server.Shutdown(shutdownCtx)
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Idempotency-Key")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

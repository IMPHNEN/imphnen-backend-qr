package main

import (
	"log"

	"github.com/IMPHNEN/imphnen-backend-qr/internal/config"
	"github.com/IMPHNEN/imphnen-backend-qr/internal/handler"
	"github.com/IMPHNEN/imphnen-backend-qr/internal/middleware"
	"github.com/IMPHNEN/imphnen-backend-qr/internal/repository"
	"github.com/IMPHNEN/imphnen-backend-qr/internal/seeder"
	"github.com/IMPHNEN/imphnen-backend-qr/internal/service"
	"github.com/IMPHNEN/imphnen-backend-qr/pkg/database"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := config.Load()
	db := database.NewPostgres(cfg.DatabaseURL)
	defer db.Close()

	// Auto-seed demo users
	seeder.Run(db)

	// Repositories
	userRepo := repository.NewUserRepository(db)
	qrCampaignRepo := repository.NewQRCampaignRepository(db)

	// Services
	authService := service.NewAuthService(userRepo, cfg)
	userService := service.NewUserService(userRepo)
	qrCampaignService := service.NewQRCampaignService(qrCampaignRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	qrCampaignHandler := handler.NewQRCampaignHandler(qrCampaignService)

	// Echo
	e := echo.New()
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
	}))

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	// Auth routes (public)
	auth := e.Group("/api/v1/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
	auth.GET("/google", authHandler.GoogleAuth)
	auth.GET("/google/callback", authHandler.GoogleCallback)
	auth.POST("/refresh", authHandler.RefreshToken)

	// User routes (protected)
	users := e.Group("/api/v1/users")
	users.Use(middleware.JWTMiddleware(cfg.JWTSecret))

	users.GET("/me", userHandler.GetProfile)
	users.PUT("/me", userHandler.UpdateProfile)

	// Admin only routes
	admin := users.Group("")
	admin.Use(middleware.RBACMiddleware("admin"))
	admin.GET("", userHandler.GetAllUsers)
	admin.PUT("/:id/role", userHandler.UpdateUserRole)
	admin.DELETE("/:id", userHandler.DeleteUser)

	// Campaign routes (admin - JWT + RBAC)
	campaigns := e.Group("/api/v1/campaigns")
	campaigns.Use(middleware.JWTMiddleware(cfg.JWTSecret))

	adminCampaigns := campaigns.Group("")
	adminCampaigns.Use(middleware.RBACMiddleware("admin"))
	adminCampaigns.POST("", qrCampaignHandler.CreateCampaign)
	adminCampaigns.GET("", qrCampaignHandler.GetAllCampaigns)
	adminCampaigns.PUT("/:id/activate", qrCampaignHandler.SetActiveCampaign)
	adminCampaigns.DELETE("/:id", qrCampaignHandler.DeleteCampaign)

	// Campaign routes (user - JWT only, all roles)
	campaigns.POST("/process-image", qrCampaignHandler.ProcessImage)

	log.Printf("server starting on port %s", cfg.Port)
	e.Logger.Fatal(e.Start(":" + cfg.Port))
}

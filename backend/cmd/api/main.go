package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"spotslog/internal/config"
	"spotslog/internal/db"
	"spotslog/internal/handlers"
	"spotslog/internal/middleware"
	"spotslog/internal/storage"
)

func main() {
	_ = godotenv.Load()
	cfg := config.Load()

	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required (see .env.example)")
	}
	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET is required (see .env.example)")
	}

	ctx := context.Background()
	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	store, err := storage.New(cfg.S3Endpoint, cfg.S3AccessKey, cfg.S3SecretKey, cfg.S3Bucket, cfg.S3UseSSL, cfg.S3PublicBase)
	if err != nil {
		log.Fatalf("failed to init photo storage: %v", err)
	}

	authH := &handlers.AuthHandler{DB: pool, JWTSecret: cfg.JWTSecret, JWTExpiryHrs: cfg.JWTExpiryHrs}
	placesH := &handlers.PlacesHandler{DB: pool, Storage: store}
	visitsH := &handlers.VisitsHandler{DB: pool, Storage: store}
	savedH := &handlers.SavedHandler{DB: pool}

	r := gin.Default()
	r.Use(middleware.CORS(cfg.CORSAllowedOrigin))

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		auth.POST("/register", authH.Register)
		auth.POST("/login", authH.Login)
		auth.GET("/me", middleware.RequireAuth(cfg.JWTSecret), authH.Me)

		api.GET("/homepage", placesH.Homepage)

		places := api.Group("/places")
		places.GET("", placesH.List)
		places.GET("/:id", placesH.Get)
		authed := places.Group("")
		authed.Use(middleware.RequireAuth(cfg.JWTSecret))
		authed.POST("", placesH.Create)
		authed.PUT("/:id", placesH.Update)
		authed.DELETE("/:id", placesH.Delete)
		authed.PATCH("/:id/visibility", placesH.PatchVisibility)
		authed.POST("/:id/photos", placesH.UploadPhoto)
		authed.DELETE("/:id/photos/:photoId", placesH.DeletePhoto)

		users := api.Group("/users/me")
		users.Use(middleware.RequireAuth(cfg.JWTSecret))
		users.GET("/visits", visitsH.List)
		users.GET("/saved", savedH.List)

		visits := api.Group("/visits")
		visits.Use(middleware.RequireAuth(cfg.JWTSecret))
		visits.POST("", visitsH.Create)
		visits.PUT("/:id", visitsH.Update)
		visits.DELETE("/:id", visitsH.Delete)
		visits.POST("/:id/photos", visitsH.UploadPhoto)
		visits.DELETE("/:id/photos/:photoId", visitsH.DeletePhoto)

		saved := api.Group("/saved")
		saved.Use(middleware.RequireAuth(cfg.JWTSecret))
		saved.POST("", savedH.Create)
		saved.DELETE("/:placeId", savedH.Delete)
	}

	port := cfg.Port
	log.Printf("spotslog api listening on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

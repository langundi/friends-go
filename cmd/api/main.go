package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/langundi/friends-go/internal/bucket"
	"github.com/langundi/friends-go/internal/db"
	"github.com/langundi/friends-go/internal/db/store"
	"github.com/langundi/friends-go/internal/handlers"
	"github.com/langundi/friends-go/internal/services"
)

type config struct {
	addr   string
	db     dbConfig
	r2     r2Config
	secret string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

type r2Config struct {
	accountId  string
	accessKey  string
	secretKey  string
	bucketName string
	publicURL  string
}

func main() {
	godotenv.Load()

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	cfg := config{
		addr: getString("ADDR", ":8080"),
		db: dbConfig{
			addr:         getString("DB_ADDR", "postgresql://admin:adminpassword@localhost:5432/friends?sslmode=disable"),
			maxOpenConns: getInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: getInt("DB_MAX_IDLE_CONNS", 2),
			maxIdleTime:  getString("DB_MAX_IDLE_TIME", "15m"),
		},
		r2: r2Config{
			accountId:  getString("R2_ACCOUNT_ID", ""),
			accessKey:  getString("R2_ACCESS_KEY_ID", ""),
			secretKey:  getString("R2_SECRET_ACCESS_KEY", ""),
			bucketName: getString("R2_BUCKET_NAME", ""),
			publicURL:  getString("R2_PUBLIC_DEV_URL", ""),
		},
		secret: getString("SECRET_KEY", ""),
	}

	// Initialize Database
	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		slog.Error("database connection failed", "error", err)
	}

	defer db.Close()

	slog.Info("database connected")

	// Create Stores
	userStore := store.NewUserStore(db)
	refreshTokenStore := store.NewRefreshTokenStore(db)
	postStore := store.NewPostStore(db)

	ctx := context.Background()

	// Create Bucket
	r2Client, err := bucket.NewR2Client(ctx, cfg.r2.accountId, cfg.r2.accessKey, cfg.r2.secretKey)
	if err != nil {
		log.Fatalf("r2 setup failed: %v", err)
	}

	// Create Services
	authService := services.NewAuthService(userStore, refreshTokenStore, cfg.secret, 1*time.Hour)
	userService := services.NewUserService(userStore)
	postService := services.NewPostService(postStore, r2Client)

	// Create Handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	postHandler := handlers.NewPostHandler(postService, cfg.r2.bucketName, cfg.r2.publicURL)
	timelineHandler := handlers.NewTimelineHandler(postService)

	handlerCfg := handlers.HandlerConfig{
		AuthHandler:     authHandler,
		UserHandler:     userHandler,
		PostHandler:     postHandler,
		TimelineHandler: timelineHandler,
	}

	// Starting Server
	srv := &http.Server{
		Addr:         getString("ADDR", ":8080"),
		Handler:      handlers.Routes(handlerCfg),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	slog.Info("server starting", "port", "8080")

	if err := srv.ListenAndServe(); err != nil {
		slog.Error("server failed to start", "error", err)
		os.Exit(1)
	}
}

func getString(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return fallback
		}

		return intValue
	}

	return fallback
}

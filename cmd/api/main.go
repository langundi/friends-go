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
	"github.com/langundi/friends-go/internal/notification"
	"github.com/langundi/friends-go/internal/services"
)

type config struct {
	addr   string
	db     dbConfig
	r2     r2Config
	apns   apnsConfig
	secret string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

type r2Config struct {
	accountID  string
	accessKey  string
	secretKey  string
	bucketName string
	publicURL  string
}

type apnsConfig struct {
	authKeyPath string
	keyID       string
	teamID      string
	topic       string
	production  string
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
			accountID:  getString("R2_ACCOUNT_ID", ""),
			accessKey:  getString("R2_ACCESS_KEY_ID", ""),
			secretKey:  getString("R2_SECRET_ACCESS_KEY", ""),
			bucketName: getString("R2_BUCKET_NAME", ""),
			publicURL:  getString("R2_PUBLIC_DEV_URL", ""),
		},
		apns: apnsConfig{
			authKeyPath: "./AuthKey_FW2PL76896.p8",
			keyID:       getString("KEY_ID", ""),
			teamID:      getString("TEAM_ID", ""),
			topic:       getString("BUNDLE_ID", ""),
			production:  getString("APNS_ENV", "development"),
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
	friendStore := store.NewFriendStore(db)
	likeStore := store.NewLikeStore(db)
	replyStore := store.NewReplyStore(db)
	deviceTokenStore := store.NewDeviceTokenStore(db)

	ctx := context.Background()

	// Create R2 Client
	r2Client, err := bucket.NewR2Client(ctx, cfg.r2.accountID, cfg.r2.accessKey, cfg.r2.secretKey)
	if err != nil {
		log.Fatalf("r2 setup failed: %v", err)
	}

	// Create APNs Client
	apns, err := notification.NewAPNsClient(cfg.apns.authKeyPath, cfg.apns.keyID, cfg.apns.teamID, cfg.apns.topic, cfg.apns.production == "production")
	if err != nil {
		log.Fatalf("apns setup failed: %v", err)
	}

	// Create Services
	notificationService := services.NewNotificationService(apns, deviceTokenStore)
	authService := services.NewAuthService(userStore, refreshTokenStore, cfg.secret, 1*time.Hour)
	userService := services.NewUserService(userStore, r2Client)
	postService := services.NewPostService(postStore, likeStore, replyStore, deviceTokenStore, r2Client, notificationService)
	friendService := services.NewFriendService(friendStore)
	deviceTokenService := services.NewDeviceTokenService(deviceTokenStore)

	// Create Handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService, cfg.r2.bucketName, cfg.r2.publicURL)
	postHandler := handlers.NewPostHandler(postService, cfg.r2.bucketName, cfg.r2.publicURL)
	friendHandler := handlers.NewFriendHandler(friendService)
	deviceTokenHandler := handlers.NewDeviceHandler(deviceTokenService)

	handlerCfg := handlers.HandlerConfig{
		AuthHandler:        authHandler,
		UserHandler:        userHandler,
		PostHandler:        postHandler,
		FriendHandler:      friendHandler,
		DeviceTokenHandler: deviceTokenHandler,
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

package main

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/langundi/friends-go/internal/bucket"
	"github.com/langundi/friends-go/internal/db"
	"github.com/langundi/friends-go/internal/db/store"
	"github.com/langundi/friends-go/internal/handlers"
	"github.com/langundi/friends-go/internal/notification"
	"github.com/langundi/friends-go/internal/ratelimiter"
	"github.com/langundi/friends-go/internal/services"
)

type config struct {
	addr         string
	db           dbConfig
	r2           r2Config
	apns         apnsConfig
	rateLimitter ratelimiter.Config
	secret       string
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
		rateLimitter: ratelimiter.Config{
			RequestPerTimeFrame: getInt("RATE_LIMITER_REQUEST_COUNT", 20),
			TimeFrame:           time.Second * 5,
			Enabled:             getBool("RATE_LIMITER_ENABLED", true),
		},
		secret: getString("SECRET_KEY", ""),
	}

	apnsKey, err := loadAPNsKey(cfg.apns.authKeyPath)
	if err != nil {
		slog.Error("error loading apns key", "error", err)
		os.Exit(1)
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
	notificationStore := store.NewNotificationStore(db)

	ctx := context.Background()

	// Create R2 Client
	r2Client, err := bucket.NewR2Client(ctx, cfg.r2.accountID, cfg.r2.accessKey, cfg.r2.secretKey)
	if err != nil {
		log.Fatalf("r2 setup failed: %v", err)
	}

	// Create APNs Client
	apns, err := notification.NewAPNsClient(apnsKey, cfg.apns.keyID, cfg.apns.teamID, cfg.apns.topic, cfg.apns.production == "production")
	if err != nil {
		log.Fatalf("apns setup failed: %v", err)
	}

	// Create Services
	notificationService := services.NewNotificationService(notificationStore, deviceTokenStore, apns)
	authService := services.NewAuthService(userStore, refreshTokenStore, cfg.secret, 24*time.Hour)
	userService := services.NewUserService(userStore, r2Client)
	postService := services.NewPostService(postStore, likeStore, replyStore, deviceTokenStore, r2Client, notificationService)
	friendService := services.NewFriendService(friendStore, notificationService)
	deviceTokenService := services.NewDeviceTokenService(deviceTokenStore)

	// Create Handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService, cfg.r2.bucketName, cfg.r2.publicURL)
	postHandler := handlers.NewPostHandler(postService, cfg.r2.bucketName, cfg.r2.publicURL)
	friendHandler := handlers.NewFriendHandler(friendService)
	deviceTokenHandler := handlers.NewDeviceHandler(deviceTokenService)
	notificationHandler := handlers.NewNotificationHandler(notificationService)

	// rate limiter
	rateLimiter := ratelimiter.NewRateLimiter(
		cfg.rateLimitter.RequestPerTimeFrame,
		cfg.rateLimitter.TimeFrame,
	)

	handlerCfg := handlers.HandlerConfig{
		AuthHandler:         authHandler,
		UserHandler:         userHandler,
		PostHandler:         postHandler,
		FriendHandler:       friendHandler,
		DeviceTokenHandler:  deviceTokenHandler,
		NotificationHandler: notificationHandler,
	}

	mux := handlers.Routes(handlerCfg, cfg.rateLimitter, rateLimiter)
	if err := run(mux, cfg); err != nil {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}
}

func run(mux http.Handler, cfg config) error {
	srv := &http.Server{
		Addr:         cfg.addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	shutdown := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		slog.Info("signal caught", "signal", s.String())

		shutdown <- srv.Shutdown(ctx)
	}()

	slog.Info("server started", "addr", cfg.addr)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdown
	if err != nil {
		return err
	}

	slog.Info("server has stopped", "addr", cfg.addr)

	return nil
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

func loadAPNsKey(path string) ([]byte, error) {
	if b64 := os.Getenv("APNS_AUTH_KEY_B64"); b64 != "" {
		data, err := base64.StdEncoding.DecodeString(b64)
		if err != nil {
			return nil, fmt.Errorf("error decoding apns key: %w", err)
		}
		return data, nil
	}
	return os.ReadFile(path)
}

func getBool(key string, fallback bool) bool {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	boolVal, err := strconv.ParseBool(val)
	if err != nil {
		return fallback
	}

	return boolVal
}

package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/langundi/friends-go/internal/db/store"
	"github.com/langundi/friends-go/internal/types"
)

type AuthService struct {
	userStore         *store.UserStore
	refreshTokenStore *store.RefeshTokenStore
	jwtSecret         []byte
	accessTokenTTL    time.Duration
}

func NewAuthService(userStore *store.UserStore, refreshTokenStore *store.RefeshTokenStore, jwtSecret string, accessTokenTTL time.Duration) *AuthService {
	return &AuthService{
		userStore:         userStore,
		refreshTokenStore: refreshTokenStore,
		accessTokenTTL:    accessTokenTTL,
		jwtSecret:         []byte(jwtSecret),
	}
}

var (
	ErrEmptyFields        = errors.New("Username, email, and password are required.")
	ErrPasswordField      = errors.New("Password must be atleast 8 characters long.")
	ErrEmailExist         = errors.New("Email has already been registered.")
	ErrInvalidCredentials = errors.New("Email or password is incorrect.")
	ErrInvalidToken       = errors.New("Token is invalid.")
	ErrExpiredToken       = errors.New("Token has expired.")
)

func (s *AuthService) RegisterUser(ctx context.Context, req types.RegisterRequest) (*store.User, error) {
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return nil, ErrEmptyFields
	}

	if len(req.Password) < 8 {
		return nil, ErrPasswordField
	}

	_, err := s.userStore.GetUserByEmail(ctx, req.Email)
	if err == nil {
		return nil, ErrEmailExist
	}

	user := &store.User{
		Username: req.Username,
		Email:    req.Email,
	}

	if err := user.SetPassword(req.Password); err != nil {
		return nil, err
	}

	if err := s.userStore.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) LoginUser(ctx context.Context, req types.LoginRequest) (string, string, error) {
	if req.Email == "" || req.Password == "" {
		return "", "", ErrEmptyFields
	}

	user, err := s.userStore.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return "", "", ErrInvalidCredentials
	}

	if err := user.CheckPassword(req.Password); err != nil {
		return "", "", ErrInvalidCredentials
	}

	accessTokenStr, err := s.generateAccessToken(user)
	if err != nil {
		return "", "", err
	}

	refreshTokenStr, err := s.generateRefreshToken()
	if err != nil {
		return "", "", err
	}

	exp := time.Now().Add(7 * 24 * time.Hour)
	refreshToken := &store.RefreshToken{
		UserID:    user.ID,
		Token:     refreshTokenStr,
		ExpiresAt: exp,
		Revoked:   false,
	}

	if err := s.refreshTokenStore.CreateRefreshToken(ctx, refreshToken); err != nil {
		return "", "", err
	}

	return accessTokenStr, refreshTokenStr, nil
}

func (s *AuthService) LogoutUser(ctx context.Context, refreshToken string) error {
	if err := s.refreshTokenStore.DeleteRefreshToken(ctx, refreshToken); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) RefreshAccessToken(ctx context.Context, refreshToken string) (string, error) {
	token, err := s.refreshTokenStore.GetRefreshToken(ctx, refreshToken)

	if err != nil {
		return "", ErrInvalidToken
	}

	if token.Revoked {
		return "", ErrInvalidToken
	}

	if time.Now().After(token.ExpiresAt) {
		return "", ErrExpiredToken
	}

	user, err := s.userStore.GetUserByID(ctx, token.UserID)
	if err != nil {
		return "", err
	}

	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (s *AuthService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}

		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

func (s *AuthService) generateAccessToken(user *store.User) (string, error) {
	exp := time.Now().Add(s.accessTokenTTL).Unix()
	claims := jwt.MapClaims{
		"sub":  user.ID,
		"name": user.Username,
		"exp":  exp,
		"iat":  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) generateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

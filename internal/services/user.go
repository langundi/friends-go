package services

import (
	"context"
	"errors"

	"github.com/langundi/friends-go/internal/db/store"
	"github.com/langundi/friends-go/internal/types"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userStore *store.UserStore
}

var (
	ErrUserNotFound      = errors.New("User not found.")
	ErrDuplicateUsername = errors.New("Username already exists.")
	ErrDuplicateEmail    = errors.New("Email already registered.")
	ErrInvalidPassword   = errors.New("Current password is incorrect.")
)

func NewUserService(userStore *store.UserStore) *UserService {
	return &UserService{
		userStore: userStore,
	}
}

// Fetch user by ID
func (s *UserService) GetUserByID(ctx context.Context, id int64) (*store.User, error) {
	user, err := s.userStore.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, store.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

// Fetch user by username
func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*store.User, error) {
	user, err := s.userStore.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, store.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (s *UserService) ChangeUsername(ctx context.Context, username string, userID int64) error {
	err := s.userStore.ChangeUsername(ctx, username, userID)
	if err != nil {
		if errors.Is(err, store.ErrDuplicateUsername) {
			return ErrDuplicateUsername
		}
		return err
	}

	return nil
}

func (s *UserService) ChangeEmail(ctx context.Context, email string, userID int64) error {
	err := s.userStore.ChangeEmail(ctx, email, userID)
	if err != nil {
		if errors.Is(err, store.ErrDuplicateEmail) {
			return ErrDuplicateEmail
		}
		return err
	}

	return nil
}

func (s *UserService) ChangePassword(ctx context.Context, req types.ChangePasswordRequest, userID int64) error {
	user, err := s.userStore.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return err
	}

	if err := user.CheckPassword(req.CurrentPassword); err != nil {
		return ErrInvalidPassword
	}

	if err := user.SetPassword(req.NewPassword); err != nil {
		return err
	}

	if err := s.userStore.ChangePassword(ctx, user.Password, userID); err != nil {
		return err
	}

	return nil
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

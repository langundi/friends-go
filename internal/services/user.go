package services

import (
	"context"
	"errors"

	"github.com/langundi/friends-go/internal/db/store"
)

type UserService struct {
	userStore *store.UserStore
}

var (
	ErrUserNotFound      = errors.New("User not found.")
	ErrDuplicateUsername = errors.New("Username already exists.")
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

func (s *UserService) UpdateUsername(ctx context.Context, username string, userID int64) error {
	err := s.userStore.UpdateUsername(ctx, username, userID)
	if err != nil {
		if errors.Is(err, store.ErrDuplicateUsername) {
			return ErrDuplicateUsername
		}
		return err
	}

	return nil
}

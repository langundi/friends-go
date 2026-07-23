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
	ErrUserNotFound = errors.New("User not found.")
)

func NewUserService(userStore *store.UserStore) *UserService {
	return &UserService{
		userStore: userStore,
	}
}

func (s *UserService) GetUserProfile(ctx context.Context, id int64) (*store.User, error) {
	user, err := s.userStore.GetUserByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

package services

import (
	"context"

	"github.com/langundi/friends-go/internal/db/store"
)

type LikeService struct {
	likeStore *store.LikeStore
}

func NewLikeService(likeStore *store.LikeStore) *LikeService {
	return &LikeService{likeStore: likeStore}
}

func (s *LikeService) LikePost(ctx context.Context, userID, postID int64) error {
	err := s.likeStore.LikePost(ctx, userID, postID)
	if err != nil {
		return err
	}

	return nil
}

func (s *LikeService) UnlikePost(ctx context.Context, userID, postID int64) error {
	err := s.likeStore.UnlikePost(ctx, userID, postID)
	if err != nil {
		return err
	}

	return nil
}

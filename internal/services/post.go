package services

import (
	"context"
	"errors"

	"github.com/langundi/friends-go/internal/db/store"
	"github.com/langundi/friends-go/internal/types"
)

type PostService struct {
	postStore *store.PostStore
}

var (
	ErrPostNotFound = errors.New("Post not found.")
)

func NewPostService(postStore *store.PostStore) *PostService {
	return &PostService{postStore: postStore}
}

func (s *PostService) NewPost(ctx context.Context, req types.NewPostRequest) (*store.Post, error) {
	post := &store.Post{
		UserID:   req.UserID,
		Caption:  req.Caption,
		ImageURL: req.ImageURL,
	}

	err := s.postStore.CreatePost(ctx, post)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (s *PostService) GetPostByID(ctx context.Context, id int64) (*store.Post, error) {
	post, err := s.postStore.GetPostByID(ctx, id)
	if err != nil {
		return nil, ErrPostNotFound
	}

	return post, nil
}

func (s *PostService) GetLatestPost(ctx context.Context) (*store.Post, error) {
	post, err := s.postStore.GetLatestPost(ctx)
	if err != nil {
		return nil, ErrPostNotFound
	}

	return post, nil
}

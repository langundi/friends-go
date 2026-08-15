package services

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/langundi/friends-go/internal/db/store"
	"github.com/langundi/friends-go/internal/types"
)

type PostService struct {
	postStore *store.PostStore
	r2Client  *s3.Client
}

var (
	ErrPostNotFound = errors.New("Post not found.")
)

func NewPostService(postStore *store.PostStore, r2Client *s3.Client) *PostService {
	return &PostService{
		postStore: postStore,
		r2Client:  r2Client,
	}
}

func (s *PostService) NewPost(ctx context.Context, req types.NewPostRequest, userID int64) (*store.Post, error) {
	post := &store.Post{
		UserID:    userID,
		Caption:   req.Caption,
		ImageURL:  req.ImageURL,
		ObjectKey: req.ObjectKey,
	}

	err := s.postStore.CreatePost(ctx, post)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (s *PostService) DeletePost(ctx context.Context, id int64) error {
	err := s.postStore.DeletePostByID(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostService) GetPostByID(ctx context.Context, id int64) (*store.Post, error) {
	post, err := s.postStore.GetPostByID(ctx, id)
	if err != nil {
		return nil, ErrPostNotFound
	}

	return post, nil
}

func (s *PostService) GetPostsByUserID(ctx context.Context, userId int64) ([]store.Post, error) {
	posts, err := s.postStore.GetPostsByUserID(ctx, userId)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (s *PostService) GetLatestPosts(ctx context.Context) ([]store.Post, error) {
	posts, err := s.postStore.GetLatestPosts(ctx)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (s *PostService) GetLatestPost(ctx context.Context) (*store.Post, error) {
	post, err := s.postStore.GetLatestPost(ctx)
	if err != nil {
		return nil, ErrPostNotFound
	}

	return post, nil
}

func (s *PostService) PresignUploadURL(ctx context.Context, bucketName, objectKey, contentType string) (string, error) {
	presignClient := s3.NewPresignClient(s.r2Client)

	req, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(objectKey),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", err
	}

	return req.URL, nil
}

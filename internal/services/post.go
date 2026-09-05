package services

import (
	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	awstypes "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/langundi/friends-go/internal/db/store"
	"github.com/langundi/friends-go/internal/types"
)

type PostService struct {
	postStore           *store.PostStore
	likeStore           *store.LikeStore
	replyStore          *store.ReplyStore
	deviceTokenStore    *store.DeviceTokenStore
	r2Client            *s3.Client
	notificationService *NotificationService
}

var (
	ErrPostNotFound = errors.New("Post not found.")
)

func NewPostService(postStore *store.PostStore, likeStore *store.LikeStore, replyStore *store.ReplyStore, deviceTokenStore *store.DeviceTokenStore, r2Client *s3.Client, notificationService *NotificationService) *PostService {
	return &PostService{
		postStore:           postStore,
		likeStore:           likeStore,
		replyStore:          replyStore,
		deviceTokenStore:    deviceTokenStore,
		r2Client:            r2Client,
		notificationService: notificationService,
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

func (s *PostService) DeleteImage(ctx context.Context, bucketName, objectKey string) error {
	_, err := s.r2Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	return err
}

func (s *PostService) DeleteAllImage(ctx context.Context, bucketName string, objectKeys []string) error {
	var objects []awstypes.ObjectIdentifier
	for _, key := range objectKeys {
		objects = append(objects, awstypes.ObjectIdentifier{
			Key: aws.String(key),
		})
	}

	_, err := s.r2Client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
		Bucket: aws.String(bucketName),
		Delete: &awstypes.Delete{
			Objects: objects,
		},
	})
	return err
}

func (s *PostService) DeletePostByID(ctx context.Context, id int64) error {
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

func (s *PostService) GetMyPosts(ctx context.Context, userID int64) ([]store.Post, error) {
	posts, err := s.postStore.GetMyPosts(ctx, userID)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (s *PostService) GetFriendPosts(ctx context.Context, friendID int64, userID int64) ([]store.Post, error) {
	posts, err := s.postStore.GetFriendPosts(ctx, friendID, userID)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (s *PostService) GetTimeline(ctx context.Context, userID int64) ([]store.Post, error) {
	posts, err := s.postStore.GetTimeline(ctx, userID)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (s *PostService) GetMoreTimeline(ctx context.Context, userID int64, lastCreatedAt time.Time) ([]store.Post, error) {
	posts, err := s.postStore.GetMoreTimeline(ctx, userID, lastCreatedAt)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (s *PostService) LikePost(ctx context.Context, senderID, postID int64, req types.LikeNotificationRequest) error {
	err := s.likeStore.LikePost(ctx, senderID, postID)
	if err != nil {
		return err
	}
	if senderID == req.ReceiverID {
		return nil
	}

	tokens, err := s.deviceTokenStore.GetDeviceTokens(ctx, req.ReceiverID)
	if err != nil {
		return err
	}
	if len(tokens) == 0 {
		return nil
	}

	go s.notificationService.NotifyLike(tokens, req.SenderUsername)
	return nil
}

func (s *PostService) UnlikePost(ctx context.Context, userID, postID int64) error {
	err := s.likeStore.UnlikePost(ctx, userID, postID)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostService) GetRepliesForPost(ctx context.Context, userID, postID int64) ([]store.Reply, error) {
	replies, err := s.replyStore.GetRepliesForPostID(ctx, userID, postID)
	if err != nil {
		return nil, err
	}
	return replies, nil
}

func (s *PostService) ReplyPost(ctx context.Context, userID, postID int64, req types.ReplyRequest) (*store.Reply, error) {
	reply := &store.Reply{
		UserID:      userID,
		PostID:      postID,
		Reply:       req.Reply,
		RepliedByMe: true,
		Username:    req.Username,
	}

	err := s.replyStore.CreateReply(ctx, reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func (s *PostService) DeleteReply(ctx context.Context, replyID int64) error {
	err := s.replyStore.DeleteReply(ctx, replyID)
	if err != nil {
		return err
	}
	return nil
}

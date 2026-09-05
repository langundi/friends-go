package services

import (
	"context"
	"fmt"

	"github.com/langundi/friends-go/internal/db/store"
	"github.com/langundi/friends-go/internal/notification"
	"github.com/langundi/friends-go/internal/types"
)

type NotificationService struct {
	deviceStore *store.DeviceTokenStore
	apns        *notification.APNsClient
}

type NotificationAction int

const (
	Like NotificationAction = iota
	Reply
)

func NewNotificationService(apns *notification.APNsClient, deviceStore *store.DeviceTokenStore) *NotificationService {
	return &NotificationService{
		apns:        apns,
		deviceStore: deviceStore,
	}
}

// Notify like
func (s *NotificationService) NotifyLike(ctx context.Context, req types.LikeNotificationRequest) error {
	tokens, err := s.deviceStore.GetDeviceTokens(ctx, req.ReceiverID)
	if err != nil {
		return err
	}

	if len(tokens) == 0 {
		return nil
	}

	message := messageBuilder(req.SenderUsername, nil, Like)

	go func() {
		for _, t := range tokens {
			s.apns.SendNotification(t.Token, message)
		}
	}()
	return nil
}

// Notify reply
func (s *NotificationService) NotifyReply(ctx context.Context, req types.ReplyRequest) error {
	tokens, err := s.deviceStore.GetDeviceTokens(ctx, req.ReceiverID)
	if err != nil {
		return err
	}

	if len(tokens) == 0 {
		return nil
	}

	message := messageBuilder(req.Username, &req.Reply, Reply)

	go func() {
		for _, t := range tokens {
			s.apns.SendNotification(t.Token, message)
		}
	}()
	return nil
}

func messageBuilder(username string, reply *string, action NotificationAction) string {
	switch action {
	case Like:
		return username + " liked your post."
	case Reply:
		return username + " replied to your post: " + *reply
	default:
		panic(fmt.Errorf("unknown action: %v", action))
	}
}

// Notify Friend Request
// Notify Accept Friend Request

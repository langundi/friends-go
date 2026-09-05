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

// Notify Like
func (s *NotificationService) NotifyLike(ctx context.Context, req types.LikeNotificationRequest) error {
	tokens, err := s.deviceStore.GetDeviceTokens(ctx, req.ReceiverID)
	if err != nil {
		return err
	}

	if len(tokens) == 0 {
		return nil
	}

	title := titleBuilder(Like)
	message := messageBuilder(req.SenderUsername, Like)

	for _, t := range tokens {
		err := s.apns.SendNotification(t.Token, title, message)
		if err != nil {
			return err
		}
	}
	return nil
}

func titleBuilder(action NotificationAction) string {
	switch action {
	case Like:
		return "Someone liked your post."
	case Reply:
		return "Someone replied to your post."
	default:
		panic(fmt.Errorf("unknown action: %v", action))
	}
}

func messageBuilder(username string, action NotificationAction) string {
	switch action {
	case Like:
		return username + " liked your post."
	case Reply:
		return username + " replied to your post."
	default:
		panic(fmt.Errorf("unknown action: %v", action))
	}
}

// Notify Reply
// Notify Friend Request
// Notify Accept Friend Request

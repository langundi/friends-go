package services

import (
	"context"
	"fmt"

	"github.com/langundi/friends-go/internal/db/store"
	"github.com/langundi/friends-go/internal/notification"
	"github.com/langundi/friends-go/internal/types"
)

type NotificationService struct {
	notificationStore *store.NotificationStore
	deviceStore       *store.DeviceTokenStore
	apns              *notification.APNsClient
}

type NotificationAction int

const (
	Like NotificationAction = iota
	Reply
	Sent
	Accept
)

func NewNotificationService(notificationStore *store.NotificationStore, deviceStore *store.DeviceTokenStore, apns *notification.APNsClient) *NotificationService {
	return &NotificationService{
		notificationStore: notificationStore,
		deviceStore:       deviceStore,
		apns:              apns,
	}
}

// Notify like
func (s *NotificationService) NotifyLike(ctx context.Context, userID, postID int64, req types.LikeNotificationRequest) error {
	tokens, err := s.deviceStore.GetDeviceTokens(ctx, req.ReceiverID)
	if err != nil {
		return err
	}

	if len(tokens) == 0 {
		return nil
	}

	message := messageBuilder(req.SenderUsername, nil, Like)
	notification := &store.Notification{
		ReceiverID: req.ReceiverID,
		SenderID:   userID,
		PostID:     postID,
		Category:   "like",
		Message:    message,
	}

	if err := s.notificationStore.CreateNotification(ctx, notification); err != nil {
		return err
	}

	go func() {
		for _, t := range tokens {
			s.apns.SendNotification(t.Token, message)
		}
	}()
	return nil
}

// Notify post owner
func (s *NotificationService) NotifyPostReply(ctx context.Context, userID, postID int64, req types.ReplyRequest) error {
	tokens, err := s.deviceStore.GetDeviceTokens(ctx, req.PostOwnerID)
	if err != nil {
		return err
	}

	if len(tokens) == 0 {
		return nil
	}

	message := messageBuilder(req.Username, &req.Reply, Reply)
	notification := &store.Notification{
		ReceiverID: req.PostOwnerID,
		SenderID:   userID,
		PostID:     postID,
		Category:   "reply",
		Message:    message,
	}

	if err := s.notificationStore.CreateNotification(ctx, notification); err != nil {
		return err
	}

	go func() {
		for _, t := range tokens {
			s.apns.SendNotification(t.Token, message)
		}
	}()

	return nil
}

// Notify user reply
func (s *NotificationService) NotifyUserReply(ctx context.Context, userID, postID int64, req types.ReplyRequest) error {
	tokens, err := s.deviceStore.GetDeviceTokens(ctx, req.ReceiverID)
	if err != nil {
		return err
	}

	if len(tokens) == 0 {
		return nil
	}

	message := messageBuilder(req.Username, &req.Reply, Reply)
	notification := &store.Notification{
		ReceiverID: req.ReceiverID,
		SenderID:   userID,
		PostID:     postID,
		Category:   "reply",
		Message:    message,
	}

	if err := s.notificationStore.CreateNotification(ctx, notification); err != nil {
		return err
	}

	go func() {
		for _, t := range tokens {
			s.apns.SendNotification(t.Token, message)
		}
	}()

	return nil
}

func (s *NotificationService) NotifySentFriendRequest(ctx context.Context, req types.SendFriendRequestNotification) error {
	tokens, err := s.deviceStore.GetDeviceTokens(ctx, req.ReceiverID)
	if err != nil {
		return err
	}

	if len(tokens) == 0 {
		return nil
	}

	message := messageBuilder(req.SenderUsername, nil, Sent)
	go func() {
		for _, t := range tokens {
			s.apns.SendNotification(t.Token, message)
		}
	}()
	return nil
}

func (s *NotificationService) NotifyAcceptFriendRequest(ctx context.Context, req types.AcceptFriendRequestNotification) error {
	tokens, err := s.deviceStore.GetDeviceTokens(ctx, req.SenderID)
	if err != nil {
		return err
	}

	if len(tokens) == 0 {
		return nil
	}

	message := messageBuilder(req.SenderUsername, nil, Accept)
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
		return username + " replied: " + *reply
	case Sent:
		return username + " sent you a friend request."
	case Accept:
		return username + " accepted your friend request."
	default:
		panic(fmt.Errorf("unknown action: %v", action))
	}
}

func (s *NotificationService) GetAllNotifications(ctx context.Context, userID int64) ([]store.Notification, error) {
	notifications, err := s.notificationStore.GetAllNotifications(ctx, userID)
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

func (s *NotificationService) ReadNotifications(ctx context.Context) error {
	err := s.notificationStore.ReadNotifications(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *NotificationService) DeleteLikeNotification(ctx context.Context, senderID, postID int64) error {
	if err := s.notificationStore.DeleteNotificationByAction(ctx, senderID, postID, "like"); err != nil {
		return err
	}
	return nil
}

func (s *NotificationService) DeleteReplyNotification(ctx context.Context, senderID, postID int64) error {
	if err := s.notificationStore.DeleteNotificationByAction(ctx, senderID, postID, "reply"); err != nil {
		return err
	}
	return nil
}

// Notify Friend Request
// Notify Accept Friend Request

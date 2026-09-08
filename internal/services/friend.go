package services

import (
	"context"

	"github.com/langundi/friends-go/internal/db/store"
	"github.com/langundi/friends-go/internal/types"
)

type FriendService struct {
	friendStore         *store.FriendStore
	notificationService *NotificationService
}

type FriendshipStatus string

const (
	StatusNotAdded FriendshipStatus = "NOT_ADDED"
	StatusSent     FriendshipStatus = "SENT"
	StatusReceived FriendshipStatus = "RECEIVED"
	StatusFriends  FriendshipStatus = "FRIENDS"
)

func NewFriendService(friendStore *store.FriendStore, notificationService *NotificationService) *FriendService {
	return &FriendService{
		friendStore:         friendStore,
		notificationService: notificationService,
	}
}

// Send a friend request with a sender and receiver ID
func (s *FriendService) SendFriendRequest(ctx context.Context, senderID int64, req types.SendFriendRequestNotification) (*store.NewFriendRequest, error) {
	friend, err := s.friendStore.CreateFriendRequest(ctx, senderID, req.ReceiverID)
	if err != nil {
		return nil, err
	}

	if err := s.notificationService.NotifySentFriendRequest(ctx, req); err != nil {
		return nil, err
	}

	return friend, nil
}

// Accept a friend request
func (s *FriendService) AcceptFriendRequestByID(ctx context.Context, id int64, req types.AcceptFriendRequestNotification) error {
	if err := s.friendStore.AcceptFriendByID(ctx, id); err != nil {
		return err
	}

	if err := s.notificationService.NotifyAcceptFriendRequest(ctx, req); err != nil {
		return err
	}

	return nil
}

// Get all friend request for a user
func (s *FriendService) GetFriendRequestsForUserID(ctx context.Context, userID int64) ([]store.FriendRequest, error) {
	list, err := s.friendStore.GetFriendRequestsForUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// Get friendship status for searching user
func (s *FriendService) GetFriendshipStatus(ctx context.Context, currentUserID, searchedUserID int64) (string, error) {
	rel, err := s.friendStore.GetFriendshipStatus(ctx, currentUserID, searchedUserID)
	if err != nil {
		return "", err
	}

	if rel == nil {
		return string(StatusNotAdded), nil
	}

	if rel.Status == "accepted" {
		return string(StatusFriends), nil
	}

	if rel.SenderID == currentUserID {
		return string(StatusSent), nil
	}

	return string(StatusReceived), nil
}

// Delete or Decline a friend request
func (s *FriendService) DeclineOrUnfriendFriendByID(ctx context.Context, id int64) error {
	err := s.friendStore.DeleteFriendByID(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

// Get friend list
func (s *FriendService) GetFriendListForUserID(ctx context.Context, userID int64) ([]store.Friend, error) {
	list, err := s.friendStore.GetFriendListForUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return list, nil
}

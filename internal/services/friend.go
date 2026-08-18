package services

import (
	"context"

	"github.com/langundi/friends-go/internal/db/store"
)

type FriendService struct {
	friendStore *store.FriendStore
}

type FriendshipStatus string

const (
	StatusNotAdded FriendshipStatus = "NOT_ADDED"
	StatusSent     FriendshipStatus = "SENT"
	StatusReceived FriendshipStatus = "RECEIVED"
	StatusFriends  FriendshipStatus = "FRIENDS"
)

func NewFriendService(friendStore *store.FriendStore) *FriendService {
	return &FriendService{friendStore: friendStore}
}

// Create a friend request with a sender and receiver ID
func (s *FriendService) CreateFriendRequest(ctx context.Context, senderID int64, receiverID int64) (*store.NewFriendRequest, error) {
	friend, err := s.friendStore.CreateFriendRequest(ctx, senderID, receiverID)
	if err != nil {
		return nil, err
	}

	return friend, nil
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
func (s *FriendService) GetFriendshipStatus(ctx context.Context, currentUserId, searchedUserId int64) (string, error) {
	rel, err := s.friendStore.GetFriendshipStatus(ctx, currentUserId, searchedUserId)
	if err != nil {
		return "", err
	}

	if rel == nil {
		return string(StatusNotAdded), nil
	}

	if rel.Status == "accepted" {
		return string(StatusFriends), nil
	}

	if rel.SenderID == currentUserId {
		return string(StatusSent), nil
	}

	return string(StatusReceived), nil
}

// Delete or Decline a friend request
func (s *FriendService) DeleteFriendRequestByID(ctx context.Context, id int64) error {
	err := s.friendStore.DeleteFriendRequestByID(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

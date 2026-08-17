package services

import (
	"context"

	"github.com/langundi/friends-go/internal/db/store"
)

type FriendService struct {
	friendStore *store.FriendStore
}

func NewFriendService(friendStore *store.FriendStore) *FriendService {
	return &FriendService{friendStore: friendStore}
}

func (s *FriendService) CreateFriendRequest(ctx context.Context, senderID int64, receiverID int64) (*store.NewFriendRequest, error) {
	friend, err := s.friendStore.CreateFriendRequest(ctx, senderID, receiverID)
	if err != nil {
		return nil, err
	}

	return friend, nil
}

func (s *FriendService) GetFriendRequestsForUserID(ctx context.Context, userID int64) ([]store.FriendRequest, error) {
	list, err := s.friendStore.GetFriendRequestsForUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (s *FriendService) DeleteFriendRequestByID(ctx context.Context, id int64) error {
	err := s.friendStore.DeleteFriendRequestByID(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

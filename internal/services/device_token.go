package services

import (
	"context"

	"github.com/langundi/friends-go/internal/db/store"
	"github.com/langundi/friends-go/internal/types"
)

type DeviceTokenService struct {
	deviceTokenStore *store.DeviceTokenStore
}

func NewDeviceTokenService(deviceTokenStore *store.DeviceTokenStore) *DeviceTokenService {
	return &DeviceTokenService{deviceTokenStore: deviceTokenStore}
}

func (s *DeviceTokenService) RegisterDeviceToken(ctx context.Context, userID int64, req types.DeviceTokenRequest) error {
	err := s.deviceTokenStore.CreateDeviceToken(ctx, userID, req.DeviceToken)
	if err != nil {
		return err
	}
	return nil
}

func (s *DeviceTokenService) DeleteDeviceToken(ctx context.Context, userID int64, req types.DeviceTokenRequest) error {
	err := s.deviceTokenStore.DeleteDeviceToken(ctx, userID, req.DeviceToken)
	if err != nil {
		return err
	}
	return nil
}

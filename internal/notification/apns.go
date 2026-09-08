package notification

import (
	"errors"
	"log/slog"

	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/payload"
	"github.com/sideshow/apns2/token"
)

type APNsClient struct {
	Client *apns2.Client
	Topic  string
}

var (
	ErrInvalidAuthKey = errors.New("Auth key is invalid.")
)

func NewAPNsClient(authKey []byte, keyID, teamID, topic string, production bool) (*APNsClient, error) {
	key, err := token.AuthKeyFromBytes(authKey)
	if err != nil {
		return nil, ErrInvalidAuthKey
	}

	token := &token.Token{
		AuthKey: key,
		KeyID:   keyID,
		TeamID:  teamID,
	}

	client := apns2.NewTokenClient(token)
	if production {
		client = client.Production()
	} else {
		client = client.Development()
	}

	apns := &APNsClient{
		Client: client,
		Topic:  topic,
	}

	return apns, nil
}

func (s *APNsClient) SendNotification(deviceToken, message string) {
	payload := payload.NewPayload()
	payload.AlertBody(message)

	notification := &apns2.Notification{
		DeviceToken: deviceToken,
		Topic:       s.Topic,
		Payload:     payload,
	}

	res, _ := s.Client.Push(notification)

	if res.Sent() {
		slog.Info("push notification sent", "id", res.ApnsID)
	} else {
		slog.Error("push notification not sent", "id", res.ApnsID, "status", res.StatusCode, "reson", res.Reason)
	}
}

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

func NewAPNsClient(authKeyPath, keyID, teamID, topic string, production bool) (*APNsClient, error) {
	authKey, err := token.AuthKeyFromFile(authKeyPath)
	if err != nil {
		return nil, ErrInvalidAuthKey
	}

	token := &token.Token{
		AuthKey: authKey,
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

func (s *APNsClient) SendNotification(deviceToken, title, message string) error {
	payload := payload.NewPayload()
	payload.AlertTitle(title)
	payload.AlertBody(message)

	notification := &apns2.Notification{
		DeviceToken: deviceToken,
		Topic:       s.Topic,
		Payload:     payload,
	}

	res, err := s.Client.Push(notification)
	if err != nil {
		return err
	}

	if res.Sent() {
		slog.Info("push notification sent", "id", res.ApnsID)
	} else {
		slog.Error("push notification not sent", "id", res.ApnsID, "status", res.StatusCode, "reson", res.Reason)
	}
	return nil
}

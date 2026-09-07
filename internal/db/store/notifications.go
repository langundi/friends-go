package store

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Notification struct {
	ID             int64
	ReceiverID     int64
	SenderID       int64
	PostID         int64
	Category       string
	Message        string
	IsRead         bool
	CreatedAt      time.Time
	ProfilePicture *string
}

type NotificationStore struct {
	db *pgxpool.Pool
}

func NewNotificationStore(db *pgxpool.Pool) *NotificationStore {
	return &NotificationStore{db: db}
}

func (s *NotificationStore) CreateNotification(ctx context.Context, notification *Notification) error {
	query := `
		INSERT INTO notifications (receiver_id, sender_id, post_id, category, message)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := s.db.Exec(ctx, query, notification.ReceiverID, notification.SenderID, notification.PostID, notification.Category, notification.Message)
	if err != nil {
		return fmt.Errorf("create notification: %w", err)
	}
	return nil
}

func (s *NotificationStore) DeleteNotification(ctx context.Context, id int64) error {
	query := `DELETE FROM notifications WHERE id = $1`
	_, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete notification: %w", err)
	}
	return nil
}

func (s *NotificationStore) DeleteNotificationByAction(ctx context.Context, senderID, postID int64, category string) error {
	query := `
		DELETE FROM notifications
		WHERE sender_id = $1
			AND post_id = $2
			AND category = $3
	`
	_, err := s.db.Exec(ctx, query, senderID, postID, category)
	if err != nil {
		return fmt.Errorf("delete notification by action: %w", err)
	}
	return nil
}

func (s *NotificationStore) GetAllNotifications(ctx context.Context, userID int64) ([]Notification, error) {
	query := `
		SELECT n.id, n.receiver_id, n.sender_id, n.post_id, n.category, n.message, n.is_read, n.created_at, u.profile_picture
		FROM notifications n
		JOIN users u ON u.id = n.sender_id
		WHERE receiver_id = $1
		ORDER BY created_at DESC
	`
	rows, err := s.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query notifications: %w", err)
	}

	notifications, err := pgx.CollectRows(rows, pgx.RowToStructByName[Notification])
	if err != nil {
		return nil, fmt.Errorf("collect notifications: %w", err)
	}

	return notifications, nil
}

func (s *NotificationStore) ReadNotifications(ctx context.Context) error {
	query := `
	 	UPDATE notifications
	 	SET is_read = TRUE
	  	WHERE is_read = FALSE
	`
	_, err := s.db.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("read notification: %w", err)
	}
	return nil
}

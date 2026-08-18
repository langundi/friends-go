package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NewFriendRequest struct {
	ID         int64
	SenderID   int64
	ReceiverID int64
	Status     string
	CreatedAt  time.Time
}

type FriendRequest struct {
	ID             int64
	SenderID       int64
	ReceiverID     int64
	Status         string
	CreatedAt      time.Time
	SenderUsername string
}

type FriendStore struct {
	db *pgxpool.Pool
}

func NewFriendStore(db *pgxpool.Pool) *FriendStore {
	return &FriendStore{db: db}
}

func (s *FriendStore) CreateFriendRequest(ctx context.Context, senderID, receiverID int64) (*NewFriendRequest, error) {
	query := `
		INSERT INTO friends (sender_id, receiver_id, status)
		VALUES ($1, $2, 'pending')
		RETURNING id, sender_id, receiver_id, status, created_at
	`

	var friend NewFriendRequest

	err := s.db.QueryRow(ctx, query, senderID, receiverID).Scan(
		&friend.ID,
		&friend.SenderID,
		&friend.ReceiverID,
		&friend.Status,
		&friend.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("insert friend request: %w", err)
	}

	return &friend, err
}

func (s *FriendStore) GetFriendRequestsForUserID(ctx context.Context, userID int64) ([]FriendRequest, error) {
	query := `
		SELECT f.id, f.sender_id, f.receiver_id, f.status, f.created_at, u.username AS sender_username
		FROM friends f
		JOIN users u ON u.id = f.sender_id
		WHERE f.receiver_id = $1 AND f.status = 'pending'
		ORDER BY f.created_at DESC
	`

	rows, err := s.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query incoming request: %w", err)
	}

	list, err := pgx.CollectRows(rows, pgx.RowToStructByName[FriendRequest])
	if err != nil {
		return nil, fmt.Errorf("collecting friend request: %w", err)
	}

	return list, nil
}

func (s *FriendStore) GetFriendshipStatus(ctx context.Context, currentUserID, searchedUserID int64) (*NewFriendRequest, error) {
	query := `
		SELECT id, sender_id, receiver_id, status, created_at
		FROM friends
		WHERE (sender_id = $1 AND receiver_id = $2)
		OR (sender_id = $2 AND receiver_id = $1)
	`

	var friendReq NewFriendRequest

	err := s.db.QueryRow(ctx, query, currentUserID, searchedUserID).Scan(
		&friendReq.ID,
		&friendReq.SenderID,
		&friendReq.ReceiverID,
		&friendReq.Status,
		&friendReq.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("checking relationship status: %w", err)
	}

	return &friendReq, nil
}

func (s *FriendStore) DeleteFriendRequestByID(ctx context.Context, id int64) error {
	query := `DELETE FROM friends WHERE id = $1`

	_, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *FriendStore) AcceptFriendRequestByID(ctx context.Context, id int64) error {
	query := `
		UPDATE friends
		SET status = 'accepted'
		WHERE id = $1
	`

	_, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

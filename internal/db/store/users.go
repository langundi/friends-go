package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"udpated_at"`
}

type UserStore struct {
	db *pgxpool.Pool
}

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrDuplicateUsername = errors.New("username already exists")
)

func NewUserStore(db *pgxpool.Pool) *UserStore {
	return &UserStore{db: db}
}

func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hash)

	return nil
}

func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

func (s *UserStore) CreateUser(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (username, email, password)
		VALUES($1, $2, $3)
		RETURNING id, created_at
	`

	err := s.db.QueryRow(ctx, query,
		user.Username,
		user.Email,
		user.Password,
	).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

func (s *UserStore) DeleteUserByID(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("delete user %v:", id)
	}

	return nil
}

func (s *UserStore) GetUserByID(ctx context.Context, id int64) (*User, error) {
	query := `
		SELECT id, username, email
		FROM users
		WHERE id = $1
	`

	var user User

	err := s.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("query user by id: %w", err)
	}

	return &user, nil
}

func (s *UserStore) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	query := `
		SELECT id, username
		FROM users
		WHERE username = $1
	`

	var user User

	err := s.db.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("query user by username: %w", err)
	}

	return &user, nil
}

func (s *UserStore) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, username, email, password, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user User

	err := s.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("query user by email %s:", email)
		}
		return nil, err
	}

	return &user, nil
}

func (s *UserStore) UpdateUsername(ctx context.Context, username string, userID int64) error {
	query := `UPDATE users SET username = $1 WHERE id = $2`

	_, err := s.db.Exec(ctx, query, username, userID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrDuplicateUsername
		}
		return err
	}

	return nil
}

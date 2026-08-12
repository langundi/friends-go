package types

import (
	"time"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type NewPostRequest struct {
	Caption  string `json:"caption"`
	ImageURL string `json:"image_url"`
}

type GetPostRequest struct {
	ID int64 `json:"id"`
}

type RegisterResponse struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshResponse struct {
	AccessToken string `json:"access_token"`
}

type PostResponse struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Caption   string    `json:"caption"`
	ImageURL  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
}

type UploadImageRequest struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
}

type UploadImageResponse struct {
	UploadURL string `json:"upload_url"`
	PublicURL string `json:"public_url"`
	ObjectKey string `json:"object_key"`
}

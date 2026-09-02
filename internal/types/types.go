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
	Caption   string `json:"caption"`
	ImageURL  string `json:"image_url"`
	ObjectKey string `json:"object_key"`
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
	ID             int64     `json:"id"`
	UserID         int64     `json:"user_id"`
	Caption        string    `json:"caption"`
	ImageURL       string    `json:"image_url"`
	ObjectKey      string    `json:"object_key"`
	LikeCount      int       `json:"like_count"`
	ReplyCount     int       `json:"reply_count"`
	CreatedAt      time.Time `json:"created_at"`
	LikedByMe      bool      `json:"liked_by_me"`
	Username       string    `json:"username"`
	ProfilePicture *string   `json:"profile_picture"`
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

type DeletePostRequest struct {
	ObjectKey string `json:"object_key"`
}

type DeleteAllImagesRequest struct {
	ObjectKeys []string `json:"object_keys"`
}

type UsernameResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

type NewFriendRequestResponse struct {
	ID         int64     `json:"id"`
	SenderID   int64     `json:"sender_id"`
	ReceiverID int64     `json:"receiver_id"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

type FriendRequestResponse struct {
	ID             int64     `json:"id"`
	SenderID       int64     `json:"sender_id"`
	SenderUsername string    `json:"sender_username"`
	ReceiverID     int64     `json:"receiver_id"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	ProfilePicture *string   `json:"profile_picture"`
}

type FriendshipStatusResponse struct {
	FriendshipStatus string `json:"friendship_status"`
}

type FriendResponse struct {
	ID             int64   `json:"id"`
	UserID         int64   `json:"user_id"`
	Username       string  `json:"username"`
	ProfilePicture *string `json:"profile_picture"`
}

type ReplyRequest struct {
	Reply    string `json:"reply"`
	Username string `json:"username"`
}

type ReplyResponse struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	PostID      int64     `json:"post_id"`
	Reply       string    `json:"reply"`
	CreatedAt   time.Time `json:"created_at"`
	RepliedByMe bool      `json:"replied_by_me"`
	Username    string    `json:"username"`
}

type SetProfilePictureRequest struct {
	ImageURL  string `json:"image_url"`
	ObjectKey string `json:"object_key"`
}

type SetProfilePictureResponse struct {
	ProfilePicture string `json:"profile_picture"`
	ObjectKey      string `json:"object_key"`
}

type ChangeUsernameRequest struct {
	Username string `json:"username"`
}

type ChangeEmailRequest struct {
	Email string `json:"email"`
}

type ChangePasswordRequest struct {
	Email           string `json:"email"`
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

type DeleteProfilePictureRequest struct {
	Objectkey string `json:"object_key"`
}

type MoreTimelineRequest struct {
	CreatedAt time.Time `json:"created_at"`
}

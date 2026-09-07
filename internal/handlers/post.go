package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"path"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/langundi/friends-go/internal/middlewares"
	"github.com/langundi/friends-go/internal/services"
	"github.com/langundi/friends-go/internal/types"
	"github.com/langundi/friends-go/internal/utils"
)

type PostHandler struct {
	postService *services.PostService
	bucketName  string
	publicURL   string
}

func NewPostHandler(postService *services.PostService, bucketName, publicURL string) *PostHandler {
	return &PostHandler{
		postService: postService,
		bucketName:  bucketName,
		publicURL:   publicURL,
	}
}

// Make a new post
func (h *PostHandler) NewPostHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserID(r)
	if !ok {
		utils.UnauthorizedError(w, r, errors.New("Unauthorized."))
		return
	}

	var req types.NewPostRequest

	if err := utils.ReadJson(w, r, &req); err != nil {
		utils.InvalidPayloadError(w, r, err)
		return
	}

	post, err := h.postService.NewPost(r.Context(), req, userID)
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	response := types.PostResponse{
		ID:        post.ID,
		UserID:    post.UserID,
		Caption:   post.Caption,
		ImageURL:  post.ImageURL,
		ObjectKey: post.ObjectKey,
		CreatedAt: post.CreatedAt,
	}

	utils.WriteJson(w, http.StatusCreated, utils.JsonResponse{
		Success: true,
		Data:    response,
	})
}

// Get a presigned URL to upload image to Cloudflare R2
func (h *PostHandler) GetPresignedURLHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserID(r)
	if !ok {
		utils.UnauthorizedError(w, r, errors.New("Unauthorized."))
		return
	}

	var req types.UploadImageRequest
	if err := utils.ReadJson(w, r, &req); err != nil {
		utils.InvalidPayloadError(w, r, err)
		return
	}

	objectKey := h.generateObjectKey(userID, req.Filename, "posts")
	url, err := h.postService.PresignUploadURL(r.Context(), h.bucketName, objectKey, req.ContentType)
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	response := types.UploadImageResponse{
		UploadURL: url,
		PublicURL: h.publicImageURL(objectKey),
		ObjectKey: objectKey,
	}

	utils.WriteJson(w, http.StatusOK, utils.JsonResponse{
		Success: true,
		Data:    response,
	})
}

// Delete existing post
func (h *PostHandler) DeletePostHadler(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	var req types.DeletePostRequest
	if err := utils.ReadJson(w, r, &req); err != nil {
		utils.InvalidPayloadError(w, r, err)
		return
	}

	// Delete the image from Cloudflare R2
	if err := h.postService.DeleteImage(r.Context(), h.bucketName, req.ObjectKey); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	// Delete post from database
	if err := h.postService.DeletePostByID(r.Context(), postID); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.JsonResponse{
		Success: true,
	})
}

// Delete all image from object storage, used for user account deleteion
func (h *PostHandler) DeleteAllImageHandler(w http.ResponseWriter, r *http.Request) {
	var req types.DeleteAllImagesRequest
	if err := utils.ReadJson(w, r, &req); err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	if err := h.postService.DeleteAllImage(r.Context(), h.bucketName, req.ObjectKeys); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.JsonResponse{
		Success: true,
	})
}

// Get timeline
func (h *PostHandler) GetTimelineHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserID(r)
	if !ok {
		utils.UnauthorizedError(w, r, ErrUnauthorized)
		return
	}

	posts, err := h.postService.GetTimeline(r.Context(), userID)
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	var data []types.PostResponse

	for _, v := range posts {
		response := types.PostResponse{
			ID:             v.ID,
			UserID:         v.UserID,
			ImageURL:       v.ImageURL,
			Caption:        v.Caption,
			ObjectKey:      v.ObjectKey,
			LikeCount:      v.LikeCount,
			ReplyCount:     v.ReplyCount,
			CreatedAt:      v.CreatedAt,
			LikedByMe:      v.LikedByMe,
			Username:       v.Username,
			ProfilePicture: v.ProfilePicture,
		}

		data = append(data, response)
	}

	utils.WriteJson(w, http.StatusOK, utils.JsonResponse{
		Success: true,
		Data:    data,
	})
}

// Get more timeline
func (h *PostHandler) GetMoreTimelineHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserID(r)
	if !ok {
		utils.UnauthorizedError(w, r, ErrUnauthorized)
		return
	}

	var req types.MoreTimelineRequest
	if err := utils.ReadJson(w, r, &req); err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	posts, err := h.postService.GetMoreTimeline(r.Context(), userID, req.CreatedAt)
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	var data []types.PostResponse

	for _, v := range posts {
		response := types.PostResponse{
			ID:             v.ID,
			UserID:         v.UserID,
			ImageURL:       v.ImageURL,
			Caption:        v.Caption,
			ObjectKey:      v.ObjectKey,
			LikeCount:      v.LikeCount,
			ReplyCount:     v.ReplyCount,
			CreatedAt:      v.CreatedAt,
			LikedByMe:      v.LikedByMe,
			Username:       v.Username,
			ProfilePicture: v.ProfilePicture,
		}

		data = append(data, response)
	}

	utils.WriteJson(w, http.StatusOK, utils.JsonResponse{
		Success: true,
		Data:    data,
	})
}

// Get a post by it's id
func (h *PostHandler) GetPostHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserID(r)
	if !ok {
		utils.UnauthorizedError(w, r, ErrUnauthorized)
		return
	}

	postID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	post, err := h.postService.GetPostByID(r.Context(), postID, userID)
	if err != nil {
		utils.NotFoundError(w, r, err)
		return
	}

	response := types.PostResponse{
		ID:             post.ID,
		UserID:         post.UserID,
		ImageURL:       post.ImageURL,
		Caption:        post.Caption,
		ObjectKey:      post.ObjectKey,
		LikeCount:      post.LikeCount,
		ReplyCount:     post.ReplyCount,
		CreatedAt:      post.CreatedAt,
		LikedByMe:      post.LikedByMe,
		Username:       post.Username,
		ProfilePicture: post.ProfilePicture,
	}

	utils.WriteJson(w, http.StatusOK, utils.JsonResponse{
		Success: true,
		Data:    response,
	})
}

func (h *PostHandler) LikePostHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserID(r)
	if !ok {
		utils.UnauthorizedError(w, r, errors.New("Unauthorized."))
		return
	}

	postID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	var req types.LikeNotificationRequest
	if err := utils.ReadJson(w, r, &req); err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	if err := h.postService.LikePost(r.Context(), userID, postID, req); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.JsonResponse{
		Success: true,
	})
}

func (h *PostHandler) UnlikePostHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserID(r)
	if !ok {
		utils.UnauthorizedError(w, r, errors.New("Unauthorized."))
		return
	}

	postID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	if err := h.postService.UnlikePost(r.Context(), userID, postID); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.JsonResponse{
		Success: true,
	})
}

func (h *PostHandler) GetPostRepliesHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserID(r)
	if !ok {
		utils.UnauthorizedError(w, r, errors.New("Unauthorized."))
		return
	}

	postID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	replies, err := h.postService.GetRepliesForPost(r.Context(), userID, postID)
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	var data []types.ReplyResponse

	for _, v := range replies {
		response := types.ReplyResponse{
			ID:          v.ID,
			UserID:      v.UserID,
			PostID:      v.PostID,
			Reply:       v.Reply,
			CreatedAt:   v.CreatedAt,
			RepliedByMe: v.RepliedByMe,
			Username:    v.Username,
		}

		data = append(data, response)
	}

	utils.WriteJson(w, http.StatusCreated, utils.JsonResponse{
		Success: true,
		Data:    data,
	})
}

func (h *PostHandler) ReplyPostHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserID(r)
	if !ok {
		utils.UnauthorizedError(w, r, errors.New("Unauthorized."))
		return
	}

	postID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	var req types.ReplyRequest
	if err := utils.ReadJson(w, r, &req); err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	reply, err := h.postService.ReplyPost(r.Context(), userID, postID, req)
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	response := types.ReplyResponse{
		ID:          reply.ID,
		UserID:      reply.UserID,
		PostID:      reply.PostID,
		Reply:       reply.Reply,
		CreatedAt:   reply.CreatedAt,
		RepliedByMe: reply.RepliedByMe,
		Username:    reply.Username,
	}

	utils.WriteJson(w, http.StatusCreated, utils.JsonResponse{
		Success: true,
		Data:    response,
	})
}

func (h *PostHandler) DeleteReplyHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserID(r)
	if !ok {
		utils.UnauthorizedError(w, r, errors.New("Unauthorized."))
		return
	}

	postID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	var req types.DeleteReplyRequest
	if err := utils.ReadJson(w, r, &req); err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	if err := h.postService.DeleteReply(r.Context(), userID, postID, req.ReplyID); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.JsonResponse{
		Success: true,
	})
}

// Get current user posts, used for profile tab
func (h *PostHandler) GetMyPostsHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserID(r)
	if !ok {
		utils.UnauthorizedError(w, r, ErrUnauthorized)
		return
	}

	posts, err := h.postService.GetMyPosts(r.Context(), userID)
	if err != nil {
		utils.NotFoundError(w, r, err)
		return
	}

	var data []types.PostResponse

	for _, v := range posts {
		response := types.PostResponse{
			ID:             v.ID,
			UserID:         v.UserID,
			Caption:        v.Caption,
			ImageURL:       v.ImageURL,
			ObjectKey:      v.ObjectKey,
			LikeCount:      v.LikeCount,
			ReplyCount:     v.ReplyCount,
			CreatedAt:      v.CreatedAt,
			LikedByMe:      v.LikedByMe,
			Username:       v.Username,
			ProfilePicture: v.ProfilePicture,
		}

		data = append(data, response)
	}

	utils.WriteJson(w, http.StatusOK, utils.JsonResponse{
		Success: true,
		Data:    data,
	})
}

// Get posts from a friend
func (h *PostHandler) GetFriendPostsHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserID(r)
	if !ok {
		utils.UnauthorizedError(w, r, errors.New("Unauthorized."))
		return
	}

	friendID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	posts, err := h.postService.GetFriendPosts(r.Context(), friendID, userID)
	if err != nil {
		utils.NotFoundError(w, r, err)
		return
	}

	var data []types.PostResponse

	for _, v := range posts {
		response := types.PostResponse{
			ID:             v.ID,
			UserID:         v.UserID,
			Caption:        v.Caption,
			ImageURL:       v.ImageURL,
			ObjectKey:      v.ObjectKey,
			LikeCount:      v.LikeCount,
			ReplyCount:     v.ReplyCount,
			CreatedAt:      v.CreatedAt,
			LikedByMe:      v.LikedByMe,
			Username:       v.Username,
			ProfilePicture: v.ProfilePicture,
		}

		data = append(data, response)
	}

	utils.WriteJson(w, http.StatusOK, utils.JsonResponse{
		Success: true,
		Data:    data,
	})
}

func (h *PostHandler) generateObjectKey(userID int64, filename, folder string) string {
	return path.Join(folder, strconv.Itoa(int(userID)), filename+".jpeg")
}

func (h *PostHandler) publicImageURL(key string) string {
	return fmt.Sprintf("%s/%s", h.publicURL, key)
}

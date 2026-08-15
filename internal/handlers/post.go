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

// Delete existing post
func (h *PostHandler) DeletePostHadler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	if err := h.postService.DeletePost(r.Context(), id); err != nil {
		utils.NotFoundError(w, r, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.JsonResponse{
		Success: true,
	})
}

// TODO:
//
// Filter posts by users mutuals.Currently fetches all posts without
// filtering followings.
func (h *PostHandler) GetTimelineHandler(w http.ResponseWriter, r *http.Request) {
	// userID, ok := middlewares.GetUserID(r)
	// if !ok {
	// 	utils.UnauthorizedError(w, r, errors.New("Unauthorized."))
	// 	return
	// }

	posts, err := h.postService.GetLatestPosts(r.Context())
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	var data []types.PostResponse

	for _, v := range posts {
		response := types.PostResponse{
			ID:        v.ID,
			UserID:    v.UserID,
			ImageURL:  v.ImageURL,
			Caption:   v.Caption,
			ObjectKey: v.ObjectKey,
			CreatedAt: v.CreatedAt,
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
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	post, err := h.postService.GetPostByID(r.Context(), id)
	if err != nil {
		utils.NotFoundError(w, r, err)
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

	utils.WriteJson(w, http.StatusOK, utils.JsonResponse{
		Success: true,
		Data:    response,
	})
}

// Get current user posts, used for profile tab
func (h *PostHandler) GetMyPostsHandler(w http.ResponseWriter, r *http.Request) {
	userId, ok := middlewares.GetUserID(r)
	if !ok {
		utils.UnauthorizedError(w, r, errors.New("Unauthorized."))
		return
	}

	posts, err := h.postService.GetPostsByUserID(r.Context(), userId)
	if err != nil {
		utils.NotFoundError(w, r, err)
		return
	}

	var data []types.PostResponse

	for _, v := range posts {
		response := types.PostResponse{
			ID:        v.ID,
			UserID:    v.UserID,
			Caption:   v.Caption,
			ImageURL:  v.ImageURL,
			ObjectKey: v.ObjectKey,
			CreatedAt: v.CreatedAt,
		}

		data = append(data, response)
	}

	utils.WriteJson(w, http.StatusOK, utils.JsonResponse{
		Success: true,
		Data:    data,
	})
}

func (h *PostHandler) GetUsersPostsHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.ParseInt(chi.URLParam(r, "userId"), 10, 64)
	if err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	posts, err := h.postService.GetPostsByUserID(r.Context(), userId)
	if err != nil {
		utils.NotFoundError(w, r, err)
		return
	}

	var data []types.PostResponse

	for _, v := range posts {
		response := types.PostResponse{
			ID:        v.ID,
			UserID:    v.UserID,
			Caption:   v.Caption,
			ImageURL:  v.ImageURL,
			ObjectKey: v.ObjectKey,
			CreatedAt: v.CreatedAt,
		}

		data = append(data, response)
	}

	utils.WriteJson(w, http.StatusOK, utils.JsonResponse{
		Success: true,
		Data:    data,
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

func (h *PostHandler) generateObjectKey(userID int64, filename, folder string) string {
	return path.Join(folder, strconv.Itoa(int(userID)), filename+".jpeg")
}

func (h *PostHandler) publicImageURL(key string) string {
	return fmt.Sprintf("%s/%s", h.publicURL, key)
}

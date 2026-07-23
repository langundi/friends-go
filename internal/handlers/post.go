package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/langundi/friends-go/internal/middlewares"
	"github.com/langundi/friends-go/internal/services"
	"github.com/langundi/friends-go/internal/types"
	"github.com/langundi/friends-go/internal/utils"
)

type PostHandler struct {
	postService *services.PostService
}

func NewPostHandler(postService *services.PostService) *PostHandler {
	return &PostHandler{postService: postService}
}

func (h *PostHandler) NewPostHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserID(r)
	if !ok {
		utils.UnauthorizedError(w, r, errors.New("Unauthorized."))
		return
	}

	var req types.NewPostRequest
	req.UserID = userID

	if err := utils.ReadJson(w, r, &req); err != nil {
		utils.InvalidPayloadError(w, r, err)
		return
	}

	post, err := h.postService.NewPost(r.Context(), req)
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	response := types.PostResponse{
		ID:        post.ID,
		UserID:    post.UserID,
		Caption:   post.Caption,
		ImageURL:  post.ImageURL,
		CreatedAt: post.CreatedAt,
	}

	utils.WriteJson(w, http.StatusCreated, utils.JsonResponse{
		Success: true,
		Data:    response,
	})
}

func (h *PostHandler) GetPostHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.InternalServerError(w, r, err)
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
		CreatedAt: post.CreatedAt,
	}

	utils.WriteJson(w, http.StatusCreated, utils.JsonResponse{
		Success: true,
		Data:    response,
	})
}

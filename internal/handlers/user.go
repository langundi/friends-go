package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/langundi/friends-go/internal/middlewares"
	"github.com/langundi/friends-go/internal/services"
	"github.com/langundi/friends-go/internal/utils"
)

type UserHandler struct {
	userService *services.UserService
}

type UserResponse struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) ProfileHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserID(r)
	if !ok {
		utils.UnauthorizedError(w, r, errors.New("Unauthorized."))
		return
	}

	user, err := h.userService.GetUserProfile(r.Context(), userID)
	if err != nil {
		utils.InternalServerError(w, r, errors.New("User not found."))
		return
	}

	response := UserResponse{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
	}

	utils.WriteJson(w, http.StatusOK, utils.JsonResponse{
		Success: true,
		Data:    response,
	})
}

func (h *UserHandler) FriendsProfileHandler(w http.ResponseWriter, r *http.Request) {
	friendID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	user, err := h.userService.GetUserProfile(r.Context(), friendID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrUserNotFound):
			utils.BadRequestError(w, r, err)
		default:
			utils.InternalServerError(w, r, err)
		}
		return
	}

	response := UserResponse{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
	}

	utils.WriteJson(w, http.StatusOK, utils.JsonResponse{
		Success: true,
		Data:    response,
	})
}

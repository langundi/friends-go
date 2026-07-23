package handlers

import (
	"errors"
	"net/http"

	"github.com/langundi/friends-go/internal/services"
	"github.com/langundi/friends-go/internal/types"
	"github.com/langundi/friends-go/internal/utils"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req types.RegisterRequest
	if err := utils.ReadJson(w, r, &req); err != nil {
		utils.InvalidPayloadError(w, r, err)
		return
	}

	user, err := h.authService.RegisterUser(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrEmptyFields), errors.Is(err, services.ErrPasswordField):
			utils.BadRequestError(w, r, err)
		case errors.Is(err, services.ErrEmailExist):
			utils.ConflictError(w, r, err)
		default:
			utils.InternalServerError(w, r, err)
		}
		return
	}

	response := types.RegisterResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	utils.WriteJson(w, http.StatusCreated, utils.JsonResponse{
		Success: true,
		Data:    response,
	})
}

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req types.LoginRequest
	if err := utils.ReadJson(w, r, &req); err != nil {
		utils.InvalidPayloadError(w, r, err)
		return
	}

	access, refresh, err := h.authService.LoginUser(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrEmptyFields):
			utils.BadRequestError(w, r, err)
		case errors.Is(err, services.ErrInvalidCredentials):
			utils.UnauthorizedError(w, r, err)
		default:
			utils.InternalServerError(w, r, err)
		}
		return
	}

	response := types.LoginResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}

	utils.WriteJson(w, http.StatusCreated, utils.JsonResponse{
		Success: true,
		Data:    response,
	})
}

func (h *AuthHandler) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	var req types.RefreshRequest
	if err := utils.ReadJson(w, r, &req); err != nil {
		utils.InvalidPayloadError(w, r, err)
		return
	}

	token, err := h.authService.RefreshAccessToken(r.Context(), req.RefreshToken)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidToken), errors.Is(err, services.ErrExpiredToken):
			utils.UnauthorizedError(w, r, err)
		default:
			utils.InternalServerError(w, r, err)
		}
		return
	}

	response := types.RefreshResponse{
		AccessToken: token,
	}

	utils.WriteJson(w, http.StatusOK, utils.JsonResponse{
		Success: true,
		Data:    response,
	})
}

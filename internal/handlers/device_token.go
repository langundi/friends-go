package handlers

import (
	"net/http"

	"github.com/langundi/friends-go/internal/middlewares"
	"github.com/langundi/friends-go/internal/services"
	"github.com/langundi/friends-go/internal/types"
	"github.com/langundi/friends-go/internal/utils"
)

type DeviceTokenHandler struct {
	deviceTokenService *services.DeviceTokenService
}

func NewDeviceHandler(deviceTokenService *services.DeviceTokenService) *DeviceTokenHandler {
	return &DeviceTokenHandler{deviceTokenService: deviceTokenService}
}

func (h *DeviceTokenHandler) RegisterDeviceToken(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserID(r)
	if !ok {
		utils.UnauthorizedError(w, r, ErrUnauthorized)
		return
	}

	var req types.DeviceTokenRequest
	if err := utils.ReadJson(w, r, &req); err != nil {
		utils.InvalidPayloadError(w, r, err)
		return
	}

	err := h.deviceTokenService.RegisterDeviceToken(r.Context(), userID, req)
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	utils.WriteJson(w, http.StatusCreated, utils.JsonResponse{
		Success: true,
	})
}

func (h *DeviceTokenHandler) DeleteDeviceToken(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserID(r)
	if !ok {
		utils.UnauthorizedError(w, r, ErrUnauthorized)
		return
	}

	var req types.DeviceTokenRequest
	if err := utils.ReadJson(w, r, &req); err != nil {
		utils.InvalidPayloadError(w, r, err)
		return
	}

	err := h.deviceTokenService.DeleteDeviceToken(r.Context(), userID, req)
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	utils.WriteJson(w, http.StatusCreated, utils.JsonResponse{
		Success: true,
	})
}

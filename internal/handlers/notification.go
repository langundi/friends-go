package handlers

import (
	"net/http"

	"github.com/langundi/friends-go/internal/middlewares"
	"github.com/langundi/friends-go/internal/services"
	"github.com/langundi/friends-go/internal/types"
	"github.com/langundi/friends-go/internal/utils"
)

type NotificationHandler struct {
	notificationService *services.NotificationService
}

func NewNotificationHandler(notificationService *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{notificationService: notificationService}
}

func (h *NotificationHandler) GetAllNotifications(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserID(r)
	if !ok {
		utils.UnauthorizedError(w, r, ErrUnauthorized)
		return
	}

	notifications, err := h.notificationService.GetAllNotifications(r.Context(), userID)
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	var data []types.NotificationResponse
	for _, n := range notifications {
		response := types.NotificationResponse{
			ID:             n.ID,
			ReceiverID:     n.ReceiverID,
			SenderID:       n.SenderID,
			Message:        n.Message,
			PostID:         n.PostID,
			IsRead:         n.IsRead,
			CreatedAt:      n.CreatedAt,
			ProfilePicture: n.ProfilePicture,
		}
		data = append(data, response)
	}

	utils.WriteJson(w, http.StatusOK, utils.JsonResponse{
		Success: true,
		Data:    data,
	})
}

func (h *NotificationHandler) ReadNotifications(w http.ResponseWriter, r *http.Request) {
	if err := h.notificationService.ReadNotifications(r.Context()); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.JsonResponse{
		Success: true,
	})
}

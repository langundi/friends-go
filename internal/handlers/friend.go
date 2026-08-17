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

type FriendHandler struct {
	friendService *services.FriendService
}

var (
	ErrUnauthorized = errors.New("You are unauthorized.")
)

func NewFriendHandler(friendService *services.FriendService) *FriendHandler {
	return &FriendHandler{friendService: friendService}
}

func (h *FriendHandler) CreateFriendRequestHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserID(r)
	if !ok {
		utils.UnauthorizedError(w, r, ErrUnauthorized)
		return
	}

	receiverId, err := strconv.ParseInt(chi.URLParam(r, "receiverId"), 10, 64)
	if err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	friendRequest, err := h.friendService.CreateFriendRequest(r.Context(), userID, receiverId)
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	response := types.NewFriendRequestResponse{
		ID:         friendRequest.ID,
		SenderID:   friendRequest.SenderID,
		ReceiverID: friendRequest.ReceiverID,
		Status:     friendRequest.Status,
		CreatedAt:  friendRequest.CreatedAt,
	}

	utils.WriteJson(w, http.StatusCreated, utils.JsonResponse{
		Success: true,
		Data:    response,
	})
}

func (h *FriendHandler) GetFriendRequests(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserID(r)
	if !ok {
		utils.UnauthorizedError(w, r, ErrUnauthorized)
		return
	}

	list, err := h.friendService.GetFriendRequestsForUserID(r.Context(), userID)
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	var data []types.FriendRequestResponse

	for _, v := range list {
		response := types.FriendRequestResponse{
			ID:             v.ID,
			SenderID:       v.SenderID,
			SenderUsername: v.SenderUsername,
			ReceiverID:     v.ReceiverID,
			Status:         v.Status,
			CreatedAt:      v.CreatedAt,
		}

		data = append(data, response)
	}

	utils.WriteJson(w, http.StatusOK, utils.JsonResponse{
		Success: true,
		Data:    data,
	})
}

func (h *FriendHandler) DeclineFriendRequestHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	if err := h.friendService.DeleteFriendRequestByID(r.Context(), id); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.JsonResponse{
		Success: true,
	})
}

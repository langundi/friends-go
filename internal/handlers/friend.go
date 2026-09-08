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

// Send a new friend request
func (h *FriendHandler) SendFriendRequestHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserID(r)
	if !ok {
		utils.UnauthorizedError(w, r, ErrUnauthorized)
		return
	}

	var req types.SendFriendRequestNotification
	if err := utils.ReadJson(w, r, &req); err != nil {
		utils.InvalidPayloadError(w, r, err)
		return
	}

	friendRequest, err := h.friendService.SendFriendRequest(r.Context(), userID, req)
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

// Accept a friend request for current user
func (h *FriendHandler) AcceptFriendRequestHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	var req types.AcceptFriendRequestNotification
	if err := utils.ReadJson(w, r, &req); err != nil {
		utils.InvalidPayloadError(w, r, err)
		return
	}

	if err := h.friendService.AcceptFriendRequestByID(r.Context(), userID, req); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.JsonResponse{
		Success: true,
	})
}

// Fetch all friend request for a current user
func (h *FriendHandler) GetFriendRequestsHandler(w http.ResponseWriter, r *http.Request) {
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
			ProfilePicture: v.ProfilePicture,
		}

		data = append(data, response)
	}

	utils.WriteJson(w, http.StatusOK, utils.JsonResponse{
		Success: true,
		Data:    data,
	})
}

// Get friendship status for current and targeted user
func (h *FriendHandler) GetFriendshipStatusHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserID(r)
	if !ok {
		utils.UnauthorizedError(w, r, ErrUnauthorized)
		return
	}

	searchedUserID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	rel, err := h.friendService.GetFriendshipStatus(r.Context(), userID, searchedUserID)
	if err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	respone := types.FriendshipStatusResponse{
		FriendshipStatus: rel,
	}

	utils.WriteJson(w, http.StatusOK, utils.JsonResponse{
		Success: true,
		Data:    respone,
	})

}

// Decline a friend request for current user
func (h *FriendHandler) DeclineOrUnfriendFriendHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	if err := h.friendService.DeclineOrUnfriendFriendByID(r.Context(), userID); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.JsonResponse{
		Success: true,
	})
}

// Get friend list for current user
func (h *FriendHandler) GetMyFriendListHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserID(r)
	if !ok {
		utils.UnauthorizedError(w, r, ErrUnauthorized)
		return
	}

	list, err := h.friendService.GetMyFriendList(r.Context(), userID)
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	var data []types.FriendResponse

	for _, v := range list {
		response := types.FriendResponse{
			ID:             v.ID,
			UserID:         v.UserID,
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

// Get friend list for user
func (h *FriendHandler) GetFriendListHandler(w http.ResponseWriter, r *http.Request) {
	currentUserID, ok := middlewares.GetUserID(r)
	if !ok {
		utils.UnauthorizedError(w, r, ErrUnauthorized)
		return
	}

	userID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	list, err := h.friendService.GetFriendList(r.Context(), userID, currentUserID)
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	var data []types.FriendsFriendResponse
	for _, v := range list {
		response := types.FriendsFriendResponse{
			ID:             v.ID,
			UserID:         v.UserID,
			Username:       v.Username,
			ProfilePicture: v.ProfilePicture,
			FriendsWithMe:  v.FriendsWithMe,
		}
		data = append(data, response)
	}

	utils.WriteJson(w, http.StatusOK, utils.JsonResponse{
		Success: true,
		Data:    data,
	})
}

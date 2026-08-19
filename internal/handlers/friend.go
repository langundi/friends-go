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

// Create a new friend request
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
	currentUserID, ok := middlewares.GetUserID(r)
	if !ok {
		utils.UnauthorizedError(w, r, ErrUnauthorized)
		return
	}

	searchedUserID, err := strconv.ParseInt(chi.URLParam(r, "userId"), 10, 64)
	if err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	rel, err := h.friendService.GetFriendshipStatus(r.Context(), currentUserID, searchedUserID)
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
func (h *FriendHandler) DeclineFriendRequestHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	if err := h.friendService.DeclineFriendRequestByID(r.Context(), id); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.JsonResponse{
		Success: true,
	})
}

// Accept a friend request for current user
func (h *FriendHandler) AcceptFriendRequestHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	if err := h.friendService.AcceptFriendRequestByID(r.Context(), id); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.JsonResponse{
		Success: true,
	})
}

// Get friend list for current user
func (h *FriendHandler) GetFriendListHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserID(r)
	if !ok {
		utils.UnauthorizedError(w, r, ErrUnauthorized)
		return
	}

	list, err := h.friendService.GetFriendListForUserID(r.Context(), userID)
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	var data []types.UsernameResponse

	for _, v := range list {
		response := types.UsernameResponse{
			ID:       v.ID,
			Username: v.Username,
		}

		data = append(data, response)
	}

	utils.WriteJson(w, http.StatusOK, utils.JsonResponse{
		Success: true,
		Data:    data,
	})
}

package handlers

import (
	"net/http"

	"github.com/langundi/friends-go/internal/services"
	"github.com/langundi/friends-go/internal/types"
	"github.com/langundi/friends-go/internal/utils"
)

type TimelineHandler struct {
	postService *services.PostService
}

func NewTimelineHandler(postService *services.PostService) *TimelineHandler {
	return &TimelineHandler{postService: postService}
}

// TODO:
//
// Filter posts by users mutuals.Currently fetches all posts without
// filtering followings.
func (h *TimelineHandler) GetTimelineHandler(w http.ResponseWriter, r *http.Request) {
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
			CreatedAt: v.CreatedAt,
		}

		data = append(data, response)
	}

	utils.WriteJson(w, http.StatusOK, utils.JsonResponse{
		Success: true,
		Data:    data,
	})
}

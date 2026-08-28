package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/langundi/friends-go/internal/middlewares"
)

type HandlerConfig struct {
	*AuthHandler
	*UserHandler
	*PostHandler
	*FriendHandler
}

func Routes(h HandlerConfig) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.ClientIPFromRemoteAddr)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", h.RegisterHandler)
		r.Post("/login", h.LoginHandler)
		r.Post("/refresh", h.RefreshTokenHandler)

		r.Group(func(r chi.Router) {
			r.Use(middlewares.AuthMiddleware(h.authService))
			r.Post("/logout", h.LogoutHandler)
		})
	})

	r.Group(func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware(h.authService))

		r.Route("/user", func(r chi.Router) {
			r.Get("/", h.GetProfileHandler)
			r.Get("/{id}", h.GetFriendProfileHandler)
			r.Get("/post/me", h.GetMyPostsHandler)
			r.Get("/post/{id}", h.GetUsersPostsHandler)
			r.Get("/search/{username}", h.SearchProfileHandler)

			r.Patch("/change/username", h.ChangeUsername)
			r.Patch("/change/email", h.ChangeEmail)
			r.Patch("/change/password", h.ChangePassword)
		})

		r.Route("/post", func(r chi.Router) {
			r.Post("/", h.NewPostHandler)
			r.Post("/upload-image", h.GetPresignedURLHandler)
			r.Post("/like/{id}", h.LikePostHandler)
			r.Post("/reply/{id}", h.ReplyPostHandler)

			r.Get("/{id}", h.GetPostHandler)
			r.Get("/timeline", h.GetTimelineHandler)
			r.Get("/reply/{id}", h.GetPostRepliesHandler)

			r.Delete("/", h.DeletePostHadler)
			r.Delete("/unlike/{id}", h.UnlikePostHandler)
			r.Delete("/reply/{id}", h.DeleteReplyHandler)
		})

		r.Route("/friend-request", func(r chi.Router) {
			r.Post("/{id}", h.CreateFriendRequestHandler)

			r.Get("/", h.GetFriendRequestsHandler)

			r.Patch("/{id}", h.AcceptFriendRequestHandler)

			r.Delete("/{id}", h.DeclineOrUnfriendFriendHandler)
		})

		r.Route("/friend", func(r chi.Router) {
			r.Get("/", h.GetFriendListHandler)
			r.Get("/{id}/status", h.GetFriendshipStatusHandler)

			r.Delete("/{id}", h.DeclineOrUnfriendFriendHandler)
		})
	})

	return r
}

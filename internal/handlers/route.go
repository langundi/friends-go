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
			r.Get("/me/posts", h.GetMyPostsHandler)
			r.Get("/search/{username}", h.SearchProfileHandler)
			r.Post("/upload-image", h.ProfilePicturePresignedURLHandler)
			r.Delete("/delete", h.DeleteAccount)

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", h.GetFriendProfileHandler)
				r.Get("/posts", h.GetUsersPostsHandler)
			})

			r.Route("/change", func(r chi.Router) {
				r.Patch("/username", h.ChangeUsername)
				r.Patch("/email", h.ChangeEmail)
				r.Patch("/password", h.ChangePassword)
			})

			r.Route("/profile-picture", func(r chi.Router) {
				r.Patch("/", h.SetProfilePicture)
				r.Delete("/", h.DeleteProfilePictureHandler)
				r.Delete("/remove", h.RemoveProfilePictureHandler)
			})
		})

		r.Route("/post", func(r chi.Router) {
			r.Post("/", h.NewPostHandler)
			r.Post("/upload-image", h.GetPresignedURLHandler)
			r.Delete("/image/all", h.DeleteAllImageHandler)

			r.Route("/timeline", func(r chi.Router) {
				r.Get("/", h.GetTimelineHandler)
				r.Post("/more", h.GetMoreTimelineHandler)
			})

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", h.GetPostHandler)
				r.Get("/replies", h.GetPostRepliesHandler)
				r.Post("/like", h.LikePostHandler)
				r.Post("/reply", h.ReplyPostHandler)
				r.Delete("/", h.DeletePostHadler)
				r.Delete("/unlike", h.UnlikePostHandler)
				r.Delete("/reply", h.DeleteReplyHandler)
			})
		})

		r.Route("/friend-request", func(r chi.Router) {
			r.Get("/", h.GetFriendRequestsHandler)
			r.Post("/{id}", h.CreateFriendRequestHandler)
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

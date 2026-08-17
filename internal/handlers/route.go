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
			r.Get("/search/{username}", h.SearchProfileHandler)
		})

		r.Route("/post", func(r chi.Router) {
			r.Post("/", h.NewPostHandler)
			r.Post("/upload", h.GetPresignedURLHandler)

			r.Get("/{id}", h.GetPostHandler)
			r.Get("/timeline", h.GetTimelineHandler)
			r.Get("/user/me", h.GetMyPostsHandler)
			r.Get("/user/{userId}", h.GetUsersPostsHandler)

			r.Delete("/", h.DeletePostHadler)
		})
	})

	return r
}

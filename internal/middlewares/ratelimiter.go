package middlewares

import (
	"net"
	"net/http"

	"github.com/langundi/friends-go/internal/ratelimiter"
	"github.com/langundi/friends-go/internal/utils"
)

func RateLimiterMiddleware(config ratelimiter.Config, rl *ratelimiter.RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				ip = r.RemoteAddr
			}
			if config.Enabled {
				if allow, retryAfter := rl.Allow(ip); !allow {
					utils.RateLimitExceed(w, r, retryAfter.String())
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

package middleware

import (
	"errors"
	"net/http"

	"github.com/robsonrg/goexpert-desafio-tecnico-rate-limiter/internal/infra/ratelimit"
)

const (
	tooManyRequestsErr = "you have reached the maximum number of requests or actions allowed within a certain time frame"
)

type WebServerRateLimiter struct {
	ratelimit *ratelimit.RateLimit
}

func NewWebServerRateLimiter(ratelimit *ratelimit.RateLimit) *WebServerRateLimiter {
	return &WebServerRateLimiter{
		ratelimit: ratelimit,
	}
}

func (m *WebServerRateLimiter) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allow, err := m.ratelimit.Allow(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !allow {
			http.Error(w, errors.New(tooManyRequestsErr).Error(), http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

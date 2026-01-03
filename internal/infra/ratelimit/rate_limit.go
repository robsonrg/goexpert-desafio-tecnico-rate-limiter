package ratelimit

import (
	"context"
	"errors"
	"net"
	"net/http"
	"strings"
)

type Limiter interface {
	Quota(ctx context.Context, key string) (int64, error)
}

type RateLimit struct {
	limiter       Limiter
	limitPerIP    int64
	limitPerToken int64
}

func NewRateLimit(limiter Limiter, limitPerIP, limitPerToken int64) *RateLimit {
	return &RateLimit{
		limiter:       limiter,
		limitPerIP:    limitPerIP,
		limitPerToken: limitPerToken,
	}
}

func (rl *RateLimit) Allow(r *http.Request) (bool, error) {
	key := rl.buildLimiterKey(r)
	quota, err := rl.limiter.Quota(r.Context(), key)
	if err != nil {
		return false, err
	}
	apiKey := r.Header.Get("API_KEY")
	if apiKey != "" {
		return quota <= rl.limitPerToken, nil
	}
	return quota <= rl.limitPerIP, nil
}

func (rl *RateLimit) buildLimiterKey(r *http.Request) string {
	ip, _ := rl.extractRequestIP(r)
	return ip
}

func (rl *RateLimit) extractRequestIP(r *http.Request) (string, error) {
	ips := r.Header.Get("X-Forwarded-For")
	splitIps := strings.Split(ips, ",")

	if len(splitIps) > 0 {
		netIP := net.ParseIP(splitIps[len(splitIps)-1])
		if netIP != nil {
			return netIP.String(), nil
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}

	netIP := net.ParseIP(ip)
	if netIP != nil {
		ip := netIP.String()
		if ip == "::1" {
			return "127.0.0.1", nil
		}
		return ip, nil
	}

	return "", errors.New("IP not found")
}

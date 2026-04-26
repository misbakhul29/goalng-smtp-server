package middleware

import (
	"net/http"
	"sync"
	"time"
)

type visitor struct {
	count    int
	windowAt time.Time
}

type RateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		limit:    limit,
		window:   window,
	}
	go rl.cleanupLoop()
	return rl
}

func (rl *RateLimiter) Limit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := realIP(r)

		rl.mu.Lock()
		v, ok := rl.visitors[ip]
		if !ok || time.Since(v.windowAt) > rl.window {
			rl.visitors[ip] = &visitor{count: 1, windowAt: time.Now()}
			rl.mu.Unlock()
			next(w, r)
			return
		}

		v.count++
		if v.count > rl.limit {
			rl.mu.Unlock()
			http.Error(w, `{"error":"too many requests, please slow down"}`, http.StatusTooManyRequests)
			return
		}
		rl.mu.Unlock()

		next(w, r)
	}
}

func (rl *RateLimiter) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.windowAt) > rl.window {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func realIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}
	if ip := r.Header.Get("X-Real-Ip"); ip != "" {
		return ip
	}
	return r.RemoteAddr
}

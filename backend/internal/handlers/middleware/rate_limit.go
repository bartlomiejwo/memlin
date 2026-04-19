package middleware

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/time/rate"
)

// Visitor tracks a rate limiter and its last access time
type Visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// RateLimiter limits requests per IP and/or user
type RateLimiter struct {
	visitors    map[string]*Visitor
	users       map[string]*Visitor
	mu          sync.Mutex
	ipRate      rate.Limit // Rate for IP-based limiting
	ipBurst     int        // Burst for IP-based limiting
	userRate    rate.Limit // Rate for user-based limiting
	userBurst   int        // Burst for user-based limiting
	jwtSecret   string     // For parsing JWT to get user ID
	useUserRate bool       // Toggle user-based limiting
}

// NewRateLimiter creates a new rate limiter and starts cleanup
func NewRateLimiter(ipRate rate.Limit, ipBurst int, userRate rate.Limit, userBurst int, jwtSecret string, useUserRate bool) *RateLimiter {
	rl := &RateLimiter{
		visitors:    make(map[string]*Visitor),
		users:       make(map[string]*Visitor),
		ipRate:      ipRate,
		ipBurst:     ipBurst,
		userRate:    userRate,
		userBurst:   userBurst,
		jwtSecret:   jwtSecret,
		useUserRate: useUserRate,
	}
	go rl.cleanupVisitors()
	return rl
}

// GetLimiter retrieves or creates a rate limiter for an IP or user
func (rl *RateLimiter) GetLimiter(key string, isUser bool) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	targetMap := rl.visitors
	rateLimit := rl.ipRate
	burst := rl.ipBurst
	if isUser {
		targetMap = rl.users
		rateLimit = rl.userRate
		burst = rl.userBurst
	}

	if visitor, exists := targetMap[key]; exists {
		visitor.lastSeen = time.Now()
		return visitor.limiter
	}

	limiter := rate.NewLimiter(rateLimit, burst)
	targetMap[key] = &Visitor{limiter, time.Now()}
	return limiter
}

// cleanupVisitors removes inactive entries
func (rl *RateLimiter) cleanupVisitors() {
	for {
		time.Sleep(10 * time.Minute)
		rl.mu.Lock()
		for ip, visitor := range rl.visitors {
			if time.Since(visitor.lastSeen) > 1*time.Hour {
				delete(rl.visitors, ip)
			}
		}
		for user, visitor := range rl.users {
			if time.Since(visitor.lastSeen) > 1*time.Hour {
				delete(rl.users, user)
			}
		}
		rl.mu.Unlock()
	}
}

// RateLimitMiddleware applies rate limiting
func RateLimitMiddleware(rl *RateLimiter, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get IP
		ip := r.Header.Get("X-Forwarded-For")
		if ip == "" {
			host, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				host = r.RemoteAddr
			}
			ip = host
		} else {
			ip = strings.TrimSpace(strings.Split(ip, ",")[0])
		}

		// Always enforce IP-based limit
		if !rl.GetLimiter(ip, false).Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		// Check user-based limit if enabled and authenticated
		if rl.useUserRate {
			authHeader := r.Header.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
				token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
					return []byte(rl.jwtSecret), nil
				})
				if err == nil && token.Valid {
					if claims, ok := token.Claims.(jwt.MapClaims); ok {
						if userID, ok := claims["sub"].(string); ok {
							if !rl.GetLimiter(userID, true).Allow() {
								http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
								return
							}
						}
					}
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}

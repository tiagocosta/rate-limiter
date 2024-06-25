package limiter

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

const (
	HEADER_IP    = "IP"
	HEADER_TOKEN = "API_KEY"
)

type Message struct {
	Status string `json:"status"`
	Body   string `json:"body"`
}

type Middleware struct {
	repository RepositoryInterface
}

func NewMiddleware(repository RepositoryInterface) (*Middleware, error) {
	if repository == nil {
		return nil, fmt.Errorf("repository can't be nil")
	}

	return &Middleware{
		repository: repository,
	}, nil
}

func (m *Middleware) Limit(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		token := r.Header.Get(HEADER_TOKEN)

		if m.isLimitedByTokenOnly() {
			if isTokenValid(token) {
				if !m.allowToken(ctx, token) {
					m.writeErrorMessage(w)
					return
				}
			} else {
				m.writeErrorMessage(w)
				return
			}
		} else if m.isLimitedByIPOnly() {
			if !m.allowIP(ctx, ip) {
				m.writeErrorMessage(w)
				return
			}
		} else if m.isLimitedByTokenAndIP() {
			if isTokenValid(token) {
				if !m.allowToken(ctx, token) {
					if !m.allowIP(ctx, ip) {
						m.writeErrorMessage(w)
						return
					}
				}
			} else {
				if !m.allowIP(ctx, ip) {
					m.writeErrorMessage(w)
					return
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}

func isLimitedByIP() bool {
	return slices.Contains(strings.Split(os.Getenv("LIMIT_BY"), ","), HEADER_IP)
}

func isLimitedByToken() bool {
	return slices.Contains(strings.Split(os.Getenv("LIMIT_BY"), ","), HEADER_TOKEN)
}

func (m *Middleware) allowIP(ctx context.Context, ip string) bool {
	expired, _ := m.repository.IsExpired(ctx, ip)
	if expired {
		return false
	}
	consumed, err := m.repository.Get(ctx, ip)
	if err != nil {
		m.repository.Set(ctx, ip, 0, getIPWindow())
	}
	if consumed < getIPLimit() {
		m.repository.Increment(ctx, ip)
		consumed, _ = m.repository.Get(ctx, ip)
		if consumed == getIPLimit() {
			m.repository.SetExpired(ctx, ip, getIPExpiration())
		}
		return true
	}
	return false
}

func (m *Middleware) allowToken(ctx context.Context, token string) bool {
	fmt.Println(token)
	expired, _ := m.repository.IsExpired(ctx, token)
	if expired {
		return false
	}
	consumed, err := m.repository.Get(ctx, token)
	if err != nil {
		m.repository.Set(ctx, token, 0, getTokenWindow())
	}
	if consumed < getTokenLimit(token) {
		m.repository.Increment(ctx, token)
		consumed, _ = m.repository.Get(ctx, token)
		if consumed == getTokenLimit(token) {
			m.repository.SetExpired(ctx, token, getTokenExpiration(token))
		}
		return true
	}
	return false
}

func (m *Middleware) isLimitedByTokenAndIP() bool {
	return isLimitedByToken() && isLimitedByIP()
}

func (m *Middleware) isLimitedByTokenOnly() bool {
	return isLimitedByToken() && !isLimitedByIP()
}

func (m *Middleware) isLimitedByIPOnly() bool {
	return !isLimitedByToken() && isLimitedByIP()
}

func (m *Middleware) writeErrorMessage(w http.ResponseWriter) {
	message := Message{
		Status: "Request Failed",
		Body:   "you have reached the maximum number of requests or actions allowed within a certain time frame",
	}
	w.WriteHeader(http.StatusTooManyRequests)
	json.NewEncoder(w).Encode(&message)
}

func getTokens() []string {
	return strings.Split(os.Getenv("TOKENS"), ",")
}

func isTokenValid(token string) bool {
	return slices.Contains(getTokens(), token)
}

func getTokenLimit(token string) int {
	tokensLimit := strings.Split(os.Getenv("TOKENS_LIMIT"), ",")
	limit, err := strconv.Atoi(tokensLimit[slices.Index(getTokens(), token)])
	if err != nil {
		panic(err)
	}
	return limit
}

func getTokenExpiration(token string) time.Duration {
	tokensExpiration := strings.Split(os.Getenv("TOKENS_EXPIRATION_INTERVAL"), ",")
	expiration, err := time.ParseDuration(tokensExpiration[slices.Index(getTokens(), token)])
	if err != nil {
		panic(err)
	}
	return expiration
}

func getTokenWindow() time.Duration {
	window, err := time.ParseDuration(os.Getenv("TOKEN_WINDOW_INTERVAL"))
	if err != nil {
		panic(err)
	}
	return window
}

func getIPLimit() int {
	limit, err := strconv.Atoi(os.Getenv("IP_LIMIT"))
	if err != nil {
		panic(err)
	}
	return limit
}

func getIPExpiration() time.Duration {
	expiration, err := time.ParseDuration(os.Getenv("IP_EXPIRATION_INTERVAL"))
	if err != nil {
		panic(err)
	}
	return expiration
}

func getIPWindow() time.Duration {
	window, err := time.ParseDuration(os.Getenv("IP_WINDOW_INTERVAL"))
	if err != nil {
		panic(err)
	}
	return window
}

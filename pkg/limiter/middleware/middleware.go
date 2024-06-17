package middleware

// import (
// 	"encoding/json"
// 	"net"
// 	"net/http"
// 	"os"
// 	"strconv"
// 	"sync"
// 	"time"

// 	"github.com/tiagocosta/rate-limiter/pkg/limiter"
// 	"golang.org/x/time/rate"
// )

// var (
// 	mu                        sync.Mutex
// 	clients                   = make(map[string]*client)
// 	ipLimit, _                = strconv.Atoi(os.Getenv("IP_LIMIT"))
// 	tokenLimit, _             = strconv.Atoi(os.Getenv("TOKEN_LIMIT"))
// 	ipInterval, _             = time.ParseDuration(os.Getenv("IP_INTERVAL"))
// 	tokenInterval, _          = time.ParseDuration(os.Getenv("TOKEN_INTERVAL"))
// 	ipRejectionInterval, _    = time.ParseDuration(os.Getenv("IP_REJECTION_INTERVAL"))
// 	tokenRejectionInterval, _ = time.ParseDuration(os.Getenv("TOKEN_REJECTION_INTERVAL"))
// )

// func RateLimiter(next func(writer http.ResponseWriter, request *http.Request)) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		ip, _, err := net.SplitHostPort(r.RemoteAddr)
// 		if err != nil {
// 			w.WriteHeader(http.StatusInternalServerError)
// 			return
// 		}
// 		mu.Lock()
// 		if _, found := clients[ip]; !found {
// 			client := &client{
// 				limiterIP:    limiter.NewLimiter(ipLimit, ipInterval, "", ""),
// 				limiterToken: rate.NewLimiter(rate.Limit(tokenLimit), tokenLimit),
// 				ip:           ip,
// 			}

// 			clients[ip] = client
// 		}
// 		if !clients[ip].limiterIP.Allow() {
// 			mu.Unlock()

// 			message := message{
// 				Status: "Request Failed",
// 				Body:   "you have reached the maximum number of requests or actions allowed within a certain time frame",
// 			}

// 			w.WriteHeader(http.StatusTooManyRequests)
// 			json.NewEncoder(w).Encode(&message)
// 			return
// 		}
// 		mu.Unlock()
// 		next(w, r)
// 	})
// }

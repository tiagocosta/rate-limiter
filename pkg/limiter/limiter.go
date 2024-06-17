package limiter

// import (
// 	"encoding/json"
// 	"net"
// 	"net/http"
// 	"os"
// 	"strconv"
// 	"sync"
// 	"time"

// 	"golang.org/x/time/rate"
// )

// type message struct {
// 	Status string `json:"status"`
// 	Body   string `json:"body"`
// }

// type client struct {
// 	ip           string
// 	token        string
// 	limiterIP    *rate.Limiter
// 	limiterToken *rate.Limiter

// 	TimerIP    *time.Timer
// 	TimerToken *time.Timer

// 	ipBlockedAt    *time.Time
// 	tokenBlockedAt *time.Time
// }

// var (
// 	mu                    sync.Mutex
// 	clients               = make(map[string]*client)
// 	rateIP, _             = strconv.Atoi(os.Getenv("RATE_IP"))
// 	rateToken, _          = strconv.Atoi(os.Getenv("RATE_TOKEN"))
// 	blockIntervalIP, _    = time.ParseDuration(os.Getenv("BOLCK_INTERVAL_IP"))
// 	blockIntervalToken, _ = time.ParseDuration(os.Getenv("BOLCK_INTERVAL_TOKEN"))
// )

// func RateLimiter(next func(writer http.ResponseWriter, request *http.Request)) http.Handler {
// 	// go func() {
// 	// 	for {
// 	// 		time.Sleep(time.Second)
// 	// 		mu.Lock()
// 	// 		for _, client := range clients {
// 	// 			if client.ipBlockedAt != nil {
// 	// 				if time.Since(*client.ipBlockedAt) >= blockIntervalIP {
// 	// 					client.ipBlockedAt = nil
// 	// 				}
// 	// 			}
// 	// 			if client.tokenBlockedAt != nil {
// 	// 				if time.Since(*client.tokenBlockedAt) >= blockIntervalToken {
// 	// 					client.tokenBlockedAt = nil
// 	// 				}
// 	// 			}
// 	// 		}
// 	// 		mu.Unlock()
// 	// 	}
// 	// }()
// 	startExpirationVerification()

// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		ip, _, err := net.SplitHostPort(r.RemoteAddr)
// 		if err != nil {
// 			w.WriteHeader(http.StatusInternalServerError)
// 			return
// 		}
// 		mu.Lock()
// 		if _, found := clients[ip]; !found {
// 			client := &client{
// 				limiterIP:    rate.NewLimiter(rate.Limit(rateIP), rateIP),
// 				limiterToken: rate.NewLimiter(rate.Limit(rateToken), rateToken),
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

// func startExpirationVerification() {
// 	go func() {
// 		for {
// 			time.Sleep(time.Second)
// 			mu.Lock()

// 			for ip, client := range clients {
// 				if time.Since(client.lastSeen) > 3*time.Minute {
// 					delete(clients, ip)
// 				}
// 			}
// 			for ip, client := range clients {
// 				if client.ipBlockedAt != nil {
// 					if time.Since(*client.ipBlockedAt) >= blockIntervalIP {
// 						client.ipBlockedAt = nil
// 					}
// 				}
// 				if client.tokenBlockedAt != nil {
// 					if time.Since(*client.tokenBlockedAt) >= blockIntervalToken {
// 						client.tokenBlockedAt = nil
// 					}
// 				}
// 			}

// 			mu.Unlock()
// 		}
// 	}()
// }

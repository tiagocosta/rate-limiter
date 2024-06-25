package main

// import (
// 	"encoding/json"
// 	"log"
// 	"net/http"

// 	"github.com/joho/godotenv"
// 	"github.com/tiagocosta/rate-limiter/pkg/limiter"
// )

// type Message struct {
// 	Status string `json:"status"`
// 	Body   string `json:"body"`
// }

// func main() {
// 	if err := godotenv.Load("cmd/.env"); err != nil {
// 		log.Fatal("error trying to load env variables")
// 		return
// 	}

// 	http.Handle("/ping", limiter.RateLimiter(endpointHandler))
// 	err := http.ListenAndServe(":8080", nil)
// 	if err != nil {
// 		log.Println("There was an error listening on port :8080", err)
// 	}
// }

// func endpointHandler(writer http.ResponseWriter, request *http.Request) {
// 	writer.Header().Set("Content-Type", "application/json")
// 	writer.WriteHeader(http.StatusOK)
// 	message := Message{
// 		Status: "Successful",
// 		Body:   "api reached",
// 	}
// 	err := json.NewEncoder(writer).Encode(&message)
// 	if err != nil {
// 		return
// 	}
// }

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/tiagocosta/rate-limiter/pkg/infra/redis_cache"
	"github.com/tiagocosta/rate-limiter/pkg/limiter"
)

type Message struct {
	Status string `json:"status"`
	Body   string `json:"body"`
}

func endpointHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	message := Message{
		Status: "Successful",
		Body:   "API reached!",
	}
	err := json.NewEncoder(w).Encode(&message)
	if err != nil {
		return
	}
}

func main() {
	err := godotenv.Load("cmd/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	repo := redis_cache.NewRedisRepository(rdb)

	rateLimiter, err := limiter.NewMiddleware(repo)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/ping", endpointHandler)

	// mux := http.NewServeMux()

	// mux.HandleFunc("/ping", endpointHandler)

	err = http.ListenAndServe(":8080", rateLimiter.Limit(endpointHandler))
	if err != nil {
		log.Println("There was an error listening on port :8080", err)
	}
}

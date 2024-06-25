package limiter

// https://www.inanzzz.com/index.php/post/xgod/testing-a-middleware-within-golang
// https://github.com/go-redis/redismock/tree/master
// https://elliotchance.medium.com/mocking-redis-in-unit-tests-in-go-28aff285b98

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/tiagocosta/rate-limiter/pkg/infra/redis_cache"
)

type MiddlewareTestSuite struct {
	suite.Suite
	middleware *Middleware
	redisC     testcontainers.Container
	client     *redis.Client
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

func (suite *MiddlewareTestSuite) SetupSuite() {
	err := godotenv.Load("../../cmd/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "redis:alpine",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}
	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("Could not start redis: %s", err)
	}
	endpoint, err := redisC.Endpoint(ctx, "")
	if err != nil {
		log.Fatal(err)
	}
	client := redis.NewClient(&redis.Options{
		Addr: endpoint,
	})
	repo := redis_cache.NewRedisRepository(client)
	suite.client = client
	suite.redisC = redisC
	rateLimiter, err := NewMiddleware(repo)
	if err != nil {
		panic(err)
	}
	suite.middleware = rateLimiter
}

func TestMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(MiddlewareTestSuite))
}

func (suite *MiddlewareTestSuite) TearDownSuite() {
	if err := suite.redisC.Terminate(context.Background()); err != nil {
		log.Fatalf("Could not stop redis: %s", err)
	}
}

func (suite *MiddlewareTestSuite) TearDownTest() {
	err := suite.client.FlushAll(context.Background()).Err()
	if err != nil {
		log.Fatalf("Could not flush redis: %s", err)
	}
}

func (suite *MiddlewareTestSuite) TestLimitByIPOnly() {
	os.Setenv("LIMIT_BY", "IP")

	// Create test HTTP server
	ts := httptest.NewServer(suite.middleware.Limit(endpointHandler))
	defer ts.Close()

	for range getIPLimit() {
		resp, _ := http.Get(fmt.Sprintf("%s/", ts.URL))
		if resp.StatusCode != http.StatusOK {
			suite.T().Errorf("status code = %d, want %d", resp.StatusCode, http.StatusOK)
		}
	}

	resp, _ := http.Get(fmt.Sprintf("%s/", ts.URL))
	if resp.StatusCode != http.StatusTooManyRequests {
		suite.T().Errorf("status code = %d, want %d", resp.StatusCode, http.StatusTooManyRequests)
	}

	time.Sleep(getIPExpiration())

	resp, _ = http.Get(fmt.Sprintf("%s/", ts.URL))
	if resp.StatusCode != http.StatusOK {
		suite.T().Errorf("status code = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

func (suite *MiddlewareTestSuite) TestLimitByTokenOnly() {
	os.Setenv("LIMIT_BY", "API_KEY")
	// Create test HTTP server
	ts := httptest.NewServer(suite.middleware.Limit(endpointHandler))
	defer ts.Close()

	token := "abc123"
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/", ts.URL), nil)
	client := &http.Client{}
	for range getTokenLimit(token) {
		req.Header.Set(HEADER_TOKEN, token)
		resp, _ := client.Do(req)
		if resp.StatusCode != http.StatusOK {
			suite.T().Errorf("status code = %d, want %d", resp.StatusCode, http.StatusOK)
		}
	}

	resp, _ := client.Do(req)
	if resp.StatusCode != http.StatusTooManyRequests {
		suite.T().Errorf("status code = %d, want %d", resp.StatusCode, http.StatusTooManyRequests)
	}

	time.Sleep(getTokenExpiration(token) + getIPExpiration())

	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusOK {
		suite.T().Errorf("status code = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

func (suite *MiddlewareTestSuite) TestLimitByTokenAndIP() {
	os.Setenv("LIMIT_BY", "IP,API_KEY")
	// Create test HTTP server
	ts := httptest.NewServer(suite.middleware.Limit(endpointHandler))
	defer ts.Close()

	token := "def456"
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/", ts.URL), nil)
	client := &http.Client{}
	for range getTokenLimit(token) + getIPLimit() {
		req.Header.Set(HEADER_TOKEN, token)
		resp, _ := client.Do(req)
		if resp.StatusCode != http.StatusOK {
			suite.T().Errorf("status code = %d, want %d", resp.StatusCode, http.StatusOK)
		}
	}

	resp, _ := client.Do(req)
	if resp.StatusCode != http.StatusTooManyRequests {
		suite.T().Errorf("status code = %d, want %d", resp.StatusCode, http.StatusTooManyRequests)
	}

	time.Sleep(getTokenExpiration(token) + getIPExpiration())

	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusOK {
		suite.T().Errorf("status code = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

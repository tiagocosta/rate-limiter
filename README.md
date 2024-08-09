# rate-limiter
Rate limiting is the process of controlling how many requests your application users can make within a specified timeframe. \
This project provides a rate limiter written in Go and uses Redis for caching. It is configured as a middleware that intercepts all requests, goes to redis cache and verifies if requests should be allowed to continue or blocked. It supports rate limiting by IP and/or Token. For the sake of simplicity, all configurations are in .env file.

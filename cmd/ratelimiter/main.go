package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/robsonrg/goexpert-desafio-tecnico-rate-limiter/internal/infra/configs"
	"github.com/robsonrg/goexpert-desafio-tecnico-rate-limiter/internal/infra/ratelimit"
	"github.com/robsonrg/goexpert-desafio-tecnico-rate-limiter/internal/infra/ratelimit/limiter"
	"github.com/robsonrg/goexpert-desafio-tecnico-rate-limiter/internal/infra/web"
	"github.com/robsonrg/goexpert-desafio-tecnico-rate-limiter/internal/infra/web/webserver"
	"github.com/robsonrg/goexpert-desafio-tecnico-rate-limiter/internal/infra/web/webserver/middleware"
)

func main() {
	done := make(chan struct{})

	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	rdb, err := redisClient(configs.RedisHost, configs.RedisPassword)
	if err != nil {
		panic(err)
	}
	defer rdb.Close()

	webserver := webserver.NewWebServer(":8080")
	webHandler := web.NewWebHandler()

	limiter := limiter.NewRedisLimiter(rdb)
	ratelimit := ratelimit.NewRateLimit(limiter, configs.RateLimitRequestPerSecondsIP, configs.RateLimitRequestPerSecondsToken)

	rl := middleware.NewWebServerRateLimiter(ratelimit)
	webserver.AddMiddleware("rate_limiter", rl.Handle)

	webserver.AddHandler("/sample", webHandler.Handle)
	fmt.Println("Starting web server on port", 8080)

	webserver.Start()
	<-done
}

func redisClient(host, password string) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       0,
	})
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}
	return rdb, nil
}

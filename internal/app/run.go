package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/C-dexTeam/codex-compiler/internal/config"
	"github.com/C-dexTeam/codex-compiler/internal/http"
	"github.com/C-dexTeam/codex-compiler/internal/http/middlewares"
	"github.com/C-dexTeam/codex-compiler/internal/http/response"
	"github.com/C-dexTeam/codex-compiler/internal/http/server"
	"github.com/redis/go-redis/v9"
)

func Run(cfg *config.Config) {
	// Redis Client Başlat
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// Redis Connection Test
	err := redisClient.Ping(context.Background()).Err()
	if err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
	}
	fmt.Println("Redis connection established successfully.")

	// Handler Initialize
	handlers := http.NewHandler(redisClient)

	// Fiber İnitialize
	fiberServer := server.NewServer(cfg, response.ResponseHandler)

	// Captcha Initialize
	go func() {
		err := fiberServer.Run(handlers.Init(cfg.Application.DevMode, middlewares.InitMiddlewares(cfg)...))
		if err != nil {
			log.Fatalf("Error while running fiber server: %v", err.Error())
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	fmt.Println("Gracefully shutting down...")
	_ = fiberServer.Shutdown(context.Background())
	fmt.Println("Fiber was successful shutdown.")
}

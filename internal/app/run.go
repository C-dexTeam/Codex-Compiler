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
	"github.com/C-dexTeam/codex-compiler/internal/services"
	validatorService "github.com/C-dexTeam/codex-compiler/pkg/validator_service"
)

func Run(cfg *config.Config) {
	// Utilities Initialize
	validatorService := validatorService.NewValidatorService()

	// Service Initialize
	allServices := services.CreateNewServices(
		validatorService,
	)

	// Handler Initialize
	handlers := http.NewHandler(allServices)

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

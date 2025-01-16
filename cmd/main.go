package main

import (
	"github.com/C-dexTeam/codex-compiler/internal/app"
	"github.com/C-dexTeam/codex-compiler/internal/config"
)

// @title API Service
// @description API Service for Codex-Compiler
// @host localhost
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in cookie
// @name session_id
func main() {
	// Setting all the configs and starting the app.
	cfg, err := config.Init("./config")
	if err != nil {
		panic(err)
	}
	app.Run(cfg)
}

package main

import (
	"github.com/C-dexTeam/codex-compiler/internal/app"
	"github.com/C-dexTeam/codex-compiler/internal/config"
)

func main() {
	cfg, err := config.Init("./config")
	if err != nil {
		panic(err)
	}
	app.Run(cfg)
}

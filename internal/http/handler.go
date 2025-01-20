package http

import (
	"github.com/C-dexTeam/codex-compiler/docs"
	"github.com/C-dexTeam/codex-compiler/internal/config"
	dto "github.com/C-dexTeam/codex-compiler/internal/http/dtos"
	"github.com/C-dexTeam/codex-compiler/internal/http/sessionStore"
	v1 "github.com/C-dexTeam/codex-compiler/internal/http/v1"
	"github.com/redis/go-redis/v9"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

type Handler struct {
	redis *redis.Client
}

func NewHandler(redis *redis.Client) *Handler {
	return &Handler{
		redis: redis,
	}
}

func (h *Handler) Init(devMode bool, middlewares ...func(*fiber.Ctx) error) *fiber.App {
	app := fiber.New()
	for i := range middlewares {
		app.Use(middlewares[i])
	}

	if devMode {
		docs.SwaggerInfo.Version = config.Version
		app.Get("/compiler-api/dev/*", swagger.New(swagger.Config{
			Title:                "Codex-Compiler Backend",
			TryItOutEnabled:      true,
			PersistAuthorization: true,
		}))
	}

	root := app.Group("/compiler-api")
	sessionStore := sessionStore.NewSessionStore()
	dtoManager := dto.CreateNewDTOManager()

	// init routes
	v1.NewV1Handler(dtoManager, h.redis).Init(root, sessionStore)

	return app
}

package http

import (
	"github.com/C-dexTeam/codex-compiler/docs"
	"github.com/C-dexTeam/codex-compiler/internal/config"
	dto "github.com/C-dexTeam/codex-compiler/internal/http/dtos"
	"github.com/C-dexTeam/codex-compiler/internal/http/sessionStore"
	v1 "github.com/C-dexTeam/codex-compiler/internal/http/v1"
	"github.com/C-dexTeam/codex-compiler/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

type Handler struct {
	services *services.Services
}

func NewHandler(services *services.Services) *Handler {
	return &Handler{
		services: services,
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
	v1.NewV1Handler(dtoManager, h.services).Init(root, sessionStore)

	return app
}

package v1

import (
	dto "github.com/C-dexTeam/codex-compiler/internal/http/dtos"
	"github.com/C-dexTeam/codex-compiler/internal/http/response"
	"github.com/C-dexTeam/codex-compiler/internal/http/v1/private"
	"github.com/redis/go-redis/v9"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type V1Handler struct {
	redis      *redis.Client
	dtoManager dto.IDTOManager
}

func NewV1Handler(dtoManager dto.IDTOManager, redis *redis.Client) *V1Handler {
	return &V1Handler{
		dtoManager: dtoManager,
	}
}

func (h *V1Handler) Init(router fiber.Router, sessionStore *session.Store) {
	root := router.Group("/v1")
	root.Get("/", func(c *fiber.Ctx) error {
		return response.Response(200, "Welcome to Codex-Compiler API (Root Zone)", nil)
	})

	// Init Fiber Session Store
	private.NewPrivateHandler(sessionStore, h.dtoManager, h.redis).Init(root)
}

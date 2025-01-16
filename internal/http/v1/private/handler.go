package private

import (
	"fmt"

	serviceErrors "github.com/C-dexTeam/codex-compiler/internal/errors"
	dto "github.com/C-dexTeam/codex-compiler/internal/http/dtos"
	"github.com/C-dexTeam/codex-compiler/internal/http/response"
	"github.com/C-dexTeam/codex-compiler/internal/http/sessionStore"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
)

// PrivateHandler struct
type PrivateHandler struct {
	sess_store *session.Store
	dtoManager dto.IDTOManager
	redis      *redis.Client
}

// NewPrivateHandler creates a new instance of PrivateHandler
func NewPrivateHandler(
	sessStore *session.Store,
	dtoManager dto.IDTOManager,
	redis *redis.Client,
) *PrivateHandler {
	return &PrivateHandler{
		sess_store: sessStore,
		dtoManager: dtoManager,
		redis:      redis,
	}
}

// Init initializes the routes
func (h *PrivateHandler) Init(router fiber.Router) {
	root := router.Group("/private")
	root.Use(h.authMiddleware)

	root.Get("/", func(c *fiber.Ctx) error {
		data := sessionStore.GetSessionData(c)
		return response.Response(200, fmt.Sprintf("Dear %s %s Welcome to Codex-Compiler API (Private Zone)", data.Name, data.Surname), nil)
	})

	// Initialize Routes
	// h.initUserRoutes(root)
}

// authMiddleware checks if the user is authenticated
func (h *PrivateHandler) authMiddleware(c *fiber.Ctx) error {
	ctx := context.Background()          // Redis context
	sessionID := c.Cookies("session_id") // Get session ID from cookies
	if sessionID == "" {
		return serviceErrors.NewServiceErrorWithMessage(401, "unauthorized")
	}

	// Redis'te session kontrolü yap
	sessionData, err := h.redis.Get(ctx, sessionID).Result()
	if err == redis.Nil {
		return serviceErrors.NewServiceErrorWithMessage(401, "unauthorized")
	} else if err != nil {
		return serviceErrors.NewServiceErrorWithMessage(500, "internal server error")
	}

	// Session verilerini çözümle (örnek olarak string olarak saklanıyor)
	var data sessionStore.SessionData
	if err := data.Unmarshal(sessionData); err != nil {
		return serviceErrors.NewServiceErrorWithMessage(500, "session data error")
	}

	// Kullanıcının rolünü kontrol et
	if data.Role == "Banned" {
		return serviceErrors.NewServiceErrorWithMessage(403, "Banned")
	}

	// Kullanıcı verilerini request'e ekle
	c.Locals("user", data)

	return c.Next()
}

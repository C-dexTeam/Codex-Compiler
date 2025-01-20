package private

import (
	"fmt"

	"github.com/C-dexTeam/codex-compiler/internal/http/response"
	"github.com/gofiber/fiber/v2"
)

func (h *PrivateHandler) initUserRoutes(root fiber.Router) {
	root.Get("/run", h.Run)
}

// @Tags Run
// @Summary Runs codes.
// @Description Runs user codes.
// @Accept json
// @Produce json
// @Success 200 {object} response.BaseResponse{}
// @Router /private/run [get]
func (h *PrivateHandler) Run(c *fiber.Ctx) error {
	fmt.Println("Selam")

	return response.Response(200, "Code Runnded", "hi")
}

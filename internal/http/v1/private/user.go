package private

import (
	"fmt"

	dto "github.com/C-dexTeam/codex-compiler/internal/http/dtos"
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
// @Param quest body dto.QuestDTO true "Chapter Quest"
// @Success 200 {object} response.BaseResponse{}
// @Router /private/run [get]
func (h *PrivateHandler) Run(c *fiber.Ctx) error {
	var quest dto.QuestDTO
	if err := c.BodyParser(&quest); err != nil {
		return err
	}
	if err := h.services.UtilService().Validator().ValidateStruct(quest); err != nil {
		return err
	}

	fmt.Println("Selam")

	return response.Response(200, "Code Runnded", "hi")
}

package private

import (
	"github.com/C-dexTeam/codex-compiler/internal/domains"
	dto "github.com/C-dexTeam/codex-compiler/internal/http/dtos"
	"github.com/C-dexTeam/codex-compiler/internal/http/response"
	"github.com/C-dexTeam/codex-compiler/internal/http/sessionStore"
	"github.com/gofiber/fiber/v2"
)

func (h *PrivateHandler) initUserRoutes(root fiber.Router) {
	root.Post("/run", h.Run)
}

// @Tags Run
// @Summary Runs codes.
// @Description Runs user codes.
// @Accept json
// @Produce json
// @Param quest body dto.QuestDTO true "Chapter Quest"
// @Success 200 {object} response.BaseResponse{}
// @Router /private/run [post]
func (h *PrivateHandler) Run(c *fiber.Ctx) error {
	userSession := sessionStore.GetSessionData(c)
	var quest dto.QuestDTO
	if err := c.BodyParser(&quest); err != nil {
		return err
	}
	if err := h.services.UtilService().Validator().ValidateStruct(quest); err != nil {
		return err
	}

	lang := domains.GetLanguage(quest.ProgrammingLanguageDTO.Name)

	// Create Necessary Directories. Like UserCode etc.
	if err := h.services.RunnerService().CreateDirectories(userSession.UserID); err != nil {
		return err
	}

	// Create Necessary files. Like UserCode etc.
	if err := h.services.RunnerService().CreateFiles(userSession.UserID, lang.DefaultName, quest.Chapter, quest.Tests); err != nil {
		return err
	}

	// TODO: Build The Code
	h.services.RunnerService().BuildCode(lang.Build)

	// TODO: Run Code
	h.services.RunnerService().RunCode(quest.ProgrammingLanguageDTO.Name, quest.Tests)

	// fmt.Println("Compiler Quest:", quest)
	return response.Response(200, "Code Runnded", "hi")
}

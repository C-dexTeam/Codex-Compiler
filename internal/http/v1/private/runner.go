package private

import (
	"github.com/C-dexTeam/codex-compiler/internal/domains"
	serviceErrors "github.com/C-dexTeam/codex-compiler/internal/errors"
	dto "github.com/C-dexTeam/codex-compiler/internal/http/dtos"
	"github.com/C-dexTeam/codex-compiler/internal/http/response"
	"github.com/C-dexTeam/codex-compiler/internal/http/sessionStore"
	"github.com/gofiber/fiber/v2"
)

func (h *PrivateHandler) initUserRoutes(root fiber.Router) {
	root.Post("/run", h.Run)
	root.Get("/getPLanguages", h.GetCompileLangs)
}

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

	if lang == nil {
		return response.Response(404, "There is no such Programming Language", nil)
	}

	// Create Necessary Directories. Like UserCode etc.
	if err := h.services.RunnerService().CreateDirectories(userSession.UserID); err != nil {
		return err
	}

	// Create Necessary files. Like UserCode etc.
	if err := h.services.RunnerService().CreateFiles(userSession.UserID, lang.DefaultName, quest.Chapter, quest.Tests); err != nil {
		return err
	}

	// Build The Code for Syntax Errors
	buildLog := h.services.RunnerService().BuildCode(lang.Build, userSession.UserID, quest.Chapter.ChapterID, lang.DefaultName)
	if buildLog.BuildError != "" {
		return response.Response(400, serviceErrors.ErrCodeBuild, buildLog)
	}

	// Run Code
	codeLog := h.services.RunnerService().RunCode(userSession.UserID, quest.Chapter.ChapterID, lang.DefaultName, lang.Run, quest.Tests)

	return response.Response(200, "Code Runnded", codeLog)
}

func (h *PrivateHandler) GetCompileLangs(c *fiber.Ctx) error {
	langNames := domains.GetLanguagesName()

	return response.Response(200, "Languages", langNames)
}

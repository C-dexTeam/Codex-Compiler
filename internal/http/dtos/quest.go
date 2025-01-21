package dto

type QuestDTOManager struct{}

func NewQuestDTOManager() QuestDTOManager {
	return QuestDTOManager{}
}

type QuestDTO struct {
	Chapter                QuestChapter         `json:"chapter"`
	Tests                  []QuestTest          `json:"tests"`
	ProgrammingLanguageDTO QuestProgrammingLang `json:"programmingLanguage"`
}

type QuestChapter struct {
	ChapterID   string `json:"id"`
	UserCode    string `json:"userCode"`
	FuncName    string `json:"funcname"`
	FrontendTmp string `json:"frontendTmp"`
	DockerTmp   string `json:"dockerTmp"`
	CheckTmp    string `json:"checkTmp"`
}

type QuestProgrammingLang struct {
	Name string `json:"name"`
}

type QuestTest struct {
	TestID string `json:"id"`
	Input  string `json:"input"`
	Output string `json:"output"`
}

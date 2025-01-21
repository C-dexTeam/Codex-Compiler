package dto

type IDTOManager interface {
	QuestManager() *QuestDTOManager
}

type DTOManager struct {
	questManager *QuestDTOManager
}

func CreateNewDTOManager() IDTOManager {
	questManager := NewQuestDTOManager()

	return &DTOManager{
		questManager: &questManager,
	}
}

func (m *DTOManager) QuestManager() *QuestDTOManager {
	return m.questManager
}

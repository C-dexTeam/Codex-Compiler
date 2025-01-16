package dto

type IDTOManager interface {
}

type DTOManager struct {
}

func CreateNewDTOManager() IDTOManager {
	return &DTOManager{}
}

package services

type IService interface {
	UtilService() IUtilService
	RunnerService() *runnerService
}

type Services struct {
	utilService   IUtilService
	runnerService *runnerService
}

func CreateNewServices(
	validatorService IValidatorService,

) *Services {
	utilsService := newUtilService(validatorService)
	runnerService := NewRunnerService(utilsService)

	return &Services{
		utilService:   utilsService,
		runnerService: runnerService,
	}
}

func (s *Services) UtilService() IUtilService {
	return s.utilService
}

func (s *Services) RunnerService() *runnerService {
	return s.runnerService
}

// ------------------------------------------------------

type IValidatorService interface {
	ValidateStruct(s any) error
}

type IUtilService interface {
	Validator() IValidatorService
}

// -------------------

type utilService struct {
	validatorService IValidatorService
}

func newUtilService(
	validatorService IValidatorService,
) IUtilService {
	return &utilService{
		validatorService: validatorService,
	}
}

func (s *utilService) Validator() IValidatorService {
	return s.validatorService
}

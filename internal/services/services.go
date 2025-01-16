package services

type IService interface {
	UtilService() IUtilService
}

type Services struct {
	utilService IUtilService
}

func CreateNewServices(
	validatorService IValidatorService,

) *Services {
	utilsService := newUtilService(validatorService)

	return &Services{
		utilService: utilsService,
	}
}

func (s *Services) UtilService() IUtilService {
	return s.utilService
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

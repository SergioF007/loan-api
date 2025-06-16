package services

import (
	"loan-api/models"
	"loan-api/repositories"
)

// LoanTypeService interface para el servicio de tipos de préstamo
type LoanTypeService interface {
	GetLoanTypesWithForms(tenantID uint) ([]models.LoanTypeResponse, error)
	GetLoanTypeByCode(tenantID uint, code string) (*models.LoanTypeResponse, error)
	GetLoanTypeByID(id uint) (*models.LoanTypeResponse, error)
}

// loanTypeService implementación del servicio
type loanTypeService struct {
	loanTypeRepo repositories.LoanTypeRepository
}

// NewLoanTypeService crea una nueva instancia del servicio
func NewLoanTypeService(loanTypeRepo repositories.LoanTypeRepository) LoanTypeService {
	return &loanTypeService{
		loanTypeRepo: loanTypeRepo,
	}
}

// GetLoanTypesWithForms obtiene todos los tipos de préstamo con formularios por tenant
func (s *loanTypeService) GetLoanTypesWithForms(tenantID uint) ([]models.LoanTypeResponse, error) {
	loanTypes, err := s.loanTypeRepo.GetActiveByTenantID(tenantID)
	if err != nil {
		return nil, err
	}

	response := make([]models.LoanTypeResponse, len(loanTypes))
	for i, loanType := range loanTypes {
		response[i] = s.buildLoanTypeResponse(loanType)
	}

	return response, nil
}

// GetLoanTypeByCode obtiene un tipo de préstamo por código y tenant
func (s *loanTypeService) GetLoanTypeByCode(tenantID uint, code string) (*models.LoanTypeResponse, error) {
	loanType, err := s.loanTypeRepo.GetByTenantIDAndCode(tenantID, code)
	if err != nil {
		return nil, err
	}

	response := s.buildLoanTypeResponse(*loanType)
	return &response, nil
}

// GetLoanTypeByID obtiene un tipo de préstamo por ID con formularios
func (s *loanTypeService) GetLoanTypeByID(id uint) (*models.LoanTypeResponse, error) {
	loanType, err := s.loanTypeRepo.GetByIDWithForms(id)
	if err != nil {
		return nil, err
	}

	response := s.buildLoanTypeResponse(*loanType)
	return &response, nil
}

// buildLoanTypeResponse construye la respuesta del tipo de préstamo
func (s *loanTypeService) buildLoanTypeResponse(loanType models.LoanType) models.LoanTypeResponse {
	response := models.LoanTypeResponse{
		ID:          loanType.ID,
		Name:        loanType.Name,
		Code:        loanType.Code,
		Description: loanType.Description,
		MinAmount:   loanType.MinAmount,
		MaxAmount:   loanType.MaxAmount,
	}

	// Tomar la versión por defecto (la primera activa)
	if len(loanType.Versions) > 0 {
		version := loanType.Versions[0]
		response.Version = models.LoanTypeVersionResponse{
			ID:          version.ID,
			Version:     version.Version,
			Description: version.Description,
			Forms:       s.buildFormsResponse(version.Forms),
			FormInputs:  s.buildFormInputsResponse(version.FormInputs),
		}
	}

	return response
}

// buildFormsResponse construye la respuesta de formularios
func (s *loanTypeService) buildFormsResponse(forms []models.LoanTypeForm) []models.LoanTypeFormResponse {
	response := make([]models.LoanTypeFormResponse, len(forms))
	for i, form := range forms {
		response[i] = models.LoanTypeFormResponse{
			ID:          form.ID,
			Label:       form.Label,
			Code:        form.Code,
			Description: form.Description,
			Order:       form.Order,
			IsRequired:  form.IsRequired,
			FormInputs:  s.buildFormInputsResponse(form.FormInputs),
		}
	}
	return response
}

// buildFormInputsResponse construye la respuesta de inputs de formulario
func (s *loanTypeService) buildFormInputsResponse(inputs []models.LoanTypeVersionFormInput) []models.LoanTypeVersionFormInputResponse {
	response := make([]models.LoanTypeVersionFormInputResponse, len(inputs))
	for i, input := range inputs {
		response[i] = models.LoanTypeVersionFormInputResponse{
			ID:              input.ID,
			Label:           input.Label,
			Code:            input.Code,
			InputType:       input.InputType,
			Placeholder:     input.Placeholder,
			DefaultValue:    input.DefaultValue,
			ValidationRules: input.ValidationRules,
			Options:         input.Options,
			Order:           input.Order,
			IsRequired:      input.IsRequired,
		}
	}
	return response
}

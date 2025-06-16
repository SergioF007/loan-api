package services

import (
	"loan-api/models"
	"loan-api/repositories"
)

// TenantService interface para el servicio de tenant
type TenantService interface {
	GetAvailableTenants() ([]models.TenantResponse, error)
	GetTenantByCode(code string) (*models.TenantResponse, error)
	ValidateTenantID(tenantID uint) (*models.Tenant, error)
}

// tenantService implementación del servicio
type tenantService struct {
	tenantRepo repositories.TenantRepository
}

// NewTenantService crea una nueva instancia del servicio
func NewTenantService(tenantRepo repositories.TenantRepository) TenantService {
	return &tenantService{
		tenantRepo: tenantRepo,
	}
}

// GetAvailableTenants obtiene todos los tenants disponibles para pruebas
func (s *tenantService) GetAvailableTenants() ([]models.TenantResponse, error) {
	tenants, err := s.tenantRepo.GetAllActive()
	if err != nil {
		return nil, err
	}

	response := make([]models.TenantResponse, len(tenants))
	for i, tenant := range tenants {
		response[i] = models.TenantResponse{
			ID:          tenant.ID,
			Name:        tenant.Name,
			Code:        tenant.Code,
			Description: tenant.Description,
			IsActive:    tenant.IsActive,
		}
	}

	return response, nil
}

// GetTenantByCode obtiene un tenant por código
func (s *tenantService) GetTenantByCode(code string) (*models.TenantResponse, error) {
	tenant, err := s.tenantRepo.GetByCode(code)
	if err != nil {
		return nil, err
	}

	response := &models.TenantResponse{
		ID:          tenant.ID,
		Name:        tenant.Name,
		Code:        tenant.Code,
		Description: tenant.Description,
		IsActive:    tenant.IsActive,
	}

	return response, nil
}

// ValidateTenantID valida que un tenant ID existe y está activo
func (s *tenantService) ValidateTenantID(tenantID uint) (*models.Tenant, error) {
	return s.tenantRepo.GetByID(tenantID)
}

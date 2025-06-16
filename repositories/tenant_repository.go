package repositories

import (
	"loan-api/models"

	"gorm.io/gorm"
)

// TenantRepository interface para operaciones de tenant
type TenantRepository interface {
	GetAllActive() ([]models.Tenant, error)
	GetByCode(code string) (*models.Tenant, error)
	GetByID(id uint) (*models.Tenant, error)
}

// tenantRepository implementación del repository
type tenantRepository struct {
	db *gorm.DB
}

// NewTenantRepository crea una nueva instancia del repository
func NewTenantRepository(db *gorm.DB) TenantRepository {
	return &tenantRepository{db: db}
}

// GetAllActive obtiene todos los tenants activos
func (r *tenantRepository) GetAllActive() ([]models.Tenant, error) {
	var tenants []models.Tenant
	err := r.db.Where("is_active = ?", true).Find(&tenants).Error
	return tenants, err
}

// GetByCode obtiene un tenant por código
func (r *tenantRepository) GetByCode(code string) (*models.Tenant, error) {
	var tenant models.Tenant
	err := r.db.Where("code = ? AND is_active = ?", code, true).First(&tenant).Error
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

// GetByID obtiene un tenant por ID
func (r *tenantRepository) GetByID(id uint) (*models.Tenant, error) {
	var tenant models.Tenant
	err := r.db.Where("id = ? AND is_active = ?", id, true).First(&tenant).Error
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

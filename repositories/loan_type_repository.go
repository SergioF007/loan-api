package repositories

import (
	"loan-api/models"

	"gorm.io/gorm"
)

// LoanTypeRepository interface para operaciones de loan type
type LoanTypeRepository interface {
	GetByTenantID(tenantID uint) ([]models.LoanType, error)
	GetByTenantIDAndCode(tenantID uint, code string) (*models.LoanType, error)
	GetByIDWithForms(id uint) (*models.LoanType, error)
	GetActiveByTenantID(tenantID uint) ([]models.LoanType, error)
}

// loanTypeRepository implementación del repository
type loanTypeRepository struct {
	db *gorm.DB
}

// NewLoanTypeRepository crea una nueva instancia del repository
func NewLoanTypeRepository(db *gorm.DB) LoanTypeRepository {
	return &loanTypeRepository{db: db}
}

// GetByTenantID obtiene todos los tipos de préstamo por tenant
func (r *loanTypeRepository) GetByTenantID(tenantID uint) ([]models.LoanType, error) {
	var loanTypes []models.LoanType
	err := r.db.Where("tenant_id = ?", tenantID).
		Preload("Versions").
		Find(&loanTypes).Error
	return loanTypes, err
}

// GetByTenantIDAndCode obtiene un tipo de préstamo por tenant y código
func (r *loanTypeRepository) GetByTenantIDAndCode(tenantID uint, code string) (*models.LoanType, error) {
	var loanType models.LoanType
	err := r.db.Where("tenant_id = ? AND code = ? AND is_active = ?", tenantID, code, true).
		Preload("Versions", "is_active = ?", true).
		Preload("Versions.Forms", "is_active = ?", true).
		Preload("Versions.FormInputs", "is_active = ?", true).
		First(&loanType).Error
	if err != nil {
		return nil, err
	}
	return &loanType, nil
}

// GetByIDWithForms obtiene un tipo de préstamo con todos sus formularios
func (r *loanTypeRepository) GetByIDWithForms(id uint) (*models.LoanType, error) {
	var loanType models.LoanType
	err := r.db.Where("id = ? AND is_active = ?", id, true).
		Preload("Versions", "is_active = ? AND is_default = ?", true, true).
		Preload("Versions.Forms", "is_active = ?", true).
		Preload("Versions.Forms.FormInputs", "is_active = ?", true).
		Preload("Versions.FormInputs", "is_active = ?", true).
		First(&loanType).Error
	if err != nil {
		return nil, err
	}
	return &loanType, nil
}

// GetActiveByTenantID obtiene todos los tipos de préstamo activos por tenant
func (r *loanTypeRepository) GetActiveByTenantID(tenantID uint) ([]models.LoanType, error) {
	var loanTypes []models.LoanType
	err := r.db.Where("tenant_id = ? AND is_active = ?", tenantID, true).
		Preload("Versions", "is_active = ? AND is_default = ?", true, true).
		Find(&loanTypes).Error
	return loanTypes, err
}

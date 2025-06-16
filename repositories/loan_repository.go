package repositories

import (
	"loan-api/models"

	"gorm.io/gorm"
)

// LoanRepository interface para operaciones de préstamo
type LoanRepository interface {
	Create(loan *models.Loan) error
	GetByID(id uint) (*models.Loan, error)
	GetByUserID(userID uint) ([]models.Loan, error)
	Update(loan *models.Loan) error
	SaveLoanData(loanData []models.LoanData) error
	GetLoanDataByLoanID(loanID uint) ([]models.LoanData, error)
	DeleteLoanDataByLoanID(loanID uint) error
}

// loanRepository implementación del repository
type loanRepository struct {
	db *gorm.DB
}

// NewLoanRepository crea una nueva instancia del repository
func NewLoanRepository(db *gorm.DB) LoanRepository {
	return &loanRepository{db: db}
}

// Create crea un nuevo préstamo
func (r *loanRepository) Create(loan *models.Loan) error {
	return r.db.Create(loan).Error
}

// GetByID obtiene un préstamo por ID con todas sus relaciones
func (r *loanRepository) GetByID(id uint) (*models.Loan, error) {
	var loan models.Loan
	err := r.db.Where("id = ?", id).
		Preload("LoanType").
		Preload("User").
		Preload("Data").
		First(&loan).Error
	if err != nil {
		return nil, err
	}
	return &loan, nil
}

// GetByUserID obtiene todos los préstamos de un usuario
func (r *loanRepository) GetByUserID(userID uint) ([]models.Loan, error) {
	var loans []models.Loan
	err := r.db.Where("user_id = ?", userID).
		Preload("LoanType").
		Preload("User").
		Preload("Data").
		Find(&loans).Error
	return loans, err
}

// Update actualiza un préstamo
func (r *loanRepository) Update(loan *models.Loan) error {
	return r.db.Save(loan).Error
}

// SaveLoanData guarda los datos de un préstamo
func (r *loanRepository) SaveLoanData(loanData []models.LoanData) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, data := range loanData {
			if err := tx.Create(&data).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// GetLoanDataByLoanID obtiene todos los datos de un préstamo
func (r *loanRepository) GetLoanDataByLoanID(loanID uint) ([]models.LoanData, error) {
	var loanData []models.LoanData
	err := r.db.Where("loan_id = ?", loanID).
		Preload("Form").
		Find(&loanData).Error
	return loanData, err
}

// DeleteLoanDataByLoanID elimina todos los datos de un préstamo
func (r *loanRepository) DeleteLoanDataByLoanID(loanID uint) error {
	return r.db.Where("loan_id = ?", loanID).Delete(&models.LoanData{}).Error
}

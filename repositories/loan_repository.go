package repositories

import (
	"errors"
	"time"

	"loan-api/app_error"
	"loan-api/models"

	"gorm.io/gorm"
)

// LoanRepository define la interfaz para operaciones de préstamo
type LoanRepository interface {
	Create(loan *models.Loan) error
	GetByID(id uint) (*models.Loan, error)
	GetByIDWithUser(id uint) (*models.Loan, error)
	Update(loan *models.Loan) error
	Delete(id uint) error
	List(limit, offset int) ([]models.Loan, int64, error)
	ListWithUsers(limit, offset int) ([]models.Loan, int64, error)
	GetByUserID(userID uint, limit, offset int) ([]models.Loan, int64, error)
	GetByStatus(status models.LoanStatus, limit, offset int) ([]models.Loan, int64, error)
	UpdateStatus(id uint, status models.LoanStatus) error
	CountByUserID(userID uint) (int64, error)
	GetPendingLoans(limit, offset int) ([]models.Loan, int64, error)
}

// loanRepository implementa LoanRepository
type loanRepository struct {
	db *gorm.DB
}

// NewLoanRepository crea una nueva instancia del repositorio de préstamos
func NewLoanRepository(db *gorm.DB) LoanRepository {
	return &loanRepository{
		db: db,
	}
}

// Create crea un nuevo préstamo
func (r *loanRepository) Create(loan *models.Loan) error {
	// Establecer fecha de solicitud si no está establecida
	if loan.RequestDate.IsZero() {
		loan.RequestDate = time.Now()
	}

	// Establecer estado por defecto si no está establecido
	if loan.Status == "" {
		loan.Status = models.LoanStatusPending
	}

	if err := r.db.Create(loan).Error; err != nil {
		return app_error.NewDatabaseError("crear préstamo", err.Error())
	}
	return nil
}

// GetByID obtiene un préstamo por su ID
func (r *loanRepository) GetByID(id uint) (*models.Loan, error) {
	var loan models.Loan
	if err := r.db.First(&loan, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, app_error.ErrLoanNotFound
		}
		return nil, app_error.NewDatabaseError("obtener préstamo", err.Error())
	}
	return &loan, nil
}

// GetByIDWithUser obtiene un préstamo por su ID incluyendo datos del usuario
func (r *loanRepository) GetByIDWithUser(id uint) (*models.Loan, error) {
	var loan models.Loan
	if err := r.db.Preload("User").First(&loan, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, app_error.ErrLoanNotFound
		}
		return nil, app_error.NewDatabaseError("obtener préstamo con usuario", err.Error())
	}
	return &loan, nil
}

// Update actualiza un préstamo existente
func (r *loanRepository) Update(loan *models.Loan) error {
	if err := r.db.Save(loan).Error; err != nil {
		return app_error.NewDatabaseError("actualizar préstamo", err.Error())
	}
	return nil
}

// Delete elimina un préstamo (soft delete)
func (r *loanRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Loan{}, id)
	if result.Error != nil {
		return app_error.NewDatabaseError("eliminar préstamo", result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return app_error.ErrLoanNotFound
	}
	return nil
}

// List obtiene una lista paginada de préstamos
func (r *loanRepository) List(limit, offset int) ([]models.Loan, int64, error) {
	var loans []models.Loan
	var total int64

	// Contar total de registros
	if err := r.db.Model(&models.Loan{}).Count(&total).Error; err != nil {
		return nil, 0, app_error.NewDatabaseError("contar préstamos", err.Error())
	}

	// Obtener préstamos paginados ordenados por fecha de solicitud descendente
	if err := r.db.Order("request_date DESC").Limit(limit).Offset(offset).Find(&loans).Error; err != nil {
		return nil, 0, app_error.NewDatabaseError("listar préstamos", err.Error())
	}

	return loans, total, nil
}

// ListWithUsers obtiene una lista paginada de préstamos con información del usuario
func (r *loanRepository) ListWithUsers(limit, offset int) ([]models.Loan, int64, error) {
	var loans []models.Loan
	var total int64

	// Contar total de registros
	if err := r.db.Model(&models.Loan{}).Count(&total).Error; err != nil {
		return nil, 0, app_error.NewDatabaseError("contar préstamos", err.Error())
	}

	// Obtener préstamos con usuarios paginados
	if err := r.db.Preload("User").Order("request_date DESC").Limit(limit).Offset(offset).Find(&loans).Error; err != nil {
		return nil, 0, app_error.NewDatabaseError("listar préstamos con usuarios", err.Error())
	}

	return loans, total, nil
}

// GetByUserID obtiene préstamos de un usuario específico
func (r *loanRepository) GetByUserID(userID uint, limit, offset int) ([]models.Loan, int64, error) {
	var loans []models.Loan
	var total int64

	// Contar total de registros para este usuario
	if err := r.db.Model(&models.Loan{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, app_error.NewDatabaseError("contar préstamos del usuario", err.Error())
	}

	// Obtener préstamos del usuario
	if err := r.db.Where("user_id = ?", userID).Order("request_date DESC").Limit(limit).Offset(offset).Find(&loans).Error; err != nil {
		return nil, 0, app_error.NewDatabaseError("obtener préstamos del usuario", err.Error())
	}

	return loans, total, nil
}

// GetByStatus obtiene préstamos por estado
func (r *loanRepository) GetByStatus(status models.LoanStatus, limit, offset int) ([]models.Loan, int64, error) {
	var loans []models.Loan
	var total int64

	// Contar total de registros con este estado
	if err := r.db.Model(&models.Loan{}).Where("status = ?", status).Count(&total).Error; err != nil {
		return nil, 0, app_error.NewDatabaseError("contar préstamos por estado", err.Error())
	}

	// Obtener préstamos por estado
	if err := r.db.Preload("User").Where("status = ?", status).Order("request_date DESC").Limit(limit).Offset(offset).Find(&loans).Error; err != nil {
		return nil, 0, app_error.NewDatabaseError("obtener préstamos por estado", err.Error())
	}

	return loans, total, nil
}

// UpdateStatus actualiza solo el estado de un préstamo
func (r *loanRepository) UpdateStatus(id uint, status models.LoanStatus) error {
	updates := map[string]interface{}{
		"status": status,
	}

	// Si el estado es aprobado o rechazado, establecer fecha de aprobación
	if status == models.LoanStatusApproved || status == models.LoanStatusRejected {
		updates["approval_date"] = time.Now()
	}

	result := r.db.Model(&models.Loan{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return app_error.NewDatabaseError("actualizar estado del préstamo", result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return app_error.ErrLoanNotFound
	}
	return nil
}

// CountByUserID cuenta el número de préstamos de un usuario
func (r *loanRepository) CountByUserID(userID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&models.Loan{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, app_error.NewDatabaseError("contar préstamos del usuario", err.Error())
	}
	return count, nil
}

// GetPendingLoans obtiene préstamos pendientes (útil para procesos de aprobación)
func (r *loanRepository) GetPendingLoans(limit, offset int) ([]models.Loan, int64, error) {
	return r.GetByStatus(models.LoanStatusPending, limit, offset)
}

// GetLoanStats obtiene estadísticas de préstamos (método adicional)
func (r *loanRepository) GetLoanStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total de préstamos
	var total int64
	if err := r.db.Model(&models.Loan{}).Count(&total).Error; err != nil {
		return nil, app_error.NewDatabaseError("obtener estadísticas de préstamos", err.Error())
	}
	stats["total"] = total

	// Préstamos por estado
	var pendingCount, approvedCount, rejectedCount int64

	if err := r.db.Model(&models.Loan{}).Where("status = ?", models.LoanStatusPending).Count(&pendingCount).Error; err != nil {
		return nil, app_error.NewDatabaseError("contar préstamos pendientes", err.Error())
	}
	stats["pending"] = pendingCount

	if err := r.db.Model(&models.Loan{}).Where("status = ?", models.LoanStatusApproved).Count(&approvedCount).Error; err != nil {
		return nil, app_error.NewDatabaseError("contar préstamos aprobados", err.Error())
	}
	stats["approved"] = approvedCount

	if err := r.db.Model(&models.Loan{}).Where("status = ?", models.LoanStatusRejected).Count(&rejectedCount).Error; err != nil {
		return nil, app_error.NewDatabaseError("contar préstamos rechazados", err.Error())
	}
	stats["rejected"] = rejectedCount

	return stats, nil
}

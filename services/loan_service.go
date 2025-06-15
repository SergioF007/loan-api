package services

import (
	"time"

	"loan-api/app_error"
	"loan-api/models"
	"loan-api/repositories"

	"github.com/go-playground/validator/v10"
)

// LoanService define la interfaz para la lógica de negocio de préstamos
type LoanService interface {
	CreateLoan(req *models.LoanRequest) (*models.Loan, error)
	GetLoanByID(id uint) (*models.Loan, error)
	GetLoanByIDWithUser(id uint) (*models.Loan, error)
	UpdateLoanStatus(id uint, req *models.LoanStatusUpdateRequest) (*models.Loan, error)
	ListLoans(page, limit int) ([]models.Loan, int64, error)
	ListLoansWithUsers(page, limit int) ([]models.Loan, int64, error)
	GetLoansByUserID(userID uint, page, limit int) ([]models.Loan, int64, error)
	GetLoansByStatus(status string, page, limit int) ([]models.Loan, int64, error)
	ValidateLoan(req *models.LoanRequest) error
	ProcessLoanApproval(loanID uint) error
	GetLoanStatistics() (map[string]interface{}, error)
}

// loanService implementa LoanService
type loanService struct {
	loanRepo  repositories.LoanRepository
	userRepo  repositories.UserRepository
	validator *validator.Validate
}

// NewLoanService crea una nueva instancia del servicio de préstamos
func NewLoanService(loanRepo repositories.LoanRepository, userRepo repositories.UserRepository) LoanService {
	return &loanService{
		loanRepo:  loanRepo,
		userRepo:  userRepo,
		validator: validator.New(),
	}
}

// CreateLoan crea una nueva solicitud de préstamo
func (s *loanService) CreateLoan(req *models.LoanRequest) (*models.Loan, error) {
	// Validar datos de entrada
	if err := s.ValidateLoan(req); err != nil {
		return nil, err
	}

	// Verificar que el usuario existe
	user, err := s.userRepo.GetByID(req.UserID)
	if err != nil {
		return nil, err
	}

	// Validar lógica de negocio para aprobación
	if err := s.validateLoanBusinessRules(req, user); err != nil {
		return nil, err
	}

	// Crear modelo de préstamo
	loan := &models.Loan{
		UserID:      req.UserID,
		Amount:      req.Amount,
		Purpose:     req.Purpose,
		Status:      models.LoanStatusPending,
		RequestDate: time.Now(),
	}

	// Guardar en base de datos
	if err := s.loanRepo.Create(loan); err != nil {
		return nil, err
	}

	// Cargar datos del usuario para la respuesta
	loan.User = user

	return loan, nil
}

// GetLoanByID obtiene un préstamo por su ID
func (s *loanService) GetLoanByID(id uint) (*models.Loan, error) {
	if id == 0 {
		return nil, app_error.NewValidationError("id", "ID de préstamo es requerido")
	}

	return s.loanRepo.GetByID(id)
}

// GetLoanByIDWithUser obtiene un préstamo por su ID incluyendo datos del usuario
func (s *loanService) GetLoanByIDWithUser(id uint) (*models.Loan, error) {
	if id == 0 {
		return nil, app_error.NewValidationError("id", "ID de préstamo es requerido")
	}

	return s.loanRepo.GetByIDWithUser(id)
}

// UpdateLoanStatus actualiza el estado de un préstamo
func (s *loanService) UpdateLoanStatus(id uint, req *models.LoanStatusUpdateRequest) (*models.Loan, error) {
	// Validar ID
	if id == 0 {
		return nil, app_error.NewValidationError("id", "ID de préstamo es requerido")
	}

	// Validar estado
	if err := s.validator.Struct(req); err != nil {
		return nil, app_error.NewValidationError("status", "Estado inválido")
	}

	// Obtener préstamo actual
	loan, err := s.loanRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Verificar si se puede actualizar el estado
	if !loan.CanUpdateStatus(req.Status) {
		return nil, app_error.ErrCannotUpdateStatus
	}

	// Actualizar estado
	if err := s.loanRepo.UpdateStatus(id, req.Status); err != nil {
		return nil, err
	}

	// Retornar préstamo actualizado
	return s.loanRepo.GetByIDWithUser(id)
}

// ListLoans obtiene una lista paginada de préstamos
func (s *loanService) ListLoans(page, limit int) ([]models.Loan, int64, error) {
	// Validar parámetros de paginación
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20 // Límite por defecto
	}

	offset := (page - 1) * limit
	return s.loanRepo.List(limit, offset)
}

// ListLoansWithUsers obtiene una lista paginada de préstamos con datos del usuario
func (s *loanService) ListLoansWithUsers(page, limit int) ([]models.Loan, int64, error) {
	// Validar parámetros de paginación
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20 // Límite por defecto
	}

	offset := (page - 1) * limit
	return s.loanRepo.ListWithUsers(limit, offset)
}

// GetLoansByUserID obtiene préstamos de un usuario específico
func (s *loanService) GetLoansByUserID(userID uint, page, limit int) ([]models.Loan, int64, error) {
	// Validar ID de usuario
	if userID == 0 {
		return nil, 0, app_error.NewValidationError("user_id", "ID de usuario es requerido")
	}

	// Verificar que el usuario existe
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, 0, err
	}

	// Validar parámetros de paginación
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20 // Límite por defecto
	}

	offset := (page - 1) * limit
	return s.loanRepo.GetByUserID(userID, limit, offset)
}

// GetLoansByStatus obtiene préstamos por estado
func (s *loanService) GetLoansByStatus(status string, page, limit int) ([]models.Loan, int64, error) {
	// Validar estado
	loanStatus := models.LoanStatus(status)
	if !isValidLoanStatus(loanStatus) {
		return nil, 0, app_error.NewValidationError("status", "Estado de préstamo inválido")
	}

	// Validar parámetros de paginación
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20 // Límite por defecto
	}

	offset := (page - 1) * limit
	return s.loanRepo.GetByStatus(loanStatus, limit, offset)
}

// ValidateLoan valida los datos de un préstamo
func (s *loanService) ValidateLoan(req *models.LoanRequest) error {
	if err := s.validator.Struct(req); err != nil {
		// Procesar errores de validación
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Tag() {
			case "required":
				validationErrors = append(validationErrors, err.Field()+" es requerido")
			case "min":
				validationErrors = append(validationErrors, err.Field()+" debe ser mayor a "+err.Param())
			case "max":
				validationErrors = append(validationErrors, err.Field()+" no puede ser mayor a "+err.Param())
			default:
				validationErrors = append(validationErrors, err.Field()+" no es válido")
			}
		}

		return app_error.NewValidationError("validación", validationErrors[0])
	}
	return nil
}

// ProcessLoanApproval procesa la aprobación automática de préstamos
func (s *loanService) ProcessLoanApproval(loanID uint) error {
	// Obtener préstamo con datos del usuario
	loan, err := s.loanRepo.GetByIDWithUser(loanID)
	if err != nil {
		return err
	}

	// Solo procesar préstamos pendientes
	if loan.Status != models.LoanStatusPending {
		return app_error.NewBusinessError("Préstamo no procesable", "Solo se pueden procesar préstamos pendientes")
	}

	// Evaluar criterios de aprobación
	newStatus := s.evaluateLoanApproval(loan)

	// Actualizar estado
	return s.loanRepo.UpdateStatus(loanID, newStatus)
}

// validateLoanBusinessRules valida reglas de negocio específicas
func (s *loanService) validateLoanBusinessRules(req *models.LoanRequest, user *models.User) error {
	// Validar monto mínimo
	if req.Amount < 1000 {
		return app_error.NewBusinessError(
			"Monto insuficiente",
			"El monto mínimo del préstamo es $1,000",
		)
	}

	// Validar monto máximo basado en ingresos
	maxAmount := user.Income * 5 // Máximo 5 veces los ingresos mensuales
	if req.Amount > maxAmount {
		return app_error.NewBusinessError(
			"Monto excesivo",
			"El monto no puede ser mayor a 5 veces sus ingresos mensuales",
		)
	}

	// Validar prontitud crediticia mínima
	if user.CreditScore < 500 {
		return app_error.NewBusinessError(
			"Puntaje crediticio insuficiente",
			"Se requiere un puntaje crediticio mínimo de 500 para solicitar un préstamo",
		)
	}

	// Verificar número máximo de préstamos activos
	activeLoans, err := s.loanRepo.CountByUserID(user.ID)
	if err != nil {
		return err
	}

	if activeLoans >= 3 {
		return app_error.NewBusinessError(
			"Límite de préstamos excedido",
			"No puede tener más de 3 préstamos activos al mismo tiempo",
		)
	}

	return nil
}

// evaluateLoanApproval evalúa si un préstamo debe ser aprobado automáticamente
func (s *loanService) evaluateLoanApproval(loan *models.Loan) models.LoanStatus {
	user := loan.User

	// Criterios de aprobación automática
	score := 0

	// Puntaje crediticio (40% del peso)
	if user.CreditScore >= 750 {
		score += 40
	} else if user.CreditScore >= 700 {
		score += 30
	} else if user.CreditScore >= 650 {
		score += 20
	} else if user.CreditScore >= 600 {
		score += 10
	}

	// Relación préstamo/ingresos (30% del peso)
	loanToIncomeRatio := loan.Amount / user.Income
	if loanToIncomeRatio <= 2 {
		score += 30
	} else if loanToIncomeRatio <= 3 {
		score += 20
	} else if loanToIncomeRatio <= 4 {
		score += 10
	}

	// Ingresos estables (20% del peso)
	if user.Income >= 5000 {
		score += 20
	} else if user.Income >= 3000 {
		score += 15
	} else if user.Income >= 2000 {
		score += 10
	}

	// Monto del préstamo (10% del peso)
	if loan.Amount <= 10000 {
		score += 10
	} else if loan.Amount <= 50000 {
		score += 5
	}

	// Decidir basado en el puntaje
	if score >= 80 {
		return models.LoanStatusApproved
	} else if score < 50 {
		return models.LoanStatusRejected
	} else {
		// Mantener pendiente para revisión manual
		return models.LoanStatusPending
	}
}

// isValidLoanStatus verifica si el estado de préstamo es válido
func isValidLoanStatus(status models.LoanStatus) bool {
	return status == models.LoanStatusPending ||
		status == models.LoanStatusApproved ||
		status == models.LoanStatusRejected
}

// GetLoanStatistics obtiene estadísticas de préstamos (método adicional)
func (s *loanService) GetLoanStatistics() (map[string]interface{}, error) {
	// Usar el método del repositorio si está disponible
	// Para este ejemplo, implementamos una versión básica
	stats := make(map[string]interface{})

	// Este método podría expandirse para incluir más estadísticas
	stats["message"] = "Estadísticas de préstamos disponibles"

	return stats, nil
}

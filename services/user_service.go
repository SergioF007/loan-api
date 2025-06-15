package services

import (
	"loan-api/app_error"
	"loan-api/models"
	"loan-api/repositories"

	"github.com/go-playground/validator/v10"
)

// UserService define la interfaz para la lógica de negocio de usuarios
type UserService interface {
	CreateUser(req *models.UserRequest) (*models.User, error)
	GetUserByID(id uint) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	UpdateUser(id uint, req *models.UserRequest) (*models.User, error)
	DeleteUser(id uint) error
	ListUsers(page, limit int) ([]models.User, int64, error)
	ValidateUser(req *models.UserRequest) error
	GetUserCreditSummary(userID uint) (map[string]interface{}, error)
}

// userService implementa UserService
type userService struct {
	userRepo  repositories.UserRepository
	validator *validator.Validate
}

// NewUserService crea una nueva instancia del servicio de usuarios
func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo:  userRepo,
		validator: validator.New(),
	}
}

// CreateUser crea un nuevo usuario
func (s *userService) CreateUser(req *models.UserRequest) (*models.User, error) {
	// Validar datos de entrada
	if err := s.ValidateUser(req); err != nil {
		return nil, err
	}

	// Verificar si el email ya existe
	exists, err := s.userRepo.ExistsByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, app_error.ErrEmailExists
	}

	// Validar lógica de negocio
	if err := s.validateUserBusinessRules(req); err != nil {
		return nil, err
	}

	// Crear modelo de usuario
	user := &models.User{
		Name:        req.Name,
		Email:       req.Email,
		Phone:       req.Phone,
		Income:      req.Income,
		CreditScore: req.CreditScore,
	}

	// Guardar en base de datos
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByID obtiene un usuario por su ID
func (s *userService) GetUserByID(id uint) (*models.User, error) {
	if id == 0 {
		return nil, app_error.NewValidationError("id", "ID de usuario es requerido")
	}

	return s.userRepo.GetByID(id)
}

// GetUserByEmail obtiene un usuario por su email
func (s *userService) GetUserByEmail(email string) (*models.User, error) {
	if email == "" {
		return nil, app_error.NewValidationError("email", "Email es requerido")
	}

	return s.userRepo.GetByEmail(email)
}

// UpdateUser actualiza un usuario existente
func (s *userService) UpdateUser(id uint, req *models.UserRequest) (*models.User, error) {
	// Validar ID
	if id == 0 {
		return nil, app_error.NewValidationError("id", "ID de usuario es requerido")
	}

	// Validar datos de entrada
	if err := s.ValidateUser(req); err != nil {
		return nil, err
	}

	// Verificar que el usuario existe
	existingUser, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Verificar si el email ya existe en otro usuario
	if req.Email != existingUser.Email {
		exists, err := s.userRepo.ExistsByEmail(req.Email)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, app_error.ErrEmailExists
		}
	}

	// Validar lógica de negocio
	if err := s.validateUserBusinessRules(req); err != nil {
		return nil, err
	}

	// Actualizar campos
	existingUser.Name = req.Name
	existingUser.Email = req.Email
	existingUser.Phone = req.Phone
	existingUser.Income = req.Income
	existingUser.CreditScore = req.CreditScore

	// Guardar cambios
	if err := s.userRepo.Update(existingUser); err != nil {
		return nil, err
	}

	return existingUser, nil
}

// DeleteUser elimina un usuario
func (s *userService) DeleteUser(id uint) error {
	if id == 0 {
		return app_error.NewValidationError("id", "ID de usuario es requerido")
	}

	// Verificar que el usuario existe
	_, err := s.userRepo.GetByID(id)
	if err != nil {
		return err
	}

	// TODO: Verificar si el usuario tiene préstamos activos
	// Esto se puede implementar cuando se tenga el loan service

	return s.userRepo.Delete(id)
}

// ListUsers obtiene una lista paginada de usuarios
func (s *userService) ListUsers(page, limit int) ([]models.User, int64, error) {
	// Validar parámetros de paginación
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20 // Límite por defecto
	}

	offset := (page - 1) * limit
	return s.userRepo.List(limit, offset)
}

// ValidateUser valida los datos de un usuario
func (s *userService) ValidateUser(req *models.UserRequest) error {
	if err := s.validator.Struct(req); err != nil {
		// Procesar errores de validación
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Tag() {
			case "required":
				validationErrors = append(validationErrors, err.Field()+" es requerido")
			case "email":
				validationErrors = append(validationErrors, "Email debe tener un formato válido")
			case "min":
				validationErrors = append(validationErrors, err.Field()+" debe tener al menos "+err.Param()+" caracteres")
			case "max":
				validationErrors = append(validationErrors, err.Field()+" no puede tener más de "+err.Param()+" caracteres")
			default:
				validationErrors = append(validationErrors, err.Field()+" no es válido")
			}
		}

		return app_error.NewValidationError("validación", validationErrors[0])
	}
	return nil
}

// validateUserBusinessRules valida reglas de negocio específicas
func (s *userService) validateUserBusinessRules(req *models.UserRequest) error {
	// Validar puntaje crediticio
	if req.CreditScore < 300 || req.CreditScore > 850 {
		return app_error.NewBusinessError(
			"Puntaje crediticio inválido",
			"El puntaje crediticio debe estar entre 300 y 850",
		)
	}

	// Validar ingresos mínimos
	if req.Income < 0 {
		return app_error.NewBusinessError(
			"Ingresos inválidos",
			"Los ingresos no pueden ser negativos",
		)
	}

	// Validar formato de teléfono básico
	if len(req.Phone) < 10 {
		return app_error.NewBusinessError(
			"Teléfono inválido",
			"El teléfono debe tener al menos 10 dígitos",
		)
	}

	return nil
}

// GetUserCreditSummary obtiene un resumen crediticio del usuario (método adicional)
func (s *userService) GetUserCreditSummary(userID uint) (map[string]interface{}, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	summary := map[string]interface{}{
		"user_id":      user.ID,
		"name":         user.Name,
		"credit_score": user.CreditScore,
		"income":       user.Income,
		"status":       s.getCreditStatus(user.CreditScore),
	}

	return summary, nil
}

// getCreditStatus determina el estado crediticio basado en el puntaje
func (s *userService) getCreditStatus(score int) string {
	if score >= 750 {
		return "Excelente"
	} else if score >= 700 {
		return "Muy Bueno"
	} else if score >= 650 {
		return "Bueno"
	} else if score >= 600 {
		return "Regular"
	} else {
		return "Pobre"
	}
}

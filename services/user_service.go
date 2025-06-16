package services

import (
	"loan-api/app_error"
	"loan-api/config"
	"loan-api/models"
	"loan-api/repositories"
	"loan-api/utils"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

// UserService define la interfaz para la lógica de negocio de usuarios
type UserService interface {
	RegisterUser(req *models.RegisterRequest, tenantID uint) (*models.User, error)
	Login(req *models.LoginRequest, tenantID uint, cfg *config.Config) (*models.LoginResponse, error)
	ValidateRegister(req *models.RegisterRequest) error
	ValidateLogin(req *models.LoginRequest) error
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

// RegisterUser registra un nuevo usuario
func (s *userService) RegisterUser(req *models.RegisterRequest, tenantID uint) (*models.User, error) {
	// Validar datos de entrada
	if err := s.ValidateRegister(req); err != nil {
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

	// Verificar si el documento ya existe
	exists, err = s.userRepo.ExistsByDocument(req.DocumentType, req.DocumentNumber)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, app_error.NewValidationError("document", "Ya existe un usuario con este documento")
	}

	// Validar lógica de negocio
	if err := s.validateRegisterBusinessRules(req); err != nil {
		return nil, err
	}

	// Hash de la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, app_error.NewDatabaseError("hash contraseña", err.Error())
	}

	// Crear modelo de usuario con tenant_id
	user := &models.User{
		TenantID:       tenantID,
		Name:           req.Name,
		Email:          req.Email,
		Phone:          req.Phone,
		DocumentType:   req.DocumentType,
		DocumentNumber: req.DocumentNumber,
		Password:       string(hashedPassword),
	}

	// Guardar en base de datos
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login autentica un usuario
func (s *userService) Login(req *models.LoginRequest, tenantID uint, cfg *config.Config) (*models.LoginResponse, error) {
	// Validar datos de entrada
	if err := s.ValidateLogin(req); err != nil {
		return nil, err
	}

	// Buscar usuario por email
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, app_error.NewValidationError("credentials", "Email o contraseña incorrectos")
	}

	// Verificar que el usuario pertenezca al tenant correcto
	if user.TenantID != tenantID {
		return nil, app_error.NewValidationError("credentials", "Usuario no pertenece a este tenant")
	}

	// Verificar contraseña
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, app_error.NewValidationError("credentials", "Email o contraseña incorrectos")
	}

	// Generar token JWT
	token, err := utils.GenerateAccessToken(user, cfg)
	if err != nil {
		return nil, app_error.NewDatabaseError("generar token", err.Error())
	}

	// Crear respuesta
	response := &models.LoginResponse{
		User:  user.ToResponse(),
		Token: token,
	}

	return response, nil
}

// ValidateRegister valida los datos de registro de un usuario
func (s *userService) ValidateRegister(req *models.RegisterRequest) error {
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

// ValidateLogin valida los datos de login de un usuario
func (s *userService) ValidateLogin(req *models.LoginRequest) error {
	if err := s.validator.Struct(req); err != nil {
		// Procesar errores de validación
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Tag() {
			case "required":
				validationErrors = append(validationErrors, err.Field()+" es requerido")
			case "email":
				validationErrors = append(validationErrors, "Email debe tener un formato válido")
			default:
				validationErrors = append(validationErrors, err.Field()+" no es válido")
			}
		}

		return app_error.NewValidationError("validación", validationErrors[0])
	}
	return nil
}

// validateRegisterBusinessRules valida reglas de negocio específicas para el registro
func (s *userService) validateRegisterBusinessRules(req *models.RegisterRequest) error {
	// Validar confirmación de contraseña
	if req.Password != req.PasswordConfirmation {
		return app_error.NewValidationError("password_confirmation", "Las contraseñas no coinciden")
	}

	// Validar tipo de documento
	if !models.IsValidDocumentType(req.DocumentType) {
		return app_error.NewValidationError("document_type", "Tipo de documento inválido")
	}

	return nil
}

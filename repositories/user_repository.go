package repositories

import (
	"errors"

	"loan-api/app_error"
	"loan-api/models"

	"gorm.io/gorm"
)

// UserRepository define la interfaz para operaciones de usuario
type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uint) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByDocument(documentType models.DocumentType, documentNumber string) (*models.User, error)
	ExistsByEmail(email string, tenantID uint) (bool, error)
	ExistsByDocument(documentType models.DocumentType, documentNumber string) (bool, error)
}

// userRepository implementa UserRepository
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository crea una nueva instancia del repositorio de usuarios
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

// Create crea un nuevo usuario
func (r *userRepository) Create(user *models.User) error {
	if err := r.db.Create(user).Error; err != nil {
		// Verificar si es un error de duplicado
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return app_error.ErrEmailExists
		}
		return app_error.NewDatabaseError("crear usuario", err.Error())
	}
	return nil
}

// GetByID obtiene un usuario por su ID
func (r *userRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, app_error.ErrUserNotFound
		}
		return nil, app_error.NewDatabaseError("obtener usuario", err.Error())
	}
	return &user, nil
}

// GetByEmail obtiene un usuario por su email
func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, app_error.ErrUserNotFound
		}
		return nil, app_error.NewDatabaseError("obtener usuario por email", err.Error())
	}
	return &user, nil
}

// GetByDocument obtiene un usuario por su tipo y número de documento
func (r *userRepository) GetByDocument(documentType models.DocumentType, documentNumber string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("document_type = ? AND document_number = ?", documentType, documentNumber).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, app_error.ErrUserNotFound
		}
		return nil, app_error.NewDatabaseError("obtener usuario por documento", err.Error())
	}
	return &user, nil
}

// ExistsByEmail verifica si existe un usuario con el email dado
func (r *userRepository) ExistsByEmail(email string, tenantID uint) (bool, error) {
	var count int64
	if err := r.db.Model(&models.User{}).Where("email = ? AND tenant_id = ?", email, tenantID).Count(&count).Error; err != nil {
		return false, app_error.NewDatabaseError("verificar email", err.Error())
	}
	return count > 0, nil
}

// ExistsByDocument verifica si existe un usuario con el documento dado
func (r *userRepository) ExistsByDocument(documentType models.DocumentType, documentNumber string) (bool, error) {
	var count int64
	if err := r.db.Model(&models.User{}).Where("document_type = ? AND document_number = ?", documentType, documentNumber).Count(&count).Error; err != nil {
		return false, app_error.NewDatabaseError("verificar documento", err.Error())
	}
	return count > 0, nil
}

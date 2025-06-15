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
	Update(user *models.User) error
	Delete(id uint) error
	List(limit, offset int) ([]models.User, int64, error)
	ExistsByEmail(email string) (bool, error)
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

// Update actualiza un usuario existente
func (r *userRepository) Update(user *models.User) error {
	if err := r.db.Save(user).Error; err != nil {
		// Verificar si es un error de duplicado
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return app_error.ErrEmailExists
		}
		return app_error.NewDatabaseError("actualizar usuario", err.Error())
	}
	return nil
}

// Delete elimina un usuario (soft delete)
func (r *userRepository) Delete(id uint) error {
	result := r.db.Delete(&models.User{}, id)
	if result.Error != nil {
		return app_error.NewDatabaseError("eliminar usuario", result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return app_error.ErrUserNotFound
	}
	return nil
}

// List obtiene una lista paginada de usuarios
func (r *userRepository) List(limit, offset int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	// Contar total de registros
	if err := r.db.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, app_error.NewDatabaseError("contar usuarios", err.Error())
	}

	// Obtener usuarios paginados
	if err := r.db.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, 0, app_error.NewDatabaseError("listar usuarios", err.Error())
	}

	return users, total, nil
}

// ExistsByEmail verifica si existe un usuario con el email dado
func (r *userRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	if err := r.db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, app_error.NewDatabaseError("verificar email", err.Error())
	}
	return count > 0, nil
}

// GetUsersWithLoans obtiene usuarios con sus préstamos (método adicional)
func (r *userRepository) GetUsersWithLoans(limit, offset int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	// Contar total de registros
	if err := r.db.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, app_error.NewDatabaseError("contar usuarios", err.Error())
	}

	// Obtener usuarios con préstamos
	if err := r.db.Preload("Loans").Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, 0, app_error.NewDatabaseError("listar usuarios con préstamos", err.Error())
	}

	return users, total, nil
}

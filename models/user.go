package models

import (
	"time"

	"gorm.io/gorm"
)

// User representa un usuario en el sistema
type User struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"type:varchar(100);not null" validate:"required,min=2,max=100"`
	Email       string         `json:"email" gorm:"type:varchar(100);uniqueIndex;not null" validate:"required,email"`
	Phone       string         `json:"phone" gorm:"type:varchar(20);not null" validate:"required,min=10,max=20"`
	Income      float64        `json:"income" gorm:"type:decimal(15,2);not null" validate:"required,min=0"`
	CreditScore int            `json:"credit_score" gorm:"type:int;not null" validate:"required,min=300,max=850"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relaciones
	Loans []Loan `json:"loans,omitempty" gorm:"foreignKey:UserID"`
}

// UserRequest representa la estructura para crear/actualizar usuarios
type UserRequest struct {
	Name        string  `json:"name" validate:"required,min=2,max=100"`
	Email       string  `json:"email" validate:"required,email"`
	Phone       string  `json:"phone" validate:"required,min=10,max=20"`
	Income      float64 `json:"income" validate:"required,min=0"`
	CreditScore int     `json:"credit_score" validate:"required,min=300,max=850"`
}

// UserResponse representa la respuesta del usuario (sin datos sensibles)
type UserResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	Income      float64   `json:"income"`
	CreditScore int       `json:"credit_score"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToResponse convierte un User a UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:          u.ID,
		Name:        u.Name,
		Email:       u.Email,
		Phone:       u.Phone,
		Income:      u.Income,
		CreditScore: u.CreditScore,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}

// TableName especifica el nombre de la tabla para GORM
func (User) TableName() string {
	return "users"
}

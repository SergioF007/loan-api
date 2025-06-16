package models

import (
	"time"

	"gorm.io/gorm"
)

// DocumentType representa los tipos de documento disponibles
type DocumentType string

const (
	DocumentTypeCedula           DocumentType = "cedula"
	DocumentTypePasaporte        DocumentType = "pasaporte"
	DocumentTypeTarjetaIdentidad DocumentType = "tarjeta_identidad"
)

// User representa un usuario en el sistema
type User struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	TenantID       uint           `json:"tenant_id" gorm:"not null;index"`
	Name           string         `json:"name" gorm:"type:varchar(100);not null" validate:"required,min=2,max=100"`
	Email          string         `json:"email" gorm:"type:varchar(100);uniqueIndex;not null" validate:"required,email"`
	Phone          string         `json:"phone" gorm:"type:varchar(20);not null" validate:"required,min=10,max=20"`
	DocumentType   DocumentType   `json:"document_type" gorm:"type:varchar(20);not null" validate:"required"`
	DocumentNumber string         `json:"document_number" gorm:"type:varchar(20);uniqueIndex;not null" validate:"required,min=5,max=20"`
	Password       string         `json:"-" gorm:"type:varchar(255);not null" validate:"required,min=8"`
	Income         *float64       `json:"income" gorm:"type:decimal(15,2);default:0"`
	IP             string         `json:"ip,omitempty" gorm:"type:varchar(45)"`
	CreatedAt      time.Time      `json:"created_at" gorm:"autoCreateTime:true"`
	UpdatedAt      time.Time      `json:"updated_at" gorm:"autoUpdateTime:true"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`

	// Relaciones
	Tenant Tenant `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`
	Loans  []Loan `json:"loans,omitempty" gorm:"foreignKey:UserID"`
}

// RegisterRequest representa la estructura para registrar usuarios
type RegisterRequest struct {
	Name                 string       `json:"name" validate:"required,min=2,max=100"`
	Email                string       `json:"email" validate:"required,email"`
	Phone                string       `json:"phone" validate:"required,min=10,max=20"`
	DocumentType         DocumentType `json:"document_type" validate:"required"`
	DocumentNumber       string       `json:"document_number" validate:"required,min=5,max=20"`
	Password             string       `json:"password" validate:"required,min=8"`
	PasswordConfirmation string       `json:"password_confirmation" validate:"required,min=8"`
}

// LoginRequest representa la estructura para login de usuarios
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UserResponse representa la respuesta del usuario (sin datos sensibles)
type UserResponse struct {
	ID             uint         `json:"id"`
	TenantID       uint         `json:"tenant_id"`
	Name           string       `json:"name"`
	Email          string       `json:"email"`
	Phone          string       `json:"phone"`
	DocumentType   DocumentType `json:"document_type"`
	DocumentNumber string       `json:"document_number"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}

// LoginResponse representa la respuesta del login
type LoginResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

// ToResponse convierte un User a UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:             u.ID,
		TenantID:       u.TenantID,
		Name:           u.Name,
		Email:          u.Email,
		Phone:          u.Phone,
		DocumentType:   u.DocumentType,
		DocumentNumber: u.DocumentNumber,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
	}
}

// IsValidDocumentType verifica si el tipo de documento es válido
func IsValidDocumentType(docType DocumentType) bool {
	return docType == DocumentTypeCedula ||
		docType == DocumentTypePasaporte ||
		docType == DocumentTypeTarjetaIdentidad
}

// TableName especifica el nombre de la tabla para GORM
func (User) TableName() string {
	return "users"
}

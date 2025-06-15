package models

import (
	"time"

	"gorm.io/gorm"
)

// LoanStatus define los posibles estados de un préstamo
type LoanStatus string

const (
	LoanStatusPending  LoanStatus = "pending"
	LoanStatusApproved LoanStatus = "approved"
	LoanStatusRejected LoanStatus = "rejected"
)

// Loan representa una solicitud de préstamo
type Loan struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	UserID       uint           `json:"user_id" gorm:"not null" validate:"required"`
	Amount       float64        `json:"amount" gorm:"type:decimal(15,2);not null" validate:"required,min=1"`
	Purpose      string         `json:"purpose" gorm:"type:varchar(500);not null" validate:"required,min=10,max=500"`
	Status       LoanStatus     `json:"status" gorm:"type:enum('pending','approved','rejected');default:'pending'" validate:"omitempty,oneof=pending approved rejected"`
	RequestDate  time.Time      `json:"request_date" gorm:"not null"`
	ApprovalDate *time.Time     `json:"approval_date,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	// Relación con User (evitamos referencia circular con puntero)
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// LoanRequest representa la estructura para crear solicitudes de préstamo
type LoanRequest struct {
	UserID  uint    `json:"user_id" validate:"required"`
	Amount  float64 `json:"amount" validate:"required,min=1"`
	Purpose string  `json:"purpose" validate:"required,min=10,max=500"`
}

// LoanStatusUpdateRequest representa la estructura para actualizar el estado
type LoanStatusUpdateRequest struct {
	Status LoanStatus `json:"status" validate:"required,oneof=pending approved rejected"`
}

// LoanResponse representa la respuesta de préstamo
type LoanResponse struct {
	ID           uint       `json:"id"`
	UserID       uint       `json:"user_id"`
	Amount       float64    `json:"amount"`
	Purpose      string     `json:"purpose"`
	Status       LoanStatus `json:"status"`
	RequestDate  time.Time  `json:"request_date"`
	ApprovalDate *time.Time `json:"approval_date,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	User         *User      `json:"user,omitempty"`
}

// ToResponse convierte un Loan a LoanResponse
func (l *Loan) ToResponse() LoanResponse {
	return LoanResponse{
		ID:           l.ID,
		UserID:       l.UserID,
		Amount:       l.Amount,
		Purpose:      l.Purpose,
		Status:       l.Status,
		RequestDate:  l.RequestDate,
		ApprovalDate: l.ApprovalDate,
		CreatedAt:    l.CreatedAt,
		UpdatedAt:    l.UpdatedAt,
		User:         l.User,
	}
}

// IsValidStatus verifica si el estado es válido
func (l *Loan) IsValidStatus(status LoanStatus) bool {
	return status == LoanStatusPending || status == LoanStatusApproved || status == LoanStatusRejected
}

// CanUpdateStatus verifica si se puede actualizar el estado
func (l *Loan) CanUpdateStatus(newStatus LoanStatus) bool {
	// No se puede cambiar el estado de préstamos ya aprobados o rechazados
	if l.Status == LoanStatusApproved || l.Status == LoanStatusRejected {
		return false
	}

	// Solo se puede cambiar de pending a approved/rejected
	if l.Status == LoanStatusPending && (newStatus == LoanStatusApproved || newStatus == LoanStatusRejected) {
		return true
	}

	return false
}

// TableName especifica el nombre de la tabla para GORM
func (Loan) TableName() string {
	return "loans"
}

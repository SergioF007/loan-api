package models

import (
	"time"

	"github.com/shopspring/decimal"
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
	ID             uint            `json:"id" gorm:"primaryKey"`
	LoanTypeID     uint            `json:"loan_type_id" gorm:"not null;index"`
	LoanType       LoanType        `json:"loan_type"`
	UserID         uint            `json:"user_id" gorm:"not null;index"`
	User           User            `json:"user"`
	Status         string          `json:"status" gorm:"size:50;default:'pending'"`
	Observation    string          `json:"observation" gorm:"type:text"`
	AmountApproved decimal.Decimal `json:"amount_approved" gorm:"type:decimal(13,2);default:0"`
	Data           []LoanData      `json:"data"`
	CreatedAt      time.Time       `json:"created_at" gorm:"autoCreateTime:true"`
	UpdatedAt      time.Time       `json:"updated_at" gorm:"autoUpdateTime:true"`
	DeletedAt      gorm.DeletedAt  `json:"-" gorm:"index"`
}

// LoanData representa los datos dinámicos de una solicitud de préstamo
type LoanData struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	LoanID    uint           `json:"loan_id" gorm:"not null;index"`
	Loan      Loan           `json:"-"`
	FormID    uint           `json:"form_id" gorm:"not null;index"`
	Form      LoanTypeForm   `json:"-"`
	Key       string         `json:"key" gorm:"size:255;not null"`
	Value     string         `json:"value" gorm:"type:text"`
	Index     uint           `json:"index" gorm:"default:0"`
	CreatedAt time.Time      `json:"-" gorm:"autoCreateTime:true"`
	UpdatedAt time.Time      `json:"-" gorm:"autoUpdateTime:true"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Requests para la nueva estructura
// CreateLoanRequest representa la estructura para crear solicitudes de préstamo
type CreateLoanRequest struct {
	LoanTypeID uint `json:"loan_type_id" validate:"required"`
}

// SaveLoanDataRequest representa la estructura para guardar datos de préstamo
type SaveLoanDataRequest struct {
	LoanID uint                  `json:"loan_id" validate:"required"`
	Data   []LoanDataItemRequest `json:"data" validate:"required,dive"`
}

// LoanDataItemRequest representa un item de datos de préstamo
type LoanDataItemRequest struct {
	FormID uint   `json:"form_id" validate:"required"`
	Key    string `json:"key" validate:"required"`
	Value  string `json:"value" validate:"required"`
	Index  uint   `json:"index"`
}

// Responses para la nueva estructura
// LoanResponse representa la respuesta de préstamo
type LoanResponse struct {
	ID             uint               `json:"id"`
	LoanTypeID     uint               `json:"loan_type_id"`
	LoanType       LoanTypeResponse   `json:"loan_type"`
	UserID         uint               `json:"user_id"`
	User           UserResponse       `json:"user"`
	Status         string             `json:"status"`
	Observation    string             `json:"observation"`
	AmountApproved decimal.Decimal    `json:"amount_approved"`
	Data           []LoanDataResponse `json:"data"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
}

// LoanDataResponse representa la respuesta de datos de préstamo
type LoanDataResponse struct {
	ID     uint   `json:"id"`
	FormID uint   `json:"form_id"`
	Key    string `json:"key"`
	Value  string `json:"value"`
	Index  uint   `json:"index"`
}

// ToResponse convierte un Loan a LoanResponse
func (l *Loan) ToResponse() LoanResponse {
	dataResponse := make([]LoanDataResponse, len(l.Data))
	for i, data := range l.Data {
		dataResponse[i] = LoanDataResponse{
			ID:     data.ID,
			FormID: data.FormID,
			Key:    data.Key,
			Value:  data.Value,
			Index:  data.Index,
		}
	}

	return LoanResponse{
		ID:             l.ID,
		LoanTypeID:     l.LoanTypeID,
		UserID:         l.UserID,
		Status:         l.Status,
		Observation:    l.Observation,
		AmountApproved: l.AmountApproved,
		Data:           dataResponse,
		CreatedAt:      l.CreatedAt,
		UpdatedAt:      l.UpdatedAt,
	}
}

// TableName especifica el nombre de la tabla para GORM
func (Loan) TableName() string {
	return "loans"
}

// TableName especifica el nombre de la tabla para GORM
func (LoanData) TableName() string {
	return "loan_data"
}

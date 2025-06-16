package models

import (
	"time"

	"gorm.io/gorm"
)

// Tenant representa una entidad crediticia
type Tenant struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"size:255;not null"`
	Code        string         `json:"code" gorm:"size:50;uniqueIndex;not null"`
	Description string         `json:"description" gorm:"type:text"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	Config      string         `json:"config" gorm:"type:json"`
	LoanTypes   []LoanType     `json:"loan_types,omitempty"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime:true"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime:true"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// TenantResponse representa la respuesta del tenant
type TenantResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

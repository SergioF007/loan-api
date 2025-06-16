package models

import (
	"time"

	"gorm.io/gorm"
)

// LoanType representa un tipo de crédito
type LoanType struct {
	ID          uint              `json:"id" gorm:"primaryKey"`
	TenantID    uint              `json:"tenant_id" gorm:"not null;index"`
	Tenant      Tenant            `json:"tenant,omitempty"`
	Name        string            `json:"name" gorm:"size:255;not null"`
	Code        string            `json:"code" gorm:"size:50;not null"`
	Description string            `json:"description" gorm:"type:text"`
	IsActive    bool              `json:"is_active" gorm:"default:true"`
	MinAmount   float64           `json:"min_amount" gorm:"type:decimal(15,2);default:0"`
	MaxAmount   float64           `json:"max_amount" gorm:"type:decimal(15,2);default:0"`
	Versions    []LoanTypeVersion `json:"versions,omitempty"`
	CreatedAt   time.Time         `json:"created_at" gorm:"autoCreateTime:true"`
	UpdatedAt   time.Time         `json:"updated_at" gorm:"autoUpdateTime:true"`
	DeletedAt   gorm.DeletedAt    `json:"-" gorm:"index"`
}

// LoanTypeVersion representa una versión de un tipo de crédito
type LoanTypeVersion struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	LoanTypeID  uint           `json:"loan_type_id" gorm:"not null;index"`
	LoanType    LoanType       `json:"loan_type,omitempty"`
	Version     string         `json:"version" gorm:"size:50;not null"`
	Description string         `json:"description" gorm:"type:text"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	IsDefault   bool           `json:"is_default" gorm:"default:false"`
	Config      string         `json:"config" gorm:"type:json"`
	Forms       []LoanTypeForm `json:"forms,omitempty"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime:true"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime:true"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// LoanTypeForm representa un formulario disponible para un tipo de crédito
type LoanTypeForm struct {
	ID                uint                       `json:"id" gorm:"primaryKey"`
	LoanTypeVersionID uint                       `json:"loan_type_version_id" gorm:"not null;index"`
	LoanTypeVersion   LoanTypeVersion            `json:"loan_type_version,omitempty"`
	Label             string                     `json:"label" gorm:"size:255;not null"`
	Code              string                     `json:"code" gorm:"size:50;not null"`
	Description       string                     `json:"description" gorm:"type:text"`
	Order             int                        `json:"order" gorm:"default:0"`
	IsRequired        bool                       `json:"is_required" gorm:"default:false"`
	IsActive          bool                       `json:"is_active" gorm:"default:true"`
	Config            string                     `json:"config" gorm:"type:json"`
	FormInputs        []LoanTypeVersionFormInput `json:"form_inputs,omitempty" gorm:"foreignKey:LoanTypeFormID"`
	CreatedAt         time.Time                  `json:"created_at" gorm:"autoCreateTime:true"`
	UpdatedAt         time.Time                  `json:"updated_at" gorm:"autoUpdateTime:true"`
	DeletedAt         gorm.DeletedAt             `json:"-" gorm:"index"`
}

// LoanTypeVersionFormInput representa un input para la configuración de formularios
type LoanTypeVersionFormInput struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	LoanTypeFormID  uint           `json:"loan_type_form_id,omitempty" gorm:"index"`
	LoanTypeForm    LoanTypeForm   `json:"loan_type_form,omitempty"`
	Label           string         `json:"label" gorm:"size:255;not null"`
	Code            string         `json:"code" gorm:"size:50;not null"`
	InputType       string         `json:"input_type" gorm:"size:50;not null"` // text, number, email, select, etc.
	Placeholder     string         `json:"placeholder" gorm:"size:255"`
	DefaultValue    string         `json:"default_value" gorm:"type:text"`
	ValidationRules string         `json:"validation_rules" gorm:"type:json"`
	Options         string         `json:"options" gorm:"type:json"`
	Order           int            `json:"order" gorm:"default:0"`
	IsRequired      bool           `json:"is_required" gorm:"default:false"`
	IsActive        bool           `json:"is_active" gorm:"default:true"`
	Config          string         `json:"config" gorm:"type:json"`
	CreatedAt       time.Time      `json:"created_at" gorm:"autoCreateTime:true"`
	UpdatedAt       time.Time      `json:"updated_at" gorm:"autoUpdateTime:true"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

// LoanTypeResponse representa la respuesta de un tipo de crédito con formularios
type LoanTypeResponse struct {
	ID          uint                    `json:"id"`
	Name        string                  `json:"name"`
	Code        string                  `json:"code"`
	Description string                  `json:"description"`
	MinAmount   float64                 `json:"min_amount"`
	MaxAmount   float64                 `json:"max_amount"`
	Version     LoanTypeVersionResponse `json:"version"`
}

// LoanTypeVersionResponse representa la respuesta de una versión con formularios
type LoanTypeVersionResponse struct {
	ID          uint                   `json:"id"`
	Version     string                 `json:"version"`
	Description string                 `json:"description"`
	Forms       []LoanTypeFormResponse `json:"forms"`
}

// LoanTypeFormResponse representa la respuesta de un formulario
type LoanTypeFormResponse struct {
	ID          uint                               `json:"id"`
	Label       string                             `json:"label"`
	Code        string                             `json:"code"`
	Description string                             `json:"description"`
	Order       int                                `json:"order"`
	IsRequired  bool                               `json:"is_required"`
	FormInputs  []LoanTypeVersionFormInputResponse `json:"form_inputs"`
}

// LoanTypeVersionFormInputResponse representa la respuesta de un input
type LoanTypeVersionFormInputResponse struct {
	ID              uint   `json:"id"`
	Label           string `json:"label"`
	Code            string `json:"code"`
	InputType       string `json:"input_type"`
	Placeholder     string `json:"placeholder"`
	DefaultValue    string `json:"default_value"`
	ValidationRules string `json:"validation_rules"`
	Options         string `json:"options"`
	Order           int    `json:"order"`
	IsRequired      bool   `json:"is_required"`
}

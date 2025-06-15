package utils

import (
	"net/http"

	"loan-api/app_error"

	"github.com/gin-gonic/gin"
)

// APIResponse representa la estructura estándar de respuesta de la API
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

// ErrorInfo contiene información detallada del error
type ErrorInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// PaginatedResponse representa una respuesta paginada
type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

// Pagination contiene información de paginación
type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// SuccessResponse envía una respuesta exitosa
func SuccessResponse(c *gin.Context, code int, message string, data interface{}) {
	response := APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
	c.JSON(code, response)
}

// ErrorResponse envía una respuesta de error
func ErrorResponse(c *gin.Context, err error) {
	// Si es un AppError, usar su información
	if appErr, ok := app_error.IsAppError(err); ok {
		response := APIResponse{
			Success: false,
			Message: "Error en la solicitud",
			Error: &ErrorInfo{
				Code:    appErr.Code,
				Message: appErr.Message,
				Details: appErr.Details,
			},
		}
		c.JSON(appErr.Code, response)
		return
	}

	// Error genérico
	response := APIResponse{
		Success: false,
		Message: "Error interno del servidor",
		Error: &ErrorInfo{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		},
	}
	c.JSON(http.StatusInternalServerError, response)
}

// ValidationErrorResponse envía una respuesta de error de validación
func ValidationErrorResponse(c *gin.Context, errors []string) {
	response := APIResponse{
		Success: false,
		Message: "Error de validación",
		Error: &ErrorInfo{
			Code:    http.StatusBadRequest,
			Message: "Los datos proporcionados no son válidos",
			Details: joinErrors(errors),
		},
	}
	c.JSON(http.StatusBadRequest, response)
}

// PaginatedSuccessResponse envía una respuesta exitosa con paginación
func PaginatedSuccessResponse(c *gin.Context, message string, data interface{}, pagination Pagination) {
	response := PaginatedResponse{
		Success:    true,
		Message:    message,
		Data:       data,
		Pagination: pagination,
	}
	c.JSON(http.StatusOK, response)
}

// CreatedResponse envía una respuesta de recurso creado
func CreatedResponse(c *gin.Context, message string, data interface{}) {
	SuccessResponse(c, http.StatusCreated, message, data)
}

// UpdatedResponse envía una respuesta de recurso actualizado
func UpdatedResponse(c *gin.Context, message string, data interface{}) {
	SuccessResponse(c, http.StatusOK, message, data)
}

// DeletedResponse envía una respuesta de recurso eliminado
func DeletedResponse(c *gin.Context, message string) {
	SuccessResponse(c, http.StatusOK, message, nil)
}

// NotFoundResponse envía una respuesta de recurso no encontrado
func NotFoundResponse(c *gin.Context, message string) {
	ErrorResponse(c, app_error.NewAppError(http.StatusNotFound, message))
}

// BadRequestResponse envía una respuesta de solicitud incorrecta
func BadRequestResponse(c *gin.Context, message string) {
	ErrorResponse(c, app_error.NewAppError(http.StatusBadRequest, message))
}

// UnauthorizedResponse envía una respuesta de no autorizado
func UnauthorizedResponse(c *gin.Context, message string) {
	ErrorResponse(c, app_error.NewAppError(http.StatusUnauthorized, message))
}

// ForbiddenResponse envía una respuesta de prohibido
func ForbiddenResponse(c *gin.Context, message string) {
	ErrorResponse(c, app_error.NewAppError(http.StatusForbidden, message))
}

// InternalServerErrorResponse envía una respuesta de error interno
func InternalServerErrorResponse(c *gin.Context, message string) {
	ErrorResponse(c, app_error.NewAppError(http.StatusInternalServerError, message))
}

// NewPagination crea una nueva estructura de paginación
func NewPagination(page, limit int, total int64) Pagination {
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return Pagination{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}

// joinErrors une múltiples errores en una cadena
func joinErrors(errors []string) string {
	if len(errors) == 0 {
		return ""
	}

	result := errors[0]
	for i := 1; i < len(errors); i++ {
		result += "; " + errors[i]
	}

	return result
}

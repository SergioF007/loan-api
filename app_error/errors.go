package app_error

import (
	"fmt"
	"net/http"
)

// AppError representa un error personalizado de la aplicación
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Error implementa la interfaz error
func (e *AppError) Error() string {
	return e.Message
}

// NewAppError crea un nuevo error personalizado
func NewAppError(code int, message string, details ...string) *AppError {
	appErr := &AppError{
		Code:    code,
		Message: message,
	}

	if len(details) > 0 {
		appErr.Details = details[0]
	}

	return appErr
}

// Errores comunes predefinidos
var (
	// Errores de validación
	ErrValidationFailed = NewAppError(http.StatusBadRequest, "Error de validación")
	ErrInvalidJSON      = NewAppError(http.StatusBadRequest, "Formato JSON inválido")
	ErrRequiredField    = NewAppError(http.StatusBadRequest, "Campo requerido faltante")

	// Errores de base de datos
	ErrDatabaseConnection = NewAppError(http.StatusInternalServerError, "Error de conexión a la base de datos")
	ErrDatabaseQuery      = NewAppError(http.StatusInternalServerError, "Error en la consulta de base de datos")
	ErrRecordNotFound     = NewAppError(http.StatusNotFound, "Registro no encontrado")
	ErrDuplicateRecord    = NewAppError(http.StatusConflict, "Registro duplicado")

	// Errores de usuarios
	ErrUserNotFound = NewAppError(http.StatusNotFound, "Usuario no encontrado")
	ErrUserExists   = NewAppError(http.StatusConflict, "El usuario ya existe")
	ErrInvalidUser  = NewAppError(http.StatusBadRequest, "Datos de usuario inválidos")
	ErrEmailExists  = NewAppError(http.StatusConflict, "El email ya está registrado")

	// Errores de préstamos
	ErrLoanNotFound       = NewAppError(http.StatusNotFound, "Préstamo no encontrado")
	ErrInvalidLoan        = NewAppError(http.StatusBadRequest, "Datos de préstamo inválidos")
	ErrLoanStatusInvalid  = NewAppError(http.StatusBadRequest, "Estado de préstamo inválido")
	ErrCannotUpdateStatus = NewAppError(http.StatusBadRequest, "No se puede actualizar el estado del préstamo")
	ErrInsufficientCredit = NewAppError(http.StatusBadRequest, "Puntaje crediticio insuficiente")
	ErrInvalidAmount      = NewAppError(http.StatusBadRequest, "Monto inválido")

	// Errores de autenticación
	ErrUnauthorized = NewAppError(http.StatusUnauthorized, "No autorizado")
	ErrForbidden    = NewAppError(http.StatusForbidden, "Acceso prohibido")
	ErrInvalidToken = NewAppError(http.StatusUnauthorized, "Token inválido")
	ErrTokenExpired = NewAppError(http.StatusUnauthorized, "Token expirado")

	// Errores del servidor
	ErrInternalServer     = NewAppError(http.StatusInternalServerError, "Error interno del servidor")
	ErrServiceUnavailable = NewAppError(http.StatusServiceUnavailable, "Servicio no disponible")
)

// NewValidationError crea un error de validación con detalles específicos
func NewValidationError(field string, message string) *AppError {
	return NewAppError(http.StatusBadRequest,
		fmt.Sprintf("Error de validación en el campo '%s'", field),
		message)
}

// NewDatabaseError crea un error de base de datos con detalles
func NewDatabaseError(operation string, details string) *AppError {
	return NewAppError(http.StatusInternalServerError,
		fmt.Sprintf("Error en la operación de base de datos: %s", operation),
		details)
}

// NewBusinessError crea un error de lógica de negocio
func NewBusinessError(message string, details string) *AppError {
	return NewAppError(http.StatusBadRequest, message, details)
}

// IsAppError verifica si un error es de tipo AppError
func IsAppError(err error) (*AppError, bool) {
	appErr, ok := err.(*AppError)
	return appErr, ok
}

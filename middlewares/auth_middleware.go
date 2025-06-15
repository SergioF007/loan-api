package middlewares

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"loan-api/app_error"
	"loan-api/config"
	"loan-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware define la estructura del middleware de autenticación
type AuthMiddleware struct {
	config *config.Config
}

// NewAuthMiddleware crea una nueva instancia del middleware de autenticación
func NewAuthMiddleware(cfg *config.Config) *AuthMiddleware {
	return &AuthMiddleware{
		config: cfg,
	}
}

// RequireAuth middleware que requiere autenticación JWT
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener el token del header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.UnauthorizedResponse(c, "Token de autenticación requerido")
			c.Abort()
			return
		}

		// Verificar formato del header (Bearer token)
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.UnauthorizedResponse(c, "Formato de token inválido")
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validar el token
		claims, err := m.validateToken(tokenString)
		if err != nil {
			utils.ErrorResponse(c, err)
			c.Abort()
			return
		}

		// Guardar información del usuario en el contexto
		c.Set("user_id", claims["user_id"])
		c.Set("user_email", claims["email"])

		c.Next()
	}
}

// OptionalAuth middleware que permite autenticación opcional
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		tokenString := parts[1]
		claims, err := m.validateToken(tokenString)
		if err == nil {
			c.Set("user_id", claims["user_id"])
			c.Set("user_email", claims["email"])
		}

		c.Next()
	}
}

// CORS middleware para manejar CORS
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin,Content-Type,Authorization,Accept,User-Agent,Cache-Control,Pragma")
		c.Header("Access-Control-Expose-Headers", "Content-Length,Access-Control-Allow-Origin,Access-Control-Allow-Headers,Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// LoggerMiddleware middleware personalizado para logging
func LoggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[%s] %s %s %d %s %s\n",
			param.TimeStamp.Format("2006-01-02 15:04:05"),
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency,
			param.ClientIP,
		)
	})
}

// ErrorHandlerMiddleware middleware para manejo centralizado de errores
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Procesar errores que puedan haber ocurrido
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			utils.ErrorResponse(c, err)
		}
	}
}

// validateToken valida un token JWT y retorna los claims
func (m *AuthMiddleware) validateToken(tokenString string) (jwt.MapClaims, error) {
	// Parsear el token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verificar el método de firma
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, app_error.NewAppError(http.StatusUnauthorized, "Método de firma inválido")
		}
		return []byte(m.config.JWTSecret), nil
	})

	if err != nil {
		return nil, app_error.ErrInvalidToken
	}

	// Verificar que el token sea válido
	if !token.Valid {
		return nil, app_error.ErrInvalidToken
	}

	// Obtener los claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, app_error.ErrInvalidToken
	}

	return claims, nil
}

// GenerateToken genera un nuevo token JWT (método de utilidad)
func (m *AuthMiddleware) GenerateToken(userID uint, email string) (string, error) {
	// Crear los claims
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(time.Hour * time.Duration(m.config.JWTExpirationHours)).Unix(),
		"iat":     time.Now().Unix(),
	}

	// Crear el token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Firmar el token
	tokenString, err := token.SignedString([]byte(m.config.JWTSecret))
	if err != nil {
		return "", app_error.NewAppError(http.StatusInternalServerError, "Error al generar token")
	}

	return tokenString, nil
}

// GetUserIDFromContext obtiene el ID del usuario del contexto
func GetUserIDFromContext(c *gin.Context) (uint, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, app_error.ErrUnauthorized
	}

	// Convertir a uint
	switch v := userID.(type) {
	case uint:
		return v, nil
	case float64:
		return uint(v), nil
	case int:
		return uint(v), nil
	default:
		return 0, app_error.ErrUnauthorized
	}
}

// GetUserEmailFromContext obtiene el email del usuario del contexto
func GetUserEmailFromContext(c *gin.Context) (string, error) {
	email, exists := c.Get("user_email")
	if !exists {
		return "", app_error.ErrUnauthorized
	}

	emailStr, ok := email.(string)
	if !ok {
		return "", app_error.ErrUnauthorized
	}

	return emailStr, nil
}

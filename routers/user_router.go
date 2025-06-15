package routers

import (
	"loan-api/controllers"
	"loan-api/middlewares"

	"github.com/gin-gonic/gin"
)

// UserRouter configura las rutas relacionadas con usuarios
type UserRouter struct {
	userController *controllers.UserController
	authMiddleware *middlewares.AuthMiddleware
}

// NewUserRouter crea una nueva instancia del router de usuarios
func NewUserRouter(userController *controllers.UserController, authMiddleware *middlewares.AuthMiddleware) *UserRouter {
	return &UserRouter{
		userController: userController,
		authMiddleware: authMiddleware,
	}
}

// SetupRoutes configura todas las rutas de usuarios
func (r *UserRouter) SetupRoutes(router *gin.RouterGroup) {
	// Grupo de rutas para usuarios
	users := router.Group("/users")
	{
		// Rutas públicas (sin autenticación)
		users.POST("", r.userController.CreateUser)           // POST /api/v1/users - Crear usuario
		users.GET("/search", r.userController.GetUserByEmail) // GET /api/v1/users/search?email=... - Buscar por email

		// Rutas que requieren autenticación opcional (para algunas funcionalidades)
		users.GET("", r.authMiddleware.OptionalAuth(), r.userController.ListUsers)   // GET /api/v1/users - Listar usuarios
		users.GET("/:id", r.authMiddleware.OptionalAuth(), r.userController.GetUser) // GET /api/v1/users/:id - Obtener usuario

		// Rutas que requieren autenticación (para operaciones sensibles)
		users.PUT("/:id", r.authMiddleware.RequireAuth(), r.userController.UpdateUser)                          // PUT /api/v1/users/:id - Actualizar usuario
		users.DELETE("/:id", r.authMiddleware.RequireAuth(), r.userController.DeleteUser)                       // DELETE /api/v1/users/:id - Eliminar usuario
		users.GET("/:id/credit-summary", r.authMiddleware.RequireAuth(), r.userController.GetUserCreditSummary) // GET /api/v1/users/:id/credit-summary - Resumen crediticio
	}
}

// SetupPublicRoutes configura rutas públicas de usuarios (sin autenticación)
func (r *UserRouter) SetupPublicRoutes(router *gin.RouterGroup) {
	// Grupo de rutas públicas para usuarios
	publicUsers := router.Group("/users")
	{
		publicUsers.POST("", r.userController.CreateUser)           // Crear usuario
		publicUsers.GET("/search", r.userController.GetUserByEmail) // Buscar por email
	}
}

// SetupProtectedRoutes configura rutas protegidas de usuarios (requieren autenticación)
func (r *UserRouter) SetupProtectedRoutes(router *gin.RouterGroup) {
	// Aplicar middleware de autenticación al grupo
	protected := router.Group("/users", r.authMiddleware.RequireAuth())
	{
		protected.GET("", r.userController.ListUsers)                               // Listar usuarios
		protected.GET("/:id", r.userController.GetUser)                             // Obtener usuario
		protected.PUT("/:id", r.userController.UpdateUser)                          // Actualizar usuario
		protected.DELETE("/:id", r.userController.DeleteUser)                       // Eliminar usuario
		protected.GET("/:id/credit-summary", r.userController.GetUserCreditSummary) // Resumen crediticio
	}
}

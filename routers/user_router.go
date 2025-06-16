package routers

import (
	"loan-api/controllers"

	"github.com/gin-gonic/gin"
)

// UserRouter configura las rutas relacionadas con usuarios
type UserRouter struct {
	userController *controllers.UserController
}

// NewUserRouter crea una nueva instancia del router de usuarios
func NewUserRouter(userController *controllers.UserController) *UserRouter {
	return &UserRouter{
		userController: userController,
	}
}

// Setup configura todas las rutas de usuarios
func (r *UserRouter) Setup(router *gin.RouterGroup) {
	// Rutas de autenticación
	auth := router.Group("/auth")
	{
		auth.POST("/register", r.userController.RegisterUser) // Registro de usuario
		auth.POST("/login", r.userController.Login)           // Login de usuario
	}
}

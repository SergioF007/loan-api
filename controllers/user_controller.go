package controllers

import (
	"loan-api/config"
	"loan-api/models"
	"loan-api/services"
	"loan-api/utils"
	"log"

	"github.com/gin-gonic/gin"
)

// UserController maneja las operaciones relacionadas con usuarios
type UserController struct {
	userService services.UserService
	config      *config.Config
}

// NewUserController crea una nueva instancia del controlador de usuarios
func NewUserController(userService services.UserService, cfg *config.Config) *UserController {
	return &UserController{
		userService: userService,
		config:      cfg,
	}
}

// RegisterUser godoc
// @Summary Registrar un nuevo usuario
// @Description Registra un nuevo usuario en el sistema
// @Tags auth
// @Accept json
// @Produce json
// @Param X-Tenant-ID header string true "ID del tenant"
// @Param user body models.RegisterRequest true "Datos del usuario"
// @Success 201 {object} utils.APIResponse{data=models.UserResponse}
// @Failure 400 {object} utils.APIResponse
// @Failure 409 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /auth/register [post]
func (ctrl *UserController) RegisterUser(c *gin.Context) {
	log.Println("UserController::RegisterUser was invoke")

	var req models.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Formato JSON inválido")
		return
	}

	tenantID, exists := c.Get("tenant_id")
	if !exists {
		utils.BadRequestResponse(c, "ID de tenant requerido")
		return
	}

	tenantIDUint, ok := tenantID.(uint)
	if !ok {
		utils.BadRequestResponse(c, "ID de tenant inválido")
		return
	}

	// Registrar usuario con tenant_id
	user, err := ctrl.userService.RegisterUser(&req, tenantIDUint)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	// Retornar respuesta exitosa
	utils.CreatedResponse(c, "Usuario registrado exitosamente", user.ToResponse())
}

// Login godoc
// @Summary Iniciar sesión
// @Description Autentica un usuario y retorna un token JWT
// @Tags auth
// @Accept json
// @Produce json
// @Param X-Tenant-ID header string true "ID del tenant"
// @Param credentials body models.LoginRequest true "Credenciales del usuario"
// @Success 200 {object} utils.APIResponse{data=models.LoginResponse}
// @Failure 400 {object} utils.APIResponse
// @Failure 401 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /auth/login [post]
func (ctrl *UserController) Login(c *gin.Context) {
	log.Println("UserController::Login was invoke")

	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Formato JSON inválido")
		return
	}

	tenantID, exists := c.Get("tenant_id")
	if !exists {
		utils.BadRequestResponse(c, "ID de tenant requerido")
		return
	}

	tenantIDUint, ok := tenantID.(uint)
	if !ok {
		utils.BadRequestResponse(c, "ID de tenant inválido")
		return
	}

	// Autenticar usuario con validación de tenant
	response, err := ctrl.userService.Login(&req, tenantIDUint, ctrl.config)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, 200, "Login exitoso", response)
}

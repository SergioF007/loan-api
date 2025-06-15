package controllers

import (
	"net/http"
	"strconv"

	"loan-api/models"
	"loan-api/services"
	"loan-api/utils"

	"github.com/gin-gonic/gin"
)

// UserController maneja las operaciones relacionadas con usuarios
type UserController struct {
	userService services.UserService
}

// NewUserController crea una nueva instancia del controlador de usuarios
func NewUserController(userService services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// CreateUser godoc
// @Summary Crear un nuevo usuario
// @Description Crea un nuevo usuario en el sistema
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.UserRequest true "Datos del usuario"
// @Success 201 {object} utils.APIResponse{data=models.UserResponse}
// @Failure 400 {object} utils.APIResponse
// @Failure 409 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /api/v1/users [post]
func (ctrl *UserController) CreateUser(c *gin.Context) {
	var req models.UserRequest

	// Parsear JSON del request
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Formato JSON inválido")
		return
	}

	// Crear usuario
	user, err := ctrl.userService.CreateUser(&req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	// Retornar respuesta exitosa
	utils.CreatedResponse(c, "Usuario creado exitosamente", user.ToResponse())
}

// GetUser godoc
// @Summary Obtener usuario por ID
// @Description Obtiene la información de un usuario por su ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID del usuario"
// @Success 200 {object} utils.APIResponse{data=models.UserResponse}
// @Failure 400 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /api/v1/users/{id} [get]
func (ctrl *UserController) GetUser(c *gin.Context) {
	// Obtener ID del parámetro URL
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID de usuario inválido")
		return
	}

	// Obtener usuario
	user, err := ctrl.userService.GetUserByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	// Retornar respuesta exitosa
	utils.SuccessResponse(c, http.StatusOK, "Usuario obtenido exitosamente", user.ToResponse())
}

// UpdateUser godoc
// @Summary Actualizar usuario por ID
// @Description Actualiza la información de un usuario existente
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID del usuario"
// @Param user body models.UserRequest true "Datos actualizados del usuario"
// @Success 200 {object} utils.APIResponse{data=models.UserResponse}
// @Failure 400 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Failure 409 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /api/v1/users/{id} [put]
func (ctrl *UserController) UpdateUser(c *gin.Context) {
	// Obtener ID del parámetro URL
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID de usuario inválido")
		return
	}

	var req models.UserRequest

	// Parsear JSON del request
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Formato JSON inválido")
		return
	}

	// Actualizar usuario
	user, err := ctrl.userService.UpdateUser(uint(id), &req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	// Retornar respuesta exitosa
	utils.UpdatedResponse(c, "Usuario actualizado exitosamente", user.ToResponse())
}

// DeleteUser godoc
// @Summary Eliminar usuario por ID
// @Description Elimina un usuario del sistema (soft delete)
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID del usuario"
// @Success 200 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /api/v1/users/{id} [delete]
func (ctrl *UserController) DeleteUser(c *gin.Context) {
	// Obtener ID del parámetro URL
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID de usuario inválido")
		return
	}

	// Eliminar usuario
	err = ctrl.userService.DeleteUser(uint(id))
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	// Retornar respuesta exitosa
	utils.DeletedResponse(c, "Usuario eliminado exitosamente")
}

// ListUsers godoc
// @Summary Listar usuarios
// @Description Obtiene una lista paginada de usuarios
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Número de página" default(1)
// @Param limit query int false "Elementos por página" default(20)
// @Success 200 {object} utils.PaginatedResponse{data=[]models.UserResponse}
// @Failure 400 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /api/v1/users [get]
func (ctrl *UserController) ListUsers(c *gin.Context) {
	// Obtener parámetros de paginación
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	// Validar parámetros
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	// Obtener usuarios
	users, total, err := ctrl.userService.ListUsers(page, limit)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	// Convertir a respuestas
	userResponses := make([]models.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = user.ToResponse()
	}

	// Crear paginación
	pagination := utils.NewPagination(page, limit, total)

	// Retornar respuesta paginada
	utils.PaginatedSuccessResponse(c, "Usuarios obtenidos exitosamente", userResponses, pagination)
}

// GetUserByEmail godoc
// @Summary Obtener usuario por email
// @Description Obtiene la información de un usuario por su email
// @Tags users
// @Accept json
// @Produce json
// @Param email query string true "Email del usuario"
// @Success 200 {object} utils.APIResponse{data=models.UserResponse}
// @Failure 400 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /api/v1/users/search [get]
func (ctrl *UserController) GetUserByEmail(c *gin.Context) {
	// Obtener email del query parameter
	email := c.Query("email")
	if email == "" {
		utils.BadRequestResponse(c, "Email es requerido")
		return
	}

	// Obtener usuario
	user, err := ctrl.userService.GetUserByEmail(email)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	// Retornar respuesta exitosa
	utils.SuccessResponse(c, http.StatusOK, "Usuario encontrado exitosamente", user.ToResponse())
}

// GetUserCreditSummary godoc
// @Summary Obtener resumen crediticio del usuario
// @Description Obtiene un resumen del estado crediticio de un usuario
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID del usuario"
// @Success 200 {object} utils.APIResponse{data=map[string]interface{}}
// @Failure 400 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /api/v1/users/{id}/credit-summary [get]
func (ctrl *UserController) GetUserCreditSummary(c *gin.Context) {
	// Obtener ID del parámetro URL
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID de usuario inválido")
		return
	}

	// Obtener resumen crediticio
	summary, err := ctrl.userService.GetUserCreditSummary(uint(id))
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	// Retornar respuesta exitosa
	utils.SuccessResponse(c, http.StatusOK, "Resumen crediticio obtenido exitosamente", summary)
}

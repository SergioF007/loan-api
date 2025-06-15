package controllers

import (
	"net/http"
	"strconv"

	"loan-api/models"
	"loan-api/services"
	"loan-api/utils"

	"github.com/gin-gonic/gin"
)

// LoanController maneja las operaciones relacionadas con préstamos
type LoanController struct {
	loanService services.LoanService
}

// NewLoanController crea una nueva instancia del controlador de préstamos
func NewLoanController(loanService services.LoanService) *LoanController {
	return &LoanController{
		loanService: loanService,
	}
}

// CreateLoan godoc
// @Summary Crear una nueva solicitud de préstamo
// @Description Crea una nueva solicitud de préstamo en el sistema
// @Tags loans
// @Accept json
// @Produce json
// @Param loan body models.LoanRequest true "Datos del préstamo"
// @Success 201 {object} utils.APIResponse{data=models.LoanResponse}
// @Failure 400 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /api/v1/loans [post]
func (ctrl *LoanController) CreateLoan(c *gin.Context) {
	var req models.LoanRequest

	// Parsear JSON del request
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Formato JSON inválido")
		return
	}

	// Crear préstamo
	loan, err := ctrl.loanService.CreateLoan(&req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	// Retornar respuesta exitosa
	utils.CreatedResponse(c, "Solicitud de préstamo creada exitosamente", loan.ToResponse())
}

// GetLoan godoc
// @Summary Obtener préstamo por ID
// @Description Obtiene la información de un préstamo por su ID
// @Tags loans
// @Accept json
// @Produce json
// @Param id path int true "ID del préstamo"
// @Success 200 {object} utils.APIResponse{data=models.LoanResponse}
// @Failure 400 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /api/v1/loans/{id} [get]
func (ctrl *LoanController) GetLoan(c *gin.Context) {
	// Obtener ID del parámetro URL
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID de préstamo inválido")
		return
	}

	// Obtener préstamo con información del usuario
	loan, err := ctrl.loanService.GetLoanByIDWithUser(uint(id))
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	// Retornar respuesta exitosa
	utils.SuccessResponse(c, http.StatusOK, "Préstamo obtenido exitosamente", loan.ToResponse())
}

// UpdateLoanStatus godoc
// @Summary Actualizar estado del préstamo
// @Description Actualiza el estado de un préstamo (pending, approved, rejected)
// @Tags loans
// @Accept json
// @Produce json
// @Param id path int true "ID del préstamo"
// @Param status body models.LoanStatusUpdateRequest true "Nuevo estado del préstamo"
// @Success 200 {object} utils.APIResponse{data=models.LoanResponse}
// @Failure 400 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /api/v1/loans/{id}/status [put]
func (ctrl *LoanController) UpdateLoanStatus(c *gin.Context) {
	// Obtener ID del parámetro URL
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID de préstamo inválido")
		return
	}

	var req models.LoanStatusUpdateRequest

	// Parsear JSON del request
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Formato JSON inválido")
		return
	}

	// Actualizar estado del préstamo
	loan, err := ctrl.loanService.UpdateLoanStatus(uint(id), &req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	// Retornar respuesta exitosa
	utils.UpdatedResponse(c, "Estado del préstamo actualizado exitosamente", loan.ToResponse())
}

// ListLoans godoc
// @Summary Listar préstamos
// @Description Obtiene una lista paginada de préstamos con información del usuario
// @Tags loans
// @Accept json
// @Produce json
// @Param page query int false "Número de página" default(1)
// @Param limit query int false "Elementos por página" default(20)
// @Success 200 {object} utils.PaginatedResponse{data=[]models.LoanResponse}
// @Failure 400 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /api/v1/loans [get]
func (ctrl *LoanController) ListLoans(c *gin.Context) {
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

	// Obtener préstamos con información del usuario
	loans, total, err := ctrl.loanService.ListLoansWithUsers(page, limit)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	// Convertir a respuestas
	loanResponses := make([]models.LoanResponse, len(loans))
	for i, loan := range loans {
		loanResponses[i] = loan.ToResponse()
	}

	// Crear paginación
	pagination := utils.NewPagination(page, limit, total)

	// Retornar respuesta paginada
	utils.PaginatedSuccessResponse(c, "Préstamos obtenidos exitosamente", loanResponses, pagination)
}

// GetLoansByUser godoc
// @Summary Obtener préstamos de un usuario
// @Description Obtiene una lista paginada de préstamos de un usuario específico
// @Tags loans
// @Accept json
// @Produce json
// @Param userId path int true "ID del usuario"
// @Param page query int false "Número de página" default(1)
// @Param limit query int false "Elementos por página" default(20)
// @Success 200 {object} utils.PaginatedResponse{data=[]models.LoanResponse}
// @Failure 400 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /api/v1/loans/user/{userId} [get]
func (ctrl *LoanController) GetLoansByUser(c *gin.Context) {
	// Obtener ID del usuario del parámetro URL
	userIDStr := c.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID de usuario inválido")
		return
	}

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

	// Obtener préstamos del usuario
	loans, total, err := ctrl.loanService.GetLoansByUserID(uint(userID), page, limit)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	// Convertir a respuestas
	loanResponses := make([]models.LoanResponse, len(loans))
	for i, loan := range loans {
		loanResponses[i] = loan.ToResponse()
	}

	// Crear paginación
	pagination := utils.NewPagination(page, limit, total)

	// Retornar respuesta paginada
	utils.PaginatedSuccessResponse(c, "Préstamos del usuario obtenidos exitosamente", loanResponses, pagination)
}

// GetLoansByStatus godoc
// @Summary Obtener préstamos por estado
// @Description Obtiene una lista paginada de préstamos filtrados por estado
// @Tags loans
// @Accept json
// @Produce json
// @Param status query string true "Estado del préstamo (pending, approved, rejected)"
// @Param page query int false "Número de página" default(1)
// @Param limit query int false "Elementos por página" default(20)
// @Success 200 {object} utils.PaginatedResponse{data=[]models.LoanResponse}
// @Failure 400 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /api/v1/loans/status [get]
func (ctrl *LoanController) GetLoansByStatus(c *gin.Context) {
	// Obtener estado del query parameter
	status := c.Query("status")
	if status == "" {
		utils.BadRequestResponse(c, "Estado es requerido")
		return
	}

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

	// Obtener préstamos por estado
	loans, total, err := ctrl.loanService.GetLoansByStatus(status, page, limit)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	// Convertir a respuestas
	loanResponses := make([]models.LoanResponse, len(loans))
	for i, loan := range loans {
		loanResponses[i] = loan.ToResponse()
	}

	// Crear paginación
	pagination := utils.NewPagination(page, limit, total)

	// Retornar respuesta paginada
	utils.PaginatedSuccessResponse(c, "Préstamos obtenidos exitosamente", loanResponses, pagination)
}

// ProcessLoanApproval godoc
// @Summary Procesar aprobación automática de préstamo
// @Description Procesa la aprobación automática de un préstamo basado en criterios predefinidos
// @Tags loans
// @Accept json
// @Produce json
// @Param id path int true "ID del préstamo"
// @Success 200 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /api/v1/loans/{id}/process [post]
func (ctrl *LoanController) ProcessLoanApproval(c *gin.Context) {
	// Obtener ID del parámetro URL
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID de préstamo inválido")
		return
	}

	// Procesar aprobación automática
	err = ctrl.loanService.ProcessLoanApproval(uint(id))
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	// Retornar respuesta exitosa
	utils.SuccessResponse(c, http.StatusOK, "Préstamo procesado exitosamente", nil)
}

// GetLoanStatistics godoc
// @Summary Obtener estadísticas de préstamos
// @Description Obtiene estadísticas generales de préstamos del sistema
// @Tags loans
// @Accept json
// @Produce json
// @Success 200 {object} utils.APIResponse{data=map[string]interface{}}
// @Failure 500 {object} utils.APIResponse
// @Router /api/v1/loans/statistics [get]
func (ctrl *LoanController) GetLoanStatistics(c *gin.Context) {
	// Obtener estadísticas
	stats, err := ctrl.loanService.GetLoanStatistics()
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	// Retornar respuesta exitosa
	utils.SuccessResponse(c, http.StatusOK, "Estadísticas obtenidas exitosamente", stats)
}

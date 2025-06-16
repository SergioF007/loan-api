package controllers

import (
	"loan-api/models"
	"loan-api/services"
	"loan-api/utils"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

// LoanController maneja las operaciones relacionadas con préstamos
type LoanController struct {
	loanService   services.LoanService
	tenantService services.TenantService
}

// NewLoanController crea una nueva instancia del controlador de préstamos
func NewLoanController(loanService services.LoanService, tenantService services.TenantService) *LoanController {
	return &LoanController{
		loanService:   loanService,
		tenantService: tenantService,
	}
}

// CreateLoan godoc
// @Summary Crear una nueva solicitud de préstamo
// @Description Crea una nueva solicitud de préstamo para un usuario autenticado
// @Tags loans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "ID del tenant"
// @Param loan body models.CreateLoanRequest true "Datos de la solicitud de préstamo"
// @Success 201 {object} utils.APIResponse{data=models.LoanResponse}
// @Failure 400 {object} utils.APIResponse
// @Failure 401 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /loans [post]
func (ctrl *LoanController) CreateLoan(c *gin.Context) {
	log.Println("LoanController::CreateLoan was invoked")

	var req models.CreateLoanRequest

	// Parsear JSON del request
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Formato JSON inválido")
		return
	}

	// Validar campos requeridos
	if req.LoanTypeID == 0 {
		utils.BadRequestResponse(c, "loan_type_id es requerido")
		return
	}

	// Obtener ID del usuario desde el contexto (middleware de autenticación)
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Token de autenticación requerido")
		return
	}

	// Crear préstamo
	loanResponse, err := ctrl.loanService.CreateLoan(userID.(uint), req)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	// Retornar respuesta exitosa
	utils.CreatedResponse(c, "Solicitud de préstamo creada exitosamente", loanResponse)
}

// SaveLoanData godoc
// @Summary Guardar datos de una solicitud de préstamo
// @Description Guarda los datos dinámicos de una solicitud de préstamo existente
// @Tags loans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "ID del tenant"
// @Param data body models.SaveLoanDataRequest true "Datos del préstamo a guardar"
// @Success 200 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Failure 401 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /loans/data [post]
func (ctrl *LoanController) SaveLoanData(c *gin.Context) {
	log.Println("LoanController::SaveLoanData was invoked")

	var req models.SaveLoanDataRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Formato JSON inválido")
		return
	}

	if req.LoanID == 0 {
		utils.BadRequestResponse(c, "loan_id es requerido")
		return
	}

	if len(req.Data) == 0 {
		utils.BadRequestResponse(c, "data es requerido")
		return
	}

	// Guardar datos del préstamo
	if err := ctrl.loanService.SaveLoanData(req); err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	// Retornar respuesta exitosa
	utils.SuccessResponse(c, 200, "Datos del préstamo guardados exitosamente", nil)
}

// GetLoan godoc
// @Summary Obtener información de un préstamo
// @Description Obtiene la información completa de un préstamo por ID
// @Tags loans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "ID del tenant"
// @Param id path int true "ID del préstamo"
// @Success 200 {object} utils.APIResponse{data=models.LoanResponse}
// @Failure 400 {object} utils.APIResponse
// @Failure 401 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /loans/{id} [get]
func (ctrl *LoanController) GetLoan(c *gin.Context) {
	log.Println("LoanController::GetLoan was invoked")

	// Obtener ID del préstamo desde los parámetros
	loanIDStr := c.Param("id")
	loanID, err := strconv.ParseUint(loanIDStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID del préstamo debe ser un número válido")
		return
	}

	// Obtener préstamo
	loanResponse, err := ctrl.loanService.GetLoanByID(uint(loanID))
	if err != nil {
		utils.NotFoundResponse(c, err.Error())
		return
	}

	// Retornar respuesta exitosa
	utils.SuccessResponse(c, 200, "Préstamo obtenido exitosamente", loanResponse)
}

// GetUserLoans godoc
// @Summary Obtener préstamos de un usuario
// @Description Obtiene todos los préstamos de un usuario autenticado
// @Tags loans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "ID del tenant"
// @Success 200 {object} utils.APIResponse{data=[]models.LoanResponse}
// @Failure 400 {object} utils.APIResponse
// @Failure 401 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /loans/user [get]
func (ctrl *LoanController) GetUserLoans(c *gin.Context) {
	log.Println("LoanController::GetUserLoans was invoked")

	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Token de autenticación requerido")
		return
	}

	loansResponse, err := ctrl.loanService.GetLoansByUserID(userID.(uint))
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	// Retornar respuesta exitosa
	utils.SuccessResponse(c, 200, "Préstamos obtenidos exitosamente", loansResponse)
}

// ProcessLoanDecision godoc
// @Summary Procesar decisión final del préstamo
// @Description Evalúa el score crediticio y verificación de identidad para aprobar/rechazar y realizar desembolso
// @Tags loans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "ID del tenant"
// @Param id path int true "ID del préstamo"
// @Success 200 {object} utils.APIResponse{data=models.LoanResponse}
// @Failure 400 {object} utils.APIResponse
// @Failure 401 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /loans/{id}/decision [post]
func (ctrl *LoanController) ProcessLoanDecision(c *gin.Context) {
	log.Println("LoanController::ProcessLoanDecision was invoked")

	loanIDStr := c.Param("id")
	loanID, err := strconv.ParseUint(loanIDStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID del préstamo debe ser un número válido")
		return
	}

	loanResponse, err := ctrl.loanService.ProcessLoanDecision(uint(loanID))
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, 200, "Decisión del préstamo procesada exitosamente", loanResponse)
}

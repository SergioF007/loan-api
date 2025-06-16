package controllers

import (
	"loan-api/services"
	"loan-api/utils"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

// LoanTypeController maneja las operaciones relacionadas con tipos de préstamo
type LoanTypeController struct {
	loanTypeService services.LoanTypeService
	tenantService   services.TenantService
}

// NewLoanTypeController crea una nueva instancia del controlador de tipos de préstamo
func NewLoanTypeController(loanTypeService services.LoanTypeService, tenantService services.TenantService) *LoanTypeController {
	return &LoanTypeController{
		loanTypeService: loanTypeService,
		tenantService:   tenantService,
	}
}

// GetLoanTypesWithForms godoc
// @Summary Obtener tipos de préstamo con formularios
// @Description Obtiene todos los tipos de préstamo con sus formularios e inputs disponibles para un tenant
// @Tags loan-types
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "ID del tenant"
// @Success 200 {object} utils.APIResponse{data=[]models.LoanTypeResponse}
// @Failure 400 {object} utils.APIResponse
// @Failure 401 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /loan-types [get]
func (ctrl *LoanTypeController) GetLoanTypesWithForms(c *gin.Context) {
	log.Println("LoanTypeController::GetLoanTypesWithForms was invoked")

	// Validar header X-Tenant-ID
	tenantIDStr := c.GetHeader("X-Tenant-ID")
	if tenantIDStr == "" {
		utils.BadRequestResponse(c, "Header X-Tenant-ID es requerido")
		return
	}

	tenantID, err := strconv.ParseUint(tenantIDStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "X-Tenant-ID debe ser un número válido")
		return
	}

	// Validar que el tenant existe
	_, err = ctrl.tenantService.ValidateTenantID(uint(tenantID))
	if err != nil {
		utils.NotFoundResponse(c, "Tenant no encontrado")
		return
	}

	// Obtener tipos de préstamo con formularios
	loanTypes, err := ctrl.loanTypeService.GetLoanTypesWithForms(uint(tenantID))
	if err != nil {
		utils.InternalServerErrorResponse(c, "Error al obtener tipos de préstamo: "+err.Error())
		return
	}

	// Retornar respuesta exitosa
	utils.SuccessResponse(c, 200, "Tipos de préstamo obtenidos exitosamente", loanTypes)
}

// GetLoanTypeByCode godoc
// @Summary Obtener tipo de préstamo por código
// @Description Obtiene un tipo de préstamo específico con sus formularios e inputs por código
// @Tags loan-types
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "ID del tenant"
// @Param code path string true "Código del tipo de préstamo"
// @Success 200 {object} utils.APIResponse{data=models.LoanTypeResponse}
// @Failure 400 {object} utils.APIResponse
// @Failure 401 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /loan-types/{code} [get]
func (ctrl *LoanTypeController) GetLoanTypeByCode(c *gin.Context) {
	log.Println("LoanTypeController::GetLoanTypeByCode was invoked")

	// Validar header X-Tenant-ID
	tenantIDStr := c.GetHeader("X-Tenant-ID")
	if tenantIDStr == "" {
		utils.BadRequestResponse(c, "Header X-Tenant-ID es requerido")
		return
	}

	tenantID, err := strconv.ParseUint(tenantIDStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "X-Tenant-ID debe ser un número válido")
		return
	}

	// Validar que el tenant existe
	_, err = ctrl.tenantService.ValidateTenantID(uint(tenantID))
	if err != nil {
		utils.NotFoundResponse(c, "Tenant no encontrado")
		return
	}

	// Obtener código del tipo de préstamo desde los parámetros
	code := c.Param("code")
	if code == "" {
		utils.BadRequestResponse(c, "Código del tipo de préstamo es requerido")
		return
	}

	// Obtener tipo de préstamo por código
	loanType, err := ctrl.loanTypeService.GetLoanTypeByCode(uint(tenantID), code)
	if err != nil {
		utils.NotFoundResponse(c, "Tipo de préstamo no encontrado: "+err.Error())
		return
	}

	// Retornar respuesta exitosa
	utils.SuccessResponse(c, 200, "Tipo de préstamo obtenido exitosamente", loanType)
}

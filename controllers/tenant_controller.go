package controllers

import (
	"loan-api/services"
	"loan-api/utils"
	"log"

	"github.com/gin-gonic/gin"
)

// TenantController maneja las operaciones relacionadas con tenants
type TenantController struct {
	tenantService services.TenantService
}

// NewTenantController crea una nueva instancia del controlador de tenants
func NewTenantController(tenantService services.TenantService) *TenantController {
	return &TenantController{
		tenantService: tenantService,
	}
}

// GetAvailableTenants godoc
// @Summary Obtener tenants disponibles
// @Description Obtiene todos los tenants disponibles para pruebas
// @Tags tenants
// @Accept json
// @Produce json
// @Success 200 {object} utils.APIResponse{data=[]models.TenantResponse}
// @Failure 500 {object} utils.APIResponse
// @Router /tenants [get]
func (ctrl *TenantController) GetAvailableTenants(c *gin.Context) {
	log.Println("TenantController::GetAvailableTenants was invoked")

	// Obtener tenants disponibles
	tenants, err := ctrl.tenantService.GetAvailableTenants()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Error al obtener tenants: "+err.Error())
		return
	}

	// Retornar respuesta exitosa
	utils.SuccessResponse(c, 200, "Tenants obtenidos exitosamente", tenants)
}

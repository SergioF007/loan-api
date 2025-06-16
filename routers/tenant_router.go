package routers

import (
	"loan-api/controllers"

	"github.com/gin-gonic/gin"
)

// TenantRouter configura las rutas relacionadas con tenants
type TenantRouter struct {
	tenantController *controllers.TenantController
}

// NewTenantRouter crea una nueva instancia del router de tenants
func NewTenantRouter(tenantController *controllers.TenantController) *TenantRouter {
	return &TenantRouter{
		tenantController: tenantController,
	}
}

// Setup configura todas las rutas de tenants
func (r *TenantRouter) Setup(router *gin.RouterGroup) {
	// Rutas públicas de tenants
	tenants := router.Group("/tenants")
	{
		tenants.GET("", r.tenantController.GetAvailableTenants) // GET /api/v1/tenants - Obtener tenants disponibles
	}
}

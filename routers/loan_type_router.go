package routers

import (
	"loan-api/controllers"
	"loan-api/middlewares"

	"github.com/gin-gonic/gin"
)

// LoanTypeRouter configura las rutas relacionadas con tipos de préstamo
type LoanTypeRouter struct {
	loanTypeController *controllers.LoanTypeController
}

// NewLoanTypeRouter crea una nueva instancia del router de tipos de préstamo
func NewLoanTypeRouter(loanTypeController *controllers.LoanTypeController) *LoanTypeRouter {
	return &LoanTypeRouter{
		loanTypeController: loanTypeController,
	}
}

// Setup configura todas las rutas de tipos de préstamo
func (r *LoanTypeRouter) Setup(router *gin.RouterGroup) {
	// Grupo de rutas para tipos de préstamo
	loanTypes := router.Group("/loan-types")
	{
		// Todas las rutas de tipos de préstamo requieren autenticación
		loanTypes.Use(middlewares.AuthMiddleware())

		loanTypes.GET("", r.loanTypeController.GetLoanTypesWithForms)   // GET /api/v1/loan-types - Obtener tipos con formularios
		loanTypes.GET("/:code", r.loanTypeController.GetLoanTypeByCode) // GET /api/v1/loan-types/{code} - Obtener tipo por código
	}
}

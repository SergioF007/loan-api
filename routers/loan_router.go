package routers

import (
	"loan-api/controllers"
	"loan-api/middlewares"

	"github.com/gin-gonic/gin"
)

// LoanRouter configura las rutas relacionadas con préstamos
type LoanRouter struct {
	loanController *controllers.LoanController
}

// NewLoanRouter crea una nueva instancia del router de préstamos
func NewLoanRouter(loanController *controllers.LoanController) *LoanRouter {
	return &LoanRouter{
		loanController: loanController,
	}
}

// Setup configura todas las rutas de préstamos
func (r *LoanRouter) Setup(router *gin.RouterGroup) {
	// Grupo de rutas para préstamos
	loans := router.Group("/loans")
	{
		// Todas las rutas de préstamos requieren autenticación
		loans.Use(middlewares.AuthMiddleware())

		loans.POST("", r.loanController.CreateLoan)        // POST /api/v1/loans - Crear préstamo
		loans.POST("/data", r.loanController.SaveLoanData) // POST /api/v1/loans/data - Guardar datos del préstamo
		loans.GET("/:id", r.loanController.GetLoan)        // GET /api/v1/loans/{id} - Obtener préstamo por ID
		loans.GET("/user", r.loanController.GetUserLoans)  // GET /api/v1/loans/user - Obtener préstamos del usuario
	}
}

package routers

import (
	"loan-api/controllers"
	"loan-api/middlewares"

	"github.com/gin-gonic/gin"
)

// LoanRouter configura las rutas relacionadas con préstamos
type LoanRouter struct {
	loanController *controllers.LoanController
	authMiddleware *middlewares.AuthMiddleware
}

// NewLoanRouter crea una nueva instancia del router de préstamos
func NewLoanRouter(loanController *controllers.LoanController, authMiddleware *middlewares.AuthMiddleware) *LoanRouter {
	return &LoanRouter{
		loanController: loanController,
		authMiddleware: authMiddleware,
	}
}

// SetupRoutes configura todas las rutas de préstamos
func (r *LoanRouter) SetupRoutes(router *gin.RouterGroup) {
	// Grupo de rutas para préstamos
	loans := router.Group("/loans")
	{
		// Rutas públicas (solo lectura)
		loans.GET("", r.authMiddleware.OptionalAuth(), r.loanController.ListLoans)                    // GET /api/v1/loans - Listar préstamos
		loans.GET("/status", r.authMiddleware.OptionalAuth(), r.loanController.GetLoansByStatus)      // GET /api/v1/loans/status?status=... - Préstamos por estado
		loans.GET("/statistics", r.authMiddleware.OptionalAuth(), r.loanController.GetLoanStatistics) // GET /api/v1/loans/statistics - Estadísticas
		loans.GET("/:id", r.authMiddleware.OptionalAuth(), r.loanController.GetLoan)                  // GET /api/v1/loans/:id - Obtener préstamo específico

		// Rutas que requieren autenticación
		loans.POST("", r.authMiddleware.RequireAuth(), r.loanController.CreateLoan)                      // POST /api/v1/loans - Crear solicitud de préstamo
		loans.PUT("/:id/status", r.authMiddleware.RequireAuth(), r.loanController.UpdateLoanStatus)      // PUT /api/v1/loans/:id/status - Actualizar estado del préstamo
		loans.POST("/:id/process", r.authMiddleware.RequireAuth(), r.loanController.ProcessLoanApproval) // POST /api/v1/loans/:id/process - Procesar aprobación automática
		loans.GET("/user/:userId", r.authMiddleware.RequireAuth(), r.loanController.GetLoansByUser)      // GET /api/v1/loans/user/:userId - Préstamos de un usuario
	}
}

// SetupPublicRoutes configura rutas públicas de préstamos (sin autenticación)
func (r *LoanRouter) SetupPublicRoutes(router *gin.RouterGroup) {
	// Grupo de rutas públicas para préstamos (solo consultas)
	publicLoans := router.Group("/loans")
	{
		publicLoans.GET("", r.loanController.ListLoans)                    // Listar préstamos
		publicLoans.GET("/status", r.loanController.GetLoansByStatus)      // Préstamos por estado
		publicLoans.GET("/statistics", r.loanController.GetLoanStatistics) // Estadísticas públicas
		publicLoans.GET("/:id", r.loanController.GetLoan)                  // Obtener préstamo específico
	}
}

// SetupProtectedRoutes configura rutas protegidas de préstamos (requieren autenticación)
func (r *LoanRouter) SetupProtectedRoutes(router *gin.RouterGroup) {
	// Aplicar middleware de autenticación al grupo
	protected := router.Group("/loans", r.authMiddleware.RequireAuth())
	{
		// Operaciones CRUD de préstamos
		protected.POST("", r.loanController.CreateLoan)                      // Crear solicitud de préstamo
		protected.PUT("/:id/status", r.loanController.UpdateLoanStatus)      // Actualizar estado del préstamo
		protected.POST("/:id/process", r.loanController.ProcessLoanApproval) // Procesar aprobación automática

		// Consultas específicas que requieren autenticación
		protected.GET("/user/:userId", r.loanController.GetLoansByUser) // Préstamos de un usuario específico

		// Operaciones administrativas
		protected.GET("", r.loanController.ListLoans)                    // Listar préstamos (con más detalles)
		protected.GET("/:id", r.loanController.GetLoan)                  // Obtener préstamo específico
		protected.GET("/status", r.loanController.GetLoansByStatus)      // Préstamos por estado
		protected.GET("/statistics", r.loanController.GetLoanStatistics) // Estadísticas completas
	}
}

// SetupAdminRoutes configura rutas administrativas de préstamos (requieren roles específicos)
// Nota: Esta funcionalidad se puede expandir cuando se implemente un sistema de roles
func (r *LoanRouter) SetupAdminRoutes(router *gin.RouterGroup) {
	// Aplicar middleware de autenticación y autorización administrativa
	admin := router.Group("/loans", r.authMiddleware.RequireAuth()) // TODO: Agregar middleware de rol admin
	{
		// Operaciones administrativas específicas
		admin.PUT("/:id/status", r.loanController.UpdateLoanStatus)      // Cambiar estado manualmente
		admin.POST("/:id/process", r.loanController.ProcessLoanApproval) // Procesar aprobaciones
		admin.GET("/statistics", r.loanController.GetLoanStatistics)     // Estadísticas completas del sistema
	}
}

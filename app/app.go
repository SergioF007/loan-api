package app

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	"loan-api/config"
	"loan-api/controllers"
	"loan-api/database"
	"loan-api/middlewares"
	"loan-api/repositories"
	"loan-api/routers"
	"loan-api/services"
)

// Validadores personalizados
var alphaValidation validator.Func = func(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(string)
	if ok {
		match, err := regexp.MatchString("^[a-zA-Z0-9 ]*$", value)
		if err != nil {
			return false
		}
		return match
	}
	return false
}

var passwordValidation validator.Func = func(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(string)
	if ok {
		// Verificar longitud mínima
		if len(value) < 8 {
			return false
		}

		var (
			hasUpper   bool
			hasLower   bool
			hasNumber  bool
			hasSpecial bool
		)

		for _, char := range value {
			switch {
			case 'A' <= char && char <= 'Z':
				hasUpper = true
			case 'a' <= char && char <= 'z':
				hasLower = true
			case '0' <= char && char <= '9':
				hasNumber = true
			case strings.ContainsRune("!@#$%^&*()-_+=[]{}|;:'\",.<>/?~", char):
				hasSpecial = true
			}
		}

		return hasUpper && hasLower && hasNumber && hasSpecial
	}
	return false
}

var creditScoreValidation validator.Func = func(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(int)
	if ok {
		return value >= 300 && value <= 850
	}
	return false
}

var loanAmountValidation validator.Func = func(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(float64)
	if ok {
		return value >= 100000 && value <= 50000000 // Entre 100K y 50M COP
	}
	return false
}

// SetupRouter configura el router de Gin específicamente para testing
func SetupRouter(cfg config.Config) *gin.Engine {
	// Configurar Gin en modo test
	gin.SetMode(gin.TestMode)

	// Registrar validaciones personalizadas
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("alpha", alphaValidation)
		v.RegisterValidation("password", passwordValidation)
		v.RegisterValidation("credit_score", creditScoreValidation)
		v.RegisterValidation("loan_amount", loanAmountValidation)
	}

	// Conectar a la base de datos (usa la configuración pasada)
	database.Connect(&cfg)

	// Ejecutar migraciones para tests
	database.Migrate()

	// Inicializar repositorios
	userRepository := repositories.NewUserRepository(database.DB)
	loanRepository := repositories.NewLoanRepository(database.DB)
	tenantRepository := repositories.NewTenantRepository(database.DB)
	loanTypeRepository := repositories.NewLoanTypeRepository(database.DB)

	// Inicializar servicios
	userService := services.NewUserService(userRepository)
	tenantService := services.NewTenantService(tenantRepository)
	loanTypeService := services.NewLoanTypeService(loanTypeRepository)
	loanService := services.NewLoanService(loanRepository, userRepository, loanTypeRepository)

	// Inicializar controladores
	userController := controllers.NewUserController(userService, &cfg)
	tenantController := controllers.NewTenantController(tenantService)
	loanTypeController := controllers.NewLoanTypeController(loanTypeService, tenantService)
	loanController := controllers.NewLoanController(loanService, tenantService)

	// Configurar servidor Gin
	router := gin.New()
	router.MaxMultipartMemory = 8 << 20 // 8 MiB

	// Configurar logger personalizado para testing
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// En modo test, logging más simple
		if cfg.AppEnv == "test" {
			return fmt.Sprintf("[TEST] %s %s %d\n",
				param.Method,
				param.Path,
				param.StatusCode,
			)
		}
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	// Configurar CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Tenant-ID"}
	corsConfig.AllowCredentials = true

	router.Use(cors.New(corsConfig))

	// Middleware de recuperación de errores
	router.Use(gin.Recovery())

	// Aplicar middleware Tenant de manera global
	router.Use(middlewares.Tenant())

	// Configurar rutas
	apiGroup := router.Group("/loan-api/api/v1")

	// Inicializar y configurar routers
	userRouter := routers.NewUserRouter(userController)
	loanRouter := routers.NewLoanRouter(loanController)
	tenantRouter := routers.NewTenantRouter(tenantController)
	loanTypeRouter := routers.NewLoanTypeRouter(loanTypeController)

	// Configurar rutas de los módulos
	userRouter.Setup(apiGroup)
	loanRouter.Setup(apiGroup)
	tenantRouter.Setup(apiGroup)
	loanTypeRouter.Setup(apiGroup)

	// Ruta de health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "success",
			"message": "El servicio de préstamos está funcionando correctamente",
		})
	})

	// Ruta de health check en el grupo API
	apiGroup.GET("/health-checker", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "success",
			"message": "El servicio de préstamos está funcionando correctamente",
		})
	})

	return router
}

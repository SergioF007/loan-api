package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"

	"loan-api/config"
	"loan-api/controllers"
	"loan-api/database"
	_ "loan-api/docs"
	"loan-api/middlewares"
	"loan-api/repositories"
	"loan-api/routers"
	"loan-api/services"
)

type Server struct {
	DB     *gorm.DB
	Config config.Config
}

var (
	server *gin.Engine
)

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

// Validación para score de crédito
var creditScoreValidation validator.Func = func(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(int)
	if ok {
		return value >= 300 && value <= 850
	}
	return false
}

// Validación para montos de préstamo
var loanAmountValidation validator.Func = func(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(float64)
	if ok {
		return value >= 100000 && value <= 50000000 // Entre 100K y 50M COP
	}
	return false
}

// @title Loan API
// @version 1.0
// @description API REST para gestión de solicitudes de préstamos
// @termsOfService http://swagger.io/terms/

// @contact.name Soporte API
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /loan-api/api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	// Registrar validaciones personalizadas
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("alpha", alphaValidation)
		v.RegisterValidation("password", passwordValidation)
		v.RegisterValidation("credit_score", creditScoreValidation)
		v.RegisterValidation("loan_amount", loanAmountValidation)
	}

	// Cargar configuración
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("No se pudieron cargar las variables de entorno", err)
	}

	// Conectar a la base de datos
	database.Connect(&config)

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
	userController := controllers.NewUserController(userService, &config)
	tenantController := controllers.NewTenantController(tenantService)
	loanTypeController := controllers.NewLoanTypeController(loanTypeService, tenantService)
	loanController := controllers.NewLoanController(loanService, tenantService)

	// Configurar servidor Gin
	server = gin.New()
	server.MaxMultipartMemory = 8 << 20 // 8 MiB

	// Configurar logger personalizado
	server.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		if param.Path == "/loan-api/api/v1/health-checker" {
			return ""
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

	server.Use(cors.New(corsConfig))

	// Aplicar middleware Tenant de manera global
	router := server.Group("/loan-api/api/v1")
	router.Use(middlewares.Tenant())

	// Configurar Swagger
	if config.AppEnv == "local" || config.AppEnv == "dev" {
		// Configurar Swagger en una ruta separada para evitar conflictos
		server.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// Configurar rutas
	userRouter := routers.NewUserRouter(userController)
	tenantRouter := routers.NewTenantRouter(tenantController)
	loanTypeRouter := routers.NewLoanTypeRouter(loanTypeController)
	loanRouter := routers.NewLoanRouter(loanController)

	// Configurar rutas de los módulos
	userRouter.Setup(router)
	tenantRouter.Setup(router)
	loanTypeRouter.Setup(router)
	loanRouter.Setup(router)

	// Ruta de health check
	router.GET("/health-checker", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "success",
			"message": "El servicio de préstamos está funcionando correctamente",
		})
	})

	// Ruta raíz con información de la API
	server.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":     "¡Bienvenido a Loan API!",
			"version":     config.AppVersion,
			"environment": config.AppEnv,
			"docs":        "/docs/index.html",
			"health":      "/loan-api/api/v1/health-checker",
		})
	})

	// Iniciar servidor
	log.Printf("🚀 Iniciando servidor en http://localhost:%s", config.ServerPort)
	log.Printf("📚 Documentación disponible en: http://localhost:%s/docs/index.html", config.ServerPort)
	log.Printf("💚 Health check: http://localhost:%s/loan-api/api/v1/health-checker", config.ServerPort)

	log.Fatal(server.Run(":" + config.ServerPort))
}

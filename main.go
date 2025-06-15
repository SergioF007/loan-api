package main

import (
	"log"
	"net/http"

	"loan-api/config"
	"loan-api/controllers"
	"loan-api/database"
	_ "loan-api/docs" // Importar docs generados por swag
	"loan-api/middlewares"
	"loan-api/repositories"
	"loan-api/routers"
	"loan-api/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

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
// @BasePath /api/v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// 1. Cargar configuración
	log.Println("🔧 Cargando configuración...")
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("❌ Error al cargar la configuración: %v", err)
	}
	log.Printf("✅ Configuración cargada - Entorno: %s", cfg.AppEnv)

	// 2. Conectar base de datos
	log.Println("🔌 Conectando a la base de datos...")
	if err := database.ConnectDB(&cfg); err != nil {
		log.Fatalf("❌ Error al conectar con la base de datos: %v", err)
	}

	// 3. Ejecutar migraciones
	log.Println("🔄 Ejecutando migraciones...")
	if err := database.AutoMigrate(); err != nil {
		log.Fatalf("❌ Error en las migraciones: %v", err)
	}

	// 4. Inicializar repositorios
	log.Println("📦 Inicializando repositorios...")
	db := database.GetDB()
	userRepo := repositories.NewUserRepository(db)
	loanRepo := repositories.NewLoanRepository(db)

	// 5. Inicializar servicios
	log.Println("⚙️  Inicializando servicios...")
	userService := services.NewUserService(userRepo)
	loanService := services.NewLoanService(loanRepo, userRepo)

	// 6. Inicializar controladores
	log.Println("🎮 Inicializando controladores...")
	userController := controllers.NewUserController(userService)
	loanController := controllers.NewLoanController(loanService)

	// 7. Inicializar middlewares
	log.Println("🛡️  Inicializando middlewares...")
	authMiddleware := middlewares.NewAuthMiddleware(&cfg)

	// 8. Configurar Gin
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 9. Crear router principal
	log.Println("🚦 Configurando rutas...")
	app := gin.New()

	// 10. Configurar middlewares globales
	app.Use(middlewares.LoggerMiddleware())
	app.Use(gin.Recovery())
	app.Use(middlewares.ErrorHandlerMiddleware())

	// Configurar CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	app.Use(cors.New(corsConfig))

	// 11. Configurar rutas de salud
	app.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "OK",
			"message": "Loan API funcionando correctamente",
			"version": cfg.AppVersion,
		})
	})

	app.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "¡Bienvenido a Loan API!",
			"version": cfg.AppVersion,
			"docs":    "/swagger/index.html",
		})
	})

	// 12. Configurar routers de la API
	apiV1 := app.Group("/api/v1")

	// Inicializar routers
	userRouter := routers.NewUserRouter(userController, authMiddleware)
	loanRouter := routers.NewLoanRouter(loanController, authMiddleware)

	// Configurar rutas
	userRouter.SetupRoutes(apiV1)
	loanRouter.SetupRoutes(apiV1)

	// 13. Configurar Swagger
	if cfg.IsDevelopment() {
		app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		log.Printf("📚 Documentación Swagger disponible en: http://%s/swagger/index.html", cfg.GetServerAddress())
	}

	// 14. Iniciar servidor
	serverAddr := cfg.GetServerAddress()
	log.Printf("🚀 Iniciando servidor en http://%s", serverAddr)
	log.Printf("📚 Documentación disponible en: http://%s/swagger/index.html", serverAddr)
	log.Printf("💚 Health check: http://%s/health", serverAddr)

	if err := app.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("❌ Error al iniciar el servidor: %v", err)
	}
}

// gracefulShutdown maneja el cierre graceful del servidor
func gracefulShutdown() {
	log.Println("🔄 Cerrando conexiones...")

	// Cerrar conexión a la base de datos
	if err := database.CloseDB(); err != nil {
		log.Printf("⚠️  Error al cerrar la base de datos: %v", err)
	}

	log.Println("✅ Servidor cerrado correctamente")
}

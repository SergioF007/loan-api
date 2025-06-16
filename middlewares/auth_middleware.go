package middlewares

import (
	"log"
	"net/http"

	"loan-api/config"
	"loan-api/database"
	"loan-api/models"
	"loan-api/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware middleware para autenticación
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			log.Println("request does not contain an access token")

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   true,
				"message": "Por favor, iniciar sesión",
			})
			return
		}

		loadConfig, err := config.LoadConfig(".")
		if err != nil {
			log.Fatal("could not load environment variables", err)

			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error":   true,
				"message": "Algo salió mal. Intentar otra vez.",
			})
			return
		}

		sub, err := utils.ValidateToken(tokenString, loadConfig.AccessTokenPublicKey)
		if err != nil {
			log.Println(err.Error())

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   true,
				"message": "Por favor, iniciar sesión",
			})
			return
		}

		// Extraer el ID del usuario del payload
		payload, ok := sub.(map[string]interface{})
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   true,
				"message": "Token inválido",
			})
			return
		}

		userID, ok := payload["id"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   true,
				"message": "Token inválido",
			})
			return
		}

		// Guardar solo el ID del usuario en el contexto
		c.Set("user_id", uint(userID))
		c.Next()
	}
}

// Tenant middleware para validar el tenant ID
func Tenant() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		loadConfig, err := config.LoadConfig(".")
		if err != nil {
			log.Fatal("could not load environment variables", err)

			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error":   true,
				"message": "Algo salió mal. Intentar otra vez.",
			})

			return
		}
		database.Connect(&loadConfig)

		tenantId := ctx.GetHeader("X-Tenant-ID")
		log.Printf("Tenant ID: %v", tenantId)
		var tenant models.Tenant

		result := database.DB.Where("id = ?", tenantId).First(&tenant)
		if result.Error != nil {
			log.Println(result.Error)
			log.Println("could not load tenant ", err)

			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":   true,
				"message": "Tenant invalido.",
			})

			return
		}

		ctx.Set("tenant_id", tenant.ID)
		ctx.Set("tenant", tenant)

		ctx.Next()
	}
}

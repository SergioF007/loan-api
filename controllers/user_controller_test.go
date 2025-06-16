package controllers_test

import (
	"encoding/json"
	"log"
	"os"
	"testing"

	"loan-api/config"
	"loan-api/database"
	"loan-api/models"
	"loan-api/test"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

var (
	DB     *gorm.DB
	CONFIG config.Config
)

func TestMain(m *testing.M) {
	cfg, err := config.LoadConfig("../")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	CONFIG = cfg
	database.Connect(&cfg)
	DB = database.DB

	code := m.Run()

	test.ClearDatabase(DB)

	os.Exit(code)
}

func TestUserController_RegisterUser(t *testing.T) {
	c := require.New(t)

	// Headers válidos
	headers := map[string]string{
		"X-Tenant-ID": "1",
	}

	t.Run("Debería registrar usuario exitosamente con datos válidos", func(t *testing.T) {

		// Datos de prueba
		email := "test@example.com"
		name := "Test User"
		phone := "3001234567"
		documentType := models.DocumentTypeCedula
		documentNumber := "12345678"
		password := "Password123!"

		// Verificar que el usuario no existe
		var userBefore models.User
		resultBefore := DB.Where("email = ?", email).First(&userBefore)
		c.Error(resultBefore.Error)
		c.Equal(gorm.ErrRecordNotFound, resultBefore.Error)

		// Construir el cuerpo de la solicitud
		requestBody := map[string]interface{}{
			"name":                  name,
			"email":                 email,
			"phone":                 phone,
			"document_type":         documentType,
			"document_number":       documentNumber,
			"password":              password,
			"password_confirmation": password,
		}

		// Realizar la petición POST
		w := test.MakePostRequest(CONFIG, "/loan-api/api/v1/auth/register", requestBody, headers)

		// Verificar código de respuesta exitoso
		c.Equal(201, w.Code)

		// Decodificar respuesta JSON
		var response map[string]interface{}
		c.NoError(json.Unmarshal(w.Body.Bytes(), &response))

		// Verificar mensaje de éxito
		c.Equal("Usuario registrado exitosamente", response["message"])

		// Verificar datos del usuario en la respuesta
		data := response["data"].(map[string]interface{})
		c.Equal(email, data["email"])
		c.Equal(name, data["name"])
		c.Equal(phone, data["phone"])
		c.Equal(string(documentType), data["document_type"])
		c.Equal(documentNumber, data["document_number"])

		// Verificar que el usuario se creó en la base de datos
		var userAfter models.User
		resultAfter := DB.Where("email = ?", email).First(&userAfter)
		c.NoError(resultAfter.Error)
		c.Equal(name, userAfter.Name)
		c.Equal(email, userAfter.Email)
		c.Equal(phone, userAfter.Phone)
		c.Equal(documentType, userAfter.DocumentType)
		c.Equal(documentNumber, userAfter.DocumentNumber)
		c.NotEmpty(userAfter.Password) // Verificar que la contraseña fue hasheada
		c.False(userAfter.CreatedAt.IsZero())
	})

	t.Run("Debería fallar con email duplicado", func(t *testing.T) {
		// Limpiar y cargar datos de prueba
		test.LoadTestData(DB)

		// Intentar crear usuario con email existente
		requestBody := map[string]interface{}{
			"name":                  "Otro Usuario",
			"email":                 "juan@example.com", // Email que ya existe en test data
			"phone":                 "3009999999",
			"document_type":         models.DocumentTypeCedula,
			"document_number":       "87654321",
			"password":              "Password123!",
			"password_confirmation": "Password123!",
		}

		// Realizar la petición POST
		w := test.MakePostRequest(CONFIG, "/loan-api/api/v1/auth/register", requestBody, headers)

		// Verificar código de respuesta de error
		c.Equal(409, w.Code)

		// Decodificar respuesta JSON
		var response map[string]interface{}
		c.NoError(json.Unmarshal(w.Body.Bytes(), &response))

		// Verificar que hay un campo error
		c.Contains(response, "error")

		// Acceder al mensaje dentro del campo error
		errorData := response["error"].(map[string]interface{})
		c.Contains(errorData["message"], "El email ya está registrado")
	})

	t.Run("Debería fallar con datos inválidos - email vacío", func(t *testing.T) {
		// Datos inválidos
		requestBody := map[string]interface{}{
			"name":                  "Test User",
			"email":                 "", // Email vacío
			"phone":                 "3001234567",
			"document_type":         models.DocumentTypeCedula,
			"document_number":       "12345678",
			"password":              "Password123!",
			"password_confirmation": "Password123!",
		}

		// Realizar la petición POST
		w := test.MakePostRequest(CONFIG, "/loan-api/api/v1/auth/register", requestBody, headers)

		// Verificar código de respuesta de error
		c.Equal(400, w.Code)

		// Decodificar respuesta JSON
		var response map[string]interface{}
		c.NoError(json.Unmarshal(w.Body.Bytes(), &response))

		// Verificar que hay un campo error (los errores de validación van en error)
		c.Contains(response, "error")

		// Acceder al mensaje dentro del campo error
		errorData := response["error"].(map[string]interface{})
		c.NotEmpty(errorData["message"])
	})

	t.Run("Debería fallar con contraseñas que no coinciden", func(t *testing.T) {
		// Datos inválidos
		requestBody := map[string]interface{}{
			"name":                  "Test User",
			"email":                 "test2@example.com",
			"phone":                 "3001234567",
			"document_type":         models.DocumentTypeCedula,
			"document_number":       "123456789",
			"password":              "Password123!",
			"password_confirmation": "DifferentPassword123!", // Contraseñas diferentes
		}

		// Realizar la petición POST
		w := test.MakePostRequest(CONFIG, "/loan-api/api/v1/auth/register", requestBody, headers)

		// Verificar código de respuesta de error
		c.Equal(400, w.Code)

		// Decodificar respuesta JSON
		var response map[string]interface{}
		c.NoError(json.Unmarshal(w.Body.Bytes(), &response))

		// Verificar que hay un campo error
		c.Contains(response, "error")

		// Acceder al mensaje dentro del campo error
		errorData := response["error"].(map[string]interface{})
		c.Contains(errorData["message"], "Error de validación en el campo 'password_confirmation'")
	})

	t.Run("Debería fallar con formato JSON inválido", func(t *testing.T) {
		// Request con formato JSON inválido (string en lugar de object)
		invalidJSON := "invalid_json_string"

		// Realizar la petición POST con JSON inválido
		w := test.MakePostRequest(CONFIG, "/loan-api/api/v1/auth/register", invalidJSON, headers)

		// Verificar código de respuesta de error
		c.Equal(400, w.Code)

		// Decodificar respuesta JSON
		var response map[string]interface{}
		c.NoError(json.Unmarshal(w.Body.Bytes(), &response))

		// Verificar que hay un campo error
		c.Contains(response, "error")

		// Acceder al mensaje dentro del campo error
		errorData := response["error"].(map[string]interface{})
		c.Contains(errorData["message"], "Formato JSON inválido")
	})
}

func TestUserController_Login(t *testing.T) {
	c := require.New(t)

	// Headers válidos
	headers := map[string]string{
		"X-Tenant-ID": "1",
	}

	t.Run("Debería hacer login exitosamente con credenciales válidas", func(t *testing.T) {
		// Cargar datos de prueba
		test.LoadTestData(DB)

		// Datos de login
		requestBody := map[string]interface{}{
			"email":    "juan@example.com",
			"password": "password123!", // Contraseña de los datos de prueba
		}

		// Realizar la petición POST
		w := test.MakePostRequest(CONFIG, "/loan-api/api/v1/auth/login", requestBody, headers)

		// Verificar código de respuesta exitoso
		c.Equal(200, w.Code)

		// Decodificar respuesta JSON
		var response map[string]interface{}
		c.NoError(json.Unmarshal(w.Body.Bytes(), &response))

		// Verificar mensaje de éxito
		c.Equal("Login exitoso", response["message"])

		// Verificar datos en la respuesta
		data := response["data"].(map[string]interface{})
		c.NotEmpty(data["token"])

		user := data["user"].(map[string]interface{})
		c.Equal("juan@example.com", user["email"])
		c.NotEmpty(user["name"])
	})

	t.Run("Debería fallar con credenciales incorrectas", func(t *testing.T) {
		// Cargar datos de prueba
		test.LoadTestData(DB)

		// Datos de login con contraseña incorrecta
		requestBody := map[string]interface{}{
			"email":    "juan@example.com",
			"password": "wrongpassword",
		}

		// Realizar la petición POST
		w := test.MakePostRequest(CONFIG, "/loan-api/api/v1/auth/login", requestBody, headers)

		// Verificar código de respuesta de error
		c.Equal(400, w.Code)

		// Decodificar respuesta JSON
		var response map[string]interface{}
		c.NoError(json.Unmarshal(w.Body.Bytes(), &response))

		// Verificar que hay un campo error
		c.Contains(response, "error")

		// Acceder al mensaje dentro del campo error
		errorData := response["error"].(map[string]interface{})
		c.Contains(errorData["details"], "Email o contraseña incorrectos")
	})

	t.Run("Debería fallar con usuario inexistente", func(t *testing.T) {
		// Limpiar base de datos
		test.ClearDatabase(DB)

		// Datos de login con email inexistente
		requestBody := map[string]interface{}{
			"email":    "noexiste@example.com",
			"password": "password123!",
		}

		// Realizar la petición POST
		w := test.MakePostRequest(CONFIG, "/loan-api/api/v1/auth/login", requestBody, headers)

		// Verificar código de respuesta de error
		c.Equal(400, w.Code)

		// Decodificar respuesta JSON
		var response map[string]interface{}
		c.NoError(json.Unmarshal(w.Body.Bytes(), &response))

		// Verificar que hay un campo error
		c.Contains(response, "error")

		// Acceder al mensaje dentro del campo error
		errorData := response["error"].(map[string]interface{})
		c.Contains(errorData["message"], "Email o contraseña incorrectos")
	})
}

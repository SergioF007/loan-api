package controllers_test

import (
	"encoding/json"
	"testing"

	"loan-api/models"
	"loan-api/test"

	"github.com/stretchr/testify/require"
)

// loginAndGetToken función helper para hacer login y obtener token
func loginAndGetToken(t *testing.T, email, password string) string {
	c := require.New(t)

	// Datos de login
	loginBody := map[string]interface{}{
		"email":    email,
		"password": password,
	}

	// Headers necesarios para login (incluir X-Tenant-ID)
	headers := map[string]string{
		"X-Tenant-ID": "1", // Tenant creado por el seed
	}

	// Hacer login
	w := test.MakePostRequest(CONFIG, "/loan-api/api/v1/auth/login", loginBody, headers)
	c.Equal(200, w.Code)

	// Extraer token de la respuesta
	var loginResponse map[string]interface{}
	c.NoError(json.Unmarshal(w.Body.Bytes(), &loginResponse))

	data := loginResponse["data"].(map[string]interface{})
	token := data["token"].(string)
	c.NotEmpty(token)

	return token
}

func TestLoanController_CreateLoan(t *testing.T) {
	c := require.New(t)

	t.Run("Debería crear préstamo exitosamente con datos válidos", func(t *testing.T) {
		test.LoadTestData(DB)

		// Hacer login para obtener token
		token := loginAndGetToken(t, "juan@example.com", "password123!")

		// Preparar headers con token de autorización y tenant
		headers := map[string]string{
			"Authorization": token,
			"X-Tenant-ID":   "1", // Tenant creado por el seed
		}

		// Verificar préstamos antes de la creación
		loansBefore := test.CountLoans(DB)

		// Datos de la solicitud
		requestBody := map[string]interface{}{
			"loan_type_id": 1, // Loan type creado por el seed
		}

		// Realizar la petición POST
		w := test.MakePostRequest(CONFIG, "/loan-api/api/v1/loans", requestBody, headers)

		// Verificar código de respuesta exitoso
		c.Equal(201, w.Code)

		// Decodificar respuesta JSON
		var response map[string]interface{}
		c.NoError(json.Unmarshal(w.Body.Bytes(), &response))

		// Verificar mensaje de éxito
		c.Equal("Solicitud de préstamo creada exitosamente", response["message"])

		// Verificar datos del préstamo en la respuesta
		data := response["data"].(map[string]interface{})
		c.NotEmpty(data["id"])
		c.Equal(float64(1), data["loan_type_id"])
		c.NotEmpty(data["user_id"])
		c.Equal("pending", data["status"])
		c.Equal("", data["observation"])
		c.Equal("0", data["amount_approved"])
		c.NotEmpty(data["created_at"])

		// Verificar datos del usuario en la respuesta
		user := data["user"].(map[string]interface{})
		c.Equal("juan@example.com", user["email"])
		c.Equal("Juan Pérez", user["name"])

		// Verificar datos del loan type en la respuesta
		loanType := data["loan_type"].(map[string]interface{})
		c.Equal("Préstamo Personal", loanType["name"])
		c.Equal("personal_loan", loanType["code"])

		// Verificar que se creó un nuevo préstamo en la base de datos
		loansAfter := test.CountLoans(DB)
		c.Equal(loansBefore+1, loansAfter)
	})

	t.Run("Debería fallar sin token de autorización", func(t *testing.T) {

		test.LoadTestData(DB)

		headers := map[string]string{
			"X-Tenant-ID": "1",
		}

		requestBody := map[string]interface{}{
			"loan_type_id": 1,
		}

		w := test.MakePostRequest(CONFIG, "/loan-api/api/v1/loans", requestBody, headers)

		c.Equal(401, w.Code)
		var response map[string]interface{}
		c.NoError(json.Unmarshal(w.Body.Bytes(), &response))
		c.Contains(response, "error")
	})

	t.Run("Debería fallar con tenant inexistente", func(t *testing.T) {
		test.LoadTestData(DB)

		token := loginAndGetToken(t, "juan@example.com", "password123!")

		headers := map[string]string{
			"Authorization": token,
			"X-Tenant-ID":   "999", // Tenant que no existe
		}

		requestBody := map[string]interface{}{
			"loan_type_id": 1,
		}

		w := test.MakePostRequest(CONFIG, "/loan-api/api/v1/loans", requestBody, headers)
		c.Equal(403, w.Code)

		var response map[string]interface{}
		c.NoError(json.Unmarshal(w.Body.Bytes(), &response))
		c.Contains(response, "error")
		c.Equal(true, response["error"])
		c.Contains(response, "message")
		c.Equal("Tenant invalido.", response["message"])
	})

	t.Run("Debería fallar con loan_type_id inexistente", func(t *testing.T) {
		test.LoadTestData(DB)

		token := loginAndGetToken(t, "juan@example.com", "password123!")

		headers := map[string]string{
			"Authorization": token,
			"X-Tenant-ID":   "1",
		}

		requestBody := map[string]interface{}{
			"loan_type_id": 999,
		}

		w := test.MakePostRequest(CONFIG, "/loan-api/api/v1/loans", requestBody, headers)

		c.Equal(500, w.Code)

		var response map[string]interface{}
		c.NoError(json.Unmarshal(w.Body.Bytes(), &response))
		c.Contains(response, "error")
		errorData := response["error"].(map[string]interface{})
		c.Contains(errorData["message"], "tipo de préstamo no encontrado")
	})

	t.Run("Debería fallar con datos inválidos - loan_type_id faltante", func(t *testing.T) {
		test.LoadTestData(DB)

		token := loginAndGetToken(t, "juan@example.com", "password123!")

		headers := map[string]string{
			"Authorization": token,
			"X-Tenant-ID":   "1",
		}

		requestBody := map[string]interface{}{}

		w := test.MakePostRequest(CONFIG, "/loan-api/api/v1/loans", requestBody, headers)

		c.Equal(400, w.Code)

		var response map[string]interface{}
		c.NoError(json.Unmarshal(w.Body.Bytes(), &response))

		c.Contains(response, "error")
		errorData := response["error"].(map[string]interface{})
		c.Contains(errorData["message"], "loan_type_id es requerido")
	})

	t.Run("Debería fallar con formato JSON inválido", func(t *testing.T) {
		test.LoadTestData(DB)

		token := loginAndGetToken(t, "juan@example.com", "password123!")

		headers := map[string]string{
			"Authorization": token,
			"X-Tenant-ID":   "1",
		}

		// JSON inválido
		invalidJSON := "invalid_json_string"

		w := test.MakePostRequest(CONFIG, "/loan-api/api/v1/loans", invalidJSON, headers)

		c.Equal(400, w.Code)

		var response map[string]interface{}
		c.NoError(json.Unmarshal(w.Body.Bytes(), &response))

		c.Contains(response, "error")
		errorData := response["error"].(map[string]interface{})
		c.Contains(errorData["message"], "Formato JSON inválido")
	})
}

func TestLoanController_SaveLoanData(t *testing.T) {
	c := require.New(t)

	t.Run("Debería guardar datos de préstamo exitosamente con datos válidos", func(t *testing.T) {
		test.LoadTestData(DB)

		token := loginAndGetToken(t, "juan@example.com", "password123!")

		headers := map[string]string{
			"Authorization": token,
			"X-Tenant-ID":   "1",
		}

		// Datos válidos para guardar
		requestBody := map[string]interface{}{
			"loan_id": 1, // Loan ID del test data
			"data": []map[string]interface{}{
				{
					"form_id": 1,
					"key":     "full_name",
					"value":   "Juan Pérez",
					"index":   0, // Mismo formulario, mismo index
				},
				{
					"form_id": 1,
					"key":     "document_type",
					"value":   "cedula",
					"index":   0,
				},
				{
					"form_id": 1,
					"key":     "document_number",
					"value":   "12345678",
					"index":   0,
				},
			},
		}

		w := test.MakePostRequest(CONFIG, "/loan-api/api/v1/loans/data", requestBody, headers)

		c.Equal(200, w.Code)

		var response map[string]interface{}
		c.NoError(json.Unmarshal(w.Body.Bytes(), &response))

		c.Equal("Datos del préstamo guardados exitosamente", response["message"])
	})

	t.Run("Debería fallar sin token de autorización", func(t *testing.T) {
		test.LoadTestData(DB)

		headers := map[string]string{
			"X-Tenant-ID": "1",
		}

		requestBody := map[string]interface{}{
			"loan_id": 1,
			"data": []map[string]interface{}{
				{
					"form_id": 1,
					"key":     "full_name",
					"value":   "Juan Pérez",
					"index":   0,
				},
			},
		}

		w := test.MakePostRequest(CONFIG, "/loan-api/api/v1/loans/data", requestBody, headers)

		c.Equal(401, w.Code)
		var response map[string]interface{}
		c.NoError(json.Unmarshal(w.Body.Bytes(), &response))
		c.Contains(response, "error")
	})

	t.Run("Debería fallar con tenant inexistente", func(t *testing.T) {
		test.LoadTestData(DB)

		token := loginAndGetToken(t, "juan@example.com", "password123!")

		headers := map[string]string{
			"Authorization": token,
			"X-Tenant-ID":   "999",
		}

		requestBody := map[string]interface{}{
			"loan_id": 1,
			"data": []map[string]interface{}{
				{
					"form_id": 1,
					"key":     "full_name",
					"value":   "Juan Pérez",
					"index":   0,
				},
			},
		}

		w := test.MakePostRequest(CONFIG, "/loan-api/api/v1/loans/data", requestBody, headers)

		c.Equal(403, w.Code)

		var response map[string]interface{}
		c.NoError(json.Unmarshal(w.Body.Bytes(), &response))
		c.Contains(response, "error")
		c.Equal(true, response["error"])
		c.Contains(response, "message")
		c.Equal("Tenant invalido.", response["message"])
	})

	t.Run("Debería fallar con loan_id faltante", func(t *testing.T) {
		test.LoadTestData(DB)

		token := loginAndGetToken(t, "juan@example.com", "password123!")

		headers := map[string]string{
			"Authorization": token,
			"X-Tenant-ID":   "1",
		}

		requestBody := map[string]interface{}{
			"data": []map[string]interface{}{
				{
					"form_id": 1,
					"key":     "full_name",
					"value":   "Juan Pérez",
					"index":   0,
				},
			},
		}

		w := test.MakePostRequest(CONFIG, "/loan-api/api/v1/loans/data", requestBody, headers)

		c.Equal(400, w.Code)

		var response map[string]interface{}
		c.NoError(json.Unmarshal(w.Body.Bytes(), &response))
		c.Contains(response, "error")
		errorData := response["error"].(map[string]interface{})
		c.Contains(errorData["message"], "loan_id es requerido")
	})

	t.Run("Debería fallar con data faltante", func(t *testing.T) {
		test.LoadTestData(DB)

		token := loginAndGetToken(t, "juan@example.com", "password123!")

		headers := map[string]string{
			"Authorization": token,
			"X-Tenant-ID":   "1",
		}

		requestBody := map[string]interface{}{
			"loan_id": 1,
		}

		w := test.MakePostRequest(CONFIG, "/loan-api/api/v1/loans/data", requestBody, headers)

		c.Equal(400, w.Code)

		var response map[string]interface{}
		c.NoError(json.Unmarshal(w.Body.Bytes(), &response))
		c.Contains(response, "error")
		errorData := response["error"].(map[string]interface{})
		c.Contains(errorData["message"], "data es requerido")
	})

	t.Run("Debería fallar con data vacío", func(t *testing.T) {
		test.LoadTestData(DB)

		token := loginAndGetToken(t, "juan@example.com", "password123!")

		headers := map[string]string{
			"Authorization": token,
			"X-Tenant-ID":   "1",
		}

		requestBody := map[string]interface{}{
			"loan_id": 1,
			"data":    []map[string]interface{}{},
		}

		w := test.MakePostRequest(CONFIG, "/loan-api/api/v1/loans/data", requestBody, headers)

		c.Equal(400, w.Code)

		var response map[string]interface{}
		c.NoError(json.Unmarshal(w.Body.Bytes(), &response))
		c.Contains(response, "error")
		errorData := response["error"].(map[string]interface{})
		c.Contains(errorData["message"], "data es requerido")
	})

	t.Run("Debería fallar con loan_id inexistente", func(t *testing.T) {
		test.LoadTestData(DB)

		token := loginAndGetToken(t, "juan@example.com", "password123!")

		headers := map[string]string{
			"Authorization": token,
			"X-Tenant-ID":   "1",
		}

		requestBody := map[string]interface{}{
			"loan_id": 999, // Loan ID que no existe
			"data": []map[string]interface{}{
				{
					"form_id": 1,
					"key":     "full_name",
					"value":   "Juan Pérez",
					"index":   0,
				},
			},
		}

		w := test.MakePostRequest(CONFIG, "/loan-api/api/v1/loans/data", requestBody, headers)

		c.Equal(500, w.Code)

		var response map[string]interface{}
		c.NoError(json.Unmarshal(w.Body.Bytes(), &response))
		c.Contains(response, "error")
		errorData := response["error"].(map[string]interface{})
		c.NotEmpty(errorData["message"])
	})

	t.Run("Debería fallar con formato JSON inválido", func(t *testing.T) {
		test.LoadTestData(DB)

		token := loginAndGetToken(t, "juan@example.com", "password123!")

		headers := map[string]string{
			"Authorization": token,
			"X-Tenant-ID":   "1",
		}

		invalidJSON := "invalid_json_string"

		w := test.MakePostRequest(CONFIG, "/loan-api/api/v1/loans/data", invalidJSON, headers)

		c.Equal(400, w.Code)

		var response map[string]interface{}
		c.NoError(json.Unmarshal(w.Body.Bytes(), &response))
		c.Contains(response, "error")
		errorData := response["error"].(map[string]interface{})
		c.Contains(errorData["message"], "Formato JSON inválido")
	})

	t.Run("Debería cambiar estado a completed con todos los campos requeridos", func(t *testing.T) {
		test.LoadTestData(DB)

		token := loginAndGetToken(t, "juan@example.com", "password123!")

		headers := map[string]string{
			"Authorization": token,
			"X-Tenant-ID":   "1",
		}

		// Datos completos para todos los formularios requeridos
		requestBody := map[string]interface{}{
			"loan_id": 1, // Loan ID del test data
			"data": []map[string]interface{}{
				// Información Personal (form_id: 1)
				{
					"form_id": 1,
					"key":     "full_name",
					"value":   "Juan Pérez",
					"index":   0,
				},
				{
					"form_id": 1,
					"key":     "document_type",
					"value":   "cedula",
					"index":   0,
				},
				{
					"form_id": 1,
					"key":     "document_number",
					"value":   "12345678",
					"index":   0,
				},
				{
					"form_id": 1,
					"key":     "age",
					"value":   "30",
					"index":   0,
				},
				// Información Financiera (form_id: 2)
				{
					"form_id": 2,
					"key":     "monthly_income",
					"value":   "5000000",
					"index":   0,
				},
				{
					"form_id": 2,
					"key":     "monthly_expenses",
					"value":   "2000000",
					"index":   0,
				},
				// Detalles del Préstamo (form_id: 3)
				{
					"form_id": 3,
					"key":     "requested_amount",
					"value":   "2000000",
					"index":   0,
				},
				{
					"form_id": 3,
					"key":     "purpose",
					"value":   "Educación",
					"index":   0,
				},
			},
		}

		w := test.MakePostRequest(CONFIG, "/loan-api/api/v1/loans/data", requestBody, headers)

		c.Equal(200, w.Code)

		var response map[string]interface{}
		c.NoError(json.Unmarshal(w.Body.Bytes(), &response))

		c.Equal("Datos del préstamo guardados exitosamente", response["message"])

		// Verificar que el préstamo cambió a estado "completed"
		// Obtener el préstamo actualizado de la base de datos
		var loan models.Loan
		err := DB.Preload("Data").First(&loan, 1).Error
		c.NoError(err)

		c.Equal("completed", loan.Status)
		c.NotNil(loan.CreditScore)
		c.NotNil(loan.IdentityVerified)
		c.Contains(loan.Observation, "Solicitud completada")
		c.Contains(loan.Observation, "Score crediticio:")
		c.Contains(loan.Observation, "Verificación de identidad:")
	})
}

func TestLoanController_ProcessLoanDecision(t *testing.T) {
	c := require.New(t)

	t.Run("Debería procesar decisión exitosamente con préstamo completed", func(t *testing.T) {
		test.LoadTestData(DB)

		// Primero crear un préstamo con estado completed
		token := loginAndGetToken(t, "juan@example.com", "password123!")

		headers := map[string]string{
			"Authorization": token,
			"X-Tenant-ID":   "1",
		}

		// Guardar datos completos para que el préstamo pase a completed
		requestBody := map[string]interface{}{
			"loan_id": 1,
			"data": []map[string]interface{}{
				// Información Personal
				{"form_id": 1, "key": "full_name", "value": "Juan Pérez", "index": 0},
				{"form_id": 1, "key": "document_type", "value": "cedula", "index": 0},
				{"form_id": 1, "key": "document_number", "value": "12345678", "index": 0},
				{"form_id": 1, "key": "age", "value": "30", "index": 0},
				// Información Financiera
				{"form_id": 2, "key": "monthly_income", "value": "5000000", "index": 0},
				{"form_id": 2, "key": "monthly_expenses", "value": "2000000", "index": 0},
				// Detalles del Préstamo
				{"form_id": 3, "key": "requested_amount", "value": "2000000", "index": 0},
				{"form_id": 3, "key": "purpose", "value": "Educación", "index": 0},
			},
		}

		// Guardar datos para completar el préstamo
		w := test.MakePostRequest(CONFIG, "/loan-api/api/v1/loans/data", requestBody, headers)
		c.Equal(200, w.Code)

		// Verificar que el préstamo esté en estado completed
		var loan models.Loan
		err := DB.First(&loan, 1).Error
		c.NoError(err)
		c.Equal("completed", loan.Status)

		// Ahora procesar la decisión
		w = test.MakePostRequest(CONFIG, "/loan-api/api/v1/loans/1/decision", nil, headers)

		c.Equal(200, w.Code)

		var response map[string]interface{}
		c.NoError(json.Unmarshal(w.Body.Bytes(), &response))

		c.Equal("Decisión del préstamo procesada exitosamente", response["message"])
		c.Contains(response, "data")

		// Verificar que el préstamo cambió a approved o rejected
		err = DB.First(&loan, 1).Error
		c.NoError(err)
		c.True(loan.Status == "approved" || loan.Status == "rejected")
		c.NotEmpty(loan.Observation)
	})

	t.Run("Debería fallar sin token de autorización", func(t *testing.T) {
		test.LoadTestData(DB)

		headers := map[string]string{
			"X-Tenant-ID": "1",
		}

		w := test.MakePostRequest(CONFIG, "/loan-api/api/v1/loans/1/decision", nil, headers)

		c.Equal(401, w.Code)
		var response map[string]interface{}
		c.NoError(json.Unmarshal(w.Body.Bytes(), &response))
		c.Contains(response, "error")
	})

	t.Run("Debería fallar con ID de préstamo inválido", func(t *testing.T) {
		test.LoadTestData(DB)

		token := loginAndGetToken(t, "juan@example.com", "password123!")

		headers := map[string]string{
			"Authorization": token,
			"X-Tenant-ID":   "1",
		}

		w := test.MakePostRequest(CONFIG, "/loan-api/api/v1/loans/invalid/decision", nil, headers)

		c.Equal(400, w.Code)

		var response map[string]interface{}
		c.NoError(json.Unmarshal(w.Body.Bytes(), &response))
		c.Equal(false, response["success"])
		c.Contains(response, "error")
		c.Contains(response, "message")
		errorData := response["error"].(map[string]interface{})
		c.Contains(errorData, "message")
		c.Contains(errorData["message"], "ID del préstamo debe ser un número válido")
	})

	t.Run("Debería fallar con préstamo inexistente", func(t *testing.T) {
		test.LoadTestData(DB)

		token := loginAndGetToken(t, "juan@example.com", "password123!")

		headers := map[string]string{
			"Authorization": token,
			"X-Tenant-ID":   "1",
		}

		w := test.MakePostRequest(CONFIG, "/loan-api/api/v1/loans/999/decision", nil, headers)

		c.Equal(500, w.Code)

		var response map[string]interface{}
		c.NoError(json.Unmarshal(w.Body.Bytes(), &response))
		c.Equal(false, response["success"])
		c.Contains(response, "error")
		c.Contains(response, "message")
		errorData := response["error"].(map[string]interface{})
		c.Contains(errorData, "message")
		c.Contains(errorData["message"], "préstamo no encontrado")
	})

	t.Run("Debería fallar con préstamo en estado pending", func(t *testing.T) {
		test.LoadTestData(DB)

		token := loginAndGetToken(t, "juan@example.com", "password123!")

		headers := map[string]string{
			"Authorization": token,
			"X-Tenant-ID":   "1",
		}

		// Verificar que el préstamo esté en estado pending (estado inicial)
		var loan models.Loan
		err := DB.First(&loan, 1).Error
		c.NoError(err)
		c.Equal("pending", loan.Status)

		// Intentar procesar decisión sin completar el préstamo
		w := test.MakePostRequest(CONFIG, "/loan-api/api/v1/loans/1/decision", nil, headers)

		c.Equal(500, w.Code)

		var response map[string]interface{}
		c.NoError(json.Unmarshal(w.Body.Bytes(), &response))
		c.Equal(false, response["success"])
		c.Contains(response, "error")
		c.Contains(response, "message")
		errorData := response["error"].(map[string]interface{})
		c.Contains(errorData, "message")
		c.Contains(errorData["message"], "solo se pueden evaluar préstamos en estado completado")
	})

	t.Run("Debería fallar con préstamo en estado on_progress", func(t *testing.T) {
		test.LoadTestData(DB)

		token := loginAndGetToken(t, "juan@example.com", "password123!")

		headers := map[string]string{
			"Authorization": token,
			"X-Tenant-ID":   "1",
		}

		// Guardar datos parciales para que el préstamo pase a on_progress
		requestBody := map[string]interface{}{
			"loan_id": 1,
			"data": []map[string]interface{}{
				{"form_id": 1, "key": "full_name", "value": "Juan Pérez", "index": 0},
				{"form_id": 1, "key": "document_type", "value": "cedula", "index": 0},
				// Faltan campos requeridos
			},
		}

		// Guardar datos parciales
		w := test.MakePostRequest(CONFIG, "/loan-api/api/v1/loans/data", requestBody, headers)
		c.Equal(200, w.Code)

		// Verificar que el préstamo esté en estado on_progress
		var loan models.Loan
		err := DB.First(&loan, 1).Error
		c.NoError(err)
		c.Equal("on_progress", loan.Status)

		// Intentar procesar decisión con préstamo incompleto
		w = test.MakePostRequest(CONFIG, "/loan-api/api/v1/loans/1/decision", nil, headers)

		c.Equal(500, w.Code)

		var response map[string]interface{}
		c.NoError(json.Unmarshal(w.Body.Bytes(), &response))
		c.Equal(false, response["success"])
		c.Contains(response, "error")
		c.Contains(response, "message")
		errorData := response["error"].(map[string]interface{})
		c.Contains(errorData, "message")
		c.Contains(errorData["message"], "solo se pueden evaluar préstamos en estado completado")
	})
}

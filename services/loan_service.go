package services

import (
	"errors"
	"loan-api/models"
	"loan-api/repositories"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

// LoanService interface para el servicio de préstamos
type LoanService interface {
	CreateLoan(userID uint, request models.CreateLoanRequest) (*models.LoanResponse, error)
	SaveLoanData(request models.SaveLoanDataRequest) error
	ProcessLoanDecision(loanID uint) (*models.LoanResponse, error)
	GetLoanByID(id uint) (*models.LoanResponse, error)
	GetLoansByUserID(userID uint) ([]models.LoanResponse, error)
}

// loanService implementación del servicio
type loanService struct {
	loanRepo     repositories.LoanRepository
	userRepo     repositories.UserRepository
	loanTypeRepo repositories.LoanTypeRepository
}

// NewLoanService crea una nueva instancia del servicio
func NewLoanService(loanRepo repositories.LoanRepository, userRepo repositories.UserRepository, loanTypeRepo repositories.LoanTypeRepository) LoanService {
	return &loanService{
		loanRepo:     loanRepo,
		userRepo:     userRepo,
		loanTypeRepo: loanTypeRepo,
	}
}

// CreateLoan crea un nuevo préstamo
func (s *loanService) CreateLoan(userID uint, request models.CreateLoanRequest) (*models.LoanResponse, error) {
	// Validar que el usuario existe
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("usuario no encontrado")
	}

	// Validar que el tipo de préstamo existe
	loanType, err := s.loanTypeRepo.GetByIDWithForms(request.LoanTypeID)
	if err != nil {
		return nil, errors.New("tipo de préstamo no encontrado")
	}

	// Crear el préstamo
	loan := &models.Loan{
		LoanTypeID:     request.LoanTypeID,
		UserID:         userID,
		Status:         "pending",
		Observation:    "",
		AmountApproved: decimal.NewFromFloat(0),
	}

	if err := s.loanRepo.Create(loan); err != nil {
		return nil, errors.New("error al crear el préstamo")
	}

	// Obtener el préstamo creado con todas sus relaciones
	createdLoan, err := s.loanRepo.GetByID(loan.ID)
	if err != nil {
		return nil, errors.New("error al obtener el préstamo creado")
	}

	// Construir la respuesta
	response := s.buildLoanResponse(*createdLoan, *user, *loanType)
	return &response, nil
}

// SaveLoanData guarda los datos de un préstamo
func (s *loanService) SaveLoanData(request models.SaveLoanDataRequest) error {
	// Validar que el préstamo existe
	loan, err := s.loanRepo.GetByID(request.LoanID)
	if err != nil {
		return errors.New("préstamo no encontrado")
	}

	// Validar que el préstamo está en estado pendiente o en progreso
	if loan.Status != "pending" && loan.Status != "on_progress" {
		return errors.New("solo se pueden actualizar préstamos en estado pendiente o en progreso")
	}

	// Extraer datos necesarios para validaciones
	documentType := s.extractLoanDataValue(request.Data, "document_type")
	documentNumber := s.extractLoanDataValue(request.Data, "document_number")
	fullName := s.extractLoanDataValue(request.Data, "full_name")

	// Procesar validaciones solo si tenemos los datos necesarios
	var creditScore *int
	var identityVerified *bool

	if documentType != "" && documentNumber != "" {
		// 1. Simulación del score crediticio
		score, err := s.simulateCreditScore(documentType, documentNumber)
		if err != nil {
			return errors.New("error al consultar el score crediticio: " + err.Error())
		}
		creditScore = &score

		// 2. Verificación de identidad
		if fullName != "" {
			verified, err := s.verifyIdentity(loan.UserID, documentType, documentNumber, fullName)
			if err != nil {
				// Solo falla si hay errores técnicos (datos insuficientes, problemas de BD, etc.)
				return errors.New("error al verificar la identidad: " + err.Error())
			}
			// Si no hay error técnico, usar el resultado de la verificación (true/false)
			identityVerified = &verified
		}
	}

	// Actualizar los campos de validación del préstamo
	if creditScore != nil {
		loan.CreditScore = creditScore
	}
	if identityVerified != nil {
		loan.IdentityVerified = identityVerified
	}

	// Eliminar datos existentes
	if err := s.loanRepo.DeleteLoanDataByLoanID(request.LoanID); err != nil {
		return errors.New("error al eliminar datos anteriores")
	}

	// Preparar los nuevos datos
	loanDataList := make([]models.LoanData, len(request.Data))
	for i, dataItem := range request.Data {
		loanDataList[i] = models.LoanData{
			LoanID:    request.LoanID,
			FormID:    dataItem.FormID,
			Key:       dataItem.Key,
			Value:     dataItem.Value,
			Index:     dataItem.Index,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	}

	// Guardar los nuevos datos
	if err := s.loanRepo.SaveLoanData(loanDataList); err != nil {
		return errors.New("error al guardar los datos del préstamo")
	}

	// DESPUÉS de guardar los datos, determinar el nuevo estado basado en la configuración real
	newStatus, err := s.determineNewLoanStatusFromDB(request.LoanID, loan.LoanTypeID, creditScore, identityVerified)
	if err != nil {
		return errors.New("error al determinar el estado del préstamo: " + err.Error())
	}

	// Actualizar estado y observación
	loan.Status = newStatus
	loan.Observation = s.generateStatusObservation(newStatus, creditScore, identityVerified)

	// Guardar los cambios finales
	if err := s.loanRepo.Update(loan); err != nil {
		return errors.New("error al actualizar el estado del préstamo")
	}

	return nil
}

// GetLoanByID obtiene un préstamo por ID
func (s *loanService) GetLoanByID(id uint) (*models.LoanResponse, error) {
	loan, err := s.loanRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("préstamo no encontrado")
	}

	// Obtener usuario y tipo de préstamo para la respuesta completa
	user, err := s.userRepo.GetByID(loan.UserID)
	if err != nil {
		return nil, errors.New("usuario no encontrado")
	}

	loanType, err := s.loanTypeRepo.GetByIDWithForms(loan.LoanTypeID)
	if err != nil {
		return nil, errors.New("tipo de préstamo no encontrado")
	}

	response := s.buildLoanResponse(*loan, *user, *loanType)
	return &response, nil
}

// GetLoansByUserID obtiene todos los préstamos de un usuario
func (s *loanService) GetLoansByUserID(userID uint) ([]models.LoanResponse, error) {
	loans, err := s.loanRepo.GetByUserID(userID)
	if err != nil {
		return nil, errors.New("error al obtener préstamos del usuario")
	}

	// Obtener usuario una sola vez
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("usuario no encontrado")
	}

	response := make([]models.LoanResponse, len(loans))
	for i, loan := range loans {
		// Obtener tipo de préstamo para cada préstamo
		loanType, err := s.loanTypeRepo.GetByIDWithForms(loan.LoanTypeID)
		if err != nil {
			continue // Saltar préstamos con tipos no encontrados
		}

		response[i] = s.buildLoanResponse(loan, *user, *loanType)
	}

	return response, nil
}

// simulateCreditScore simula la consulta de score crediticio basado en el tipo y número de documento
func (s *loanService) simulateCreditScore(documentType, documentNumber string) (int, error) {
	// Simulación basada en el número de documento para tener resultados consistentes
	// En un escenario real, esto sería una llamada a un servicio externo

	// Generar un score basado en el hash del número de documento para consistencia
	score := 300 // Score mínimo

	// Usar los últimos dígitos del documento para generar variabilidad
	if len(documentNumber) > 0 {
		lastDigit := documentNumber[len(documentNumber)-1:]
		if digit, err := strconv.Atoi(lastDigit); err == nil {
			// Mapear el último dígito a un rango de score
			switch digit {
			case 0, 1:
				score = 300 + (digit * 50) // 300-350 (score bajo)
			case 2, 3, 4:
				score = 400 + ((digit - 2) * 50) // 400-500 (score medio-bajo)
			case 5, 6, 7:
				score = 550 + ((digit - 5) * 50) // 550-650 (score medio)
			case 8, 9:
				score = 700 + ((digit - 8) * 50) // 700-750 (score alto)
			}
		}
	}

	// Agregar algo de variabilidad aleatoria (+/- 25 puntos)
	variation := rand.Intn(51) - 25 // -25 a +25
	score += variation

	// Asegurar que el score esté en el rango válido (300-850)
	if score < 300 {
		score = 300
	}
	if score > 850 {
		score = 850
	}

	return score, nil
}

// verifyIdentity verifica la identidad del solicitante comparando con los datos del registro
func (s *loanService) verifyIdentity(userID uint, documentType, documentNumber, fullName string) (bool, error) {
	// Validaciones básicas - error técnico si faltan datos
	if documentType == "" || documentNumber == "" || fullName == "" {
		return false, errors.New("datos insuficientes para verificación de identidad")
	}

	// Obtener los datos del usuario registrado - error técnico si no se puede obtener
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return false, errors.New("error al obtener datos del usuario registrado")
	}

	// Las siguientes son verificaciones de identidad, no errores técnicos
	// Si fallan, retornan false pero no error

	// Verificar que el tipo de documento coincida
	if string(user.DocumentType) != documentType {
		return false, nil // Fallo de verificación: tipo de documento no coincide
	}

	// Verificar que el número de documento coincida
	if user.DocumentNumber != documentNumber {
		return false, nil // Fallo de verificación: número de documento no coincide
	}

	// Verificar que el nombre del registro esté contenido en el nombre completo
	userNameNormalized := strings.ToLower(strings.TrimSpace(user.Name))
	inputNameNormalized := strings.ToLower(strings.TrimSpace(fullName))

	if !strings.Contains(inputNameNormalized, userNameNormalized) {
		return false, nil // Fallo de verificación: nombre no está contenido
	}

	// Todas las verificaciones pasaron exitosamente
	return true, nil
}

// extractLoanDataValue extrae un valor específico de los datos del préstamo
func (s *loanService) extractLoanDataValue(data []models.LoanDataItemRequest, key string) string {
	for _, item := range data {
		if item.Key == key {
			return item.Value
		}
	}
	return ""
}

// buildLoanResponse construye la respuesta del préstamo
func (s *loanService) buildLoanResponse(loan models.Loan, user models.User, loanType models.LoanType) models.LoanResponse {
	// Construir datos del préstamo
	dataResponse := make([]models.LoanDataResponse, len(loan.Data))
	for i, data := range loan.Data {
		dataResponse[i] = models.LoanDataResponse{
			ID:     data.ID,
			FormID: data.FormID,
			Key:    data.Key,
			Value:  data.Value,
			Index:  data.Index,
		}
	}

	// Construir respuesta del usuario
	userResponse := models.UserResponse{
		ID:             user.ID,
		Name:           user.Name,
		Email:          user.Email,
		Phone:          user.Phone,
		DocumentType:   user.DocumentType,
		DocumentNumber: user.DocumentNumber,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}

	// Construir respuesta del tipo de préstamo
	loanTypeResponse := models.LoanTypeResponse{
		ID:          loanType.ID,
		Name:        loanType.Name,
		Code:        loanType.Code,
		Description: loanType.Description,
		MinAmount:   loanType.MinAmount,
		MaxAmount:   loanType.MaxAmount,
	}

	return models.LoanResponse{
		ID:               loan.ID,
		LoanTypeID:       loan.LoanTypeID,
		LoanType:         loanTypeResponse,
		UserID:           loan.UserID,
		User:             userResponse,
		Status:           loan.Status,
		Observation:      loan.Observation,
		AmountApproved:   loan.AmountApproved,
		CreditScore:      loan.CreditScore,
		IdentityVerified: loan.IdentityVerified,
		Data:             dataResponse,
		CreatedAt:        loan.CreatedAt,
		UpdatedAt:        loan.UpdatedAt,
	}
}

// ProcessLoanDecision evalúa el préstamo y toma la decisión final de aprobación/rechazo
func (s *loanService) ProcessLoanDecision(loanID uint) (*models.LoanResponse, error) {
	// Obtener el préstamo con todos sus datos
	loan, err := s.loanRepo.GetByID(loanID)
	if err != nil {
		return nil, errors.New("préstamo no encontrado")
	}

	// Validar que el préstamo esté en estado completed (listo para evaluación)
	if loan.Status != "completed" {
		return nil, errors.New("solo se pueden evaluar préstamos en estado completado")
	}

	// Validar que se hayan completado las validaciones previas
	if loan.CreditScore == nil {
		return nil, errors.New("el préstamo debe tener un score crediticio calculado")
	}

	if loan.IdentityVerified == nil {
		return nil, errors.New("el préstamo debe tener la verificación de identidad procesada")
	}

	// Obtener información adicional necesaria
	requestedAmount := s.extractLoanDataFromLoan(*loan, "requested_amount")
	monthlyIncome := s.extractLoanDataFromLoan(*loan, "monthly_income")

	// Aplicar reglas de negocio para la decisión
	decision, reason := s.evaluateLoanApproval(*loan.CreditScore, *loan.IdentityVerified, requestedAmount, monthlyIncome)

	// Actualizar el estado del préstamo
	loan.Status = decision
	loan.Observation = reason

	// Si es aprobado, calcular monto aprobado y simular desembolso
	if decision == "approved" {
		approvedAmount := s.calculateApprovedAmount(requestedAmount, decimal.NewFromInt(int64(*loan.CreditScore)), monthlyIncome)
		loan.AmountApproved = approvedAmount

		// Simular desembolso
		disbursementSuccess := s.simulateDisbursement(loan.UserID, approvedAmount)
		if !disbursementSuccess {
			// Si falla el desembolso, rechazar el préstamo
			loan.Status = "rejected"
			loan.Observation = "Préstamo aprobado pero falló el desembolso. Contacte soporte."
		} else {
			loan.Observation += " - Desembolso realizado exitosamente"
		}
	}

	// Guardar los cambios
	if err := s.loanRepo.Update(loan); err != nil {
		return nil, errors.New("error al actualizar el estado del préstamo")
	}

	// Obtener préstamo actualizado para la respuesta
	updatedLoan, err := s.loanRepo.GetByID(loanID)
	if err != nil {
		return nil, errors.New("error al obtener el préstamo actualizado")
	}

	// Construir respuesta
	user, err := s.userRepo.GetByID(updatedLoan.UserID)
	if err != nil {
		return nil, errors.New("error al obtener usuario")
	}

	loanType, err := s.loanTypeRepo.GetByIDWithForms(updatedLoan.LoanTypeID)
	if err != nil {
		return nil, errors.New("error al obtener tipo de préstamo")
	}

	response := s.buildLoanResponse(*updatedLoan, *user, *loanType)
	return &response, nil
}

// extractLoanDataFromLoan extrae un valor específico de los datos del préstamo
func (s *loanService) extractLoanDataFromLoan(loan models.Loan, key string) decimal.Decimal {
	for _, data := range loan.Data {
		if data.Key == key {
			value, _ := decimal.NewFromString(data.Value)
			return value
		}
	}
	return decimal.NewFromFloat(0)
}

// evaluateLoanApproval evalúa la aprobación del préstamo basado en el score y la verificación de identidad
func (s *loanService) evaluateLoanApproval(creditScore int, identityVerified bool, requestedAmount, monthlyIncome decimal.Decimal) (string, string) {
	// 1. Verificar que la identidad esté verificada (requisito obligatorio)
	if !identityVerified {
		return "rejected", "Préstamo rechazado: verificación de identidad fallida"
	}

	// 2. Evaluar score crediticio mínimo
	if creditScore < 400 {
		return "rejected", "Préstamo rechazado: score crediticio muy bajo (" + strconv.Itoa(creditScore) + ")"
	}

	// 3. Verificar capacidad de pago (máximo 50% del ingreso mensual)
	maxAmount := monthlyIncome.Mul(decimal.NewFromFloat(0.5))
	if requestedAmount.GreaterThan(maxAmount) && creditScore < 650 {
		return "rejected", "Préstamo rechazado: monto solicitado excede capacidad de pago y score insuficiente"
	}

	// 4. Verificar monto mínimo
	if requestedAmount.LessThan(decimal.NewFromFloat(100000)) {
		return "rejected", "Préstamo rechazado: monto mínimo no alcanzado"
	}

	// 5. Evaluación por rangos de score
	var reason string
	switch {
	case creditScore >= 700:
		reason = "Préstamo aprobado: excelente score crediticio (" + strconv.Itoa(creditScore) + ")"
	case creditScore >= 600:
		reason = "Préstamo aprobado: buen score crediticio (" + strconv.Itoa(creditScore) + ")"
	case creditScore >= 500:
		reason = "Préstamo aprobado: score crediticio aceptable (" + strconv.Itoa(creditScore) + ")"
	default:
		reason = "Préstamo aprobado: score crediticio bajo pero dentro del rango aceptable (" + strconv.Itoa(creditScore) + ")"
	}

	return "approved", reason
}

// calculateApprovedAmount calcula el monto aprobado del préstamo
func (s *loanService) calculateApprovedAmount(requestedAmount, creditScore, monthlyIncome decimal.Decimal) decimal.Decimal {
	// Implementa la lógica para calcular el monto aprobado del préstamo basado en el score y el ingreso mensual
	// Este es un ejemplo básico, puedes implementar una lógica más robusta basada en tus reglas de negocio

	// Ejemplo: Si el monto solicitado es menor o igual al 50% del ingreso mensual, aprobar el monto completo
	maxAmount := monthlyIncome.Mul(decimal.NewFromFloat(0.5))
	if requestedAmount.LessThanOrEqual(maxAmount) {
		return requestedAmount
	}

	// Si el monto solicitado excede el 50% del ingreso, aprobar solo el monto máximo
	return maxAmount
}

// simulateDisbursement simula el desembolso del préstamo
func (s *loanService) simulateDisbursement(userID uint, amount decimal.Decimal) bool {
	// Simulación de desembolso más realista
	// En un escenario real, esto sería una llamada a un servicio de pagos/bancario

	// 1. Verificar que el monto sea válido
	if amount.LessThanOrEqual(decimal.NewFromFloat(0)) {
		return false // Falla: monto inválido
	}

	// 2. Simulación basada en el ID del usuario para consistencia en pruebas
	userIDStr := strconv.Itoa(int(userID))
	lastDigit := userIDStr[len(userIDStr)-1:]

	// 3. Simular fallos de desembolso ocasionales (10% de probabilidad)
	if lastDigit == "0" {
		return false // Simular falla del sistema bancario
	}

	// 4. Simular límites de desembolso diario
	if amount.GreaterThan(decimal.NewFromFloat(50000000)) { // 50 millones
		return false // Excede límite diario de desembolso
	}

	// 5. Desembolso exitoso para todos los otros casos
	// En un escenario real, aquí se haría:
	// - Llamada al API bancario
	// - Registro de la transacción
	// - Notificación al usuario
	return true
}

// determineNewLoanStatusFromDB determina el nuevo estado del préstamo basado en la configuración real de formularios
func (s *loanService) determineNewLoanStatusFromDB(loanID, loanTypeID uint, creditScore *int, identityVerified *bool) (string, error) {
	// Obtener el préstamo actualizado con todos sus datos
	loan, err := s.loanRepo.GetByID(loanID)
	if err != nil {
		return "", errors.New("error al obtener el préstamo")
	}

	// Si no hay datos guardados, mantener pending
	if len(loan.Data) == 0 {
		return "pending", nil
	}

	// Obtener la configuración de formularios para este loan type
	loanType, err := s.loanTypeRepo.GetByIDWithForms(loanTypeID)
	if err != nil {
		return "", errors.New("error al obtener configuración de formularios")
	}

	// Verificar si todos los campos requeridos están completos
	allRequiredFieldsComplete := s.checkAllRequiredFieldsComplete(*loan, *loanType)

	// Determinar estado basado en completitud de campos y validaciones
	if allRequiredFieldsComplete && creditScore != nil && identityVerified != nil {
		return "completed", nil
	} else if len(loan.Data) > 0 {
		return "on_progress", nil
	}

	return "pending", nil
}

// checkAllRequiredFieldsComplete verifica si todos los campos requeridos están completos
func (s *loanService) checkAllRequiredFieldsComplete(loan models.Loan, loanType models.LoanType) bool {
	// Crear un mapa de los datos guardados para búsqueda rápida
	savedData := make(map[string]map[uint]string) // key -> index -> value
	for _, data := range loan.Data {
		if savedData[data.Key] == nil {
			savedData[data.Key] = make(map[uint]string)
		}
		savedData[data.Key][data.Index] = data.Value
	}

	// Verificar cada formulario y sus inputs requeridos
	for _, version := range loanType.Versions {
		if !version.IsDefault {
			continue // Solo verificar la versión por defecto
		}

		for _, form := range version.Forms {
			if !form.IsRequired {
				continue // Solo verificar formularios requeridos
			}

			for _, input := range form.FormInputs {
				if !input.IsRequired {
					continue // Solo verificar inputs requeridos
				}

				// Verificar si este input tiene al menos un valor guardado
				inputValues, exists := savedData[input.Code]
				if !exists || len(inputValues) == 0 {
					return false // Falta un campo requerido
				}

				// Verificar que al menos un valor no esté vacío
				hasNonEmptyValue := false
				for _, value := range inputValues {
					if strings.TrimSpace(value) != "" {
						hasNonEmptyValue = true
						break
					}
				}

				if !hasNonEmptyValue {
					return false // Campo requerido está vacío
				}
			}
		}
	}

	return true // Todos los campos requeridos están completos
}

// generateStatusObservation genera la observación basada en el estado y validaciones
func (s *loanService) generateStatusObservation(status string, creditScore *int, identityVerified *bool) string {
	switch status {
	case "pending":
		return "Solicitud creada, esperando datos"
	case "on_progress":
		return "Datos parciales guardados, completar información faltante"
	case "completed":
		observation := "Solicitud completada."
		if creditScore != nil {
			observation += " Score crediticio: " + strconv.Itoa(*creditScore) + "."
		}
		if identityVerified != nil {
			if *identityVerified {
				observation += " Verificación de identidad: exitosa."
			} else {
				observation += " Verificación de identidad: fallida."
			}
		}
		return observation
	default:
		return ""
	}
}

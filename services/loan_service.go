package services

import (
	"errors"
	"loan-api/models"
	"loan-api/repositories"
	"time"

	"github.com/shopspring/decimal"
)

// LoanService interface para el servicio de préstamos
type LoanService interface {
	CreateLoan(userID uint, request models.CreateLoanRequest) (*models.LoanResponse, error)
	SaveLoanData(request models.SaveLoanDataRequest) error
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
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
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

	// Validar que el préstamo está en estado pendiente
	if loan.Status != "pending" {
		return errors.New("solo se pueden actualizar préstamos en estado pendiente")
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
		ID:             loan.ID,
		LoanTypeID:     loan.LoanTypeID,
		LoanType:       loanTypeResponse,
		UserID:         loan.UserID,
		User:           userResponse,
		Status:         loan.Status,
		Observation:    loan.Observation,
		AmountApproved: loan.AmountApproved,
		Data:           dataResponse,
		CreatedAt:      loan.CreatedAt,
		UpdatedAt:      loan.UpdatedAt,
	}
}

package test

import (
	"log"
	"time"

	"github.com/shopspring/decimal"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"loan-api/models"
)

// ClearDatabase limpia todas las tablas y resetea los AUTO_INCREMENT
func ClearDatabase(DB *gorm.DB) {
	if DB == nil {
		log.Println("Warning: Database connection is nil")
		return
	}

	DB.Exec("SET foreign_key_checks = 0")

	// Limpiar préstamos primero (tiene FK a users)
	DB.Exec("DELETE FROM loans")
	DB.Exec("ALTER TABLE loans AUTO_INCREMENT = 1")

	// Limpiar usuarios
	DB.Exec("DELETE FROM users")
	DB.Exec("ALTER TABLE users AUTO_INCREMENT = 1")

	DB.Exec("SET foreign_key_checks = 1")

	log.Println("Database cleared successfully")
}

// LoadUsers carga usuarios de prueba en la base de datos
func LoadUsers(DB *gorm.DB) {
	if DB == nil {
		log.Println("Warning: Database connection is nil")
		return
	}

	// Hash de contraseña de prueba
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	users := []models.User{
		{
			Name:           "Juan Pérez",
			Email:          "juan@example.com",
			Phone:          "3001234567",
			DocumentType:   models.DocumentTypeCedula,
			DocumentNumber: "12345678",
			Password:       string(hashedPassword),
			IP:             "192.168.1.1",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
		{
			Name:           "María García",
			Email:          "maria@example.com",
			Phone:          "3009876543",
			DocumentType:   models.DocumentTypeCedula,
			DocumentNumber: "87654321",
			Password:       string(hashedPassword),
			IP:             "192.168.1.2",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
		{
			Name:           "Carlos López",
			Email:          "carlos@example.com",
			Phone:          "3012345678",
			DocumentType:   models.DocumentTypePasaporte,
			DocumentNumber: "AB123456",
			Password:       string(hashedPassword),
			IP:             "192.168.1.3",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
		{
			Name:           "Ana Rodríguez",
			Email:          "ana@example.com",
			Phone:          "3098765432",
			DocumentType:   models.DocumentTypeCedula,
			DocumentNumber: "11223344",
			Password:       string(hashedPassword),
			IP:             "192.168.1.4",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
		{
			Name:           "Diego Martínez",
			Email:          "diego@example.com",
			Phone:          "3011111111",
			DocumentType:   models.DocumentTypeTarjetaIdentidad,
			DocumentNumber: "TI556677",
			Password:       string(hashedPassword),
			IP:             "192.168.1.5",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
	}

	for _, user := range users {
		if err := DB.Create(&user).Error; err != nil {
			log.Printf("Error creating user %s: %v", user.Email, err)
		}
	}

	log.Printf("Loaded %d test users", len(users))
}

// LoadLoans carga préstamos de prueba en la base de datos
func LoadLoans(DB *gorm.DB) {
	if DB == nil {
		log.Println("Warning: Database connection is nil")
		return
	}

	loans := []models.Loan{
		{
			UserID:         1,
			AmountApproved: decimal.NewFromFloat(10000000),
			Status:         "pending",
		},
		{
			UserID:         2,
			AmountApproved: decimal.NewFromFloat(5000000),
			Status:         "approved",
		},
		{
			UserID:         3,
			AmountApproved: decimal.NewFromFloat(15000000),
			Status:         "rejected",
		},
		{
			UserID:         4,
			AmountApproved: decimal.NewFromFloat(8000000),
			Status:         "approved",
		},
		{
			UserID:         5,
			AmountApproved: decimal.NewFromFloat(3000000),
			Status:         "pending",
		},
		{
			UserID:         1,
			AmountApproved: decimal.NewFromFloat(12000000),
			Status:         "pending",
		},
	}

	for _, loan := range loans {
		if err := DB.Create(&loan).Error; err != nil {
			log.Printf("Error creating loan for user %d: %v", loan.UserID, err)
		}
	}

	log.Printf("Loaded %d test loans", len(loans))
}

// LoadTestData carga todos los datos de prueba en orden correcto
func LoadTestData(DB *gorm.DB) {
	if DB == nil {
		log.Println("Warning: Database connection is nil")
		return
	}

	log.Println("Loading test data...")

	// Limpiar base de datos primero
	ClearDatabase(DB)

	// Cargar datos en orden correcto (usuarios primero, luego préstamos)
	LoadUsers(DB)
	LoadLoans(DB)

	log.Println("Test data loaded successfully")
}

// GetTestUser retorna un usuario de prueba específico por ID
func GetTestUser(DB *gorm.DB, userID uint) (*models.User, error) {
	var user models.User
	err := DB.First(&user, userID).Error
	return &user, err
}

// GetTestLoan retorna un préstamo de prueba específico por ID
func GetTestLoan(DB *gorm.DB, loanID uint) (*models.Loan, error) {
	var loan models.Loan
	err := DB.First(&loan, loanID).Error
	return &loan, err
}

// CountUsers cuenta el número total de usuarios en la base de datos
func CountUsers(DB *gorm.DB) int64 {
	var count int64
	DB.Model(&models.User{}).Count(&count)
	return count
}

// CountLoans cuenta el número total de préstamos en la base de datos
func CountLoans(DB *gorm.DB) int64 {
	var count int64
	DB.Model(&models.Loan{}).Count(&count)
	return count
}

// CountLoansByStatus cuenta préstamos por estado
func CountLoansByStatus(DB *gorm.DB, status string) int64 {
	var count int64
	DB.Model(&models.Loan{}).Where("status = ?", status).Count(&count)
	return count
}

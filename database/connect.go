package database

import (
	"fmt"
	"log"
	"time"

	"loan-api/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// ConnectDB establece la conexión con la base de datos MySQL
func ConnectDB(cfg *config.Config) error {
	var err error

	// Configurar el logger de GORM
	logLevel := logger.Silent
	if cfg.IsDevelopment() {
		logLevel = logger.Info
	}

	// Configuración de GORM
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	}

	// Conectar a la base de datos
	DB, err = gorm.Open(mysql.Open(cfg.GetDSN()), gormConfig)
	if err != nil {
		return fmt.Errorf("error al conectar con la base de datos: %w", err)
	}

	// Obtener la conexión SQL subyacente para configurar el pool
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("error al obtener la conexión SQL: %w", err)
	}

	// Configurar el pool de conexiones
	sqlDB.SetMaxIdleConns(10)           // Conexiones inactivas
	sqlDB.SetMaxOpenConns(30)           // Conexiones máximas abiertas
	sqlDB.SetConnMaxLifetime(time.Hour) // Tiempo de vida de las conexiones

	// Verificar la conexión
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("error al hacer ping a la base de datos: %w", err)
	}

	log.Println("✅ Conexión a la base de datos establecida exitosamente")
	return nil
}

// GetDB retorna la instancia de la base de datos
func GetDB() *gorm.DB {
	return DB
}

// CloseDB cierra la conexión con la base de datos
func CloseDB() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

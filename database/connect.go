package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"loan-api/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Connect inicializa la conexión a la base de datos con configuración dinámica
func Connect(cfg *config.Config) {
	if DB != nil {
		return
	}

	dsn := buildDSN(cfg)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(getLogLevel(cfg)),
	})
	if err != nil {
		log.Fatal("❌ Error al conectar a la base de datos:", err)
	}

	DB = db

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("❌ Error al obtener DB SQL:", err)
	}

	configureConnectionPool(sqlDB, cfg)

	if err := sqlDB.Ping(); err != nil {
		log.Fatal("❌ Error al hacer ping a la base de datos:", err)
	}

	log.Println("✅ Conexión a la base de datos establecida")

	if err := Migrate(); err != nil {
		log.Fatal("❌ Error al ejecutar migraciones:", err)
	}
}

// buildDSN genera el Data Source Name para MySQL
func buildDSN(cfg *config.Config) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)
}

// getLogLevel retorna el nivel de logging según el entorno
func getLogLevel(cfg *config.Config) logger.LogLevel {
	if cfg.IsDevelopment() {
		return logger.Info
	}
	return logger.Silent
}

// configureConnectionPool configura los parámetros del pool de conexiones
func configureConnectionPool(sqlDB *sql.DB, cfg *config.Config) {
	if cfg.IsDevelopment() {
		sqlDB.SetMaxIdleConns(5)
		sqlDB.SetMaxOpenConns(10)
		sqlDB.SetConnMaxLifetime(30 * time.Minute)
	} else {
		sqlDB.SetMaxIdleConns(20)
		sqlDB.SetMaxOpenConns(50)
		sqlDB.SetConnMaxLifetime(2 * time.Hour)
	}
}

// GetDB retorna la instancia GORM de la base de datos
func GetDB() *gorm.DB {
	return DB
}

// CloseDB cierra la conexión activa
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

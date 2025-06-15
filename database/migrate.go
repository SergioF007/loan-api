package database

import (
	"log"

	"loan-api/models"
)

// AutoMigrate ejecuta las migraciones automáticas de GORM
func AutoMigrate() error {
	log.Println("🔄 Iniciando migraciones de base de datos...")

	// Ejecutar migraciones para todos los modelos
	err := DB.AutoMigrate(
		&models.User{},
		&models.Loan{},
	)

	if err != nil {
		log.Printf("❌ Error en las migraciones: %v", err)
		return err
	}

	log.Println("✅ Migraciones completadas exitosamente")
	return nil
}

// DropTables elimina todas las tablas (útil para desarrollo)
// ⚠️ CUIDADO: Solo usar en desarrollo
func DropTables() error {
	log.Println("⚠️  Eliminando todas las tablas...")

	err := DB.Migrator().DropTable(
		&models.Loan{},
		&models.User{},
	)

	if err != nil {
		log.Printf("❌ Error al eliminar tablas: %v", err)
		return err
	}

	log.Println("✅ Tablas eliminadas exitosamente")
	return nil
}

// ResetDatabase elimina y vuelve a crear todas las tablas
// ⚠️ CUIDADO: Solo usar en desarrollo
func ResetDatabase() error {
	log.Println("🔄 Reiniciando base de datos...")

	// Eliminar tablas
	if err := DropTables(); err != nil {
		return err
	}

	// Crear tablas nuevamente
	if err := AutoMigrate(); err != nil {
		return err
	}

	log.Println("✅ Base de datos reiniciada exitosamente")
	return nil
}

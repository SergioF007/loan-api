-- Script de inicialización para la base de datos loan_api
-- Ejecutar este script en MySQL Workbench o desde la línea de comandos

-- Crear la base de datos si no existe
CREATE DATABASE IF NOT EXISTS loan_api 
CHARACTER SET utf8mb4 
COLLATE utf8mb4_unicode_ci;

-- Usar la base de datos
USE loan_api;

-- Verificar que la base de datos fue creada correctamente
SELECT 'Base de datos loan_api creada correctamente' AS mensaje;

-- Las tablas se crearán automáticamente cuando ejecutes la aplicación
-- gracias a la función AutoMigrate de GORM 
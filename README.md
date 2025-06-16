# Loan API - Sistema de Solicitudes de Préstamo

API REST para el manejo completo del flujo de solicitudes de préstamo, desde la inicialización hasta la decisión final y desembolso.

## 📋 Descripción del Proyecto

Sistema que implementa un flujo completo de solicitud de préstamo con las siguientes funcionalidades:

1. **Inicialización de solicitud** - Crear solicitud con datos básicos
2. **Consulta de score crediticio** - Simulación de consulta a servicios externos
3. **Verificación de identidad** - Validación contra registros internos
4. **Decisión final y desembolso** - Evaluación y procesamiento de aprobación/rechazo

## 🏗️ Arquitectura

### Estructura del Proyecto
```
loan-api/
├── app/                    # Configuración principal de la aplicación
├── config/                 # Configuración de la aplicación
├── controllers/            # Controladores HTTP
├── database/              # Migraciones y configuración de BD
├── docs/                  # Documentación Swagger
├── middlewares/           # Middlewares de autenticación y validación
├── models/                # Modelos de datos
├── repositories/          # Capa de acceso a datos
├── routers/               # Configuración de rutas
├── services/              # Lógica de negocio
├── test/                  # Utilidades para testing
├── utils/                 # Utilidades generales
├── app.env                # Variables de entorno
├── go.mod                 # Dependencias Go
└── main.go               # Punto de entrada
```

### Tecnologías Utilizadas
- **Go 1.23+** - Lenguaje de programación
- **Gin** - Framework web
- **GORM** - ORM para base de datos
- **MySQL** - Base de datos
- **JWT** - Autenticación
- **Swagger** - Documentación de API
- **Testify** - Testing

## 🚀 Instalación y Configuración

### Prerrequisitos
- Go 1.23 o superior
- MySQL 8.0 o superior
- Git

### 1. Clonar el repositorio
```bash
git clone https://github.com/SergioF007/loan-api.git
cd loan-api
```

### 2. Instalar dependencias
```bash
go mod download
```

### 3. Configurar base de datos
Crear una base de datos MySQL:
```sql
CREATE DATABASE loan_api;
```

### 4. Configurar variables de entorno
Copiar y configurar el archivo `app.env`:
```bash
cp app.env.example app.env
```

Configurar las siguientes variables en `app.env`:
```env
# Database
DB_HOST=localhost
DB_PORT=3306
DB_USER=tu_usuario
DB_PASSWORD=tu_password
DB_NAME=loan_api

# Server
SERVER_ADDRESS=0.0.0.0:8080
GIN_MODE=debug

# JWT
ACCESS_TOKEN_PRIVATE_KEY=tu_clave_privada_rsa
ACCESS_TOKEN_PUBLIC_KEY=tu_clave_publica_rsa
ACCESS_TOKEN_EXPIRED_IN=15m
ACCESS_TOKEN_MAXAGE=900

# CORS
CLIENT_ORIGIN=http://localhost:3000
```

### 5. Generar claves RSA para JWT
```bash
# Generar clave privada
openssl genrsa -out private_key.pem 2048

# Generar clave pública
openssl rsa -in private_key.pem -pubout -out public_key.pem

# Convertir a formato de una línea para el .env
awk 'NF {sub(/\r/, ""); printf "%s\\n",$0;}' private_key.pem
awk 'NF {sub(/\r/, ""); printf "%s\\n",$0;}' public_key.pem
```

## 🏃‍♂️ Ejecución

### Desarrollo
```bash
# Ejecutar con recarga automática
go run main.go

# O usando air (si está instalado)
air
```

### Producción
```bash
# Compilar
go build -o loan-api main.go

# Ejecutar
./loan-api
```

La API estará disponible en: `http://localhost:8080`

## 📚 Documentación de la API

### Swagger UI
Una vez ejecutando la aplicación, la documentación interactiva estará disponible en:
- **Swagger UI**: `http://localhost:8080/swagger/index.html`

### Endpoints Principales

#### Autenticación
- `POST /api/v1/auth/register` - Registro de usuario
- `POST /api/v1/auth/login` - Inicio de sesión

#### Préstamos
- `POST /api/v1/loans` - Crear solicitud de préstamo
- `POST /api/v1/loans/data` - Guardar datos del préstamo
- `POST /api/v1/loans/{id}/decision` - Procesar decisión final
- `GET /api/v1/loans/{id}` - Obtener préstamo por ID
- `GET /api/v1/loans/user` - Obtener préstamos del usuario

#### Tipos de Préstamo
- `GET /api/v1/loan-types` - Listar tipos de préstamo disponibles

## 🧪 Pruebas

### Ejecutar todas las pruebas
```bash
go test ./...
```

### Ejecutar pruebas con cobertura
```bash
go test -cover ./...
```

### Ejecutar pruebas específicas
```bash
# Pruebas de controladores
go test -v ./controllers

# Pruebas de servicios
go test -v ./services

# Prueba específica
go test -v ./controllers -run TestLoanController_CreateLoan
```

### Pruebas de integración
```bash
# Ejecutar con base de datos de prueba
go test -v ./controllers -run TestLoanController
```

## 🔄 Flujo de Uso

### 1. Registro y Autenticación
```bash
# Registrar usuario
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: 1" \
  -d '{
    "name": "Juan Pérez",
    "email": "juan@example.com",
    "phone": "3001234567",
    "document_type": "cedula",
    "document_number": "12345678",
    "password": "password123!",
    "password_confirmation": "password123!"
  }'

# Iniciar sesión
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: 1" \
  -d '{
    "email": "juan@example.com",
    "password": "password123!"
  }'
```

### 2. Crear Solicitud de Préstamo
```bash
curl -X POST http://localhost:8080/api/v1/loans \
  -H "Content-Type: application/json" \
  -H "Authorization: tu_jwt_token" \
  -H "X-Tenant-ID: 1" \
  -d '{
    "loan_type_id": 1
  }'
```

### 3. Guardar Datos del Préstamo
```bash
curl -X POST http://localhost:8080/api/v1/loans/data \
  -H "Content-Type: application/json" \
  -H "Authorization: tu_jwt_token" \
  -H "X-Tenant-ID: 1" \
  -d '{
    "loan_id": 1,
    "data": [
      {"form_id": 1, "key": "full_name", "value": "Juan Pérez", "index": 0},
      {"form_id": 1, "key": "document_type", "value": "cedula", "index": 0},
      {"form_id": 1, "key": "document_number", "value": "12345678", "index": 0},
      {"form_id": 1, "key": "age", "value": "30", "index": 0},
      {"form_id": 2, "key": "monthly_income", "value": "5000000", "index": 0},
      {"form_id": 2, "key": "monthly_expenses", "value": "2000000", "index": 0},
      {"form_id": 3, "key": "requested_amount", "value": "2000000", "index": 0},
      {"form_id": 3, "key": "purpose", "value": "Educación", "index": 0}
    ]
  }'
```

### 4. Procesar Decisión Final
```bash
curl -X POST http://localhost:8080/api/v1/loans/1/decision \
  -H "Authorization: tu_jwt_token" \
  -H "X-Tenant-ID: 1"
```

## 📊 Estados del Préstamo

| Estado | Descripción |
|--------|-------------|
| `pending` | Solicitud creada, sin datos |
| `on_progress` | Datos parciales guardados |
| `completed` | Datos completos + validaciones realizadas |
| `approved` | Préstamo aprobado y desembolsado |
| `rejected` | Préstamo rechazado |

## 🔧 Configuración Avanzada

### Variables de Entorno Completas
```env
# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=loan_api

# Server Configuration
SERVER_ADDRESS=0.0.0.0:8080
GIN_MODE=debug

# JWT Configuration
ACCESS_TOKEN_PRIVATE_KEY=-----BEGIN RSA PRIVATE KEY-----\n...\n-----END RSA PRIVATE KEY-----
ACCESS_TOKEN_PUBLIC_KEY=-----BEGIN PUBLIC KEY-----\n...\n-----END PUBLIC KEY-----
ACCESS_TOKEN_EXPIRED_IN=15m
ACCESS_TOKEN_MAXAGE=900

# CORS Configuration
CLIENT_ORIGIN=http://localhost:3000
```

### Configuración de Base de Datos
La aplicación incluye migraciones automáticas y datos semilla. Al ejecutar por primera vez:
1. Se crean todas las tablas necesarias
2. Se insertan datos iniciales (tenants, loan types, formularios)
3. Se configuran usuarios de prueba

## 🐛 Solución de Problemas

### Error de conexión a base de datos
```bash
# Verificar que MySQL esté ejecutándose
sudo systemctl status mysql

# Verificar credenciales en app.env
# Verificar que la base de datos existe
```

### Error de claves JWT
```bash
# Regenerar claves RSA
openssl genrsa -out private_key.pem 2048
openssl rsa -in private_key.pem -pubout -out public_key.pem
```

### Problemas con migraciones
```bash
# Limpiar y recrear base de datos
DROP DATABASE loan_api;
CREATE DATABASE loan_api;
# Reiniciar aplicación
```

## 🤝 Contribución

1. Fork el proyecto
2. Crear rama para feature (`git checkout -b feature/nueva-funcionalidad`)
3. Commit cambios (`git commit -am 'Agregar nueva funcionalidad'`)
4. Push a la rama (`git push origin feature/nueva-funcionalidad`)
5. Crear Pull Request

## 📄 Licencia

Este proyecto está bajo la Licencia MIT. Ver el archivo `LICENSE` para más detalles.

## 👥 Autor

**Sergio F007** - [GitHub](https://github.com/SergioF007)

## 🔗 Enlaces Útiles

- [Documentación de Gin](https://gin-gonic.com/)
- [Documentación de GORM](https://gorm.io/)
- [Documentación de JWT](https://jwt.io/)
- [Swagger/OpenAPI](https://swagger.io/) 
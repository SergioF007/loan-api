# Loan API

API REST para gestión de solicitudes de préstamos desarrollada en Go con arquitectura de capas.

## 🏗️ Arquitectura

Esta API sigue el patrón de arquitectura de capas (Clean Architecture) con las siguientes capas:

- **Capa de Presentación**: Controladores y routers (HTTP handlers)
- **Capa de Aplicación**: Servicios (lógica de negocio)
- **Capa de Dominio**: Modelos y entidades
- **Capa de Infraestructura**: Repositorios y base de datos

## 🚀 Características

- ✅ API REST completa para usuarios y préstamos
- ✅ Autenticación JWT
- ✅ Validaciones de datos
- ✅ Paginación
- ✅ Manejo de errores centralizado
- ✅ Logging estructurado
- ✅ Documentación Swagger
- ✅ Configuración con variables de entorno
- ✅ Docker support
- ✅ Middleware CORS
- ✅ Base de datos MySQL con GORM

## 📋 Requisitos

- Go 1.21+
- MySQL 8.0+
- Docker (opcional)

### Configuración de Base de Datos

**Para desarrolladores que clonen el proyecto:**

1. **MySQL debe estar instalado y ejecutándose localmente**
2. **Tener credenciales de acceso a MySQL (usuario y contraseña)**
3. **Verificar la conectividad:**
   ```bash
   mysql -u tu_usuario -p
   ```
4. **El proyecto incluye un script de inicialización:** `database/init.sql`

## 🛠️ Instalación

### Desarrollo Local

1. **Clonar el repositorio**
```bash
git clone <repository-url>
cd loan-api
```

2. **Instalar dependencias**
```bash
go mod download
```

3. **Configurar variables de entorno**
```bash
cp app.env.example app.env
# Editar app.env con tus configuraciones
```

4. **Configurar base de datos MySQL**

Tienes dos opciones para configurar la base de datos:

**Opción A: Usando MySQL Workbench (Recomendado)**
1. Abre MySQL Workbench
2. Conéctate a tu servidor MySQL local
3. Ejecuta el script `database/init.sql` que se encuentra en el proyecto
4. Verifica que la base de datos `loan_api` fue creada correctamente

**Opción B: Desde línea de comandos**
```bash
# Conectar a MySQL
mysql -u vedana -p

# Ejecutar el script de inicialización
source database/init.sql;
```

**Configuración del archivo app.env:**
Asegúrate de configurar correctamente las credenciales de tu base de datos local en el archivo `app.env`:

```env
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=tu_usuario_mysql
DB_PASSWORD=tu_contraseña_mysql
DB_NAME=loan_api
```

**Nota:** Las tablas se crearán automáticamente cuando ejecutes la aplicación por primera vez gracias a GORM AutoMigrate.

5. **Ejecutar la aplicación**
```bash
go run main.go
```

### Con Docker

1. **Construir imagen**
```bash
docker build -t loan-api .
```

2. **Ejecutar contenedor**
```bash
docker run -p 8080:8080 --env-file app.env loan-api
```

## ⚙️ Configuración

Las variables de entorno disponibles son:

| Variable | Descripción | Valor por defecto |
|----------|-------------|-------------------|
| `DB_HOST` | Host de la base de datos | localhost |
| `DB_PORT` | Puerto de la base de datos | 3306 |
| `DB_USER` | Usuario de la base de datos | root |
| `DB_PASSWORD` | Contraseña de la base de datos | password |
| `DB_NAME` | Nombre de la base de datos | loan_api |
| `SERVER_HOST` | Host del servidor | localhost |
| `SERVER_PORT` | Puerto del servidor | 8080 |
| `JWT_SECRET` | Clave secreta para JWT | (requerido) |
| `JWT_EXPIRATION_HOURS` | Horas de expiración del JWT | 24 |
| `APP_ENV` | Ambiente de la aplicación | development |
| `APP_NAME` | Nombre de la aplicación | Loan API |
| `APP_VERSION` | Versión de la aplicación | 1.0.0 |

## 📚 API Endpoints

### Usuarios

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| POST | `/api/v1/users` | Crear usuario |
| GET | `/api/v1/users` | Listar usuarios |
| GET | `/api/v1/users/:id` | Obtener usuario |
| PUT | `/api/v1/users/:id` | Actualizar usuario |
| DELETE | `/api/v1/users/:id` | Eliminar usuario |
| GET | `/api/v1/users/search?email=` | Buscar por email |
| GET | `/api/v1/users/:id/credit-summary` | Resumen crediticio |

### Préstamos

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| POST | `/api/v1/loans` | Crear solicitud de préstamo |
| GET | `/api/v1/loans` | Listar préstamos |
| GET | `/api/v1/loans/:id` | Obtener préstamo |
| PUT | `/api/v1/loans/:id/status` | Actualizar estado |
| GET | `/api/v1/loans/user/:userId` | Préstamos de usuario |
| GET | `/api/v1/loans/status?status=` | Filtrar por estado |
| POST | `/api/v1/loans/:id/process` | Procesar aprobación |
| GET | `/api/v1/loans/statistics` | Estadísticas |

### Utilidad

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| GET | `/health` | Health check |
| GET | `/` | Información de la API |

## 📝 Ejemplos de Uso

### Crear Usuario

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Juan Pérez",
    "email": "juan@example.com",
    "phone": "1234567890",
    "income": 5000.00,
    "credit_score": 750
  }'
```

### Crear Solicitud de Préstamo

```bash
curl -X POST http://localhost:8080/api/v1/loans \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "amount": 10000.00,
    "purpose": "Compra de vehículo para trabajo"
  }'
```

### Actualizar Estado de Préstamo

```bash
curl -X PUT http://localhost:8080/api/v1/loans/1/status \
  -H "Content-Type: application/json" \
  -d '{
    "status": "approved"
  }'
```

## 🔒 Autenticación

La API utiliza JWT para autenticación. Para endpoints protegidos, incluye el token en el header:

```bash
Authorization: Bearer <tu-jwt-token>
```

## 📖 Documentación

### Generar Documentación Swagger

Para generar la documentación Swagger automáticamente:

```bash
# Instalar swag (si no está instalado)
go install github.com/swaggo/swag/cmd/swag@latest

# Generar documentación
swag init
```

La documentación Swagger estará disponible en:
- Desarrollo: `http://localhost:8080/swagger/index.html`

La documentación se genera automáticamente desde las anotaciones de comentarios en los controladores.

## 🧪 Testing

```bash
# Ejecutar tests
go test ./...

# Ejecutar tests con cobertura
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 🐛 Debug

Para habilitar logs detallados, configura:
```
APP_ENV=development
```

## 📊 Monitoreo

### Health Check

```bash
curl http://localhost:8080/health
```

Respuesta:
```json
{
  "status": "OK",
  "message": "Loan API funcionando correctamente",
  "version": "1.0.0"
}
```

## 🏗️ Estructura del Proyecto

```
loan-api/
├── main.go                 # Punto de entrada
├── go.mod                  # Dependencias
├── go.sum                  # Checksums de dependencias
├── app.env                 # Variables de entorno
├── app.env.example         # Ejemplo de configuración
├── Dockerfile              # Configuración Docker
├── README.md               # Este archivo
├── .gitignore              # Archivos ignorados por Git
├── config/                 # Configuración
│   └── config.go
├── database/               # Base de datos
│   ├── connect.go
│   └── migrate.go
├── models/                 # Modelos de datos
│   ├── loan.go
│   └── user.go
├── repositories/           # Capa de datos
│   ├── loan_repository.go
│   └── user_repository.go
├── services/               # Lógica de negocio
│   ├── loan_service.go
│   └── user_service.go
├── controllers/            # Controladores HTTP
│   ├── loan_controller.go
│   └── user_controller.go
├── routers/                # Configuración de rutas
│   ├── loan_router.go
│   └── user_router.go
├── middlewares/            # Middlewares
│   └── auth_middleware.go
├── utils/                  # Utilidades
│   └── response.go
├── app_error/              # Manejo de errores
│   └── errors.go
└── docs/                   # Documentación
    └── swagger.json
```

## 🔄 Estados de Préstamos

- `pending`: Solicitud pendiente de revisión
- `approved`: Préstamo aprobado
- `rejected`: Préstamo rechazado

## 📈 Criterios de Aprobación Automática

El sistema evalúa automáticamente las solicitudes basándose en:

1. **Puntaje Crediticio (40%)**
   - ≥750: 40 puntos
   - ≥700: 30 puntos
   - ≥650: 20 puntos
   - ≥600: 10 puntos

2. **Relación Préstamo/Ingresos (30%)**
   - ≤2x: 30 puntos
   - ≤3x: 20 puntos
   - ≤4x: 10 puntos

3. **Ingresos Estables (20%)**
   - ≥$5,000: 20 puntos
   - ≥$3,000: 15 puntos
   - ≥$2,000: 10 puntos

4. **Monto del Préstamo (10%)**
   - ≤$10,000: 10 puntos
   - ≤$50,000: 5 puntos

**Decisión:**
- ≥80 puntos: Aprobado automáticamente
- <50 puntos: Rechazado automáticamente
- 50-79 puntos: Pendiente para revisión manual

## 🤝 Contribución

1. Fork el proyecto
2. Crea una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## 📄 Licencia

Este proyecto está bajo la Licencia MIT - ver el archivo [LICENSE](LICENSE) para detalles.

## 👨‍💻 Autor

Desarrollado como prueba técnica para demostrar arquitectura de microservicios en Go.

## 🆘 Soporte

Si encuentras algún problema o tienes preguntas:

1. Revisa la documentación
2. Consulta los logs de la aplicación
3. Verifica la configuración de variables de entorno
4. Asegúrate de que la base de datos esté ejecutándose

---

**¡Gracias por usar Loan API! 🚀** 
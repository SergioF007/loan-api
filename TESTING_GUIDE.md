# 🧪 Guía de Pruebas - Loan API

Esta guía te llevará paso a paso a través del flujo completo de la API de préstamos, desde la creación de un usuario hasta la decisión final del préstamo.

## 📋 Prerrequisitos

1. **Aplicación ejecutándose**: `http://localhost:8080`
2. **Base de datos configurada** con datos semilla
3. **Herramienta para hacer requests HTTP**:
   - **Postman** (recomendado)
   - **cURL** (línea de comandos)
   - **Thunder Client** (VS Code)
   - **Insomnia**

## 🔄 Flujo Completo de Pruebas

### **Paso 1: Registro de Usuario**

Primero necesitamos crear un usuario en el sistema.

#### **Endpoint:** `POST /api/v1/auth/register`

**Headers:**
```
Content-Type: application/json
X-Tenant-ID: 1
```

**Body (JSON):**
```json
{
  "name": "Juan Pérez García",
  "email": "juan.perez@example.com",
  "phone": "3001234567",
  "document_type": "cedula",
  "document_number": "12345678",
  "password": "MiPassword123!",
  "password_confirmation": "MiPassword123!"
}
```

**Respuesta Esperada (201):**
```json
{
  "message": "Usuario registrado exitosamente",
  "data": {
    "id": 6,
    "name": "Juan Pérez García",
    "email": "juan.perez@example.com",
    "phone": "3001234567",
    "document_type": "cedula",
    "document_number": "12345678",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

#### **cURL:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: 1" \
  -d '{
    "name": "Juan Pérez García",
    "email": "juan.perez@example.com",
    "phone": "3001234567",
    "document_type": "cedula",
    "document_number": "12345678",
    "password": "MiPassword123!",
    "password_confirmation": "MiPassword123!"
  }'
```

---

### **Paso 2: Inicio de Sesión**

Ahora iniciamos sesión para obtener el token JWT.

#### **Endpoint:** `POST /api/v1/auth/login`

**Headers:**
```
Content-Type: application/json
X-Tenant-ID: 1
```

**Body (JSON):**
```json
{
  "email": "juan.perez@example.com",
  "password": "MiPassword123!"
}
```

**Respuesta Esperada (200):**
```json
{
  "message": "Inicio de sesión exitoso",
  "data": {
    "token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 6,
      "name": "Juan Pérez García",
      "email": "juan.perez@example.com",
      "phone": "3001234567",
      "document_type": "cedula",
      "document_number": "12345678",
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    }
  }
}
```

#### **cURL:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: 1" \
  -d '{
    "email": "juan.perez@example.com",
    "password": "MiPassword123!"
  }'
```

**⚠️ IMPORTANTE:** Guarda el token JWT de la respuesta, lo necesitarás para los siguientes pasos.

---

### **Paso 3: Crear Solicitud de Préstamo**

Creamos una nueva solicitud de préstamo.

#### **Endpoint:** `POST /api/v1/loans`

**Headers:**
```
Content-Type: application/json
Authorization: eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...
X-Tenant-ID: 1
```

**Body (JSON):**
```json
{
  "loan_type_id": 1
}
```

**Respuesta Esperada (201):**
```json
{
  "message": "Solicitud de préstamo creada exitosamente",
  "data": {
    "id": 7,
    "loan_type_id": 1,
    "loan_type": {
      "id": 1,
      "name": "Préstamo Personal",
      "code": "personal_loan",
      "description": "Préstamo personal para gastos diversos",
      "min_amount": 100000,
      "max_amount": 10000000
    },
    "user_id": 6,
    "user": {
      "id": 6,
      "name": "Juan Pérez García",
      "email": "juan.perez@example.com",
      "phone": "3001234567",
      "document_type": "cedula",
      "document_number": "12345678",
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    },
    "status": "pending",
    "observation": "",
    "amount_approved": "0.00",
    "data": [],
    "created_at": "2024-01-15T10:35:00Z",
    "updated_at": "2024-01-15T10:35:00Z"
  }
}
```

#### **cURL:**
```bash
curl -X POST http://localhost:8080/api/v1/loans \
  -H "Content-Type: application/json" \
  -H "Authorization: TU_TOKEN_JWT_AQUI" \
  -H "X-Tenant-ID: 1" \
  -d '{
    "loan_type_id": 1
  }'
```

**⚠️ IMPORTANTE:** Guarda el `id` del préstamo (en este ejemplo es `7`), lo necesitarás para los siguientes pasos.

---

### **Paso 4: Guardar Datos del Préstamo**

Ahora guardamos todos los datos necesarios del préstamo. Este paso incluye la consulta automática del score crediticio y la verificación de identidad.

#### **Endpoint:** `POST /api/v1/loans/data`

**Headers:**
```
Content-Type: application/json
Authorization: eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...
X-Tenant-ID: 1
```

**Body (JSON):**
```json
{
  "loan_id": 7,
  "data": [
    {
      "form_id": 1,
      "key": "full_name",
      "value": "Juan Pérez García",
      "index": 0
    },
    {
      "form_id": 1,
      "key": "document_type",
      "value": "cedula",
      "index": 0
    },
    {
      "form_id": 1,
      "key": "document_number",
      "value": "12345678",
      "index": 0
    },
    {
      "form_id": 1,
      "key": "age",
      "value": "30",
      "index": 0
    },
    {
      "form_id": 2,
      "key": "monthly_income",
      "value": "5000000",
      "index": 0
    },
    {
      "form_id": 2,
      "key": "monthly_expenses",
      "value": "2000000",
      "index": 0
    },
    {
      "form_id": 3,
      "key": "requested_amount",
      "value": "2500000",
      "index": 0
    },
    {
      "form_id": 3,
      "key": "purpose",
      "value": "Educación",
      "index": 0
    }
  ]
}
```

**Respuesta Esperada (200):**
```json
{
  "message": "Datos del préstamo guardados exitosamente"
}
```

#### **cURL:**
```bash
curl -X POST http://localhost:8080/api/v1/loans/data \
  -H "Content-Type: application/json" \
  -H "Authorization: TU_TOKEN_JWT_AQUI" \
  -H "X-Tenant-ID: 1" \
  -d '{
    "loan_id": 7,
    "data": [
      {"form_id": 1, "key": "full_name", "value": "Juan Pérez García", "index": 0},
      {"form_id": 1, "key": "document_type", "value": "cedula", "index": 0},
      {"form_id": 1, "key": "document_number", "value": "12345678", "index": 0},
      {"form_id": 1, "key": "age", "value": "30", "index": 0},
      {"form_id": 2, "key": "monthly_income", "value": "5000000", "index": 0},
      {"form_id": 2, "key": "monthly_expenses", "value": "2000000", "index": 0},
      {"form_id": 3, "key": "requested_amount", "value": "2500000", "index": 0},
      {"form_id": 3, "key": "purpose", "value": "Educación", "index": 0}
    ]
  }'
```

**🔍 ¿Qué sucede internamente?**
1. Se guardan todos los datos del formulario
2. Se simula la consulta del **score crediticio** (basado en document_number)
3. Se realiza la **verificación de identidad** (comparando con datos registrados)
4. El estado del préstamo cambia a `completed` si todos los campos requeridos están completos

---

### **Paso 5: Verificar Estado del Préstamo**

Antes de procesar la decisión final, verifiquemos que el préstamo esté en estado `completed`.

#### **Endpoint:** `GET /api/v1/loans/{id}`

**Headers:**
```
Authorization: eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...
X-Tenant-ID: 1
```

**Respuesta Esperada (200):**
```json
{
  "message": "Préstamo obtenido exitosamente",
  "data": {
    "id": 7,
    "loan_type_id": 1,
    "loan_type": {
      "id": 1,
      "name": "Préstamo Personal",
      "code": "personal_loan",
      "description": "Préstamo personal para gastos diversos",
      "min_amount": 100000,
      "max_amount": 10000000
    },
    "user_id": 6,
    "user": {
      "id": 6,
      "name": "Juan Pérez García",
      "email": "juan.perez@example.com",
      "phone": "3001234567",
      "document_type": "cedula",
      "document_number": "12345678",
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    },
    "status": "completed",
    "observation": "Solicitud completada. Score crediticio: 678. Verificación de identidad: exitosa.",
    "amount_approved": "0.00",
    "credit_score": 678,
    "identity_verified": true,
    "data": [
      {
        "id": 1,
        "form_id": 1,
        "key": "full_name",
        "value": "Juan Pérez García",
        "index": 0
      },
      {
        "id": 2,
        "form_id": 1,
        "key": "document_type",
        "value": "cedula",
        "index": 0
      }
    ],
    "created_at": "2024-01-15T10:35:00Z",
    "updated_at": "2024-01-15T10:40:00Z"
  }
}
```

#### **cURL:**
```bash
curl -X GET http://localhost:8080/api/v1/loans/7 \
  -H "Authorization: TU_TOKEN_JWT_AQUI" \
  -H "X-Tenant-ID: 1"
```

**✅ Verificaciones importantes:**
- `status` debe ser `"completed"`
- `credit_score` debe tener un valor (ej: 678)
- `identity_verified` debe ser `true`

---

### **Paso 6: Procesar Decisión Final del Préstamo**

Finalmente, procesamos la decisión final que evaluará el score crediticio y la verificación de identidad para aprobar o rechazar el préstamo.

#### **Endpoint:** `POST /api/v1/loans/{id}/decision`

**Headers:**
```
Authorization: eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...
X-Tenant-ID: 1
```

**Body:** No requiere body

**Respuesta Esperada - Caso Aprobado (200):**
```json
{
  "message": "Decisión del préstamo procesada exitosamente",
  "data": {
    "id": 7,
    "loan_type_id": 1,
    "loan_type": {
      "id": 1,
      "name": "Préstamo Personal",
      "code": "personal_loan",
      "description": "Préstamo personal para gastos diversos",
      "min_amount": 100000,
      "max_amount": 10000000
    },
    "user_id": 6,
    "user": {
      "id": 6,
      "name": "Juan Pérez García",
      "email": "juan.perez@example.com",
      "phone": "3001234567",
      "document_type": "cedula",
      "document_number": "12345678",
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    },
    "status": "approved",
    "observation": "Préstamo aprobado: buen score crediticio (678) - Desembolso realizado exitosamente",
    "amount_approved": "2500000.00",
    "credit_score": 678,
    "identity_verified": true,
    "data": [],
    "created_at": "2024-01-15T10:35:00Z",
    "updated_at": "2024-01-15T10:45:00Z"
  }
}
```

**Respuesta Esperada - Caso Rechazado (200):**
```json
{
  "message": "Decisión del préstamo procesada exitosamente",
  "data": {
    "id": 7,
    "status": "rejected",
    "observation": "Préstamo rechazado: score crediticio muy bajo (350)",
    "amount_approved": "0.00",
    "credit_score": 350,
    "identity_verified": true
  }
}
```

#### **cURL:**
```bash
curl -X POST http://localhost:8080/api/v1/loans/7/decision \
  -H "Authorization: TU_TOKEN_JWT_AQUI" \
  -H "X-Tenant-ID: 1"
```

**🔍 ¿Qué sucede internamente?**
1. Se valida que el préstamo esté en estado `completed`
2. Se evalúan las **reglas de negocio**:
   - Verificación de identidad (obligatoria)
   - Score crediticio mínimo (≥400)
   - Capacidad de pago (máximo 50% del ingreso)
   - Monto mínimo ($100,000)
3. Se calcula el **monto aprobado**
4. Se simula el **desembolso** (puede fallar ocasionalmente)
5. Se actualiza el estado final: `approved` o `rejected`

---

## 🎯 Casos de Prueba Adicionales

### **Caso 1: Préstamo Rechazado por Score Bajo**

Usa un `document_number` que genere un score bajo (ej: "00000001"):

```json
{
  "form_id": 1,
  "key": "document_number",
  "value": "00000001",
  "index": 0
}
```

### **Caso 2: Préstamo Rechazado por Verificación de Identidad**

Usa un nombre diferente al registrado:

```json
{
  "form_id": 1,
  "key": "full_name",
  "value": "Nombre Diferente",
  "index": 0
}
```

### **Caso 3: Préstamo Rechazado por Monto Excesivo**

Solicita un monto mayor al 50% del ingreso mensual:

```json
{
  "form_id": 2,
  "key": "monthly_income",
  "value": "2000000",
  "index": 0
},
{
  "form_id": 3,
  "key": "requested_amount",
  "value": "1500000",
  "index": 0
}
```

### **Caso 4: Fallo en Desembolso**

Usa un `user_id` que termine en 0 para simular fallo de desembolso:
- Registra un usuario con email que genere ID terminado en 0
- El préstamo se aprobará pero fallará el desembolso

---

## 📊 Estados del Préstamo

| Estado | Descripción |
|--------|-------------|
| `pending` | Solicitud creada, sin datos |
| `on_progress` | Datos parciales guardados |
| `completed` | Datos completos + validaciones realizadas |
| `approved` | Préstamo aprobado y desembolsado |
| `rejected` | Préstamo rechazado |

---

## 🔧 Reglas de Negocio Implementadas

### **Score Crediticio:**
- **Generación:** Basada en `document_number` para consistencia
- **Rango:** 300-850 puntos
- **Mínimo aceptable:** 400 puntos

### **Verificación de Identidad:**
- **Documento:** Debe coincidir exactamente
- **Nombre:** El nombre registrado debe estar contenido en el nombre proporcionado

### **Evaluación de Aprobación:**
1. **Identidad verificada** (obligatorio)
2. **Score ≥ 400** (obligatorio)
3. **Monto ≤ 50% ingreso mensual** (si score < 650)
4. **Monto ≥ $100,000** (obligatorio)

### **Simulación de Desembolso:**
- **Éxito:** 90% de los casos
- **Fallo:** 10% de los casos (user_id terminado en 0)
- **Límite:** Máximo $50,000,000

---

## 🐛 Posibles Errores y Soluciones

### **Error 401 - Unauthorized**
```json
{
  "error": true,
  "message": "Token de autenticación requerido"
}
```
**Solución:** Verifica que el token JWT esté en el header `Authorization`

### **Error 403 - Forbidden**
```json
{
  "error": true,
  "message": "Tenant invalido."
}
```
**Solución:** Verifica que el header `X-Tenant-ID` sea `1`

### **Error 500 - Solo préstamos completed**
```json
{
  "success": false,
  "message": "Error en la solicitud",
  "error": {
    "message": "solo se pueden evaluar préstamos en estado completado"
  }
}
```
**Solución:** Asegúrate de completar el Paso 4 antes del Paso 6

---

## 📝 Notas Importantes

1. **Orden de ejecución:** Los pasos deben ejecutarse en orden secuencial
2. **Tokens JWT:** Tienen expiración de 15 minutos
3. **Tenant ID:** Siempre usar `1` para las pruebas
4. **Document numbers:** Diferentes números generan diferentes scores
5. **Nombres:** Deben coincidir parcialmente para verificación de identidad

---

## 🎉 ¡Felicitaciones!

Si has completado todos los pasos exitosamente, has probado el flujo completo de la API de préstamos:

✅ Registro de usuario  
✅ Autenticación JWT  
✅ Creación de solicitud  
✅ Consulta de score crediticio  
✅ Verificación de identidad  
✅ Decisión final y desembolso  

El sistema está funcionando correctamente y cumple con todos los requerimientos de la prueba técnica. 
# Etapa de construcción
FROM golang:1.21-alpine AS builder

# Instalar dependencias del sistema
RUN apk add --no-cache git ca-certificates tzdata

# Establecer directorio de trabajo
WORKDIR /app

# Copiar archivos de dependencias
COPY go.mod go.sum ./

# Descargar dependencias
RUN go mod download

# Copiar código fuente
COPY . .

# Construir la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Etapa de producción
FROM alpine:latest

# Instalar certificados SSL y zona horaria
RUN apk --no-cache add ca-certificates tzdata

# Crear usuario no-root
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Establecer directorio de trabajo
WORKDIR /root/

# Copiar binario desde la etapa de construcción
COPY --from=builder /app/main .

# Copiar archivo de configuración de ejemplo
COPY --from=builder /app/app.env.example ./app.env

# Cambiar propietario de los archivos
RUN chown -R appuser:appgroup /root/

# Cambiar a usuario no-root
USER appuser

# Exponer puerto
EXPOSE 8080

# Comando por defecto
CMD ["./main"] 
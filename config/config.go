package config

import (
	"log"

	"github.com/spf13/viper"
)

// Config contiene toda la configuración de la aplicación
type Config struct {
	// Base de datos
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`

	// Servidor
	ServerHost string `mapstructure:"SERVER_HOST"`
	ServerPort string `mapstructure:"SERVER_PORT"`

	// JWT
	JWTSecret          string `mapstructure:"JWT_SECRET"`
	JWTExpirationHours int    `mapstructure:"JWT_EXPIRATION_HOURS"`

	// Aplicación
	AppEnv     string `mapstructure:"APP_ENV"`
	AppName    string `mapstructure:"APP_NAME"`
	AppVersion string `mapstructure:"APP_VERSION"`
}

// LoadConfig carga la configuración desde variables de entorno
func LoadConfig(path string) (config Config, err error) {
	// Configurar Viper
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	// Permitir que Viper lea variables de entorno
	viper.AutomaticEnv()

	// Leer el archivo de configuración
	err = viper.ReadInConfig()
	if err != nil {
		log.Printf("Warning: No se pudo leer el archivo de configuración: %v", err)
		log.Println("Usando solo variables de entorno...")
	}

	// Mapear las variables a la estructura Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return config, err
	}

	// Validar configuración crítica
	if config.DBHost == "" {
		config.DBHost = "localhost"
	}
	if config.DBPort == "" {
		config.DBPort = "3306"
	}
	if config.ServerPort == "" {
		config.ServerPort = "8080"
	}
	if config.JWTSecret == "" {
		log.Fatal("JWT_SECRET es requerido")
	}
	if config.JWTExpirationHours == 0 {
		config.JWTExpirationHours = 24
	}

	return config, nil
}

// GetDSN retorna la cadena de conexión para MySQL
func (c *Config) GetDSN() string {
	return c.DBUser + ":" + c.DBPassword + "@tcp(" + c.DBHost + ":" + c.DBPort + ")/" + c.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"
}

// GetServerAddress retorna la dirección completa del servidor
func (c *Config) GetServerAddress() string {
	return c.ServerHost + ":" + c.ServerPort
}

// IsDevelopment verifica si la aplicación está en modo desarrollo
func (c *Config) IsDevelopment() bool {
	return c.AppEnv == "development"
}

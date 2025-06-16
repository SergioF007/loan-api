package config

import (
	"fmt"
	"log"
	"time"

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

	// Token Configuration (RSA)
	AccessTokenPrivateKey string        `mapstructure:"ACCESS_TOKEN_PRIVATE_KEY"`
	AccessTokenPublicKey  string        `mapstructure:"ACCESS_TOKEN_PUBLIC_KEY"`
	AccessTokenExpiresIn  time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRED_IN"`
	AccessTokenMaxAge     int           `mapstructure:"ACCESS_TOKEN_MAXAGE"`

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
		log.Printf("Warning: No se pudo leer el archivo de configuración desde %s: %v", path, err)
		log.Println("Usando solo variables de entorno...")
	} else {
		log.Printf("Configuración cargada desde: %s", viper.ConfigFileUsed())
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
	if config.AccessTokenPrivateKey == "" {
		return config, fmt.Errorf("ACCESS_TOKEN_PRIVATE_KEY es requerido")
	}
	if config.AccessTokenPublicKey == "" {
		return config, fmt.Errorf("ACCESS_TOKEN_PUBLIC_KEY es requerido")
	}
	if config.AccessTokenExpiresIn == 0 {
		config.AccessTokenExpiresIn = 600 * time.Minute
	}
	if config.AccessTokenMaxAge == 0 {
		config.AccessTokenMaxAge = 43800
	}

	log.Printf("Configuración cargada: DB=%s:%s/%s, AppEnv=%s", config.DBHost, config.DBPort, config.DBName, config.AppEnv)

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
	return c.AppEnv == "local" || c.AppEnv == "dev"
}

// IsProduction verifica si la aplicación está en modo producción
func (c *Config) IsProduction() bool {
	return c.AppEnv == "prod"
}

// IsLocal verifica si la aplicación está en modo local
func (c *Config) IsLocal() bool {
	return c.AppEnv == "local"
}

// IsQA verifica si la aplicación está en modo QA
func (c *Config) IsQA() bool {
	return c.AppEnv == "qa"
}

// GetEnvironment retorna el entorno actual
func (c *Config) GetEnvironment() string {
	return c.AppEnv
}

// UseRSATokens verifica si debe usar tokens RSA
func (c *Config) UseRSATokens() bool {
	return c.AccessTokenPrivateKey != "" && c.AccessTokenPublicKey != ""
}

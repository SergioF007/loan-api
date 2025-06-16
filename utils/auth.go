package utils

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"time"

	"loan-api/config"
	"loan-api/models"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrKeyMustBePEMEncoded = errors.New("invalid key: key must be a PEM encoded PKCS1 or PKCS8 key")
	ErrNotRSAPrivateKey    = errors.New("key is not a valid RSA private key")
)

// ParseRSAPrivateKeyFromPEM parses a PEM encoded PKCS1 or PKCS8 private key
func ParseRSAPrivateKeyFromPEM(key []byte) (*rsa.PrivateKey, error) {
	var err error

	// Parse PEM block
	var block *pem.Block
	if block, _ = pem.Decode(key); block == nil {
		return nil, ErrKeyMustBePEMEncoded
	}

	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS1PrivateKey(block.Bytes); err != nil {
		if parsedKey, err = x509.ParsePKCS8PrivateKey(block.Bytes); err != nil {
			return nil, err
		}
	}

	var pkey *rsa.PrivateKey
	var ok bool
	if pkey, ok = parsedKey.(*rsa.PrivateKey); !ok {
		return nil, ErrNotRSAPrivateKey
	}

	return pkey, nil
}

// CreateToken genera un token JWT usando RSA
func CreateToken(ttl time.Duration, payload interface{}, privateKey string) (string, error) {
	decodePrivateKey, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return "", fmt.Errorf("could not decode key: %w", err)
	}

	key, err := ParseRSAPrivateKeyFromPEM(decodePrivateKey)
	if err != nil {
		return "", fmt.Errorf("create: parse key: %w", err)
	}

	now := time.Now().UTC()

	claims := make(jwt.MapClaims)
	claims["sub"] = payload
	claims["exp"] = now.Add(ttl).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil {
		return "", fmt.Errorf("create: sign token: %w", err)
	}

	return token, nil
}

// GenerateAccessToken genera un token JWT para acceso
func GenerateAccessToken(user *models.User, cfg *config.Config) (string, error) {
	log.Printf("GenerateAccessToken - Generando token para usuario ID: %d - Email: %s", user.ID, user.Email)

	payload := map[string]interface{}{
		"id":    user.ID,
		"email": user.Email,
	}

	ttl := cfg.AccessTokenExpiresIn
	if ttl == 0 {
		ttl = time.Hour // Valor por defecto: 1 hora
	}

	return CreateToken(ttl, payload, cfg.AccessTokenPrivateKey)
}

// ValidateToken valida un token JWT
func ValidateToken(token string, publicKey string) (interface{}, error) {
	decodePublicKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return nil, fmt.Errorf("could not decode: %w", err)
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(decodePublicKey)
	if err != nil {
		return "", fmt.Errorf("validate: parse key: %w", err)
	}

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
		}

		return key, nil
	})
	if err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, fmt.Errorf("validate: invalid token")
	}

	return claims["sub"], nil
}

// ExtractIPFromForwardedHeader extrae la IP del header Forwarded
func ExtractIPFromForwardedHeader(forwardedHeader string) string {
	if forwardedHeader == "" {
		return ""
	}

	// Esta es una implementación básica
	// En producción, se debería hacer un parsing más robusto
	return forwardedHeader
}

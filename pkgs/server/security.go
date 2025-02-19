package server

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"

	jwt "github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	privateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

type Claims struct {
	UserID  string
	GuildID string

	jwt.RegisteredClaims
}

// If no keypair exists, generate a new one
func NewJWT() *JWTService {
	pkFile, err := os.ReadFile("/app/certs/private.key")
	if err != nil {
		log.Fatal(err)
	}
	KeyBlock, _ := pem.Decode(pkFile)
	privateKey, err := x509.ParsePKCS1PrivateKey(KeyBlock.Bytes)
	if err != nil {
		log.Fatal(err)
	}
	publicKey := &privateKey.PublicKey
	return &JWTService{privateKey: privateKey, PublicKey: publicKey}
}

// Create new JWT token
func (t *JWTService) CreateToken(claims Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = "1"
	tokenString, err := token.SignedString(t.privateKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// Verify JWT token and only return error
func (t *JWTService) VerifyToken(tokenString string) error {
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return t.PublicKey, nil
	})
	return err
}

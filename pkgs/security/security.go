package security

import (
	"crypto/rand"
	"os"

	jwt "github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secretKey []byte
}

type Claims struct {
	UserID  string
	GuildID string
	Admin   bool

	jwt.RegisteredClaims
}

func NewJWT() (*JWTService, error) {
	var secretKey []byte
	var err error

	if os.Getenv("JWT_SECRET_KEY") != "" {
		keyBytes := []byte(os.Getenv("JWT_SECRET_KEY"))
		if err != nil {
			return nil, err
		}
		secretKey = keyBytes
	} else {
		secretKey, err = generateSecretKey()
		if err != nil {
			return nil, err
		}

	}
	return &JWTService{secretKey: secretKey}, nil
}

// Create new JWT token
func (t *JWTService) CreateToken(claims Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(t.secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// Verify JWT token and only return error
func (t *JWTService) VerifyToken(tokenString string) error {
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return t.secretKey, nil
	})
	return err
}

func (t *JWTService) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return t.secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, err
	}
	return claims, nil
}

func generateSecretKey() ([]byte, error) {
	secretKey := make([]byte, 32)
	if _, err := rand.Read(secretKey); err != nil {
		return nil, err
	}
	return secretKey, nil
}

package service

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

// jwt service
type JWTService interface {
	GenerateToken(userId, name, email string) string
	ValidateToken(token string) (*jwt.Token, error)
}

type AuthCustomClaims struct {
	UserId string `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	jwt.StandardClaims
}

type JwtServicesImpl struct {
	secretKey string
	issure    string
}

func getSecretKey() string {
	secret := os.Getenv("TOKEN_SECRET_KEY")
	if secret == "" {
		secret = "secret"
	}
	return secret
}

func (service *JwtServicesImpl) GenerateToken(userId, name, email string) string {
	claims := &AuthCustomClaims{
		userId,
		email,
		name,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
			Issuer:    service.issure,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//encoded string
	t, err := token.SignedString([]byte(service.secretKey))
	if err != nil {
		panic(err)
	}
	return t
}

func (service *JwtServicesImpl) ValidateToken(encodedToken string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(encodedToken, &AuthCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, isValid := token.Method.(*jwt.SigningMethodHMAC); !isValid {
			return nil, errors.New("invalid token")

		}

		return []byte(service.secretKey), nil
	})

}

func NewJWTAuthService() JWTService {
	return &JwtServicesImpl{
		secretKey: getSecretKey(),
		issure:    "Bikash",
	}
}

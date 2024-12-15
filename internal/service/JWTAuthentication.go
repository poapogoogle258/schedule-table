package service

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

// jwt service
type JWTService interface {
	GenerateToken(userId, name, email string) string
	ValidateToken(token string) (*jwt.Token, error)
}
type authCustomClaims struct {
	UserId string `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	jwt.StandardClaims
}

type jwtServices struct {
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

func (service *jwtServices) GenerateToken(userId, name, email string) string {
	claims := &authCustomClaims{
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

func (service *jwtServices) ValidateToken(encodedToken string) (*jwt.Token, error) {
	return jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, isvalid := token.Method.(*jwt.SigningMethodHMAC); !isvalid {
			return nil, fmt.Errorf("Invalid token", token.Header["alg"])
		}
		return []byte(service.secretKey), nil
	})

}

func NewJWTAuthService() JWTService {
	return &jwtServices{
		secretKey: getSecretKey(),
		issure:    "Bikash",
	}
}

package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const hash_cost int = 12

type AuthService struct {
	SigningKey []byte
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), hash_cost)
	if err != nil {
		log.Fatal(err)
	}
	return string(hashedPassword), nil
}

func CheckPasswordHash(hash, password string) error {
	passOk := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if passOk != nil {
		log.Println("Password mismatch")
	}
	return passOk
}

// JWT
func GetBearerToken(headers http.Header) (string, error) {
	authHeaders, found := headers["Authorization"]
	if !found {
		return "", fmt.Errorf("request Header doesn't contain the Bearer token")
	}
	for _, authHeader := range authHeaders {
		parts := strings.Split(authHeader, " ")
		if strings.TrimSpace(parts[0]) == "Bearer" {
			if len(parts) < 2 || strings.TrimSpace(parts[1]) == "" {
				return "", fmt.Errorf("bearer token is empty")
			}
			return strings.TrimSpace(parts[1]), nil
		}
	}
	return "", fmt.Errorf("request Header doesn't contain the Bearer token")
}

func MakeJWT(userID uuid.UUID, tokenSecret string) (string, error) {
	claims := &jwt.RegisteredClaims{
		Issuer:    "chirpy",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 60)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Subject:   userID.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	if !token.Valid {
		return uuid.Nil, fmt.Errorf("invalid token")
	}

	fmt.Println(claims.Subject)
	return uuid.MustParse(claims.Subject), nil
}

func MakeRefreshToken() string {
	refreshToken := make([]byte, 32)
	rand.Read(refreshToken)
	return hex.EncodeToString(refreshToken)
}

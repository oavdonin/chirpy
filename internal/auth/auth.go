package auth

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

const hash_cost int = 12

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

package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	// First, we need to create some hashed passwords for testing
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name     string
		password string
		hash     string
		wantErr  bool
	}{
		{
			name:     "Correct password",
			password: password1,
			hash:     hash1,
			wantErr:  false,
		},
		{
			name:     "Incorrect password",
			password: "wrongPassword",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Password doesn't match different hash",
			password: password1,
			hash:     hash2,
			wantErr:  true,
		},
		{
			name:     "Empty password",
			password: "",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Invalid hash",
			password: password1,
			hash:     "invalidhash",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPasswordHash(tt.hash, tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// JWT Tests
func TestMakeAndValidateJWT(t *testing.T) {
	tokenSecret := "supersecret"

	userID := uuid.New()
	expiresIn := time.Minute * 5

	// Создаём токен
	tokenString, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("failed to create JWT: %v", err)
	}

	// Проверяем валидацию этого токена
	validatedID, err := ValidateJWT(tokenString, tokenSecret)
	if err != nil {
		t.Fatalf("failed to validate JWT: %v", err)
	}

	if validatedID != userID {
		t.Errorf("validated UUID does not match. Expected %v, got %v", userID, validatedID)
	}
}

func TestExpiredJWT(t *testing.T) {
	tokenSecret := "supersecret"

	userID := uuid.New()

	// Токен с истёкшим временем
	tokenString, err := MakeJWT(userID, tokenSecret, -time.Minute)
	if err != nil {
		t.Fatalf("failed to create expired JWT: %v", err)
	}

	// Проверяем, что такой токен невалиден
	_, err = ValidateJWT(tokenString, tokenSecret)
	if err == nil {
		t.Error("expected error when validating expired JWT, got none")
	}
}

func TestInvalidSignatureJWT(t *testing.T) {
	tokenSecret := "supersecret"
	wrongSecret := "wrongsecret"

	userID := uuid.New()
	expiresIn := time.Minute * 5

	// Создаём токен с правильным ключом
	tokenString, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("failed to create JWT: %v", err)
	}

	// Проверяем, что валидация с другим ключом не проходит
	_, err = ValidateJWT(tokenString, wrongSecret)
	if err == nil {
		t.Error("expected error when validating JWT with invalid signature, got none")
	}
}

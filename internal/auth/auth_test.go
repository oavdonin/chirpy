package auth

import (
	"net/http"
	"testing"

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

	tokenString, err := MakeJWT(userID, tokenSecret)
	if err != nil {
		t.Fatalf("failed to create JWT: %v", err)
	}

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

	tokenString, err := MakeJWT(userID, tokenSecret)
	if err != nil {
		t.Fatalf("failed to create expired JWT: %v", err)
	}

	_, err = ValidateJWT(tokenString, tokenSecret)
	if err == nil {
		t.Error("expected error when validating expired JWT, got none")
	}
}

func TestInvalidSignatureJWT(t *testing.T) {
	tokenSecret := "supersecret"
	wrongSecret := "wrongsecret"

	userID := uuid.New()

	tokenString, err := MakeJWT(userID, tokenSecret)
	if err != nil {
		t.Fatalf("failed to create JWT: %v", err)
	}

	_, err = ValidateJWT(tokenString, wrongSecret)
	if err == nil {
		t.Error("expected error when validating JWT with invalid signature, got none")
	}
}

// testing getBearerToken
func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name      string
		headers   http.Header
		wantToken string
		expectErr bool
	}{
		{
			name: "valid bearer token",
			headers: http.Header{
				"Authorization": []string{"Bearer my-token-123"},
			},
			wantToken: "my-token-123",
			expectErr: false,
		},
		{
			name:      "missing authorization header",
			headers:   http.Header{},
			wantToken: "",
			expectErr: true,
		},
		{
			name: "authorization header without bearer",
			headers: http.Header{
				"Authorization": []string{"Basic abcdefgh"},
			},
			wantToken: "",
			expectErr: true,
		},
		{
			name: "bearer without token",
			headers: http.Header{
				"Authorization": []string{"Bearer "},
			},
			wantToken: "",
			expectErr: true,
		},
		{
			name: "multiple authorization headers with bearer",
			headers: http.Header{
				"Authorization": []string{
					"Basic abcdefgh",
					"Bearer multi-token-456",
				},
			},
			wantToken: "multi-token-456",
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GetBearerToken(tt.headers)
			if (err != nil) != tt.expectErr {
				t.Errorf("expected error: %v, got: %v", tt.expectErr, err)
			}
			if token != tt.wantToken {
				t.Errorf("expected token: %s, got: %s", tt.wantToken, token)
			}
		})
	}
}

package auth

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeAndValidateJWT(t *testing.T) {
	// Setup test data
	userID := uuid.New()
	secret := "your-test-secret"

	tests := []struct {
		name     string
		duration time.Duration
		secret   string
		wantErr  bool
	}{
		{
			name:     "valid token",
			duration: time.Hour,
			secret:   secret,
			wantErr:  false,
		},
		{
			name:     "expired token",
			duration: -time.Hour, // negative duration creates already-expired token
			secret:   secret,
			wantErr:  true,
		},
		{
			name:     "wrong secret",
			duration: time.Hour,
			secret:   "wrong-secret",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			str, err := MakeJWT(userID, secret)
			if err != nil {
				t.Fatalf("Error making JWT token")
			}

			id, err := ValidateJWT(str, tt.secret)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if id != userID {
				t.Fatalf("UUIDs don't match. Expected: %v, Got: %v", userID, id)
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {
	headers := http.Header{}
	headers.Add("Authorization", "Bearer MY_TOKEN_VAL")

	tests := []struct {
		header http.Header
		token  string
	}{
		{
			header: headers,
			token:  "MY_TOKEN_VAL",
		},
		{
			header: nil,
			token:  "",
		},
	}

	for _, tt := range tests {
		token, err := GetBearerToken(tt.header)
		if err != nil {
			if !strings.Contains(err.Error(), "Authorization header is not present") {
				t.Fatal("Incorrect error")
			}
		}

		if token != tt.token {
			t.Fatalf("Expected toekn to equal %v.  Actual:%v", tt.token, token)
		}
	}
}

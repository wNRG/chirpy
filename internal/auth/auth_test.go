package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestPasswordHash(t *testing.T) {
	password := "lebronjames23"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	match, err := CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("failed to compare password: %v", err)
	}

	if !match {
		t.Fatal("expected password to match")
	}

	match, err = CheckPasswordHash("lebwrongjames23", hash)
	if err != nil {
		t.Fatalf("failed to compare password: %v", err)
	}

	if match {
		t.Fatal("expected password not to match")
	}
}

func TestJWT(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"

	t.Run("valid token", func(t *testing.T) {
		token, err := MakeJWT(userID, secret, time.Hour)
		if err != nil {
			t.Fatalf("failed to create token: %v", err)
		}

		validatedID, err := ValidateJWT(token, secret)
		if err != nil {
			t.Fatalf("failed to validate token: %v", err)
		}

		if validatedID != userID {
			t.Errorf("expected user ID %s, got %s", userID, validatedID)
		}
	})

	t.Run("exipred token", func(t *testing.T) {
		token, err := MakeJWT(userID, secret, -time.Hour)
		if err != nil {
			t.Fatalf("failed to create token: %v", err)
		}

		_, err = ValidateJWT(token, secret)
		if err == nil {
			t.Fatal("expected expired token to be rejected")
		}
	})

	t.Run("wrong secret", func(t *testing.T) {
		token, err := MakeJWT(userID, secret, -time.Hour)
		if err != nil {
			t.Fatalf("failed to create token: %v", err)
		}

		_, err = ValidateJWT(token, "wrong-secret")
		if err == nil {
			t.Fatal("expected token signed with wrong secret to be rejected")
		}
	})
}

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name    string
		header  string
		want    string
		wantErr bool
	}{
		{
			name:   "valid bearer token",
			header: "Bearer test-token",
			want:   "test-token",
		},
		{
			name:    "missing authorization header",
			header:  "",
			wantErr: true,
		},
		{
			name:    "invalid authorization scheme",
			header:  "Basic test-token",
			wantErr: true,
		},
		{
			name:    "missing token",
			header:  "Bearer",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers := http.Header{}
			if tt.header != "" {
				headers.Set("Authorization", tt.header)
			}

			got, err := GetBearerToken(headers)

			if (err != nil) != tt.wantErr {
				t.Fatalf("GetBearerToken() error = %v, wantErr %v", err, tt.wantErr)
			}

			if got != tt.want {
				t.Errorf("GetBearerToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

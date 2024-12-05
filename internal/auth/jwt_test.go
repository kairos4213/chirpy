package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWT(t *testing.T) {
	tokenSecret := "secretSecret"
	userID, _ := uuid.Parse("57039a6f-85f6-41c3-83f8-4258aa7bdf7f")
	expiresIn := 2 * time.Hour

	valid, _ := MakeJWT(userID, tokenSecret, expiresIn)
	expired, _ := MakeJWT(userID, tokenSecret, -expiresIn)
	incorrectSecret, _ := MakeJWT(userID, "secret", expiresIn)

	tests := []struct {
		name     string
		secret   string
		tokenStr string
		wantErr  bool
	}{
		{
			name:     "Valid JWT",
			secret:   tokenSecret,
			tokenStr: valid,
			wantErr:  false,
		},
		{
			name:     "Expired JWT",
			secret:   tokenSecret,
			tokenStr: expired,
			wantErr:  true,
		},
		{
			name:     "Incorrect Secret",
			secret:   tokenSecret,
			tokenStr: incorrectSecret,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ValidateJWT(tt.tokenStr, tt.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

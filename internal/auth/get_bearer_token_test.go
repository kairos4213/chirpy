package auth

import (
	"net/http"
	"testing"
)

func TestGetBearerToken(t *testing.T) {
	h1 := make(http.Header)
	h1.Add("Authorization", "Bearer TOKEN_STRING")

	h2 := make(http.Header)
	h2.Add("Authorization", "")

	tests := []struct {
		name     string
		testCase http.Header
		wantStr  string
		wantErr  bool
	}{
		{
			name:     "Bearer Token Exists",
			testCase: h1,
			wantStr:  "TOKEN_STRING",
			wantErr:  false,
		},
		{
			name:     "Null Auth Header",
			testCase: h2,
			wantStr:  "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetBearerToken(tt.testCase)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBearerToken() error = %v, wantErr %v", err, tt.wantErr)
			}
			if result != tt.wantStr {
				t.Errorf("GetBearerToken() result = %v, want %v", result, tt.wantStr)
			}
		})
	}
}

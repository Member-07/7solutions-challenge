package jwt

import (
	"testing"
	"time"
)

func TestTokenServiceGenerateAndParse(t *testing.T) {
	svc, err := NewJwtManager("test-secret", time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	token, err := svc.Generate("user-1")
	if err != nil {
		t.Fatal(err)
	}

	userID, err := svc.Parse(token)
	if err != nil {
		t.Fatal(err)
	}
	if userID != "user-1" {
		t.Fatalf("got %s", userID)
	}
}

func TestNewTokenServiceRequiresSecret(t *testing.T) {
	if _, err := NewJwtManager("", time.Hour); err == nil {
		t.Fatal("expected error")
	}
}

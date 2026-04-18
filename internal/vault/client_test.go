package vault

import (
	"testing"
)

func TestNewClient_MissingAddress(t *testing.T) {
	_, err := NewClient("", "token")
	if err == nil {
		t.Fatal("expected error for empty address, got nil")
	}
}

func TestNewClient_MissingToken(t *testing.T) {
	_, err := NewClient("http://127.0.0.1:8200", "")
	if err == nil {
		t.Fatal("expected error for empty token, got nil")
	}
}

func TestNewClient_Valid(t *testing.T) {
	c, err := NewClient("http://127.0.0.1:8200", "mytoken")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Address != "http://127.0.0.1:8200" {
		t.Errorf("expected address to be set, got %q", c.Address)
	}
	if c.Token != "mytoken" {
		t.Errorf("expected token to be set, got %q", c.Token)
	}
}

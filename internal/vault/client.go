package vault

import (
	"fmt"
	"net/http"
	"time"
)

// Client wraps an HTTP client and Vault connection details.
type Client struct {
	Address string
	Token   string
	http    *http.Client
}

// NewClient creates a new Vault client with the given address and token.
func NewClient(address, token string) (*Client, error) {
	if address == "" {
		return nil, fmt.Errorf("vault address must not be empty")
	}
	if token == "" {
		return nil, fmt.Errorf("vault token must not be empty")
	}
	return &Client{
		Address: address,
		Token:   token,
		http: &http.Client{
			Timeout: 10 * time.Second,
		},
	}, nil
}

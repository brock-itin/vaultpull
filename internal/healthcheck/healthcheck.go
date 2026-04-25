// Package healthcheck provides utilities for verifying Vault connectivity
// and token validity before performing sync operations.
package healthcheck

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Status represents the result of a health check.
type Status struct {
	Reachable    bool
	TokenValid   bool
	VaultVersion string
	Error        error
}

// Options configures health check behaviour.
type Options struct {
	Timeout time.Duration
}

// DefaultOptions returns sensible defaults for health checks.
func DefaultOptions() Options {
	return Options{
		Timeout: 5 * time.Second,
	}
}

// Check verifies that the Vault server at addr is reachable and that
// the provided token is valid. addr should include scheme and host.
func Check(ctx context.Context, addr, token string, opts Options) Status {
	client := &http.Client{Timeout: opts.Timeout}

	sysURL := addr + "/v1/sys/health"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, sysURL, nil)
	if err != nil {
		return Status{Error: fmt.Errorf("building health request: %w", err)}
	}

	resp, err := client.Do(req)
	if err != nil {
		return Status{Error: fmt.Errorf("vault unreachable: %w", err)}
	}
	defer resp.Body.Close()

	// Vault returns 200, 429, 472, 473, 501, 503 from sys/health — all mean
	// the server responded; only a network error means truly unreachable.
	s := Status{Reachable: true}
	s.VaultVersion = resp.Header.Get("X-Vault-Version")

	// Validate the token via /v1/auth/token/lookup-self.
	tokenURL := addr + "/v1/auth/token/lookup-self"
	treq, err := http.NewRequestWithContext(ctx, http.MethodGet, tokenURL, nil)
	if err != nil {
		s.Error = fmt.Errorf("building token request: %w", err)
		return s
	}
	treq.Header.Set("X-Vault-Token", token)

	tresp, err := client.Do(treq)
	if err != nil {
		s.Error = fmt.Errorf("token lookup failed: %w", err)
		return s
	}
	defer tresp.Body.Close()

	if tresp.StatusCode == http.StatusOK {
		s.TokenValid = true
	} else {
		s.Error = fmt.Errorf("token invalid or expired (HTTP %d)", tresp.StatusCode)
	}

	return s
}

// OK returns true when both the server is reachable and the token is valid.
func OK(s Status) bool {
	return s.Reachable && s.TokenValid
}

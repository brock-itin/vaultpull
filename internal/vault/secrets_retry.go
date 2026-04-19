package vault

import (
	"fmt"

	"github.com/example/vaultpull/internal/retry"
)

// GetSecretsWithRetry fetches secrets from the given path, retrying on
// transient errors according to opts.
func (c *Client) GetSecretsWithRetry(path string, opts retry.Options) (map[string]string, error) {
	var result map[string]string
	err := retry.Do(opts, func(attempt int) error {
		secrets, err := c.GetSecrets(path)
		if err != nil {
			return fmt.Errorf("attempt %d: %w", attempt, err)
		}
		result = secrets
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

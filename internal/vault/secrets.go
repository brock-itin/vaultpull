package vault

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type secretResponse struct {
	Data struct {
		Data map[string]string `json:"data"`
	} `json:"data"`
}

// GetSecrets fetches key/value secrets from a KV v2 path in Vault.
func (c *Client) GetSecrets(path string) (map[string]string, error) {
	url := fmt.Sprintf("%s/v1/%s", c.Address, path)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.Token)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden {
		return nil, fmt.Errorf("permission denied: check vault token")
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("secret path not found: %s", path)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d from vault", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	var result secretResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parsing vault response: %w", err)
	}

	return result.Data.Data, nil
}

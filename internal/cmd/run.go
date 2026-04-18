package cmd

import (
	"fmt"
	"os"

	"github.com/vaultpull/internal/config"
	"github.com/vaultpull/internal/env"
	"github.com/vaultpull/internal/vault"
)

// Run executes the main vaultpull sync flow:
// load config → fetch secrets → write .env file.
func Run(outputPath string, overwrite bool) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	client, err := vault.NewClient(cfg.VaultAddr, cfg.VaultToken)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	secrets, err := vault.GetSecrets(client, cfg.SecretPath)
	if err != nil {
		return fmt.Errorf("vault secrets: %w", err)
	}

	if outputPath == "" {
		outputPath = ".env"
	}

	if err := env.Write(outputPath, secrets, overwrite); err != nil {
		return fmt.Errorf("write env: %w", err)
	}

	fmt.Fprintf(os.Stdout, "✓ Synced %d secret(s) to %s\n", len(secrets), outputPath)
	return nil
}

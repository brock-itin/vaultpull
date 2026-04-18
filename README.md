# vaultpull

> CLI tool to sync secrets from HashiCorp Vault into local `.env` files safely

---

## Installation

```bash
go install github.com/yourusername/vaultpull@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourusername/vaultpull/releases).

---

## Usage

Authenticate with your Vault instance and run `vaultpull` from your project root:

```bash
export VAULT_ADDR="https://vault.example.com"
export VAULT_TOKEN="s.your-vault-token"

vaultpull --path secret/myapp --output .env
```

This will pull all key-value pairs from the specified Vault path and write them to your local `.env` file.

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--path` | Vault secret path to read from | *(required)* |
| `--output` | Output file path | `.env` |
| `--overwrite` | Overwrite existing file without prompting | `false` |
| `--addr` | Vault server address | `$VAULT_ADDR` |

### Example Output

```env
DB_HOST=prod-db.internal
DB_PASSWORD=supersecret
API_KEY=abc123xyz
```

---

## Requirements

- Go 1.21+
- A running HashiCorp Vault instance
- A valid Vault token or supported auth method

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

---

## License

[MIT](LICENSE) © 2024 yourusername
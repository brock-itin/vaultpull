// Package lineage tracks the provenance of secrets fetched from HashiCorp Vault.
//
// Each secret key is annotated with the Vault path it was read from, the
// secret version at the time of the fetch, and a timestamp. This information
// is persisted as a JSON sidecar file alongside the generated .env file so
// that operators can audit where every value originated and detect when a
// local copy has drifted from a newer Vault version.
//
// # File format
//
// The sidecar file is a JSON object whose top-level keys mirror the secret
// keys written to the .env file. Each value is an object with the fields:
//
//	{
//	  "path":      "secret/data/myapp",
//	  "version":   5,
//	  "fetched_at": "2024-01-15T10:30:00Z"
//	}
//
// # Typical usage
//
//	record := lineage.Build(secrets, "secret/data/myapp", 5)
//	if err := lineage.Save(".vaultpull/lineage.json", record); err != nil {
//		log.Fatal(err)
//	}
//
// To load and inspect an existing lineage file:
//
//	record, err := lineage.Load(".vaultpull/lineage.json")
//	if err != nil {
//		log.Fatal(err)
//	}
package lineage

// Package lineage tracks the provenance of secrets fetched from HashiCorp Vault.
//
// Each secret key is annotated with the Vault path it was read from, the
// secret version at the time of the fetch, and a timestamp. This information
// is persisted as a JSON sidecar file alongside the generated .env file so
// that operators can audit where every value originated and detect when a
// local copy has drifted from a newer Vault version.
//
// Typical usage:
//
//	record := lineage.Build(secrets, "secret/data/myapp", 5)
//	if err := lineage.Save(".vaultpull/lineage.json", record); err != nil {
//		log.Fatal(err)
//	}
package lineage

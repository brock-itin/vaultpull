// Package inject provides utilities for injecting Vault secrets into
// process environments and child command environments.
//
// # Injecting into the current process
//
//	err := inject.IntoProcess(secrets, inject.DefaultOptions())
//
// # Injecting into a child command
//
//	cmd := exec.Command("myapp")
//	err := inject.IntoCommand(cmd, secrets, inject.DefaultOptions())
//	cmd.Run()
//
// By default existing variables are preserved. Set Options.Overwrite = true
// to replace them. Use Options.Prefix to namespace injected keys.
package inject

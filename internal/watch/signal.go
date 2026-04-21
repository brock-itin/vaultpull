package watch

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
)

// RunWithSignals starts the watcher and cancels it gracefully on SIGINT or SIGTERM.
// It blocks until the watcher exits.
func RunWithSignals(w *Watcher, out io.Writer) error {
	if out == nil {
		out = os.Stderr
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	errCh := make(chan error, 1)
	go func() {
		errCh <- w.Run(ctx)
	}()

	select {
	case sig := <-sigCh:
		fmt.Fprintf(out, "[watch] received signal %s, shutting down\n", sig)
		cancel()
		// wait for watcher to finish
		<-errCh
		return nil
	case err := <-errCh:
		if err == context.Canceled {
			return nil
		}
		return err
	}
}

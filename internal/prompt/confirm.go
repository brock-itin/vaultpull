// Package prompt provides interactive CLI prompts for user confirmation.
package prompt

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Confirmer is the interface for asking yes/no questions.
type Confirmer interface {
	Confirm(question string) (bool, error)
}

// Prompt reads user input from a reader and writes to a writer.
type Prompt struct {
	In  io.Reader
	Out io.Writer
}

// New creates a Prompt with the given reader and writer.
func New(in io.Reader, out io.Writer) *Prompt {
	return &Prompt{In: in, Out: out}
}

// Confirm asks the user a yes/no question and returns true if they answer yes.
// Accepts: y, yes (case-insensitive). Anything else is treated as no.
func (p *Prompt) Confirm(question string) (bool, error) {
	fmt.Fprintf(p.Out, "%s [y/N]: ", question)

	scanner := bufio.NewScanner(p.In)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return false, fmt.Errorf("prompt: read error: %w", err)
		}
		// EOF — treat as no
		return false, nil
	}

	answer := strings.TrimSpace(strings.ToLower(scanner.Text()))
	return answer == "y" || answer == "yes", nil
}

// AutoConfirm always returns true without prompting. Useful for --yes/-y flags.
type AutoConfirm struct{}

// Confirm always returns true.
func (a AutoConfirm) Confirm(_ string) (bool, error) {
	return true, nil
}

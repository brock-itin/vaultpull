// Package chain provides a pipeline for applying sequential transformations
// and filters to a map of secret key-value pairs.
package chain

import "fmt"

// Step represents a single transformation step in the pipeline.
type Step struct {
	Name string
	Fn   func(map[string]string) (map[string]string, error)
}

// Pipeline holds an ordered list of steps to apply to secrets.
type Pipeline struct {
	steps []Step
}

// New returns an empty Pipeline.
func New() *Pipeline {
	return &Pipeline{}
}

// Add appends a named step to the pipeline.
func (p *Pipeline) Add(name string, fn func(map[string]string) (map[string]string, error)) *Pipeline {
	p.steps = append(p.steps, Step{Name: name, Fn: fn})
	return p
}

// Result holds the output of a pipeline run.
type Result struct {
	Secrets  map[string]string
	Applied  []string
	Skipped  []string
}

// Run executes each step in order, passing the output of one step as the input
// to the next. If a step returns an error the pipeline halts and the error is
// wrapped with the step name.
func (p *Pipeline) Run(secrets map[string]string) (Result, error) {
	current := copyMap(secrets)
	result := Result{Skipped: []string{}}

	for _, step := range p.steps {
		out, err := step.Fn(current)
		if err != nil {
			return result, fmt.Errorf("chain step %q: %w", step.Name, err)
		}
		if out == nil {
			result.Skipped = append(result.Skipped, step.Name)
			continue
		}
		current = out
		result.Applied = append(result.Applied, step.Name)
	}

	result.Secrets = current
	return result, nil
}

func copyMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

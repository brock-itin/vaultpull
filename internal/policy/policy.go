// Package policy provides rule-based access control and enforcement
// for secrets fetched from Vault. Policies define which keys are allowed,
// required, or forbidden for a given environment or path.
package policy

import (
	"errors"
	"fmt"
	"strings"
)

// Rule defines a single policy rule applied to a secret key.
type Rule struct {
	// Key is the secret key this rule applies to (case-insensitive).
	Key string `json:"key"`
	// Required indicates the key must be present.
	Required bool `json:"required,omitempty"`
	// Forbidden indicates the key must not be present.
	Forbidden bool `json:"forbidden,omitempty"`
	// AllowEmpty permits the key to have an empty value.
	AllowEmpty bool `json:"allow_empty,omitempty"`
}

// Policy holds a named set of rules.
type Policy struct {
	Name  string `json:"name"`
	Rules []Rule `json:"rules"`
}

// Violation describes a single policy breach.
type Violation struct {
	Key     string
	Message string
}

// Error implements the error interface.
func (v Violation) Error() string {
	return fmt.Sprintf("policy violation [%s]: %s", v.Key, v.Message)
}

// Result holds the outcome of a policy check.
type Result struct {
	Policy     string
	Violations []Violation
}

// OK returns true when no violations were found.
func (r Result) OK() bool {
	return len(r.Violations) == 0
}

// Error returns a combined error string for all violations, or nil.
func (r Result) Error() error {
	if r.OK() {
		return nil
	}
	msgs := make([]string, len(r.Violations))
	for i, v := range r.Violations {
		msgs[i] = v.Error()
	}
	return errors.New(strings.Join(msgs, "; "))
}

// Check evaluates the given secrets map against the policy rules and
// returns a Result containing any violations found.
func Check(p Policy, secrets map[string]string) Result {
	result := Result{Policy: p.Name}

	// Build a normalised lookup of present keys.
	norm := make(map[string]string, len(secrets))
	for k, v := range secrets {
		norm[strings.ToUpper(k)] = v
	}

	for _, rule := range p.Rules {
		normKey := strings.ToUpper(rule.Key)
		val, present := norm[normKey]

		switch {
		case rule.Forbidden && present:
			result.Violations = append(result.Violations, Violation{
				Key:     rule.Key,
				Message: "key is forbidden but present",
			})

		case rule.Required && !present:
			result.Violations = append(result.Violations, Violation{
				Key:     rule.Key,
				Message: "key is required but missing",
			})

		case rule.Required && present && !rule.AllowEmpty && val == "":
			result.Violations = append(result.Violations, Violation{
				Key:     rule.Key,
				Message: "key is required but has an empty value",
			})
		}
	}

	return result
}

// HasViolations returns true when the result contains at least one violation.
func HasViolations(r Result) bool {
	return !r.OK()
}

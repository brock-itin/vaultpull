// Package template renders .env content from a Go text/template string,
// allowing dynamic secret injection into custom file formats.
package template

import (
	"bytes"
	"fmt"
	"text/template"
)

// Options configures template rendering.
type Options struct {
	// LeftDelim and RightDelim override the default {{ }} delimiters.
	LeftDelim  string
	RightDelim string
}

// Render executes the given template string with secrets as the data map.
// Returns the rendered output as a string.
func Render(tmpl string, secrets map[string]string, opts *Options) (string, error) {
	if opts == nil {
		opts = &Options{}
	}

	left := opts.LeftDelim
	right := opts.RightDelim
	if left == "" {
		left = "{{"
	}
	if right == "" {
		right = "}}"
	}

	t, err := template.New("env").
		Delims(left, right).
		Funcs(funcMap()).
		Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("template parse error: %w", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, secrets); err != nil {
		return "", fmt.Errorf("template execute error: %w", err)
	}

	return buf.String(), nil
}

// funcMap returns helper functions available inside templates.
func funcMap() template.FuncMap {
	return template.FuncMap{
		"required": func(key string, m map[string]string) (string, error) {
			v, ok := m[key]
			if !ok || v == "" {
				return "", fmt.Errorf("required secret %q is missing or empty", key)
			}
			return v, nil
		},
		"default": func(def, key string, m map[string]string) string {
			if v, ok := m[key]; ok && v != "" {
				return v
			}
			return def
		},
	}
}

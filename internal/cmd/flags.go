package cmd

import "flag"

// Options holds parsed CLI flag values.
type Options struct {
	Output    string
	Overwrite bool
}

// ParseFlags parses command-line flags and returns Options.
func ParseFlags(args []string) (*Options, error) {
	fs := flag.NewFlagSet("vaultpull", flag.ContinueOnError)

	output := fs.String("output", ".env", "Path to the output .env file")
	overwrite := fs.Bool("overwrite", false, "Overwrite existing values in the .env file")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	return &Options{
		Output:    *output,
		Overwrite: *overwrite,
	}, nil
}

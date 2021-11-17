package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/operator-framework/combo/pkg/combination"
	generate "github.com/operator-framework/combo/pkg/generator"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	ErrEmptyFile = errors.New("empty file")
)

func init() {
	evalCmd.Flags().StringToStringP("replacements", "r", map[string]string{}, "Key value pair of comma delimited values. Example: 'NAMESPACE=foo,bar'")
	if err := evalCmd.MarkFlagRequired("replacements"); err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize eval: %v", err)
		os.Exit(1)
	}
}

// formatReplacements takes a map[string]string from the args and formats them
// in a way that the combinations package wants
func formatReplacements(replacements map[string]string) map[string][]string {
	formattedReplacements := make(map[string][]string)
	for key, val := range replacements {
		formattedReplacements[key] = strings.Split(val, ",")
	}
	return formattedReplacements
}

// validateFile is a simple wrapper to ensure the input/output of valid YAML
func validateFile(file []byte) error {
	var holder interface{}
	err := yaml.Unmarshal(file, &holder)
	if holder == nil {
		return ErrEmptyFile
	}
	return err
}

var (
	evalCmd = &cobra.Command{
		Use:   "eval [file]",
		Short: "Evaluate the combinations for a file at the given path",
		Long: `Evaluate the combinations for a file at the given path. The file provided must be valid YAML.

Note: the combo binary requires the --replacement flag to be explicitly set.

The replacements flag allows users to specify a series of key value pairs in the form of KEY=VALUES.

Example: combo eval -r REPLACE_ME=1,2,3 path/to/file
	`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			replacements, err := cmd.Flags().GetStringToString("replacements")
			if err != nil {
				return fmt.Errorf("failed to access replacements flag: %w", err)
			}

			file, err := os.ReadFile(args[0])
			if err != nil {
				return fmt.Errorf("failed to read file specified: %w", err)
			}

			if err := validateFile(file); err != nil {
				return fmt.Errorf("failed to validate file specified: %w", err)
			}

			combinations := combination.NewStream(
				combination.WithArgs(formatReplacements(replacements)),
				combination.WithSolveAhead(),
			)

			generator := generate.NewGenerator(string(file), combinations)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			generatedDocuments, err := generator.Evaluate(ctx)
			if err != nil {
				return fmt.Errorf("failed to generate combinations: %w", err)
			}

			generatedFile := "---\n" + strings.Join(generatedDocuments, "\n---\n")

			if err := validateFile([]byte(generatedFile)); err != nil {
				return fmt.Errorf("failed to validate file generated: %w", err)
			}

			fmt.Println(generatedFile)

			return nil
		},
	}
)

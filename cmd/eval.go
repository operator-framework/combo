package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/operator-framework/combo/pkg/combination"
	"github.com/operator-framework/combo/pkg/generate"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	ErrEmptyFile        = errors.New("empty file")
	ErrCouldNotReadFile = errors.New("could not read file")
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

// validateFile is a simple wrapper to ensure the file we're using exists, is readable,
// and is valid YAML
func validateFile(file io.Reader) error {
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return ErrCouldNotReadFile
	}

	if len(fileBytes) == 0 {
		return ErrEmptyFile
	}

	var holder interface{}
	return yaml.Unmarshal(fileBytes, &holder)
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

			file, err := os.Open(args[0])
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

			generator, err := generate.NewGenerator(file, combinations)
			if err != nil {
				return fmt.Errorf("failed to construct generator: %w", err)
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			generatedDocuments, err := generator.Generate(ctx)
			if err != nil {
				return fmt.Errorf("failed to generate combinations: %w", err)
			}

			generatedFile := "---\n" + strings.Join(generatedDocuments, "\n---\n")

			if err := validateFile(strings.NewReader(generatedFile)); err != nil {
				return fmt.Errorf("failed to validate file generated: %w", err)
			}

			fmt.Println(generatedFile)

			return nil
		},
	}
)

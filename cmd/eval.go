package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/operator-framework/combo/pkg/combination"
	"github.com/operator-framework/combo/pkg/generator"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	replacements map[string]string
	evalCmd      = &cobra.Command{
		Use:   "eval [file]",
		Short: "Evaluate the combinations for a file at the given path",
		Long: `
Evaluate the combinations for a file at the given path. The file provided must be valid YAML.

Note: the combo binary requires the --replacement flag to be explicitly set.

The replacements flag allows users to specify a series of key value pairs in the form of KEY=VALUE.

Example: combo eval -r REPLACE_ME=1,2,3 path/to/file
	`,
		RunE: run,
	}
)

func init() {
	rootCmd.AddCommand(evalCmd)
	evalCmd.Flags().StringToStringVarP(&replacements, "replacement", "r", map[string]string{}, "Key value pair of comma delimited values. Example: 'NAMESPACE=foo,bar'")
	if err := evalCmd.MarkFlagRequired("replacement"); err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize eval: %v", err)
		os.Exit(1)
	}

}

// run is used during the actual execution of the command to generate
// combinations.
func run(cmd *cobra.Command, args []string) error {
	file, err := os.ReadFile(args[0])
	if err != nil {
		return err
	}

	if err := validateFile(file); err != nil {
		return fmt.Errorf("failed to validate file specified: %v", err)
	}

	combinations := combination.NewStream(
		combination.WithArgs(formatReplacements(replacements)),
		combination.WithSolveAhead(),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	generatedFile, err := generator.Evaluate(ctx, string(file), combinations)
	if err != nil {
		return err
	}

	if err := validateFile([]byte(generatedFile)); err != nil {
		return fmt.Errorf("failed to validate file generated: %v", err)
	}

	fmt.Println(generatedFile)

	return nil
}

// formatReplacements takes a map[string]string from the args and formats them
// in a way that the combinations package wants
func formatReplacements(replacementsInput map[string]string) map[string][]string {
	formattedReplacements := make(map[string][]string)
	for key, val := range replacementsInput {
		formattedReplacements[key] = strings.Split(val, ",")
	}
	return formattedReplacements
}

// validateFile is a simple wrapper to ensure the input/output of valid YAML
func validateFile(file []byte) error {
	var holder interface{}
	return yaml.Unmarshal(file, &holder)
}

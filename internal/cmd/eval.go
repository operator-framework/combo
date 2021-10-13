package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/operator-framework/combo/pkg/combinator"
	"github.com/operator-framework/combo/pkg/generator"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	replacements map[string]string
)

// evalCmd represents the eval command
var evalCmd = &cobra.Command{
	Use:   "eval",
	Short: "Evaluate the combinations for a file at the given path",
	Long: `Evaluate the combinations for a file at the given path. The file provided must be valid YAML.

Example: combo eval -r REPLACE_ME=1,2,3 path/to/file
	`,
	RunE: run,
	Args: cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(evalCmd)
	evalCmd.Flags().StringToStringVarP(&replacements, "replacement", "r", map[string]string{}, "Key value pair of comma delimited values. Example: 'NAMESPACE=foo,bar'")
	evalCmd.MarkFlagRequired("replacement")
}

func run(cmd *cobra.Command, args []string) error {
	file, err := os.ReadFile(args[0])
	if err != nil {
		return err
	}

	if err := validateFile(file); err != nil {
		return fmt.Errorf("failed to validate file specified: %v", err)
	}

	combinations := combinator.Solve(formatReplacements())
	generatedFile, err := generator.Generate(combinations, file)
	if err != nil {
		return err
	}

	if err := validateFile(generatedFile); err != nil {
		return fmt.Errorf("failed to validate file generated: %v", err)
	}

	fmt.Println(string(generatedFile))

	return nil
}

func formatReplacements() map[string][]string {
	formattedReplacements := make(map[string][]string)
	for key, val := range replacements {
		formattedReplacements[key] = strings.Split(val, ",")
	}
	return formattedReplacements
}

func validateFile(file []byte) error {
	var holder interface{}
	if err := yaml.Unmarshal(file, &holder); err != nil {
		return err
	}
	return nil
}

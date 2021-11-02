package eval

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
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "eval",
		Short: "Evaluate the combinations for a file at the given path",
		Long: `Evaluate the combinations for a file at the given path. The file provided must be valid YAML.
	
	Example: combo eval -r REPLACE_ME=1,2,3 path/to/file
		`,
		RunE: run,
		Args: cobra.ExactArgs(1),
	}

	cmd.Flags().StringToStringVarP(&replacements, "replacement", "r", map[string]string{}, "Key value pair of comma delimited values. Example: 'NAMESPACE=foo,bar'")
	cmd.MarkFlagRequired("replacement")

	return cmd
}

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

func formatReplacements(replacementsInput map[string]string) map[string][]string {
	formattedReplacements := make(map[string][]string)
	for key, val := range replacementsInput {
		formattedReplacements[key] = strings.Split(val, ",")
	}
	return formattedReplacements
}

func validateFile(file []byte) error {
	var holder interface{}
	return yaml.Unmarshal(file, &holder)
}

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
	comboErrors "github.com/operator-framework/combo/pkg/errors"
	"github.com/operator-framework/combo/pkg/template"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	ErrEmptyFile      = errors.New("empty file")
	FilePathArgsIndex = 0
	useSolvedAhead    bool
)

func init() {
	evalCmd.Flags().StringToStringP("replacements", "r", map[string]string{}, "Key value pair of comma delimited values. Example: 'NAMESPACE=foo,bar'")
	evalCmd.Flags().BoolVarP(&useSolvedAhead, "presolve", "p", false, "Toggles how combinations are generated. When applied combinations are generated all at once compared to iteratively.")
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
		return fmt.Errorf("%w: %s", comboErrors.ErrCouldNotReadFile, err.Error())
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

			templateFile, err := os.Open(args[FilePathArgsIndex])
			if err != nil {
				return fmt.Errorf("failed to read file specified: %w", err)
			}

			if err := validateFile(templateFile); err != nil {
				return fmt.Errorf("failed to validate file specified: %w", err)
			}

			var combinations combination.Stream
			if useSolvedAhead {
				combinations = combination.NewStream(
					combination.WithArgs(formatReplacements(replacements)),
					combination.WithSolveAhead(),
				)
			} else {
				combinations = combination.NewStream(
					combination.WithArgs(formatReplacements(replacements)),
				)
			}

			templateBuilder, err := template.NewBuilder(templateFile, combinations)
			if err != nil {
				return fmt.Errorf("failed to construct builder: %w", err)
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			combinedTemplateManifests, err := templateBuilder.Build(ctx)
			if err != nil {
				return fmt.Errorf("failed to build manifests with combinations: %w", err)
			}

			combinedTemplate := "---\n" + strings.Join(combinedTemplateManifests, "\n---\n")

			if err := validateFile(strings.NewReader(combinedTemplate)); err != nil {
				return fmt.Errorf("failed to validate combined template constructed: %w", err)
			}

			fmt.Println(combinedTemplate)

			return nil
		},
	}
)

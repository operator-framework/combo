package cmd

import (
	"fmt"

	"github.com/operator-framework/combo/pkg/version"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolP("full", "f", false, "Returns the fully constructed version")
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Return the current version of Combo and its most recent git commit",
	RunE: func(cmd *cobra.Command, args []string) error {
		fullEnabled, err := cmd.Flags().GetBool("full")
		if err != nil {
			return fmt.Errorf("failed to access full flag: %w", err)
		}

		output := version.String()
		if fullEnabled {
			output = version.Full()
		}

		fmt.Println(output)
		return nil
	},
}

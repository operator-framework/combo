package cmd

import (
	"fmt"

	"github.com/operator-framework/combo/pkg/version"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Return the current version of Combo and its most recent git commit",
	Run:   func(cmd *cobra.Command, args []string) { fmt.Println(version.String()) },
}

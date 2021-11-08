package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

// controllerCmd represents the controller command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run Combo as a controller on the cluster",
	Long:  `add long description`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("run called")
	},
}

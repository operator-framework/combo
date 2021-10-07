package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	replacements map[string]string
)

// evalCmd represents the eval command
var evalCmd = &cobra.Command{
	Use:   "eval",
	Short: "Insert short description",
	Long:  `Insert longer description`,
	Run:   run,
}

func init() {
	rootCmd.AddCommand(evalCmd)
	evalCmd.Flags().StringToStringVarP(&replacements, "replacement", "r", map[string]string{}, "Key value pair of comma delimited values. Example: 'NAMESPACE=foo,bar'")
}

func run(cmd *cobra.Command, args []string) {
	fmt.Println("eval called: need to implement")
}

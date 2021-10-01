package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	arguments []string
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
	evalCmd.Flags().StringArrayVarP(&arguments, "argument", "a", []string{}, "Key value pair of comma delimited values. Example: 'NAMESPACE=foo,bar'")
}

func run(cmd *cobra.Command, args []string) {
	fmt.Println("Eval was ran! Woo!")
}

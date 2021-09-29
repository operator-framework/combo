package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
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
}

func run(cmd *cobra.Command, args []string) {
	fmt.Println("Eval was ran! Woo!")
}

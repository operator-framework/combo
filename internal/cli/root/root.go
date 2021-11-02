package root

import (
	"github.com/operator-framework/combo/internal/cli/eval"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "combo",
		Short: "Create combinations of Kubernetes manifests",
	}

	cmd.CompletionOptions.DisableDefaultCmd = true

	cmd.AddCommand(
		eval.NewCommand(),
	)

	return cmd
}

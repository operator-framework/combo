package cmd

import (
	"github.com/operator-framework/combo/pkg/controller"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

func init() {
	runCmd.Flags().Int("verbosity", 1, "Sets verbosity level of combo CR controller with default verbosity 1. Verbosity decreases as the value given increases.")
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run Combo as a controller on the cluster",
	Long: `Run Combo as a controller on the cluster to begin reconciling new events.

This will reconcile any events for the Combination and Template resources.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctrl.SetLogger(rootLog)

		mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
			Scheme: runtime.NewScheme(),
		})
		if err != nil {
			return err
		}

		verbosityLevel, err := cmd.Flags().GetInt("verbosity")
		if err != nil {
			return err
		}

		c, err := controller.NewController(
			mgr.GetClient(),
			ctrl.Log.V(verbosityLevel).WithName("run"),
		)
		if err != nil {
			return nil
		}

		if err = c.ManageWith(mgr); err != nil {
			return err
		}

		return mgr.Start(signals.SetupSignalHandler())
	},
}

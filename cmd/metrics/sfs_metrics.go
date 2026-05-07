package metrics

import (
	"fmt"

	"github.com/spf13/cobra"
)

// TODO: SFS metrics API is not yet available in the SDK.
// This command is a placeholder and will be implemented when the API is added.

func newSFSMetrics() *cobra.Command {
	c := &cobra.Command{
		Use:   "sfs.metrics <metric> [sfs-name]",
		Short: "Query SFS metrics (TODO: not yet implemented)",
		Long:  "Query SFS metrics. This feature is not yet available in the CloudTower SDK.",
		Args:  func(cmd *cobra.Command, args []string) error {
			if listFlag {
				return nil
			}
			return cobra.RangeArgs(1, 2)(cmd, args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if listFlag {
				return ListMetrics(cmd.OutOrStdout(), "sfs")
			}
			return fmt.Errorf("SFS metrics API is not yet implemented in the CloudTower SDK")
		},
	}
	return c
}

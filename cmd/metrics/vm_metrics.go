package metrics

import (
	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/metrics"
	"github.com/spf13/cobra"
	"github.com/smartxworks/cloudtower-go-sdk/v2/models"
)

func newVMMetrics() *cobra.Command {
	c := &cobra.Command{
		Use:   "vm.metrics <metric> [vm-name]",
		Short: "Query VM metrics (elf_*)",
		Long:  "Query VM metrics with optional VM name filter. Example: vm.metrics elf_cpu_usage vm001",
		Args:  func(cmd *cobra.Command, args []string) error {
			if listFlag {
				return nil
			}
			return cobra.RangeArgs(1, 2)(cmd, args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if listFlag {
				return ListMetrics(cmd.OutOrStdout(), "vm")
			}
			cli := client.From(cmd.Context())
			mc := metrics.NewMetricsClient(cli)

			vmName := ""
			if len(args) > 1 {
				vmName = args[1]
			}

			input := &models.GetVMMetricInput{
				Metrics: []string{args[0]},
				Range:   &rangeFlag,
				Vms:     &models.VMWhereInput{Name: &vmName},
			}

			results, err := mc.GetVMMetrics(cmd.Context(), input)
			if err != nil {
				return err
			}

			return renderMetricsResults(cmd.OutOrStdout(), results, vmName, args[0], "vm")
		},
	}
	return c
}

package metrics

import (
	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/metrics"
	"github.com/spf13/cobra"
	"github.com/smartxworks/cloudtower-go-sdk/v2/models"
)

func newVMMetrics() *cobra.Command {
	c := &cobra.Command{
		Use:     "vm.metrics <metric> [vm-name]",
		Short:   "Query VM metrics (elf_*)",
		GroupID: "metrics",
		Long:  "Query VM metrics with optional VM name filter. Example: vm.metrics elf_cpu_usage vm001",
		Args: func(cmd *cobra.Command, args []string) error {
			if listFlag {
				return nil
			}
			return cobra.RangeArgs(1, 2)(cmd, args)
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) == 0 {
				names, err := GetMetricNames("vm")
				if err != nil {
					return nil, cobra.ShellCompDirectiveNoFileComp
				}
				return names, cobra.ShellCompDirectiveNoFileComp
			}
			return nil, cobra.ShellCompDirectiveNoFileComp
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

func RegisterVMMetrics(root *cobra.Command) {
	c := newVMMetrics()
	c.Flags().BoolVar(&listFlag, "list", false, "List available metrics")
	c.Flags().StringVar(&rangeFlag, "range", "5m", "Time range: 5m, 1h, 1d, 7d")
	c.Flags().BoolVar(&latestFlag, "latest", false, "Show only latest value")
	c.Flags().StringVar(&formatFlag, "format", "table", "Output format: table, json, chart")
	root.AddCommand(c)
}
package metrics

import (
	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/metrics"
	"github.com/spf13/cobra"
	"github.com/smartxworks/cloudtower-go-sdk/v2/models"
)

func newVmVolumeMetrics() *cobra.Command {
	c := &cobra.Command{
		Use:     "vm.volume <metric> [vm-name]",
		Short:   "Query VM volume metrics (elf_vm_disk_overall_*)",
		GroupID: "metrics",
		Long:  "Query VM volume metrics by VM name. Metrics: elf_vm_disk_overall_logical_size_bytes, elf_vm_disk_overall_read_iops, etc. Example: vm.volume elf_vm_disk_overall_logical_size_bytes my-vm",
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
				var filtered []string
				for _, n := range names {
					if len(n) >= 20 && n[:20] == "elf_vm_disk_overall_" {
						filtered = append(filtered, n)
					}
				}
				return filtered, cobra.ShellCompDirectiveNoFileComp
			}
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if listFlag {
				return ListMetrics(cmd.OutOrStdout(), "volume")
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

			return renderMetricsResults(cmd.OutOrStdout(), results, vmName, args[0], "volume")
		},
	}
	return c
}

func RegisterVmVolumeMetrics(root *cobra.Command) {
	c := newVmVolumeMetrics()
	c.Flags().BoolVar(&listFlag, "list", false, "List available metrics")
	c.Flags().StringVar(&rangeFlag, "range", "5m", "Time range: 5m, 1h, 1d, 7d")
	c.Flags().BoolVar(&latestFlag, "latest", false, "Show only latest value")
	c.Flags().StringVar(&formatFlag, "format", "table", "Output format: table, json, chart")
	root.AddCommand(c)
}
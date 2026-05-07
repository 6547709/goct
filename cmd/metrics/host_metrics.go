package metrics

import (
	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/metrics"
	"github.com/spf13/cobra"
	"github.com/smartxworks/cloudtower-go-sdk/v2/models"
)

func newHostMetrics() *cobra.Command {
	c := &cobra.Command{
		Use:   "host.metrics <metric> [host-name]",
		Short: "Query Host metrics",
		Long:  "Query Host metrics with optional host name filter. Example: host.metrics elf_host_cpu_usage host001",
		Args:  func(cmd *cobra.Command, args []string) error {
			if listFlag {
				return nil
			}
			return cobra.RangeArgs(1, 2)(cmd, args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if listFlag {
				return ListMetrics(cmd.OutOrStdout(), "host")
			}
			cli := client.From(cmd.Context())
			mc := metrics.NewMetricsClient(cli)

			hostName := ""
			if len(args) > 1 {
				hostName = args[1]
			}

			input := &models.GetHostMetricInput{
				Metrics: []string{args[0]},
				Range:   &rangeFlag,
				Hosts:   &models.HostWhereInput{Name: &hostName},
			}

			results, err := mc.GetHostMetrics(cmd.Context(), input)
			if err != nil {
				return err
			}

			return renderMetricsResults(cmd.OutOrStdout(), results, hostName, args[0], "host")
		},
	}
	return c
}

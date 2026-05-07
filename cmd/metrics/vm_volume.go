package metrics

import (
	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/metrics"
	"github.com/spf13/cobra"
	"github.com/smartxworks/cloudtower-go-sdk/v2/models"
)

func newVmVolumeMetrics() *cobra.Command {
	c := &cobra.Command{
		Use:   "vm.volume <metric> [volume-name]",
		Short: "Query VM volume metrics (zbs_volume_*)",
		Long:  "Query VM volume metrics with optional volume name filter. Example: vm.volume zbs_volume_read_iops my-volume",
		Args:  func(cmd *cobra.Command, args []string) error {
			if listFlag {
				return nil
			}
			return cobra.RangeArgs(1, 2)(cmd, args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if listFlag {
				return ListMetrics(cmd.OutOrStdout(), "volume")
			}
			cli := client.From(cmd.Context())
			mc := metrics.NewMetricsClient(cli)

			volumeName := ""
			if len(args) > 1 {
				volumeName = args[1]
			}

			input := &models.GetVMVolumeMetricInput{
				Metrics:   []string{args[0]},
				Range:     &rangeFlag,
				VMVolumes: &models.VMVolumeWhereInput{Name: &volumeName},
			}

			results, err := mc.GetVmVolumeMetrics(cmd.Context(), input)
			if err != nil {
				return err
			}

			return renderMetricsResults(cmd.OutOrStdout(), results, volumeName, args[0], "volume")
		},
	}
	return c
}

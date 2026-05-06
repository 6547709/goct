package metrics

import (
	"github.com/spf13/cobra"
)

var (
	listFlag   bool
	latestFlag bool
	rangeFlag  string
	formatFlag string
)

var rootCmd = &cobra.Command{
	Use:   "metrics",
	Short: "Query CloudTower metrics",
	Long:  "Query VM, Host, Volume, Cluster and SFS metrics",
}

func Register(root *cobra.Command) {
	root.AddCommand(rootCmd)

	rootCmd.PersistentFlags().BoolVar(&listFlag, "list", false, "List available metrics")
	rootCmd.PersistentFlags().StringVar(&rangeFlag, "range", "5m", "Time range: 5m, 1h, 1d, 7d")
	rootCmd.PersistentFlags().BoolVar(&latestFlag, "latest", false, "Show only latest value")
	rootCmd.PersistentFlags().StringVar(&formatFlag, "format", "table", "Output format: table, json, chart")
}
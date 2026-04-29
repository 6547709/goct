package datastore

import (
	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/output"
	"github.com/6547709/goct/pkg/service"
	"github.com/spf13/cobra"
)

func newDiskLs() *cobra.Command {
	var hostID string
	c := &cobra.Command{
		Use: "datastore.disk.ls", Short: "List physical disks", GroupID: groupID,
		RunE: func(c *cobra.Command, _ []string) error {
			cli := client.From(c.Context())
			disks, err := service.NewDatastore(cli).ListDisks(c.Context(), hostID)
			if err != nil { return err }
			out := make([]any, len(disks))
			for i := range disks { out[i] = disks[i] }
			format, _ := c.Flags().GetString("format")
			return output.Render(c.OutOrStdout(), out, format, output.DiskListColumns)
		},
	}
	c.Flags().StringVar(&hostID, "host", "", "Filter by host ID")
	return c
}

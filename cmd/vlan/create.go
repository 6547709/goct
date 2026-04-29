package vlan

import (
	"fmt"
	"github.com/6547709/goct/pkg/adapter"
	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/service"
	"github.com/6547709/goct/pkg/task"
	"github.com/spf13/cobra"
)

func newCreate() *cobra.Command {
	var name, vdsID string
	c := &cobra.Command{
		Use: "vlan.create", Short: "Create a VLAN", GroupID: groupID,
		RunE: func(c *cobra.Command, _ []string) error {
			cli := client.From(c.Context())
			ref, err := service.NewVLAN(cli).Create(c.Context(), adapter.VLANCreateSpec{Name: name, VdsID: vdsID})
			if err != nil { return err }
			if ref.IsSync() { fmt.Fprintln(c.OutOrStdout(), "VLAN created (sync)"); return nil }
			return task.New(cli, task.Options{Out: c.OutOrStderr()}).Watch(c.Context(), ref.ID)
		},
	}
	c.Flags().StringVar(&name, "name", "", "VLAN name (required)")
	c.Flags().StringVar(&vdsID, "vds", "", "VDS ID (required)")
	_ = c.MarkFlagRequired("name"); _ = c.MarkFlagRequired("vds")
	return c
}

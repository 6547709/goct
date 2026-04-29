package vlan

import (
	"fmt"
	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/service"
	"github.com/6547709/goct/pkg/task"
	"github.com/spf13/cobra"
)

func newDestroy() *cobra.Command {
	return &cobra.Command{
		Use: "vlan.destroy <name|id>", Short: "Delete a VLAN", GroupID: groupID,
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			ref, err := service.NewVLAN(cli).Delete(c.Context(), args[0])
			if err != nil { return err }
			if ref.IsSync() { fmt.Fprintln(c.OutOrStdout(), "VLAN deleted (sync)"); return nil }
			return task.New(cli, task.Options{Out: c.OutOrStderr()}).Watch(c.Context(), ref.ID)
		},
	}
}

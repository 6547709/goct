package vm

import (
	"fmt"

	"github.com/6547709/goct/pkg/adapter"
	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/service"
	"github.com/6547709/goct/pkg/task"
	"github.com/spf13/cobra"
)

func newExport() *cobra.Command {
	var keepMAC bool
	c := &cobra.Command{
		Use: "vm.export [name|id]", Short: "Export a VM as OVF", GroupID: groupID,
		Args: cobra.MaximumNArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			id, err := resolveVMArg(args)
			if err != nil {
				return err
			}
			ref, err := service.NewVM(cli).Export(c.Context(), id, adapter.VMExportSpec{
				KeepMAC: keepMAC,
			})
			if err != nil {
				return err
			}
			if ref.IsSync() {
				fmt.Fprintln(c.OutOrStdout(), "VM exported (sync)")
				return nil
			}
			w := task.New(cli, task.Options{Out: c.OutOrStderr()})
			return w.Watch(c.Context(), ref.ID)
		},
	}
	c.Flags().BoolVar(&keepMAC, "keep-mac", false, "Keep MAC addresses in exported OVF")
	return c
}

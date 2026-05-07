package vm

import (
	"github.com/6547709/goct/pkg/adapter"
	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/service"
	"github.com/6547709/goct/pkg/task"
	"github.com/spf13/cobra"
)

func newUpdate() *cobra.Command {
	var name, description string
	c := &cobra.Command{
		Use:   "vm.update [name|id]",
		Short: "Update VM name or description",
		Long: `Update VM basic information (name, description).

Examples:
  goct vm.update myvm --name newname
  goct vm.update myvm --description "production server"`,
		GroupID: groupID,
		Args:    cobra.MaximumNArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			id, err := resolveVMArg(args)
			if err != nil {
				return err
			}
			ref, err := service.NewVM(cli).Update(c.Context(), id, adapter.VMUpdateSpec{
				Name:        name,
				Description: description,
			})
			if err != nil {
				return err
			}
			if ref.IsSync() {
				return nil
			}
			w := task.New(cli, task.Options{Out: c.OutOrStderr()})
			return w.Watch(c.Context(), ref.ID)
		},
	}
	c.Flags().StringVar(&name, "name", "", "New VM name")
	c.Flags().StringVar(&description, "description", "", "New VM description")
	return c
}

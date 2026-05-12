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

Use --description='' (empty string) to explicitly clear the description.
Flags that are not provided are not sent and keep the original value.

Examples:
  goct vm.update myvm --name newname
  goct vm.update myvm --description "production server"
  goct vm.update myvm --description ""             # clear description`,
		GroupID: groupID,
		Args:    cobra.MaximumNArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			id, err := resolveVMArg(args)
			if err != nil {
				return err
			}
			spec := adapter.VMUpdateSpec{}
			if c.Flags().Changed("name") {
				spec.Name = &name
			}
			if c.Flags().Changed("description") {
				spec.Description = &description
			}
			ref, err := service.NewVM(cli).Update(c.Context(), id, spec)
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
	c.Flags().StringVar(&name, "name", "", "New VM name (omit to keep, '' is rejected by CloudTower)")
	c.Flags().StringVar(&description, "description", "", "New VM description ('' clears it)")
	return c
}

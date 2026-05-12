package vm

import (
	"fmt"

	"github.com/6547709/goct/pkg/adapter"
	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/service"
	"github.com/6547709/goct/pkg/task"
	"github.com/spf13/cobra"
)

func newRebuild() *cobra.Command {
	var name, clusterID, hostID string
	c := &cobra.Command{
		Use:     "vm.rebuild [vm-name|vm-id] [snapshot-id]",
		Short:   "Rebuild VM from snapshot",
		GroupID: "vm",
		Long: `Rebuild a virtual machine from a snapshot.

Examples:
  goct vm rebuild my-vm snap-uuid
  goct vm rebuild my-vm snap-uuid --name new-name`,
		Args: cobra.ExactArgs(2),
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			vmID := args[0]
			snapshotID := args[1]
			ref, err := service.NewVM(cli).RebuildVM(c.Context(), vmID, adapter.RebuildVMSpec{
				SnapshotID: snapshotID,
				Name:       name,
				ClusterID:  clusterID,
				HostID:     hostID,
			})
			if err != nil {
				return err
			}
			if ref.IsSync() {
				_, _ = fmt.Fprintln(c.OutOrStdout(), "vm rebuilt")
			} else {
				_, _ = fmt.Fprintf(c.OutOrStdout(), "task: %s\n", ref.ID)
				w := task.New(cli, task.Options{Out: c.OutOrStderr()})
				return w.Watch(c.Context(), ref.ID)
			}
			return nil
		},
	}
	c.Flags().StringVar(&name, "name", "", "New VM name")
	c.Flags().StringVar(&clusterID, "cluster", "", "Target cluster ID")
	c.Flags().StringVar(&hostID, "host", "", "Target host ID")
	return c
}

func newResetPassword() *cobra.Command {
	var username, password string
	c := &cobra.Command{
		Use:     "vm.reset-password [vm-name|vm-id]",
		Short:   "Reset guest OS password",
		GroupID: "vm",
		Long: `Reset the guest OS administrator password on a VM.

Examples:
  goct vm reset-password my-vm --username admin --password newpass123`,
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			if username == "" || password == "" {
				return fmt.Errorf("--username and --password required")
			}
			ref, err := service.NewVM(cli).ResetPassword(c.Context(), args[0], username, password)
			if err != nil {
				return err
			}
			if ref.IsSync() {
				_, _ = fmt.Fprintln(c.OutOrStdout(), "password reset")
			} else {
				_, _ = fmt.Fprintf(c.OutOrStdout(), "task: %s\n", ref.ID)
			}
			return nil
		},
	}
	c.Flags().StringVar(&username, "username", "", "Guest OS username (required)")
	c.Flags().StringVar(&password, "password", "", "New password (required)")
	return c
}

func newMigrateAbort() *cobra.Command {
	c := &cobra.Command{
		Use:     "vm.migrate.abort [vm-name|vm-id]",
		Short:   "Abort cross-cluster migration",
		GroupID: "vm",
		Long: `Abort an in-progress cross-cluster migration.

Examples:
  goct vm migrate.abort my-vm`,
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			ref, err := service.NewVM(cli).AbortMigrateAcrossCluster(c.Context(), args[0])
			if err != nil {
				return err
			}
			if ref.IsSync() {
				_, _ = fmt.Fprintln(c.OutOrStdout(), "migration aborted")
			} else {
				_, _ = fmt.Fprintf(c.OutOrStdout(), "task: %s\n", ref.ID)
			}
			return nil
		},
	}
	return c
}

func newConvertToVM() *cobra.Command {
	var newName string
	c := &cobra.Command{
		Use:   "convert-to-vm [template-name|template-id]",
		Short: "Convert template to VM",
		Long: `Convert a content library template to a virtual machine.

Examples:
  goct vm convert-to-vm my-template
  goct vm convert-to-vm my-template --name my-new-vm`,
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			ref, err := service.NewVM(cli).ConvertToVM(c.Context(), args[0], newName)
			if err != nil {
				return err
			}
			if ref.IsSync() {
				_, _ = fmt.Fprintln(c.OutOrStdout(), "template converted to vm")
			} else {
				_, _ = fmt.Fprintf(c.OutOrStdout(), "task: %s\n", ref.ID)
			}
			return nil
		},
	}
	c.Flags().StringVar(&newName, "name", "", "Name for the converted VM (default: <template>-vm)")
	return c
}
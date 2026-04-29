package task

import (
	"github.com/6547709/goct/pkg/client"
	pkgtask "github.com/6547709/goct/pkg/task"
	"github.com/spf13/cobra"
)

func newWait() *cobra.Command {
	return &cobra.Command{
		Use: "task.wait <id>", Short: "Wait for a task to complete", GroupID: groupID,
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			w := pkgtask.New(cli, pkgtask.Options{Out: c.OutOrStderr()})
			return w.Watch(c.Context(), args[0])
		},
	}
}

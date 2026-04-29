package alert

import (
	"fmt"
	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/service"
	"github.com/6547709/goct/pkg/task"
	"github.com/spf13/cobra"
)

func newAck() *cobra.Command {
	return &cobra.Command{
		Use: "alert.ack <id>", Short: "Acknowledge (resolve) an alert", GroupID: groupID,
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			ref, err := service.NewAlert(cli).Ack(c.Context(), args[0])
			if err != nil { return err }
			if ref.IsSync() { fmt.Fprintln(c.OutOrStdout(), "Alert acknowledged"); return nil }
			return task.New(cli, task.Options{Out: c.OutOrStderr()}).Watch(c.Context(), ref.ID)
		},
	}
}

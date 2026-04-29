package task

import (
	"fmt"
	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/output"
	"github.com/6547709/goct/pkg/service"
	"github.com/spf13/cobra"
)

func newInfo() *cobra.Command {
	return &cobra.Command{
		Use: "task.info <id>", Short: "Show task details", GroupID: groupID,
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			t, err := service.NewTask(cli).Get(c.Context(), args[0])
			if err != nil { return err }
			format, _ := c.Flags().GetString("format")
			if format == "json" { return output.Render(c.OutOrStdout(), []any{*t}, format, nil) }
			fmt.Fprintf(c.OutOrStdout(), "%-14s %s\n", "ID:", t.ID)
			fmt.Fprintf(c.OutOrStdout(), "%-14s %s\n", "Status:", t.Status)
			fmt.Fprintf(c.OutOrStdout(), "%-14s %d%%\n", "Progress:", t.Progress)
			fmt.Fprintf(c.OutOrStdout(), "%-14s %s\n", "Description:", t.Description)
			if t.ErrorMessage != "" {
				fmt.Fprintf(c.OutOrStdout(), "%-14s %s\n", "Error:", t.ErrorMessage)
			}
			fmt.Fprintf(c.OutOrStdout(), "%-14s %s\n", "Created:", t.CreatedAt)
			return nil
		},
	}
}

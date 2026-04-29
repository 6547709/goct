package task

import (
	"fmt"
	"github.com/6547709/goct/pkg/adapter"
	"github.com/spf13/cobra"
)

func newCancel() *cobra.Command {
	return &cobra.Command{
		Use: "task.cancel <id>", Short: "Cancel a task (not supported by SDK)", GroupID: groupID,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, _ []string) error {
			return fmt.Errorf("task.cancel: %w", adapter.ErrUnsupported)
		},
	}
}

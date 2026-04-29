package system

import (
	"fmt"

	"github.com/6547709/goct/internal/version"
	"github.com/spf13/cobra"
)

func newVersion() *cobra.Command {
	return &cobra.Command{
		Use:               "version",
		Short:             "Print goct version",
		GroupID:           groupID,
		Annotations:       map[string]string{"nologin": "true"},
		Run: func(c *cobra.Command, _ []string) {
			fmt.Fprintln(c.OutOrStdout(), version.Full())
		},
	}
}

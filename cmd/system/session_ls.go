package system

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/6547709/goct/pkg/session"
	"github.com/spf13/cobra"
)

func newSessionLs() *cobra.Command {
	return &cobra.Command{
		Use:               "session.ls",
		Short:             "List locally cached session files",
		GroupID:           groupID,
		Annotations:       map[string]string{"nologin": "true"},
		RunE: func(c *cobra.Command, _ []string) error {
			paths, err := session.List()
			if err != nil {
				fmt.Fprintln(c.OutOrStdout(), "No cached sessions found.")
				return nil
			}
			if len(paths) == 0 {
				fmt.Fprintln(c.OutOrStdout(), "No cached sessions found.")
				return nil
			}
			fmt.Fprintf(c.OutOrStdout(), "%-48s  %s\n", "SESSION FILE", "MODIFIED")
			fmt.Fprintln(c.OutOrStdout(), "------------------------------------------------------------------------")
			for _, p := range paths {
				info, err := os.Stat(p)
				modTime := "unknown"
				if err == nil {
					modTime = info.ModTime().Format("2006-01-02 15:04:05")
				}
				fmt.Fprintf(c.OutOrStdout(), "%-48s  %s\n", filepath.Base(p), modTime)
			}
			return nil
		},
	}
}

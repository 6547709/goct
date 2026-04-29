package system

import (
	"fmt"
	"net/url"

	"github.com/6547709/goct/pkg/client"
	"github.com/spf13/cobra"
)

// sessionLoginCommand forces a fresh login and saves the new token to cache.
// It lets PersistentPreRunE perform the normal login (cache miss → login path),
// then confirms the session is usable by calling About().
func newSessionLogin() *cobra.Command {
	return &cobra.Command{
		Use:     "session.login",
		Short:   "Force login and save session token to local cache",
		GroupID: groupID,
		RunE: func(c *cobra.Command, _ []string) error {
			cli := client.From(c.Context())
			if cli == nil {
				// Should not happen — nologin=false → PersistentPreRunE runs.
				return nil
			}
			// Lightweight call to confirm the session is usable.
			info, err := cli.About(c.Context())
			if err != nil {
				return err
			}
			host := hostFlag(c)
			fmt.Fprintf(c.OutOrStdout(),
				"Login successful: host=%s version=%s\n", host, info.Version)
			return nil
		},
	}
}

// hostFlag extracts the host from the --url flag.
func hostFlag(c *cobra.Command) string {
	if u, _ := c.Flags().GetString("url"); u != "" {
		if pu, err := url.Parse(u); err == nil && pu.Host != "" {
			return pu.Host
		}
		return u
	}
	return "(unknown host)"
}

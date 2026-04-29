package system

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/6547709/goct/pkg/session"
	"github.com/spf13/cobra"
)

func newSessionLogout() *cobra.Command {
	var urlFlag, userFlag string
	cmd := &cobra.Command{
		Use:               "session.logout",
		Short:             "Delete cached session token for a given host and user",
		GroupID:           groupID,
		Annotations:       map[string]string{"nologin": "true"},
		RunE: func(c *cobra.Command, _ []string) error {
			if urlFlag == "" || userFlag == "" {
				return errors.New("session.logout: --url and --user are required")
			}
			host := urlFlag
			if pu, err := url.Parse(urlFlag); err == nil && pu.Host != "" {
				host = pu.Host
			}
			if err := session.Delete(host, userFlag); err != nil {
				return fmt.Errorf("session.logout: %w", err)
			}
			fmt.Fprintf(c.OutOrStdout(), "Session deleted: host=%s user=%s\n", host, userFlag)
			return nil
		},
	}
	cmd.Flags().StringVar(&urlFlag, "url", "", "CloudTower URL")
	cmd.Flags().StringVar(&userFlag, "user", "", "Username")
	_ = cmd.MarkFlagRequired("url")
	_ = cmd.MarkFlagRequired("user")
	return cmd
}

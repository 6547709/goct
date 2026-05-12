package ntp

import (
	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/output"
	"github.com/6547709/goct/pkg/service"
	"github.com/spf13/cobra"
)

func newGet() *cobra.Command {
	c := &cobra.Command{
		Use:   "ntp.get",
		Short: "Get NTP service URL",
		GroupID: groupID,
		RunE: func(c *cobra.Command, _ []string) error {
			cli := client.From(c.Context())
			settings, err := service.NewNTP(cli).GetSettings(c.Context())
			if err != nil {
				return err
			}
			out := make([]any, 0)
			for _, url := range settings.URLs {
				out = append(out, url)
			}
			format, _ := c.Flags().GetString("format")
			return output.Render(c.OutOrStdout(), out, format, []output.Column{
				{Header: "NTP URLs", Get: func(v any) string { return v.(string) }},
			})
		},
	}
	return c
}

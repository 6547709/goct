package license

import (
	"github.com/6547709/goct/pkg/adapter"
	"github.com/6547709/goct/pkg/client"
	gflags "github.com/6547709/goct/pkg/flags"
	"github.com/6547709/goct/pkg/output"
	"github.com/6547709/goct/pkg/service"
	"github.com/spf13/cobra"
)

var licenseListColumns = []output.Column{
	{Header: "ID", Get: func(v any) string { return v.(adapter.License).ID }},
	{Header: "Serial", Get: func(v any) string { return v.(adapter.License).LicenseSerial }},
	{Header: "Type", Get: func(v any) string { return v.(adapter.License).Type }},
	{Header: "Edition", Get: func(v any) string { return v.(adapter.License).SoftwareEdition }},
	{Header: "Expire Date", Get: func(v any) string { return v.(adapter.License).ExpireDate }},
}

func newLs() *cobra.Command {
	var sf gflags.SearchFlags
	c := &cobra.Command{
		Use:   "license.ls",
		Short: "List licenses",
		GroupID: groupID,
		RunE: func(c *cobra.Command, _ []string) error {
			cli := client.From(c.Context())
			items, err := service.NewLicense(cli).List(c.Context(), adapter.ListOpts{NameContains: sf.Name, Limit: sf.Limit})
			if err != nil {
				return err
			}
			out := make([]any, len(items))
			for i := range items {
				out[i] = items[i]
			}
			format, _ := c.Flags().GetString("format")
			return output.Render(c.OutOrStdout(), out, format, licenseListColumns)
		},
	}
	sf.Register(c)
	return c
}

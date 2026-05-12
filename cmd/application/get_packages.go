package application

import (
	"github.com/6547709/goct/pkg/adapter"
	"github.com/6547709/goct/pkg/client"
	gflags "github.com/6547709/goct/pkg/flags"
	"github.com/6547709/goct/pkg/output"
	"github.com/6547709/goct/pkg/service"
	"github.com/spf13/cobra"
)

var cloudTowerApplicationPackageListColumns = []output.Column{
	{Header: "ID", Get: func(v any) string { return v.(adapter.CloudTowerApplicationPackage).ID }},
	{Header: "Name", Get: func(v any) string { return v.(adapter.CloudTowerApplicationPackage).Name }},
	{Header: "Version", Get: func(v any) string { return v.(adapter.CloudTowerApplicationPackage).Version }},
	{Header: "Architecture", Get: func(v any) string { return v.(adapter.CloudTowerApplicationPackage).Architecture }},
	{Header: "SCOS Version", Get: func(v any) string { return v.(adapter.CloudTowerApplicationPackage).ScosVersion }},
}

func newGetPackages() *cobra.Command {
	var sf gflags.SearchFlags
	c := &cobra.Command{
		Use:   "get-cloudtower-application-packages",
		Short: "List CloudTower application packages",
		GroupID: groupID,
		RunE: func(c *cobra.Command, _ []string) error {
			cli := client.From(c.Context())
			items, err := service.NewCloudTowerApplication(cli).ListPackages(c.Context(), adapter.ListOpts{NameContains: sf.Name, Limit: sf.Limit})
			if err != nil {
				return err
			}
			out := make([]any, len(items))
			for i := range items {
				out[i] = items[i]
			}
			format, _ := c.Flags().GetString("format")
			return output.Render(c.OutOrStdout(), out, format, cloudTowerApplicationPackageListColumns)
		},
	}
	sf.Register(c)
	return c
}

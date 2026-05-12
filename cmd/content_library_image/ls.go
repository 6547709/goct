package content_library_image

import (
	"fmt"

	"github.com/6547709/goct/pkg/adapter"
	"github.com/6547709/goct/pkg/client"
	gflags "github.com/6547709/goct/pkg/flags"
	"github.com/6547709/goct/pkg/output"
	"github.com/6547709/goct/pkg/service"
	"github.com/spf13/cobra"
)

var imageListColumns = []output.Column{
	{Header: "ID", Get: func(v any) string { return v.(adapter.ContentLibraryImage).ID }},
	{Header: "Name", Get: func(v any) string { return v.(adapter.ContentLibraryImage).Name }},
	{Header: "Size", Get: func(v any) string { return formatSize(v.(adapter.ContentLibraryImage).Size) }},
	{Header: "Created", Get: func(v any) string { return v.(adapter.ContentLibraryImage).CreatedAt }},
}

func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%dB", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func newLs() *cobra.Command {
	var sf gflags.SearchFlags
	c := &cobra.Command{
		Use:   "content-library-image.ls",
		Short: "List content library images",
		GroupID: groupID,
		RunE: func(c *cobra.Command, _ []string) error {
			cli := client.From(c.Context())
			items, err := service.NewContentLibraryImage(cli).List(c.Context(), adapter.ListOpts{NameContains: sf.Name, Limit: sf.Limit})
			if err != nil {
				return err
			}
			out := make([]any, len(items))
			for i := range items {
				out[i] = items[i]
			}
			format, _ := c.Flags().GetString("format")
			return output.Render(c.OutOrStdout(), out, format, imageListColumns)
		},
	}
	sf.Register(c)
	return c
}

package content_library_image

import (
	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/service"
	"github.com/spf13/cobra"
)

func newDelete() *cobra.Command {
	c := &cobra.Command{
		Use:   "content-library-image.delete <id>",
		Short: "Delete content library image",
		GroupID: groupID,
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			return service.NewContentLibraryImage(cli).Delete(c.Context(), args[0])
		},
	}
	return c
}

package host

import (
	"github.com/6547709/goct/pkg/service"
	"github.com/spf13/cobra"
)

func newDisconnect() *cobra.Command {
	return &cobra.Command{
		Use: "host.disconnect [name|id]", Short: "Disconnect a host (not supported by SDK)", GroupID: groupID,
		Args: cobra.MaximumNArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			id, err := resolveHostArg(args)
			if err != nil {
				return err
			}
			_, err = service.NewHost(nil).Disconnect(c.Context(), id)
			return err
		},
	}
}

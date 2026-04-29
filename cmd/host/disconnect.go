package host

import (
	"github.com/6547709/goct/pkg/service"
	"github.com/spf13/cobra"
)

func newDisconnect() *cobra.Command {
	return &cobra.Command{
		Use: "host.disconnect <name|id>", Short: "Disconnect a host (not supported by SDK)", GroupID: groupID,
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			_, err := service.NewHost(nil).Disconnect(c.Context(), args[0])
			return err
		},
	}
}

package vm

import (
	"fmt"

	"github.com/6547709/goct/pkg/adapter"
	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/service"
	"github.com/openlyinc/pointy"
	"github.com/spf13/cobra"
)

func newNicUpdate() *cobra.Command {
	var nicIndex int32
	var connectVlanID string
	var enableFlag, disableFlag bool
	var gateway, ipAddress, macAddress, model, subnetMask string

	c := &cobra.Command{
		Use:   "nic.update",
		Short: "Update VM NIC configuration",
		Long: `Update a NIC configuration on a VM.

Examples:
  goct vm nic.update --nic-index 1 --gateway 192.168.1.1
  goct vm nic.update --nic-index 1 --model VIRTIO
  goct vm nic.update --nic-index 1 --connect-vlan-id vlan-uuid`,
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			spec := adapter.VMNicUpdateSpec{
				NicIndex: nicIndex,
			}
			if connectVlanID != "" {
				spec.ConnectVlanID = connectVlanID
			}
			if gateway != "" {
				spec.Gateway = gateway
			}
			if ipAddress != "" {
				spec.IPAddress = ipAddress
			}
			if macAddress != "" {
				spec.MacAddress = macAddress
			}
			if model != "" {
				spec.Model = model
			}
			if subnetMask != "" {
				spec.SubnetMask = subnetMask
			}
			// Handle enable/disable flags
			if enableFlag && !disableFlag {
				spec.Enabled = pointy.Bool(true)
			} else if disableFlag && !enableFlag {
				spec.Enabled = pointy.Bool(false)
			}
			ref, err := service.NewVM(cli).UpdateNic(c.Context(), "", spec)
			if err != nil {
				return err
			}
			if ref.IsSync() {
				_, _ = fmt.Fprintln(c.OutOrStdout(), "nic updated")
			} else {
				_, _ = fmt.Fprintf(c.OutOrStdout(), "task: %s\n", ref.ID)
			}
			return nil
		},
	}
	c.Flags().Int32Var(&nicIndex, "nic-index", 0, "NIC index (LocalID)")
	c.Flags().StringVar(&connectVlanID, "connect-vlan-id", "", "Connect to VLAN ID")
	c.Flags().StringVar(&gateway, "gateway", "", "Gateway IP")
	c.Flags().StringVar(&ipAddress, "ip", "", "IP address")
	c.Flags().StringVar(&macAddress, "mac", "", "MAC address")
	c.Flags().StringVar(&model, "model", "", "Model (RTL8139/E1000/VIRTIO)")
	c.Flags().StringVar(&subnetMask, "subnet-mask", "", "Subnet mask")
	c.Flags().BoolVar(&enableFlag, "enable", false, "Enable NIC")
	c.Flags().BoolVar(&disableFlag, "disable", false, "Disable NIC")
	return c
}
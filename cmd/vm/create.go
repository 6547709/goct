package vm

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/6547709/goct/pkg/adapter"
	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/service"
	"github.com/6547709/goct/pkg/task"
	"github.com/spf13/cobra"
)

// newCreate 创建 VM。
//
// 设计原则（v0.2.1 修复）：
//   - 不再隐式塞默认磁盘 / 默认 NIC：用户必须用 --disk / --nic 显式声明，与 govc vm.create 对齐。
//   - --disk 和 --nic 都是可重复 flag，可多次出现以构造多盘多卡。
//   - --ha 默认不下发（用 CloudTower 默认行为），传 --ha=true|false 时显式覆盖。
//
// Flag 语法：
//
//	--disk size=10g[,bus=SCSI][,name=diskN][,index=N][,boot=N][,iops=N]
//	--nic  vlan=<id|name>[,model=VIRTIO|E1000][,type=VLAN|VPC]
func newCreate() *cobra.Command {
	var (
		name              string
		clusterID         string
		vcpu              int32
		memoryMiB         int64
		firmware          string
		description       string
		fromTemplate      string
		isFullCopy        bool
		nicType           string
		nicModel          string
		nicVlan           string
		ha                string
		disks             []string
		nics              []string
		cloudInitHostname string
		cloudInitPassword string
		cloudInitSSHKey   []string
		cloudInitDNSServers []string
		cloudInitUserData string
		cloudInitNetworks []string
	)
	c := &cobra.Command{
		Use: "vm.create", Short: "Create a new VM", GroupID: groupID,
		Long: `Create a new virtual machine.

Two modes:
  1) From template (--from-template): use existing content library template.
  2) From scratch: build a VM with --disk / --nic flags.

Cloud-init options (for --from-template):
  --hostname string       Cloud-init hostname
  --password string      Default user password
  --ssh-key stringArray  SSH public key (literal string or @/path/to/key.pub, repeatable)
  --dns stringArray      DNS nameserver (repeatable)
  --user-data string     Cloud-init user_data (@/path/to/file.yaml or literal YAML)
  --network stringArray  NIC config: nic=0[,ip=x][,netmask=x][,gateway=x][,route=x][,type=IPV4|DHCP]

Examples:
  # From template with cloud-init (static IP)
  goct vm.create --name web1 --cluster c1 --from-template tpl1 \
                 --hostname web1 --password 'Pass123' \
                 --ssh-key @/home/user/.ssh/id_rsa.pub \
                 --dns 8.8.8.8 --dns 8.8.4.4 \
                 --network nic=0,ip=192.168.1.100,netmask=255.255.255.0,gateway=192.168.1.1

  # From template with cloud-init (DHCP)
  goct vm.create --name web2 --cluster c1 --from-template tpl1 \
                 --hostname web2 --ssh-key "ssh-rsa AAAA..." \
                 --network nic=0,type=DHCP

  # From template with custom user_data
  goct vm.create --name web3 --cluster c1 --from-template tpl1 \
                 --user-data @/path/to/cloud-config.yaml

  # From scratch (no cloud-init)
  goct vm.create --name web4 --cluster c1 --vcpu 2 --memory 2048 \
                 --disk size=20g,bus=SCSI --nic vlan=vlan0,model=VIRTIO`,
		RunE: func(c *cobra.Command, _ []string) error {
			cli := client.From(c.Context())

			if fromTemplate != "" {
				// Build CloudInit spec if any cloud-init flag is set
				var cloudInit *adapter.CloudInitSpec
				if cloudInitHostname != "" || cloudInitPassword != "" || len(cloudInitSSHKey) > 0 || len(cloudInitDNSServers) > 0 || cloudInitUserData != "" || len(cloudInitNetworks) > 0 {
					cloudInit = &adapter.CloudInitSpec{
						Hostname:            cloudInitHostname,
						DefaultUserPassword: cloudInitPassword,
						DNSServers:          cloudInitDNSServers,
					}
					// Resolve @file paths for ssh-key and user-data
					for _, key := range cloudInitSSHKey {
						if resolved := resolveValue(key); resolved != "" {
							cloudInit.PublicKeys = append(cloudInit.PublicKeys, resolved)
						}
					}
					if cloudInitUserData != "" {
						cloudInit.UserData = resolveValue(cloudInitUserData)
					}
					if len(cloudInitNetworks) > 0 {
						cloudInit.Networks = parseCloudInitNetworks(cloudInitNetworks)
					}
				}

				ref, err := service.NewVM(cli).CreateFromTemplate(c.Context(), adapter.VMCreateFromTemplateSpec{
					TemplateID:  fromTemplate,
					Name:        name,
					ClusterID:   clusterID,
					VCPU:        vcpu,
					MemoryBytes: memoryMiB * 1024 * 1024,
					Firmware:    firmware,
					Description: description,
					IsFullCopy:  isFullCopy,
					NIC: adapter.NicConfig{
						Type:   nicType,
						Model:  nicModel,
						VlanID: nicVlan,
					},
					CloudInit: cloudInit,
				})
				if err != nil {
					return err
				}
				if ref.IsSync() {
					fmt.Fprintln(c.OutOrStdout(), "VM created from template (sync)")
					return nil
				}
				w := task.New(cli, task.Options{Out: c.OutOrStderr()})
				return w.Watch(c.Context(), ref.ID)
			}

			parsedDisks, err := parseDiskFlags(disks)
			if err != nil {
				return err
			}
			parsedNics, err := parseNicFlags(nics)
			if err != nil {
				return err
			}
			haPtr, err := parseHaFlag(ha)
			if err != nil {
				return err
			}

			ref, err := service.NewVM(cli).Create(c.Context(), adapter.VMCreateSpec{
				Name:        name,
				ClusterID:   clusterID,
				VCPU:        vcpu,
				MemoryBytes: memoryMiB * 1024 * 1024, // MiB → bytes
				Firmware:    firmware,
				Description: description,
				HA:          haPtr,
				Disks:       parsedDisks,
				Nics:        parsedNics,
			})
			if err != nil {
				return err
			}
			if ref.IsSync() {
				fmt.Fprintln(c.OutOrStdout(), "VM created (sync)")
				return nil
			}
			w := task.New(cli, task.Options{Out: c.OutOrStderr()})
			return w.Watch(c.Context(), ref.ID)
		},
	}
	c.Flags().StringVar(&name, "name", "", "VM name")
	c.Flags().StringVar(&clusterID, "cluster", "", "Target cluster ID")
	c.Flags().Int32Var(&vcpu, "vcpu", 0, "Number of vCPUs (0 = use template default)")
	c.Flags().Int64Var(&memoryMiB, "memory", 0, "Memory in MiB (0 = use template default)")
	c.Flags().StringVar(&firmware, "firmware", "", "Firmware: BIOS or UEFI (empty = use template default)")
	c.Flags().StringVar(&description, "description", "", "Description")
	c.Flags().StringVar(&ha, "ha", "", "Enable HA: true|false (empty = use cluster default)")
	c.Flags().StringArrayVar(&disks, "disk", nil, "Disk spec, repeatable: size=10g[,bus=SCSI][,name=diskN][,index=N][,boot=N][,iops=N]")
	c.Flags().StringArrayVar(&nics, "nic", nil, "NIC spec, repeatable: vlan=<id|name>[,model=VIRTIO][,type=VLAN|VPC]")
	c.Flags().StringVar(&fromTemplate, "from-template", "", "Template ID to create VM from")
	c.Flags().BoolVar(&isFullCopy, "full-copy", false, "Full copy when creating from template")
	c.Flags().StringVar(&nicType, "nic-type", "", "NIC type: VLAN or VPC (only for --from-template)")
	c.Flags().StringVar(&nicModel, "nic-model", "", "NIC model: E1000, SRIOV, VIRTIO (only for --from-template)")
	c.Flags().StringVar(&nicVlan, "nic-vlan", "", "VLAN ID (only for --from-template)")
	c.Flags().StringVar(&cloudInitHostname, "hostname", "", "Cloud-init hostname")
	c.Flags().StringVar(&cloudInitPassword, "password", "", "Default user password (cloud-init)")
	c.Flags().StringArrayVar(&cloudInitSSHKey, "ssh-key", nil, "SSH public key (literal or @/path/to/key.pub)")
	c.Flags().StringArrayVar(&cloudInitDNSServers, "dns", nil, "DNS nameserver (repeatable, e.g. --dns 8.8.8.8)")
	c.Flags().StringVar(&cloudInitUserData, "user-data", "", "Cloud-init user_data (@/path/to/file.yaml or literal YAML)")
	c.Flags().StringArrayVar(&cloudInitNetworks, "network", nil, "NIC config: nic=0[,ip=x][,netmask=x][,gateway=x][,route=10.0.0.0/8:192.168.1.1][,type=IPV4|DHCP]")
	return c
}

// parseDiskFlags 把 --disk k=v[,k=v...] 字符串数组解析为 DiskAddSpec 列表。
func parseDiskFlags(raw []string) ([]adapter.DiskAddSpec, error) {
	out := make([]adapter.DiskAddSpec, 0, len(raw))
	for _, s := range raw {
		kv, err := parseKVList(s)
		if err != nil {
			return nil, fmt.Errorf("invalid --disk %q: %w", s, err)
		}
		spec := adapter.DiskAddSpec{}
		if v, ok := kv["size"]; ok {
			n, err := parseSize(v)
			if err != nil {
				return nil, fmt.Errorf("invalid disk size %q: %w", v, err)
			}
			spec.SizeBytes = n
		}
		if spec.SizeBytes == 0 {
			return nil, fmt.Errorf("--disk %q: size required (e.g. size=10g)", s)
		}
		spec.Bus = strings.ToUpper(strings.TrimSpace(kv["bus"]))
		spec.Name = kv["name"]
		if v, ok := kv["index"]; ok {
			n, err := strconv.ParseInt(v, 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid disk index %q", v)
			}
			spec.Index = int32(n)
		}
		if v, ok := kv["boot"]; ok {
			n, err := strconv.ParseInt(v, 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid disk boot %q", v)
			}
			spec.Boot = int32(n)
		}
		if v, ok := kv["iops"]; ok {
			n, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid disk iops %q", v)
			}
			spec.IOPSMax = n
		}
		out = append(out, spec)
	}
	return out, nil
}

// parseNicFlags 把 --nic k=v[,k=v...] 字符串数组解析为 NicAddSpec 列表。
func parseNicFlags(raw []string) ([]adapter.NicAddSpec, error) {
	out := make([]adapter.NicAddSpec, 0, len(raw))
	for _, s := range raw {
		kv, err := parseKVList(s)
		if err != nil {
			return nil, fmt.Errorf("invalid --nic %q: %w", s, err)
		}
		out = append(out, adapter.NicAddSpec{
			Type:   strings.ToUpper(strings.TrimSpace(kv["type"])),
			Model:  strings.ToUpper(strings.TrimSpace(kv["model"])),
			VlanID: kv["vlan"],
		})
	}
	return out, nil
}

// parseHaFlag 解析 --ha 字符串：空 → nil（用默认）；true/false → 显式 *bool。
func parseHaFlag(s string) (*bool, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	if s == "" {
		return nil, nil
	}
	switch s {
	case "true", "yes", "1", "on":
		v := true
		return &v, nil
	case "false", "no", "0", "off":
		v := false
		return &v, nil
	}
	return nil, fmt.Errorf("invalid --ha %q (want true|false)", s)
}

// parseCloudInitNetworks parses --network flags into NicStaticConfig.
// Flag format: nic=0[,ip=x][,netmask=x][,gateway=x][,route=10.0.0.0/8:192.168.1.1][,type=IPV4|DHCP]
// Examples:
//   nic=0,ip=192.168.1.100,netmask=255.255.255.0,gateway=192.168.1.1,type=IPV4
//   nic=0,type=DHCP
//   nic=0,ip=192.168.1.100,netmask=255.255.255.0,gateway=192.168.1.1,route=10.0.0.0/8:192.168.1.1
func parseCloudInitNetworks(raw []string) []adapter.NicStaticConfig {
	out := make([]adapter.NicStaticConfig, 0, len(raw))
	for _, s := range raw {
		kv, err := parseKVList(s)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warn: invalid --network %q: %v\n", s, err)
			continue
		}
		cfg := adapter.NicStaticConfig{}

		// NIC index (required)
		if v, ok := kv["nic"]; ok {
			if n, err := strconv.ParseInt(v, 10, 32); err == nil {
				cfg.Index = int32(n)
			}
		}

		// Static IP config
		cfg.IP = kv["ip"]
		cfg.Netmask = kv["netmask"]
		cfg.Gateway = kv["gateway"]

		// Type: explicit IPV4/DHCP or auto-detect from ip presence
		if t := strings.ToUpper(kv["type"]); t != "" {
			cfg.Type = t
		} else if cfg.IP != "" && cfg.Netmask != "" {
			cfg.Type = "IPV4"
		} else {
			cfg.Type = "IPV4_DHCP"
		}

		// Custom static routes (format: network/netmask:gateway, repeatable with semicolon)
		if routeStr := kv["route"]; routeStr != "" {
			for _, route := range strings.Split(routeStr, ";") {
				route = strings.TrimSpace(route)
				if route == "" {
					continue
				}
				// Format: network/netmask:gateway (e.g. 10.0.0.0/8:192.168.1.1)
				parts := strings.Split(route, ":")
				if len(parts) != 2 {
					continue
				}
				netParts := strings.Split(parts[0], "/")
				if len(netParts) != 2 {
					continue
				}
				cfg.Routes = append(cfg.Routes, adapter.StaticRoute{
					Network: netParts[0],
					Netmask: netParts[1],
					Gateway: strings.TrimSpace(parts[1]),
				})
			}
		}

		out = append(out, cfg)
	}
	return out
}

// resolveValue resolves @file paths or returns the literal value.
// If s starts with "@", reads the file at the path and returns its content.
// Otherwise returns s as-is.
func resolveValue(s string) string {
	if strings.HasPrefix(s, "@") {
		path := strings.TrimPrefix(s, "@")
		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warn: failed to read file %q: %v\n", path, err)
			return ""
		}
		return string(data)
	}
	return s
}

// parseKVList 解析 "k1=v1,k2=v2" 字符串到 map；空字符串返回空 map。
func parseKVList(s string) (map[string]string, error) {
	out := map[string]string{}
	s = strings.TrimSpace(s)
	if s == "" {
		return out, nil
	}
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		eq := strings.IndexByte(part, '=')
		if eq <= 0 {
			return nil, fmt.Errorf("expected k=v, got %q", part)
		}
		k := strings.TrimSpace(strings.ToLower(part[:eq]))
		v := strings.TrimSpace(part[eq+1:])
		out[k] = v
	}
	return out, nil
}

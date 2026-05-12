package vm

import (
	"fmt"
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
		name         string
		clusterID    string
		vcpu         int32
		memoryMiB    int64
		firmware     string
		description  string
		fromTemplate string
		isFullCopy   bool
		nicType      string
		nicModel     string
		nicVlan      string
		ha           string
		disks        []string
		nics         []string
	)
	c := &cobra.Command{
		Use: "vm.create", Short: "Create a new VM", GroupID: groupID,
		Long: `Create a new virtual machine.

Two modes:
  1) From template (--from-template): use existing content library template.
  2) From scratch: build a VM with --disk / --nic flags.

Examples:
  goct vm.create --name web1 --cluster c1 --vcpu 2 --memory 2048 \
                 --disk size=20g,bus=SCSI --nic vlan=vlan0,model=VIRTIO

  goct vm.create --name win1 --cluster c1 --vcpu 4 --memory 4096 \
                 --firmware UEFI --ha=true \
                 --disk size=60g,bus=SCSI --disk size=200g,bus=SCSI,name=data \
                 --nic vlan=vlan0

  goct vm.create --name from-tpl --cluster c1 --from-template tpl1 \
                 --nic-vlan vlan0`,
		RunE: func(c *cobra.Command, _ []string) error {
			cli := client.From(c.Context())

			if fromTemplate != "" {
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

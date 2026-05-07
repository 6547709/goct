package output

import (
	"fmt"
	"strings"

	"github.com/6547709/goct/pkg/adapter"
)

// VMListColumns 是 vm.ls 命令的表格列定义。
var VMListColumns = []Column{
	{Header: "ID", Get: func(v any) string {
		return v.(adapter.VM).ID
	}},
	{Header: "NAME", Get: func(v any) string {
		return v.(adapter.VM).Name
	}},
	{Header: "STATUS", Get: func(v any) string {
		return v.(adapter.VM).Status
	}},
	{Header: "VCPU", Get: func(v any) string {
		return fmt.Sprintf("%d", v.(adapter.VM).VCPU)
	}},
	{Header: "MEMORY", Get: func(v any) string {
		return HumanBytes(v.(adapter.VM).MemoryBytes)
	}},
	{Header: "IPS", Get: func(v any) string {
		ips := v.(adapter.VM).IPs
		if len(ips) == 0 {
			return "-"
		}
		return strings.Join(ips, ", ")
	}},
	{Header: "HOST", Get: func(v any) string {
		h := v.(adapter.VM).HostName
		if h == "" {
			return "-"
		}
		return h
	}},
}

// TaskListColumns 是 task.ls 的表格列定义。
var TaskListColumns = []Column{
	{Header: "ID", Get: func(v any) string { return v.(adapter.Task).ID }},
	{Header: "STATUS", Get: func(v any) string { return v.(adapter.Task).Status }},
	{Header: "PROGRESS", Get: func(v any) string {
		p := v.(adapter.Task).Progress
		if p == 0 {
			return "-"
		}
		return fmt.Sprintf("%d%%", p)
	}},
	{Header: "DESCRIPTION", Get: func(v any) string { return v.(adapter.Task).Description }},
	{Header: "CREATED", Get: func(v any) string { return v.(adapter.Task).CreatedAt }},
	{Header: "STARTED", Get: func(v any) string {
		s := v.(adapter.Task).StartedAt
		if s == "" {
			return "-"
		}
		return s
	}},
	{Header: "FINISHED", Get: func(v any) string {
		s := v.(adapter.Task).FinishedAt
		if s == "" {
			return "-"
		}
		return s
	}},
}

// AlertListColumns 是 alert.ls 的表格列定义。
var AlertListColumns = []Column{
	{Header: "ID", Get: func(v any) string { return v.(adapter.Alert).ID }},
	{Header: "SEVERITY", Get: func(v any) string { return v.(adapter.Alert).Severity }},
	{Header: "MESSAGE", Get: func(v any) string { return v.(adapter.Alert).Message }},
	{Header: "CAUSE", Get: func(v any) string { return v.(adapter.Alert).Cause }},
}

// UserListColumns 是 user.ls 的表格列定义。
var UserListColumns = []Column{
	{Header: "ID", Get: func(v any) string { return v.(adapter.User).ID }},
	{Header: "NAME", Get: func(v any) string { return v.(adapter.User).Name }},
	{Header: "USERNAME", Get: func(v any) string { return v.(adapter.User).Username }},
	{Header: "ROLE", Get: func(v any) string { return v.(adapter.User).Role }},
	{Header: "SOURCE", Get: func(v any) string { return v.(adapter.User).Source }},
}

// NetworkListColumns 是 network.ls 的表格列定义。
var NetworkListColumns = []Column{
	{Header: "ID", Get: func(v any) string { return v.(adapter.Network).ID }},
	{Header: "NAME", Get: func(v any) string { return v.(adapter.Network).Name }},
	{Header: "TYPE", Get: func(v any) string { return v.(adapter.Network).Type }},
	{Header: "CLUSTER", Get: func(v any) string { return v.(adapter.Network).ClusterID }},
}

// VLANListColumns 是 vlan.ls 的表格列定义。
var VLANListColumns = []Column{
	{Header: "ID", Get: func(v any) string { return v.(adapter.VLAN).ID }},
	{Header: "NAME", Get: func(v any) string { return v.(adapter.VLAN).Name }},
	{Header: "VLAN TAG", Get: func(v any) string { return fmt.Sprintf("%d", v.(adapter.VLAN).VlanTag) }},
	{Header: "TYPE", Get: func(v any) string { return v.(adapter.VLAN).Type }},
	{Header: "VDS", Get: func(v any) string { return v.(adapter.VLAN).VdsID }},
}

// ClusterListColumns 是 cluster.ls 的表格列定义。
var ClusterListColumns = []Column{
	{Header: "ID", Get: func(v any) string { return v.(adapter.Cluster).ID }},
	{Header: "NAME", Get: func(v any) string { return v.(adapter.Cluster).Name }},
	{Header: "CPU CORES", Get: func(v any) string { return fmt.Sprintf("%d", v.(adapter.Cluster).TotalCPUCores) }},
	{Header: "MEMORY", Get: func(v any) string { return HumanBytes(v.(adapter.Cluster).TotalMemoryBytes) }},
	{Header: "STORAGE", Get: func(v any) string { return HumanBytes(v.(adapter.Cluster).TotalDataCapacity) }},
	{Header: "VMs", Get: func(v any) string { return fmt.Sprintf("%d", v.(adapter.Cluster).RunningVMs) }},
}

// DatastoreListColumns 是 datastore.ls 的表格列定义。
var DatastoreListColumns = []Column{
	{Header: "ID", Get: func(v any) string { return v.(adapter.Datastore).ID }},
	{Header: "NAME", Get: func(v any) string { return v.(adapter.Datastore).Name }},
	{Header: "TYPE", Get: func(v any) string { return v.(adapter.Datastore).Type }},
	{Header: "CLUSTER", Get: func(v any) string { return v.(adapter.Datastore).ClusterID }},
}

// DiskPoolListColumns 是 storage.pool.ls 的表格列定义（超融合存储池）。
var DiskPoolListColumns = []Column{
	{Header: "ID", Get: func(v any) string { return v.(adapter.DiskPool).ID }},
	{Header: "HOST", Get: func(v any) string { return v.(adapter.DiskPool).HostName }},
	{Header: "STATUS", Get: func(v any) string { return v.(adapter.DiskPool).Status }},
	{Header: "USE", Get: func(v any) string { return v.(adapter.DiskPool).UseState }},
	{Header: "CAPACITY", Get: func(v any) string {
		return HumanBytes(v.(adapter.DiskPool).TotalDataBytes)
	}},
	{Header: "USED", Get: func(v any) string {
		return HumanBytes(v.(adapter.DiskPool).UsedDataBytes)
	}},
	{Header: "CACHE", Get: func(v any) string {
		return HumanBytes(v.(adapter.DiskPool).TotalCacheBytes)
	}},
	{Header: "HDD", Get: func(v any) string { return fmt.Sprintf("%d", v.(adapter.DiskPool).HddCount) }},
	{Header: "NVME", Get: func(v any) string { return fmt.Sprintf("%d", v.(adapter.DiskPool).NvmeCount) }},
	{Header: "SATA", Get: func(v any) string { return fmt.Sprintf("%d", v.(adapter.DiskPool).SataCount) }},
}

// DiskListColumns 是 datastore.disk.ls 的表格列定义。
var DiskListColumns = []Column{
	{Header: "ID", Get: func(v any) string { return v.(adapter.Disk).ID }},
	{Header: "NAME", Get: func(v any) string { return v.(adapter.Disk).Name }},
	{Header: "TYPE", Get: func(v any) string { return v.(adapter.Disk).Type }},
	{Header: "SIZE", Get: func(v any) string { return HumanBytes(v.(adapter.Disk).SizeBytes) }},
	{Header: "HOST", Get: func(v any) string { return v.(adapter.Disk).HostName }},
	{Header: "PATH", Get: func(v any) string { return v.(adapter.Disk).Path }},
}

// HostListColumns 是 host.ls 的表格列定义。
var HostListColumns = []Column{
	{Header: "ID", Get: func(v any) string { return v.(adapter.Host).ID }},
	{Header: "NAME", Get: func(v any) string { return v.(adapter.Host).Name }},
	{Header: "STATUS", Get: func(v any) string { return v.(adapter.Host).Status }},
	{Header: "MGMT IP", Get: func(v any) string { return v.(adapter.Host).ManagementIP }},
	{Header: "MEMORY", Get: func(v any) string { return HumanBytes(v.(adapter.Host).TotalMemoryBytes) }},
	{Header: "VMs", Get: func(v any) string { return fmt.Sprintf("%d", v.(adapter.Host).RunningVMs) }},
}

// HostInfoRows 返回 Host 的 key-value 行。
func HostInfoRows(h adapter.Host) [][]string {
	return [][]string{
		{"ID", h.ID},
		{"Name", h.Name},
		{"Status", h.Status},
		{"Management IP", h.ManagementIP},
		{"Data IP", h.DataIP},
		{"CPU Model", h.CPUModel},
		{"Total Memory", HumanBytes(h.TotalMemoryBytes)},
		{"Running VMs", fmt.Sprintf("%d", h.RunningVMs)},
		{"Cluster", h.ClusterID},
	}
}

// SnapshotListColumns 是 vm.snapshot.ls 的表格列定义。
var SnapshotListColumns = []Column{
	{Header: "ID", Get: func(v any) string { return v.(adapter.Snapshot).ID }},
	{Header: "NAME", Get: func(v any) string { return v.(adapter.Snapshot).Name }},
	{Header: "VM", Get: func(v any) string { return v.(adapter.Snapshot).VMID }},
	{Header: "CREATED", Get: func(v any) string {
		c := v.(adapter.Snapshot).CreatedAt
		if c == "" {
			return "-"
		}
		return c
	}},
	{Header: "DESCRIPTION", Get: func(v any) string { return v.(adapter.Snapshot).Description }},
}

// VMDiskListColumns 是 vm.disk.ls 的表格列定义。
var VMDiskListColumns = []Column{
	{Header: "ID", Get: func(v any) string { return v.(adapter.VMDisk).ID }},
	{Header: "VM", Get: func(v any) string { return v.(adapter.VMDisk).VMID }},
	{Header: "BOOT", Get: func(v any) string { return fmt.Sprintf("%d", v.(adapter.VMDisk).Boot) }},
	{Header: "BUS", Get: func(v any) string { return v.(adapter.VMDisk).Bus }},
	{Header: "TYPE", Get: func(v any) string { return v.(adapter.VMDisk).Type }},
	{Header: "VOLUME", Get: func(v any) string { return v.(adapter.VMDisk).VolumeName }},
	{Header: "SIZE", Get: func(v any) string { return HumanBytes(uint64(v.(adapter.VMDisk).VolumeSize)) }},
	{Header: "MAX BW", Get: func(v any) string {
		if v.(adapter.VMDisk).MaxBandwidth == nil {
			return "-"
		}
		return fmt.Sprintf("%d", *v.(adapter.VMDisk).MaxBandwidth)
	}},
	{Header: "MAX IOPS", Get: func(v any) string {
		if v.(adapter.VMDisk).MaxIops == nil {
			return "-"
		}
		return fmt.Sprintf("%d", *v.(adapter.VMDisk).MaxIops)
	}},
}

// VMNicListColumns 是 vm.nic.ls 的表格列定义。
var VMNicListColumns = []Column{
	{Header: "ID", Get: func(v any) string { return v.(adapter.VMNic).ID }},
	{Header: "VM", Get: func(v any) string { return v.(adapter.VMNic).VMID }},
	{Header: "LOCAL ID", Get: func(v any) string { return v.(adapter.VMNic).LocalID }},
	{Header: "MAC", Get: func(v any) string { return v.(adapter.VMNic).MacAddress }},
	{Header: "MODEL", Get: func(v any) string { return v.(adapter.VMNic).Model }},
	{Header: "TYPE", Get: func(v any) string { return v.(adapter.VMNic).Type }},
	{Header: "VLAN", Get: func(v any) string { return v.(adapter.VMNic).VlanID }},
	{Header: "IP", Get: func(v any) string { return v.(adapter.VMNic).IPAddress }},
	{Header: "GATEWAY", Get: func(v any) string { return v.(adapter.VMNic).Gateway }},
	{Header: "ENABLED", Get: func(v any) string {
		if v.(adapter.VMNic).Enabled {
			return "true"
		}
		return "false"
	}},
}

// GpuDeviceListColumns 是 vm.gpu.ls 的表格列定义。
var GpuDeviceListColumns = []Column{
	{Header: "ID", Get: func(v any) string { return v.(adapter.GpuDevice).ID }},
	{Header: "NAME", Get: func(v any) string { return v.(adapter.GpuDevice).Name }},
}

var VMInfoColumns = []Column{
	{Header: "FIELD", Get: func(_ any) string { return "" }},
	{Header: "VALUE", Get: func(_ any) string { return "" }},
}

// VMInfoRows 返回 VM 的 key-value 行用于 info 展示（~22 行，操作友好）。
func VMInfoRows(v adapter.VM) [][]string {
	ips := v.IPs
	if len(ips) == 0 {
		ips = []string{"-"}
	}
	haStr := "false"
	if v.Ha {
		haStr = "true"
	}
	protectedStr := "false"
	if v.Protected {
		protectedStr = "true"
	}
	inBinStr := "false"
	if v.InRecycleBin {
		inBinStr = "true"
	}
	desc := v.Description
	if desc == "" {
		desc = "-"
	}
	dns := v.DNSServers
	if dns == "" {
		dns = "-"
	}
	hostname := v.Hostname
	if hostname == "" {
		hostname = "-"
	}
	clusterVal := v.ClusterID
	if v.ClusterName != "" {
		clusterVal = fmt.Sprintf("%s (%s)", v.ClusterName, v.ClusterID)
	}
	hostVal := v.HostName
	if hostVal == "" {
		hostVal = "-"
	} else if v.HostIP != "" {
		hostVal = fmt.Sprintf("%s [%s]", v.HostName, v.HostIP)
	}
	return [][]string{
		{"Name", v.Name},
		{"ID", v.ID},
		{"Status", v.Status},
		{"VCPU", fmt.Sprintf("%d", v.VCPU)},
		{"Memory", HumanBytes(v.MemoryBytes)},
		{"Firmware", v.Firmware},
		{"Guest OS", v.GuestOS},
		{"HA", haStr},
		{"IPs", strings.Join(ips, ", ")},
		{"Hostname", hostname},
		{"DNS Servers", dns},
		{"VMTools", v.VMToolsStatus},
		{"VMTools Version", v.VMToolsVersion},
		{"CPU Model", v.CPUModel},
		{"Host", hostVal},
		{"Cluster", clusterVal},
		{"Disks", fmt.Sprintf("%d", v.DiskCount)},
		{"NICs", fmt.Sprintf("%d", v.NicCount)},
		{"Provisioned", HumanBytes(v.ProvisionedBytes)},
		{"Used", HumanBytes(v.UsedBytes)},
		{"In Recycle Bin", inBinStr},
		{"Protected", protectedStr},
		{"Created", v.CreatedAt},
		{"Description", desc},
	}
}

// VMDetailRows 返回 VM 的详细 key-value 行用于 info --detail 展示（包含额外字段）。
func VMDetailRows(v adapter.VM) [][]string {
	// Start with the base info rows
	rows := VMInfoRows(v)

	// Helper to format optional float
	formatFloat := func(f float64) string {
		if f == 0 {
			return "-"
		}
		return fmt.Sprintf("%.1f%%", f)
	}

	// Helper for bool string
	formatBool := func(b bool) string {
		if b {
			return "true"
		}
		return "false"
	}

	// Append detail fields
	nestedVirtStr := formatBool(v.NestedVirt)
	if !v.NestedVirt {
		nestedVirtStr = "false"
	}

	cloudInitStr := "false"
	if v.CloudInit {
		cloudInitStr = "true"
	}

	haPriority := v.HaPriority
	if haPriority == "" {
		haPriority = "-"
	}

	biosUUID := v.BiosUUID
	if biosUUID == "" {
		biosUUID = "-"
	}

	labelsStr := "-"
	if len(v.Labels) > 0 {
		labelsStr = strings.Join(v.Labels, ", ")
	}

	gpuStr := "-"
	if len(v.GpuDevices) > 0 {
		var names []string
		for _, d := range v.GpuDevices {
			names = append(names, d.Name)
		}
		gpuStr = strings.Join(names, ", ")
	}

	usbStr := "-"
	if len(v.UsbDevices) > 0 {
		var names []string
		for _, d := range v.UsbDevices {
			names = append(names, d.Name)
		}
		usbStr = strings.Join(names, ", ")
	}

	videoType := v.VideoType
	if videoType == "" {
		videoType = "-"
	}

	detailRows := [][]string{
		{"---", "---"},
		{"Bios UUID", biosUUID},
		{"CPU Usage", formatFloat(v.CPUUsage)},
		{"Memory Usage", formatFloat(v.MemoryUsage)},
		{"Guest Size Usage", formatFloat(v.GuestSizeUsage)},
		{"Guest Used Size", HumanBytes(uint64(v.GuestUsedSize))},
		{"Logical Size", HumanBytes(uint64(v.LogicalSizeBytes))},
		{"Nested Virt", nestedVirtStr},
		{"CloudInit", cloudInitStr},
		{"HA Priority", haPriority},
		{"Video Type", videoType},
		{"GPU Devices", gpuStr},
		{"USB Devices", usbStr},
		{"Labels", labelsStr},
	}
	return append(rows, detailRows...)
}

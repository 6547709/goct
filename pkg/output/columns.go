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
var VMInfoColumns = []Column{
	{Header: "FIELD", Get: func(_ any) string { return "" }},
	{Header: "VALUE", Get: func(_ any) string { return "" }},
}

// VMInfoRows 返回 VM 的 key-value 行用于 info 展示。
func VMInfoRows(v adapter.VM) [][]string {
	return [][]string{
		{"ID", v.ID},
		{"Name", v.Name},
		{"Status", v.Status},
		{"VCPU", fmt.Sprintf("%d", v.VCPU)},
		{"Memory", HumanBytes(v.MemoryBytes)},
		{"IPs", strings.Join(v.IPs, ", ")},
		{"Host", v.HostName},
		{"Cluster", v.ClusterID},
		{"Description", v.Description},
	}
}

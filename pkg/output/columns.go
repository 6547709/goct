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

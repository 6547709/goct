// Package find 实现 govc 风格的统一查找命令 `goct find`。
//
// 设计目标：
//   - 提供"按类型 + 名称模式 + 数量限制"的统一入口，避免用户为每种资源记忆 ls 命令；
//   - 输出统一为 "TYPE  ID  NAME"（或 --id-only / --json），便于和 xargs 拼接。
//
// 与 govc find 的差异：
//   - govc find 走 inventory path（/DC/vm/...），CloudTower 没有这个概念，goct find 直接在
//     全集群内按资源类型扫描（可用 --cluster 过滤）。
//   - 类型名对齐 govc：m=VM、h=Host、c=Cluster、d=Datastore、n=Network、p=Pool（CloudTower 暂无），
//     额外补充 CloudTower 独有：vlan、folder、pg（placement group）、template、label、user、task、alert。
package find

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/6547709/goct/pkg/adapter"
	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/service"
	"github.com/spf13/cobra"
)

const groupID = "system"

// Register 把 find 子命令挂到 root。
func Register(root *cobra.Command) {
	root.AddCommand(newFind())
}

// resourceType 描述一种可被 find 枚举的资源。
type resourceType struct {
	short    string
	longName string
	list     func(ctx context.Context, cli adapter.Client, name, clusterID string) ([]row, error)
}

// row 是 find 输出的统一行结构。
type row struct {
	Type string `json:"type"`
	ID   string `json:"id"`
	Name string `json:"name"`
}

// types 列出所有 find 支持的类型。
// short 与 govc 兼容：m=VirtualMachine、h=HostSystem、c=ClusterComputeResource、d=Datastore、n=Network。
// CloudTower 独有的资源用全名（如 vlan / folder / pg）。
func types() []resourceType {
	return []resourceType{
		{"m", "vm", listVMs},
		{"h", "host", listHosts},
		{"c", "cluster", listClusters},
		{"d", "datastore", listDatastores},
		{"n", "network", listNetworks},
		{"v", "vlan", listVLANs},
		{"f", "folder", listFolders},
		{"g", "pg", listPlacementGroups},
		{"t", "template", listTemplates},
		{"l", "label", listLabels},
		{"u", "user", listUsers},
		{"a", "alert", listAlerts},
	}
}

func newFind() *cobra.Command {
	var (
		typeFilter string
		nameFilter string
		clusterID  string
		idOnly     bool
		formatJSON bool
		limit      int32
	)
	c := &cobra.Command{
		Use:     "find",
		Short:   "Find resources across the inventory (govc-style)",
		GroupID: groupID,
		Long: `Find resources of a given type matching a name pattern.

Type shorthand (compatible with govc):
  m   VM (virtual machine)
  h   Host
  c   Cluster
  d   Datastore
  n   Network
  v   VLAN                 (CloudTower)
  f   VM Folder            (CloudTower)
  g   VM Placement Group   (CloudTower)
  t   Template             (CloudTower content library)
  l   Label                (CloudTower)
  u   User
  a   Alert

Long names are also accepted (vm, host, cluster, datastore, network, vlan,
folder, pg, template, label, user, alert).

If --type is omitted, find scans every type.

Examples:
  goct find --type m --name web           # all VMs whose name contains "web"
  goct find --type h --cluster Cluster01  # all hosts in cluster c1
  goct find --name prod                   # all resources whose name contains "prod" (scans all types)
  goct find --type m --id-only | xargs -I{} goct vm.power.on {}`,
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			selected, err := resolveTypes(typeFilter)
			if err != nil {
				return err
			}

			var rows []row
			for _, t := range selected {
				items, err := t.list(c.Context(), cli, nameFilter, clusterID)
				if err != nil {
					// 单一类型的失败不应该 abort 整个 find（用户可能没权限看某类资源），
					// 仅打 warning 到 stderr。
					fmt.Fprintf(c.ErrOrStderr(), "warn: list %s: %v\n", t.longName, err)
					continue
				}
				rows = append(rows, items...)
				if limit > 0 && int32(len(rows)) >= limit {
					rows = rows[:limit]
					break
				}
			}
			sort.SliceStable(rows, func(i, j int) bool {
				if rows[i].Type != rows[j].Type {
					return rows[i].Type < rows[j].Type
				}
				return rows[i].Name < rows[j].Name
			})

			if idOnly {
				for _, r := range rows {
					fmt.Fprintln(c.OutOrStdout(), r.ID)
				}
				return nil
			}
			if formatJSON {
				return writeJSON(c.OutOrStdout(), rows)
			}
			return writeTable(c.OutOrStdout(), rows)
		},
	}
	c.Flags().StringVar(&typeFilter, "type", "", "Resource type filter: m|h|c|d|n|v|f|g|t|l|u|a or full name (default: all)")
	c.Flags().StringVar(&nameFilter, "name", "", "Filter by name (substring match)")
	c.Flags().StringVar(&clusterID, "cluster", "", "Restrict to this cluster (ID or name); ignored for cluster-less resource types")
	c.Flags().BoolVar(&idOnly, "id-only", false, "Print only IDs (one per line)")
	c.Flags().BoolVar(&formatJSON, "json", false, "JSON output")
	c.Flags().Int32Var(&limit, "limit", 0, "Maximum results (0 = unlimited)")
	return c
}

// resolveTypes 把 --type 字符串解析成要扫描的资源类型列表。
// 空值表示全部。多类型可用逗号分隔。
func resolveTypes(s string) ([]resourceType, error) {
	all := types()
	if s == "" {
		return all, nil
	}
	wanted := map[string]bool{}
	for _, part := range strings.Split(s, ",") {
		p := strings.TrimSpace(strings.ToLower(part))
		if p == "" {
			continue
		}
		wanted[p] = true
	}
	out := make([]resourceType, 0, len(wanted))
	matched := map[string]bool{}
	for _, t := range all {
		if wanted[t.short] || wanted[t.longName] {
			out = append(out, t)
			matched[t.short] = true
			matched[t.longName] = true
		}
	}
	for k := range wanted {
		if !matched[k] {
			return nil, fmt.Errorf("unknown type %q (want one of: m,h,c,d,n,v,f,g,t,l,u,a)", k)
		}
	}
	return out, nil
}

// ---------- per-type listers ----------

func listVMs(ctx context.Context, cli adapter.Client, name, clusterID string) ([]row, error) {
	opts := adapter.ListOpts{NameContains: name, ClusterID: clusterID}
	items, err := service.NewVM(cli).List(ctx, opts)
	if err != nil {
		return nil, err
	}
	out := make([]row, 0, len(items))
	for _, v := range items {
		out = append(out, row{Type: "vm", ID: v.ID, Name: v.Name})
	}
	return out, nil
}

func listHosts(ctx context.Context, cli adapter.Client, name, clusterID string) ([]row, error) {
	items, err := cli.ListHosts(ctx, adapter.ListOpts{NameContains: name, ClusterID: clusterID})
	if err != nil {
		return nil, err
	}
	out := make([]row, 0, len(items))
	for _, h := range items {
		out = append(out, row{Type: "host", ID: h.ID, Name: h.Name})
	}
	return out, nil
}

func listClusters(ctx context.Context, cli adapter.Client, name, _ string) ([]row, error) {
	items, err := cli.ListClusters(ctx, adapter.ListOpts{NameContains: name})
	if err != nil {
		return nil, err
	}
	out := make([]row, 0, len(items))
	for _, c := range items {
		out = append(out, row{Type: "cluster", ID: c.ID, Name: c.Name})
	}
	return out, nil
}

func listDatastores(ctx context.Context, cli adapter.Client, name, clusterID string) ([]row, error) {
	items, err := cli.ListDatastores(ctx, adapter.ListOpts{NameContains: name, ClusterID: clusterID})
	if err != nil {
		return nil, err
	}
	out := make([]row, 0, len(items))
	for _, d := range items {
		out = append(out, row{Type: "datastore", ID: d.ID, Name: d.Name})
	}
	return out, nil
}

func listNetworks(ctx context.Context, cli adapter.Client, name, clusterID string) ([]row, error) {
	items, err := cli.ListNetworks(ctx, adapter.ListOpts{NameContains: name, ClusterID: clusterID})
	if err != nil {
		return nil, err
	}
	out := make([]row, 0, len(items))
	for _, n := range items {
		out = append(out, row{Type: "network", ID: n.ID, Name: n.Name})
	}
	return out, nil
}

func listVLANs(ctx context.Context, cli adapter.Client, name, clusterID string) ([]row, error) {
	items, err := cli.ListVLANs(ctx, adapter.ListOpts{NameContains: name, ClusterID: clusterID})
	if err != nil {
		return nil, err
	}
	out := make([]row, 0, len(items))
	for _, v := range items {
		out = append(out, row{Type: "vlan", ID: v.ID, Name: v.Name})
	}
	return out, nil
}

func listFolders(ctx context.Context, cli adapter.Client, name, clusterID string) ([]row, error) {
	items, err := cli.ListVMFolders(ctx, adapter.ListOpts{NameContains: name, ClusterID: clusterID})
	if err != nil {
		return nil, err
	}
	out := make([]row, 0, len(items))
	for _, f := range items {
		out = append(out, row{Type: "folder", ID: f.ID, Name: f.Name})
	}
	return out, nil
}

func listPlacementGroups(ctx context.Context, cli adapter.Client, name, clusterID string) ([]row, error) {
	items, err := cli.ListVMPlacementGroups(ctx, adapter.ListOpts{NameContains: name, ClusterID: clusterID})
	if err != nil {
		return nil, err
	}
	out := make([]row, 0, len(items))
	for _, g := range items {
		out = append(out, row{Type: "pg", ID: g.ID, Name: g.Name})
	}
	return out, nil
}

func listTemplates(ctx context.Context, cli adapter.Client, name, _ string) ([]row, error) {
	items, err := cli.ListContentLibraryTemplates(ctx, adapter.ListOpts{NameContains: name})
	if err != nil {
		return nil, err
	}
	out := make([]row, 0, len(items))
	for _, t := range items {
		out = append(out, row{Type: "template", ID: t.ID, Name: t.Name})
	}
	return out, nil
}

func listLabels(ctx context.Context, cli adapter.Client, name, _ string) ([]row, error) {
	items, err := cli.ListLabels(ctx, adapter.ListOpts{NameContains: name})
	if err != nil {
		return nil, err
	}
	out := make([]row, 0, len(items))
	for _, l := range items {
		// Label uses Key/Value as identity; expose Key as Name for find.
		display := l.Key
		if l.Value != "" {
			display = l.Key + "=" + l.Value
		}
		out = append(out, row{Type: "label", ID: l.ID, Name: display})
	}
	return out, nil
}

func listUsers(ctx context.Context, cli adapter.Client, name, _ string) ([]row, error) {
	items, err := cli.ListUsers(ctx, adapter.ListOpts{NameContains: name})
	if err != nil {
		return nil, err
	}
	out := make([]row, 0, len(items))
	for _, u := range items {
		out = append(out, row{Type: "user", ID: u.ID, Name: u.Name})
	}
	return out, nil
}

func listAlerts(ctx context.Context, cli adapter.Client, name, _ string) ([]row, error) {
	items, err := cli.ListAlerts(ctx, adapter.ListOpts{NameContains: name})
	if err != nil {
		return nil, err
	}
	out := make([]row, 0, len(items))
	for _, a := range items {
		// Alert has no Name; use first 60 chars of Message as display.
		msg := a.Message
		if len(msg) > 60 {
			msg = msg[:60] + "..."
		}
		out = append(out, row{Type: "alert", ID: a.ID, Name: msg})
	}
	return out, nil
}

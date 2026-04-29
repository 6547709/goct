// Package flags 定义可被多个命令复用的通用 flag 结构体。
//
// ConnectionFlags / OutputFlags / SearchFlags 通过 Register(cmd) 把字段
// 绑定到 cobra 命令；其中 ConnectionFlags 与 OutputFlags 注册到
// PersistentFlags（向所有子命令传播），SearchFlags 注册到本地 Flags
// （仅当前命令可见）。
package flags

import "github.com/spf13/cobra"

// ConnectionFlags 描述 CloudTower 连接所需的全局参数。
// 这些值最终会与环境变量与配置文件合并（见 pkg/config）。
type ConnectionFlags struct {
	URL      string // 例：https://tower.example.com
	Username string
	Password string
	Cluster  string // 默认集群（按 ID 或 Name 解析）
	Insecure bool   // 跳过 TLS 校验，仅用于自签名内网环境
	Source   string // 登录源：local|ldap|sso|authn，默认 local
}

// Register 把 ConnectionFlags 的字段注册到 cmd 的 PersistentFlags。
func (f *ConnectionFlags) Register(c *cobra.Command) {
	pf := c.PersistentFlags()
	pf.StringVar(&f.URL, "url", "", "CloudTower endpoint URL (env: GOCT_URL)")
	pf.StringVar(&f.Username, "username", "", "Login username (env: GOCT_USERNAME)")
	pf.StringVar(&f.Password, "password", "", "Login password (env: GOCT_PASSWORD)")
	pf.BoolVar(&f.Insecure, "insecure", false, "Skip TLS certificate verification (env: GOCT_INSECURE)")
	pf.StringVar(&f.Cluster, "cluster", "", "Default cluster ID or name (env: GOCT_CLUSTER)")
	pf.StringVar(&f.Source, "source", "", "Login source: local|ldap|sso|authn (env: GOCT_SOURCE, default local)")
}

// OutputFlags 控制命令输出渲染方式。
type OutputFlags struct {
	Format string // table | json
}

// Register 把 --format 注册到 cmd 的 PersistentFlags，默认 table。
func (f *OutputFlags) Register(c *cobra.Command) {
	if f.Format == "" {
		f.Format = "table"
	}
	c.PersistentFlags().StringVar(&f.Format, "format", f.Format, "Output format: table|json")
}

// SearchFlags 是 list 类命令通用的过滤与分页 flag。
type SearchFlags struct {
	Name  string // 子串匹配
	Limit int32  // 0 表示不限
	Skip  int32
}

// Register 把 --name/--limit/--skip 注册到当前 cmd 的本地 Flags。
func (f *SearchFlags) Register(c *cobra.Command) {
	c.Flags().StringVar(&f.Name, "name", "", "Filter by name (substring match)")
	c.Flags().Int32Var(&f.Limit, "limit", 0, "Limit number of results (0 = no limit)")
	c.Flags().Int32Var(&f.Skip, "skip", 0, "Skip N results (pagination)")
}

// Package cmd 定义 goct 的所有 CLI 命令。
// root.go 提供根命令、全局 flag、PersistentPreRunE 构建 client。
package cmd

import (
	"context"
	"net/http"
	"os"

	"github.com/6547709/goct/cmd/alert"
	"github.com/6547709/goct/cmd/alert_rule"
	"github.com/6547709/goct/cmd/application"
	"github.com/6547709/goct/cmd/cluster"
	"github.com/6547709/goct/cmd/cluster_settings"
	"github.com/6547709/goct/cmd/content_library_image"
	"github.com/6547709/goct/cmd/datastore"
	"github.com/6547709/goct/cmd/deploy"
	"github.com/6547709/goct/cmd/events"
	"github.com/6547709/goct/cmd/find"
	"github.com/6547709/goct/cmd/host"
	"github.com/6547709/goct/cmd/license"
	"github.com/6547709/goct/cmd/metrics"
	"github.com/6547709/goct/cmd/network"
	"github.com/6547709/goct/cmd/ntp"
	"github.com/6547709/goct/cmd/system"
	"github.com/6547709/goct/cmd/task"
	"github.com/6547709/goct/cmd/user"
	"github.com/6547709/goct/cmd/vlan"
	vmcmd "github.com/6547709/goct/cmd/vm"
	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/config"
	"github.com/6547709/goct/pkg/debug"
	"github.com/6547709/goct/pkg/flags"
	"github.com/spf13/cobra"
)

var (
	connFlags   flags.ConnectionFlags
	outputFlags flags.OutputFlags
	debugFlags  debug.DebugFlags
)

// rootCmd 是 goct 的根命令，子命令通过各子包注册。
var rootCmd = &cobra.Command{
	Use:           "goct",
	Short:         "CloudTower CLI (govc-style)",
	Long:          "goct is a govc-style command-line client for SmartX CloudTower.",
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRunE: func(c *cobra.Command, _ []string) error {
		// 初始化分级日志（GOCT_LOG=TRACE|DEBUG|INFO|WARN|ERROR）
		debug.Init()

		// 子命令可通过 Annotations["nologin"]="true" 跳过登录步骤。
		// 典型场景：version / completion / 纯本地的 session.ls 等。
		if v, ok := c.Annotations["nologin"]; ok && v == "true" {
			return nil
		}

		opts := debugFlags.Resolve()

		// 如果启用了 trace，创建 TraceRoundTripper
		var traceTransport http.RoundTripper
		if opts.TraceLevel > debug.TraceLevelOff {
			traceTransport = debug.NewTraceRoundTripper(nil, opts.TraceLevel, os.Stderr)
		}

		cfg, err := config.Resolve(config.Override{
			URL:         connFlags.URL,
			Username:    connFlags.Username,
			Password:    connFlags.Password,
			Cluster:     connFlags.Cluster,
			Source:      connFlags.Source,
			Insecure:    connFlags.Insecure,
			InsecureSet: c.Flags().Changed("insecure"),
		})
		if err != nil {
			return err
		}
		cli, err := client.New(c.Context(), cfg, traceTransport)
		if err != nil {
			return err
		}
		c.SetContext(client.With(c.Context(), cli))
		return nil
	},
}

func init() {
	rootCmd.SetContext(context.Background())
	connFlags.Register(rootCmd)
	outputFlags.Register(rootCmd)
	debugFlags.Register(rootCmd)

	rootCmd.AddGroup(
		&cobra.Group{ID: "system", Title: "System:"},
		&cobra.Group{ID: "application", Title: "Applications:"},
		&cobra.Group{ID: "content_library_image", Title: "Content Library:"},
		&cobra.Group{ID: "vm", Title: "Virtual Machines:"},
		&cobra.Group{ID: "host", Title: "Hosts:"},
		&cobra.Group{ID: "cluster", Title: "Clusters:"},
		&cobra.Group{ID: "datastore", Title: "Datastores:"},
		&cobra.Group{ID: "network", Title: "Networks:"},
		&cobra.Group{ID: "task", Title: "Tasks:"},
		&cobra.Group{ID: "alert", Title: "Alerts:"},
		&cobra.Group{ID: "user", Title: "Users:"},
		&cobra.Group{ID: "metrics", Title: "Metrics:"},
	)

	alert.Register(rootCmd)
	alert_rule.Register(rootCmd)
	application.Register(rootCmd)
	cluster.Register(rootCmd)
	cluster_settings.Register(rootCmd)
	content_library_image.Register(rootCmd)
	datastore.Register(rootCmd)
	deploy.Register(rootCmd)
	events.Register(rootCmd)
	find.Register(rootCmd)
	host.Register(rootCmd)
	license.Register(rootCmd)
	network.Register(rootCmd)
	ntp.Register(rootCmd)
	system.Register(rootCmd)
	task.Register(rootCmd)
	user.Register(rootCmd)
	vlan.Register(rootCmd)
	vmcmd.Register(rootCmd)

	// metrics 子命令直接注册到根命令（去除 metrics 前缀）
	metrics.RegisterVMMetrics(rootCmd)
	metrics.RegisterHostMetrics(rootCmd)
	metrics.RegisterVmVolumeMetrics(rootCmd)
	metrics.RegisterClusterMetrics(rootCmd)
	metrics.RegisterVolumeMetrics(rootCmd)
	metrics.RegisterSFSMetrics(rootCmd)
}

// Execute 是 main.go 的唯一入口。
func Execute() error {
	err := rootCmd.Execute()
	if err != nil {
		errorExit(err)
	}
	return nil
}

// OutputFormat 返回当前全局 --json flag 值，供命令层调用。
func OutputFormat() string { return outputFlags.Format }

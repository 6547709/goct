// Package cmd 定义 goct 的所有 CLI 命令。
// root.go 提供根命令、全局 flag、PersistentPreRunE 构建 client。
package cmd

import (
	"context"

	"github.com/6547709/goct/cmd/alert"
	"github.com/6547709/goct/cmd/cluster"
	"github.com/6547709/goct/cmd/datastore"
	"github.com/6547709/goct/cmd/host"
	"github.com/6547709/goct/cmd/metrics"
	"github.com/6547709/goct/cmd/network"
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

		// 带 nologin 注解的命令跳过 client 构造（如 goct version/session.*）
		if c.Annotations["nologin"] == "true" {
			debug.Debugf("skipping login for command %q (nologin annotation)", c.Name())
			return nil
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
		debug.Debugf("resolved config: url=%s user=%s cluster=%s insecure=%v source=%s",
			cfg.URL, cfg.Username, cfg.Cluster, cfg.Insecure, cfg.Source)
		cli, err := client.New(c.Context(), cfg)
		if err != nil {
			return err
		}
		debug.Info("client created, login successful")
		c.SetContext(client.With(c.Context(), cli))
		return nil
	},
}

func init() {
	rootCmd.SetContext(context.Background())
	connFlags.Register(rootCmd)
	outputFlags.Register(rootCmd)

	rootCmd.AddGroup(
		&cobra.Group{ID: "system", Title: "System:"},
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
	cluster.Register(rootCmd)
	datastore.Register(rootCmd)
	host.Register(rootCmd)
	network.Register(rootCmd)
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

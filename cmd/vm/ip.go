package vm

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/service"
	"github.com/spf13/cobra"
)

// newIp 实现 govc 风格的 vm.ip：
//
//   - 默认行为：等到 VM 上报到至少一个 IP 后才输出（受 --wait 控制超时）。
//   - --a / --all 输出所有 IP（每行一个）；不加该 flag 时只输出第一个。
//   - --v4 / --v6：只输出 IPv4 / IPv6 地址。
//   - --no-wait：不等，直接读当前缓存值（即使为空也立即返回 0）。
//
// 与 govc 的差异：
//   - govc 用 vSphere 的 WaitForUpdates 走 push；CloudTower 没有等价 API，
//     只能客户端轮询 GetVM。轮询间隔 2s，默认超时 5min。
func newIp() *cobra.Command {
	var all, v4, v6, noWait bool
	var waitTimeout time.Duration

	c := &cobra.Command{
		Use:   "vm.ip [name|id]",
		Short: "Show VM IP addresses (waits for VM tools by default)",
		Long: `Output the IP address(es) of a VM, one per line.
Designed for scripting: pipe VM list IDs to get IPs.

By default the command waits up to --wait (default 5m) for the VM to report
at least one IP. Use --no-wait to read the cached value immediately.

Examples:
  goct vm.ip myvm                      # first IP, wait up to 5m
  goct vm.ip myvm -a                   # all IPs, one per line
  goct vm.ip myvm --v4                 # only IPv4
  goct vm.ip myvm --no-wait            # no wait, exit 0 even if empty
  goct vm.ls --id-only | xargs -I{} goct vm.ip {}`,
		GroupID: groupID,
		Args:    cobra.MaximumNArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			id, err := resolveVMArg(args)
			if err != nil {
				return err
			}
			svc := service.NewVM(cli)

			ctx := c.Context()
			if !noWait && waitTimeout > 0 {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, waitTimeout)
				defer cancel()
			}

			ips, err := waitForIPs(ctx, svc, id, !noWait)
			if err != nil {
				return err
			}

			ips = filterIPs(ips, v4, v6)
			if len(ips) == 0 {
				return nil
			}
			if !all {
				fmt.Fprintln(c.OutOrStdout(), ips[0])
				return nil
			}
			for _, ip := range ips {
				fmt.Fprintln(c.OutOrStdout(), ip)
			}
			return nil
		},
	}
	c.Flags().BoolVarP(&all, "all", "a", false, "Output all IPs, one per line")
	c.Flags().BoolVar(&v4, "v4", false, "IPv4 only")
	c.Flags().BoolVar(&v6, "v6", false, "IPv6 only")
	c.Flags().BoolVar(&noWait, "no-wait", false, "Do not wait for VM tools to populate IPs")
	c.Flags().DurationVar(&waitTimeout, "wait", 5*time.Minute, "Maximum wait duration when waiting for IPs")
	return c
}

// waitForIPs 在 wait=true 时轮询 GetVM 直到拿到至少一个 IP 或 ctx 终止。
// wait=false 时只读一次。
func waitForIPs(ctx context.Context, svc *service.VMService, id string, wait bool) ([]string, error) {
	const interval = 2 * time.Second
	for {
		v, err := svc.Resolve(ctx, id)
		if err != nil {
			return nil, err
		}
		if len(v.IPs) > 0 || !wait {
			return v.IPs, nil
		}
		select {
		case <-ctx.Done():
			// 超时不算错误：与 govc 行为一致，安静退出。
			return nil, nil
		case <-time.After(interval):
		}
	}
}

// filterIPs 按 v4 / v6 过滤；两者都未指定时返回全部。
func filterIPs(ips []string, v4Only, v6Only bool) []string {
	if !v4Only && !v6Only {
		return ips
	}
	out := make([]string, 0, len(ips))
	for _, raw := range ips {
		s := strings.TrimSpace(raw)
		ip := net.ParseIP(s)
		if ip == nil {
			continue
		}
		isV4 := ip.To4() != nil
		switch {
		case v4Only && isV4:
			out = append(out, s)
		case v6Only && !isV4:
			out = append(out, s)
		}
	}
	return out
}

package vm

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

// resolveVMArg 从位置参数、GOCT_VM 环境变量或 stdin 获取 VM 标识。
// 优先级：位置参数 > 环境变量 > stdin（管道）。
// stdin 仅在非 TTY 时读取（即 echo id | goct vm.info）。
func resolveVMArg(args []string) (string, error) {
	if len(args) > 0 && args[0] != "" {
		return args[0], nil
	}
	if v := os.Getenv("GOCT_VM"); v != "" {
		return v, nil
	}
	// 尝试从 stdin 管道读取（非 TTY 时）
	if stat, _ := os.Stdin.Stat(); stat.Mode()&os.ModeCharDevice == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			if line := strings.TrimSpace(scanner.Text()); line != "" {
				return line, nil
			}
		}
	}
	return "", errors.New("VM not specified: use positional arg, set GOCT_VM, or pipe via stdin")
}

package output

import "fmt"

// HumanBytes 用 IEC 二进制单位（KiB/MiB/GiB/TiB/PiB）格式化字节数。
// 小于 1 KiB 时按 "N B" 输出整数，无小数。
func HumanBytes(n uint64) string {
	const unit = uint64(1024)
	if n < unit {
		return fmt.Sprintf("%d B", n)
	}
	div, exp := unit, 0
	for x := n / unit; x >= unit; x /= unit {
		div *= unit
		exp++
	}
	suffixes := []string{"KiB", "MiB", "GiB", "TiB", "PiB"}
	if exp >= len(suffixes) {
		exp = len(suffixes) - 1
	}
	return fmt.Sprintf("%.1f %s", float64(n)/float64(div), suffixes[exp])
}

// JoinIPs 把 IP 列表合并为逗号分隔字符串；空列表返回 "-"。
func JoinIPs(ips []string) string {
	if len(ips) == 0 {
		return "-"
	}
	out := ips[0]
	for _, ip := range ips[1:] {
		out += "," + ip
	}
	return out
}

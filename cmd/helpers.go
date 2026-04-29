package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/6547709/goct/pkg/adapter"
)

// errorExit 把命令最终错误映射到 user-friendly 输出 + 进程 exit code：
//
//	0  成功（不会进入此函数）
//	2  ErrAuth（认证/登录失败）
//	3  ErrNotFound（资源未找到）
//	4  ErrTaskFailed（异步 task 终态 FAILED）
//	1  其它
//
// 错误信息一律走 stderr，业务输出走 stdout，避免污染管道。
func errorExit(err error) {
	if err == nil {
		return
	}
	fmt.Fprintln(os.Stderr, "Error:", err)
	switch {
	case errors.Is(err, adapter.ErrAuth):
		os.Exit(2)
	case errors.Is(err, adapter.ErrNotFound):
		os.Exit(3)
	case errors.Is(err, adapter.ErrTaskFailed):
		os.Exit(4)
	default:
		os.Exit(1)
	}
}

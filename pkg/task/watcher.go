// Package task 提供异步任务进度的轮询与展示。
//
// CloudTower 写操作返回 task ID，本包封装统一的 Watcher：
//   - 1s 间隔轮询 GetTaskProgress（间隔可调）
//   - TTY 下用 \r 单行刷新进度，非 TTY / Quiet 静默
//   - 终结状态映射：SUCCESSED → nil；FAILED → ErrFailed 包装
//   - ctx 超时 / Ctrl-C 优雅退出
//
// 重要：CloudTower SDK 实际成功状态字符串是 "SUCCESSED"（拼写错但已固化），
// 不要写成 "SUCCEEDED"。
package task

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

// ErrFailed 表示远端 task 状态为 FAILED。可用 errors.Is 匹配。
var ErrFailed = errors.New("task failed")

// 终结态字符串（与 SDK 保持一致）。
const (
	statusSuccess = "SUCCESSED"
	statusFailed  = "FAILED"
)

// Ops 是 watcher 依赖的最小接口；由 pkg/adapter 实现。
//
// percent 取值 [0,100]；status 为 SDK 原始字符串。
type Ops interface {
	GetTaskProgress(ctx context.Context, id string) (percent int, status string, err error)
}

// Options 控制 watcher 行为。
type Options struct {
	Out      io.Writer     // 进度输出目标，nil 默认 os.Stderr
	Interval time.Duration // 轮询间隔，<=0 默认 1s
	Quiet    bool          // 静默模式（JSON 输出 / 非 TTY 时使用）
}

// Watcher 持有依赖与配置；通过 New 构造。
type Watcher struct {
	ops  Ops
	opts Options
}

// New 构造 Watcher，自动填充默认值。
func New(ops Ops, opts Options) *Watcher {
	if opts.Out == nil {
		opts.Out = os.Stderr
	}
	if opts.Interval <= 0 {
		opts.Interval = time.Second
	}
	return &Watcher{ops: ops, opts: opts}
}

// Watch 阻塞轮询直到 task 终结或 ctx 取消。
//
//   - 终态 SUCCESSED → 返回 nil
//   - 终态 FAILED    → 返回 wrap(ErrFailed)
//   - ctx 取消/超时   → 返回 ctx.Err()
//   - 查询出错        → 返回 wrap(原错误)
func (w *Watcher) Watch(ctx context.Context, id string) error {
	t := time.NewTicker(w.opts.Interval)
	defer t.Stop()

	for {
		// 立即查一次，避免首次轮询前先 sleep 一个 Interval。
		percent, status, err := w.ops.GetTaskProgress(ctx, id)
		if err != nil {
			return fmt.Errorf("query task %s: %w", id, err)
		}
		w.print(id, percent, status)

		switch status {
		case statusSuccess:
			w.newline()
			return nil
		case statusFailed:
			w.newline()
			return fmt.Errorf("task %s: %w", id, ErrFailed)
		}

		select {
		case <-ctx.Done():
			w.newline()
			return ctx.Err()
		case <-t.C:
		}
	}
}

// print 输出一行进度（非 Quiet 模式）。
func (w *Watcher) print(id string, percent int, status string) {
	if w.opts.Quiet {
		return
	}
	fmt.Fprintf(w.opts.Out, "\rTask %s: %d%% (%s)   ", id, percent, status)
}

// newline 在终结时换行，避免下一行业务输出与进度行黏在一起。
func (w *Watcher) newline() {
	if w.opts.Quiet {
		return
	}
	fmt.Fprintln(w.opts.Out)
}

package task_test

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/6547709/goct/pkg/task"
)

// fakeOps 是 watcher 的依赖 stub；按 progresses/statuses 顺序返回，超出后返回末尾稳定值。
type fakeOps struct {
	progresses []int
	statuses   []string
	idx        int
}

func (f *fakeOps) GetTaskProgress(_ context.Context, _ string) (int, string, error) {
	if f.idx >= len(f.progresses) {
		return f.progresses[len(f.progresses)-1], f.statuses[len(f.statuses)-1], nil
	}
	p, s := f.progresses[f.idx], f.statuses[f.idx]
	f.idx++
	return p, s, nil
}

// TestWatcher_Success 验证状态走到 SUCCESSED 时返回 nil。
func TestWatcher_Success(t *testing.T) {
	f := &fakeOps{
		progresses: []int{10, 50, 100},
		statuses:   []string{"EXECUTING", "EXECUTING", "SUCCESSED"},
	}
	w := task.New(f, task.Options{Out: &bytes.Buffer{}, Interval: time.Millisecond, Quiet: true})
	if err := w.Watch(context.Background(), "tid"); err != nil {
		t.Fatal(err)
	}
}

// TestWatcher_Failed 验证 FAILED 状态返回 ErrFailed 链。
func TestWatcher_Failed(t *testing.T) {
	f := &fakeOps{
		progresses: []int{50, 100},
		statuses:   []string{"EXECUTING", "FAILED"},
	}
	w := task.New(f, task.Options{Out: &bytes.Buffer{}, Interval: time.Millisecond, Quiet: true})
	err := w.Watch(context.Background(), "tid")
	if !errors.Is(err, task.ErrFailed) {
		t.Fatalf("err=%v, want wrap ErrFailed", err)
	}
}

// TestWatcher_Cancel 验证 ctx 超时/取消时返回 ctx.Err。
func TestWatcher_Cancel(t *testing.T) {
	f := &fakeOps{progresses: []int{10}, statuses: []string{"EXECUTING"}}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()
	w := task.New(f, task.Options{Interval: time.Millisecond, Quiet: true})
	err := w.Watch(ctx, "tid")
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("err=%v, want DeadlineExceeded", err)
	}
}

// TestWatcher_QuietNoOutput 验证 Quiet=true 时不向 Out 写入任何内容。
func TestWatcher_QuietNoOutput(t *testing.T) {
	var buf bytes.Buffer
	f := &fakeOps{progresses: []int{100}, statuses: []string{"SUCCESSED"}}
	w := task.New(f, task.Options{Out: &buf, Interval: time.Millisecond, Quiet: true})
	if err := w.Watch(context.Background(), "tid"); err != nil {
		t.Fatal(err)
	}
	if buf.Len() != 0 {
		t.Fatalf("expected silent output, got %q", buf.String())
	}
}

// TestWatcher_VerboseWritesProgress 验证非 Quiet 模式写入了进度信息。
func TestWatcher_VerboseWritesProgress(t *testing.T) {
	var buf bytes.Buffer
	f := &fakeOps{progresses: []int{100}, statuses: []string{"SUCCESSED"}}
	w := task.New(f, task.Options{Out: &buf, Interval: time.Millisecond})
	if err := w.Watch(context.Background(), "tid"); err != nil {
		t.Fatal(err)
	}
	if buf.Len() == 0 {
		t.Fatal("expected progress output, got empty")
	}
}

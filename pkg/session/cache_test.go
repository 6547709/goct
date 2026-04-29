package session_test

import (
	"os"
	"testing"
	"time"

	"github.com/6547709/goct/pkg/session"
)

// TestSaveAndLoad 验证 save → load 闭环正确。
func TestSaveAndLoad(t *testing.T) {
	t.Setenv("XDG_CACHE_HOME", t.TempDir())
	tok := session.Token{Value: "abc", ExpireAt: time.Now().Add(time.Hour)}
	if err := session.Save("h:443", "u", tok); err != nil {
		t.Fatal(err)
	}
	got, err := session.Load("h:443", "u")
	if err != nil {
		t.Fatal(err)
	}
	if got.Value != "abc" {
		t.Fatalf("token = %q, want abc", got.Value)
	}
}

// TestLoadMiss 验证文件不存在时 Load 返回 os.IsNotExist 可识别错误。
func TestLoadMiss(t *testing.T) {
	t.Setenv("XDG_CACHE_HOME", t.TempDir())
	_, err := session.Load("nope", "x")
	if !os.IsNotExist(err) {
		t.Fatalf("err=%v, want IsNotExist", err)
	}
}

// TestExpired 验证过期 token 不允许加载。
func TestExpired(t *testing.T) {
	t.Setenv("XDG_CACHE_HOME", t.TempDir())
	if err := session.Save("h", "u", session.Token{Value: "x", ExpireAt: time.Now().Add(-time.Hour)}); err != nil {
		t.Fatal(err)
	}
	if _, err := session.Load("h", "u"); err == nil {
		t.Fatal("expected expiry error")
	}
}

// TestPerm0600 验证写入文件权限为 0600（防 token 泄露）。
func TestPerm0600(t *testing.T) {
	t.Setenv("XDG_CACHE_HOME", t.TempDir())
	if err := session.Save("h", "u", session.Token{Value: "x", ExpireAt: time.Now().Add(time.Hour)}); err != nil {
		t.Fatal(err)
	}
	st, err := os.Stat(session.PathFor("h", "u"))
	if err != nil {
		t.Fatal(err)
	}
	if st.Mode().Perm() != 0o600 {
		t.Fatalf("perm=%v, want 0600", st.Mode().Perm())
	}
}

// TestDelete 验证 Delete 幂等且能真正移除文件。
func TestDelete(t *testing.T) {
	t.Setenv("XDG_CACHE_HOME", t.TempDir())
	_ = session.Save("h", "u", session.Token{Value: "x", ExpireAt: time.Now().Add(time.Hour)})
	if err := session.Delete("h", "u"); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(session.PathFor("h", "u")); !os.IsNotExist(err) {
		t.Fatalf("file still exists: %v", err)
	}
	// 再删一次应当无错（幂等）
	if err := session.Delete("h", "u"); err != nil {
		t.Fatalf("second delete err=%v", err)
	}
}

// TestList 验证 List 列出当前 cache 目录下所有 session 文件路径。
func TestList(t *testing.T) {
	t.Setenv("XDG_CACHE_HOME", t.TempDir())
	_ = session.Save("h1", "u1", session.Token{Value: "a", ExpireAt: time.Now().Add(time.Hour)})
	_ = session.Save("h2", "u2", session.Token{Value: "b", ExpireAt: time.Now().Add(time.Hour)})
	got, err := session.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Fatalf("List len=%d, want 2", len(got))
	}
}

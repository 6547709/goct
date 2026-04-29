package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/6547709/goct/pkg/config"
)

// TestResolve_FlagOverEnv 验证 CLI flag 覆盖环境变量。
func TestResolve_FlagOverEnv(t *testing.T) {
	t.Setenv("GOCT_URL", "https://from-env")
	got, err := config.Resolve(config.Override{URL: "https://from-flag"})
	if err != nil {
		t.Fatal(err)
	}
	if got.URL != "https://from-flag" {
		t.Fatalf("URL = %q, want https://from-flag", got.URL)
	}
}

// TestResolve_EnvOverFile 验证环境变量覆盖配置文件。
func TestResolve_EnvOverFile(t *testing.T) {
	dir := t.TempDir()
	cfg := filepath.Join(dir, ".goct.yaml")
	if err := os.WriteFile(cfg, []byte("url: https://from-file\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("GOCT_URL", "https://from-env")
	got, err := config.Resolve(config.Override{ConfigFile: cfg})
	if err != nil {
		t.Fatal(err)
	}
	if got.URL != "https://from-env" {
		t.Fatalf("URL = %q, want https://from-env", got.URL)
	}
}

// TestResolve_FileOnly 验证仅靠配置文件也能解析出 URL。
func TestResolve_FileOnly(t *testing.T) {
	dir := t.TempDir()
	cfg := filepath.Join(dir, ".goct.yaml")
	body := "url: https://from-file\nusername: alice\ninsecure: true\n"
	if err := os.WriteFile(cfg, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	// 隔离环境变量，避免主机污染
	t.Setenv("GOCT_URL", "")
	t.Setenv("GOCT_USERNAME", "")
	t.Setenv("GOCT_INSECURE", "")
	got, err := config.Resolve(config.Override{ConfigFile: cfg})
	if err != nil {
		t.Fatal(err)
	}
	if got.URL != "https://from-file" || got.Username != "alice" || !got.Insecure {
		t.Fatalf("got=%+v", got)
	}
}

// TestResolve_DefaultSource 验证 source 缺省为 local。
func TestResolve_DefaultSource(t *testing.T) {
	t.Setenv("GOCT_SOURCE", "")
	got, err := config.Resolve(config.Override{})
	if err != nil {
		t.Fatal(err)
	}
	if got.Source != "local" {
		t.Fatalf("source = %q, want local", got.Source)
	}
}

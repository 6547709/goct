package flags_test

import (
	"testing"

	"github.com/6547709/goct/pkg/flags"
	"github.com/spf13/cobra"
)

// TestConnectionFlags_Register 验证连接相关 flag 都注册到 PersistentFlags 上。
func TestConnectionFlags_Register(t *testing.T) {
	c := &cobra.Command{Use: "x"}
	cf := &flags.ConnectionFlags{}
	cf.Register(c)
	for _, name := range []string{"url", "username", "password", "insecure", "cluster"} {
		if c.PersistentFlags().Lookup(name) == nil {
			t.Fatalf("missing flag: %s", name)
		}
	}
}

// TestOutputFlags_Default 验证 OutputFlags 默认 format 为 table。
func TestOutputFlags_Default(t *testing.T) {
	c := &cobra.Command{Use: "x"}
	of := &flags.OutputFlags{}
	of.Register(c)
	if of.Format != "table" {
		t.Fatalf("default format = %q, want table", of.Format)
	}
}

// TestSearchFlags_Register 验证 SearchFlags 注册了 name/limit/skip。
func TestSearchFlags_Register(t *testing.T) {
	c := &cobra.Command{Use: "x"}
	sf := &flags.SearchFlags{}
	sf.Register(c)
	for _, name := range []string{"name", "limit", "skip"} {
		if c.Flags().Lookup(name) == nil {
			t.Fatalf("missing flag: %s", name)
		}
	}
}

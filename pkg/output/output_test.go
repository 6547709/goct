package output_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/6547709/goct/pkg/output"
)

type row struct{ Name, Status string }

var cols = []output.Column{
	{Header: "NAME", Get: func(v any) string { return v.(row).Name }},
	{Header: "STATUS", Get: func(v any) string { return v.(row).Status }},
}

// TestRender_Table 验证 table 渲染包含表头与数据。
func TestRender_Table(t *testing.T) {
	var buf bytes.Buffer
	data := []any{row{"vm1", "RUNNING"}, row{"vm2", "STOPPED"}}
	if err := output.Render(&buf, data, "table", cols); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	for _, want := range []string{"NAME", "STATUS", "vm1", "RUNNING", "vm2", "STOPPED"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
}

// TestRender_TableEmpty 验证空数据集只渲染表头不报错。
func TestRender_TableEmpty(t *testing.T) {
	var buf bytes.Buffer
	if err := output.Render(&buf, []any{}, "table", cols); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "NAME") {
		t.Fatalf("missing header in empty render:\n%s", buf.String())
	}
}

// TestRender_JSON 验证 json 模式直接编码 data。
func TestRender_JSON(t *testing.T) {
	var buf bytes.Buffer
	if err := output.Render(&buf, []row{{"vm1", "RUNNING"}}, "json", nil); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, `"Name": "vm1"`) || !strings.Contains(out, `"Status": "RUNNING"`) {
		t.Fatalf("unexpected json:\n%s", out)
	}
}

// TestRender_DefaultIsTable 验证空 format 默认走 table 路径。
func TestRender_DefaultIsTable(t *testing.T) {
	var buf bytes.Buffer
	if err := output.Render(&buf, []any{row{"x", "y"}}, "", cols); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "NAME") {
		t.Fatalf("default did not render table header:\n%s", buf.String())
	}
}

// TestRender_BadFormat 验证未知 format 返回错误。
func TestRender_BadFormat(t *testing.T) {
	if err := output.Render(&bytes.Buffer{}, []any{}, "yaml", cols); err == nil {
		t.Fatal("expected error for unknown format")
	}
}

// TestHumanBytes 验证字节单位换算。
func TestHumanBytes(t *testing.T) {
	const (
		KiB uint64 = 1024
		MiB        = KiB * 1024
		GiB        = MiB * 1024
		TiB        = GiB * 1024
	)
	cases := map[uint64]string{
		0:        "0 B",
		512:      "512 B",
		KiB:      "1.0 KiB",
		MiB:      "1.0 MiB",
		GiB:      "1.0 GiB",
		3 * GiB / 2: "1.5 GiB",
		TiB:      "1.0 TiB",
	}
	for in, want := range cases {
		if got := output.HumanBytes(in); got != want {
			t.Errorf("HumanBytes(%d)=%q want %q", in, got, want)
		}
	}
}

// TestJoinIPs 验证 IP 列表合并；空列表返回 "-"。
func TestJoinIPs(t *testing.T) {
	if got := output.JoinIPs(nil); got != "-" {
		t.Errorf("empty: got %q want -", got)
	}
	if got := output.JoinIPs([]string{"10.0.0.1"}); got != "10.0.0.1" {
		t.Errorf("single: got %q", got)
	}
	if got := output.JoinIPs([]string{"10.0.0.1", "10.0.0.2"}); got != "10.0.0.1,10.0.0.2" {
		t.Errorf("multi: got %q", got)
	}
}

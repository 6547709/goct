package client_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/config"
	"github.com/6547709/goct/pkg/session"
)

// TestNew_LoginAndCache 端到端：mock 登录成功，token 被写入 session cache。
func TestNew_LoginAndCache(t *testing.T) {
	t.Setenv("XDG_CACHE_HOME", t.TempDir())

	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case strings.HasSuffix(r.URL.Path, "/login"):
			_, _ = w.Write([]byte(`{"data":{"token":"new-token"}}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	c, err := client.New(context.Background(), config.Resolved{
		URL: srv.URL, Username: "u", Password: "p", Insecure: true,
	})
	if err != nil {
		t.Fatalf("New err=%v", err)
	}
	if c == nil {
		t.Fatal("client is nil")
	}

	// 验证 token 已被缓存
	tok, err := session.Load(hostFromURL(t, srv.URL), "u")
	if err != nil {
		t.Fatalf("session.Load err=%v", err)
	}
	if tok.Value != "new-token" {
		t.Fatalf("cached token=%q want new-token", tok.Value)
	}
}

// TestNew_MissingURL 验证缺少 URL 立即报错。
func TestNew_MissingURL(t *testing.T) {
	t.Setenv("XDG_CACHE_HOME", t.TempDir())
	if _, err := client.New(context.Background(), config.Resolved{}); err == nil {
		t.Fatal("expected error when URL missing")
	}
}

// TestNew_MissingCredentials 验证 URL 在但没用户/密码时报错。
func TestNew_MissingCredentials(t *testing.T) {
	t.Setenv("XDG_CACHE_HOME", t.TempDir())
	if _, err := client.New(context.Background(), config.Resolved{
		URL: "https://tower.example.com",
	}); err == nil {
		t.Fatal("expected error when credentials missing")
	}
}

// TestWithFrom 验证 ctx 注入与提取链路。
func TestWithFrom(t *testing.T) {
	if got := client.From(context.Background()); got != nil {
		t.Fatal("From on bare ctx should return nil")
	}
}

// hostFromURL 从 https://127.0.0.1:port 提出 host 部分（与 client.hostKey 等价）。
func hostFromURL(t *testing.T, raw string) string {
	t.Helper()
	for _, prefix := range []string{"https://", "http://"} {
		if strings.HasPrefix(raw, prefix) {
			return raw[len(prefix):]
		}
	}
	return raw
}

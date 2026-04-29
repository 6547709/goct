package adapter_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/6547709/goct/pkg/adapter"
)

// newFakeTower 返回一个最小可用的 CloudTower mock：
//   - /v2/api/login 返回带 token 的 WithTask_LoginResponse
//   - /v2/api/get-version 返回裸 JSON 字符串 "4.8.1"
func newFakeTower(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case strings.HasSuffix(r.URL.Path, "/login"):
			// CloudTower WithTask_LoginResponse_：data 含 token，task_id 可省
			_, _ = w.Write([]byte(`{"data":{"token":"fake-token"}}`))
		case strings.HasSuffix(r.URL.Path, "/get-version"):
			// SDK GetAPIVersionOK.Payload 是裸 string；OpenAPI 把整个响应体当 string 反序列化
			_, _ = w.Write([]byte(`"4.8.1"`))
		default:
			http.NotFound(w, r)
		}
	}))
}

// TestNewClient_LoginAndAbout 端到端验证：mock 登录 → About 拿到版本。
func TestNewClient_LoginAndAbout(t *testing.T) {
	srv := newFakeTower(t)
	defer srv.Close()

	c, tok, err := adapter.NewClient(context.Background(), adapter.Options{
		URL: srv.URL, Username: "u", Password: "p", Insecure: true,
	})
	if err != nil {
		t.Fatalf("NewClient err=%v", err)
	}
	if tok.Value != "fake-token" {
		t.Fatalf("token = %q, want fake-token", tok.Value)
	}

	info, err := c.About(context.Background())
	if err != nil {
		t.Fatalf("About err=%v", err)
	}
	if info.Version != "4.8.1" {
		t.Fatalf("version = %q, want 4.8.1", info.Version)
	}
}

// TestNewClient_TokenSkipsLogin 验证 Token 提供时不调 /login（mock 中无 /login 也能成功）。
func TestNewClient_TokenSkipsLogin(t *testing.T) {
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/login") {
			t.Errorf("login should not be called when token provided; got path=%s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`"4.8.1"`))
	}))
	defer srv.Close()

	c, tok, err := adapter.NewClient(context.Background(), adapter.Options{
		URL: srv.URL, Token: "preset", Insecure: true,
	})
	if err != nil {
		t.Fatalf("err=%v", err)
	}
	if tok.Value != "preset" {
		t.Fatalf("token = %q", tok.Value)
	}
	if _, err := c.About(context.Background()); err != nil {
		t.Fatalf("About err=%v", err)
	}
}

// TestNewClient_BadURL 验证非法 URL 立即报错，不发请求。
func TestNewClient_BadURL(t *testing.T) {
	for _, bad := range []string{"", "not-a-url", "ftp://x.y", "no-scheme.example.com"} {
		_, _, err := adapter.NewClient(context.Background(), adapter.Options{URL: bad})
		if err == nil {
			t.Errorf("NewClient(%q) want error, got nil", bad)
		}
	}
}

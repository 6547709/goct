package debug

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

// mockTransport records all requests that pass through, for verifying trace output.
type mockTransport struct {
	Requests []*http.Request
	Response *http.Response
	Err      error
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	m.Requests = append(m.Requests, req)
	return m.Response, m.Err
}

func TestTraceRoundTripper_TraceLevel(t *testing.T) {
	// 1. Create mockTransport
	mock := &mockTransport{
		Response: &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(`{"result":"ok"}`)),
			Header:     http.Header{},
		},
	}

	// 2. Create TraceRoundTripper with TraceLevelTrace
	buf := &bytes.Buffer{}
	tracer := NewTraceRoundTripper(mock, TraceLevelTrace, buf)

	// 3.发起 fake 请求
	req, _ := http.NewRequest("GET", "http://example.com/api/vms", nil)
	req.Header.Set("Content-Type", "application/json")
	_, err := tracer.RoundTrip(req)

	if err != nil {
		t.Fatalf("RoundTrip returned error: %v", err)
	}

	// 4. 验证 trace 输出 JSON 包含 method/path/status/duration_ms
	output := buf.String()
	if output == "" {
		t.Fatal("expected trace output, got empty string")
	}

	var entry TraceEntry
	if err := json.Unmarshal([]byte(output), &entry); err != nil {
		t.Fatalf("failed to parse trace JSON: %v", err)
	}

	if entry.Level != "trace" {
		t.Errorf("expected level 'trace', got '%s'", entry.Level)
	}
	if entry.Method != "GET" {
		t.Errorf("expected method 'GET', got '%s'", entry.Method)
	}
	if entry.Path != "/api/vms" {
		t.Errorf("expected path '/api/vms', got '%s'", entry.Path)
	}
	if entry.RespStatus != 200 {
		t.Errorf("expected status 200, got %d", entry.RespStatus)
	}
	if entry.DurationMs < 0 {
		t.Errorf("expected non-negative duration_ms, got %d", entry.DurationMs)
	}

	// At trace level, headers and body should be empty
	if entry.ReqHeaders != nil {
		t.Errorf("expected nil req_headers at trace level, got %v", entry.ReqHeaders)
	}
	if entry.RespHeaders != nil {
		t.Errorf("expected nil resp_headers at trace level, got %v", entry.RespHeaders)
	}
	if entry.ReqBody != "" {
		t.Errorf("expected empty req_body at trace level, got '%s'", entry.ReqBody)
	}
	if entry.RespBody != "" {
		t.Errorf("expected empty resp_body at trace level, got '%s'", entry.RespBody)
	}
}

func TestTraceRoundTripper_VerboseLevel(t *testing.T) {
	// 1. Create mockTransport returning response with body
	body := `{"data":[{"id":"vm-1","name":"test-vm"}]}`
	mock := &mockTransport{
		Response: &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     http.Header{"Content-Type": []string{"application/json"}},
		},
	}

	// 2. Create TraceRoundTripper with TraceLevelVerbose
	buf := &bytes.Buffer{}
	tracer := NewTraceRoundTripper(mock, TraceLevelVerbose, buf)

	// Create request with headers and body
	reqBody := `{"name":"test-vm","password":"secret"}`
	req, _ := http.NewRequest("POST", "http://example.com/api/vms", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer token123")
	_, err := tracer.RoundTrip(req)

	if err != nil {
		t.Fatalf("RoundTrip returned error: %v", err)
	}

	// 3. Verify trace output contains req_headers/resp_headers/req_body/resp_body
	output := buf.String()
	if output == "" {
		t.Fatal("expected trace output, got empty string")
	}

	var entry TraceEntry
	if err := json.Unmarshal([]byte(output), &entry); err != nil {
		t.Fatalf("failed to parse trace JSON: %v", err)
	}

	// Verify req_headers
	if entry.ReqHeaders == nil {
		t.Fatal("expected req_headers at verbose level, got nil")
	}
	if entry.ReqHeaders["Content-Type"] != "application/json" {
		t.Errorf("expected Content-Type 'application/json', got '%s'", entry.ReqHeaders["Content-Type"])
	}

	// 4. Verify Authorization header is filtered to ***
	if entry.ReqHeaders["Authorization"] != "***" {
		t.Errorf("expected Authorization '***', got '%s'", entry.ReqHeaders["Authorization"])
	}

	// Verify resp_headers
	if entry.RespHeaders == nil {
		t.Fatal("expected resp_headers at verbose level, got nil")
	}
	if entry.RespHeaders["Content-Type"] != "application/json" {
		t.Errorf("expected Content-Type 'application/json', got '%s'", entry.RespHeaders["Content-Type"])
	}

	// Verify req_body and resp_body
	if entry.ReqBody == "" {
		t.Error("expected non-empty req_body at verbose level")
	}
	// Verify password is filtered in body
	if strings.Contains(entry.ReqBody, "secret") {
		t.Errorf("expected password to be filtered, got req_body: '%s'", entry.ReqBody)
	}

	if entry.RespBody == "" {
		t.Error("expected non-empty resp_body at verbose level")
	}
}

func TestTraceRoundTripper_BodyTruncation(t *testing.T) {
	// 1. Create a body > 64KB
	largeBody := strings.Repeat(`{"id":"vm-1","name":"test-vm","data":"`+"x"+`"}`, 15000) // ~65KB
	// Make sure it's valid JSON
	largeBody = "[" + largeBody + "]"

	mock := &mockTransport{
		Response: &http.Response{
			StatusCode:   200,
			Body:          io.NopCloser(strings.NewReader(largeBody)),
			Header:        http.Header{"Content-Type": []string{"application/json"}},
			ContentLength: int64(len(largeBody)),
		},
	}

	// 2. Create TraceRoundTripper with TraceLevelVerbose
	buf := &bytes.Buffer{}
	tracer := NewTraceRoundTripper(mock, TraceLevelVerbose, buf)

	req, _ := http.NewRequest("GET", "http://example.com/api/vms", nil)
	_, err := tracer.RoundTrip(req)

	if err != nil {
		t.Fatalf("RoundTrip returned error: %v", err)
	}

	output := buf.String()
	var entry TraceEntry
	if err := json.Unmarshal([]byte(output), &entry); err != nil {
		t.Fatalf("failed to parse trace JSON: %v", err)
	}

	// 3. Verify body_truncated=true
	if !entry.BodyTruncated {
		t.Error("expected body_truncated=true for large body at verbose level")
	}

	// 4. Create TraceRoundTripper with TraceLevelDump
	buf2 := &bytes.Buffer{}
	tracer2 := NewTraceRoundTripper(mock, TraceLevelDump, buf2)

	// Create new request since body was consumed
	req2, _ := http.NewRequest("GET", "http://example.com/api/vms", nil)
	_, err = tracer2.RoundTrip(req2)

	if err != nil {
		t.Fatalf("RoundTrip returned error: %v", err)
	}

	output2 := buf2.String()
	var entry2 TraceEntry
	if err := json.Unmarshal([]byte(output2), &entry2); err != nil {
		t.Fatalf("failed to parse trace JSON: %v", err)
	}

	// 5. Verify body_truncated=false
	if entry2.BodyTruncated {
		t.Error("expected body_truncated=false at dump level")
	}
}

func TestTraceRoundTripper_Error(t *testing.T) {
	mock := &mockTransport{
		Err: http.ErrServerClosed,
	}

	buf := &bytes.Buffer{}
	tracer := NewTraceRoundTripper(mock, TraceLevelTrace, buf)

	req, _ := http.NewRequest("GET", "http://example.com/api/vms", nil)
	_, err := tracer.RoundTrip(req)

	if err != http.ErrServerClosed {
		t.Errorf("expected http.ErrServerClosed, got %v", err)
	}

	output := buf.String()
	var entry TraceEntry
	if err := json.Unmarshal([]byte(output), &entry); err != nil {
		t.Fatalf("failed to parse trace JSON: %v", err)
	}

	if entry.Error != http.ErrServerClosed.Error() {
		t.Errorf("expected error '%s', got '%s'", http.ErrServerClosed.Error(), entry.Error)
	}
	if entry.Level != "trace" {
		t.Errorf("expected level 'trace', got '%s'", entry.Level)
	}
}

func TestTraceRoundTripper_Host(t *testing.T) {
	mock := &mockTransport{
		Response: &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(`{}`)),
		},
	}

	buf := &bytes.Buffer{}
	tracer := NewTraceRoundTripper(mock, TraceLevelTrace, buf)

	req, _ := http.NewRequest("GET", "http://example.com:8080/api/vms", nil)
	_, _ = tracer.RoundTrip(req)

	output := buf.String()
	var entry TraceEntry
	if err := json.Unmarshal([]byte(output), &entry); err != nil {
		t.Fatalf("failed to parse trace JSON: %v", err)
	}

	if entry.Host != "example.com:8080" {
		t.Errorf("expected host 'example.com:8080', got '%s'", entry.Host)
	}
}

// TestTraceRoundTripper_SetBase 锁定 Bug 1 修复：SetBase 后 trace 应使用新 base。
//
// 历史 Bug：cmd/root.go 调 NewTraceRoundTripper(nil, ...)，nil base 会回退到
// http.DefaultTransport（做 TLS 校验），导致 --trace --insecure 同开时仍校验自签名。
// 修复方案是 adapter.newTransport 在装配时 SetBase 注入 insecure transport。
func TestTraceRoundTripper_SetBase(t *testing.T) {
	first := &mockTransport{Response: &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(""))}}
	tr := NewTraceRoundTripper(nil, TraceLevelTrace, &bytes.Buffer{})

	// SetBase 之前默认是 http.DefaultTransport（不可控）。
	// SetBase(first) 之后，所有请求必须从 first 走。
	tr.SetBase(first)

	req, _ := http.NewRequest("GET", "https://example.com/x", nil)
	if _, err := tr.RoundTrip(req); err != nil {
		t.Fatalf("RoundTrip err=%v", err)
	}
	if len(first.Requests) != 1 {
		t.Fatalf("expected request to go through first transport, got %d", len(first.Requests))
	}

	// SetBase(nil) 不应该清空已有 base（防止误覆盖）。
	tr.SetBase(nil)
	if _, err := tr.RoundTrip(req); err != nil {
		t.Fatalf("RoundTrip err=%v", err)
	}
	if len(first.Requests) != 2 {
		t.Fatalf("expected 2 requests through first after SetBase(nil), got %d", len(first.Requests))
	}
}

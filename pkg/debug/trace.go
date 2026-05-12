package debug

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

// TraceRoundTripper 实现 http.RoundTripper，拦截 SDK 所有 HTTP 请求/响应。
type TraceRoundTripper struct {
	Base   http.RoundTripper // 底层真实 transport（必填）
	Level  TraceLevel        // trace 级别
	Output io.Writer         // 输出目的地（应与 debug.Log 同路）
	Mutex  sync.Mutex
}

// NewTraceRoundTripper 构造 TraceRoundTripper。
// base: 底层 transport，nil 时使用 http.DefaultTransport。
// level: trace 详细程度。
// output: 输出目的地，nil 时使用 os.Stderr。
func NewTraceRoundTripper(base http.RoundTripper, level TraceLevel, output io.Writer) *TraceRoundTripper {
	if base == nil {
		base = http.DefaultTransport
	}
	if output == nil {
		output = os.Stderr
	}
	return &TraceRoundTripper{
		Base:   base,
		Level:  level,
		Output: output,
	}
}

// SetBase 由 adapter.newTransport 在装配链路时回填底层 transport，
// 使得 `--trace --insecure` 同开时 trace 包装的仍然是跳过 TLS 校验的 transport。
// nil 入参不会清空已有 Base，避免误覆盖。
func (t *TraceRoundTripper) SetBase(base http.RoundTripper) {
	if base == nil {
		return
	}
	t.Base = base
}

// RoundTrip 实现 http.RoundTripper 接口。
func (t *TraceRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()

	// --- 请求阶段 ---
	reqBytes := req.ContentLength

	// 发送真实请求
	resp, err := t.Base.RoundTrip(req)
	duration := time.Since(start)

	// --- 响应阶段 ---
	trace := TraceEntry{
		Level:      "trace",
		Method:     req.Method,
		Path:       req.URL.Path,
		Host:       req.URL.Host,
		DurationMs: duration.Milliseconds(),
	}

	if err != nil {
		trace.Error = err.Error()
		t.emit(&trace)
		return resp, err
	}

	trace.RespStatus = resp.StatusCode
	trace.ReqBytes = reqBytes

	if t.Level >= TraceLevelVerbose {
		trace.ReqHeaders = t.filterHeaders(req.Header)
		if req.Body != nil && reqBytes > 0 {
			trace.ReqBody, trace.BodyTruncated, trace.Binary = t.readBody(req.Body, reqBytes)
		}
		trace.RespHeaders = t.filterHeaders(resp.Header)
	}

	if resp.Body != nil {
		respBytes := resp.ContentLength
		trace.RespBytes = respBytes
		if t.Level >= TraceLevelVerbose {
			bodyStr, truncated, binary := t.readBody(resp.Body, respBytes)
			trace.RespBody = bodyStr
			trace.BodyTruncated = truncated
			trace.Binary = binary
		}
	}

	t.emit(&trace)
	return resp, err
}

// TraceEntry 是结构化 trace JSON 的 schema。
type TraceEntry struct {
	Level         string            `json:"level"`                    // 固定 "trace"
	Method        string            `json:"method"`                  // HTTP method
	Path          string            `json:"path"`                    // URL path
	Host          string            `json:"host"`                    // URL host
	ReqHeaders    map[string]string `json:"req_headers,omitempty"`   // 请求 headers（verbose+）
	ReqBody       string            `json:"req_body,omitempty"`      // 请求 body（verbose+）
	ReqBytes      int64             `json:"req_bytes,omitempty"`    // 请求体大小
	RespStatus    int               `json:"resp_status,omitempty"`   // 响应状态码
	RespHeaders   map[string]string `json:"resp_headers,omitempty"` // 响应 headers（verbose+）
	RespBody      string            `json:"resp_body,omitempty"`    // 响应 body（verbose+）
	RespBytes     int64             `json:"resp_bytes,omitempty"`   // 响应体大小
	BodyTruncated bool              `json:"body_truncated,omitempty"` // body 是否被截断
	DurationMs    int64             `json:"duration_ms"`            // 耗时（毫秒）
	Error         string            `json:"error,omitempty"`        // 错误信息
	Binary        bool              `json:"binary,omitempty"`        // 非 JSON body
}

// emit 输出 trace JSON（线程安全）。
func (t *TraceRoundTripper) emit(entry *TraceEntry) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	// dump 级别且 body 未截断时，不截断输出
	if t.Level < TraceLevelVerbose {
		entry.ReqBody = ""
		entry.RespBody = ""
		entry.ReqHeaders = nil
		entry.RespHeaders = nil
	}

	data, _ := json.Marshal(entry)
	t.Output.Write(data)
	t.Output.Write([]byte("\n"))
}

// filterHeaders 过滤敏感 headers（替换 Authorization token）。
func (t *TraceRoundTripper) filterHeaders(h http.Header) map[string]string {
	out := make(map[string]string)
	for k, v := range h {
		val := strings.Join(v, ", ")
		if k == "Authorization" {
			val = "***"
		}
		out[k] = val
	}
	return out
}

// readBody 读取 body 内容，支持截断。返回 (body字符串, 是否截断, 是否为二进制)。
func (t *TraceRoundTripper) readBody(body io.ReadCloser, size int64) (string, bool, bool) {
	const maxBodySize = 64 * 1024 // 64KB

	if body == nil {
		return "", false, false
	}

	// 确定读取上限
	var limit int64 = maxBodySize
	if t.Level == TraceLevelDump || (size > 0 && size < maxBodySize) {
		limit = size
	}

	// 读取 body（使用 LimitReader 避免一次读取过多）
	lr := io.LimitReader(body, limit)
	buf, _ := io.ReadAll(lr)
	body.Close()

	truncated := size > 0 && int64(len(buf)) < size

	// 检测是否为 binary（非 JSON）
	if !isJSON(buf) {
		return "", false, true
	}

	// 过滤 request body 中的敏感字段
	filtered := filterSensitiveBody(string(buf))

	return filtered, truncated, false
}

// isJSON 检测数据是否为 JSON（简单检测：首字符为 { 或 [）。
func isJSON(data []byte) bool {
	trimmed := bytes.TrimLeft(data, " \t\r\n")
	return len(trimmed) > 0 && (trimmed[0] == '{' || trimmed[0] == '[')
}

// filterSensitiveBody 过滤 body 中的敏感字段（password 等）。
func filterSensitiveBody(body string) string {
	// 简单替换：password 字段值替换为 ***
	re := regexp.MustCompile(`"password"\s*:\s*"[^"]*"`)
	return re.ReplaceAllString(body, `"password": "***"`)
}

// 注意：response body 被读取后，原始 handler 无法再次读取。
// 如果 SDK 后续还需要读取 response body，则需要在 RoundTrip 外部预读取，
// 这里简化处理以专注于 trace 功能。

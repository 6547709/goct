// Package session 缓存 CloudTower 登录 token，避免每次命令都重新登录。
//
// 缓存路径：$XDG_CACHE_HOME/goct/session-<sha1(host|user)>.json
//
//	未设置时回退到 ~/.cache/goct/，最终回退到 os.TempDir()/goct/
//
// 文件权限固定 0600，目录 0700，防止 token 泄露。
//
// 过期判定基于 Token.ExpireAt；过期 token 调 Load 会返回错误。
package session

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Token 是缓存的最小载体。
type Token struct {
	Value    string    `json:"value"`
	ExpireAt time.Time `json:"expire_at"`
}

// Expired 报告 token 是否已过期。
func (t Token) Expired() bool { return time.Now().After(t.ExpireAt) }

// PathFor 返回 host+user 对应的 cache 文件绝对路径。
// 同一(host, user)始终映射到同一路径；不同主机/用户互不影响。
func PathFor(host, user string) string {
	sum := sha1.Sum([]byte(host + "|" + user))
	return filepath.Join(cacheDir(), "session-"+hex.EncodeToString(sum[:])+".json")
}

// Save 把 token 持久化到文件（0600 权限）。
func Save(host, user string, tok Token) error {
	p := PathFor(host, user)
	if err := os.MkdirAll(filepath.Dir(p), 0o700); err != nil {
		return fmt.Errorf("mkdir cache dir: %w", err)
	}
	data, err := json.Marshal(tok)
	if err != nil {
		return fmt.Errorf("marshal token: %w", err)
	}
	if err := os.WriteFile(p, data, 0o600); err != nil {
		return fmt.Errorf("write %s: %w", p, err)
	}
	return nil
}

// Load 从 cache 加载 token；
// 文件不存在时返回 os.IsNotExist 可识别的错误；
// 文件存在但已过期返回 ErrExpired 包装的错误；
// JSON 解析失败返回包装的错误（不擦除底层链）。
func Load(host, user string) (Token, error) {
	p := PathFor(host, user)
	data, err := os.ReadFile(p)
	if err != nil {
		return Token{}, err // 保留 os.IsNotExist 语义
	}
	var tok Token
	if err := json.Unmarshal(data, &tok); err != nil {
		return Token{}, fmt.Errorf("decode session %s: %w", p, err)
	}
	if tok.Expired() {
		return Token{}, fmt.Errorf("session token at %s: %w", p, ErrExpired)
	}
	return tok, nil
}

// Delete 移除指定 host/user 的缓存文件，幂等。
func Delete(host, user string) error {
	err := os.Remove(PathFor(host, user))
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("delete session: %w", err)
	}
	return nil
}

// List 返回当前 cache 目录下所有 session 文件的绝对路径。
// 目录不存在视为空列表，不返回错误。
func List() ([]string, error) {
	dir := cacheDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read cache dir %s: %w", dir, err)
	}
	out := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasPrefix(name, "session-") || !strings.HasSuffix(name, ".json") {
			continue
		}
		out = append(out, filepath.Join(dir, name))
	}
	return out, nil
}

// ErrExpired 标识 token 已过期。
var ErrExpired = errors.New("expired")

// cacheDir 选择最合适的 cache 目录，优先 XDG，回退到 ~/.cache，最后回退 TempDir。
func cacheDir() string {
	if x := os.Getenv("XDG_CACHE_HOME"); x != "" {
		return filepath.Join(x, "goct")
	}
	if h, err := os.UserHomeDir(); err == nil && h != "" {
		return filepath.Join(h, ".cache", "goct")
	}
	return filepath.Join(os.TempDir(), "goct")
}

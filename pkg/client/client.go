// Package client 把 adapter.Client 与 session token cache 粘合。
//
// 工作流程：
//  1. 命中本地 session cache → 用 token 构造 adapter.Client（跳过登录）
//  2. token 已过期/不可用 → 删除 cache 后走完整登录
//  3. 登录成功后把新 token 写回 cache（0600）
//
// adapter.Client 通过 With/From 注入到 cobra cmd 的 context，
// 所有命令统一从 ctx 取 client，避免显式传参污染签名。
package client

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/6547709/goct/pkg/adapter"
	"github.com/6547709/goct/pkg/config"
	"github.com/6547709/goct/pkg/session"
)

// New 按"命中 cache → 否则登录 → 写回 cache"流程构造 adapter.Client。
//
// cfg.URL 必须存在；cfg.Username 仅在没有缓存或缓存失效时使用。
func New(ctx context.Context, cfg config.Resolved) (adapter.Client, error) {
	if cfg.URL == "" {
		return nil, errors.New("missing --url / GOCT_URL / config 'url'")
	}
	host := hostKey(cfg.URL)

	// 步骤 1：尝试 session cache
	if cfg.Username != "" {
		if tok, err := session.Load(host, cfg.Username); err == nil {
			c, _, e := adapter.NewClient(ctx, adapter.Options{
				URL:      cfg.URL,
				Insecure: cfg.Insecure,
				Token:    tok.Value,
			})
			if e == nil {
				return c, nil
			}
			// token 不可用（401 或网络问题）→ 清掉缓存，走完整登录
			_ = session.Delete(host, cfg.Username)
		}
	}

	// 步骤 2：完整登录
	if cfg.Username == "" || cfg.Password == "" {
		return nil, errors.New("missing --username / --password (cannot login without credentials)")
	}
	c, tok, err := adapter.NewClient(ctx, adapter.Options{
		URL:      cfg.URL,
		Username: cfg.Username,
		Password: cfg.Password,
		Source:   cfg.Source,
		Insecure: cfg.Insecure,
	})
	if err != nil {
		return nil, fmt.Errorf("login %s: %w", cfg.URL, err)
	}

	// 步骤 3：写回 cache（失败不致命，仅打 trace）
	_ = session.Save(host, cfg.Username, session.Token{
		Value:    tok.Value,
		ExpireAt: tok.ExpireAt,
	})
	return c, nil
}

// hostKey 把 URL 归一到 host[:port]，作为 cache key 的一部分。
// 解析失败时退回原始字符串，保证 key 仍然稳定。
func hostKey(raw string) string {
	if u, err := url.Parse(raw); err == nil && u.Host != "" {
		return u.Host
	}
	return raw
}

// ctxKey 是私有类型，避免与其他包的 context key 冲突。
type ctxKey struct{}

// With 把 client 注入 ctx，供子命令通过 From 提取。
func With(ctx context.Context, c adapter.Client) context.Context {
	return context.WithValue(ctx, ctxKey{}, c)
}

// From 从 ctx 提取 client；未注入时返回 nil。
// 命令层应在使用前判 nil 并返回 user-friendly 错误。
func From(ctx context.Context) adapter.Client {
	c, _ := ctx.Value(ctxKey{}).(adapter.Client)
	return c
}

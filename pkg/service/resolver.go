// Package service 实现资源级别的业务编排逻辑。
//
// service 层坐在 cmd 与 adapter 之间，负责：
// - name|id 统一解析（Resolve）
// - 调 adapter 后 Watch task（对写操作）
// - 串联多步骤逻辑（如 clone 前查 src、migrate 前查目标 host）
package service

import (
	"context"
	"fmt"
	"regexp"

	"github.com/6547709/goct/pkg/adapter"
)

// uuidRe 匹配标准 UUID（36 字符）。
var uuidRe = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// cuidRe 匹配 CloudTower 资源 ID（cuid 格式）。
//
// 格式：c + 小写字母 + 22-24 个小写字母或数字 = 24-26 字符
// 当前已知前缀：cm（25字符，如 cmowkefxp039m0818ld8uli5x）
// 未来可能扩展：cn、ck 等新前缀
//
// 排除：ca（易与 "ca" 开头的人类可读名称混淆）
var cuidRe = regexp.MustCompile(`^c[a-np-z][0-9a-z]{22,24}$`)

// IsID 报告 s 是否看起来像一个 ID（UUID 或 cuid），而非用户可读的名称。
// CloudTower 同时使用两种 ID 格式：
//   - 标准 UUID：5e52cf6e-1e8c-4a0a-9e3a-1b2c3d4e5f6a
//   - cuid：cl5k7g2xo04070822fhxjfsev9q
func IsID(s string) bool {
	return uuidRe.MatchString(s) || cuidRe.MatchString(s)
}

// IsUUID 保留向后兼容。
func IsUUID(s string) bool { return uuidRe.MatchString(s) }

// Resolve 是通用 name|id 解析器。
//
// 解析路径（v0.2.1 优化）：
//  1. 看起来像 ID（UUID 或 cuid）→ 直接 get(id)
//  2. 否则 → list(Name=<exact>) 走服务端精确过滤
//  3. 0 命中 → list(NameContains) 二次模糊匹配做 fallback（兼容部分服务端不支持 Name= 的字段）
//  4. 客户端再做一次精确过滤（防止 NameContains 服务端语义偏宽）
//  5. 0 / 1 / >1 → NotFound / OK / Ambiguous
func Resolve[T any](
	ctx context.Context,
	list func(ctx context.Context, opts adapter.ListOpts) ([]T, error),
	get func(ctx context.Context, id string) (*T, error),
	extract func(T) (id, name string),
	idOrName string,
) (*T, error) {
	if IsID(idOrName) {
		return get(ctx, idOrName)
	}

	// 优先走服务端精确匹配
	candidates, err := list(ctx, adapter.ListOpts{Name: idOrName})
	if err != nil {
		return nil, err
	}

	// fallback：服务端忽略 Name 字段时（部分资源 ListOpts 实现不支持），
	// 再用 NameContains 拉一遍。
	if len(candidates) == 0 {
		candidates, err = list(ctx, adapter.ListOpts{NameContains: idOrName})
		if err != nil {
			return nil, err
		}
	}

	matches := make([]T, 0, 1)
	for _, v := range candidates {
		_, n := extract(v)
		if n == idOrName {
			matches = append(matches, v)
		}
	}
	switch len(matches) {
	case 0:
		return nil, fmt.Errorf("%w: %q", adapter.ErrNotFound, idOrName)
	case 1:
		return &matches[0], nil
	default:
		return nil, fmt.Errorf("ambiguous name %q (%d matches)", idOrName, len(matches))
	}
}

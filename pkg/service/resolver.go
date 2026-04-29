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

// cuidRe 匹配 CloudTower cuid 格式（以 cl 开头，25-27 字符字母数字）。
var cuidRe = regexp.MustCompile(`^cl[0-9a-z]{23,25}$`)

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
// 如果 idOrName 看起来像 ID（UUID 或 cuid）→ 直接 get(id)
// 否则 → list(nameContains) → 精确匹配 → 0=NotFound / 1=OK / >1=Ambiguous
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
	all, err := list(ctx, adapter.ListOpts{NameContains: idOrName})
	if err != nil {
		return nil, err
	}
	matches := make([]T, 0, 1)
	for _, v := range all {
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

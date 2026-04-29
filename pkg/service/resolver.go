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

// uuidRe 匹配标准 UUID（含/不含连字符，36/32 位均可）。
// CloudTower ID 是标准 36 字符 UUID。
var uuidRe = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// IsUUID 报告 s 是否为标准 UUID 格式。
func IsUUID(s string) bool { return uuidRe.MatchString(s) }

// Resolve 是通用 name|id 解析器。
//
// 如果 idOrName 是 UUID → 直接 get(id)
// 否则 → list(nameContains) → 精确匹配 → 0=NotFound / 1=OK / >1=Ambiguous
func Resolve[T any](
	ctx context.Context,
	list func(ctx context.Context, opts adapter.ListOpts) ([]T, error),
	get func(ctx context.Context, id string) (*T, error),
	extract func(T) (id, name string),
	idOrName string,
) (*T, error) {
	if IsUUID(idOrName) {
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

package service

import (
	"context"

	"github.com/6547709/goct/pkg/adapter"
)

// VMService 封装 VM 相关的业务逻辑。
type VMService struct{ c adapter.VMOps }

// NewVM 构造 VMService。c 一般来自 client.From(ctx)（实现 adapter.Client→VMOps）。
func NewVM(c adapter.VMOps) *VMService { return &VMService{c: c} }

// List 列出虚拟机。
func (s *VMService) List(ctx context.Context, opts adapter.ListOpts) ([]adapter.VM, error) {
	return s.c.ListVMs(ctx, opts)
}

// Resolve 按 name 或 UUID 解析到单个 VM。
func (s *VMService) Resolve(ctx context.Context, idOrName string) (*adapter.VM, error) {
	return Resolve(ctx, s.c.ListVMs, s.c.GetVM,
		func(v adapter.VM) (string, string) { return v.ID, v.Name },
		idOrName)
}

// Create 创建 VM，返回 task 引用。
func (s *VMService) Create(ctx context.Context, spec adapter.VMCreateSpec) (adapter.TaskRef, error) {
	return s.c.CreateVM(ctx, spec)
}

// Clone 克隆 VM。srcIDOrName 先 Resolve 再调 adapter。
func (s *VMService) Clone(ctx context.Context, srcIDOrName string, spec adapter.VMCloneSpec) (adapter.TaskRef, error) {
	src, err := s.Resolve(ctx, srcIDOrName)
	if err != nil {
		return adapter.TaskRef{}, err
	}
	return s.c.CloneVM(ctx, src.ID, spec)
}

// Destroy 销毁 VM。
func (s *VMService) Destroy(ctx context.Context, idOrName string, force bool) (adapter.TaskRef, error) {
	v, err := s.Resolve(ctx, idOrName)
	if err != nil {
		return adapter.TaskRef{}, err
	}
	return s.c.DestroyVM(ctx, v.ID, force)
}

// Migrate 迁移 VM 到另一台主机。
func (s *VMService) Migrate(ctx context.Context, vmIDOrName, hostID string) (adapter.TaskRef, error) {
	v, err := s.Resolve(ctx, vmIDOrName)
	if err != nil {
		return adapter.TaskRef{}, err
	}
	return s.c.MigrateVM(ctx, v.ID, hostID)
}

// Export 导出 VM。
func (s *VMService) Export(ctx context.Context, idOrName string, spec adapter.VMExportSpec) (adapter.TaskRef, error) {
	v, err := s.Resolve(ctx, idOrName)
	if err != nil {
		return adapter.TaskRef{}, err
	}
	return s.c.ExportVM(ctx, v.ID, spec)
}

// Power 执行 VM 电源操作。
func (s *VMService) Power(ctx context.Context, idOrName string, action adapter.PowerAction, force bool) (adapter.TaskRef, error) {
	v, err := s.Resolve(ctx, idOrName)
	if err != nil {
		return adapter.TaskRef{}, err
	}
	return s.c.PowerVM(ctx, v.ID, action, force)
}

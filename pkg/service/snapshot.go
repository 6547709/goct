package service

import (
	"context"

	"github.com/6547709/goct/pkg/adapter"
)

// SnapshotService 封装 VM 快照业务逻辑。
type SnapshotService struct {
	c  adapter.SnapshotOps
	vm adapter.VMOps
}

// NewSnapshot 构造 SnapshotService。cli 通常同时实现 VMOps 和 SnapshotOps。
func NewSnapshot(cli interface {
	adapter.SnapshotOps
	adapter.VMOps
}) *SnapshotService {
	return &SnapshotService{c: cli, vm: cli}
}

// List 列出指定 VM 的快照。vmIDOrName 先 resolve 再 list。
func (s *SnapshotService) List(ctx context.Context, vmIDOrName string) ([]adapter.Snapshot, error) {
	v, err := s.resolveVM(ctx, vmIDOrName)
	if err != nil {
		return nil, err
	}
	return s.c.ListSnapshots(ctx, v.ID)
}

// Create 创建快照。
func (s *SnapshotService) Create(ctx context.Context, vmIDOrName, name string) (adapter.TaskRef, error) {
	v, err := s.resolveVM(ctx, vmIDOrName)
	if err != nil {
		return adapter.TaskRef{}, err
	}
	return s.c.CreateSnapshot(ctx, v.ID, name)
}

// Revert 回滚到指定快照。需要 vmID（从快照 list 推断或用户指定）。
func (s *SnapshotService) Revert(ctx context.Context, vmIDOrName, snapID string) (adapter.TaskRef, error) {
	v, err := s.resolveVM(ctx, vmIDOrName)
	if err != nil {
		return adapter.TaskRef{}, err
	}
	return s.c.RevertSnapshot(ctx, v.ID, snapID)
}

// Delete 删除快照（不需要 vmID，按 snapshot ID 即可）。
func (s *SnapshotService) Delete(ctx context.Context, snapID string) (adapter.TaskRef, error) {
	return s.c.DeleteSnapshot(ctx, snapID)
}

func (s *SnapshotService) resolveVM(ctx context.Context, idOrName string) (*adapter.VM, error) {
	return Resolve(ctx, s.vm.ListVMs, s.vm.GetVM,
		func(v adapter.VM) (string, string) { return v.ID, v.Name },
		idOrName)
}

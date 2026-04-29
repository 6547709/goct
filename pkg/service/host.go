package service

import (
	"context"
	"fmt"

	"github.com/6547709/goct/pkg/adapter"
)

// HostService 封装主机相关业务逻辑。
type HostService struct{ c adapter.HostOps }

func NewHost(c adapter.HostOps) *HostService { return &HostService{c: c} }

func (s *HostService) List(ctx context.Context, opts adapter.ListOpts) ([]adapter.Host, error) {
	return s.c.ListHosts(ctx, opts)
}

func (s *HostService) Resolve(ctx context.Context, idOrName string) (*adapter.Host, error) {
	return Resolve(ctx, s.c.ListHosts, s.c.GetHost,
		func(h adapter.Host) (string, string) { return h.ID, h.Name },
		idOrName)
}

func (s *HostService) EnterMaintenance(ctx context.Context, idOrName string) (adapter.TaskRef, error) {
	h, err := s.Resolve(ctx, idOrName)
	if err != nil {
		return adapter.TaskRef{}, err
	}
	return s.c.EnterMaintenanceMode(ctx, h.ID)
}

func (s *HostService) ExitMaintenance(ctx context.Context, idOrName string) (adapter.TaskRef, error) {
	h, err := s.Resolve(ctx, idOrName)
	if err != nil {
		return adapter.TaskRef{}, err
	}
	return s.c.ExitMaintenanceMode(ctx, h.ID)
}

func (s *HostService) Shutdown(ctx context.Context, idOrName string, force bool) (adapter.TaskRef, error) {
	h, err := s.Resolve(ctx, idOrName)
	if err != nil {
		return adapter.TaskRef{}, err
	}
	return s.c.ShutdownHost(ctx, h.ID, force)
}

func (s *HostService) Reboot(ctx context.Context, idOrName string, force bool) (adapter.TaskRef, error) {
	h, err := s.Resolve(ctx, idOrName)
	if err != nil {
		return adapter.TaskRef{}, err
	}
	return s.c.RebootHost(ctx, h.ID, force)
}

// Reconnect 和 Disconnect CloudTower SDK 不支持，返回 ErrUnsupported。
func (s *HostService) Reconnect(_ context.Context, _ string) (adapter.TaskRef, error) {
	return adapter.TaskRef{}, fmt.Errorf("host reconnect: %w", adapter.ErrUnsupported)
}

func (s *HostService) Disconnect(_ context.Context, _ string) (adapter.TaskRef, error) {
	return adapter.TaskRef{}, fmt.Errorf("host disconnect: %w", adapter.ErrUnsupported)
}

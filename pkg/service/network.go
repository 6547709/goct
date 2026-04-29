package service

import (
	"context"
	"github.com/6547709/goct/pkg/adapter"
)

type NetworkService struct{ c adapter.NetworkOps }
func NewNetwork(c adapter.NetworkOps) *NetworkService { return &NetworkService{c: c} }

func (s *NetworkService) List(ctx context.Context, opts adapter.ListOpts) ([]adapter.Network, error) {
	return s.c.ListNetworks(ctx, opts)
}

func (s *NetworkService) Resolve(ctx context.Context, idOrName string) (*adapter.Network, error) {
	return Resolve(ctx, s.c.ListNetworks, s.c.GetNetwork,
		func(n adapter.Network) (string, string) { return n.ID, n.Name },
		idOrName)
}

type VLANService struct{ c adapter.VLANOps }
func NewVLAN(c adapter.VLANOps) *VLANService { return &VLANService{c: c} }

func (s *VLANService) List(ctx context.Context, opts adapter.ListOpts) ([]adapter.VLAN, error) {
	return s.c.ListVLANs(ctx, opts)
}

func (s *VLANService) Resolve(ctx context.Context, idOrName string) (*adapter.VLAN, error) {
	return Resolve(ctx, s.c.ListVLANs, s.c.GetVLAN,
		func(v adapter.VLAN) (string, string) { return v.ID, v.Name },
		idOrName)
}

func (s *VLANService) Create(ctx context.Context, spec adapter.VLANCreateSpec) (adapter.TaskRef, error) {
	return s.c.CreateVLAN(ctx, spec)
}

func (s *VLANService) Delete(ctx context.Context, idOrName string) (adapter.TaskRef, error) {
	v, err := s.Resolve(ctx, idOrName)
	if err != nil { return adapter.TaskRef{}, err }
	return s.c.DeleteVLAN(ctx, v.ID)
}

package service

import (
	"context"
	"github.com/6547709/goct/pkg/adapter"
)

type DatastoreService struct{ c adapter.DatastoreOps }
func NewDatastore(c adapter.DatastoreOps) *DatastoreService { return &DatastoreService{c: c} }

func (s *DatastoreService) List(ctx context.Context, opts adapter.ListOpts) ([]adapter.Datastore, error) {
	return s.c.ListDatastores(ctx, opts)
}

func (s *DatastoreService) Resolve(ctx context.Context, idOrName string) (*adapter.Datastore, error) {
	return Resolve(ctx, s.c.ListDatastores, s.c.GetDatastore,
		func(d adapter.Datastore) (string, string) { return d.ID, d.Name },
		idOrName)
}

func (s *DatastoreService) ListDisks(ctx context.Context, hostID string) ([]adapter.Disk, error) {
	return s.c.ListDisks(ctx, hostID)
}

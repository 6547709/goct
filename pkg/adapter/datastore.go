package adapter

import (
	"context"
	"fmt"

	"github.com/openlyinc/pointy"
	sdkds "github.com/smartxworks/cloudtower-go-sdk/v2/client/elf_data_store"
	sdkdisk "github.com/smartxworks/cloudtower-go-sdk/v2/client/disk"
	"github.com/smartxworks/cloudtower-go-sdk/v2/models"
)

// DatastoreOps 定义数据存储操作。
type DatastoreOps interface {
	ListDatastores(ctx context.Context, opts ListOpts) ([]Datastore, error)
	GetDatastore(ctx context.Context, id string) (*Datastore, error)
	ListDisks(ctx context.Context, hostID string) ([]Disk, error)
}

func (c *defaultClient) ListDatastores(ctx context.Context, opts ListOpts) ([]Datastore, error) {
	params := sdkds.NewGetElfDataStoresParams()
	params.SetContext(ctx)
	body := &models.GetElfDataStoresRequestBody{}
	where := &models.ElfDataStoreWhereInput{}
	hasWhere := false
	if opts.NameContains != "" {
		where.NameContains = pointy.String(opts.NameContains)
		hasWhere = true
	}
	if opts.ClusterID != "" {
		where.Cluster = &models.ClusterWhereInput{ID: pointy.String(opts.ClusterID)}
		hasWhere = true
	}
	if hasWhere {
		body.Where = where
	}
	if opts.Limit > 0 {
		body.First = pointy.Int32(opts.Limit)
	}
	params.SetRequestBody(body)
	resp, err := c.api.ElfDataStore.GetElfDataStores(params)
	if err != nil {
		return nil, fmt.Errorf("list datastores: %w", err)
	}
	out := make([]Datastore, 0, len(resp.Payload))
	for _, d := range resp.Payload {
		out = append(out, toDatastore(d))
	}
	return out, nil
}

func (c *defaultClient) GetDatastore(ctx context.Context, id string) (*Datastore, error) {
	params := sdkds.NewGetElfDataStoresParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.GetElfDataStoresRequestBody{
		Where: &models.ElfDataStoreWhereInput{ID: pointy.String(id)},
	})
	resp, err := c.api.ElfDataStore.GetElfDataStores(params)
	if err != nil {
		return nil, fmt.Errorf("get datastore %s: %w", id, err)
	}
	if len(resp.Payload) == 0 {
		return nil, fmt.Errorf("get datastore %s: %w", id, ErrNotFound)
	}
	d := toDatastore(resp.Payload[0])
	return &d, nil
}

func (c *defaultClient) ListDisks(ctx context.Context, hostID string) ([]Disk, error) {
	params := sdkdisk.NewGetDisksParams()
	params.SetContext(ctx)
	body := &models.GetDisksRequestBody{}
	if hostID != "" {
		body.Where = &models.DiskWhereInput{
			Host: &models.HostWhereInput{ID: pointy.String(hostID)},
		}
	}
	params.SetRequestBody(body)
	resp, err := c.api.Disk.GetDisks(params)
	if err != nil {
		return nil, fmt.Errorf("list disks: %w", err)
	}
	out := make([]Disk, 0, len(resp.Payload))
	for _, d := range resp.Payload {
		out = append(out, toDisk(d))
	}
	return out, nil
}

func toDatastore(d *models.ElfDataStore) Datastore {
	out := Datastore{}
	if d.ID != nil {
		out.ID = *d.ID
	}
	if d.Name != nil {
		out.Name = *d.Name
	}
	if d.Type != nil {
		out.Type = string(*d.Type)
	}
	if d.Internal != nil {
		out.Internal = *d.Internal
	}
	if d.Cluster != nil && d.Cluster.ID != nil {
		out.ClusterID = *d.Cluster.ID
	}
	return out
}

func toDisk(d *models.Disk) Disk {
	out := Disk{}
	if d.ID != nil {
		out.ID = *d.ID
	}
	if d.Name != nil {
		out.Name = *d.Name
	}
	if d.Type != nil {
		out.Type = string(*d.Type)
	}
	if d.Size != nil {
		out.SizeBytes = uint64(*d.Size)
	}
	if d.Path != nil {
		out.Path = *d.Path
	}
	if d.Host != nil && d.Host.Name != nil {
		out.HostName = *d.Host.Name
	}
	return out
}

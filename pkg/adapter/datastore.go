package adapter

import (
	"context"
	"fmt"

	"github.com/openlyinc/pointy"
	sdkds "github.com/smartxworks/cloudtower-go-sdk/v2/client/elf_data_store"
	sdkdisk "github.com/smartxworks/cloudtower-go-sdk/v2/client/disk"
	sdkdiskpool "github.com/smartxworks/cloudtower-go-sdk/v2/client/disk_pool"
	"github.com/smartxworks/cloudtower-go-sdk/v2/models"
)

// DatastoreOps 定义数据存储操作。
type DatastoreOps interface {
	ListDatastores(ctx context.Context, opts ListOpts) ([]Datastore, error)
	GetDatastore(ctx context.Context, id string) (*Datastore, error)
	ListDisks(ctx context.Context, hostID string) ([]Disk, error)
	ListDiskPools(ctx context.Context, opts ListOpts) ([]DiskPool, error)
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

func (c *defaultClient) ListDiskPools(ctx context.Context, opts ListOpts) ([]DiskPool, error) {
	params := sdkdiskpool.NewGetDiskPoolsParams()
	params.SetContext(ctx)
	body := &models.GetDiskPoolsRequestBody{}
	if opts.ClusterID != "" {
		where := &models.DiskPoolWhereInput{
			Host: &models.HostWhereInput{Cluster: &models.ClusterWhereInput{ID: pointy.String(opts.ClusterID)}},
		}
		body.Where = where
	}
	if opts.Limit > 0 {
		body.First = pointy.Int32(opts.Limit)
	}
	params.SetRequestBody(body)
	resp, err := c.api.DiskPool.GetDiskPools(params)
	if err != nil {
		return nil, fmt.Errorf("list disk pools: %w", err)
	}
	out := make([]DiskPool, 0, len(resp.Payload))
	for _, dp := range resp.Payload {
		out = append(out, toDiskPool(dp))
	}
	return out, nil
}

func toDiskPool(dp *models.DiskPool) DiskPool {
	out := DiskPool{}
	if dp.ID != nil {
		out.ID = *dp.ID
	}
	if dp.Status != nil {
		out.Status = string(*dp.Status)
	}
	if dp.UseState != nil {
		out.UseState = string(*dp.UseState)
	}
	if dp.TotalDataCapacity != nil {
		out.TotalDataBytes = uint64(*dp.TotalDataCapacity)
	}
	if dp.UsedDataSpace != nil {
		out.UsedDataBytes = uint64(*dp.UsedDataSpace)
	}
	if dp.TotalCacheCapacity != nil {
		out.TotalCacheBytes = uint64(*dp.TotalCacheCapacity)
	}
	if dp.UsedCacheSpace != nil {
		out.UsedCacheBytes = uint64(*dp.UsedCacheSpace)
	}
	if dp.HddDiskCount != nil {
		out.HddCount = *dp.HddDiskCount
	}
	if dp.NvmeSsdDiskCount != nil {
		out.NvmeCount = *dp.NvmeSsdDiskCount
	}
	if dp.SataOrSasSsdDiskCount != nil {
		out.SataCount = *dp.SataOrSasSsdDiskCount
	}
	if dp.Host != nil {
		if dp.Host.ID != nil {
			out.HostID = *dp.Host.ID
		}
		if dp.Host.Name != nil {
			out.HostName = *dp.Host.Name
		}
	}
	return out
}

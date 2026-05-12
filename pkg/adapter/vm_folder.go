package adapter

import (
	"context"
	"fmt"

	"github.com/openlyinc/pointy"
	"github.com/smartxworks/cloudtower-go-sdk/v2/client/vm_folder"
	"github.com/smartxworks/cloudtower-go-sdk/v2/models"
)

// VMFolderOps 定义 VM 文件夹操作。
type VMFolderOps interface {
	ListVMFolders(ctx context.Context, opts ListOpts) ([]VMFolder, error)
	CreateVMFolder(ctx context.Context, spec VMFolderCreateSpec) ([]VMFolder, error)
	UpdateVMFolder(ctx context.Context, id string, spec VMFolderUpdateSpec) ([]VMFolder, error)
	DeleteVMFolder(ctx context.Context, id string) error
}

// ---------- ListVMFolders ----------

func (c *defaultClient) ListVMFolders(ctx context.Context, opts ListOpts) ([]VMFolder, error) {
	params := vm_folder.NewGetVMFoldersParams()
	params.SetContext(ctx)

	where := &models.VMFolderWhereInput{}
	hasWhere := false
	if opts.NameContains != "" {
		where.NameContains = pointy.String(opts.NameContains)
		hasWhere = true
	}
	if opts.ClusterID != "" {
		where.Cluster = &models.ClusterWhereInput{ID: pointy.String(opts.ClusterID)}
		hasWhere = true
	}

	body := &models.GetVMFoldersRequestBody{}
	if hasWhere {
		body.Where = where
	}
	if opts.Limit > 0 {
		body.First = pointy.Int32(opts.Limit)
	}
	if opts.Skip > 0 {
		body.Skip = pointy.Int32(opts.Skip)
	}
	params.SetRequestBody(body)

	resp, err := c.api.VMFolder.GetVMFolders(params)
	if err != nil {
		return nil, fmt.Errorf("list vm folders: %w", err)
	}
	out := make([]VMFolder, 0, len(resp.Payload))
	for _, f := range resp.Payload {
		out = append(out, toVMFolder(f))
	}
	return out, nil
}

// ---------- CreateVMFolder ----------

func (c *defaultClient) CreateVMFolder(ctx context.Context, spec VMFolderCreateSpec) ([]VMFolder, error) {
	params := vm_folder.NewCreateVMFolderParams()
	params.SetContext(ctx)
	params.SetRequestBody([]*models.VMFolderCreationParams{
		{
			ClusterID: pointy.String(spec.ClusterID),
			Name:      pointy.String(spec.Name),
		},
	})

	resp, err := c.api.VMFolder.CreateVMFolder(params)
	if err != nil {
		return nil, fmt.Errorf("create vm folder: %w", err)
	}
	out := make([]VMFolder, 0, len(resp.Payload))
	for _, f := range resp.Payload {
		if f.Data != nil {
			out = append(out, toVMFolder(f.Data))
		}
	}
	return out, nil
}

// ---------- UpdateVMFolder ----------

func (c *defaultClient) UpdateVMFolder(ctx context.Context, id string, spec VMFolderUpdateSpec) ([]VMFolder, error) {
	data := &models.VMFolderUpdationParamsData{}
	if spec.Name != "" {
		data.Name = pointy.String(spec.Name)
	}

	params := vm_folder.NewUpdateVMFolderParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.VMFolderUpdationParams{
		Where: &models.VMFolderWhereInput{ID: pointy.String(id)},
		Data:  data,
	})

	resp, err := c.api.VMFolder.UpdateVMFolder(params)
	if err != nil {
		return nil, fmt.Errorf("update vm folder: %w", err)
	}
	out := make([]VMFolder, 0, len(resp.Payload))
	for _, f := range resp.Payload {
		if f.Data != nil {
			out = append(out, toVMFolder(f.Data))
		}
	}
	return out, nil
}

// ---------- DeleteVMFolder ----------

func (c *defaultClient) DeleteVMFolder(ctx context.Context, id string) error {
	params := vm_folder.NewDeleteVMFolderParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.VMFolderDeletionParams{
		Where: &models.VMFolderWhereInput{ID: pointy.String(id)},
	})

	_, err := c.api.VMFolder.DeleteVMFolder(params)
	if err != nil {
		return fmt.Errorf("delete vm folder: %w", err)
	}
	return nil
}

// toVMFolder 把 SDK models.VMFolder 转成内部 VMFolder 模型。
func toVMFolder(f *models.VMFolder) VMFolder {
	out := VMFolder{}
	if f.ID != nil {
		out.ID = *f.ID
	}
	if f.Name != nil {
		out.Name = *f.Name
	}
	if f.Cluster != nil && f.Cluster.ID != nil {
		out.ClusterID = *f.Cluster.ID
	}
	return out
}
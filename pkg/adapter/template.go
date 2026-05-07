package adapter

import (
	"context"
	"fmt"

	"github.com/openlyinc/pointy"
	"github.com/smartxworks/cloudtower-go-sdk/v2/client/content_library_vm_template"
	"github.com/smartxworks/cloudtower-go-sdk/v2/models"
)

// TemplateOps 定义内容库模板操作。
type TemplateOps interface {
	ListContentLibraryTemplates(ctx context.Context, opts ListOpts) ([]ContentLibraryTemplate, error)
	GetContentLibraryTemplateByName(ctx context.Context, name string) (*ContentLibraryTemplate, error)
}

// ContentLibraryTemplate 是 CLI 内部用的内容库模板视图。
type ContentLibraryTemplate struct {
	ID                 string
	Name               string
	Description        string
	VMID               string
	CloudInitSupported bool
}

func (c *defaultClient) ListContentLibraryTemplates(ctx context.Context, opts ListOpts) ([]ContentLibraryTemplate, error) {
	params := content_library_vm_template.NewGetContentLibraryVMTemplatesParams()
	params.SetContext(ctx)
	body := &models.GetContentLibraryVMTemplatesRequestBody{}
	where := &models.ContentLibraryVMTemplateWhereInput{}
	hasWhere := false
	if opts.NameContains != "" {
		where.Name = pointy.String(opts.NameContains)
		hasWhere = true
	}
	if hasWhere {
		body.Where = where
	}
	if opts.Limit > 0 {
		body.First = pointy.Int32(opts.Limit)
	}
	params.SetRequestBody(body)
	resp, err := c.api.ContentLibraryVMTemplate.GetContentLibraryVMTemplates(params)
	if err != nil {
		return nil, fmt.Errorf("list content library templates: %w", err)
	}
	out := make([]ContentLibraryTemplate, 0, len(resp.Payload))
	for _, t := range resp.Payload {
		out = append(out, toContentLibraryTemplate(t))
	}
	return out, nil
}

func (c *defaultClient) GetContentLibraryTemplateByName(ctx context.Context, name string) (*ContentLibraryTemplate, error) {
	params := content_library_vm_template.NewGetContentLibraryVMTemplatesParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.GetContentLibraryVMTemplatesRequestBody{
		Where: &models.ContentLibraryVMTemplateWhereInput{
			Name: &name,
		},
	})
	resp, err := c.api.ContentLibraryVMTemplate.GetContentLibraryVMTemplates(params)
	if err != nil {
		return nil, fmt.Errorf("get content library template %s: %w", name, err)
	}
	if len(resp.Payload) == 0 {
		return nil, fmt.Errorf("content library template not found: %s: %w", name, ErrNotFound)
	}
	t := toContentLibraryTemplate(resp.Payload[0])
	return &t, nil
}

func toContentLibraryTemplate(t *models.ContentLibraryVMTemplate) ContentLibraryTemplate {
	out := ContentLibraryTemplate{}
	if t.ID != nil {
		out.ID = *t.ID
	}
	if t.Name != nil {
		out.Name = *t.Name
	}
	if t.Description != nil {
		out.Description = *t.Description
	}
	if t.CloudInitSupported != nil {
		out.CloudInitSupported = *t.CloudInitSupported
	}
	// 内层 vm_templates[0].id 是真正的模板 ID
	if t.VMTemplates != nil && len(t.VMTemplates) > 0 && t.VMTemplates[0].ID != nil {
		out.VMID = *t.VMTemplates[0].ID
	}
	return out
}
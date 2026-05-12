package adapter

import (
	"context"
	"fmt"

	"github.com/openlyinc/pointy"
	"github.com/smartxworks/cloudtower-go-sdk/v2/client/label"
	"github.com/smartxworks/cloudtower-go-sdk/v2/models"
)

// LabelOps 定义标签相关的 SDK 操作。
type LabelOps interface {
	ListLabels(ctx context.Context, opts ListOpts) ([]Label, error)
	CreateLabel(ctx context.Context, spec LabelCreateSpec) ([]Label, error)
	UpdateLabel(ctx context.Context, id string, spec LabelUpdateSpec) ([]Label, error)
	DeleteLabel(ctx context.Context, id string) error
	AttachLabel(ctx context.Context, labelID string, spec LabelAttachSpec) error
	DetachLabel(ctx context.Context, labelID string, spec LabelDetachSpec) error
}

// ---------- ListLabels ----------

func (c *defaultClient) ListLabels(ctx context.Context, opts ListOpts) ([]Label, error) {
	params := label.NewGetLabelsParams()
	params.SetContext(ctx)

	where := &models.LabelWhereInput{}
	hasWhere := false
	if opts.NameContains != "" {
		where.KeyContains = pointy.String(opts.NameContains)
		hasWhere = true
	}

	body := &models.GetLabelsRequestBody{}
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

	resp, err := c.api.Label.GetLabels(params)
	if err != nil {
		return nil, fmt.Errorf("list labels: %w", err)
	}
	out := make([]Label, 0, len(resp.Payload))
	for _, l := range resp.Payload {
		out = append(out, toLabel(l))
	}
	return out, nil
}

// ---------- CreateLabel ----------

func (c *defaultClient) CreateLabel(ctx context.Context, spec LabelCreateSpec) ([]Label, error) {
	params := label.NewCreateLabelParams()
	params.SetContext(ctx)
	params.SetRequestBody([]*models.LabelCreationParams{
		{
			Key:   pointy.String(spec.Key),
			Value: pointy.String(spec.Value),
		},
	})

	resp, err := c.api.Label.CreateLabel(params)
	if err != nil {
		return nil, fmt.Errorf("create label: %w", err)
	}
	out := make([]Label, 0, len(resp.Payload))
	for _, l := range resp.Payload {
		if l.Data != nil {
			out = append(out, toLabel(l.Data))
		}
	}
	return out, nil
}

// ---------- UpdateLabel ----------

func (c *defaultClient) UpdateLabel(ctx context.Context, id string, spec LabelUpdateSpec) ([]Label, error) {
	data := &models.LabelUpdationParamsData{}
	if spec.Key != "" {
		data.Key = pointy.String(spec.Key)
	}
	if spec.Value != "" {
		data.Value = pointy.String(spec.Value)
	}

	params := label.NewUpdateLabelParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.LabelUpdationParams{
		Where: &models.LabelWhereInput{ID: pointy.String(id)},
		Data:  data,
	})

	resp, err := c.api.Label.UpdateLabel(params)
	if err != nil {
		return nil, fmt.Errorf("update label: %w", err)
	}
	out := make([]Label, 0, len(resp.Payload))
	for _, l := range resp.Payload {
		if l.Data != nil {
			out = append(out, toLabel(l.Data))
		}
	}
	return out, nil
}

// ---------- DeleteLabel ----------

func (c *defaultClient) DeleteLabel(ctx context.Context, id string) error {
	params := label.NewDeleteLabelParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.LabelDeletionParams{
		Where: &models.LabelWhereInput{ID: pointy.String(id)},
	})

	_, err := c.api.Label.DeleteLabel(params)
	if err != nil {
		return fmt.Errorf("delete label: %w", err)
	}
	return nil
}

// ---------- AttachLabel ----------

func (c *defaultClient) AttachLabel(ctx context.Context, labelID string, spec LabelAttachSpec) error {
	params := label.NewAddLabelsToResourcesParams()
	params.SetContext(ctx)

	data := buildLabelResourceData(spec.ResourceKind, spec.ResourceID)
	params.SetRequestBody(&models.AddLabelsToResourcesParams{
		Data:  data,
		Where: &models.LabelWhereInput{ID: pointy.String(labelID)},
	})

	_, err := c.api.Label.AddLabelsToResources(params)
	if err != nil {
		return fmt.Errorf("attach label: %w", err)
	}
	return nil
}

// ---------- DetachLabel ----------

func (c *defaultClient) DetachLabel(ctx context.Context, labelID string, spec LabelDetachSpec) error {
	params := label.NewRemoveLabelsFromResourcesParams()
	params.SetContext(ctx)

	data := buildLabelResourceData(spec.ResourceKind, spec.ResourceID)
	params.SetRequestBody(&models.RemoveLabelsFromResourcesParams{
		AddLabelsToResourcesParams: models.AddLabelsToResourcesParams{
			Data:  data,
			Where: &models.LabelWhereInput{ID: pointy.String(labelID)},
		},
	})

	_, err := c.api.Label.RemoveLabelsFromResources(params)
	if err != nil {
		return fmt.Errorf("detach label: %w", err)
	}
	return nil
}

// buildLabelResourceData 根据资源类型构建标签资源数据。
func buildLabelResourceData(kind, id string) *models.AddLabelsToResourcesParamsData {
	data := &models.AddLabelsToResourcesParamsData{}
	switch kind {
	case "vm":
		data.Vms = &models.VMWhereInput{ID: pointy.String(id)}
	case "host":
		data.Hosts = &models.HostWhereInput{ID: pointy.String(id)}
	case "cluster":
		data.Clusters = &models.ClusterWhereInput{ID: pointy.String(id)}
	case "volume":
		data.VMVolumes = &models.VMVolumeWhereInput{ID: pointy.String(id)}
	default:
		// 其他资源类型暂不支持
	}
	return data
}

// toLabel 把 SDK models.Label 转成内部 Label 模型。
func toLabel(l *models.Label) Label {
	out := Label{}
	if l.ID != nil {
		out.ID = *l.ID
	}
	if l.Key != nil {
		out.Key = *l.Key
	}
	if l.Value != nil {
		out.Value = *l.Value
	}
	if l.CreatedAt != nil {
		out.CreatedAt = *l.CreatedAt
	}
	return out
}
package adapter

import (
	"context"
	"fmt"

	"github.com/openlyinc/pointy"
	"github.com/smartxworks/cloudtower-go-sdk/v2/client/license"
	"github.com/smartxworks/cloudtower-go-sdk/v2/models"
)

// LicenseOps 定义许可证操作。
type LicenseOps interface {
	ListLicenses(ctx context.Context, opts ListOpts) ([]License, error)
	UpdateDeploy(ctx context.Context, licenseKey string) (TaskRef, error)
}

// ---------- ListLicenses ----------

func (c *defaultClient) ListLicenses(ctx context.Context, opts ListOpts) ([]License, error) {
	params := license.NewGetLicensesParams()
	params.SetContext(ctx)

	body := &models.GetLicensesRequestBody{}
	if opts.Limit > 0 {
		body.First = pointy.Int32(opts.Limit)
	}
	if opts.Skip > 0 {
		body.Skip = pointy.Int32(opts.Skip)
	}
	params.SetRequestBody(body)

	resp, err := c.api.License.GetLicenses(params)
	if err != nil {
		return nil, fmt.Errorf("list licenses: %w", err)
	}
	out := make([]License, 0, len(resp.Payload))
	for _, l := range resp.Payload {
		out = append(out, toLicense(l))
	}
	return out, nil
}

// ---------- UpdateDeploy ----------

func (c *defaultClient) UpdateDeploy(ctx context.Context, licenseKey string) (TaskRef, error) {
	params := license.NewUpdateDeployParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.LicenseUpdationParams{
		Data: &models.LicenseUpdationParamsData{
			License: pointy.String(licenseKey),
		},
	})

	resp, err := c.api.License.UpdateDeploy(params)
	if err != nil {
		return TaskRef{}, fmt.Errorf("update deploy: %w", err)
	}
	if resp.Payload == nil {
		return TaskRef{}, nil
	}
	if resp.Payload.TaskID == nil {
		return TaskRef{}, nil
	}
	return TaskRef{ID: *resp.Payload.TaskID}, nil
}

// toLicense 把 SDK models.License 转成内部 License 模型。
func toLicense(l *models.License) License {
	out := License{}
	if l.ID != nil {
		out.ID = *l.ID
	}
	if l.ExpireDate != nil {
		out.ExpireDate = *l.ExpireDate
	}
	if l.LicenseSerial != nil {
		out.LicenseSerial = *l.LicenseSerial
	}
	if l.MaintenanceEndDate != nil {
		out.MaintenanceEndDate = *l.MaintenanceEndDate
	}
	if l.MaintenanceStartDate != nil {
		out.MaintenanceStartDate = *l.MaintenanceStartDate
	}
	if l.MaxChunkNum != nil {
		out.MaxChunkNum = *l.MaxChunkNum
	}
	if l.MaxClusterNum != nil {
		out.MaxClusterNum = *l.MaxClusterNum
	}
	if l.SignDate != nil {
		out.SignDate = *l.SignDate
	}
	if l.SoftwareEdition != nil {
		out.SoftwareEdition = string(*l.SoftwareEdition)
	}
	if l.Type != nil {
		out.Type = string(*l.Type)
	}
	return out
}

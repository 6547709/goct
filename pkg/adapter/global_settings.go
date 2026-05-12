package adapter

import (
	"context"
	"fmt"

	"github.com/openlyinc/pointy"
	"github.com/smartxworks/cloudtower-go-sdk/v2/client/global_settings"
	"github.com/smartxworks/cloudtower-go-sdk/v2/models"
)

// GlobalSettingsOps 定义全局设置操作。
type GlobalSettingsOps interface {
	GetGlobalSettings(ctx context.Context) (*GlobalSettings, error)
	UpdateSessionTimeout(ctx context.Context, timeout int32) error
}

// ---------- GetGlobalSettings ----------

func (c *defaultClient) GetGlobalSettings(ctx context.Context) (*GlobalSettings, error) {
	params := global_settings.NewGetGlobalSettingsesParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.GetGlobalSettingsesRequestBody{
		First: pointy.Int32(1),
	})

	resp, err := c.api.GlobalSettings.GetGlobalSettingses(params)
	if err != nil {
		return nil, fmt.Errorf("get global settings: %w", err)
	}
	if len(resp.Payload) == 0 {
		return nil, fmt.Errorf("global settings not found")
	}
	return toGlobalSettings(resp.Payload[0]), nil
}

// ---------- UpdateSessionTimeout ----------

func (c *defaultClient) UpdateSessionTimeout(ctx context.Context, timeout int32) error {
	params := global_settings.NewUpdateSessionTimeoutParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.UpdateSessionTimeoutParams{
		SessionMaxAge: pointy.Int32(timeout),
	})

	_, err := c.api.GlobalSettings.UpdateSessionTimeout(params)
	if err != nil {
		return fmt.Errorf("update session timeout: %w", err)
	}
	return nil
}

// toGlobalSettings 把 SDK models.GlobalSettings 转成内部 GlobalSettings 模型。
func toGlobalSettings(p *models.GlobalSettings) *GlobalSettings {
	out := &GlobalSettings{}
	if p.ID != nil {
		out.ID = *p.ID
	}
	if p.Auth != nil && p.Auth.SessionMaxAge != nil {
		out.SessionMaxAge = *p.Auth.SessionMaxAge
	}
	if p.VMRecycleBin != nil {
		if p.VMRecycleBin.Retain != nil {
			out.VMRecycleBin.RetainPeriod = *p.VMRecycleBin.Retain
		}
		if p.VMRecycleBin.Enabled != nil {
			out.VMRecycleBin.Enabled = *p.VMRecycleBin.Enabled
		}
	}
	return out
}

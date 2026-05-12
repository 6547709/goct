package adapter

import (
	"context"
	"fmt"

	"github.com/smartxworks/cloudtower-go-sdk/v2/client/ntp"
)

// NtpOps 定义 NTP 操作。
type NtpOps interface {
	GetNtpServiceURL(ctx context.Context) (*NtpSettings, error)
}

// ---------- GetNtpServiceURL ----------

func (c *defaultClient) GetNtpServiceURL(ctx context.Context) (*NtpSettings, error) {
	params := ntp.NewGetNtpServiceURLParams()
	params.SetContext(ctx)

	resp, err := c.api.Ntp.GetNtpServiceURL(params)
	if err != nil {
		return nil, fmt.Errorf("get ntp service url: %w", err)
	}
	out := &NtpSettings{}
	if resp.Payload != nil && resp.Payload.NtpServiceURL != nil {
		out.URLs = append(out.URLs, *resp.Payload.NtpServiceURL)
	}
	return out, nil
}

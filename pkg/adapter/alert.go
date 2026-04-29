package adapter

import (
	"context"
	"fmt"

	"github.com/openlyinc/pointy"
	sdkalert "github.com/smartxworks/cloudtower-go-sdk/v2/client/alert"
	"github.com/smartxworks/cloudtower-go-sdk/v2/models"
)

// AlertOps 定义告警操作。
type AlertOps interface {
	ListAlerts(ctx context.Context, opts ListOpts) ([]Alert, error)
	GetAlert(ctx context.Context, id string) (*Alert, error)
	AckAlert(ctx context.Context, id string) (TaskRef, error)
}

func (c *defaultClient) ListAlerts(ctx context.Context, opts ListOpts) ([]Alert, error) {
	params := sdkalert.NewGetAlertsParams()
	params.SetContext(ctx)
	body := &models.GetAlertsRequestBody{}
	if opts.Limit > 0 {
		body.First = pointy.Int32(opts.Limit)
	}
	params.SetRequestBody(body)
	resp, err := c.api.Alert.GetAlerts(params)
	if err != nil {
		return nil, fmt.Errorf("list alerts: %w", err)
	}
	out := make([]Alert, 0, len(resp.Payload))
	for _, a := range resp.Payload {
		out = append(out, toAlert(a))
	}
	return out, nil
}

func (c *defaultClient) GetAlert(ctx context.Context, id string) (*Alert, error) {
	params := sdkalert.NewGetAlertsParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.GetAlertsRequestBody{
		Where: &models.AlertWhereInput{ID: pointy.String(id)},
	})
	resp, err := c.api.Alert.GetAlerts(params)
	if err != nil {
		return nil, fmt.Errorf("get alert %s: %w", id, err)
	}
	if len(resp.Payload) == 0 {
		return nil, fmt.Errorf("get alert %s: %w", id, ErrNotFound)
	}
	a := toAlert(resp.Payload[0])
	return &a, nil
}

func (c *defaultClient) AckAlert(ctx context.Context, id string) (TaskRef, error) {
	params := sdkalert.NewResolveAlertParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.ResolveAlertParams{
		Where: &models.AlertWhereInput{ID: pointy.String(id)},
	})
	resp, err := c.api.Alert.ResolveAlert(params)
	if err != nil {
		return TaskRef{}, fmt.Errorf("ack alert %s: %w", id, err)
	}
	return firstAlertTaskRef(resp.Payload), nil
}

func toAlert(a *models.Alert) Alert {
	out := Alert{}
	if a.ID != nil { out.ID = *a.ID }
	if a.Message != nil { out.Message = *a.Message }
	if a.Severity != nil { out.Severity = *a.Severity }
	if a.Cause != nil { out.Cause = *a.Cause }
	return out
}

func firstAlertTaskRef(items []*models.WithTaskAlert) TaskRef {
	if len(items) == 0 { return TaskRef{} }
	ref := TaskRef{EntityKind: "Alert"}
	if items[0].TaskID != nil { ref.ID = *items[0].TaskID }
	if items[0].Data != nil && items[0].Data.ID != nil { ref.EntityID = *items[0].Data.ID }
	return ref
}

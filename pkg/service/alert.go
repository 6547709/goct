package service

import (
	"context"
	"github.com/6547709/goct/pkg/adapter"
)

type AlertService struct{ c adapter.AlertOps }
func NewAlert(c adapter.AlertOps) *AlertService { return &AlertService{c: c} }

func (s *AlertService) List(ctx context.Context, opts adapter.ListOpts) ([]adapter.Alert, error) {
	return s.c.ListAlerts(ctx, opts)
}

func (s *AlertService) Get(ctx context.Context, id string) (*adapter.Alert, error) {
	return s.c.GetAlert(ctx, id)
}

func (s *AlertService) Ack(ctx context.Context, id string) (adapter.TaskRef, error) {
	return s.c.AckAlert(ctx, id)
}

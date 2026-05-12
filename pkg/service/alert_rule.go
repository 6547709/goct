package service

import (
	"context"

	"github.com/6547709/goct/pkg/adapter"
)

type AlertRuleService struct{ c adapter.AlertRuleOps }
func NewAlertRule(c adapter.AlertRuleOps) *AlertRuleService { return &AlertRuleService{c: c} }

func (s *AlertRuleService) List(ctx context.Context, opts adapter.ListOpts) ([]adapter.AlertRule, error) {
	return s.c.ListAlertRules(ctx, opts)
}
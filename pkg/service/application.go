package service

import (
	"context"

	"github.com/6547709/goct/pkg/adapter"
)

type ApplicationService struct{ c adapter.ApplicationOps }
func NewApplication(c adapter.ApplicationOps) *ApplicationService { return &ApplicationService{c: c} }

func (s *ApplicationService) List(ctx context.Context, opts adapter.ListOpts) ([]adapter.Application, error) {
	return s.c.ListApplications(ctx, opts)
}
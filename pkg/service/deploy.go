package service

import (
	"context"

	"github.com/6547709/goct/pkg/adapter"
)

type DeployService struct{ c adapter.DeployOps }
func NewDeploy(c adapter.DeployOps) *DeployService { return &DeployService{c: c} }

func (s *DeployService) List(ctx context.Context, opts adapter.ListOpts) ([]adapter.Deploy, error) {
	return s.c.ListDeploys(ctx, opts)
}
package service

import (
	"context"

	"github.com/6547709/goct/pkg/adapter"
)

type LicenseService struct{ c adapter.LicenseOps }
func NewLicense(c adapter.LicenseOps) *LicenseService { return &LicenseService{c: c} }

func (s *LicenseService) List(ctx context.Context, opts adapter.ListOpts) ([]adapter.License, error) {
	return s.c.ListLicenses(ctx, opts)
}

func (s *LicenseService) UpdateDeploy(ctx context.Context, licenseKey string) (adapter.TaskRef, error) {
	return s.c.UpdateDeploy(ctx, licenseKey)
}
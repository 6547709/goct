package service

import (
	"context"

	"github.com/6547709/goct/pkg/adapter"
	"github.com/smartxworks/cloudtower-go-sdk/v2/models"
)

type CloudTowerApplicationService struct{ c adapter.CloudTowerApplicationOps }
func NewCloudTowerApplication(c adapter.CloudTowerApplicationOps) *CloudTowerApplicationService { return &CloudTowerApplicationService{c: c} }

func (s *CloudTowerApplicationService) List(ctx context.Context, opts adapter.ListOpts) ([]adapter.CloudTowerApplication, error) {
	return s.c.ListCloudTowerApplications(ctx, opts)
}

func (s *CloudTowerApplicationService) ListPackages(ctx context.Context, opts adapter.ListOpts) ([]adapter.CloudTowerApplicationPackage, error) {
	return s.c.ListCloudTowerApplicationPackages(ctx, opts)
}

func (s *CloudTowerApplicationService) UploadPackage(ctx context.Context, path string, name string) (adapter.TaskRef, error) {
	return s.c.UploadCloudTowerApplicationPackage(ctx, path, name)
}

func (s *CloudTowerApplicationService) DeletePackage(ctx context.Context, id string) error {
	return s.c.DeleteCloudTowerApplicationPackage(ctx, id)
}

func (s *CloudTowerApplicationService) Deploy(ctx context.Context, name string, targetPackage string, vmSpec *models.ApplicationVMSpecDefinition) (adapter.TaskRef, error) {
	return s.c.DeployCloudTowerApplication(ctx, name, targetPackage, vmSpec)
}

package service

import (
	"context"

	"github.com/6547709/goct/pkg/adapter"
)

type ContentLibraryImageService struct{ c adapter.ContentLibraryImageOps }
func NewContentLibraryImage(c adapter.ContentLibraryImageOps) *ContentLibraryImageService { return &ContentLibraryImageService{c: c} }

func (s *ContentLibraryImageService) List(ctx context.Context, opts adapter.ListOpts) ([]adapter.ContentLibraryImage, error) {
	return s.c.ListContentLibraryImages(ctx, opts)
}

func (s *ContentLibraryImageService) Delete(ctx context.Context, id string) error {
	return s.c.DeleteContentLibraryImage(ctx, id)
}

func (s *ContentLibraryImageService) Import(ctx context.Context, path string, name string, clusterID string) (adapter.TaskRef, error) {
	return s.c.ImportContentLibraryImage(ctx, path, name, clusterID)
}

func (s *ContentLibraryImageService) Distribute(ctx context.Context, id string, clusterIDs []string) (adapter.TaskRef, error) {
	return s.c.DistributeContentLibraryImage(ctx, id, clusterIDs)
}
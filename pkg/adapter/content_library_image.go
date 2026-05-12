package adapter

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/go-openapi/runtime"
	"github.com/openlyinc/pointy"
	"github.com/smartxworks/cloudtower-go-sdk/v2/client/content_library_image"
	"github.com/smartxworks/cloudtower-go-sdk/v2/models"
)

const (
	uploadChunkSize = 8 * 1024 * 1024 // 8MB
)

// ContentLibraryImageOps 定义内容库镜像操作。
type ContentLibraryImageOps interface {
	ListContentLibraryImages(ctx context.Context, opts ListOpts) ([]ContentLibraryImage, error)
	DeleteContentLibraryImage(ctx context.Context, id string) error
	ImportContentLibraryImage(ctx context.Context, path string, name string, clusterID string) (TaskRef, error)
	DistributeContentLibraryImage(ctx context.Context, id string, clusterIDs []string) (TaskRef, error)
}

// ---------- ListContentLibraryImages ----------

func (c *defaultClient) ListContentLibraryImages(ctx context.Context, opts ListOpts) ([]ContentLibraryImage, error) {
	params := content_library_image.NewGetContentLibraryImagesParams()
	params.SetContext(ctx)

	where := &models.ContentLibraryImageWhereInput{}
	hasWhere := false
	if opts.NameContains != "" {
		where.NameContains = pointy.String(opts.NameContains)
		hasWhere = true
	}

	body := &models.GetContentLibraryImagesRequestBody{}
	if hasWhere {
		body.Where = where
	}
	if opts.Limit > 0 {
		body.First = pointy.Int32(opts.Limit)
	}
	if opts.Skip > 0 {
		body.Skip = pointy.Int32(opts.Skip)
	}
	params.SetRequestBody(body)

	resp, err := c.api.ContentLibraryImage.GetContentLibraryImages(params)
	if err != nil {
		return nil, fmt.Errorf("list content library images: %w", err)
	}
	out := make([]ContentLibraryImage, 0, len(resp.Payload))
	for _, img := range resp.Payload {
		out = append(out, toContentLibraryImage(img))
	}
	return out, nil
}

// ---------- DeleteContentLibraryImage ----------

func (c *defaultClient) DeleteContentLibraryImage(ctx context.Context, id string) error {
	params := content_library_image.NewDeleteContentLibraryImageParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.ContentLibraryImageDeletionParams{
		Where: &models.ContentLibraryImageWhereInput{ID: pointy.String(id)},
	})

	_, err := c.api.ContentLibraryImage.DeleteContentLibraryImage(params)
	if err != nil {
		return fmt.Errorf("delete content library image: %w", err)
	}
	return nil
}

// cloudTowerUploadWriter 实现 io.Writer，用于分片上传文件到 CloudTower。
type cloudTowerUploadWriter struct {
	pos          int
	uploadTaskID string
	name         string
	size         int64
	clusterID    string
	client       content_library_image.ClientService
}

func (w *cloudTowerUploadWriter) Write(p []byte) (n int, err error) {
	params := content_library_image.NewCreateContentLibraryImageParamsWithTimeout(24 * time.Hour)
	clusterWhere, _ := json.Marshal(map[string]string{"id": w.clusterID})
	params.Clusters = string(clusterWhere)

	if w.pos == 0 {
		params.Name = pointy.String(w.name)
		params.Size = pointy.String(strconv.FormatInt(w.size, 10))
		params.File = runtime.NamedReader("chunk", io.NopCloser(bytes.NewReader(p)))
		createResp, err := w.client.CreateContentLibraryImage(params)
		if err != nil {
			return 0, fmt.Errorf("create content library image: %w", err)
		}
		if len(createResp.Payload) == 0 || createResp.Payload[0].ID == nil {
			return 0, fmt.Errorf("no upload task returned")
		}
		w.uploadTaskID = *createResp.Payload[0].ID
	} else {
		params.UploadTaskID = &w.uploadTaskID
		params.File = runtime.NamedReader("chunk", io.NopCloser(bytes.NewReader(p)))
		_, err := w.client.CreateContentLibraryImage(params)
		if err != nil {
			return 0, fmt.Errorf("upload chunk: %w", err)
		}
	}
	w.pos += len(p)
	return len(p), nil
}

// ---------- ImportContentLibraryImage ----------

func (c *defaultClient) ImportContentLibraryImage(ctx context.Context, path string, name string, clusterID string) (TaskRef, error) {
	file, err := os.Open(path)
	if err != nil {
		return TaskRef{}, fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return TaskRef{}, fmt.Errorf("stat file: %w", err)
	}

	writer := &cloudTowerUploadWriter{
		pos:       0,
		clusterID: clusterID,
		name:      name,
		size:      fileInfo.Size(),
		client:    c.api.ContentLibraryImage,
	}

	bufWriter := bufio.NewWriterSize(writer, uploadChunkSize)
	_, err = io.Copy(bufWriter, file)
	if err != nil {
		return TaskRef{}, fmt.Errorf("copy file: %w", err)
	}
	err = bufWriter.Flush()
	if err != nil {
		return TaskRef{}, fmt.Errorf("flush writer: %w", err)
	}

	if writer.uploadTaskID == "" {
		return TaskRef{}, fmt.Errorf("no upload task created")
	}
	return TaskRef{ID: writer.uploadTaskID, EntityKind: "upload_task"}, nil
}

// ---------- DistributeContentLibraryImage ----------

func (c *defaultClient) DistributeContentLibraryImage(ctx context.Context, id string, clusterIDs []string) (TaskRef, error) {
	params := content_library_image.NewDistributeContentLibraryImageClustersParams()
	params.SetContext(ctx)

	where := &models.ContentLibraryImageWhereInput{ID: pointy.String(id)}
	data := &models.ContentLibraryImageUpdationClusterParamsData{
		Clusters: &models.ClusterWhereInput{IDIn: clusterIDs},
	}

	params.SetRequestBody(&models.ContentLibraryImageUpdationClusterParams{
		Where: where,
		Data:  data,
	})

	resp, err := c.api.ContentLibraryImage.DistributeContentLibraryImageClusters(params)
	if err != nil {
		return TaskRef{}, fmt.Errorf("distribute content library image: %w", err)
	}
	return toWithTaskContentLibraryImageTaskRef(resp.Payload), nil
}

// toContentLibraryImage 把 SDK models.ContentLibraryImage 转成内部 ContentLibraryImage 模型。
func toContentLibraryImage(img *models.ContentLibraryImage) ContentLibraryImage {
	out := ContentLibraryImage{}
	if img.ID != nil {
		out.ID = *img.ID
	}
	if img.Name != nil {
		out.Name = *img.Name
	}
	if img.Description != nil {
		out.Description = *img.Description
	}
	if img.Path != nil {
		out.Path = *img.Path
	}
	if img.Size != nil {
		out.Size = *img.Size
	}
	if img.CreatedAt != nil {
		out.CreatedAt = *img.CreatedAt
	}
	for _, cluster := range img.Clusters {
		if cluster.ID != nil {
			out.ClusterIDs = append(out.ClusterIDs, *cluster.ID)
		}
		if cluster.Name != nil {
			out.ClusterNames = append(out.ClusterNames, *cluster.Name)
		}
	}
	return out
}

// toWithTaskContentLibraryImageTaskRef 从 WithTaskContentLibraryImage 提取 TaskRef。
func toWithTaskContentLibraryImageTaskRef(items []*models.WithTaskContentLibraryImage) TaskRef {
	if len(items) == 0 {
		return TaskRef{}
	}
	item := items[0]
	if item.TaskID == nil {
		return TaskRef{}
	}
	ref := TaskRef{ID: *item.TaskID}
	if item.Data != nil && item.Data.ID != nil {
		ref.EntityID = *item.Data.ID
	}
	ref.EntityKind = "content_library_image"
	return ref
}

package adapter

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/go-openapi/runtime"
	"github.com/openlyinc/pointy"
	"github.com/smartxworks/cloudtower-go-sdk/v2/client/cloud_tower_application"
	"github.com/smartxworks/cloudtower-go-sdk/v2/client/cloud_tower_application_package"
	"github.com/smartxworks/cloudtower-go-sdk/v2/models"
)

const (
	applicationUploadChunkSize = 8 * 1024 * 1024 // 8MB
)

// CloudTowerApplicationOps 定义 CloudTower 应用操作。
type CloudTowerApplicationOps interface {
	ListCloudTowerApplications(ctx context.Context, opts ListOpts) ([]CloudTowerApplication, error)
	ListCloudTowerApplicationPackages(ctx context.Context, opts ListOpts) ([]CloudTowerApplicationPackage, error)
	UploadCloudTowerApplicationPackage(ctx context.Context, path string, name string) (TaskRef, error)
	DeleteCloudTowerApplicationPackage(ctx context.Context, id string) error
	DeployCloudTowerApplication(ctx context.Context, name string, targetPackage string, vmSpec *models.ApplicationVMSpecDefinition) (TaskRef, error)
}

// ---------- ListCloudTowerApplications ----------

func (c *defaultClient) ListCloudTowerApplications(ctx context.Context, opts ListOpts) ([]CloudTowerApplication, error) {
	params := cloud_tower_application.NewGetCloudTowerApplicationsParams()
	params.SetContext(ctx)

	where := &models.CloudTowerApplicationWhereInput{}
	hasWhere := false
	if opts.NameContains != "" {
		where.NameContains = pointy.String(opts.NameContains)
		hasWhere = true
	}

	body := &models.GetCloudTowerApplicationsRequestBody{}
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

	resp, err := c.api.CloudTowerApplication.GetCloudTowerApplications(params)
	if err != nil {
		return nil, fmt.Errorf("list cloudtower applications: %w", err)
	}
	out := make([]CloudTowerApplication, 0, len(resp.Payload))
	for _, a := range resp.Payload {
		out = append(out, toCloudTowerApplication(a))
	}
	return out, nil
}

// ---------- ListCloudTowerApplicationPackages ----------

func (c *defaultClient) ListCloudTowerApplicationPackages(ctx context.Context, opts ListOpts) ([]CloudTowerApplicationPackage, error) {
	params := cloud_tower_application_package.NewGetCloudTowerApplicationPackagesParams()
	params.SetContext(ctx)

	where := &models.CloudTowerApplicationPackageWhereInput{}
	hasWhere := false
	if opts.NameContains != "" {
		where.NameContains = pointy.String(opts.NameContains)
		hasWhere = true
	}

	body := &models.GetCloudTowerApplicationPackagesRequestBody{}
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

	resp, err := c.api.CloudTowerApplicationPackage.GetCloudTowerApplicationPackages(params)
	if err != nil {
		return nil, fmt.Errorf("list cloudtower application packages: %w", err)
	}
	out := make([]CloudTowerApplicationPackage, 0, len(resp.Payload))
	for _, p := range resp.Payload {
		out = append(out, toCloudTowerApplicationPackage(p))
	}
	return out, nil
}

// cloudTowerApplicationUploadWriter 实现 io.Writer，用于分片上传应用包到 CloudTower。
type cloudTowerApplicationUploadWriter struct {
	pos          int
	uploadTaskID string
	name         string
	size         int64
	client       cloud_tower_application.ClientService
}

func (w *cloudTowerApplicationUploadWriter) Write(p []byte) (n int, err error) {
	params := cloud_tower_application.NewUploadCloudTowerApplicationPackageParamsWithTimeout(24 * time.Hour)

	if w.pos == 0 {
		params.Name = pointy.String(w.name)
		params.Size = pointy.String(strconv.FormatInt(w.size, 10))
		params.File = runtime.NamedReader("chunk", io.NopCloser(bytes.NewReader(p)))
		createResp, err := w.client.UploadCloudTowerApplicationPackage(params)
		if err != nil {
			return 0, fmt.Errorf("upload cloudtower application package: %w", err)
		}
		if len(createResp.Payload) == 0 || createResp.Payload[0].ID == nil {
			return 0, fmt.Errorf("no upload task returned")
		}
		w.uploadTaskID = *createResp.Payload[0].ID
	} else {
		params.UploadTaskID = &w.uploadTaskID
		params.File = runtime.NamedReader("chunk", io.NopCloser(bytes.NewReader(p)))
		_, err := w.client.UploadCloudTowerApplicationPackage(params)
		if err != nil {
			return 0, fmt.Errorf("upload chunk: %w", err)
		}
	}
	w.pos += len(p)
	return len(p), nil
}

// ---------- UploadCloudTowerApplicationPackage ----------

func (c *defaultClient) UploadCloudTowerApplicationPackage(ctx context.Context, path string, name string) (TaskRef, error) {
	file, err := os.Open(path)
	if err != nil {
		return TaskRef{}, fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return TaskRef{}, fmt.Errorf("stat file: %w", err)
	}

	writer := &cloudTowerApplicationUploadWriter{
		pos:    0,
		name:   name,
		size:   fileInfo.Size(),
		client: c.api.CloudTowerApplication,
	}

	bufWriter := bufio.NewWriterSize(writer, applicationUploadChunkSize)
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

// ---------- DeleteCloudTowerApplicationPackage ----------

func (c *defaultClient) DeleteCloudTowerApplicationPackage(ctx context.Context, id string) error {
	params := cloud_tower_application.NewDeleteCloudTowerApplicationPackageParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.DeleteCloudTowerApplicationPackageParams{
		Where: &models.CloudTowerApplicationPackageWhereInput{ID: pointy.String(id)},
	})

	_, err := c.api.CloudTowerApplication.DeleteCloudTowerApplicationPackage(params)
	if err != nil {
		return fmt.Errorf("delete cloudtower application package: %w", err)
	}
	return nil
}

// ---------- DeployCloudTowerApplication ----------

func (c *defaultClient) DeployCloudTowerApplication(ctx context.Context, name string, targetPackage string, vmSpec *models.ApplicationVMSpecDefinition) (TaskRef, error) {
	params := cloud_tower_application.NewDeployCloudTowerApplicationParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.DeployCloudTowerApplicationParams{
		Name:         pointy.String(name),
		TargetPackage: pointy.String(targetPackage),
		VMSpec:        vmSpec,
	})

	resp, err := c.api.CloudTowerApplication.DeployCloudTowerApplication(params)
	if err != nil {
		return TaskRef{}, fmt.Errorf("deploy cloudtower application: %w", err)
	}
	if resp.Payload == nil || resp.Payload.ID == nil {
		return TaskRef{}, nil
	}
	return TaskRef{ID: *resp.Payload.ID, EntityKind: "cloudtower_application"}, nil
}

// toCloudTowerApplication 把 SDK models.CloudTowerApplication 转成内部 CloudTowerApplication 模型。
func toCloudTowerApplication(a *models.CloudTowerApplication) CloudTowerApplication {
	out := CloudTowerApplication{}
	if a.ID != nil {
		out.ID = *a.ID
	}
	if a.Name != nil {
		out.Name = *a.Name
	}
	if a.State != nil {
		out.State = string(*a.State)
	}
	if a.TargetPackage != nil {
		out.TargetPackage = *a.TargetPackage
	}
	if a.ResourceVersion != nil {
		out.ResourceVersion = *a.ResourceVersion
	}
	return out
}

// toCloudTowerApplicationPackage 把 SDK models.CloudTowerApplicationPackage 转成内部 CloudTowerApplicationPackage 模型。
func toCloudTowerApplicationPackage(p *models.CloudTowerApplicationPackage) CloudTowerApplicationPackage {
	out := CloudTowerApplicationPackage{}
	if p.ID != nil {
		out.ID = *p.ID
	}
	if p.Name != nil {
		out.Name = *p.Name
	}
	if p.Version != nil {
		out.Version = *p.Version
	}
	if p.Architecture != nil {
		out.Architecture = string(*p.Architecture)
	}
	if p.ScosVersion != nil {
		out.ScosVersion = *p.ScosVersion
	}
	return out
}

// openFile is a helper to open files, allowing for testing
var openFile = func(name string) (io.ReadCloser, error) {
	return os.Open(name)
}

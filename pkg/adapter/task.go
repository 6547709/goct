package adapter

import (
	"context"
	"fmt"

	"github.com/openlyinc/pointy"
	sdktask "github.com/smartxworks/cloudtower-go-sdk/v2/client/task"
	"github.com/smartxworks/cloudtower-go-sdk/v2/models"
)

// TaskOps 定义任务操作（list/get 用于命令层，GetTaskProgress 已在 task_progress.go）。
type TaskOps interface {
	ListTasks(ctx context.Context, opts ListOpts) ([]Task, error)
	GetTask(ctx context.Context, id string) (*Task, error)
}

func (c *defaultClient) ListTasks(ctx context.Context, opts ListOpts) ([]Task, error) {
	params := sdktask.NewGetTasksParams()
	params.SetContext(ctx)
	body := &models.GetTasksRequestBody{}
	if opts.Limit > 0 {
		body.First = pointy.Int32(opts.Limit)
	}
	params.SetRequestBody(body)
	resp, err := c.api.Task.GetTasks(params)
	if err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}
	out := make([]Task, 0, len(resp.Payload))
	for _, t := range resp.Payload {
		out = append(out, toTask(t))
	}
	return out, nil
}

func (c *defaultClient) GetTask(ctx context.Context, id string) (*Task, error) {
	params := sdktask.NewGetTasksParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.GetTasksRequestBody{
		Where: &models.TaskWhereInput{ID: pointy.String(id)},
	})
	resp, err := c.api.Task.GetTasks(params)
	if err != nil {
		return nil, fmt.Errorf("get task %s: %w", id, err)
	}
	if len(resp.Payload) == 0 {
		return nil, fmt.Errorf("get task %s: %w", id, ErrNotFound)
	}
	t := toTask(resp.Payload[0])
	return &t, nil
}

func toTask(t *models.Task) Task {
	out := Task{}
	if t.ID != nil { out.ID = *t.ID }
	if t.Description != nil { out.Description = *t.Description }
	if t.Status != nil { out.Status = string(*t.Status) }
	if t.Progress != nil { out.Progress = int(*t.Progress) }
	if t.ErrorMessage != nil { out.ErrorMessage = *t.ErrorMessage }
	if t.LocalCreatedAt != nil { out.CreatedAt = *t.LocalCreatedAt }
	if t.StartedAt != nil { out.StartedAt = *t.StartedAt }
	if t.FinishedAt != nil { out.FinishedAt = *t.FinishedAt }
	return out
}

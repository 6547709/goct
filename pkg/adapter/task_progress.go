package adapter

import (
	"context"
	"fmt"

	"github.com/openlyinc/pointy"
	sdktask "github.com/smartxworks/cloudtower-go-sdk/v2/client/task"
	"github.com/smartxworks/cloudtower-go-sdk/v2/models"
)

// GetTaskProgress 实现 task.Ops 接口，允许 watcher 轮询。
func (c *defaultClient) GetTaskProgress(ctx context.Context, id string) (percent int, status string, err error) {
	params := sdktask.NewGetTasksParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.GetTasksRequestBody{
		Where: &models.TaskWhereInput{ID: pointy.String(id)},
	})
	resp, err := c.api.Task.GetTasks(params)
	if err != nil {
		return 0, "", fmt.Errorf("get task %s: %w", id, err)
	}
	if len(resp.Payload) == 0 {
		return 0, "", fmt.Errorf("get task %s: %w", id, ErrNotFound)
	}
	t := resp.Payload[0]
	pct := 0
	if t.Progress != nil {
		pct = int(*t.Progress)
	}
	st := ""
	if t.Status != nil {
		st = string(*t.Status)
	}
	return pct, st, nil
}

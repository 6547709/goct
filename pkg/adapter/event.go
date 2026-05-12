// Package adapter — event.go 实现 govc 风格 `events` 命令的 SDK 包装。
//
// CloudTower 没有"事件流（events）"概念，最接近的是 user_audit_log（用户审计日志）：
// 谁在什么时间做了什么 action，针对哪个 resource。这刚好对应 govc events 的展示信息。
//
// 我们把 user_audit_log 适配成 Event 类型，让 cmd/events 命令直接消费。
package adapter

import (
	"context"
	"fmt"

	"github.com/openlyinc/pointy"
	userAudit "github.com/smartxworks/cloudtower-go-sdk/v2/client/user_audit_log"
	"github.com/smartxworks/cloudtower-go-sdk/v2/models"
)

// Event 是 CLI 内部用的事件视图（基于 user_audit_log）。
//
// 字段语义与 vSphere event 大致对齐：
//   - CreatedAt：事件发生时间（govc events 展示的 "Date"）
//   - Username：触发者（"User"）
//   - Action：动作名（"Type"），如 "create_vm" / "power_on_vm"
//   - ResourceType / ResourceID：受影响实体（"Target"）
//   - Message：人类可读描述
//   - Status：成功 / 失败 / 进行中
type Event struct {
	ID           string
	CreatedAt    string
	StartedAt    string
	FinishedAt   string
	Username     string
	IPAddress    string
	Action       string
	ResourceType string
	ResourceID   string
	Message      string
	Status       string
	ClusterID    string
	ClusterName  string
}

// EventListOpts 是 ListEvents 的过滤条件。
type EventListOpts struct {
	ResourceID   string // 精确匹配某个资源（如 vmID）的事件
	ResourceType string // 限定资源类型（VM / HOST / CLUSTER ...）
	Username     string // 限定触发者
	ActionLike   string // action 字段的 contains 匹配
	Limit        int32  // 0 = SDK 默认（通常 100）
	Skip         int32
}

// EventOps 暴露 events 相关的最小接口。
type EventOps interface {
	ListEvents(ctx context.Context, opts EventListOpts) ([]Event, error)
}

// ListEvents 拉取 user_audit_log 并转换成 Event。
//
// 默认按 createdAt DESC 排序（最新事件在前），与 govc events 默认行为一致。
func (c *defaultClient) ListEvents(ctx context.Context, opts EventListOpts) ([]Event, error) {
	params := userAudit.NewGetUserAuditLogsParams()
	params.SetContext(ctx)

	body := &models.GetUserAuditLogsRequestBody{
		OrderBy: models.NewUserAuditLogOrderByInput(models.UserAuditLogOrderByInputCreatedAtDESC),
	}
	if opts.Limit > 0 {
		body.First = pointy.Int32(opts.Limit)
	}
	if opts.Skip > 0 {
		body.Skip = pointy.Int32(opts.Skip)
	}

	where := &models.UserAuditLogWhereInput{}
	hasWhere := false
	if opts.ResourceID != "" {
		where.ResourceID = pointy.String(opts.ResourceID)
		hasWhere = true
	}
	if opts.ResourceType != "" {
		where.ResourceType = pointy.String(opts.ResourceType)
		hasWhere = true
	}
	if opts.Username != "" {
		where.Username = pointy.String(opts.Username)
		hasWhere = true
	}
	if opts.ActionLike != "" {
		where.ActionContains = pointy.String(opts.ActionLike)
		hasWhere = true
	}
	if hasWhere {
		body.Where = where
	}
	params.SetRequestBody(body)

	resp, err := c.api.UserAuditLog.GetUserAuditLogs(params)
	if err != nil {
		return nil, fmt.Errorf("list events: %w", err)
	}
	out := make([]Event, 0, len(resp.Payload))
	for _, e := range resp.Payload {
		out = append(out, toEvent(e))
	}
	return out, nil
}

func toEvent(e *models.UserAuditLog) Event {
	out := Event{}
	if e == nil {
		return out
	}
	if e.ID != nil {
		out.ID = *e.ID
	}
	if e.CreatedAt != nil {
		out.CreatedAt = *e.CreatedAt
	}
	if e.StartedAt != nil {
		out.StartedAt = *e.StartedAt
	}
	if e.FinishedAt != nil {
		out.FinishedAt = *e.FinishedAt
	}
	if e.Username != nil {
		out.Username = *e.Username
	}
	if e.IPAddress != nil {
		out.IPAddress = *e.IPAddress
	}
	if e.Action != nil {
		out.Action = *e.Action
	}
	if e.ResourceType != nil {
		out.ResourceType = *e.ResourceType
	}
	if e.ResourceID != nil {
		out.ResourceID = *e.ResourceID
	}
	if e.Message != nil {
		out.Message = *e.Message
	}
	if e.Status != nil {
		out.Status = string(*e.Status)
	}
	if e.Cluster != nil {
		if e.Cluster.ID != nil {
			out.ClusterID = *e.Cluster.ID
		}
		if e.Cluster.Name != nil {
			out.ClusterName = *e.Cluster.Name
		}
	}
	return out
}

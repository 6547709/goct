// Package adapter 是 goct 的 SDK 防腐层（Anti-Corruption Layer）。
//
// 这是整个项目中**唯一**允许 import `cloudtower-go-sdk/v2` 的位置。
// 所有 service / cmd 层只看 adapter 暴露的内部模型与接口，
// 与 SDK 的 schema 演进完全解耦。
//
// 设计要点：
//   - 内部模型只保留 CLI 实际需要的字段（不是 SDK 模型的镜像）
//   - 资源 sub-interface（VMOps / HostOps / ...）在 client.go 内嵌进 Client，
//     便于 service 层只依赖单一接口便于测试 mock
//   - 写操作返回 TaskRef；TaskRef.ID == "" 表示同步完成（不需要 watch）
package adapter

import "time"

// TowerInfo 描述 CloudTower 实例的版本与构建信息。
// CloudTower API 的 GetAPIVersion 仅返回裸版本字符串，Build 字段当前为空。
type TowerInfo struct {
	Version string
	Build   string
}

// VM 是 CLI 内部用的虚拟机视图。
type VM struct {
	ID          string
	Name        string
	Status      string // 例：RUNNING / STOPPED / SUSPENDED
	ClusterID   string
	VCPU        int32
	MemoryBytes uint64
	IPs         []string
	Description string
}

// TaskRef 是写操作返回的 task 引用。
//
//	ID == ""           表示该操作同步完成，不需要 watch
//	ID != ""           需要通过 task watcher 等待结束
//	EntityID/Kind 可选，用于错误/日志展示
type TaskRef struct {
	ID         string
	EntityID   string
	EntityKind string
}

// IsSync 报告该 TaskRef 是否表示同步操作。
func (r TaskRef) IsSync() bool { return r.ID == "" }

// ListOpts 是 list 类操作的统一过滤条件。
// 各 sub-interface 根据自身能力可忽略不支持的字段。
type ListOpts struct {
	NameContains string
	ClusterID    string
	Limit        int32
	Skip         int32
}

// PowerAction 抽象 VM 电源操作动作。
// 强制语义由各方法的 force 参数承载，避免动作枚举膨胀。
type PowerAction string

const (
	PowerOn      PowerAction = "ON"
	PowerOff     PowerAction = "OFF"
	PowerReset   PowerAction = "RESET"
	PowerSuspend PowerAction = "SUSPEND"
	PowerResume  PowerAction = "RESUME"
)

// SessionToken 是 adapter 暴露给 client 层的鉴权凭据。
// 与 pkg/session.Token 同构，避免 adapter 反向依赖 session 包。
type SessionToken struct {
	Value    string
	ExpireAt time.Time
}

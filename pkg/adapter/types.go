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
	ID              string
	Name            string
	Status          string // 例：RUNNING / STOPPED / SUSPENDED
	ClusterID       string
	ClusterName     string
	HostID          string
	HostName        string
	HostIP          string
	VCPU            int32
	MemoryBytes     uint64
	Firmware        string // BIOS / UEFI
	Ha              bool
	GuestOS         string // LINUX / WINDOWS / UNKNOWN (guest_os_type enum)
	OS              string // Actual OS name, e.g. "Ubuntu 22.04" (from os field)
	VMToolsStatus   string // RUNNING / NOT_RUNNING / NOT_INSTALLED / RESTRICTION
	VMToolsVersion  string
	CPUModel        string
	IPs             []string
	DNSServers      string
	Hostname        string
	DiskCount       int
	NicCount        int
	ProvisionedBytes uint64
	UsedBytes       uint64
	InRecycleBin    bool
	Protected       bool
	CreatedAt       string
	Description     string

	// Detail fields (available via --detail)
	BiosUUID         string
	CPUUsage         float64
	MemoryUsage      float64
	GuestSizeUsage   float64
	GuestUsedSize    int64
	LogicalSizeBytes int64
	UsbDevices       []UsbDevice
	GpuDevices       []GpuDevice
	VideoType        string
	NestedVirt       bool
	HaPriority       string
	Labels           []string
	CloudInit        bool
}

// GpuDevice represents a GPU device attached to a VM.
type GpuDevice struct {
	ID   string
	Name string
}

// UsbDevice represents a USB device attached to a VM.
type UsbDevice struct {
	ID   string
	Name string
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

// Task 是 CLI 内部用的任务视图。
type Task struct {
	ID           string
	Description  string
	Status       string
	Progress     int
	ErrorMessage string
	CreatedAt    string
	FinishedAt   string
	StartedAt    string
}

// Alert 是 CLI 内部用的告警视图。
type Alert struct {
	ID       string
	Message  string
	Severity string
	Cause    string
}

// User 是 CLI 内部用的用户视图。
type User struct {
	ID       string
	Name     string
	Username string
	Source   string
	Role     string
	Email    string
}

// UserCreateSpec 是 user.create 需要的参数。
type UserCreateSpec struct {
	Name     string
	Username string
	Password string
	RoleID   string
	Email    string
}

// Network 是 CLI 内部用的虚拟交换机（VDS）视图。
type Network struct {
	ID        string
	Name      string
	Type      string
	ClusterID string
}

// VLAN 是 CLI 内部用的 VLAN 视图。
type VLAN struct {
	ID      string
	Name    string
	VlanTag int32
	Type    string
	VdsID   string
}

// VLANCreateSpec 是 vlan.create 需要的参数。
type VLANCreateSpec struct {
	Name  string
	VdsID string
}

// Cluster 是 CLI 内部用的集群视图。
type Cluster struct {
	ID               string
	Name             string
	TotalMemoryBytes uint64
	TotalDataCapacity uint64
	UsedDataSpace    uint64
	TotalCPUCores    int32
	RunningVMs       int32
}

// Datastore 是 CLI 内部用的数据存储视图。
type Datastore struct {
	ID        string
	Name      string
	Type      string
	Internal  bool
	ClusterID string
}

// DiskPool 是超融合存储池视图（每个Host上的本地存储聚合）。
type DiskPool struct {
	ID              string
	HostID          string
	HostName        string
	ClusterID       string
	Status          string
	UseState        string
	TotalDataBytes  uint64
	UsedDataBytes   uint64
	TotalCacheBytes uint64
	UsedCacheBytes  uint64
	HddCount       int32
	NvmeCount      int32
	SataCount      int32
}

// Disk 是 CLI 内部用的磁盘视图。
type Disk struct {
	ID        string
	Name      string
	Type      string
	SizeBytes uint64
	Path      string
	HostName  string
}

// Host 是 CLI 内部用的主机视图。
type Host struct {
	ID              string
	Name            string
	Status          string
	ManagementIP    string
	DataIP          string
	CPUModel        string
	TotalMemoryBytes uint64
	RunningVMs      int32
	ClusterID       string
}

// Snapshot 是 CLI 内部用的快照视图。
type Snapshot struct {
	ID          string
	Name        string
	VMID        string
	Description string
	CreatedAt   string // local_created_at 原样传出
}

// VMCreateSpec 是 vm.create 命令需要的参数集合。
type VMCreateSpec struct {
	Name        string
	ClusterID   string
	VCPU        int32
	MemoryBytes int64  // bytes
	Firmware    string // BIOS / UEFI, default BIOS
	Description string
}

// VMCreateFromTemplateSpec 是 vm.create --from-template 命令的参数集合。
type VMCreateFromTemplateSpec struct {
	TemplateID  string
	Name        string
	ClusterID   string
	HostID      string
	VCPU        int32
	MemoryBytes int64
	Firmware    string
	Description string
	IsFullCopy  bool
	NIC         NicConfig // 网卡配置
}

// NicConfig 网卡配置
type NicConfig struct {
	Type   string // VLAN / VPC
	Model  string // E1000 / SRIOV / VIRTIO
	VlanID string
}

// VMCloneSpec 是 vm.clone 命令需要的参数集合。
type VMCloneSpec struct {
	Name            string
	TargetClusterID string // 可选；空则同集群
}

// VMExportSpec 是 vm.export 命令的参数。
type VMExportSpec struct {
	FileType string // OVF (default)
	KeepMAC  bool
}

// VMUpdateSpec 是 vm.update 命令的参数。
type VMUpdateSpec struct {
	Name        string
	Description string
}

// DiskAddSpec 是 vm disk.add 命令的参数。
type DiskAddSpec struct {
	Name      string
	SizeBytes int64
	Bus       string // SCSI / SATA / NVMe / IDE / VIRTIO
	Index     int32
	Boot      int32
	IOPSMax   int64
}

// CdRomAddSpec 是 vm cdrom.add 命令的参数。
type CdRomAddSpec struct {
	Boot int32
	Path string // ISO path or content library image ID
}

// NicAddSpec 是 vm nic.add 命令的参数。
type NicAddSpec struct {
	Type    string // NIC_TYPE_NORMAL / NIC_TYPE_DIRECT
	Model   string // RTL8139 / E1000 / VIRTIO
	VlanID  string
}

// VMNic 是 VM 网卡视图（用于列表和详情）。
type VMNic struct {
	ID              string
	VMID            string
	LocalID         string // NIC index (string in SDK)
	MacAddress      string
	Model           string
	Type            string // VLAN / VPC
	VlanID          string
	VlanName        string
	Gateway         string
	SubnetMask      string
	IPAddress       string
	Enabled         bool
	IngressRateLimit  *int64
	EgressRateLimit   *int64
}

// VMNicUpdateSpec 是 vm nic.update 命令的参数。
type VMNicUpdateSpec struct {
	NicIndex      int32
	ConnectVlanID string
	Enabled       *bool
	Gateway       string
	IPAddress     string
	MacAddress    string
	Model         string
	SubnetMask    string
}

// VMDisk 是 VM 磁盘视图（用于列表，包含 CD-ROM）。
type VMDisk struct {
	ID              string
	VMID            string
	Boot            int32
	Bus             string // SCSI / SATA / NVMe / IDE / VIRTIO
	Key             int32
	MaxBandwidth    *int64
	MaxIops         *int64
	Type            string // DISK / CD_ROM
	VolumeID        string
	VolumeName      string
	VolumeSize      int64
	ElfImageID      string
	ElfImageName    string
}

// DiskUpdateSpec 是 vm disk.update 命令的参数。
type DiskUpdateSpec struct {
	MaxBandwidth *int64
	MaxIops      *int64
}

// CdRomToggleSpec 是 vm cdrom.toggle 命令的参数。
type CdRomToggleSpec struct {
	Disabled bool
}

// ResetPasswordSpec 是 vm reset-password 命令的参数。
type ResetPasswordSpec struct {
	Username string
	Password string
}

// RebuildVMSpec 是 vm rebuild 命令的参数。
type RebuildVMSpec struct {
	SnapshotID string
	Name      string
	ClusterID string
	HostID    string
}

// ImportVMSpec 是 vm import 命令的参数。
type ImportVMSpec struct {
	ClusterID  string
	Name      string
	CPUCores  int32
	CPUSockets int32
	Memory    int64
	Vcpu      int32
	Ha        bool
	HostID    string
}

// VNCInfo 是 VM VNC 连接信息。
type VNCInfo struct {
	ClusterIP string
	Redirect  string
	Terminal string
	Direct   string
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

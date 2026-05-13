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
//
// Name 与 NameContains 互斥：Name 走精确匹配（比 NameContains 更高效，避免大列表），
// 同时存在时优先 Name。
type ListOpts struct {
	Name         string // 精确匹配，服务端走 = 过滤
	NameContains string // 模糊匹配，服务端走 contains 过滤
	ClusterID    string
	Limit        int32
	Skip         int32
	InRecycleBin *bool // nil=不过滤, true=仅回收站, false=仅正常VM
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
//
// Disks/Nics 为空时，adapter 不会下发任何默认磁盘/网卡（v0.2.1 之前会强制塞 10 GB SCSI 磁盘 +
// 无 VLAN 的 VIRTIO 网卡，导致 CloudTower 拒绝或行为不可控）。
// 由 cmd / service 层根据用户传入的 --disk / --nic 显式构造。
//
// HA 为 nil 时不下发该字段（用 CloudTower 的默认值）；显式 *bool 区分"未指定"与"显式关闭"。
type VMCreateSpec struct {
	Name        string
	ClusterID   string
	VCPU        int32
	MemoryBytes int64  // bytes
	Firmware    string // BIOS / UEFI, default BIOS
	Description string
	HA          *bool
	Disks       []DiskAddSpec
	Nics        []NicAddSpec
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
	NIC         NicConfig        // 网卡配置
	CloudInit   *CloudInitSpec   // cloud-init 配置
}

// NicConfig 网卡配置
type NicConfig struct {
	Type   string // VLAN / VPC
	Model  string // E1000 / SRIOV / VIRTIO
	VlanID string
}

// CloudInitSpec describes cloud-init configuration for VM creation from template.
type CloudInitSpec struct {
	Hostname            string            // VM hostname
	DefaultUserPassword string            // Default user password
	PublicKeys          []string          // SSH public keys (authorized_keys)
	DNSServers          []string          // Global DNS nameservers
	UserData            string            // Custom cloud-init user_data (YAML script or #cloud-config)
	Networks            []NicStaticConfig // Per-NIC network configuration
}

// NicStaticConfig describes static IP configuration for one NIC.
type NicStaticConfig struct {
	Index   int32         // NIC index (0-based), required
	IP      string        // Static IP address (e.g. "192.168.1.100")
	Netmask string        // Netmask in dotted notation (e.g. "255.255.255.0")
	Gateway string        // Default gateway IP
	Type    string        // "IPV4" (static) or "IPV4_DHCP" (DHCP)
	Routes  []StaticRoute // Custom static routes (excluding default route)
}

// StaticRoute describes a static route entry.
type StaticRoute struct {
	Network string // Destination network (e.g. "10.0.0.0/8")
	Netmask string // Netmask for the route (e.g. "255.0.0.0")
	Gateway string // Next hop gateway
}

// VMCloneSpec 是 vm.clone 命令需要的参数集合。
type VMCloneSpec struct {
	Name            string
	TargetClusterID string // 可选；空则同集群
	Linked          bool  // true=linked clone, false=full clone
}

// VMExportSpec 是 vm.export 命令的参数。
type VMExportSpec struct {
	FileType string // OVF (default)
	KeepMAC  bool
}

// VMUpdateSpec 是 vm.update 命令的参数。
//
// 字段语义（v0.2.1 修复）：
//   - nil  → 不下发该字段（保持原值）
//   - 非 nil（含空字符串）→ 显式下发，允许把 description 清空为 ""
//
// 之前用 string + `if != "" { set }`，导致用户没办法把 description 清空。
type VMUpdateSpec struct {
	Name        *string
	Description *string
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
//
// VMID（v0.2.1 新增）：CloudTower update-vm-nic API 要求 Where 指向 VM。
type VMNicUpdateSpec struct {
	VMID          string
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
//
// 注意：CloudTower 的 update-vm-disk API（VMUpdateDiskParamsData）只支持下列字段：
//   - bus
//   - vm_volume_id（替换底层 volume）
//   - elf_image_id / content_library_image_id（CD-ROM 换 ISO 用）
//
// 不支持 size / iops / bandwidth 在线修改：扩容请用 vm.disk.expand，
// QoS 请用 vm.disk QoS 专用 API（SDK v2.22.1 暂未导出）。
//
// VMID 必填：CloudTower 的 update-vm-disk 用 VM 维度的 Where 过滤。
type DiskUpdateSpec struct {
	VMID                  string
	Bus                   string // SCSI / IDE / VIRTIO；空表示不改
	VMVolumeID            string // 替换底层 volume；空表示不改
	ElfImageID            string // 替换 ISO；空表示不改
	ContentLibraryImageID string // 替换内容库镜像；空表示不改

	// 保留字段：CLI 历史 flag 兼容，当前 SDK 不支持，传入会被忽略并打 warning。
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

// Label 是 CLI 内部用的标签视图。
type Label struct {
	ID        string
	Key       string
	Value     string
	CreatedAt string
}

// LabelCreateSpec 是 label.create 命令的参数。
type LabelCreateSpec struct {
	Key   string
	Value string
}

// LabelUpdateSpec 是 label.update 命令的参数。
type LabelUpdateSpec struct {
	Key   string
	Value string
}

// LabelAttachSpec 是 label attach 命令的参数。
type LabelAttachSpec struct {
	ResourceKind string // vm, host, cluster, etc.
	ResourceID   string
}

// LabelDetachSpec 是 label detach 命令的参数。
type LabelDetachSpec struct {
	ResourceKind string
	ResourceID   string
}

// VMFolder 是 CLI 内部用的 VM 文件夹视图。
type VMFolder struct {
	ID        string
	Name      string
	ClusterID string
}

// VMFolderCreateSpec 是 vm.folder.create 命令的参数。
type VMFolderCreateSpec struct {
	Name      string
	ClusterID string
}

// VMFolderUpdateSpec 是 vm.folder.update 命令的参数。
type VMFolderUpdateSpec struct {
	Name string
}

// VMPlacementGroup 是 CLI 内部用的 VM 放置组视图。
type VMPlacementGroup struct {
	ID        string
	Name      string
	ClusterID string
}

// VMPlacementGroupCreateSpec 是 vm.placement-group.create 命令的参数。
type VMPlacementGroupCreateSpec struct {
	Name      string
	ClusterID string
}

// SnapshotPlan 是 CLI 内部用的快照计划视图。
type SnapshotPlan struct {
	ID               string
	Name             string
	ClusterID        string
	Status           string
	PlanType         string
	Retention        int32
	StartTime        string
	EndTime          string
	ExecHM           string
	ExecuteIntervals []int32
}

// SnapshotPlanCreateSpec 是 snapshot.plan.create 命令的参数。
type SnapshotPlanCreateSpec struct {
	Name      string
	ClusterID string
	PlanType  string
	Retention int32
	StartTime string
	EndTime   string
	VMIDs     []string
}

// ElfStoragePolicy 是 CLI 内部用的存储策略视图。
type ElfStoragePolicy struct {
	ID           string
	Name         string
	Description string
	ClusterID    string
	ClusterName  string
	LocalID      string
	ReplicaNum   int32
	StripeNum    int32
	StripeSize   int64
	ThinProvision bool
}

// GlobalSettings 是 CLI 内部用的全局设置视图。
type GlobalSettings struct {
	ID              string
	SessionMaxAge   int32
	VMRecycleBin    VMRecycleBinSetting
}

// VMRecycleBinSetting 是回收站设置视图。
type VMRecycleBinSetting struct {
	RetainPeriod int32
	Enabled      bool
}

// UsbDevice 是 CLI 内部用的 USB 设备视图。
type UsbDevice struct {
	ID            string
	Name          string
	Description   string
	LocalID       string
	Manufacturer  string
	Status        string
	UsbType       string
	Size          int64
	Binded        bool
	HostID        string
	HostName      string
	VMID          string
	VMName        string
	LocalCreatedAt string
}

// Application 是 CLI 内部用的应用视图。
type Application struct {
	ID         string
	LocalID    string
	ImageName  string
	Memory     int64
	State      string
	StorageIP  string
	ClusterID  string
	ClusterName string
	ErrorMessage string
}

// Deploy 是 CLI 内部用的部署视图。
type Deploy struct {
	ID        string
	LocalID   string
	VMID      string
	VMName    string
	State     string
	StartedAt string
}

// License 是 CLI 内部用的许可证视图。
type License struct {
	ID                   string
	ExpireDate           string
	LicenseSerial        string
	MaintenanceEndDate   string
	MaintenanceStartDate string
	MaxChunkNum          int32
	MaxClusterNum        int32
	SignDate             string
	SoftwareEdition      string
	Type                string
}

// ClusterSettings 是 CLI 内部用的集群设置视图。
type ClusterSettings struct {
	ID        string
	ClusterID string
}

// NtpSettings 是 CLI 内部用的 NTP 设置视图。
type NtpSettings struct {
	URLs []string
}

// AlertRule 是 CLI 内部用的告警规则视图。
type AlertRule struct {
	ID          string
	Name        string
	Enabled     bool
	Expression  string
	Duration    int32
	Severity    string
	TargetKind  string
	TargetID    string
}

// ContentLibraryImage 是 CLI 内部用的内容库镜像视图。
type ContentLibraryImage struct {
	ID          string
	Name        string
	Description string
	Path        string
	Size        int64
	CreatedAt   string
	ClusterIDs  []string
	ClusterNames []string
}

// CloudTowerApplication 是 CLI 内部用的 CloudTower 应用视图。
type CloudTowerApplication struct {
	ID           string
	Name         string
	State        string
	TargetPackage string
	ResourceVersion int32
	ClusterID    string
	ClusterName  string
}

// CloudTowerApplicationPackage 是 CLI 内部用的 CloudTower 应用包视图。
type CloudTowerApplicationPackage struct {
	ID           string
	Name         string
	Version      string
	Architecture string
	ScosVersion  string
}

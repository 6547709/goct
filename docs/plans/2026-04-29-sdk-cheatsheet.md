# CloudTower Go SDK v2.22.1 — Adapter 速查表

> 本文档供 `goct` 的 adapter 层（唯一允许 import SDK 的层）查阅。
> 所有签名/字段均从 `~/go/pkg/mod/github.com/smartxworks/cloudtower-go-sdk/v2@v2.22.1/` 实际源码读出。
> 探查日期：2026-04-29

---

## 0. 通用模式

### 0.1 客户端构造（推荐入口）

```go
import (
    apiclient "github.com/smartxworks/cloudtower-go-sdk/v2/client"
    "github.com/smartxworks/cloudtower-go-sdk/v2/models"
)

cli, err := apiclient.NewWithUserConfig(
    apiclient.ClientConfig{
        Host:     "tower.example.com:443", // 不带 scheme
        BasePath: "/v2/api",
        Schemes:  []string{"https"},
    },
    apiclient.UserConfig{
        Name:     "admin",
        Password: "xxx",
        Source:   models.UserSourceLOCAL,
    },
)
```

底层等价：`apiclient.New(transport, strfmt.Default)` + `client.User.Login(...)` + `transport.DefaultAuthentication = httptransport.APIKeyAuth("Authorization", "header", *resp.Payload.Data.Token)`。

### 0.2 自签名证书（goct --insecure）

```go
import httptransport "github.com/go-openapi/runtime/client"

httpClient := &http.Client{Transport: &http.Transport{
    TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}}
tr := httptransport.NewWithClient(host, "/v2/api", []string{"https"}, httpClient)
cli := apiclient.New(tr, strfmt.Default)
// 然后手动 Login + 注入 token
```

### 0.3 Cloudtower 顶层字段（仅 Tier-1）

| 字段 | 子包 |
| --- | --- |
| `cli.APIInfo` | `client/api_info` |
| `cli.User` | `client/user` |
| `cli.UserRoleNext` | `client/user_role_next` |
| `cli.VM` | `client/vm` |
| `cli.VMSnapshot` | `client/vm_snapshot` |
| `cli.Host` | `client/host` |
| `cli.Cluster` | `client/cluster` |
| `cli.ElfDataStore` | `client/elf_data_store` |
| `cli.Disk` | `client/disk` |
| `cli.Vds` | `client/vds` |
| `cli.Nic` | `client/nic` |
| `cli.Vlan` | `client/vlan` |
| `cli.Task` | `client/task` |
| `cli.Alert` | `client/alert` |
| `cli.AlertRule` | `client/alert_rule` |

### 0.4 调用模板（GET 类）

```go
p := vm.NewGetVmsParams()
p.RequestBody = &models.GetVmsRequestBody{
    Where: &models.VMWhereInput{NameContains: pointy.String("foo")},
    First: pointy.Int32(20),
}
res, err := cli.VM.GetVms(p)
list := res.GetPayload() // []*models.VM
```

### 0.5 调用模板（写操作 + 等待 task）

```go
p := vm.NewShutDownVMParams()
p.RequestBody = &models.VMOperateParams{
    Where: &models.VMWhereInput{ID: pointy.String("vm-xxx")},
}
res, err := cli.VM.ShutDownVM(p) // []*models.WithTaskVM
ids := make([]string, 0, len(res.Payload))
for _, item := range res.Payload {
    if item.TaskID != nil {
        ids = append(ids, *item.TaskID)
    }
}
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()
return utils.WaitTasks(ctx, cli, ids, 2*time.Second)
```

### 0.6 关键枚举

**`models.UserSource`** （`models/user_source.go`，类型 `string`，自带 `Pointer()`）：
- `UserSourceAUTHN` = "AUTHN"
- `UserSourceLDAP`  = "LDAP"
- `UserSourceLOCAL` = "LOCAL"
- `UserSourceSSO`   = "SSO"

**`models.TaskStatus`** （`models/task_status.go`）：
- `TaskStatusEXECUTING` = "EXECUTING"
- `TaskStatusFAILED`    = "FAILED"
- `TaskStatusPAUSED`    = "PAUSED"
- `TaskStatusPENDING`   = "PENDING"
- `TaskStatusSUCCESSED` = "SUCCESSED"  ← **注意：是 SUCCESSED 不是 SUCCEEDED**

**`models.VMStatus`** （`models/vm_status.go`）：
- `VMStatusDELETED`   = "DELETED"
- `VMStatusRUNNING`   = "RUNNING"
- `VMStatusSTOPPED`   = "STOPPED"
- `VMStatusSUSPENDED` = "SUSPENDED"
- `VMStatusUNKNOWN`   = "UNKNOWN"

**`models.HostStatus`** （`models/host_status.go`）：
- `HostStatusCONNECTEDERROR`   = "CONNECTED_ERROR"
- `HostStatusCONNECTEDHEALTHY` = "CONNECTED_HEALTHY"
- `HostStatusCONNECTEDWARNING` = "CONNECTED_WARNING"
- `HostStatusCONNECTING`       = "CONNECTING"
- `HostStatusINITIALIZING`     = "INITIALIZING"
- `HostStatusSESSIONEXPIRED`   = "SESSION_EXPIRED"

**`models.OperateActionEnum`** （host 电源操作）：
- `OperateActionEnumPoweroff` = "poweroff"
- `OperateActionEnumReboot`   = "reboot"

### 0.7 utils 子包真实签名（task_utils.go）

```go
func WaitTask(ctx context.Context, client *apiclient.Cloudtower, id *string, interval time.Duration) error
func WaitTasks(ctx context.Context, client *apiclient.Cloudtower, ids []string, interval time.Duration) error
```

- 第三参数是 `*string`，**不是** `string`
- 内部会强制 `interval >= 1*time.Second`
- 失败时返回 `fmt.Errorf("task %s failed: %s", *id, *errMsg)`，goct 上层用 `errors.Is(err, adapter.ErrTaskFailed)` 时需自己判定状态字符串
- adapter 层若想暴露 `GetTaskProgress(id)`，自己用 `cli.Task.GetTasks` + `TaskWhereInput{ID: pointy.String(id)}` 实现，progress 字段在 `*models.Task.Progress *float64`

### 0.8 注意事项

1. 几乎所有 model 字段都是 `*string` / `*int32` / `*bool`，必须用 `pointy.String/Int32/Bool` 解引用
2. WhereInput 通常都有 `AND/OR/NOT []*XxxWhereInput` 复合条件字段
3. `Cloudtower` 类型在 `client/rest_client.go`，构造器是 `client.New(transport, formats)` 与 `client.NewWithUserConfig(...)`
4. 没有 `api_info.GetAPIInfo` 这种方法。**真实**：`apiInfo.NewGetAPIVersionParams()` → `cli.APIInfo.GetAPIVersion(p)` → `*GetAPIVersionOK{Payload string}`（payload 是裸 string，不是结构体）

---

## 1. client/api_info — Tower 版本

| 项 | 值 |
| --- | --- |
| Cloudtower 字段 | `cli.APIInfo` |
| 唯一方法 | `GetAPIVersion(*GetAPIVersionParams, ...)` `(*GetAPIVersionOK, error)` |
| Params 构造 | `apiInfo.NewGetAPIVersionParams()` |
| RequestBody | 无 |
| 返回 | `*GetAPIVersionOK{ Payload string }` —— 裸版本字符串，`/get-version` 端点 |

```go
import apiInfo "github.com/smartxworks/cloudtower-go-sdk/v2/client/api_info"

p := apiInfo.NewGetAPIVersionParams().WithContext(ctx)
res, err := cli.APIInfo.GetAPIVersion(p)
version := res.Payload  // 已是 string
```

> **adapter.About 实现要点**：直接把 `res.Payload` 当 Version 用；没有 Build 字段。`adapter.TowerInfo.Build` 留空即可。

---

## 2. client/user — 登录与用户

**ClientService 方法**：
- `CreateRootUser(*CreateRootUserParams) (*CreateRootUserOK, error)`
- `CreateUser(*CreateUserParams) (*CreateUserOK, error)`
- `DeleteUser(*DeleteUserParams) (*DeleteUserOK, error)`
- `GetMyInfo(*GetMyInfoParams) (*GetMyInfoOK, error)`
- `GetUsers(*GetUsersParams) (*GetUsersOK, error)`
- `GetUsersConnection(*GetUsersConnectionParams) (*GetUsersConnectionOK, error)`
- `Login(*LoginParams) (*LoginOK, error)`
- `UpdateUser(*UpdateUserParams) (*UpdateUserOK, error)`

| 操作 | RequestBody | Payload |
| --- | --- | --- |
| Login | `*models.LoginInput` | `*models.WithTaskLoginResponse{ Data *LoginResponse{ Token *string }, TaskID *string }` |
| GetUsers | `*models.GetUsersRequestBody` | `[]*models.User` |
| GetMyInfo | `*models.GetMyInfoRequestBody` | `*models.User` |
| CreateUser | `[]*models.UserCreationParams` | `[]*models.WithTaskUser` |
| UpdateUser | `*models.UserUpdationParams` | `[]*models.WithTaskUser` |
| DeleteUser | `*models.UserDeletionParams` | `[]*models.WithTaskDeleteUser` |
| CreateRootUser | `*models.RootUserCreationParams` | `*models.WithTaskUser` |

**`models.LoginInput`**：
```go
type LoginInput struct {
    AuthConfigID *string
    MfaType      *MfaType
    Password     *string  // Required
    Source       *UserSource
    Username     *string  // Required
}
```

**`models.UserCreationParams`** 字段：`AuthConfigID`/`EmailAddress`/`Internal`/`LdapDn`/`MobilePhone`/`Name`/`Password`/`RoleID`/...

**`models.UserDeletionParams`**：`Where *UserWhereInput`

**`UserWhereInput`** 关键字段：`ID *string` / `IDIn []string` / `Name *string` / `NameContains *string` / `Username *string`

---

## 3. client/user_role_next — 角色

**方法**：CreateRole / DeleteRole / GetUserRoleNexts / GetUserRoleNextsConnection / UpdateRole

| 操作 | RequestBody | Payload |
| --- | --- | --- |
| GetUserRoleNexts | `*models.GetUserRoleNextsRequestBody` | `[]*models.UserRoleNext` |
| CreateRole | `[]*models.RoleCreationParams` | `[]*models.WithTaskUserRoleNext` |
| UpdateRole | `*models.RoleUpdationParams` | `[]*models.WithTaskUserRoleNext` |
| DeleteRole | `*models.RoleDeletionParams` | `[]*models.WithTaskDeleteRole` |

**`models.RoleCreationParams`**：
```go
type RoleCreationParams struct {
    Actions []ROLEACTION  // Required
    Name    *string       // Required, MinLength=1
}
```

**`models.RoleDeletionParams`**：`Where *UserRoleNextWhereInput`

---

## 4. client/vm — 虚拟机（核心）

**18 个 ClientService 方法**（仅列 Tier-1 关心的）：

| goct 命令 | SDK 方法 | RequestBody | Payload | 说明 |
| --- | --- | --- | --- | --- |
| vm.ls | `GetVms` | `*models.GetVmsRequestBody` | `[]*models.VM` | |
| vm.info | `GetVms` (按 ID 过滤) | 同上 | 同上 | 复用 GetVms |
| vm.create | `CreateVM` | `[]*models.VMCreationParams` | `[]*models.WithTaskVM` | **接收数组** |
| vm.clone | `CloneVM` | `[]*models.VMCloneParams` | `*models.CloneVMOK{ Payload []*WithTaskVM }` | **接收数组** |
| vm.destroy | `DeleteVM` | `*models.VMDeleteParams` | `*models.DeleteVMOK` | 注意 Payload 名为 DeleteVMOK，包内仅一份；`Effect *VMDeleteParamsEffect{IncludeSnapshots *bool}` |
| vm.destroy (软删) | `MoveVMToRecycleBin` | `*models.VMOperateParams` | `[]*models.WithTaskVM` | |
| vm.migrate | `MigrateVM` | `*models.VMMigrateParams` | `[]*models.WithTaskVM` | `Data *VMMigrateParamsData{HostID *string}` |
| vm.export | `ExportVM` | `*models.VMExportParams`（自己查，非 Tier-1 重点） | | |
| vm.power.on | `StartVM` | `*models.VMStartParams` | `[]*models.WithTaskVM` | `Data *VMStartParamsData{HostID *string}` 可选指定主机 |
| vm.power.off | `ShutDownVM` | `*models.VMOperateParams` | `[]*models.WithTaskVM` | 优雅关机 |
| vm.power.off --force | `PoweroffVM` | `*models.VMOperateParams` | `[]*models.WithTaskVM` | 强制断电 |
| vm.power.reset | `RestartVM` | `*models.VMOperateParams` | `[]*models.WithTaskVM` | |
| vm.power.reset --force | `ForceRestartVM` | `*models.VMOperateParams` | `[]*models.WithTaskVM` | |
| vm.power.suspend | `SuspendVM` | `*models.VMOperateParams` | `[]*models.WithTaskVM` | |
| vm.power.resume | `ResumeVM` | `*models.VMOperateParams` | `[]*models.WithTaskVM` | |
| (rollback to snap) | `RollbackVM` | `*models.VMRollbackParams` | `[]*models.WithTaskVM` | `Data *VMRollbackParamsData{SnapshotID *string}` |
| (recycle bin recover) | `RecoverVMFromRecycleBin` | `*models.VMOperateParams` | `[]*models.WithTaskVM` | |

**`models.VMOperateParams`**：
```go
type VMOperateParams struct {
    Where *VMWhereInput  // Required
}
```

**`models.VMStartParams`**：
```go
type VMStartParams struct {
    Data  *VMStartParamsData  // 可选，指定 HostID
    Where *VMWhereInput        // Required
}
type VMStartParamsData struct { HostID *string /* Required */ }
```

**`models.VMMigrateParams`**：
```go
type VMMigrateParams struct {
    Data  *VMMigrateParamsData  // {HostID *string}
    Where *VMWhereInput          // Required
}
```

**`models.VMDeleteParams`**：
```go
type VMDeleteParams struct {
    Effect *VMDeleteParamsEffect  // {IncludeSnapshots *bool}
    Where  *VMWhereInput           // Required
}
```

**`models.GetVmsRequestBody`**：
```go
type GetVmsRequestBody struct {
    After   *string
    Before  *string
    First   *int32
    Last    *int32
    OrderBy *VMOrderByInput
    Skip    *int32
    Where   *VMWhereInput
}
```

**`VMWhereInput`** 关键字段（共数百个，列高频）：
- `ID *string` / `IDIn []string`
- `Name *string` / `NameContains *string` / `NameIn []string`
- `Status *VMStatus`
- `Description *string`
- `InRecycleBin *bool`
- `Cluster *ClusterWhereInput` (嵌套)
- `Folder *VMFolderWhereInput` (嵌套)
- `Host *HostWhereInput` (嵌套)
- 复合：`AND/OR/NOT []*VMWhereInput`

**`models.VM`** 关键字段（用于显示）：
- `ID *string`、`Name *string`、`Status *VMStatus`、`Description *string`
- `Cluster *NestedCluster`（含 ID/Name）、`Host *NestedHost`
- `CPU *NestedCPU`、`CPUModel *string`、`Vcpu *int32`
- `Memory *int64`（字节）
- `LocalID *string`、`NodeIP *string`
- `GuestOsType *VMGuestsOperationSystem`

**`models.WithTaskVM`**：`Data *VM` + `TaskID *string`

**`models.VMCreationParams`** 关键字段：
- `Name *string`、`ClusterID *string`、`Status *VMStatus`、`Ha *bool`
- `CPUCores *int32`、`CPUSockets *int32`、`Vcpu *int32`
- `Memory *int64`、`MemoryUnit *ByteUnit`
- `Firmware *VMFirmware`（必填）
- `VMDisks *VMDiskParams`（必填）、`VMNics []*VMNicParams`（必填）
- `Description *string`、`HostID *string`、`FolderID *string`、`Owner *VMOwnerParams`

> 注：`vm.create` 完整实现复杂，本次仅暴露最少 flag（name/cluster/cpu/memory），其他字段后续迭代。

---

## 5. client/vm_snapshot — VM 快照

**ClientService 方法**：
- `CreateVMSnapshot(*CreateVMSnapshotParams) (*CreateVMSnapshotOK, error)`
- `DeleteVMSnapshot(*DeleteVMSnapshotParams) (*DeleteVMSnapshotOK, error)`
- `GetVMSnapshots(*GetVMSnapshotsParams) (*GetVMSnapshotsOK, error)`
- `GetVMSnapshotsConnection(...)`

| goct 命令 | SDK 方法 | RequestBody | Payload |
| --- | --- | --- | --- |
| vm.snapshot.ls | `GetVMSnapshots` | `*models.GetVMSnapshotsRequestBody` | `[]*models.VMSnapshot` |
| vm.snapshot.create | `CreateVMSnapshot` | `*models.VMSnapshotCreationParams` | `[]*models.WithTaskVMSnapshot` |
| vm.snapshot.rm | `DeleteVMSnapshot` | `*models.VMSnapshotDeletionParams` | `[]*models.WithTaskDeleteVMSnapshot` |
| vm.snapshot.revert | **`vm.RollbackVM`** | `*models.VMRollbackParams` | `[]*models.WithTaskVM` | ⚠️ 注意：revert 不在 vm_snapshot 子包，在 vm 子包的 RollbackVM |

**`VMSnapshotCreationParams`**：
```go
type VMSnapshotCreationParams struct {
    Data []*VMSnapshotCreationParamsData  // Required
}
type VMSnapshotCreationParamsData struct {
    ConsistentType *ConsistentType
    Name           *string  // Required
    VMID           *string  // Required
}
```

**`VMSnapshotDeletionParams`**：`Where *VMSnapshotWhereInput`

**`VMSnapshotWhereInput`**：`ID *string` / `IDIn []string` / `Name *string` / `NameContains *string` / `VM *VMWhereInput`

---

## 6. client/host — 主机

**ClientService 方法**（11 个）：
- `CreateHost`, `EnterMaintenanceMode`, `EnterMaintenanceModePreCheck`, `EnterMaintenanceModePrecheckResult`, `ExitMaintenanceMode`, `ExitMaintenanceModePrecheckResult`, `GetHosts`, `GetHostsConnection`, `PowerOffHost`, `TriggerDiskBlink`, `UpdateHost`

| goct 命令 | SDK 方法 | RequestBody | Payload |
| --- | --- | --- | --- |
| host.ls | `GetHosts` | `*models.GetHostsRequestBody` | `[]*models.Host` |
| host.info | 同上 by-ID 过滤 | | |
| host.maintenance.enter | `EnterMaintenanceMode` | `*models.EnterMaintenanceModeParams` | `*models.WithTaskHost`（**单个，不是数组**） |
| host.maintenance.exit | `ExitMaintenanceMode` | `*models.ExitMaintenanceModeParams` | `*models.WithTaskHost` |
| host.shutdown / host.reboot | `PowerOffHost` | `*models.OperateHostPowerParams` | `[]*models.WithTaskHost` | Action 字段决定具体操作 |
| host.reconnect | 没有专属 method ⚠️ | 用 `UpdateHost` 改 status，或暂不实现 | |
| host.disconnect | 没有专属 method ⚠️ | 同上 | |
| (DiskBlink) | `TriggerDiskBlink` | `*models.TriggerDiskBlinkParams` | `[]*models.WithTaskHost` |

**`OperateHostPowerParams`**（**重要：决定 host.shutdown/host.reboot 的实现路径**）：
```go
type OperateHostPowerParams struct {
    Data  *OperateHostPowerData          // Required
    Where *OperateHostPowerParamsWhere   // Required
}
type OperateHostPowerData struct {
    Action *OperateActionEnum   // "poweroff" | "reboot"
    Force  *bool                 // Required
    Reason *string               // ...
}
type OperateHostPowerParamsWhere struct {
    HostID *string  // Required（不是 Where input，是直接给 host_id）
}
```

**`EnterMaintenanceModeParams`**：
```go
type EnterMaintenanceModeParams struct {
    Data  *EnterMaintenanceModeInput          // {ShutdownVms []string}
    Where *EnterMaintenanceModeParamsWhere    // {HostID *string}
}
```

**`Host`** 关键字段：
- `ID *string`、`Name *string`、`Status *HostStatus`
- `Cluster *NestedCluster`、`ManagementIP *string`、`LocalID *string`
- 没有直接 reconnect/disconnect 字段，goct 这两个命令初版可以返回 "not yet implemented in SDK Tier-1，请用 host.update" 占位

**`HostWhereInput`**：`ID *string` / `IDIn []string` / `Name *string` / `NameContains *string` / `NameIn []string` / `Status *HostStatus`

> **goct 设计妥协**：`host.reconnect` / `host.disconnect` 在 SDK v2.22.1 没有直接的方法。adapter 这两个 op 可以返回 `errors.New("host.reconnect not supported by SDK")`，cmd 层将错误映射为 "feature unavailable"，但 ls/info/maintenance/shutdown/reboot 5 个命令完全可用。

---

## 7. client/cluster — 集群

**ClientService 方法**（15 个，仅列 Tier-1 关心）：
- `GetClusters`, `GetClustersConnection`, `GetMetaLeader`
- `ConnectCluster`（添加集群）、`DeleteCluster`、`UpdateCluster`
- 多个 Update 子命令：`UpdateClusterDisablePinInPerformance`、`UpdateClusterEnableISCSISetting`、`UpdateClusterEnablePinInPerformance`、`UpdateClusterHaSetting`、`UpdateClusterLicense`、`UpdateClusterNetworkSetting`、`UpdateClusterVirtualizationSetting`
- `GetClusterPinInPerformanceInfo`、`GetClusterStorageInfo`

| goct 命令 | SDK 方法 | RequestBody | Payload |
| --- | --- | --- | --- |
| cluster.ls | `GetClusters` | `*models.GetClustersRequestBody` | `[]*models.Cluster` |
| cluster.info | 同上 by-ID | | |
| cluster.change | `UpdateCluster` | `*models.ClusterUpdationParams` | `[]*models.WithTaskCluster` |

**`ClusterWhereInput`** 关键字段：`ID *string` / `IDIn []string` / `Name *string` / `NameContains *string`

---

## 8. client/elf_data_store — Datastore

**ClientService 方法**：`GetElfDataStores`、`GetElfDataStoresConnection`（**只有读，无写**）

| goct 命令 | SDK 方法 | RequestBody | Payload |
| --- | --- | --- | --- |
| datastore.ls | `GetElfDataStores` | `*models.GetElfDataStoresRequestBody` | `[]*models.ElfDataStore` |
| datastore.info | 同上 by-ID | | |

**`ElfDataStoreWhereInput`**：`ID` / `IDIn` / `Name` / `NameContains`

---

## 9. client/disk — 物理盘

**ClientService 方法**：`GetDisks`, `GetDisksConnection`, `MountDisk`, `UnmountDisk`

| goct 命令 | SDK 方法 | RequestBody | Payload |
| --- | --- | --- | --- |
| datastore.disk.ls | `GetDisks` | `*models.GetDisksRequestBody` | `[]*models.Disk` |
| (mount) | `MountDisk` | `*models.DiskMountParams` | `[]*models.WithTaskDisk` |
| (unmount) | `UnmountDisk` | `*models.DiskUnmountParams` | `[]*models.WithTaskDisk` |

> Tier-1 只用 `GetDisks` 实现 `datastore.disk.ls`，挂载/卸载留 Tier-2。

---

## 10. client/vds — 分布式交换机（network）

**ClientService 方法**（7 个）：`CreateVds`, `CreateVdsWithAccessVlan`, `CreateVdsWithMigrateVlan`, `DeleteVds`, `GetVdses`, `GetVdsesConnection`, `UpdateVds`

| goct 命令 | SDK 方法 | RequestBody | Payload |
| --- | --- | --- | --- |
| network.ls | `GetVdses` | `*models.GetVdsesRequestBody` | `[]*models.Vds` |
| network.info | 同上 by-ID | | |

**`VdsWhereInput`**：`ID` / `IDIn` / `Name` / `NameContains` / `Cluster *ClusterWhereInput`

---

## 11. client/nic — 网卡（network 辅助）

**ClientService 方法**：`GetNics`, `GetNicsConnection`, `UpdateNic`

| 方法 | RequestBody | Payload |
| --- | --- | --- |
| GetNics | `*models.GetNicsRequestBody` | `[]*models.Nic` |
| UpdateNic | `*models.NicUpdationParams` | `[]*models.WithTaskNic` |

---

## 12. client/vlan — VLAN

**ClientService 方法**（7 个）：`CreateVMVlan`, `DeleteVlan`, `GetVlans`, `GetVlansConnection`, `UpdateManagementVlan`, `UpdateMigrationVlan`, `UpdateVlan`

| goct 命令 | SDK 方法 | RequestBody | Payload |
| --- | --- | --- | --- |
| vlan.ls | `GetVlans` | `*models.GetVlansRequestBody` | `[]*models.Vlan` |
| vlan.info | 同上 by-ID | | |
| vlan.create | `CreateVMVlan` | `[]*models.VMVlanCreationParams` | `[]*models.WithTaskVlan` |
| vlan.destroy | `DeleteVlan` | `*models.VlanDeletionParams` | `[]*models.WithTaskDeleteVlan` |

**`VMVlanCreationParams`**：
```go
type VMVlanCreationParams struct {
    ModeType   *VlanModeType
    Name       *string  // Required
    NetworkIds []string
    QosBurst   *int64
    // ...
}
```

**`VlanDeletionParams`**：`Where *VlanWhereInput`

**`VlanWhereInput`**：`ID` / `IDIn` / `Name` / `NameContains` / `Type *NetworkType` / `VlanID *int32`

---

## 13. client/task — 任务

**ClientService 方法**：`CreateTask`, `GetTasks`, `GetTasksConnection`, `UpdateTask`

| goct 命令 | SDK 方法 | RequestBody | Payload |
| --- | --- | --- | --- |
| task.ls | `GetTasks` | `*models.GetTasksRequestBody` | `[]*models.Task` |
| task.info | 同上 by-ID | | |
| task.cancel | ⚠️ SDK 无原生 cancel；可用 `UpdateTask` 改状态，但通常不允许 | 返回 "task.cancel not supported by tower" | |
| task.wait | 直接调 `utils.WaitTask` | (id *string, interval) | error |

**`Task`** 关键字段（用于显示）：
- `ID *string`、`Description *string`、`Status *TaskStatus`、`Progress *float64`
- `ErrorCode *string`、`ErrorMessage *string`
- `LocalCreatedAt *string`、`StartedAt *string`、`FinishedAt *string`
- `Cluster *NestedCluster`、`User *NestedUser`、`Type *TaskType`
- `ResourceID *string`、`ResourceType *string`、`ResourceMutation *string`
- `Steps []*NestedStep`

**`TaskWhereInput`**：`ID *string` / `IDIn []string` / `Status *TaskStatus`

---

## 14. client/alert — 告警

**ClientService 方法**：`GetAlerts`, `GetAlertsConnection`, `ResolveAlert`

| goct 命令 | SDK 方法 | RequestBody | Payload |
| --- | --- | --- | --- |
| alert.ls | `GetAlerts` | `*models.GetAlertsRequestBody` | `[]*models.Alert` |
| alert.info | 同上 by-ID | | |
| alert.ack | `ResolveAlert` | `*models.ResolveAlertParams` | `[]*models.WithTaskAlert` | （SDK 用 "Resolve"，goct 命令用 govc 风格的 "ack"） |

**`Alert`** 关键字段：
- `ID *string`、`Cause *string`、`Severity *string`、`Message *string`
- `Solution *string`、`Impact *string`、`Threshold *float64`、`Value *float64`
- `Ended *bool`、`CreateTime *string`、`LocalStartTime *string`、`LocalEndTime *string`
- `Cluster *NestedCluster`、`Host *NestedHost`、`Disk *NestedDisk`、`Vms []*NestedVM`
- `AlertRule *NestedAlertRule`、`Labels interface{}`

**`ResolveAlertParams`**：`Where *AlertWhereInput`

**`AlertWhereInput`**：`ID *string` / `IDIn []string` / `Severity *string` / `Cluster` / `Host` / `Disk`（嵌套）/ `Ended *bool` 等

---

## 15. client/alert_rule — 告警规则

**ClientService 方法**：`GetAlertRules`, `GetAlertRulesConnection`

只读（增删改在 global_alert_rule 子包）。Tier-1 不暴露独立命令，可作为 alert.info 的辅助查询。

---

## 附录 A：与 README 的差异（重要）

| 项 | README 声明 | 实际源码 | 行动 |
| --- | --- | --- | --- |
| api_info 方法 | `GetAPIInfo` | `GetAPIVersion` | adapter.About 改用 `GetAPIVersion` |
| api_info 返回 | 含 version/build 字段对象 | 裸 `string` | adapter.TowerInfo.Build 留空 |
| TaskStatus 成功 | 文档常写 `SUCCEEDED` | 实际 `SUCCESSED`（拼写错误但已固化） | task.Watcher 用 `SUCCESSED`，不要双兼容 |
| utils.WaitTask 第三参 | README 示例 `string` | 实际 `*string` | 调用时 `pointy.String(id)` |
| host.reconnect/disconnect | 期望有专属方法 | SDK 无 | adapter 暴露占位返回 ErrUnsupported |
| vm.snapshot.revert | 期望 vm_snapshot 子包提供 | 实际在 `vm.RollbackVM` | adapter 路由到 vm 子包 |
| task.cancel | govc 等价 | SDK 无原生 cancel | adapter 暴露占位返回 ErrUnsupported |
| ShutDownVM vs PoweroffVM | 同操作 | 实际两个不同方法：ShutDownVM=优雅；PoweroffVM=强制 | adapter.PowerVM 用 force 参数路由 |
| RestartVM vs ForceRestartVM | 同操作 | 同上分两个方法 | adapter.PowerVM Reset action + force |

## 附录 B：需要 Tier-2 才补的命令

下列 Tier-1 列表中的命令在 SDK v2.22.1 没有直接对应方法，建议放入 Tier-2 或本次返回 "unsupported"：
- `host.reconnect` / `host.disconnect`（SDK 无对应 op）
- `task.cancel`（SDK 不支持取消）
- `vm.export`（SDK `ExportVM` 行为复杂，需异步导出文件，本次仅占位）
- `vm.create` 完整实现（VMCreationParams 需 disk/nic 子结构，本次先支持必要字段）

## 附录 C：utils 子包小工具

- `utils.BigIntToString(*big.Int) string`
- `utils.StringToBigInt(string) (*big.Int, bool)`
- `utils.CompareBigIntStrings(*string, *string) (int, error)`
- iso_utils.go：包含一个简短的 ISO 工具函数（本次未涉及）

---

**审计追踪**

- 2026-04-29 10:55 通过 code-explorer subagent 探查 SDK v2.22.1 全部 Tier-1 子包真实签名后写入。

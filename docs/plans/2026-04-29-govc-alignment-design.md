# goct govc 对齐设计：对象传递 / 脚本化 / 信息完整度

| 字段 | 值 |
| --- | --- |
| 日期 | 2026-04-29 |
| 状态 | Approved |
| 触发 | 用户反馈：uuid 使用方式、脚本上下传递、vm.info 返回信息不够完整 |
| 补充确认 | ① 需支持 stdin 管道读取；② vm.info 字段由开发者判断运维所需；③ 模版使用内容库模版（ContentLibraryVmTemplate）而非传统 VmTemplate |

---

## 1. 核心需求理解

用户提出三个方面的改进需求：

1. **参考 govc 的对象标识方式**：govc 用 `GOVC_VM`、`-vm` flag 等方式传递 VM，而非只接受位置参数
2. **脚本上下传递**：govc 用 `find | xargs`、`ls --id-only | xargs` 等管道模式在命令间传递对象
3. **vm.info 信息完整度**：当前只输出 9 个字段，远不够运维使用

---

## 2. 改动方案总览

### 2.1 对象定位方式对齐 govc

**govc 模式**：
```bash
# 方式 1: 位置参数（当前 goct 已支持）
govc vm.info my-vm
govc vm.info cl5k7g2xo04070822fhxjfsev9q

# 方式 2: -vm flag（govc 核心特性 — goct 要加）
govc vm.power.on -vm my-vm

# 方式 3: GOVC_VM 环境变量（govc 核心特性 — goct 要加）
export GOVC_VM=my-vm
govc vm.info
govc vm.power.on
govc snapshot.create my-snap
```

**goct 对应设计**：

| 资源 | 环境变量 | 命令 flag | 位置参数 | 优先级 |
| --- | --- | --- | --- | --- |
| VM | `GOCT_VM` | — | `<name\|id>` | 位置参数 > 环境变量 |
| Host | `GOCT_HOST` | — | `<name\|id>` | 位置参数 > 环境变量 |

**实现规则**：
- 所有需要单个资源标识的 `vm.*` 命令（info/destroy/migrate/export/power.*/snapshot.*/clone），`Args` 从 `ExactArgs(1)` 改为 `MaximumNArgs(1)`
- 当 `args[0]` 为空时，fallback 到 `os.Getenv("GOCT_VM")`
- 抽取通用 helper 函数，避免每个命令重复代码

```go
// cmd/vm/helpers.go
package vm

import (
    "bufio"
    "errors"
    "os"
)

// resolveVMArg 从位置参数、GOCT_VM 环境变量或 stdin 获取 VM 标识。
// 优先级：位置参数 > 环境变量 > stdin（管道）。
// stdin 仅在非 TTY 时读取（即 echo id | goct vm.info）。
func resolveVMArg(args []string) (string, error) {
    if len(args) > 0 && args[0] != "" {
        return args[0], nil
    }
    if v := os.Getenv("GOCT_VM"); v != "" {
        return v, nil
    }
    // 尝试从 stdin 管道读取（非 TTY 时）
    if stat, _ := os.Stdin.Stat(); stat.Mode()&os.ModeCharDevice == 0 {
        scanner := bufio.NewScanner(os.Stdin)
        if scanner.Scan() {
            if line := strings.TrimSpace(scanner.Text()); line != "" {
                return line, nil
            }
        }
    }
    return "", errors.New("VM not specified: use positional arg, set GOCT_VM, or pipe via stdin")
}
```

同理 `resolveHostArg` 在 `cmd/host/helpers.go`。

### 2.2 脚本化管道传递

**govc 模式**：
```bash
# 列出 ID → 管道传递给其他命令
govc find . -type m -name "web-*" | xargs -n1 govc vm.info
govc ls -t VirtualMachine | head -5 | xargs -n1 govc vm.power.on
```

**goct 已有基础**：
- `vm.ls --id-only` 输出纯 ID 列表（仅 vm.ls 有）

**需要补全**：

1. **所有 ls 命令统一加 `--id-only`**：host.ls / cluster.ls / datastore.ls / network.ls / vlan.ls / task.ls / alert.ls / user.ls / vm.snapshot.ls

2. **写操作输出 entity ID**：变更命令（create/clone/destroy/power.*）成功后输出被操作的 entity ID 到 stdout，便于管道消费
   ```bash
   # 创建 VM，输出新 VM ID 供后续使用
   NEW_VM=$(goct vm.create --name web-1 --cluster prod --vcpu 4 --memory 8192)
   goct vm.power.on $NEW_VM
   ```

3. **抽取通用 `--id-only` 注册 helper**：
   ```go
   // pkg/flags/flags.go 追加
   func RegisterIDOnly(cmd *cobra.Command, idOnly *bool) {
       cmd.Flags().BoolVar(idOnly, "id-only", false,
           "Output only IDs, one per line (for scripting)")
   }
   ```

### 2.3 vm.info 信息完整度

**现状**（9 个字段）：
```
ID:          cl5k7g2xo...
Name:        my-vm
Status:      RUNNING
VCPU:        4
Memory:      8.0 GiB
IPs:         10.0.0.1, 10.0.0.2
Host:        host-1
Cluster:     cl5k7g2xo...   ← 只有 ID，没有名字
Description: My test VM
```

**目标**（对标 govc vm.info，约 22 个字段）：
```
ID:            cl5k7g2xo04070822fhxjfsev9q
Name:          my-vm
Status:        RUNNING
VCPU:          4
Memory:        8.0 GiB
Firmware:      BIOS
HA:            true
Guest OS:      LINUX
VM Tools:      RUNNING (v2.3.0)
IPs:           10.0.0.1, 10.0.0.2
DNS:           8.8.8.8
Host:          host-1 (10.0.0.100)
Cluster:       prod-cluster (cl5k7g...)
CPU Model:     Intel Xeon E5-2680
Disks:         3
NICs:          2
Provisioned:   120.0 GiB
Used:          45.2 GiB
In Recycle:    false
Protected:     false
Created:       2026-04-28 10:30:00
Description:   My test VM
```

**需要扩展的 `adapter.VM` 字段**：

当前 `adapter.VM` 只有 9 个字段。需要补充从 `models.VM` 可提取的关键字段：

| 新字段 | 类型 | 来源 SDK 字段 | 用途 |
| --- | --- | --- | --- |
| `Firmware` | `string` | `*VMFirmware` | BIOS/UEFI |
| `Ha` | `bool` | `*bool` | 高可用状态 |
| `GuestOS` | `string` | `*VMGuestsOperationSystem` | Guest OS 类型 |
| `VMToolsStatus` | `string` | `*VMToolsStatus` | VMTools 状态 |
| `VMToolsVersion` | `string` | `*string` | VMTools 版本 |
| `CPUModel` | `string` | `*string` | CPU 型号 |
| `DiskCount` | `int` | `len(VMDisks)` | 磁盘数量 |
| `NicCount` | `int` | `len(VMNics)` | 网卡数量 |
| `ProvisionedBytes` | `uint64` | `*int64` | 预分配大小 |
| `UsedBytes` | `uint64` | `*int64` | 实际使用大小 |
| `InRecycleBin` | `bool` | `*bool` | 是否在回收站 |
| `Protected` | `bool` | `*bool` | 是否受保护 |
| `DNSServers` | `string` | `*string` | DNS 服务器 |
| `Hostname` | `string` | `*string` | 主机名 |
| `CreatedAt` | `string` | `*string` (local_created_at) | 创建时间 |
| `ClusterName` | `string` | `NestedCluster.Name` | 集群名称（已有 ClusterID） |
| `HostID` | `string` | `NestedHost.ID` | 主机 ID（已有 HostName） |

---

## 3. 影响范围分析

### 3.1 需要修改的文件

| 层 | 文件 | 改动 |
| --- | --- | --- |
| **adapter** | `types.go` | VM struct 新增 ~15 个字段 |
| **adapter** | `vm.go` | `toVM()` 函数扩展，提取新字段 |
| **output** | `columns.go` | `VMInfoRows()` 扩展到 ~22 行，`VMListColumns` 可选增加 |
| **cmd/vm** | `helpers.go` | 新建，`resolveVMArg()` |
| **cmd/vm** | `info.go` | Args 改 MaximumNArgs(1)，用 resolveVMArg |
| **cmd/vm** | 所有需要 VM 参数的命令 | Args 改 MaximumNArgs(1)，用 resolveVMArg |
| **cmd/host** | `helpers.go` | 新建，`resolveHostArg()` |
| **cmd/host** | 所有需要 Host 参数的命令 | Args 改 MaximumNArgs(1)，用 resolveHostArg |
| **cmd/\*/ls.go** | 所有 ls 命令 | 统一加 `--id-only` |
| **pkg/flags** | `flags.go` | 新增 `RegisterIDOnly()` helper |
| **README.md** | | 更新环境变量表、使用示例 |

### 3.2 不需要修改的部分

- `pkg/service/resolver.go` — Resolve 泛型已经支持 name|id，无需改动
- `pkg/adapter/client.go` — 接口不变
- `pkg/task/` — 不影响
- `pkg/session/`, `pkg/config/` — 不影响

---

## 4. 实施步骤建议

分 3 个 commit 推进：

### Commit 1: `feat(ux): add GOCT_VM/GOCT_HOST env vars and --id-only for all ls`
- `cmd/vm/helpers.go` 新建 `resolveVMArg()`
- `cmd/host/helpers.go` 新建 `resolveHostArg()`
- 所有 vm.* 命令：Args → MaximumNArgs(1)，调用 resolveVMArg
- 所有 host.* 命令：类似改造
- 所有 ls 命令：统一加 `--id-only`
- `pkg/flags/flags.go` 新增 `RegisterIDOnly()`

### Commit 2: `feat(vm): enrich VM model and vm.info output`
- `adapter/types.go` 扩展 VM struct
- `adapter/vm.go` 扩展 `toVM()` 提取新字段
- `output/columns.go` 扩展 `VMInfoRows()` 到 ~22 行
- `cmd/vm/info.go` 更新 KV 格式化宽度

### Commit 3: `docs: update README with GOCT_VM/GOCT_HOST and scripting examples`
- 更新 README 的环境变量表
- 添加管道编排使用示例

---

## 5. 脚本编排示例（目标状态）

```bash
# 设置默认 VM，后续命令无需重复指定
export GOCT_VM=my-vm
goct vm.info
goct vm.power.on
goct vm.snapshot.create --name before-upgrade

# stdin 管道：直接传入 ID
echo "cl5k7g2xo04070822fhxjfsev9q" | goct vm.info
goct vm.ls --name web --id-only | head -1 | goct vm.info

# xargs 编排：批量关机所有 web 开头的 VM
goct vm.ls --name web --id-only | xargs -I{} goct vm.power.off {}

# 管道编排：获取所有 RUNNING 状态 VM 的详情
goct vm.ls --id-only | while read id; do goct vm.info "$id"; done

# 创建后立即操作
VM_ID=$(goct vm.clone source-vm --name new-vm)
goct vm.power.on $VM_ID

# JSON + jq 管道
goct vm.info my-vm --format json | jq '.IPs'
goct host.ls --format json | jq -r '.[].ID' | head -1

# 跨资源编排
HOST=$(goct host.ls --id-only | head -1)
goct vm.migrate my-vm --host $HOST
```

---

## 6. 关于虚拟机模版的备注

CloudTower 使用**内容库模版**（`ContentLibraryVmTemplate`）而非传统的 `VmTemplate`。
Phase 2 的 `vm.template.*` 命令应使用 `client/content_library_vm_template` 子包，而非 `client/vm_template`。
SDK 对应方法：`CloneContentLibraryVmTemplateFromVm` / `ConvertContentLibraryVmTemplateFromVm` / `DeleteContentLibraryVmTemplate` 等。

---

**已确认，开始实施。**

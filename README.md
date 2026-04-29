# goct — CloudTower CLI

> govc 风格的 SmartX CloudTower 命令行工具

## 特性

- **52 个 Tier-1 命令**，覆盖 VM / 快照 / 主机 / 集群 / 存储 / 网络 / VLAN / 任务 / 告警 / 用户
- **govc 点分命名**：`goct vm.ls`、`goct host.maintenance.enter`、`goct vm.power.on`
- **govc 风格对象传递**：位置参数 > `GOCT_VM`/`GOCT_HOST` 环境变量 > stdin 管道，支持脚本批量操作
- **三层架构**：`cmd/` → `pkg/service/` → `pkg/adapter/`，仅 adapter 层 import SDK
- **自动化友好**：table（默认）/ JSON 双输出；`--id-only` 每行输出 ID；错误走 stderr、数据走 stdout
- **Session 缓存**：自动缓存 token 到 `$XDG_CACHE_HOME/goct/`，401 自动刷新
- **自签名 TLS**：`--insecure` 跳过证书验证（内网常见场景）
- **Task Watcher**：变更命令自动追踪异步任务进度

## 安装

```bash
go install github.com/6547709/goct@latest
```

或从源码编译：

```bash
git clone https://github.com/6547709/goct.git
cd goct
go build -o goct .
```

## 配置

### 方式一：命令行参数

```bash
goct --url https://tower.example.com --username admin --password secret vm.ls
```

### 方式二：环境变量

```bash
export GOCT_URL=https://tower.example.com
export GOCT_USERNAME=admin
export GOCT_PASSWORD=secret
export GOCT_INSECURE=true    # 自签名证书
export GOCT_CLUSTER=clusterX # 默认集群（ID 或名称），设置后无需每次传 --cluster
export GOCT_SOURCE=local     # 登录源：local|ldap|sso|authn
goct vm.ls
```

> 💡 **设置 `GOCT_CLUSTER` 后**，所有需要集群参数的命令（如 `vm.create`）会自动使用该默认值，无需每次传递 `--cluster`。

### 方式三：配置文件 `~/.goct.yaml`

```yaml
url: https://tower.example.com
username: admin
password: secret
insecure: true
source: local   # local | ldap | sso | authn
```

**优先级**：CLI flag > 环境变量 > 配置文件

### 完整环境变量列表

| 变量 | 说明 | 默认值 |
|---|---|---|
| `GOCT_URL` | CloudTower 地址 | - |
| `GOCT_USERNAME` | 登录用户名 | - |
| `GOCT_PASSWORD` | 登录密码 | - |
| `GOCT_CLUSTER` | 默认集群 ID 或名称 | - |
| `GOCT_INSECURE` | 跳过 TLS 验证（`true`/`false`） | `false` |
| `GOCT_SOURCE` | 登录源（`local`/`ldap`/`sso`/`authn`） | `local` |
| `GOCT_LOG` | 日志级别（`TRACE`/`DEBUG`/`INFO`/`WARN`/`ERROR`） | 关闭 |
| `GOCT_LOG_FILE` | 日志文件路径（默认 stderr） | stderr |
| `GOCT_VM` | VM 标识（ID 或名称），可替代命令位置参数 | - |
| `GOCT_HOST` | 主机标识（ID 或名称），可替代命令位置参数 | - |

## 命令矩阵

| 资源 | 命令 |
|---|---|
| **系统** (5) | `about` `version` `session.login` `session.logout` `session.ls` |
| **VM** (12) | `vm.ls` `vm.info` `vm.create` `vm.clone` `vm.destroy` `vm.migrate` `vm.export` `vm.power.on` `vm.power.off` `vm.power.reset` `vm.power.suspend` `vm.power.resume` |
| **快照** (4) | `vm.snapshot.ls` `vm.snapshot.create` `vm.snapshot.revert` `vm.snapshot.rm` |
| **主机** (8) | `host.ls` `host.info` `host.maintenance.enter` `host.maintenance.exit` `host.shutdown` `host.reboot` `host.reconnect` `host.disconnect` |
| **集群** (2) | `cluster.ls` `cluster.info` |
| **存储** (3) | `datastore.ls` `datastore.info` `datastore.disk.ls` |
| **网络** (2) | `network.ls` `network.info` |
| **VLAN** (4) | `vlan.ls` `vlan.info` `vlan.create` `vlan.destroy` |
| **任务** (4) | `task.ls` `task.info` `task.cancel` `task.wait` |
| **告警** (3) | `alert.ls` `alert.info` `alert.ack` |
| **用户** (4) | `user.ls` `user.info` `user.create` `user.destroy` |

## 使用示例

```bash
# 查看 CloudTower 版本
goct about

# 列出所有虚拟机
goct vm.ls

# 按名称筛选
goct vm.ls --name web-server

# JSON 输出（管道友好）
goct vm.ls --format json | jq '.[] | .Name'

# 查看 VM 详情
goct vm.info my-vm

# 开机
goct vm.power.on my-vm

# 强制关机
goct vm.power.off my-vm --force

# 创建 VM
goct vm.create --name new-vm --cluster <cluster-id> --vcpu 4 --memory 8192

# 克隆 VM
goct vm.clone source-vm --name cloned-vm

# 创建快照
goct vm.snapshot.create my-vm --name before-upgrade

# 回滚快照
goct vm.snapshot.revert <snapshot-id> --vm my-vm

# 进入维护模式
goct host.maintenance.enter my-host

# 等待任务完成
goct task.wait <task-id>

# 确认告警
goct alert.ack <alert-id>

# 管理 session
goct session.login                    # 强制重新登录
goct session.ls                       # 查看缓存
goct session.logout --url https://tower.example.com --user admin
```

## 脚本化

goct 模仿 govc 的 UX 模式，支持对象在脚本中上下传递。

### 对象解析优先级

对于 `vm.*` 和 `host.*` 命令，标识符解析顺序为：

1. **命令行位置参数**（如 `goct vm.info my-vm`）
2. **环境变量** `GOCT_VM` / `GOCT_HOST`
3. **stdin 管道**（如 `echo my-vm | goct vm.info`）

> 与 govc 完全一致：`GOVC_VM` > stdin > 位置参数（govc 是 env 优先于位置参数）

### GOCT_VM / GOCT_HOST 环境变量

设置后，命令无需每次传 VM/Host 标识：

```bash
export GOCT_VM=my-vm
goct vm.info              # 读 GOCT_VM
goct vm.power.off --force # 读 GOCT_VM
unset GOCT_VM
```

### stdin 管道

**仅当 stdin 非 TTY（管道/重定向）时才会读取**，不会阻塞交互式终端：

```bash
# 从 vm.ls --id-only 获取 ID，传递给 vm.info
goct vm.ls --id-only | while read id; do
  echo "=== VM Info ==="
  echo "$id" | goct vm.info
done

# 结合 jq 过滤后再 pipe
goct vm.ls --format json | jq -r '.[] | select(.Status=="RUNNING") | .ID' | goct vm.power.off --force

# 批量从文件读取
cat vms.txt | goct vm.info
```

### --id-only 输出

所有 `*.ls` 命令支持 `--id-only`，每行仅输出 ID，适合脚本处理：

```bash
# 导出所有 VM ID
goct vm.ls --id-only > vm-ids.txt

# 批量操作
for id in $(goct vm.ls --id-only); do
  goct vm.power.off $id --force
done

# 结合 grep 过滤
goct host.ls --id-only | grep "192.168."
```

## 退出码

| 码 | 含义 |
|---|---|
| 0 | 成功 |
| 1 | 一般错误 |
| 2 | 认证失败 |
| 3 | 资源未找到 |
| 4 | 异步任务失败 |

## 调试日志

通过环境变量 `GOCT_LOG` 启用分级日志（参考 Packer 的 `PACKER_LOG` 设计）：

```bash
# 关闭日志（默认）
goct vm.ls

# INFO 级别：登录成功、命令执行等关键事件
GOCT_LOG=INFO goct vm.ls

# DEBUG 级别：config 解析、session 缓存命中/miss、adapter 调用详情
GOCT_LOG=DEBUG goct vm.ls

# TRACE 级别：最详细，含 SDK 请求/响应细节
GOCT_LOG=TRACE goct vm.ls

# 将日志写入文件（而非 stderr）
GOCT_LOG=DEBUG GOCT_LOG_FILE=/tmp/goct.log goct vm.ls
```

| 级别 | 内容 |
|---|---|
| `TRACE` | SDK 请求/响应、参数展开 |
| `DEBUG` | config 解析、session 缓存命中/miss、adapter 调用 |
| `INFO` | 登录成功、命令执行 |
| `WARN` | session 过期、token 刷新 |
| `ERROR` | 仅错误 |

> 日志全部走 stderr（或 `GOCT_LOG_FILE`），**不会污染 stdout 的业务输出**，`goct vm.ls | jq` 安全。

## 架构

```
cmd/          → cobra 命令定义（每命令独立 .go 文件）
pkg/service/  → 业务编排（name|id 解析、调 adapter、watch task）
pkg/adapter/  → SDK 防腐层（唯一 import cloudtower-go-sdk）
pkg/output/   → table / JSON 双模式渲染
pkg/task/     → 异步任务进度 watcher
pkg/flags/    → 全局 flag 定义
pkg/config/   → 三级配置合并（flag > env > file）
pkg/session/  → XDG cache token 管理
pkg/client/   → session 命中 → 登录 → cache 回写
```

## 依赖

- Go 1.21+
- [cloudtower-go-sdk/v2](https://github.com/smartxworks/cloudtower-go-sdk) v2.22.1
- [cobra](https://github.com/spf13/cobra) / [viper](https://github.com/spf13/viper)
- [tablewriter](https://github.com/olekukonko/tablewriter) v1.x

## License

MIT

# goct 测试脚本

本目录的脚本用来对一台真实的 CloudTower 实例做端到端测试。
单元测试请直接 `go test ./...`（CI 会跑）。

## 脚本一览

| 脚本 | 是否改资源 | 用途 |
|---|---|---|
| `lib.sh` | — | 公共工具函数（日志、断言、配置归一化），其他脚本 source。 |
| `test_smoke.sh` | 只读 | 冒烟测试：`ls / info / find / events / vm.ip` 等只读路径，覆盖几乎所有命令。 |
| `test_regression.sh` | 只读 + 1 个 VM 描述写 | 锁定 `docs/REVIEW-2026-05-12.md` 中 14 个 Bug 的运行时回归。 |
| `test_vm_lifecycle.sh` | **创建/销毁 VM** | 端到端 VM 生命周期：create → update → disk/nic/cdrom → power → snapshot → recycle → destroy。 |
| `test.sh` | 取决子脚本 | 一键 runner，先跑 `go test ./...`，再按模式跑子脚本。 |

## 使用

### 配置环境变量

```bash
export GOCT_URL=https://tower.example.com
export GOCT_USERNAME=admin
export GOCT_PASSWORD=secret
export GOCT_INSECURE=true     # 自签名证书
# 仅 lifecycle 需要：
export GOCT_CLUSTER=cluster0
export GOCT_VLAN=vlan0        # 或者 export GOCT_TEMPLATE=tpl0
```

### 运行

```bash
# 默认 = smoke + regression（不会改资源）
./scripts/test.sh

# 只跑冒烟
./scripts/test.sh smoke
./scripts/test_smoke.sh demo01 host01 cluster0   # 也可以直接调

# 只跑回归
./scripts/test.sh regression

# 跑 VM 生命周期（**会创建并销毁 VM**）
./scripts/test.sh lifecycle

# 全部跑（smoke + regression + lifecycle）
./scripts/test.sh all
```

### 调试

```bash
DEBUG=true ./scripts/test_smoke.sh
KEEP_VM=true ./scripts/test_vm_lifecycle.sh    # 保留测试 VM 用于人工检查
SKIP_NETWORK=true ./scripts/test_smoke.sh      # 不连 CloudTower，只跑 --help / 退出码用例
```

## 退出码

| Code | 含义 |
|---|---|
| 0 | 全部通过 |
| 1 | 至少一项 FAIL |
| 2 | 配置错误（缺 `GOCT_URL` 等） |

## 设计要点

- 所有脚本走 `set -uo pipefail`，**不**用 `-e`：单个用例失败不应中断脚本，必须把全部用例跑完后给出聚合报告。
- 断言函数（`run_cmd / expect_fail / expect_exit / expect_json / expect_contains`）由 `lib.sh` 统一提供。
- 所有命令前缀都是 `$GOCT`（`init_goct_env` 自动拼好 `--insecure`）。
- 名字识别走 v0.2.1 修复后的 `IsID`（`^cl[0-9a-z]{25}$` 或 UUID），脚本里用 `--format json | jq` 拿 ID 后再传给后续命令，避免依赖表格列号。

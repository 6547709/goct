#!/usr/bin/env bash
# scripts/test_smoke.sh — 只读冒烟测试，**不会改变任何资源**。
#
# 覆盖：
#   - 系统命令（version、session、about）
#   - 各资源 ls / info / 多种输出格式
#   - 新增命令 find / events / 增强后的 vm.ip
#   - --help / 退出码 / cli flag 解析（不打远端）
#
# 用法：
#   GOCT_URL=... GOCT_USERNAME=... GOCT_PASSWORD=... GOCT_INSECURE=true \
#       ./scripts/test_smoke.sh [vm-name] [host-name] [cluster-name]
#
# 环境变量：
#   DEBUG=true        显示调试日志
#   SKIP_NETWORK=true 跳过所有需要登录 CloudTower 的测试（仅做 --help / 退出码 / CLI 形态校验）
#
# 退出码：
#   0 全部通过
#   1 至少一个用例 FAIL
#   2 配置错误（缺 GOCT_URL 等）

set -uo pipefail
cd "$(dirname "$0")/.."

source "scripts/lib.sh"

# 测试目标，按位置参数 / 环境变量给出，留空则跳过相关用例。
TEST_VM="${1:-${TEST_VM:-}}"
TEST_HOST="${2:-${TEST_HOST:-}}"
TEST_CLUSTER="${3:-${TEST_CLUSTER:-}}"

init_goct_env

# =========================================================
# Phase A — 离线测试（不需要 CloudTower 连接）
# =========================================================

log_section "Phase A — Offline (no network)"
expect_help_no_error "goct --help"            # root help
expect_help_no_error "version --help"         version
expect_help_no_error "find --help"            find
expect_help_no_error "events --help"          events
expect_help_no_error "vm.ip --help"           vm.ip
expect_help_no_error "vm.create --help"       vm.create
expect_help_no_error "vm.update --help"       vm.update
expect_help_no_error "vm.disk.update --help"  vm.disk.update
expect_help_no_error "vm.nic.update --help"   vm.nic.update
expect_help_no_error "vm.migrate --help"      vm.migrate
expect_help_no_error "events --help"          events

# version 必须能离线运行
run_cmd "version (offline)" $GOCT version

# 不带 --url 应当报参数缺失（除非用户已经 export GOCT_URL）。这里我们已经 init_goct_env，
# 所以 about 应该能跑通；--help 必然 0。
expect_exit 0 "about --help" $GOCT about --help

if [[ "${SKIP_NETWORK:-false}" == "true" ]]; then
    log_warn "SKIP_NETWORK=true → skipping online tests"
    print_summary; exit $?
fi

# =========================================================
# Phase B — 系统命令
# =========================================================

log_section "Phase B — System / Session"
run_cmd  "session.login"              $GOCT session.login
run_cmd  "session.ls"                 $GOCT session.ls
run_cmd  "about"                      $GOCT about
expect_json "about --format json"     $GOCT about --format json

# 自动获取 cluster ID 供后续命令使用
FIRST_CLUSTER_ID=$($GOCT cluster.ls --format json 2>/dev/null | jq -r '.[0].ID // empty' 2>/dev/null || true)

# =========================================================
# Phase C — 各资源 ls + JSON 输出
# =========================================================

log_section "Phase C — Resource listings"

declare -a LS_CMDS=(
    "vm.ls"
    "host.ls"
    "cluster.ls"
    "datastore.ls"
    "datastore.disk.ls"
    "storage.pool.ls"
    "network.ls"
    "vlan.ls"
    "task.ls"
    "alert.ls"
    "alert-rule.ls"
    "user.ls"
    "deploy.ls"
    "license.ls"
    "ntp.get"
    "content-library-image.ls"
)

for cmd in "${LS_CMDS[@]}"; do
    run_cmd     "$cmd"             $GOCT $cmd
    expect_json "$cmd --json"      $GOCT $cmd --format json
done

# id-only 输出格式
run_cmd  "vm.ls --id-only"   $GOCT vm.ls --id-only
run_cmd  "host.ls --id-only" $GOCT host.ls --id-only

# =========================================================
# Phase D — 新增命令：find / events
# =========================================================

log_section "Phase D — find / events (new commands, govc-style)"

# find 跨类型枚举
run_cmd     "find (all types, default limit)"        $GOCT find --limit 5
run_cmd     "find -type m"                           $GOCT find --type m --limit 5
run_cmd     "find -type h"                           $GOCT find --type h --limit 5
run_cmd     "find -type c"                           $GOCT find --type c --limit 5
run_cmd     "find -type d"                           $GOCT find --type d --limit 5
run_cmd     "find -type n"                           $GOCT find --type n --limit 5
run_cmd     "find --id-only --type m --limit 3"      $GOCT find --type m --limit 3 --id-only
expect_json "find --json --type m --limit 3"         $GOCT find --type m --limit 3 --json
# 错误的类型应当报错
expect_fail "find -type bogus → fail"                $GOCT find --type bogus

# events 来自 user_audit_log
run_cmd     "events (last 5)"                        $GOCT events --limit 5
expect_json "events --json --limit 3"                $GOCT events --limit 3 --json

# =========================================================
# Phase E — VM 详情 / IP / 设备列表（依赖 TEST_VM）
# =========================================================

if [[ -n "$TEST_VM" ]]; then
    log_section "Phase E — VM-targeted (TEST_VM=$TEST_VM)"
    run_cmd     "vm.info $TEST_VM"                       $GOCT vm.info "$TEST_VM"
    expect_json "vm.info $TEST_VM --json"                $GOCT vm.info "$TEST_VM" --format json
    run_cmd     "vm.ip $TEST_VM --no-wait"               $GOCT vm.ip "$TEST_VM" --no-wait
    run_cmd     "vm.ip $TEST_VM -a --no-wait"            $GOCT vm.ip "$TEST_VM" --all --no-wait
    run_cmd     "vm.ip $TEST_VM --v4 --no-wait"          $GOCT vm.ip "$TEST_VM" --v4 --no-wait
    run_cmd     "vm.ip $TEST_VM --v6 --no-wait"          $GOCT vm.ip "$TEST_VM" --v6 --no-wait
    run_cmd     "vm.disk.ls --vm $TEST_VM"               $GOCT vm.disk.ls --vm "$TEST_VM"
    run_cmd     "vm.nic.ls --vm $TEST_VM"                $GOCT vm.nic.ls --vm "$TEST_VM"
    run_cmd     "vm.gpu.ls --vm $TEST_VM"                $GOCT vm.gpu.ls --vm "$TEST_VM"
    run_cmd     "vm.cdrom.ls --vm $TEST_VM"              $GOCT vm.cdrom.ls --vm "$TEST_VM"
    run_cmd     "vm.snapshot.ls $TEST_VM"                $GOCT vm.snapshot.ls "$TEST_VM"
    run_cmd     "vm.vnc $TEST_VM"                        $GOCT vm.vnc "$TEST_VM"
    run_cmd     "events $TEST_VM (resource scoped)"      $GOCT events "$TEST_VM" --limit 5
    run_cmd     "find -name $TEST_VM"                    $GOCT find --name "$TEST_VM"

    # 不存在的 VM → ErrNotFound → exit 3
    expect_exit 3 "vm.info <nonexistent> → exit 3"       $GOCT vm.info "__nonexistent_vm_xyzzy__"
else
    log_skip "Phase E skipped (set TEST_VM=<name> to enable)"
fi

# =========================================================
# Phase F — Host / Cluster 定向用例
# =========================================================

if [[ -n "$TEST_HOST" ]]; then
    log_section "Phase F1 — Host (TEST_HOST=$TEST_HOST)"
    run_cmd     "host.info $TEST_HOST"          $GOCT host.info "$TEST_HOST"
    expect_json "host.info $TEST_HOST --json"   $GOCT host.info "$TEST_HOST" --format json
else
    log_skip "Phase F1 skipped (set TEST_HOST to enable)"
fi

if [[ -n "$TEST_CLUSTER" ]]; then
    log_section "Phase F2 — Cluster (TEST_CLUSTER=$TEST_CLUSTER)"
    run_cmd     "cluster.info $TEST_CLUSTER"          $GOCT cluster.info "$TEST_CLUSTER"
    expect_json "cluster.info $TEST_CLUSTER --json"   $GOCT cluster.info "$TEST_CLUSTER" --format json
else
    log_skip "Phase F2 skipped (set TEST_CLUSTER to enable)"
fi

# 自动选第一个 host / cluster 做 ID vs Name 解析回归
log_section "Phase G — ID vs Name auto-resolve (Bug 13 regression)"
FIRST_HOST_ID=$($GOCT host.ls --format json 2>/dev/null | jq -r '.[0].ID // empty' 2>/dev/null || true)
FIRST_HOST_NAME=$($GOCT host.ls --format json 2>/dev/null | jq -r '.[0].Name // empty' 2>/dev/null || true)
if [[ -n "$FIRST_HOST_ID" && -n "$FIRST_HOST_NAME" ]]; then
    run_cmd "host.info by ID  ($FIRST_HOST_ID)"   $GOCT host.info "$FIRST_HOST_ID"
    run_cmd "host.info by Name ($FIRST_HOST_NAME)" $GOCT host.info "$FIRST_HOST_NAME"
fi
FIRST_CLUSTER_ID=$($GOCT cluster.ls --format json 2>/dev/null | jq -r '.[0].ID // empty' 2>/dev/null || true)
FIRST_CLUSTER_NAME=$($GOCT cluster.ls --format json 2>/dev/null | jq -r '.[0].Name // empty' 2>/dev/null || true)
if [[ -n "$FIRST_CLUSTER_ID" && -n "$FIRST_CLUSTER_NAME" ]]; then
    run_cmd "cluster.info by ID   ($FIRST_CLUSTER_ID)"   $GOCT cluster.info "$FIRST_CLUSTER_ID"
    run_cmd "cluster.info by Name ($FIRST_CLUSTER_NAME)" $GOCT cluster.info "$FIRST_CLUSTER_NAME"
fi

# =========================================================
# Phase G2 — cluster-settings.get（依赖 FIRST_CLUSTER_ID）
# =========================================================

if [[ -n "$FIRST_CLUSTER_ID" ]]; then
    log_section "Phase G2 — cluster-settings.get"
    run_cmd     "cluster-settings.get --cluster $FIRST_CLUSTER_ID"   $GOCT cluster-settings.get --cluster "$FIRST_CLUSTER_ID"
    run_cmd     "cluster-settings.get --json --cluster $FIRST_CLUSTER_ID"   $GOCT cluster-settings.get --cluster "$FIRST_CLUSTER_ID" --format json
fi

# =========================================================
# Phase H — Metrics
# =========================================================

log_section "Phase H — Metrics (read-only)"
run_cmd     "vm.metrics --list"      $GOCT vm.metrics --list
run_cmd     "host.metrics --list"    $GOCT host.metrics --list
run_cmd     "cluster.metrics --list" $GOCT cluster.metrics --list
run_cmd     "volume.metrics --list"  $GOCT volume.metrics --list

# =========================================================
# Summary
# =========================================================

print_summary

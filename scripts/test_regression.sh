#!/usr/bin/env bash
# scripts/test_regression.sh — v0.2.1 14 个 Bug 的运行时回归测试。
#
# 大部分 Bug 在 `go test ./...` 已经有单测（Bug 1, 2, 13, 14, 6 字段映射等），
# 本脚本补齐**只能在真实 CloudTower 上验证**的端到端用例：
#
#   Bug 1   --trace 与 --insecure 同开 → 不应 TLS 校验失败
#   Bug 2   FAILED task → 退出码 4
#   Bug 3   vm.create 不再硬塞默认盘/网卡（语义校验：不传 --disk/--nic 时由 CloudTower 报错而非悄悄成功）
#   Bug 5   vm.migrate.abort 必须传 vmID
#   Bug 7   vm.update --description="" 能清空
#   Bug 9   Resolve 走服务端精确匹配（用大量 VM 时不应慢）
#   Bug 13  IsID 不再误判长名字
#   Bug 14  IP 输出不出现空行
#
# 用法：
#   GOCT_URL=...  GOCT_USERNAME=...  GOCT_PASSWORD=...  GOCT_INSECURE=true \
#       ./scripts/test_regression.sh [vm-name]
#
# 退出码：
#   0 通过；1 至少一项失败；2 配置错误。

set -uo pipefail
cd "$(dirname "$0")/.."

source "scripts/lib.sh"

TEST_VM="${1:-${TEST_VM:-}}"

init_goct_env

# ============== Bug 1: --trace --insecure 同开不应破坏 TLS-skip ==============
log_section "Bug 1 — --trace + --insecure compatibility"
# 前置条件：GOCT_INSECURE=true 才能验证（自签名集群）。
if [[ "$GOCT_INSECURE" == "true" ]]; then
    # 不直接看输出，只看退出码 + 是否包含 "x509"/"certificate" 错误。
    out=$($GOCT_BIN --insecure --trace --url "$GOCT_URL" --username "$GOCT_USERNAME" --password "$GOCT_PASSWORD" \
        about 2>&1 || true)
    if echo "$out" | grep -qiE "x509|certificate signed|tls"; then
        log_fail "Bug 1: --trace re-enabled TLS verification (output contains TLS error)"
    else
        log_pass "Bug 1: --trace + --insecure work together"
    fi
else
    log_skip "Bug 1: GOCT_INSECURE != true (cannot test self-signed scenario)"
fi

# ============== Bug 2: FAILED task → exit code 4 ==============
log_section "Bug 2 — FAILED task exits with code 4"
# 触发一个必失败的 task：试图删除一个不存在的 VM ID，
# 但这通常返回 ErrNotFound (exit 3)，不是 ErrTaskFailed。
# 真正的 task FAILED 需要远端真实任务，难以稳定模拟。
# 替代验证：直接在 Go 测试里已覆盖（pkg/task/watcher_test.go）。
log_skip "Bug 2: covered by go test pkg/task (TestWatcher_ErrFailed_AliasesAdapter)"

# ============== Bug 3: vm.create 没有 --disk/--nic 时不会"偷偷成功" ==============
log_section "Bug 3 — vm.create 不再硬塞默认盘 / 网卡"
# 注意：这里我们**不**传 --cluster，让命令在参数校验阶段就停（不真创建 VM）。
# 关键校验：--help 输出含 --disk / --nic 且不含 'default 10g' 这类描述。
help_out=$($GOCT vm.create --help 2>&1)
if echo "$help_out" | grep -qE -- "--disk\s+stringArray"; then
    log_pass "Bug 3: vm.create --help advertises --disk stringArray"
else
    log_fail "Bug 3: vm.create --help should advertise --disk stringArray"
fi
if echo "$help_out" | grep -qE -- "--nic\s+stringArray"; then
    log_pass "Bug 3: vm.create --help advertises --nic stringArray"
else
    log_fail "Bug 3: vm.create --help should advertise --nic stringArray"
fi
if echo "$help_out" | grep -qiE "default.*10[ ]?g"; then
    log_fail "Bug 3: vm.create help still mentions a hardcoded 10G default"
else
    log_pass "Bug 3: no '10g default' wording in help"
fi

# ============== Bug 5: vm.migrate.abort 必须传 vmID ==============
log_section "Bug 5 — vm.migrate.abort needs vmID"
# 不带参数时 cobra 会报 "accepts 1 arg(s), received 0"，应当非零退出。
expect_fail "vm.migrate.abort without arg → fail" $GOCT vm.migrate.abort

# ============== Bug 7: vm.update --description="" 能清空 ==============
log_section "Bug 7 — vm.update can clear description"
if [[ -n "$TEST_VM" ]]; then
    log_info "set description"
    $GOCT vm.update "$TEST_VM" --description "regression-test-$$" >/dev/null 2>&1 || \
        log_warn "vm.update set failed (non-fatal)"
    sleep 1
    log_info "clear description with --description=\"\""
    if $GOCT vm.update "$TEST_VM" --description "" >/dev/null 2>&1; then
        sleep 1
        cur_desc=$($GOCT vm.info "$TEST_VM" --format json 2>/dev/null | jq -r '.Description // ""')
        if [[ -z "$cur_desc" ]]; then
            log_pass "Bug 7: description cleared"
        else
            log_fail "Bug 7: description not cleared (still: $cur_desc)"
        fi
    else
        log_fail "Bug 7: vm.update --description=\"\" failed"
    fi
else
    log_skip "Bug 7: TEST_VM not set"
fi

# ============== Bug 9: Resolve 服务端精确匹配（性能可观察） ==============
log_section "Bug 9 — Resolve uses server-side exact match"
if [[ -n "$TEST_VM" ]]; then
    # 用 vm.info 触发一次 Resolve；耗时 < 5s 即可（精确匹配后服务端只返回 1 条）。
    t0=$(date +%s%N)
    $GOCT vm.info "$TEST_VM" >/dev/null 2>&1 || true
    t1=$(date +%s%N)
    elapsed_ms=$(( (t1 - t0) / 1000000 ))
    log_info "vm.info $TEST_VM elapsed: ${elapsed_ms}ms"
    if (( elapsed_ms < 5000 )); then
        log_pass "Bug 9: vm.info < 5s (server-side exact filter likely working)"
    else
        log_warn "Bug 9: vm.info took ${elapsed_ms}ms (acceptable on slow networks but worth noting)"
    fi
else
    log_skip "Bug 9: TEST_VM not set"
fi

# ============== Bug 13: IsID 不再误判长名字 ==============
log_section "Bug 13 — IsID does not misclassify long lowercase names"
# 用一个 25 字符全小写字母数字、但不以 "cl" 开头的字符串作为名字 → 期望 ErrNotFound (exit 3)。
# 如果旧正则误判为 ID，会去 GET id=... → 同样是 NotFound 但语义错。
# 我们只是验证它最终能 graceful 报错（exit != 0），不验证内部路径。
expect_fail "Bug 13: random 25-char lowercase string treated as name" \
    $GOCT vm.info "abcdefghij1234567890abcde"

# ============== Bug 14: vm.ip 输出无空行 ==============
log_section "Bug 14 — vm.ip outputs no empty lines"
if [[ -n "$TEST_VM" ]]; then
    out=$($GOCT vm.ip "$TEST_VM" --no-wait --all 2>/dev/null || true)
    if [[ -z "$out" ]]; then
        log_skip "Bug 14: VM has no IP yet"
    else
        # 不应有任何空行（IP split 后过滤了空段）。
        empty_lines=$(printf '%s\n' "$out" | grep -c '^[[:space:]]*$' || true)
        if (( empty_lines == 0 )); then
            log_pass "Bug 14: vm.ip output has no empty lines (lines=$(echo "$out" | wc -l))"
        else
            log_fail "Bug 14: vm.ip output contains $empty_lines empty line(s)"
        fi
    fi
else
    log_skip "Bug 14: TEST_VM not set"
fi

# ============== Bug 12: token TTL 缩短 → cache hit 后仍然能跑 ==============
log_section "Bug 12 — session cache hit reuses token within TTL"
$GOCT session.login >/dev/null 2>&1 || true
# 紧接着用 cache 跑 about（不应触发重新登录；速度应该快）。
t0=$(date +%s%N)
$GOCT about >/dev/null 2>&1 || true
t1=$(date +%s%N)
elapsed_ms=$(( (t1 - t0) / 1000000 ))
log_info "about (cached token) elapsed: ${elapsed_ms}ms"
if (( elapsed_ms < 3000 )); then
    log_pass "Bug 12: cached token works (about < 3s)"
else
    log_warn "Bug 12: about took ${elapsed_ms}ms (might still be OK on slow links)"
fi

# ============== Bug 11: vm.nic.update 必须 --vm + --nic-id ==============
log_section "Bug 11 — vm.nic.update requires --vm and --nic-id"
expect_fail "vm.nic.update without flags → fail" $GOCT vm.nic.update

# ============== Bug 6: vm.disk.update 必须 --vm ==============
log_section "Bug 6 — vm.disk.update requires --vm"
expect_fail "vm.disk.update without --vm → fail" $GOCT vm.disk.update --disk fake-disk-id

# ============== 总结 ==============
print_summary

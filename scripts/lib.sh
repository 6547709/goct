#!/usr/bin/env bash
# scripts/lib.sh — 公共工具函数：日志、计数、命令包装、配置归一化。
#
# 所有 scripts/test_*.sh 都通过 `source "$(dirname "$0")/lib.sh"` 加载本文件。
# 本文件不直接执行，仅暴露函数。
#
# 设计要点：
#   - 不使用 `set -e`：单个命令失败应记录到 FAIL 计数器而不是直接退出（脚本要尽量跑完所有用例）。
#   - run_cmd / expect_fail / expect_exit 是核心断言，每次自动累计 PASS/FAIL/SKIP。
#   - run_section 用于分块输出，方便从大量日志里定位。
#   - GOCT 变量由调用者初始化（指向 ./goct 或预编译路径）。

# 防止重复 source。
if [[ -n "${_GOCT_LIB_LOADED:-}" ]]; then
    return 0
fi
_GOCT_LIB_LOADED=1

# --------- 颜色 ---------
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
DIM='\033[2m'
NC='\033[0m'

# --------- 计数器 ---------
PASS=0
FAIL=0
SKIP=0
FAILED_TESTS=()

# --------- 日志 ---------
log_info()    { echo -e "${BLUE}[INFO]${NC} $*"; }
log_section() { echo -e "\n${CYAN}=== $* ===${NC}"; }
log_pass()    { echo -e "${GREEN}[PASS]${NC} $*"; PASS=$((PASS+1)); }
log_fail()    { echo -e "${RED}[FAIL]${NC} $*"; FAIL=$((FAIL+1)); FAILED_TESTS+=("$*"); }
log_skip()    { echo -e "${YELLOW}[SKIP]${NC} $*"; SKIP=$((SKIP+1)); }
log_warn()    { echo -e "${YELLOW}[WARN]${NC} $*"; }
log_debug()   { [[ "${DEBUG:-false}" == "true" ]] && echo -e "${DIM}[DEBUG]${NC} $*" || true; }

# --------- 命令运行包装 ---------
# run_cmd "描述" command args...
#   命令成功 → PASS；失败 → FAIL（输出错误到 stderr）。
run_cmd() {
    local desc="$1"; shift
    log_debug "+ $*"
    local out
    if out=$("$@" 2>&1); then
        log_pass "$desc"
        log_debug "  $(echo "$out" | head -3)"
        return 0
    fi
    log_fail "$desc"
    echo "  ↳ $out" | head -8
    return 1
}

# expect_fail "描述" command args...
#   命令应当失败（非零退出）；成功时算 FAIL。
expect_fail() {
    local desc="$1"; shift
    log_debug "+ (expect-fail) $*"
    if "$@" >/dev/null 2>&1; then
        log_fail "$desc — expected failure but succeeded"
        return 1
    fi
    log_pass "$desc — failed as expected"
    return 0
}

# expect_exit <expected_code> "描述" command args...
#   命令必须以 expected_code 退出；不符合则 FAIL。
expect_exit() {
    local want="$1"; shift
    local desc="$1"; shift
    log_debug "+ (expect exit=$want) $*"
    "$@" >/dev/null 2>&1
    local got=$?
    if [[ $got -eq $want ]]; then
        log_pass "$desc — exit=$got"
    else
        log_fail "$desc — expected exit=$want, got=$got"
    fi
}

# expect_json "描述" command args...
#   命令应当成功且输出合法 JSON。
expect_json() {
    local desc="$1"; shift
    log_debug "+ (expect-json) $*"
    local out
    if ! out=$("$@" 2>&1); then
        log_fail "$desc — command failed"
        echo "  ↳ $out" | head -5
        return 1
    fi
    if echo "$out" | jq empty >/dev/null 2>&1; then
        log_pass "$desc — valid JSON"
    else
        log_fail "$desc — invalid JSON"
        echo "  ↳ $(echo "$out" | head -3)"
    fi
}

# expect_contains "描述" "要包含的字符串" command args...
expect_contains() {
    local desc="$1"; shift
    local needle="$1"; shift
    log_debug "+ (expect-contains '$needle') $*"
    local out
    if ! out=$("$@" 2>&1); then
        log_fail "$desc — command failed"
        echo "  ↳ $out" | head -3
        return 1
    fi
    if echo "$out" | grep -q -F -- "$needle"; then
        log_pass "$desc"
    else
        log_fail "$desc — output missing $(printf %q "$needle")"
        echo "  ↳ $(echo "$out" | head -3)"
    fi
}

# expect_help_no_error "描述" subcmd
#   验证 `goct <subcmd> --help` 能正常打印（zero exit），不真正调用网络。
expect_help_no_error() {
    local desc="$1"; shift
    log_debug "+ ${GOCT} $* --help"
    if ${GOCT} "$@" --help >/dev/null 2>&1; then
        log_pass "$desc"
    else
        log_fail "$desc — --help should succeed"
    fi
}

# --------- 总结 ---------
print_summary() {
    echo ""
    echo "=========================================="
    echo "  TEST SUMMARY"
    echo "=========================================="
    echo -e "  ${GREEN}PASS:${NC} $PASS"
    echo -e "  ${RED}FAIL:${NC} $FAIL"
    echo -e "  ${YELLOW}SKIP:${NC} $SKIP"
    if (( FAIL > 0 )); then
        echo ""
        echo "Failed tests:"
        for t in "${FAILED_TESTS[@]}"; do
            echo "  - $t"
        done
        return 1
    fi
    return 0
}

# --------- 配置归一化 ---------
# init_goct_env：从环境变量解析连接配置；找到/编译 goct；导出 GOCT。
#
# 行为：
#   1. 必填 GOCT_URL；GOCT_USERNAME / GOCT_PASSWORD 缺失则提示但不强制（如果有 session 缓存就够用）。
#   2. 优先级：./goct → $PATH 中的 goct → 临时编译。
#   3. 把 --insecure 拼到 GOCT 里（如果 GOCT_INSECURE=true）。
init_goct_env() {
    : "${GOCT_URL:=}"
    : "${GOCT_USERNAME:=}"
    : "${GOCT_PASSWORD:=}"
    : "${GOCT_SOURCE:=local}"
    : "${GOCT_INSECURE:=false}"

    if [[ -z "$GOCT_URL" ]]; then
        echo -e "${RED}ERROR${NC}: GOCT_URL not set" >&2
        echo "Please export GOCT_URL / GOCT_USERNAME / GOCT_PASSWORD before running."
        exit 2
    fi
    export GOCT_URL GOCT_USERNAME GOCT_PASSWORD GOCT_SOURCE GOCT_INSECURE

    # locate / build binary
    local repo_root
    repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
    if [[ -x "$repo_root/goct" ]]; then
        GOCT_BIN="$repo_root/goct"
    elif command -v goct >/dev/null 2>&1; then
        GOCT_BIN="$(command -v goct)"
    else
        log_info "Building goct from source ..."
        ( cd "$repo_root" && go build -o goct . ) || {
            echo -e "${RED}ERROR${NC}: failed to build goct" >&2
            exit 2
        }
        GOCT_BIN="$repo_root/goct"
    fi

    # Construct GOCT command prefix.
    local insecure_flag=""
    if [[ "$GOCT_INSECURE" == "true" || "$GOCT_INSECURE" == "1" ]]; then
        insecure_flag="--insecure"
    fi
    GOCT="$GOCT_BIN $insecure_flag"
    export GOCT GOCT_BIN

    log_info "Binary : $GOCT_BIN"
    log_info "URL    : $GOCT_URL"
    log_info "User   : $GOCT_USERNAME"
    log_info "Source : $GOCT_SOURCE"
    log_info "TLS    : $([[ "$GOCT_INSECURE" == "true" ]] && echo skip || echo verify)"
}

# wait_task <task-id> [timeout-sec]
#   轮询 task.ls 直到 SUCCESSED / FAILED / 超时。
wait_task() {
    local task_id="$1"
    local timeout_sec="${2:-300}"
    local elapsed=0
    while (( elapsed < timeout_sec )); do
        local row
        row=$(${GOCT} task.ls --format json 2>/dev/null | jq -r --arg id "$task_id" \
            '.[] | select(.ID == $id) | .Status' 2>/dev/null || true)
        case "$row" in
            SUCCESSED|SUCCEEDED) return 0 ;;
            FAILED|ERROR)        return 1 ;;
        esac
        sleep 3
        elapsed=$((elapsed + 3))
    done
    log_warn "wait_task: timeout after ${timeout_sec}s (task=$task_id)"
    return 2
}

# extract_task_id <stdout>
#   从命令输出里抓取最后一个 task-* / cl* 形式的 ID。
extract_task_id() {
    grep -oE '(task-[a-z0-9]+|cl[0-9a-z]{25})' <<< "$1" | tail -1
}

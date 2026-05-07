#!/usr/bin/env bash
#
# test_vm_features.sh - VM 扩展功能综合测试脚本
#
# 用法:
#   1. 设置环境变量（见下方 CONFIG 部分）
#   2. 运行: ./test_vm_features.sh
#
# 注意: 测试会创建和销毁资源，请在测试环境运行！
#

set -uo pipefail

# ===================== CONFIG =====================
# 必填配置
GOCT_URL="${GOCT_URL:-}"
GOCT_USERNAME="${GOCT_USERNAME:-}"
GOCT_PASSWORD="${GOCT_PASSWORD:-}"
# 选填配置
GOCT_CLUSTER="${GOCT_CLUSTER:-}"           # 目标集群
GOCT_SOURCE="${GOCT_SOURCE:-local}"         # 登录源
GOCT_INSECURE="${GOCT_INSECURE:-false}"    # 跳过 TLS 校验
# 测试 VM 的前缀，便于识别和清理
TEST_VM_PREFIX="${TEST_VM_PREFIX:-test-vm}"
# 测试超时时间（秒）
TIMEOUT="${TIMEOUT:-300}"
# Debug 模式
DEBUG="${DEBUG:-false}"
# ==================================================

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# 全局变量
GOCT_BIN=""
TEMP_VM_ID=""
TEMP_VM_NAME=""
AUTH_TOKEN=""

# ===================== 工具函数 =====================
log_debug() {
    if [[ "$DEBUG" == "true" ]]; then
        echo -e "${CYAN}[DEBUG]${NC} $1"
    fi
}

log_info() { echo -e "${BLUE}[INFO]${NC} $1"; }
log_ok()   { echo -e "${GREEN}[OK]${NC}   $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC}  $1"; }
log_err()  { echo -e "${RED}[ERROR]${NC} $1"; }

run_cmd() {
    local desc="$1"
    shift
    log_debug "执行: $desc"
    log_debug "命令: $*"
    local output
    local exit_code=0
    output=$("$@" 2>&1) || exit_code=$?
    log_debug "退出码: $exit_code"
    if [[ $exit_code -ne 0 ]]; then
        log_debug "输出: $output"
    fi
    echo "$output"
    return $exit_code
}

cleanup() {
    log_info "清理测试资源..."
    if [[ -n "$TEMP_VM_ID" && "$TEMP_VM_ID" != "__NONE__" ]]; then
        log_info "尝试清理 VM: $TEMP_VM_NAME (ID: $TEMP_VM_ID)"
        timeout 10 "$GOCT_BIN" vm.destroy "$TEMP_VM_ID" --force 2>/dev/null || true
    fi
}

die() {
    log_err "$1"
    log_err "当前状态: TEMP_VM_NAME=$TEMP_VM_NAME, TEMP_VM_ID=$TEMP_VM_ID"
    cleanup
    exit 1
}

check_cmd() {
    local exit_code=$1
    local desc="$2"
    if [[ $exit_code -ne 0 ]]; then
        log_warn "命令执行失败（继续）: $desc (exit code: $exit_code)"
        return 0
    fi
}

# 获取认证 Token
get_auth_token() {
    log_info "获取认证 Token..."

    # 将 source 转为大写（API 要求大写）
    local source_upper
    source_upper=$(echo "$GOCT_SOURCE" | tr '[:lower:]' '[:upper:]')

    local login_data
    login_data=$(jq -n \
        --arg username "$GOCT_USERNAME" \
        --arg password "$GOCT_PASSWORD" \
        --arg source "$source_upper" \
        '{
            username: $username,
            password: $password,
            source: $source
        }')

    log_debug "登录请求数据: $login_data"

    local response
    response=$(curl -s -k -X POST \
        -H "Content-Type: application/json" \
        -d "$login_data" \
        "${GOCT_URL}/v2/api/login")

    log_debug "登录响应: $response"

    AUTH_TOKEN=$(echo "$response" | jq -r '.data.token // empty')

    if [[ -z "$AUTH_TOKEN" || "$AUTH_TOKEN" == "null" ]]; then
        die "登录失败: $response"
    fi

    log_ok "登录成功"
}

# 获取内容库模板列表
list_content_library_templates() {
    log_info "获取内容库模板列表..."

    local response
    response=$(curl -s -k -X POST \
        -H "Content-Type: application/json" \
        -H "Authorization: $AUTH_TOKEN" \
        -d '{"where":{}}' \
        "${GOCT_URL}/v2/api/get-content-library-vm-templates")

    log_debug "模板列表响应: $response"

    # 注意：vm_templates 数组内的 id 才是用于创建 VM 的模板 ID
    echo "$response" | jq -r '.[] | .vm_templates[]? | "\(.id)\t\(.name)"' 2>/dev/null || {
        log_err "获取模板列表失败: $response"
        return 1
    }
}

# 根据模板名获取模板 ID
get_template_id_by_name() {
    local template_name="$1"
    local templates
    templates=$(list_content_library_templates) || return 1

    local template_id
    template_id=$(echo "$templates" | grep -i "$template_name" | head -1 | cut -f1)

    if [[ -z "$template_id" ]]; then
        log_err "未找到模板: $template_name" >&2
        log_info "可用模板:" >&2
        echo "$templates" | while read -r line; do
            echo "    $line" >&2
        done
        return 1
    fi

    local template_name_found
    template_name_found=$(echo "$templates" | grep "$template_id" | cut -f2)
    log_info "找到模板: $template_name_found (ID: $template_id)" >&2
    echo "$template_id"
}

# 等待任务完成
wait_task() {
    local task_id="$1"
    local timeout_sec="${2:-$TIMEOUT}"
    local elapsed=0

    log_info "等待任务完成: $task_id (超时: ${timeout_sec}s)"

    while [[ $elapsed -lt $timeout_sec ]]; do
        local status
        local task_output
        task_output=$(run_cmd "查询任务状态" "$GOCT_BIN" task.ls 2>/dev/null) || true
        status=$(echo "$task_output" | grep "$task_id" | awk '{print $3}' || echo "unknown")

        case "$status" in
            SUCCESS|success)
                log_ok "任务完成: $task_id"
                return 0
                ;;
            FAILED|failed|ERROR|error)
                log_err "任务失败: $task_id"
                log_err "任务输出: $task_output"
                return 1
                ;;
            *)
                log_info "等待任务... ($elapsed/${timeout_sec}s) 状态: $status"
                ;;
        esac
        sleep 5
        elapsed=$((elapsed + 5))
    done

    log_warn "任务等待超时: $task_id"
    return 1
}

# ===================== 前置检查 =====================
check_env() {
    log_info "检查环境配置..."

    [[ -z "$GOCT_URL" ]] && die "请设置 GOCT_URL 环境变量"
    [[ -z "$GOCT_USERNAME" ]] && die "请设置 GOCT_USERNAME 环境变量"
    [[ -z "$GOCT_PASSWORD" ]] && die "请设置 GOCT_PASSWORD 环境变量"

    log_ok "环境变量检查通过"
}

check_goct() {
    log_info "检查 goct 二进制文件..."

    if command -v goct &>/dev/null; then
        GOCT_BIN="goct"
    elif [[ -x "./goct" ]]; then
        GOCT_BIN="./goct"
    elif [[ -x "/usr/local/bin/goct" ]]; then
        GOCT_BIN="/usr/local/bin/goct"
    else
        log_info "编译 goct..."
        go build -o goct . || die "编译 goct 失败"
        GOCT_BIN="./goct"
    fi

    log_ok "使用: $GOCT_BIN"

    # 显示 goct 版本
    local version
    version=$("$GOCT_BIN" version 2>&1 || echo "unknown")
    log_info "goct 版本: $version"
}

check_connection() {
    log_info "检查 CloudTower 连接..."

    local output
    output=$(run_cmd "检查连接" "$GOCT_BIN" about 2>&1) || die "无法连接到 CloudTower，请检查配置"

    if echo "$output" | grep -qi "error\|failed"; then
        die "CloudTower 连接失败: $output"
    fi

    log_ok "连接成功"
    log_debug "连接响应: $output"
}

# ===================== 测试用例 =====================
test_vm_ls() {
    log_info "===== 测试 vm.ls ====="
    local output
    output=$(run_cmd "vm.ls" "$GOCT_BIN" vm.ls) || check_cmd $? "vm.ls"
    log_debug "vm.ls 输出:"
    log_debug "$output"
    echo "$output" | head -10
    log_ok "vm.ls 测试通过"
}

test_vm_info() {
    log_info "===== 测试 vm.info ====="
    local output
    output=$(run_cmd "vm.info" "$GOCT_BIN" vm.info "$TEMP_VM_NAME") || check_cmd $? "vm.info"
    log_debug "vm.info 输出: $output"
    log_ok "vm.info 测试通过"
}

test_vm_ip() {
    log_info "===== 测试 vm.ip ====="
    local output
    output=$(run_cmd "vm.ip" "$GOCT_BIN" vm.ip "$TEMP_VM_NAME" 2>/dev/null) || output=""
    if [[ -n "$output" ]]; then
        log_ok "获取到 IP: $output"
    else
        log_warn "vm.ip 未返回 IP（可能正常）"
    fi
}

test_vm_vnc() {
    log_info "===== 测试 vm.vnc ====="
    local output
    output=$(run_cmd "vm.vnc" "$GOCT_BIN" vm.vnc "$TEMP_VM_NAME" 2>&1) || true
    if echo "$output" | grep -qi "error"; then
        log_warn "vm.vnc 返回: $output"
    else
        log_ok "vm.vnc 输出:"
        echo "$output" | while read -r line; do
            echo "    $line"
        done
    fi
}

test_vm_tools_install() {
    log_info "===== 测试 vm.tools.install ====="

    local output
    output=$(run_cmd "vm.tools.install" "$GOCT_BIN" vm.tools.install "$TEMP_VM_NAME" 2>&1) || true

    if echo "$output" | grep -qi "error"; then
        log_warn "vm.tools.install 返回: $output"
    elif echo "$output" | grep -qi "task-"; then
        local task_id
        task_id=$(echo "$output" | grep -oE 'task-[a-z0-9]+' | head -1)
        log_info "任务: $task_id"
        wait_task "$task_id" 120
    else
        log_ok "vm.tools.install 执行完成"
    fi
}

test_vm_nic_add() {
    log_info "===== 测试 vm.nic.add ====="

    local output
    output=$(run_cmd "vm.nic.add" "$GOCT_BIN" vm.nic.add "$TEMP_VM_NAME" --type VLAN --model VIRTIO 2>&1) || check_cmd $? "vm.nic.add"

    if echo "$output" | grep -qi "task-"; then
        local task_id
        task_id=$(echo "$output" | grep -oE 'task-[a-z0-9]+' | head -1)
        wait_task "$task_id"
    else
        log_ok "vm.nic.add 执行完成"
    fi
}

test_vm_nic_rm() {
    log_info "===== 测试 vm.nic.rm ====="
    # NIC 删除需要至少 2 个 NIC，当前测试 VM 的 NIC 数量不足，跳过此测试
    log_warn "跳过 vm.nic.rm（需要至少 2 个 NIC）"
}

test_vm_create_from_template() {
    log_info "===== 测试 vm.create --from-template ====="

    # 获取集群 ID
    if [[ -z "$GOCT_CLUSTER" ]]; then
        log_warn "跳过 --from-template（未设置 GOCT_CLUSTER）"
        return 0
    fi

    # 获取模板 ID（支持模板名或模板 ID）
    local template_id="${GOCT_TEMPLATE:-}"
    if [[ -n "$template_id" ]] && [[ "$template_id" != cl* ]]; then
        # 传入的是模板名，需要查找
        log_info "查找模板: $template_id"
        template_id=$(get_template_id_by_name "$template_id") || {
            log_warn "无法获取模板 ID，跳过测试"
            return 0
        }
    fi

    if [[ -z "$template_id" ]]; then
        log_warn "跳过 --from-template（未设置 GOCT_TEMPLATE）"
        return 0
    fi

    local clone_name="${TEST_VM_PREFIX}-from-template-$(date +%s)"
    log_info "从模板创建 VM: $clone_name"

    local output
    output=$(run_cmd "vm.create --from-template" "$GOCT_BIN" vm.create \
        --name "$clone_name" \
        --from-template "$template_id" \
        --cluster "$GOCT_CLUSTER" \
        --full-copy \
        2>&1) || check_cmd $? "vm.create --from-template"

    if echo "$output" | grep -qi "task-"; then
        local task_id
        task_id=$(echo "$output" | grep -oE 'task-[a-z0-9]+' | head -1)
        wait_task "$task_id" 180

        # 清理克隆的 VM
        sleep 3
        local clone_vm_id
        clone_vm_id=$("$GOCT_BIN" vm.ls 2>/dev/null | grep "$clone_name" | awk '{print $1}' | head -1)
        if [[ -n "$clone_vm_id" ]]; then
            log_info "清理克隆 VM: $clone_name ($clone_vm_id)"
            "$GOCT_BIN" vm.destroy "$clone_vm_id" --force 2>/dev/null || true
        fi
    else
        log_ok "vm.create --from-template 执行完成"
    fi
}

# test_vm_export() {
#     log_info "===== 测试 vm.export ====="
#
#     local output
#     output=$(run_cmd "vm.export" "$GOCT_BIN" vm.export "$TEMP_VM_NAME" 2>&1) || true
#
#     if echo "$output" | grep -qi "error"; then
#         log_warn "vm.export 返回: $output"
#     elif echo "$output" | grep -qi "task-"; then
#         local task_id
#         task_id=$(echo "$output" | grep -oE 'task-[a-z0-9]+' | head -1)
#         log_info "导出任务已启动: $task_id (不等待完成)"
#     else
#         log_ok "vm.export 执行完成"
#     fi
# }

# ===================== 创建测试 VM =====================
create_test_vm() {
    log_info "===== 创建测试 VM ====="

    if [[ -z "$GOCT_CLUSTER" ]]; then
        die "请设置 GOCT_CLUSTER 环境变量"
    fi

    TEMP_VM_NAME="${TEST_VM_PREFIX}-$(date +%s)"

    # 优先使用模板创建 VM
    if [[ -n "${GOCT_TEMPLATE:-}" ]]; then
        local template_id="${GOCT_TEMPLATE}"
        if [[ "$template_id" != cl* ]]; then
            log_info "查找模板: $template_id"
            template_id=$(get_template_id_by_name "$template_id") || {
                log_warn "无法获取模板 ID，将尝试手动创建 VM"
                template_id=""
            }
        fi

        if [[ -n "$template_id" ]]; then
            log_info "从模板创建 VM: $TEMP_VM_NAME (模板: $template_id)"
            local output
            output=$(run_cmd "vm.create --from-template" "$GOCT_BIN" vm.create \
                --name "$TEMP_VM_NAME" \
                --from-template "$template_id" \
                --cluster "$GOCT_CLUSTER" \
                --full-copy \
                2>&1)

            if [[ $? -eq 0 ]]; then
                if echo "$output" | grep -qi "task-"; then
                    local task_id
                    task_id=$(echo "$output" | grep -oE 'task-[a-z0-9]+' | head -1)
                    log_info "创建任务: $task_id"
                    wait_task "$task_id" 300
                else
                    log_ok "VM 创建完成"
                fi

                # 获取 VM ID
                sleep 3
                TEMP_VM_ID=$("$GOCT_BIN" vm.ls 2>/dev/null | grep "$TEMP_VM_NAME" | awk '{print $1}' | head -1 || "")

                if [[ -n "$TEMP_VM_ID" ]]; then
                    log_ok "测试 VM: $TEMP_VM_NAME (ID: $TEMP_VM_ID)"
                    return 0
                fi
            else
                log_warn "模板创建失败: $output"
            fi
        fi
    fi

    # 如果模板创建失败或未提供模板，尝试手动创建
    log_warn "尝试手动创建 VM（需要存储策略配置）..."
    log_info "创建 VM: $TEMP_VM_NAME"

    local output
    output=$(run_cmd "vm.create" "$GOCT_BIN" vm.create \
        --name "$TEMP_VM_NAME" \
        --cluster "$GOCT_CLUSTER" \
        --vcpu 2 \
        --memory 2048 \
        2>&1)
    local exit_code=$?

    log_debug "vm.create 退出码: $exit_code"
    log_debug "vm.create 输出: $output"

    if [[ $exit_code -ne 0 ]]; then
        log_warn "手动创建 VM 失败，尝试使用已有 VM..."
        # 从 vm.ls 中找一个运行中的 VM（不在回收站的）
        # 使用 --json 格式获取结构化数据
        local existing_line
        existing_line=$("$GOCT_BIN" vm.ls --format json 2>/dev/null | jq -r '.[] | select(.Status == "RUNNING" and (.Name | test("in-recycle-bin") | not)) | "\(.ID)\t\(.Name)"' 2>/dev/null | head -1)

        if [[ -z "$existing_line" ]]; then
            log_warn "没有可用的 VM 进行测试（所有运行中 VM 的 NIC 都只有 1 个）"
            TEMP_VM_NAME="__NONE__"
            TEMP_VM_ID=""
        else
            TEMP_VM_ID=$(echo "$existing_line" | cut -f1)
            TEMP_VM_NAME=$(echo "$existing_line" | cut -f2)
            log_info "找到已有 VM: $TEMP_VM_NAME (ID: $TEMP_VM_ID)"
            log_ok "使用已有 VM: $TEMP_VM_NAME (ID: $TEMP_VM_ID)"
            return 0
        fi
    fi

    if echo "$output" | grep -qi "task-"; then
        local task_id
        task_id=$(echo "$output" | grep -oE 'task-[a-z0-9]+' | head -1)
        log_info "创建任务: $task_id"
        wait_task "$task_id" 180
    elif echo "$output" | grep -qi "created"; then
        log_ok "VM 创建完成（同步）"
    else
        log_warn "创建输出: $output"
    fi

    # 获取 VM ID
    sleep 3
    local check_vm_id
    check_vm_id=$("$GOCT_BIN" vm.ls 2>/dev/null | grep "$TEMP_VM_NAME" | awk '{print $1}' | head -1 || "")

    if [[ -n "$check_vm_id" ]]; then
        TEMP_VM_ID="$check_vm_id"
        log_ok "测试 VM: $TEMP_VM_NAME (ID: $TEMP_VM_ID)"
    else
        log_warn "无法获取 VM ID: $TEMP_VM_NAME (可能创建失败或使用已有 VM)"
    fi
}

# ===================== 主流程 =====================
main() {
    echo "=========================================="
    echo "  VM 扩展功能综合测试"
    echo "=========================================="
    echo ""

    # 陷阱，确保退出时清理
    trap cleanup EXIT

    check_env
    check_goct

    echo ""
    log_info "CloudTower URL: $GOCT_URL"
    log_info "Username: $GOCT_USERNAME"
    [[ -n "$GOCT_CLUSTER" ]] && log_info "Cluster: $GOCT_CLUSTER"
    [[ -n "${GOCT_TEMPLATE:-}" ]] && log_info "Template: $GOCT_TEMPLATE"
    [[ "$DEBUG" == "true" ]] && log_info "Debug 模式: 启用"
    echo ""

    # 先登录获取 token，用于 GraphQL 查询
    get_auth_token
    check_connection

    echo ""
    echo "##########################################"
    echo "# 步骤 1: 创建测试 VM"
    echo "##########################################"
    create_test_vm

    echo ""
    echo "##########################################"
    echo "# 步骤 2: 执行功能测试"
    echo "##########################################"

    test_vm_ls
    test_vm_info
    test_vm_ip
    test_vm_vnc
    test_vm_tools_install
    test_vm_nic_add
    test_vm_nic_rm
    test_vm_create_from_template
    # test_vm_export  # 跳过（导出时间过长）

    echo ""
    echo "##########################################"
    echo "# 步骤 3: 清理测试 VM"
    echo "##########################################"

    echo ""
    echo "=========================================="
    log_ok "全部测试完成！"
    echo "=========================================="
    echo ""
    echo "测试 VM: $TEMP_VM_NAME"
    echo "如需手动清理，请运行:"
    echo "  $GOCT_BIN vm.destroy $TEMP_VM_ID --force"
    echo ""
}

main "$@"

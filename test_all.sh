#!/bin/bash
# goct comprehensive test script
# Usage: ./test_goct.sh [vm-name] [host-name] [cluster-name]

set -o pipefail

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Config - can be overridden by env vars
GOCT_URL="${GOCT_URL:-https://10.0.50.210:443}"
GOCT_USERNAME="${GOCT_USERNAME:-root}"
GOCT_PASSWORD="${GOCT_PASSWORD:-VMware1!}"
GOCT_SOURCE="${GOCT_SOURCE:-local}"
GOCT_INSECURE="${GOCT_INSECURE:-true}"

# Test targets - use arguments or defaults
TEST_VM="${1:-demo01}"
TEST_HOST="${2:-}"
TEST_CLUSTER="${3:-}"

# Counters
PASS=0
FAIL=0
SKIP=0

# Functions
log_info() { echo -e "${BLUE}[INFO]${NC} $1"; }
log_pass() { echo -e "${GREEN}[PASS]${NC} $1"; ((PASS++)); }
log_fail() { echo -e "${RED}[FAIL]${NC} $1"; ((FAIL++)); }
log_skip() { echo -e "${YELLOW}[SKIP]${NC} $1"; ((SKIP++)); }

run_cmd() {
    local description="$1"
    local cmd="$2"
    local expected_fail="${3:-false}"

    log_info "$description"
    log_info "  Command: $cmd"

    if output=$($cmd 2>&1); then
        if [ "$expected_fail" = "true" ]; then
            log_fail "$description - expected failure but succeeded"
            echo "  Output: $output"
        else
            log_pass "$description"
            if [ -n "$output" ]; then
                echo "  Output: $output" | head -5
            fi
        fi
    else
        if [ "$expected_fail" = "true" ]; then
            log_pass "$description - expected failure"
        else
            log_fail "$description"
            echo "  Error: $output"
        fi
    fi
}

run_cmd_json() {
    local description="$1"
    local cmd="$2"

    log_info "$description"
    log_info "  Command: $cmd"

    if output=$($cmd 2>&1); then
        if echo "$output" | jq -e . >/dev/null 2>&1; then
            log_pass "$description - valid JSON"
        else
            log_fail "$description - invalid JSON"
            echo "  Output: $output"
        fi
    else
        log_fail "$description"
        echo "  Error: $output"
    fi
}

check_field() {
    local description="$1"
    local cmd="$2"
    local field="$3"
    local expected="$4"

    log_info "$description"
    log_info "  Command: $cmd"

    if output=$($cmd 2>&1); then
        actual=$(echo "$output" | jq -r ".$field" 2>/dev/null || echo "$output" | grep "$field" | awk '{print $2}')
        if [ "$actual" = "$expected" ]; then
            log_pass "$description (got: $actual)"
        else
            log_fail "$description (expected: $expected, got: $actual)"
        fi
    else
        log_fail "$description - command failed"
        echo "  Error: $output"
    fi
}

# Export for subcommands
export GOCT_URL GOCT_USERNAME GOCT_PASSWORD GOCT_SOURCE GOCT_INSECURE

# Build goct command with common flags
GOCT="./goct --insecure"
export GOCT_URL GOCT_USERNAME GOCT_PASSWORD GOCT_SOURCE GOCT_INSECURE

echo "=========================================="
echo "  goct Comprehensive Test Suite"
echo "=========================================="
echo ""
echo "Config:"
echo "  URL: $GOCT_URL"
echo "  Username: $GOCT_USERNAME"
echo "  Source: $GOCT_SOURCE"
echo "  Insecure: $GOCT_INSECURE"
echo "  Test VM: $TEST_VM"
echo ""

#==========================================
# System Commands (no login required)
#==========================================
echo ""
echo "=========================================="
echo "  SYSTEM COMMANDS (no login required)"
echo "=========================================="

run_cmd "Version command" "$GOCT version" false
run_cmd "Session list" "$GOCT session.ls" false

#==========================================
# Session & Connection
#==========================================
echo ""
echo "=========================================="
echo "  SESSION & CONNECTION"
echo "=========================================="

run_cmd "Session login" "$GOCT session.login" false
run_cmd "About (server info)" "$GOCT about" false

#==========================================
# VM Commands
#==========================================
echo ""
echo "=========================================="
echo "  VM COMMANDS"
echo "=========================================="

run_cmd "VM list" "$GOCT vm.ls" false
run_cmd "VM list (ID only)" "$GOCT vm.ls --id-only" false
run_cmd "VM list (JSON)" "$GOCT vm.ls --format json" false
run_cmd_json "VM list JSON validation" "$GOCT vm.ls --format json"

if [ -n "$TEST_VM" ]; then
    run_cmd "VM info: $TEST_VM" "$GOCT vm.info $TEST_VM" false
    run_cmd "VM info JSON: $TEST_VM" "$GOCT vm.info $TEST_VM --format json" false
    run_cmd "VM IP: $TEST_VM" "$GOCT vm.ip $TEST_VM" false
    run_cmd "VM disk list: $TEST_VM" "$GOCT disk.ls $TEST_VM" false
    run_cmd "VM NIC list: $TEST_VM" "$GOCT nic.ls $TEST_VM" false
    run_cmd "VM snapshot list: $TEST_VM" "$GOCT vm.snapshot.ls $TEST_VM" false
    run_cmd "VM VNC: $TEST_VM" "$GOCT vm.vnc $TEST_VM" false

    # Check OS field in VM info
    log_info "Checking OS field in vm.info for $TEST_VM"
    os_output=$($GOCT vm.info $TEST_VM 2>&1)
    if echo "$os_output" | grep -q "OS:"; then
        os_value=$(echo "$os_output" | grep "OS:" | awk '{print $2}')
        if [ -n "$os_value" ] && [ "$os_value" != "-" ]; then
            log_pass "VM OS field populated: $os_value"
        else
            log_fail "VM OS field is empty or '-'"
        fi
    else
        log_skip "VM OS field not found in output"
    fi

    # VM Power State
    vm_status=$($GOCT vm.ls 2>/dev/null | grep "$TEST_VM" | awk '{print $3}')

    # Test metrics (requires running VM)
    if [ "$vm_status" = "RUNNING" ]; then
        run_cmd "VM metrics list" "$GOCT vm.metrics --list" false
        run_cmd "VM metrics (CPU): $TEST_VM" "$GOCT vm.metrics elf_vm_cpu_overall_usage_percent $TEST_VM" false
        run_cmd "VM metrics (memory): $TEST_VM" "$GOCT vm.metrics elf_vm_memory_usage_percent $TEST_VM" false
        run_cmd "VM volume metrics list" "$GOCT vm.volume --list" false
    fi
fi

#==========================================
# Host Commands
#==========================================
echo ""
echo "=========================================="
echo "  HOST COMMANDS"
echo "=========================================="

run_cmd "Host list" "$GOCT host.ls" false
run_cmd "Host list (ID only)" "$GOCT host.ls --id-only" false

if [ -n "$TEST_HOST" ]; then
    run_cmd "Host info: $TEST_HOST" "$GOCT host.info $TEST_HOST" false
fi

# Get first host for testing (NAME is column 3 in table)
FIRST_HOST=$($GOCT host.ls 2>/dev/null | tail -n +4 | head -1 | awk -F'│' '{gsub(/ /, "", $3); print $3}')
if [ -n "$FIRST_HOST" ]; then
    run_cmd "Host info (first host) by name" "$GOCT host.info $FIRST_HOST" false
fi

# Test host.info with ID (verify ID auto-detection works)
FIRST_HOST_ID=$($GOCT host.ls 2>/dev/null | tail -n +4 | head -1 | awk -F'│' '{gsub(/ /, "", $2); print $2}')
if [ -n "$FIRST_HOST_ID" ]; then
    run_cmd "Host info (first host) by ID" "$GOCT host.info $FIRST_HOST_ID" false
fi

#==========================================
# Cluster Commands
#==========================================
echo ""
echo "=========================================="
echo "  CLUSTER COMMANDS"
echo "=========================================="

run_cmd "Cluster list" "$GOCT cluster.ls" false
run_cmd "Cluster list (ID only)" "$GOCT cluster.ls --id-only" false

if [ -n "$TEST_CLUSTER" ]; then
    run_cmd "Cluster info: $TEST_CLUSTER" "$GOCT cluster.info $TEST_CLUSTER" false
fi

# Get first cluster for testing (NAME is column 3 in table)
FIRST_CLUSTER=$($GOCT cluster.ls 2>/dev/null | tail -n +4 | head -1 | awk -F'│' '{gsub(/ /, "", $3); print $3}')
if [ -n "$FIRST_CLUSTER" ]; then
    run_cmd "Cluster info (first cluster) by name" "$GOCT cluster.info $FIRST_CLUSTER" false
fi

# Test cluster.info with ID (verify ID auto-detection works)
FIRST_CLUSTER_ID=$($GOCT cluster.ls 2>/dev/null | tail -n +4 | head -1 | awk -F'│' '{gsub(/ /, "", $2); print $2}')
if [ -n "$FIRST_CLUSTER_ID" ]; then
    run_cmd "Cluster info (first cluster) by ID" "$GOCT cluster.info $FIRST_CLUSTER_ID" false
fi

#==========================================
# Datastore Commands
#==========================================
echo ""
echo "=========================================="
echo "  DATASTORE COMMANDS"
echo "=========================================="

run_cmd "Datastore list" "$GOCT datastore.ls" false
run_cmd "Datastore disk list" "$GOCT datastore.disk.ls" false

#==========================================
# Network Commands
#==========================================
echo ""
echo "=========================================="
echo "  NETWORK COMMANDS"
echo "=========================================="

run_cmd "Network list" "$GOCT network.ls" false
run_cmd "VLAN list" "$GOCT vlan.ls" false

#==========================================
# Task Commands
#==========================================
echo ""
echo "=========================================="
echo "  TASK COMMANDS"
echo "=========================================="

run_cmd "Task list" "$GOCT task.ls" false
run_cmd "Task list (JSON)" "$GOCT task.ls --format json" false

#==========================================
# Alert Commands
#==========================================
echo ""
echo "=========================================="
echo "  ALERT COMMANDS"
echo "=========================================="

run_cmd "Alert list" "$GOCT alert.ls" false
run_cmd "Alert list (JSON)" "$GOCT alert.ls --format json" false

#==========================================
# User Commands
#==========================================
echo ""
echo "=========================================="
echo "  USER COMMANDS"
echo "=========================================="

run_cmd "User list" "$GOCT user.ls" false
run_cmd "User list (JSON)" "$GOCT user.ls --format json" false

#==========================================
# Volume Metrics
#==========================================
echo ""
echo "=========================================="
echo "  VOLUME METRICS"
echo "=========================================="

run_cmd "Volume metrics list" "$GOCT volume.metrics --list" false

#==========================================
# Summary
#==========================================
echo ""
echo "=========================================="
echo "  TEST SUMMARY"
echo "=========================================="
echo -e "${GREEN}PASS:${NC} $PASS"
echo -e "${RED}FAIL:${NC} $FAIL"
echo -e "${YELLOW}SKIP:${NC} $SKIP"
echo ""

if [ $FAIL -eq 0 ]; then
    echo -e "${GREEN}All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}Some tests failed!${NC}"
    exit 1
fi
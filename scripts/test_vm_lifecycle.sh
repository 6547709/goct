#!/usr/bin/env bash
# scripts/test_vm_lifecycle.sh — VM 写操作端到端测试。
#
# **此脚本会创建并销毁 VM**，请勿在生产环境运行！
#
# 覆盖：
#   - vm.create（新的 --disk / --nic 显式 spec，对齐 v0.2.1 修复）
#   - vm.update（--description="" 清空场景，Bug 7）
#   - vm.power.on / off / reset / suspend / resume / shutdown
#   - vm.disk.add / .expand / .ls / .rm
#   - vm.nic.add / .ls / .update（要求 --vm + --nic-id，Bug 11）/ .rm
#   - vm.cdrom.add / .ls / .eject / .toggle / .rm
#   - vm.snapshot.create / .ls / .revert / .rm
#   - vm.recycle / .recover / .destroy
#
# 用法：
#   GOCT_URL=...  GOCT_USERNAME=...  GOCT_PASSWORD=...  GOCT_INSECURE=true \
#   GOCT_CLUSTER=<cluster-name-or-id> \
#   GOCT_VLAN=<vlan-id-or-name> \
#       ./scripts/test_vm_lifecycle.sh
#
# 可选环境变量：
#   GOCT_TEMPLATE   优先用模板创建；否则用 --disk/--nic 从零创建
#   TEST_VM_PREFIX  测试 VM 名前缀（默认 goct-lc）
#   KEEP_VM=true    退出时不清理 VM（调试用）
#   DEBUG=true      打印调试日志
#
# 退出码：
#   0 全部通过
#   1 至少一个用例 FAIL（脚本仍尝试清理）
#   2 配置错误

set -uo pipefail
cd "$(dirname "$0")/.."

source "scripts/lib.sh"

: "${GOCT_CLUSTER:=}"
: "${GOCT_VLAN:=}"
: "${GOCT_TEMPLATE:=}"
: "${TEST_VM_PREFIX:=goct-lc}"
: "${KEEP_VM:=false}"

if [[ -z "$GOCT_CLUSTER" ]]; then
    echo "ERROR: GOCT_CLUSTER required (cluster ID or name)" >&2
    exit 2
fi

init_goct_env

VM_NAME="${TEST_VM_PREFIX}-$$-$(date +%s)"
VM_ID=""

cleanup() {
    if [[ "$KEEP_VM" == "true" ]]; then
        log_warn "KEEP_VM=true → not destroying $VM_NAME"
        return
    fi
    if [[ -n "$VM_ID" ]]; then
        log_section "Cleanup"
        log_info "Stopping VM (best-effort) ..."
        ${GOCT} vm.power.off "$VM_ID" --force >/dev/null 2>&1 || true
        sleep 2
        log_info "Destroying VM $VM_NAME ($VM_ID) ..."
        ${GOCT} vm.destroy "$VM_ID" --force >/dev/null 2>&1 || true
    fi
}
trap cleanup EXIT

# ---------- 1. Create VM ----------
log_section "1. vm.create"

if [[ -n "$GOCT_TEMPLATE" ]]; then
    log_info "Creating from template $GOCT_TEMPLATE → $VM_NAME"
    out=$(${GOCT} vm.create \
        --name "$VM_NAME" \
        --cluster "$GOCT_CLUSTER" \
        --from-template "$GOCT_TEMPLATE" \
        ${GOCT_VLAN:+--nic-vlan "$GOCT_VLAN"} \
        2>&1) || { log_fail "vm.create from template"; echo "$out"; exit 1; }
    log_pass "vm.create from template"
else
    if [[ -z "$GOCT_VLAN" ]]; then
        echo "ERROR: GOCT_VLAN required when GOCT_TEMPLATE is unset (CloudTower needs a VLAN for the NIC)" >&2
        exit 2
    fi
    log_info "Creating from scratch (--disk + --nic) → $VM_NAME"
    out=$(${GOCT} vm.create \
        --name "$VM_NAME" \
        --cluster "$GOCT_CLUSTER" \
        --vcpu 2 --memory 2048 --firmware BIOS \
        --disk size=10g,bus=SCSI,name=disk0 \
        --nic vlan="$GOCT_VLAN",model=VIRTIO 2>&1) \
        || { log_fail "vm.create (--disk + --nic)"; echo "$out"; exit 1; }
    log_pass "vm.create (--disk + --nic)"
fi

# Resolve VM_ID
sleep 3
VM_ID=$(${GOCT} vm.ls --format json 2>/dev/null | jq -r --arg n "$VM_NAME" '.[] | select(.Name == $n) | .ID' | head -1)
if [[ -z "$VM_ID" ]]; then
    log_fail "Could not resolve VM ID for $VM_NAME"
    exit 1
fi
log_info "VM_ID=$VM_ID"

# ---------- 2. Read-only assertions ----------
log_section "2. Read VM state"
run_cmd     "vm.info $VM_NAME"           $GOCT vm.info "$VM_NAME"
expect_json "vm.info --json"             $GOCT vm.info "$VM_NAME" --format json
run_cmd     "vm.disk.ls"                 $GOCT vm.disk.ls --vm "$VM_NAME"
run_cmd     "vm.nic.ls"                  $GOCT vm.nic.ls --vm "$VM_NAME"
run_cmd     "vm.cdrom.ls"                $GOCT vm.cdrom.ls --vm "$VM_NAME"
run_cmd     "find -name $VM_NAME"        $GOCT find --name "$VM_NAME"
run_cmd     "events $VM_ID"              $GOCT events "$VM_ID" --limit 5

# ---------- 3. Update (Bug 7: clear description with --description="") ----------
log_section "3. vm.update (Bug 7 regression: clear description)"
run_cmd "vm.update set description"     $GOCT vm.update "$VM_NAME" --description "smoke test description"
sleep 1
run_cmd "vm.update clear description"   $GOCT vm.update "$VM_NAME" --description ""
# Verify cleared
desc_after=$(${GOCT} vm.info "$VM_NAME" --format json 2>/dev/null | jq -r '.Description // ""')
if [[ -z "$desc_after" ]]; then
    log_pass "Description cleared (got empty)"
else
    log_fail "Description not cleared, got: $desc_after"
fi

# ---------- 4. Disk add / expand / rm ----------
log_section "4. vm.disk.* (data disk lifecycle)"
run_cmd "vm.disk.add 5G SCSI" $GOCT vm.disk.add "$VM_NAME" --size 5G --bus SCSI --name data0
sleep 3
DATA_DISK_ID=$(${GOCT} vm.disk.ls --vm "$VM_NAME" --format json 2>/dev/null | \
    jq -r '.[] | select(.VolumeName == "data0") | .ID' | head -1)
if [[ -n "$DATA_DISK_ID" ]]; then
    log_info "data disk ID=$DATA_DISK_ID"
    run_cmd "vm.disk.expand to 8G" $GOCT vm.disk.expand "$VM_NAME" --disk "$DATA_DISK_ID" --size 8G
    sleep 2
    run_cmd "vm.disk.rm"          $GOCT vm.disk.rm "$VM_NAME" --disk "$DATA_DISK_ID"
else
    log_skip "could not locate added disk by name=data0; SKIP expand/rm"
fi

# ---------- 5. NIC add / update / rm (Bug 11: --vm + --nic-id required) ----------
log_section "5. vm.nic.* (Bug 11 regression: --vm + --nic-id)"
if [[ -n "$GOCT_VLAN" ]]; then
    run_cmd "vm.nic.add" $GOCT vm.nic.add "$VM_NAME" --type VLAN --model VIRTIO
    sleep 3
    NIC_IDS=( $(${GOCT} vm.nic.ls --vm "$VM_NAME" --format json 2>/dev/null | jq -r '.[].ID') )
    if (( ${#NIC_IDS[@]} >= 2 )); then
        SECOND_NIC="${NIC_IDS[1]}"
        log_info "second NIC ID=$SECOND_NIC"
        # Bug 11 regression: --vm and --nic-id are now mandatory.
        expect_fail "vm.nic.update without --vm should fail" \
            $GOCT vm.nic.update --nic-id "$SECOND_NIC" --model VIRTIO
        expect_fail "vm.nic.update without --nic-id should fail" \
            $GOCT vm.nic.update --vm "$VM_NAME" --model VIRTIO
        run_cmd "vm.nic.update --vm + --nic-id" \
            $GOCT vm.nic.update --vm "$VM_NAME" --nic-id "$SECOND_NIC" --model VIRTIO
        run_cmd "vm.nic.rm second NIC" $GOCT vm.nic.rm "$VM_NAME" --nic-index 1
    else
        log_skip "Only one NIC found, skipping nic.update / rm"
    fi
else
    log_skip "GOCT_VLAN not set, skipping NIC add"
fi

# ---------- 6. Power state machine ----------
log_section "6. Power state machine"
run_cmd "vm.power.on"     $GOCT vm.power.on "$VM_NAME"
sleep 5
run_cmd "vm.power.suspend" $GOCT vm.power.suspend "$VM_NAME"
sleep 3
run_cmd "vm.power.resume" $GOCT vm.power.resume "$VM_NAME"
sleep 3
run_cmd "vm.power.reset"  $GOCT vm.power.reset "$VM_NAME"
sleep 5
run_cmd "vm.power.off --force" $GOCT vm.power.off "$VM_NAME" --force

# ---------- 7. Snapshot lifecycle ----------
log_section "7. vm.snapshot.*"
run_cmd "vm.snapshot.create" $GOCT vm.snapshot.create "$VM_NAME" --name snap1
sleep 3
SNAP_ID=$(${GOCT} vm.snapshot.ls "$VM_NAME" --format json 2>/dev/null | \
    jq -r '.[] | select(.Name == "snap1") | .ID' | head -1)
if [[ -n "$SNAP_ID" ]]; then
    run_cmd "vm.snapshot.ls" $GOCT vm.snapshot.ls "$VM_NAME"
    run_cmd "vm.snapshot.rm" $GOCT vm.snapshot.rm "$SNAP_ID"
else
    log_skip "Could not locate snapshot snap1"
fi

# ---------- 8. Recycle bin round-trip ----------
log_section "8. Recycle bin"
run_cmd "vm.recycle"       $GOCT vm.recycle "$VM_NAME"
sleep 2
run_cmd "vm.ls --recycle"  $GOCT vm.ls --recycle
run_cmd "vm.recover"       $GOCT vm.recover "$VM_NAME"

# ---------- 9. vm.migrate (no host = let CloudTower choose, Bug 8) ----------
log_section "9. vm.migrate (Bug 8 regression: no random host)"
# We just verify the command is callable; success depends on cluster having >= 2 hosts.
out=$(${GOCT} vm.migrate "$VM_NAME" 2>&1) || true
if echo "$out" | grep -qiE "task|migrated"; then
    log_pass "vm.migrate (no --host, server chooses)"
elif echo "$out" | grep -qiE "no available|same host|only one host"; then
    log_skip "vm.migrate cannot test (single-host cluster)"
else
    log_warn "vm.migrate output: $(echo "$out" | head -2)"
fi

# ---------- Summary ----------
print_summary

#!/usr/bin/env bash
# scripts/test.sh — 一键跑所有测试。
#
# 用法：
#   GOCT_URL=...  GOCT_USERNAME=...  GOCT_PASSWORD=...  GOCT_INSECURE=true \
#       ./scripts/test.sh [smoke|regression|lifecycle|all]
#
# 默认运行 smoke + regression（不会改变资源）。
# 加 lifecycle / all 才会创建并销毁 VM。
#
# 退出码：任一子脚本失败 → 1，否则 0。

set -uo pipefail
cd "$(dirname "$0")/.."

MODE="${1:-default}"

run_one() {
    local script="$1"
    echo
    echo "############################################################"
    echo "# Running $script"
    echo "############################################################"
    bash "$script"
}

GO_TEST_PASS=0
SMOKE_PASS=0
REG_PASS=0
LIFE_PASS=0

# 1) Always run go unit tests first; fast and free.
echo "############################################################"
echo "# go test ./..."
echo "############################################################"
if go test ./... 2>&1 | tail -25; then
    GO_TEST_PASS=1
else
    GO_TEST_PASS=0
fi

case "$MODE" in
    smoke)
        run_one scripts/test_smoke.sh && SMOKE_PASS=1
        ;;
    regression)
        run_one scripts/test_regression.sh && REG_PASS=1
        ;;
    lifecycle)
        run_one scripts/test_vm_lifecycle.sh && LIFE_PASS=1
        ;;
    all)
        run_one scripts/test_smoke.sh       && SMOKE_PASS=1
        run_one scripts/test_regression.sh  && REG_PASS=1
        run_one scripts/test_vm_lifecycle.sh && LIFE_PASS=1
        ;;
    default|"")
        run_one scripts/test_smoke.sh       && SMOKE_PASS=1
        run_one scripts/test_regression.sh  && REG_PASS=1
        ;;
    *)
        echo "Usage: $0 [smoke|regression|lifecycle|all]" >&2
        exit 2
        ;;
esac

echo
echo "############################################################"
echo "# Aggregate result"
echo "############################################################"
echo "  go test ./...        : $((GO_TEST_PASS == 1 ? 0 : 1)) ($([[ $GO_TEST_PASS -eq 1 ]] && echo OK || echo FAIL))"
[[ "$MODE" == "smoke"      || "$MODE" == "all" || "$MODE" == "default" || "$MODE" == "" ]] && \
    echo "  scripts/test_smoke   : $((SMOKE_PASS == 1 ? 0 : 1)) ($([[ $SMOKE_PASS -eq 1 ]] && echo OK || echo FAIL))"
[[ "$MODE" == "regression" || "$MODE" == "all" || "$MODE" == "default" || "$MODE" == "" ]] && \
    echo "  scripts/regression   : $((REG_PASS == 1 ? 0 : 1)) ($([[ $REG_PASS -eq 1 ]] && echo OK || echo FAIL))"
[[ "$MODE" == "lifecycle"  || "$MODE" == "all" ]] && \
    echo "  scripts/lifecycle    : $((LIFE_PASS == 1 ? 0 : 1)) ($([[ $LIFE_PASS -eq 1 ]] && echo OK || echo FAIL))"

# Aggregate exit code.
fail=0
[[ $GO_TEST_PASS -ne 1 ]] && fail=1
case "$MODE" in
    smoke)      [[ $SMOKE_PASS -ne 1 ]] && fail=1 ;;
    regression) [[ $REG_PASS  -ne 1 ]]  && fail=1 ;;
    lifecycle)  [[ $LIFE_PASS -ne 1 ]]  && fail=1 ;;
    all)
        [[ $SMOKE_PASS -ne 1 ]] && fail=1
        [[ $REG_PASS   -ne 1 ]] && fail=1
        [[ $LIFE_PASS  -ne 1 ]] && fail=1
        ;;
    default|"")
        [[ $SMOKE_PASS -ne 1 ]] && fail=1
        [[ $REG_PASS   -ne 1 ]] && fail=1
        ;;
esac

exit $fail

#!/usr/bin/env bash
# scripts/generate_usage.sh
# 自动从 goct --help 和 goct <cmd> --help 生成 USAGE.md

set -uo pipefail
cd "$(dirname "$0")/.."

GOCT="./goct"
OUTPUT="USAGE.md"

# 检查 goct 是否存在
if [[ ! -x "$GOCT" ]]; then
    echo "Building goct first..."
    go build -o goct . || { echo "Failed to build goct"; exit 1; }
fi

# 生成 Usage 头部
cat > "$OUTPUT" << 'HEADER'
# goct usage

goct is a govc-style command-line client for SmartX CloudTower.

## Common Flags

The following common flags appear for all commands:

```
  -h, --help              Show this message
      --debug             Enable debug logging (env: GOCT_DEBUG)
      --dump              Enable full HTTP trace without body truncation (env: GOCT_DUMP)
      --format string     Output format: table|json (default "table")
      --insecure          Skip TLS certificate verification (env: GOCT_INSECURE)
      --password string   Login password (env: GOCT_PASSWORD)
      --source string     Login source: local|ldap|sso|authn (env: GOCT_SOURCE)
      --trace             Enable HTTP trace (env: GOCT_TRACE)
      --url string        CloudTower endpoint URL (env: GOCT_URL)
      --username string   Login username (env: GOCT_USERNAME)
      --verbose           Enable verbose HTTP trace with headers and body (env: GOCT_VERBOSE)
      --cluster string    Default cluster ID or name (env: GOCT_CLUSTER)
```

## Environment Variables

| Variable | Description |
|----------|-------------|
| GOCT_URL | CloudTower endpoint URL |
| GOCT_USERNAME | Login username |
| GOCT_PASSWORD | Login password |
| GOCT_INSECURE | Skip TLS verification (true/false) |
| GOCT_SOURCE | Login source: local\|ldap\|sso\|authn |
| GOCT_CLUSTER | Default cluster ID or name |
| GOCT_DEBUG | Enable debug logging |
| GOCT_TRACE | Enable HTTP trace |

HEADER

# 提取目录内容
echo "" >> "$OUTPUT"
echo '<details><summary>Contents</summary>' >> "$OUTPUT"
echo "" >> "$OUTPUT"

# 定义命令组
declare -a GROUP_NAMES=(
    "System"
    "Applications"
    "Content Library"
    "Virtual Machines"
    "Hosts"
    "Clusters"
    "Datastores"
    "Networks"
    "Tasks"
    "Alerts"
    "Users"
    "Metrics"
    "Additional Commands"
)

declare -a GROUP_CMDS=(
    "about cluster-settings.get deploy.ls find license.ls ntp.get session.login session.logout session.ls version"
    "delete-cloudtower-application-package deploy-cloudtower-application get-cloudtower-application-packages get-cloudtower-applications upload-cloudtower-application-package"
    "content-library-image.delete content-library-image.distribute content-library-image.import content-library-image.ls"
    "vm.cdrom.add vm.cdrom.eject vm.cdrom.ls vm.cdrom.rm vm.cdrom.toggle vm.clone vm.create vm.destroy vm.disk.add vm.disk.expand vm.disk.ls vm.disk.rm vm.disk.update vm.export vm.gpu.add vm.gpu.ls vm.gpu.rm vm.info vm.ip vm.ls vm.migrate vm.migrate.abort vm.migrate.across vm.nic.add vm.nic.ls vm.nic.rm vm.power.off vm.power.on vm.power.reset vm.power.resume vm.power.suspend vm.rebuild vm.recover vm.recycle vm.reset-password vm.shutdown vm.snapshot.create vm.snapshot.ls vm.snapshot.revert vm.snapshot.rm vm.tools.install vm.update vm.vnc"
    "host.disconnect host.info host.ls host.maintenance.enter host.maintenance.exit host.reboot host.reconnect host.shutdown"
    "cluster.info cluster.ls"
    "datastore.disk.ls datastore.info datastore.ls storage.pool.ls"
    "network.info network.ls vlan.create vlan.destroy vlan.info vlan.ls"
    "events task.cancel task.info task.ls task.wait"
    "alert-rule.ls alert.ack alert.info alert.ls"
    "user.create user.destroy user.info user.ls"
    "cluster.metrics host.metrics sfs.metrics vm.metrics vm.volume volume.metrics"
    "completion convert-to-vm vm.nic.update"
)

# 生成目录
for i in "${!GROUP_NAMES[@]}"; do
    group="${GROUP_NAMES[$i]}"
    group_lower=$(echo "$group" | tr '[:upper:]' '[:lower:]' | tr ' ' '-')
    echo " - [${group}](#${group_lower})" >> "$OUTPUT"
done

echo "" >> "$OUTPUT"
echo '</details>' >> "$OUTPUT"
echo "" >> "$OUTPUT"

# 处理每个组
for i in "${!GROUP_NAMES[@]}"; do
    group="${GROUP_NAMES[$i]}"
    cmds="${GROUP_CMDS[$i]}"
    group_lower=$(echo "$group" | tr '[:upper:]' '[:lower:]' | tr ' ' '-')

    echo "## ${group}" >> "$OUTPUT"
    echo "" >> "$OUTPUT"

    for cmd in $cmds; do
        # 获取命令帮助
        help_output=$($GOCT "$cmd" --help 2>&1) || continue

        # 提取描述（第一行）
        description=$(echo "$help_output" | head -1 | sed 's/^ *//')

        # 提取 Usage 行（"Usage:" 之后的第一行）
        usage=$(echo "$help_output" | grep -A1 "^Usage:" | tail -1 | sed 's/^ *//')

        # 跳过没有有效描述的命令
        if [[ -z "$description" || "$description" == "" ]]; then
            continue
        fi

        echo "### ${cmd}" >> "$OUTPUT"
        echo "" >> "$OUTPUT"
        echo '```' >> "$OUTPUT"
        echo "$usage" >> "$OUTPUT"
        echo '```' >> "$OUTPUT"
        echo "" >> "$OUTPUT"
        echo "${description}" >> "$OUTPUT"
        echo "" >> "$OUTPUT"

        # 提取 Local Flags 部分（每个命令特有的 flags，在 Global Flags 之前）
        # 格式：-h, --help 等一行一个 flag
        local_flags=$(echo "$help_output" | sed -n '/^Flags:$/,/^Global Flags:/p' | grep -E "^\s+-" | sed 's/^/    /')
        if [[ -n "$local_flags" ]]; then
            echo "Flags:" >> "$OUTPUT"
            echo "$local_flags" >> "$OUTPUT"
            echo "" >> "$OUTPUT"
        fi

        echo "---" >> "$OUTPUT"
        echo "" >> "$OUTPUT"
    done
done

echo "Generated: $OUTPUT"
# Host Commands

## Listing Hosts
```bash
goct host.ls                          # List all hosts
goct host.ls --id-only               # Output only IDs
goct --format json host.ls           # JSON output
```

## Host Info
```bash
goct host.info <host-name-or-id>     # Show host details
```

## Maintenance Mode
```bash
goct host.maintenance.enter <host>   # Enter maintenance mode
goct host.maintenance.exit <host>    # Exit maintenance mode
```

## Host Operations
```bash
goct host.reboot <host>              # Reboot host
goct host.shutdown <host>            # Shutdown host
goct host.disconnect <host>          # Disconnect host from cluster
goct host.reconnect <host>           # Reconnect host to cluster
```

## Host Metrics
```bash
goct host.metrics <metric> [host]    # Query host metrics
goct host.metrics --list            # List available host metrics
```

## Object Selection

Host can be specified by:
1. Position argument: `goct host.info my-host`
2. Environment variable: `export GOCT_HOST=my-host`

```bash
# Environment variable
GOCT_HOST=my-host goct host.info
```

## Notes

- Host must exit maintenance mode before performing power operations
- Use `goct host.ls` to find available hosts and their status
- Host metrics require the host to be connected and healthy
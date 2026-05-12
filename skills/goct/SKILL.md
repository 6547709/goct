---
name: goct
description: goct CLI tool for SmartX CloudTower management. Use when interacting with CloudTower via goct CLI commands. Trigger scenarios: VM/Host/Cluster/Datastore/Network operations, metrics queries, task management, CloudTower resource management via CLI.
---

# goct — CloudTower CLI

## Overview

goct is a govc-style command-line client for SmartX CloudTower, providing 52+ tier-1 commands covering VM, snapshot, host, cluster, storage, network, VLAN, task, alert, and user management.

## How to Use This Skill

This skill is split into multiple reference files for on-demand loading.

**Directory structure:**
```
references/
├── vm.md           # VM operations (power, snapshot, clone, migrate, disk/nic/gpu)
├── host.md         # Host operations (maintenance mode, connect/disconnect)
├── cluster.md      # Cluster operations (info, list)
├── datastore.md    # Datastore operations (list, info, physical disks)
├── network.md      # Network/VLAN operations
├── metrics.md      # Metrics queries (vm.metrics, host.metrics, cluster.metrics, etc.)
└── task.md         # Task operations (list, wait, cancel)
```

**Navigation:**
1. Find the object type you need in the command reference below
2. Read `references/<object>.md` for detailed command usage
3. Commands follow govc-style dot notation: `goct vm.power.on`, `goct host.maintenance.enter`

## Connection Configuration

### Environment Variables (Recommended)
```bash
export GOCT_URL=https://tower.example.com
export GOCT_USERNAME=admin
export GOCT_PASSWORD=secret
export GOCT_CLUSTER=clusterName    # Default cluster, avoids passing --cluster each time
export GOCT_SOURCE=local           # Auth source: local|ldap|sso|authn (default: local)
export GOCT_INSECURE=true          # Skip TLS verification for self-signed certs
goct vm.ls
```

### CLI Flags
```bash
goct --url https://tower.example.com --username admin --password secret --insecure vm.ls
```

### Configuration File (~/.goct.yaml)
```yaml
url: https://tower.example.com
username: admin
password: secret
insecure: true
source: local
```

**Priority:** CLI flags > Environment variables > Config file

## Command Index

### VM (Virtual Machines)
| Command | Description |
|---------|-------------|
| `goct vm.ls` | List VMs |
| `goct vm.info <vm>` | Show VM details |
| `goct vm.ip <vm>` | Show VM IP addresses |
| `goct vm.power.on <vm>` | Power on VM |
| `goct vm.power.off <vm>` | Power off VM |
| `goct vm.shutdown <vm>` | Graceful shutdown |
| `goct vm.create ...` | Create VM |
| `goct vm.clone ...` | Clone VM |
| `goct vm.migrate <vm>` | Migrate to another host |
| `goct vm.snapshot.*` | Snapshot operations |
| `goct vm.disk.*` | Disk operations (add/expand/rm) |
| `goct vm.nic.*` | NIC operations (add/rm/update) |
| `goct vm.gpu.*` | GPU operations (add/rm) |

### Host
| Command | Description |
|---------|-------------|
| `goct host.ls` | List hosts |
| `goct host.info <host>` | Show host details |
| `goct host.maintenance.enter <host>` | Enter maintenance mode |
| `goct host.maintenance.exit <host>` | Exit maintenance mode |
| `goct host.reboot <host>` | Reboot host |
| `goct host.shutdown <host>` | Shutdown host |

### Cluster
| Command | Description |
|---------|-------------|
| `goct cluster.ls` | List clusters |
| `goct cluster.info <cluster>` | Show cluster details |

### Datastore
| Command | Description |
|---------|-------------|
| `goct datastore.ls` | List datastores |
| `goct datastore.info <datastore>` | Show datastore details |
| `goct datastore.disk.ls` | List physical disks |

### Network
| Command | Description |
|---------|-------------|
| `goct network.ls` | List virtual switches (VDS) |
| `goct vlan.ls` | List VLANs |
| `goct vlan.create ...` | Create VLAN |
| `goct vlan.destroy <vlan>` | Delete VLAN |

### Metrics
| Command | Description |
|---------|-------------|
| `goct vm.metrics <metric> [vm]` | Query VM metrics (elf_vm_*, elf_shared_vm_disk_*) |
| `goct host.metrics <metric> [host]` | Query host metrics (host_*) |
| `goct cluster.metrics <metric> [cluster]` | Query cluster metrics (zbs_cluster_*) |
| `goct volume.metrics <metric> [volume]` | Query volume metrics (zbs_volume_*) |
| `goct vm.volume <metric> [vm]` | Query VM volume metrics (elf_vm_disk_overall_*) |
| `goct <obj>.metrics --list` | List available metrics (no connection required) |

### Task
| Command | Description |
|---------|-------------|
| `goct task.ls` | List tasks |
| `goct task.wait <task>` | Wait for task completion |

### Session
| Command | Description |
|---------|-------------|
| `goct session.login` | Force login and cache token |
| `goct session.ls` | List cached session files |
| `goct about` | Show server version and connection info |

## Output Formats

```bash
# Table output (default)
goct vm.ls

# JSON output
goct --format json vm.ls

# ID only (for scripting)
goct vm.ls --id-only
```

## Quick Examples

```bash
# List all VMs
goct vm.ls

# Query VM metrics
goct vm.metrics elf_vm_cpu_overall_usage_percent my-vm --latest

# List available VM metrics (no auth required)
goct vm.metrics --list

# Power on a VM
goct vm.power.on my-vm

# Create snapshot
goct vm.snapshot.create my-vm --name "before-upgrade"

# Migrate VM to another host
goct vm.migrate my-vm --host target-host

# Check tasks
goct task.ls
```

## Common Issues

- **Empty metrics results**: Check if metric name is correct. Use `--list` to see valid metrics for each object type.
- **Auth failures**: Ensure `--source local` is set if using local authentication.
- **Connection errors**: Use `--insecure` for self-signed certificates.
# Metrics Commands

Metrics commands allow querying performance data for VMs, hosts, clusters, and volumes. Each metrics type has its own namespace and available metrics list.

## VM Metrics
```bash
# Query VM CPU/memory/network metrics
goct vm.metrics elf_vm_cpu_overall_usage_percent <vm>
goct vm.metrics elf_vm_memory_usage_percent <vm>
goct vm.metrics elf_vm_network_receive_bytes <vm>

# Query VM disk metrics (shared volumes)
goct vm.metrics elf_shared_vm_disk_read_iops <vm>
goct vm.metrics elf_shared_vm_disk_avg_readwrite_latency_ns <vm>

# List available VM metrics (no auth required)
goct vm.metrics --list
```

## VM Volume Metrics (elf_vm_disk_overall_*)
```bash
# Query aggregated VM volume metrics
goct vm.volume elf_vm_disk_overall_read_iops <vm>
goct vm.volume elf_vm_disk_overall_logical_size_bytes <vm>
goct vm.volume elf_vm_disk_overall_avg_readwrite_latency_ns <vm>

# List available VM volume metrics
goct vm.volume --list
```

## Host Metrics
```bash
# Query host metrics
goct host.metrics host_cpu_overall_usage_percent <host>
goct host.metrics host_memory_usage_percent <host>
goct host.metrics host_disk_read_iops <host>

# List available host metrics
goct host.metrics --list
```

## Cluster Metrics
```bash
# Query cluster (ZBS) metrics
goct cluster.metrics zbs_cluster_used_data_space_bytes <cluster>
goct cluster.metrics zbs_cluster_data_reduction_ratio <cluster>
goct cluster.metrics zbs_cluster_readwrite_iops <cluster>

# List available cluster metrics
goct cluster.metrics --list
```

## Volume Metrics (Independent Volumes)
```bash
# Query independent volume metrics (not VM-attached)
goct volume.metrics zbs_volume_read_iops <volume>
goct volume.metrics zbs_volume_avg_readwrite_latency_ns <volume>

# List available volume metrics
goct volume.metrics --list
```

## Common Options
```bash
--range <time>    # Time range: 5m, 1h, 1d, 7d (default: 5m)
--latest         # Show only the latest value
--format <fmt>   # Output format: table (default), json
```

## Metric Naming Convention

| Prefix | Object | Example |
|--------|--------|---------|
| `elf_vm_*` | VM | `elf_vm_cpu_overall_usage_percent` |
| `elf_shared_vm_disk_*` | VM shared volumes | `elf_shared_vm_disk_read_iops` |
| `elf_vm_disk_overall_*` | VM all volumes aggregated | `elf_vm_disk_overall_read_iops` |
| `host_*` | Host | `host_cpu_overall_usage_percent` |
| `zbs_cluster_*` | Cluster (ZBS) | `zbs_cluster_used_data_space_bytes` |
| `zbs_volume_*` | Independent volume | `zbs_volume_read_iops` |

## Notes

- Use `--list` flag to see all available metrics without requiring CloudTower connection
- Metric names are case-sensitive
- Empty results may indicate: wrong metric name, no data in time range, or VM/host not running
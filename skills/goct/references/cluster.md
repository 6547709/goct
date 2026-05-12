# Cluster Commands

## Listing Clusters
```bash
goct cluster.ls                        # List all clusters
goct cluster.ls --id-only             # Output only IDs
goct --format json cluster.ls          # JSON output
```

## Cluster Info
```bash
goct cluster.info <cluster-name-or-id>  # Show cluster details
```

## Cluster Metrics
```bash
goct cluster.metrics <metric> [cluster]  # Query cluster metrics
goct cluster.metrics --list               # List available cluster metrics
```

## Notes

- Cluster is often used as the default scope for VM operations
- Set `GOCT_CLUSTER` environment variable to avoid passing `--cluster` flag repeatedly
- Cluster metrics cover storage (ZBS) aspects like capacity, performance, data reduction
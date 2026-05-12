# VM Commands

## Listing VMs
```bash
goct vm.ls                          # List all VMs
goct vm.ls --id-only               # Output only IDs (for scripting)
goct --format json vm.ls            # JSON output
```

## VM Info
```bash
goct vm.info <vm-name-or-id>       # Show VM details
goct vm.ip <vm-name-or-id>          # Show VM IP addresses
```

## Power Operations
```bash
goct vm.power.on <vm>               # Power on
goct vm.power.off <vm>              # Power off (force)
goct vm.shutdown <vm>               # Graceful shutdown (guest OS)
goct vm.power.reset <vm>            # Force reset
goct vm.power.suspend <vm>          # Suspend
goct vm.power.resume <vm>           # Resume from suspend
```

## Snapshot Operations
```bash
goct vm.snapshot.create <vm> --name <snapshot-name>    # Create snapshot
goct vm.snapshot.ls <vm>                                # List snapshots
goct vm.snapshot.revert <vm> --snapshot <name>          # Revert to snapshot
goct vm.snapshot.rm <vm> --snapshot <name>              # Delete snapshot
```

## VM Lifecycle
```bash
goct vm.create --name <name> --cluster <cluster> --cpu 2 --memory 4GiB --disk 40GiB
goct vm.clone <vm> --name <new-vm>
goct vm.destroy <vm>                                    # Delete VM
goct vm.recycle <vm>                                    # Move to recycle bin
goct vm.recover <vm>                                    # Recover from recycle bin
goct vm.update <vm> --name <new-name>                  # Update name/description
```

## Disk Operations
```bash
goct vm.disk.add <vm> --size 100GiB                    # Add disk
goct vm.disk.expand <vm> --disk <index> --size 200GiB  # Expand disk
goct vm.disk.rm <vm> --disk <index>                    # Remove disk
goct disk.ls <vm>                                       # List VM disks
goct disk.update <vm> --disk <index> --options ...     # Update disk settings
```

## NIC Operations
```bash
goct vm.nic.add <vm> --network <network-name>          # Add NIC
goct vm.nic.rm <vm> --index <nic-index>               # Remove NIC
goct nic.ls <vm>                                       # List NICs
goct nic.update <vm> --index <nic-index> --options    # Update NIC config
```

## GPU Operations
```bash
goct vm.gpu.add <vm> --gpu <gpu-device-id>            # Add GPU
goct vm.gpu.rm <vm> --index <gpu-index>              # Remove GPU
goct gpu.ls <vm>                                       # List GPU devices
```

## Migration
```bash
goct vm.migrate <vm>                      # Migrate to random host
goct vm.migrate <vm> --host <target-host> # Migrate to specific host
goct vm.migrate.across <vm> --cluster <target-cluster> # Cross-cluster migration
```

## Other
```bash
goct vm.vnc <vm>                          # Get VNC connection info
goct vm.tools.install <vm>                # Install VMware Tools
goct vm.export <vm> --name <export-name>  # Export as OVF
```

## Object Selection

VM can be specified by:
1. Position argument: `goct vm.power.on my-vm`
2. Environment variable: `export GOCT_VM=my-vm`
3. stdin pipe (one ID per line)

```bash
# Environment variable
GOCT_VM=my-vm goct vm.power.off

# Pipe multiple VMs
echo -e "vm1\nvm2\nvm3" | goct vm.power.on
```
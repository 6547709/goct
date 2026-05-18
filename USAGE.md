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


<details><summary>Contents</summary>

 - [System](#system)
 - [Applications](#applications)
 - [Content Library](#content-library)
 - [Virtual Machines](#virtual-machines)
 - [Hosts](#hosts)
 - [Clusters](#clusters)
 - [Datastores](#datastores)
 - [Networks](#networks)
 - [Tasks](#tasks)
 - [Alerts](#alerts)
 - [Users](#users)
 - [Metrics](#metrics)
 - [Additional Commands](#additional-commands)

</details>

## System

### about

```
goct about [flags]
```

Show CloudTower server version and connection info

Flags:
      -h, --help   help for about

---

### cluster-settings.get

```
goct cluster-settings.get [flags]
```

Get cluster settings

Flags:
          --cluster string   Cluster ID
      -h, --help             help for cluster-settings.get

---

### deploy.ls

```
goct deploy.ls [flags]
```

List deploys

Flags:
      -h, --help          help for deploy.ls
          --limit int32   Limit number of results (0 = no limit)
          --name string   Filter by name (substring match)
          --skip int32    Skip N results (pagination)

---

### find

```
goct find [flags]
```

Find resources of a given type matching a name pattern.

Flags:
          --cluster string   Restrict to this cluster (ID or name); ignored for cluster-less resource types
      -h, --help             help for find
          --id-only          Print only IDs (one per line)
          --json             JSON output
          --limit int32      Maximum results (0 = unlimited)
          --name string      Filter by name (substring match)
          --type string      Resource type filter: m|h|c|d|n|v|f|g|t|l|u|a or full name (default: all)

---

### license.ls

```
goct license.ls [flags]
```

List licenses

Flags:
      -h, --help          help for license.ls
          --limit int32   Limit number of results (0 = no limit)
          --name string   Filter by name (substring match)
          --skip int32    Skip N results (pagination)

---

### ntp.get

```
goct ntp.get [flags]
```

Get NTP service URL

Flags:
      -h, --help   help for ntp.get

---

### session.login

```
goct session.login [flags]
```

Force login and save session token to local cache

Flags:
      -h, --help   help for session.login

---

### session.logout

```
goct session.logout [flags]
```

Delete cached session token for a given host and user

Flags:
      -h, --help          help for session.logout
          --url string    CloudTower URL
          --user string   Username

---

### session.ls

```
goct session.ls [flags]
```

List locally cached session files

Flags:
      -h, --help   help for session.ls

---

### version

```
goct version [flags]
```

Print goct version

Flags:
      -h, --help   help for version

---

## Applications

### delete-cloudtower-application-package

```
goct delete-cloudtower-application-package <id> [flags]
```

Delete CloudTower application package

Flags:
      -h, --help   help for delete-cloudtower-application-package

---

### deploy-cloudtower-application

```
goct deploy-cloudtower-application <name> [flags]
```

Deploy CloudTower application

Flags:
      -h, --help             help for deploy-cloudtower-application
          --package string   Target package ID

---

### get-cloudtower-application-packages

```
goct get-cloudtower-application-packages [flags]
```

List CloudTower application packages

Flags:
      -h, --help          help for get-cloudtower-application-packages
          --limit int32   Limit number of results (0 = no limit)
          --name string   Filter by name (substring match)
          --skip int32    Skip N results (pagination)

---

### get-cloudtower-applications

```
goct get-cloudtower-applications [flags]
```

List CloudTower applications

Flags:
      -h, --help          help for get-cloudtower-applications
          --limit int32   Limit number of results (0 = no limit)
          --name string   Filter by name (substring match)
          --skip int32    Skip N results (pagination)

---

### upload-cloudtower-application-package

```
goct upload-cloudtower-application-package <path> <name> [flags]
```

Upload CloudTower application package

Flags:
      -h, --help   help for upload-cloudtower-application-package

---

## Content Library

### content-library-image.delete

```
goct content-library-image.delete <id> [flags]
```

Delete content library image

Flags:
      -h, --help   help for content-library-image.delete

---

### content-library-image.distribute

```
goct content-library-image.distribute <id> [flags]
```

Distribute content library image to clusters

Flags:
          --cluster strings   Target cluster IDs
      -h, --help              help for content-library-image.distribute

---

### content-library-image.import

```
goct content-library-image.import <path> <name> [flags]
```

Import content library image

Flags:
          --cluster string   Target cluster ID
      -h, --help             help for content-library-image.import

---

### content-library-image.ls

```
goct content-library-image.ls [flags]
```

List content library images

Flags:
      -h, --help          help for content-library-image.ls
          --limit int32   Limit number of results (0 = no limit)
          --name string   Filter by name (substring match)
          --skip int32    Skip N results (pagination)

---

## Virtual Machines

### vm.cdrom.add

```
goct vm.cdrom.add [vm-name|vm-id] [flags]
```

Add a CD-ROM drive to a VM, optionally with an ISO mounted.

Flags:
      -h, --help         help for vm.cdrom.add
          --iso string   ISO path or image ID to mount

---

### vm.cdrom.eject

```
goct vm.cdrom.eject [cdrom-id] [flags]
```

Eject the currently mounted ISO from a CD-ROM drive.

Flags:
      -h, --help   help for vm.cdrom.eject

---

### vm.cdrom.ls

```
goct vm.cdrom.ls [flags]
```

List all CD-ROM devices attached to a VM.

Flags:
      -h, --help        help for vm.cdrom.ls
          --vm string   VM name or ID

---

### vm.cdrom.rm

```
goct vm.cdrom.rm [vm-name|vm-id] [cdrom-id] [flags]
```

Remove a CD-ROM drive from a virtual machine.

Flags:
      -h, --help   help for vm.cdrom.rm

---

### vm.cdrom.toggle

```
goct vm.cdrom.toggle [cdrom-id] [flags]
```

Enable or disable a CD-ROM device.

Flags:
          --disabled   Disable the CD-ROM (default: enable)
      -h, --help       help for vm.cdrom.toggle

---

### vm.clone

```
goct vm.clone [source-name|id] [flags]
```

Clone a VM

Flags:
          --cluster string   Target cluster ID (optional, default same)
      -h, --help             help for vm.clone
          --name string      Name for cloned VM (required)

---

### vm.create

```
goct vm.create [flags]
```

Create a new virtual machine.

Flags:
          --cluster string         Target cluster ID
          --description string     Description
          --disk stringArray       Disk spec, repeatable: size=10g[,bus=SCSI][,name=diskN][,index=N][,boot=N][,iops=N]
          --dns stringArray        DNS nameserver (repeatable, e.g. --dns 8.8.8.8)
          --firmware string        Firmware: BIOS or UEFI (empty = use template default)
          --from-template string   Template ID to create VM from
          --full-copy              Full copy when creating from template
          --ha string              Enable HA: true|false (empty = use cluster default)
      -h, --help                   help for vm.create
          --hostname string        Cloud-init hostname
          --memory int             Memory in MiB (0 = use template default)
          --name string            VM name
          --network stringArray    NIC config: nic=0[,ip=x][,netmask=x][,gateway=x][,route=10.0.0.0/8:192.168.1.1][,type=IPV4|DHCP]
          --nic stringArray        NIC spec, repeatable: vlan=<id|name>[,model=VIRTIO][,type=VLAN|VPC]
          --nic-model string       NIC model: E1000, SRIOV, VIRTIO (only for --from-template)
          --nic-type string        NIC type: VLAN or VPC (only for --from-template)
          --nic-vlan string        VLAN ID (only for --from-template)
          --password string        Default user password (cloud-init)
          --ssh-key stringArray    SSH public key (literal or @/path/to/key.pub)
          --user-data string       Cloud-init user_data (@/path/to/file.yaml or literal YAML)
          --vcpu int32             Number of vCPUs (0 = use template default)

---

### vm.destroy

```
goct vm.destroy [name|id] [flags]
```

Destroy (delete) a VM

Flags:
          --force   Force destroy without graceful shutdown
      -h, --help    help for vm.destroy

---

### vm.disk.add

```
goct vm.disk.add [vm-name|vm-id] [flags]
```

Add a new disk to a virtual machine.

Flags:
          --bus string              Disk bus type (SCSI, IDE, VIRTIO) (default "SCSI")
      -h, --help                    help for vm.disk.add
          --name string             Disk name
          --size string             Disk size (e.g. 100G, 50G)
          --storage-policy string   Storage policy (e.g. ELF_CP_REPLICA_2_THICK_PROVISION, ELF_CP_REPLICA_3_THICK_PROVISION)

---

### vm.disk.expand

```
goct vm.disk.expand [vm-name|vm-id] [disk-id] [flags]
```

Expand (resize) a disk attached to a VM.

Flags:
      -h, --help          help for vm.disk.expand
          --size string   New disk size (e.g. 200G)

---

### vm.disk.ls

```
goct vm.disk.ls [flags]
```

List all disks (including CD-ROMs) attached to a VM.

Flags:
      -h, --help        help for vm.disk.ls
          --id-only     Output only IDs, one per line (for scripting)
          --vm string   VM name or ID

---

### vm.disk.rm

```
goct vm.disk.rm [vm-name|vm-id] [disk-id] [flags]
```

Remove a disk from a virtual machine.

Flags:
      -h, --help   help for vm.disk.rm

---

### vm.disk.update

```
goct vm.disk.update [flags]
```

Update a VM disk configuration.

Flags:
          --bus string             New bus: SCSI / IDE / VIRTIO
          --content-image string   Mount content library image
          --disk string            Disk ID
          --elf-image string       Mount ELF image (CD-ROM)
      -h, --help                   help for vm.disk.update
          --vm string              VM name or ID (required)
          --volume string          Replace underlying VM volume by ID

---

### vm.export

```
goct vm.export [name|id] [flags]
```

Export a VM as OVF

Flags:
      -h, --help       help for vm.export
          --keep-mac   Keep MAC addresses in exported OVF

---

### vm.gpu.add

```
goct vm.gpu.add [vm-name|vm-id] <gpu-device-id> [flags]
```

Add a GPU device to a virtual machine.

Flags:
      -h, --help   help for vm.gpu.add

---

### vm.gpu.ls

```
goct vm.gpu.ls [flags]
```

List all GPU devices attached to a VM.

Flags:
      -h, --help        help for vm.gpu.ls
          --vm string   VM name or ID

---

### vm.gpu.rm

```
goct vm.gpu.rm [vm-name|vm-id] <gpu-device-id> [flags]
```

Remove a GPU device from a virtual machine.

Flags:
      -h, --help   help for vm.gpu.rm

---

### vm.info

```
goct vm.info [name|id] [flags]
```

Show VM details

Flags:
          --detail   Show detailed VM information (includes BIOS UUID, GPU/USB devices, usage stats, etc.)
      -h, --help     help for vm.info

---

### vm.ip

```
goct vm.ip [name|id] [flags]
```

Output the IP address(es) of a VM, one per line.

Flags:
      -a, --all             Output all IPs, one per line
      -h, --help            help for vm.ip
          --no-wait         Do not wait for VM tools to populate IPs
          --v4              IPv4 only
          --v6              IPv6 only
          --wait duration   Maximum wait duration when waiting for IPs (default 5m0s)

---

### vm.ls

```
goct vm.ls [flags]
```

List virtual machines with optional filtering.

Flags:
      -h, --help          help for vm.ls
          --id-only       Output only IDs, one per line (for scripting)
          --limit int32   Limit number of results (0 = no limit)
          --name string   Filter by name (substring match)
          --recycle       Show VMs in recycle bin
          --skip int32    Skip N results (pagination)

---

### vm.migrate

```
goct vm.migrate [name|id] [flags]
```

Migrate a VM to another host within the same cluster.

Flags:
      -h, --help          help for vm.migrate
          --host string   Target host name or ID (omit = let CloudTower choose)

---

### vm.migrate.abort

```
goct vm.migrate.abort [vm-name|vm-id] [flags]
```

Abort an in-progress cross-cluster migration.

Flags:
      -h, --help   help for vm.migrate.abort

---

### vm.migrate.across

```
goct vm.migrate.across [name|id] <cluster-name|id> [flags]
```

Migrate a VM to another cluster, optionally specifying a target host.

Flags:
      -h, --help          help for vm.migrate.across
          --host string   Target host name or ID (optional, auto-select if not specified)

---

### vm.nic.add

```
goct vm.nic.add [vm-name|vm-id] [flags]
```

Add a new NIC to a virtual machine.

Flags:
      -h, --help           help for vm.nic.add
          --model string   NIC model (VIRTIO, E1000, SRIOV) (default "VIRTIO")
          --type string    NIC type (VLAN, VPC) (default "VLAN")

---

### vm.nic.ls

```
goct vm.nic.ls [flags]
```

List all network interfaces attached to a VM.

Flags:
      -h, --help        help for vm.nic.ls
          --id-only     Output only IDs, one per line (for scripting)
          --vm string   VM name or ID

---

### vm.nic.rm

```
goct vm.nic.rm [vm-name|vm-id] <nic-index> [flags]
```

Remove a NIC from a virtual machine by its index.

Flags:
      -h, --help   help for vm.nic.rm

---

### vm.power.off

```
goct vm.power.off [name|id] [flags]
```

Shut down / power off a VM

Flags:
          --force   Force power off (skip graceful shutdown)
      -h, --help    help for vm.power.off

---

### vm.power.on

```
goct vm.power.on [name|id] [flags]
```

Power on a VM

Flags:
      -h, --help   help for vm.power.on

---

### vm.power.reset

```
goct vm.power.reset [name|id] [flags]
```

Restart / force reset a VM

Flags:
          --force   Force reset (skip graceful restart)
      -h, --help    help for vm.power.reset

---

### vm.power.resume

```
goct vm.power.resume [name|id] [flags]
```

Resume a suspended VM

Flags:
      -h, --help   help for vm.power.resume

---

### vm.power.suspend

```
goct vm.power.suspend [name|id] [flags]
```

Suspend a VM

Flags:
      -h, --help   help for vm.power.suspend

---

### vm.rebuild

```
goct vm.rebuild [vm-name|vm-id] [snapshot-id] [flags]
```

Rebuild a virtual machine from a snapshot.

Flags:
          --cluster string   Target cluster ID
      -h, --help             help for vm.rebuild
          --host string      Target host ID
          --name string      New VM name

---

### vm.recover

```
goct vm.recover [name|id] [flags]
```

Recover a VM from the recycle bin back to normal state.

Flags:
      -h, --help   help for vm.recover

---

### vm.recycle

```
goct vm.recycle [name|id] [flags]
```

Move a VM to the recycle bin. The VM can be recovered later with vm.recover.

Flags:
      -h, --help   help for vm.recycle

---

### vm.reset-password

```
goct vm.reset-password [vm-name|vm-id] [flags]
```

Reset the guest OS administrator password on a VM.

Flags:
      -h, --help              help for vm.reset-password
          --password string   New password (required)
          --username string   Guest OS username (required)

---

### vm.shutdown

```
goct vm.shutdown [name|id] [flags]
```

Send a graceful shutdown request to the VM's guest OS.

Flags:
      -h, --help   help for vm.shutdown

---

### vm.snapshot.create

```
goct vm.snapshot.create [vm-name|id] [flags]
```

Create a snapshot

Flags:
      -h, --help          help for vm.snapshot.create
          --name string   Snapshot name (required)

---

### vm.snapshot.ls

```
goct vm.snapshot.ls [vm-name|id] [flags]
```

List snapshots of a VM

Flags:
      -h, --help   help for vm.snapshot.ls

---

### vm.snapshot.revert

```
goct vm.snapshot.revert <snapshot-id> [flags]
```

Revert VM to a snapshot

Flags:
      -h, --help        help for vm.snapshot.revert
          --vm string   VM name or ID (required)

---

### vm.snapshot.rm

```
goct vm.snapshot.rm <snapshot-id> [flags]
```

Delete a snapshot

Flags:
      -h, --help   help for vm.snapshot.rm

---

### vm.tools.install

```
goct vm.tools.install [vm-name|vm-id] [flags]
```

Install VMware Tools on a virtual machine.

Flags:
      -h, --help   help for vm.tools.install

---

### vm.update

```
goct vm.update [name|id] [flags]
```

Update VM basic information (name, description).

Flags:
          --description string   New VM description ('' clears it)
      -h, --help                 help for vm.update
          --name string          New VM name (omit to keep, '' is rejected by CloudTower)

---

### vm.vnc

```
goct vm.vnc [vm-name|vm-id] [flags]
```

Get VNC connection information for a virtual machine.

Flags:
      -h, --help   help for vm.vnc

---

## Hosts

### host.disconnect

```
goct host.disconnect [name|id] [flags]
```

Disconnect a host (not supported by SDK)

Flags:
      -h, --help   help for host.disconnect

---

### host.info

```
goct host.info [name|id] [flags]
```

Show host details

Flags:
      -h, --help   help for host.info

---

### host.ls

```
goct host.ls [flags]
```

List hosts

Flags:
      -h, --help          help for host.ls
          --id-only       Output only IDs, one per line (for scripting)
          --limit int32   Limit number of results (0 = no limit)
          --name string   Filter by name (substring match)
          --skip int32    Skip N results (pagination)

---

### host.maintenance.enter

```
goct host.maintenance.enter [name|id] [flags]
```

Enter maintenance mode

Flags:
      -h, --help   help for host.maintenance.enter

---

### host.maintenance.exit

```
goct host.maintenance.exit [name|id] [flags]
```

Exit maintenance mode

Flags:
      -h, --help   help for host.maintenance.exit

---

### host.reboot

```
goct host.reboot [name|id] [flags]
```

Reboot a host

Flags:
          --force   Force reboot
      -h, --help    help for host.reboot

---

### host.reconnect

```
goct host.reconnect [name|id] [flags]
```

Reconnect a host (not supported by SDK)

Flags:
      -h, --help   help for host.reconnect

---

### host.shutdown

```
goct host.shutdown [name|id] [flags]
```

Shut down a host

Flags:
          --force   Force shutdown
      -h, --help    help for host.shutdown

---

## Clusters

### cluster.info

```
goct cluster.info <name|id> [flags]
```

Show cluster details

Flags:
      -h, --help   help for cluster.info

---

### cluster.ls

```
goct cluster.ls [flags]
```

List clusters

Flags:
      -h, --help          help for cluster.ls
          --id-only       Output only IDs, one per line (for scripting)
          --limit int32   Limit number of results (0 = no limit)
          --name string   Filter by name (substring match)
          --skip int32    Skip N results (pagination)

---

## Datastores

### datastore.disk.ls

```
goct datastore.disk.ls [flags]
```

List physical disks

Flags:
      -h, --help          help for datastore.disk.ls
          --host string   Filter by host ID

---

### datastore.info

```
goct datastore.info <name|id> [flags]
```

Show datastore details

Flags:
      -h, --help   help for datastore.info

---

### datastore.ls

```
goct datastore.ls [flags]
```

List datastores

Flags:
      -h, --help          help for datastore.ls
          --id-only       Output only IDs, one per line (for scripting)
          --limit int32   Limit number of results (0 = no limit)
          --name string   Filter by name (substring match)
          --skip int32    Skip N results (pagination)

---

### storage.pool.ls

```
goct storage.pool.ls [flags]
```

List DiskPools - the hyperconverged storage pools on each host.

Flags:
          --cluster string   Filter by cluster name or ID
      -h, --help             help for storage.pool.ls
          --id-only          Output only IDs, one per line (for scripting)
          --limit int32      Limit number of results (0 = no limit)
          --name string      Filter by name (substring match)
          --skip int32       Skip N results (pagination)

---

## Networks

### network.info

```
goct network.info <name|id> [flags]
```

Show virtual switch details

Flags:
      -h, --help   help for network.info

---

### network.ls

```
goct network.ls [flags]
```

List virtual switches (VDS)

Flags:
      -h, --help          help for network.ls
          --id-only       Output only IDs, one per line (for scripting)
          --limit int32   Limit number of results (0 = no limit)
          --name string   Filter by name (substring match)
          --skip int32    Skip N results (pagination)

---

### vlan.create

```
goct vlan.create [flags]
```

Create a VLAN

Flags:
      -h, --help          help for vlan.create
          --name string   VLAN name (required)
          --vds string    VDS ID (required)

---

### vlan.destroy

```
goct vlan.destroy <name|id> [flags]
```

Delete a VLAN

Flags:
      -h, --help   help for vlan.destroy

---

### vlan.info

```
goct vlan.info <name|id> [flags]
```

Show VLAN details

Flags:
      -h, --help   help for vlan.info

---

### vlan.ls

```
goct vlan.ls [flags]
```

List VLANs

Flags:
      -h, --help          help for vlan.ls
          --id-only       Output only IDs, one per line (for scripting)
          --limit int32   Limit number of results (0 = no limit)
          --name string   Filter by name (substring match)
          --skip int32    Skip N results (pagination)

---

## Tasks

### events

```
goct events [resource-id] [flags]
```

Show CloudTower events (backed by user_audit_log).

Flags:
          --action string     Filter by action substring (e.g. create_vm)
      -f, --follow            Follow new events (poll every 2s)
      -h, --help              help for events
          --json              JSON output
      -n, --limit int32       Max events to show (ignored in --follow) (default 50)
          --resource string   Filter by resource ID (positional arg also accepted)
          --type string       Filter by resource type (VM/HOST/CLUSTER/...)
          --user string       Filter by triggering username

---

### task.cancel

```
goct task.cancel <id> [flags]
```

Cancel a task (not supported by SDK)

Flags:
      -h, --help   help for task.cancel

---

### task.info

```
goct task.info <id> [flags]
```

Show task details

Flags:
      -h, --help   help for task.info

---

### task.ls

```
goct task.ls [flags]
```

List tasks

Flags:
      -h, --help          help for task.ls
          --id-only       Output only IDs, one per line (for scripting)
          --limit int32   Limit number of results (0 = no limit)
          --name string   Filter by name (substring match)
          --skip int32    Skip N results (pagination)

---

### task.wait

```
goct task.wait <id> [flags]
```

Wait for a task to complete

Flags:
      -h, --help   help for task.wait

---

## Alerts

### alert-rule.ls

```
goct alert-rule.ls [flags]
```

List alert rules

Flags:
          --cluster string   Filter by cluster ID
      -h, --help             help for alert-rule.ls

---

### alert.ack

```
goct alert.ack <id> [flags]
```

Acknowledge (resolve) an alert

Flags:
      -h, --help   help for alert.ack

---

### alert.info

```
goct alert.info <id> [flags]
```

Show alert details

Flags:
      -h, --help   help for alert.info

---

### alert.ls

```
goct alert.ls [flags]
```

List alerts

Flags:
      -h, --help          help for alert.ls
          --id-only       Output only IDs, one per line (for scripting)
          --limit int32   Limit number of results (0 = no limit)
          --name string   Filter by name (substring match)
          --skip int32    Skip N results (pagination)

---

## Users

### user.create

```
goct user.create [flags]
```

Create a user

Flags:
          --email string      Email address
      -h, --help              help for user.create
          --name string       Display name (required)
          --password string   Password (required)
          --role string       Role ID (required)
          --username string   Login username (required)

---

### user.destroy

```
goct user.destroy <name|id> [flags]
```

Delete a user

Flags:
      -h, --help   help for user.destroy

---

### user.info

```
goct user.info <name|id> [flags]
```

Show user details

Flags:
      -h, --help   help for user.info

---

### user.ls

```
goct user.ls [flags]
```

List users

Flags:
      -h, --help          help for user.ls
          --id-only       Output only IDs, one per line (for scripting)
          --limit int32   Limit number of results (0 = no limit)
          --name string   Filter by name (substring match)
          --skip int32    Skip N results (pagination)

---

## Metrics

### cluster.metrics

```
goct cluster.metrics <metric> [cluster-name] [flags]
```

Query Cluster metrics with optional cluster name filter. Example: cluster.metrics zbs_cluster_usage cluster001

Flags:
          --format string   Output format: table, json, chart (default "table")
      -h, --help            help for cluster.metrics
          --latest          Show only latest value
          --list            List available metrics
          --range string    Time range: 5m, 1h, 1d, 7d (default "5m")

---

### host.metrics

```
goct host.metrics <metric> [host-name] [flags]
```

Query Host metrics with optional host name filter. Example: host.metrics elf_host_cpu_usage host001

Flags:
          --format string   Output format: table, json, chart (default "table")
      -h, --help            help for host.metrics
          --latest          Show only latest value
          --list            List available metrics
          --range string    Time range: 5m, 1h, 1d, 7d (default "5m")

---

### sfs.metrics

```
goct sfs.metrics <metric> [sfs-name] [flags]
```

Query SFS metrics. This feature is not yet available in the CloudTower SDK.

Flags:
          --format string   Output format: table, json, chart (default "table")
      -h, --help            help for sfs.metrics
          --latest          Show only latest value
          --list            List available metrics
          --range string    Time range: 5m, 1h, 1d, 7d (default "5m")

---

### vm.metrics

```
goct vm.metrics <metric> [vm-name] [flags]
```

Query VM metrics with optional VM name filter. Example: vm.metrics elf_cpu_usage vm001

Flags:
          --format string   Output format: table, json, chart (default "table")
      -h, --help            help for vm.metrics
          --latest          Show only latest value
          --list            List available metrics
          --range string    Time range: 5m, 1h, 1d, 7d (default "5m")

---

### vm.volume

```
goct vm.volume <metric> [vm-name] [flags]
```

Query VM volume metrics by VM name. Metrics: elf_vm_disk_overall_logical_size_bytes, elf_vm_disk_overall_read_iops, etc. Example: vm.volume elf_vm_disk_overall_logical_size_bytes my-vm

Flags:
          --format string   Output format: table, json, chart (default "table")
      -h, --help            help for vm.volume
          --latest          Show only latest value
          --list            List available metrics
          --range string    Time range: 5m, 1h, 1d, 7d (default "5m")

---

### volume.metrics

```
goct volume.metrics <metric> [volume-name] [flags]
```

Query independent volume metrics with optional volume name filter. Example: volume.metrics zbs_volume_read_iops volume001

Flags:
          --format string   Output format: table, json, chart (default "table")
      -h, --help            help for volume.metrics
          --latest          Show only latest value
          --list            List available metrics
          --range string    Time range: 5m, 1h, 1d, 7d (default "5m")

---

## Additional Commands

### completion

```
goct completion [command]
```

Generate the autocompletion script for goct for the specified shell.

Flags:
      -h, --help   help for completion

---

### convert-to-vm

```
goct convert-to-vm [template-name|template-id] [flags]
```

Convert a content library template to a virtual machine.

Flags:
      -h, --help          help for convert-to-vm
          --name string   Name for the converted VM (default: <template>-vm)

---

### vm.nic.update

```
goct vm.nic.update [flags]
```

Update a NIC configuration on a VM.

Flags:
          --connect-vlan-id string   Connect to VLAN ID
          --disable                  Disable NIC
          --enable                   Enable NIC
          --gateway string           Gateway IP
      -h, --help                     help for vm.nic.update
          --ip string                IP address
          --mac string               MAC address
          --model string             Model (RTL8139/E1000/VIRTIO)
          --nic-id string            NIC ID (required)
          --nic-index int32          NIC index (LocalID, optional)
          --subnet-mask string       Subnet mask
          --vm string                VM name or ID (required)

---


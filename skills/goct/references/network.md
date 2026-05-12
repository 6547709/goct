# Network and VLAN Commands

## Virtual Switches (VDS)
```bash
goct network.ls                         # List virtual switches
goct network.info <switch-name-or-id>   # Show switch details
```

## VLAN Management
```bash
goct vlan.ls                            # List VLANs
goct vlan.info <vlan-name-or-id>        # Show VLAN details
goct vlan.create --name <vlan-name> --vlan-id <vlan-id> [options]
goct vlan.destroy <vlan-name-or-id>    # Delete VLAN
```

## VLAN Creation Examples
```bash
# Create VLAN with specific ID
goct vlan.create --name "web-tier" --vlan-id 100

# Create VLAN with description
goct vlan.create --name "app-tier" --vlan-id 200 --description "Application tier VLAN"
```

## Notes

- VLANs are used to segment network traffic in CloudTower
- Each VLAN has a unique VLAN ID (1-4094)
- VLANs must be associated with a virtual switch to be usable
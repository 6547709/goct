# Datastore Commands

## Listing Datastores
```bash
goct datastore.ls                        # List all datastores
goct datastore.ls --id-only             # Output only IDs
goct --format json datastore.ls         # JSON output
```

## Datastore Info
```bash
goct datastore.info <datastore-name-or-id>  # Show datastore details
```

## Physical Disk Management
```bash
goct datastore.disk.ls                  # List physical disks across all datastores
goct datastore.disk.ls --datastore <ds>  # List disks in specific datastore
```

## Notes

- Datastores in CloudTower are typically hyperconverged (ZBS)
- Physical disks can be viewed but management operations are limited
- Storage pool information: `goct storage.pool.ls`
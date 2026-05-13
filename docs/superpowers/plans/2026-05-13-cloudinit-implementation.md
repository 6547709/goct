# Cloud-Init Support for VM Create from Template

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add cloud-init support to `goct vm.create --from-template`, enabling SSH public keys, default user password, hostname, custom user_data, and static IP network configuration.

**Architecture:** Cloud-init 配置分为两部分：(1) 全局配置（hostname, password, ssh keys, user_data）；(2) 网卡静态 IP 配置（每张网卡可单独指定 IP/Netmask/Gateway）。CLI 通过可重复 flag 传递网卡配置，适配层组装 `TemplateCloudInit` 结构体。

**Tech Stack:** cloudtower-go-sdk v2 (`TemplateCloudInit`, `CloudInitNetWork`, `CloudInitNetWorkRoute`, `CloudInitNetworkTypeEnum`)

---

## File Structure

```
pkg/adapter/
  types.go           # Add CloudInitSpec, NicStaticConfig types
  vm.go              # Add CloudInit to VMCreateFromTemplateSpec; wire into CreateVMFromTemplate

cmd/vm/
  create.go          # Add --hostname, --password, --ssh-key, --user-data, --network flags

cloudtower-go-sdk/   (read-only, SDK types already exist)
  models/template_cloud_init.go
  models/cloud_init_net_work.go
  models/cloud_init_net_work_route.go
  models/cloud_init_network_type_enum.go
```

---

## Task 1: Add adapter types for CloudInit

**Files:**
- Modify: `pkg/adapter/types.go`

- [ ] **Step 1: Add CloudInitSpec and NicStaticConfig types**

```go
// CloudInitSpec describes cloud-init configuration for VM creation from template.
type CloudInitSpec struct {
    Hostname            string   // VM hostname
    DefaultUserPassword string   // Default user password
    PublicKeys          []string // SSH public keys (one per line, or repeatable flag)
    UserData            string   // Custom cloud-init user_data script (plaintext or base64)
    Networks            []NicStaticConfig // Static IP config per NIC
}

// NicStaticConfig describes static IP configuration for one NIC.
type NicStaticConfig struct {
    Index    int32   // NIC index (0-based), required
    IP       string  // Static IP address, e.g. "192.168.1.100"
    Netmask  string  // Netmask in dotted notation, e.g. "255.255.255.0"
    Gateway  string  // Default gateway IP
    Type     string  // "IPV4" (static) or "IPV4_DHCP" (DHCP); default IPV4
}

// Add to VMCreateFromTemplateSpec:
type VMCreateFromTemplateSpec struct {
    // ... existing fields ...
    CloudInit *CloudInitSpec // nil = no cloud-init, CloudTower uses template default
}
```

- [ ] **Step 2: Build to verify types compile**

Run: `go build ./pkg/adapter/...`

- [ ] **Step 3: Commit**

```bash
git add pkg/adapter/types.go
git commit -m "feat(adapter): add CloudInitSpec and NicStaticConfig types"
```

---

## Task 2: Wire CloudInit into CreateVMFromTemplate

**Files:**
- Modify: `pkg/adapter/vm.go:732-791` (CreateVMFromTemplate function)
- Modify: `pkg/adapter/vm.go:51` (VMCreateFromTemplateSpec struct to add CloudInit field)
- Modify: `pkg/adapter/client.go` if needed (fakeClient stub)

- [ ] **Step 1: Add CloudInit field to VMCreateFromTemplateSpec interface**

Find the `VMCreateFromTemplateSpec` struct definition and add:
```go
CloudInit *CloudInitSpec
```

- [ ] **Step 2: Build to see the error (CloudInit field not yet in struct)**

Run: `go build ./... 2>&1 | head -20`
Expected: FAIL (struct doesn't have CloudInit field yet - we added it to interface but not impl)

- [ ] **Step 3: Add CloudInit to CreateVMFromTemplate function**

After the existing NIC setup block (around line 780), add:
```go
// Cloud-init configuration
if spec.CloudInit != nil {
    ci := &models.TemplateCloudInit{}
    if spec.CloudInit.Hostname != "" {
        ci.Hostname = pointy.String(spec.CloudInit.Hostname)
    }
    if spec.CloudInit.DefaultUserPassword != "" {
        ci.DefaultUserPassword = pointy.String(spec.CloudInit.DefaultUserPassword)
    }
    if len(spec.CloudInit.PublicKeys) > 0 {
        ci.PublicKeys = spec.CloudInit.PublicKeys
    }
    if spec.CloudInit.UserData != "" {
        ci.UserData = pointy.String(spec.CloudInit.UserData)
    }
    if len(spec.CloudInit.Networks) > 0 {
        ci.Networks = make([]*models.CloudInitNetWork, 0, len(spec.CloudInit.Networks))
        for _, n := range spec.CloudInit.Networks {
            net := &models.CloudInitNetWork{
                NicIndex: pointy.Int32(n.Index),
            }
            // Type: IPV4 (static) or IPV4_DHCP
            switch strings.ToUpper(n.Type) {
            case "DHCP", "IPV4_DHCP":
                net.Type = models.CloudInitNetworkTypeEnumIPV4DHCP.Pointer()
            default: // IPV4 or empty
                net.Type = models.CloudInitNetworkTypeEnumIPV4.Pointer()
                if n.IP != "" {
                    net.IPAddress = pointy.String(n.IP)
                }
                if n.Netmask != "" {
                    net.Netmask = pointy.String(n.Netmask)
                }
                if n.Gateway != "" {
                    net.Routes = []*models.CloudInitNetWorkRoute{
                        {
                            Gateway: pointy.String(n.Gateway),
                            Netmask: pointy.String(n.Netmask),
                            Network: pointy.String("0.0.0.0"), // default route
                        },
                    }
                }
            }
            ci.Networks = append(ci.Networks, net)
        }
    }
    p.CloudInit = ci
}
```

- [ ] **Step 4: Add CloudInitNetworkTypeEnumIPV4 import if needed**

Check imports at top of vm.go for `models` package usage. Should already be present.

- [ ] **Step 5: Build to verify**

Run: `go build ./...`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add pkg/adapter/vm.go
git commit -m "feat(adapter): wire CloudInit into CreateVMFromTemplate"
```

---

## Task 3: Add CLI flags for cloud-init in vm.create

**Files:**
- Modify: `cmd/vm/create.go`

- [ ] **Step 1: Add cloud-init flag variables**

In `newCreate()` function, add:
```go
var (
    // ... existing vars ...
    cloudInitHostname string
    cloudInitPassword string
    cloudInitSSHKey   []string
    cloudInitUserData string
    cloudInitNetworks []string
)
```

- [ ] **Step 2: Register flags after existing flags**

```go
c.Flags().StringVar(&cloudInitHostname, "hostname", "", "Cloud-init hostname")
c.Flags().StringVar(&cloudInitPassword, "password", "", "Default user password (cloud-init)")
c.Flags().StringArrayVar(&cloudInitSSHKey, "ssh-key", nil, "SSH public key (repeatable, or multiline in file)")
c.Flags().StringVar(&cloudInitUserData, "user-data", "", "Cloud-init user_data script (path to file or literal string)")
c.Flags().StringArrayVar(&cloudInitNetworks, "network", nil, "Static IP config: nic=0,ip=192.168.1.100,netmask=255.255.255.0,gateway=192.168.1.1")
```

- [ ] **Step 3: Wire flags into service call**

In `CreateFromTemplate` branch (after line 66), before calling `service.NewVM(cli).CreateFromTemplate`:
```go
// Build CloudInit spec if any cloud-init flag is set
var cloudInit *adapter.CloudInitSpec
if cloudInitHostname != "" || cloudInitPassword != "" || len(cloudInitSSHKey) > 0 || cloudInitUserData != "" || len(cloudInitNetworks) > 0 {
    cloudInit = &adapter.CloudInitSpec{
        Hostname:            cloudInitHostname,
        DefaultUserPassword: cloudInitPassword,
        PublicKeys:          cloudInitSSHKey,
        UserData:            cloudInitUserData,
    }
    // Parse --network flags
    if len(cloudInitNetworks) > 0 {
        cloudInit.Networks = parseCloudInitNetworks(cloudInitNetworks)
    }
}
```

- [ ] **Step 4: Add parseCloudInitNetworks helper**

```go
// parseCloudInitNetworks parses --network flags into NicStaticConfig.
// Flag format: nic=0,ip=192.168.1.100,netmask=255.255.255.0,gateway=192.168.1.1
// Type is auto-detected: if ip/netmask/gateway provided → IPV4 (static); otherwise → IPV4_DHCP
func parseCloudInitNetworks(raw []string) []adapter.NicStaticConfig {
    out := make([]adapter.NicStaticConfig, 0, len(raw))
    for _, s := range raw {
        kv, err := parseKVList(s)
        if err != nil {
            continue // skip invalid, error already logged by parseKVList
        }
        cfg := adapter.NicStaticConfig{}
        if v, ok := kv["nic"]; ok {
            if n, err := strconv.ParseInt(v, 10, 32); err == nil {
                cfg.Index = int32(n)
            }
        }
        cfg.IP = kv["ip"]
        cfg.Netmask = kv["netmask"]
        cfg.Gateway = kv["gateway"]
        if cfg.IP != "" && cfg.Netmask != "" {
            cfg.Type = "IPV4"
        } else {
            cfg.Type = "IPV4_DHCP"
        }
        out = append(out, cfg)
    }
    return out
}
```

Need to add `strconv` to imports if not present.

- [ ] **Step 5: Pass CloudInit to service**

In `CreateFromTemplate` call, add `CloudInit: cloudInit` to the spec.

- [ ] **Step 6: Build and test --help**

Run: `go build -o goct . && ./goct vm.create --help`
Expected: New flags visible in help output

- [ ] **Step 7: Commit**

```bash
git add cmd/vm/create.go
git commit -m "feat(vm.create): add cloud-init flags --hostname --password --ssh-key --user-data --network"
```

---

## Task 4: Test cloud-init flow (smoke test)

**Files:**
- None (existing test infrastructure)

- [ ] **Step 1: Build and verify flags exist**

Run: `goct vm.create --help | grep -E "hostname|password|ssh-key|user-data|network"`

- [ ] **Step 2: Smoke test (dry run with wrong args)**

Run: `goct vm.create --from-template templ123 --cluster c1 --hostname test-vm --ssh-key "ssh-rsa AAAA..." --network nic=0,ip=192.168.1.100,netmask=255.255.255.0,gateway=192.168.1.1 2>&1 | head -5`
Expected: Error about missing VM or auth (not flag parse error)

- [ ] **Step 3: Commit test note**

```bash
git commit --allow-empty -m "test: cloud-init flags smoke test verified"
```

---

## Usage Examples (for documentation)

```bash
# Basic cloud-init with hostname and SSH key
goct vm.create --from-template my-template --name my-vm --cluster c1 \
    --hostname my-vm \
    --ssh-key "ssh-rsa AAAAB3NzaC1..." \
    --password "YourSecurePassword123"

# Static IP on first NIC
goct vm.create --from-template my-template --name my-vm --cluster c1 \
    --hostname my-vm \
    --ssh-key "ssh-rsa AAAAB3NzaC1..." \
    --network nic=0,ip=192.168.1.100,netmask=255.255.255.0,gateway=192.168.1.1

# DHCP (no static IP config needed)
goct vm.create --from-template my-template --name my-vm --cluster c1 \
    --network nic=0,type=dhcp

# Custom user_data script from file
goct vm.create --from-template my-template --name my-vm --cluster c1 \
    --user-data /path/to/cloud-init-config.yaml

# Multiple SSH keys
goct vm.create --from-template my-template --name my-vm --cluster c1 \
    --ssh-key "ssh-rsa AAAAB3NzaC1...@email1" \
    --ssh-key "ssh-rsa BBBB3NzaC1...@email2"
```
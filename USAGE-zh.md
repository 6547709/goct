# goct 使用指南

goct 是 SmartX CloudTower 的 govc 风格命令行工具。

## 公共参数

所有命令都支持以下参数：

```
  -h, --help              显示帮助信息
      --debug             启用调试日志 (env: GOCT_DEBUG)
      --dump              启用完整 HTTP 跟踪，不截断 body (env: GOCT_DUMP)
      --format string     输出格式：table|json (默认 "table")
      --insecure          跳过 TLS 证书验证 (env: GOCT_INSECURE)
      --password string   登录密码 (env: GOCT_PASSWORD)
      --source string     认证源：local|ldap|sso|authn (env: GOCT_SOURCE)
      --trace             启用 HTTP 跟踪 (env: GOCT_TRACE)
      --url string        CloudTower 端点 URL (env: GOCT_URL)
      --username string   登录用户名 (env: GOCT_USERNAME)
      --verbose           启用详细的 HTTP 跟踪，包含 headers 和 body (env: GOCT_VERBOSE)
      --cluster string    默认集群 ID 或名称 (env: GOCT_CLUSTER)
```

## 环境变量

| 变量 | 说明 |
|----------|-------------|
| GOCT_URL | CloudTower 端点 URL |
| GOCT_USERNAME | 登录用户名 |
| GOCT_PASSWORD | 登录密码 |
| GOCT_INSECURE | 跳过 TLS 验证 (true/false) |
| GOCT_SOURCE | 认证源：local\|ldap\|sso\|authn |
| GOCT_CLUSTER | 默认集群 ID 或名称 |
| GOCT_DEBUG | 启用调试日志 |
| GOCT_TRACE | 启用 HTTP 跟踪 |

<details><summary>目录</summary>

 - [系统命令](#系统命令)
 - [应用程序](#应用程序)
 - [内容库](#内容库)
 - [虚拟机](#虚拟机)
 - [主机](#主机)
 - [集群](#集群)
 - [存储](#存储)
 - [网络](#网络)
 - [任务](#任务)
 - [告警](#告警)
 - [用户](#用户)
 - [指标](#指标)
 - [其他命令](#其他命令)

</details>

## 系统命令

### about

```
goct about [flags]
```

显示 CloudTower 服务器版本和连接信息。

参数：
      -h, --help   help for about

---

### cluster-settings.get

```
goct cluster-settings.get [flags]
```

获取集群设置。

参数：
          --cluster string   集群 ID
      -h, --help             help for cluster-settings.get

---

### deploy.ls

```
goct deploy.ls [flags]
```

列出部署。

参数：
      -h, --help   help for deploy.ls

---

### find

```
goct find [flags]
```

跨清单查找资源（govc 风格）。

参数：
      --cluster string   按集群过滤
  -h, --help            help for find
      --id-only          只输出 ID
      --json             JSON 输出
      --limit int        最大结果数 (default 20)
      --name string      按名称过滤
      --type string      资源类型：m=虚拟机, h=主机, c=集群, d=存储, n=网络, v=VLAN, f=文件夹, g=放置组, t=模板, l=标签, u=用户, a=告警

---

### license.ls

```
goct license.ls [flags]
```

列出许可证。

参数：
      -h, --help   help for license.ls

---

### license.update-deploy

```
goct license.update-deploy [flags]
```

更新部署许可证。

参数：
      -h, --help   help for license.update-deploy

---

### ntp.get

```
goct ntp.get [flags]
```

获取 NTP 服务 URL。

参数：
      -h, --help   help for ntp.get

---

### session.login

```
goct session.login [flags]
```

强制登录并将会话令牌保存到本地缓存。

参数：
      -h, --help   help for session.login

---

### session.logout

```
goct session.logout [flags]
```

删除给定主机和用户的缓存会话令牌。

参数：
      -h, --help   help for session.logout

---

### session.ls

```
goct session.ls [flags]
```

列出本地缓存的会话文件。

参数：
      -h, --help   help for session.ls

---

### version

```
goct version [flags]
```

打印 goct 版本。

参数：
      -h, --help   help for version

---

## 应用程序

### delete-cloudtower-application-package

```
goct delete-cloudtower-application-package [flags]
```

删除 CloudTower 应用程序包。

参数：
      -h, --help   help for delete-cloudtower-application-package

---

### deploy-cloudtower-application

```
goct deploy-cloudtower-application [flags]
```

部署 CloudTower 应用程序。

参数：
      -h, --help   help for deploy-cloudtower-application

---

### get-cloudtower-application-packages

```
goct get-cloudtower-application-packages [flags]
```

列出 CloudTower 应用程序包。

参数：
      -h, --help   help for get-cloudtower-application-packages

---

### get-cloudtower-applications

```
goct get-cloudtower-applications [flags]
```

列出 CloudTower 应用程序。

参数：
      -h, --help   help for get-cloudtower-applications

---

### upload-cloudtower-application-package

```
goct upload-cloudtower-application-package [flags]
```

上传 CloudTower 应用程序包。

参数：
      -h, --help   help for upload-cloudtower-application-package

---

## 内容库

### content-library-image.delete

```
goct content-library-image.delete [flags]
```

删除内容库镜像。

参数：
      -h, --help   help for content-library-image.delete

---

### content-library-image.distribute

```
goct content-library-image.distribute [flags]
```

分发内容库镜像到集群。

参数：
      -h, --help   help for content-library-image.distribute

---

### content-library-image.import

```
goct content-library-image.import [flags]
```

导入内容库镜像。

参数：
      -h, --help   help for content-library-image.import

---

### content-library-image.ls

```
goct content-library-image.ls [flags]
```

列出内容库镜像。

参数：
      -h, --help   help for content-library-image.ls

---

## 虚拟机

### vm.cdrom.add

```
goct vm.cdrom.add [flags]
```

为 VM 添加 CD-ROM 驱动器。

参数：
      -h, --help   help for vm.cdrom.add

---

### vm.cdrom.eject

```
goct vm.cdrom.eject [flags]
```

弹出 CD-ROM 中的 ISO。

参数：
      -h, --help   help for vm.cdrom.eject

---

### vm.cdrom.ls

```
goct vm.cdrom.ls [flags]
```

列出 VM CD-ROM。

参数：
      -h, --help   help for vm.cdrom.ls

---

### vm.cdrom.rm

```
goct vm.cdrom.rm [flags]
```

移除 VM 的 CD-ROM 驱动器。

参数：
      -h, --help   help for vm.cdrom.rm

---

### vm.cdrom.toggle

```
goct vm.cdrom.toggle [flags]
```

启用或禁用 CD-ROM。

参数：
      -h, --help   help for vm.cdrom.toggle

---

### vm.clone

```
goct vm.clone [source-name|id] [flags]
```

克隆 VM。

参数：
          --cluster string   目标集群 ID（可选，默认相同）
      -h, --help             help for vm.clone
          --name string      克隆 VM 的名称（必填）

---

### vm.create

```
goct vm.create [flags]
```

创建新虚拟机。

参数：
          --cluster string         目标集群 ID
          --description string    描述
          --disk stringArray      磁盘规格，可重复：size=10g[,bus=SCSI][,name=diskN][,index=N][,boot=N][,iops=N]
          --dns stringArray       DNS 域名服务器（可重复，如 --dns 8.8.8.8）
          --firmware string       固件：BIOS 或 UEFI（空 = 使用模板默认值）
          --from-template string   从模板创建 VM 的模板 ID
          --full-copy             从模板创建时完整复制
          --ha string             启用 HA：true|false（空 = 使用集群默认值）
      -h, --help                   help for vm.create
          --hostname string        Cloud-init 主机名
          --memory int             内存 MiB（0 = 使用模板默认值）
          --name string            VM 名称
          --network stringArray    NIC 配置：nic=0[,ip=x][,netmask=x][,gateway=x][,route=10.0.0.0/8:192.168.1.1][,type=IPV4|DHCP]
          --nic stringArray        NIC 规格，可重复：vlan=<id|name>[,model=VIRTIO][,type=VLAN|VPC]
          --nic-model string       NIC 型号：E1000, SRIOV, VIRTIO（仅用于 --from-template）
          --nic-type string        NIC 类型：VLAN 或 VPC（仅用于 --from-template）
          --ssh-key string         Cloud-init SSH 公钥
          --user-data string       Cloud-init 用户数据
      -h, --help                   help for vm.create

---

### vm.destroy

```
goct vm.destroy [vm-name|id] [flags]
```

销毁（删除）VM。

参数：
      -h, --help   help for vm.destroy

---

### vm.disk.add

```
goct vm.disk.add [flags]
```

为 VM 添加磁盘。

参数：
      -h, --help   help for vm.disk.add

---

### vm.disk.expand

```
goct vm.disk.expand [flags]
```

扩展 VM 磁盘。

参数：
      -h, --help   help for vm.disk.expand

---

### vm.disk.ls

```
goct vm.disk.ls [vm-name|id] [flags]
```

列出 VM 磁盘。

参数：
      -h, --help   help for vm.disk.ls

---

### vm.disk.rm

```
goct vm.disk.rm [flags]
```

移除 VM 磁盘。

参数：
      -h, --help   help for vm.disk.rm

---

### vm.disk.update

```
goct vm.disk.update [flags]
```

更新 VM 磁盘设置。

参数：
          --bus string           磁盘总线：SCSI, SATA, NVMe
      -h, --help                 help for vm.disk.update
          --iops int             IOPS 限制（0 = 无限制）
          --name string          磁盘名称
          --size string          磁盘大小（如 10g）
          --vm string            VM 名称或 ID

---

### vm.export

```
goct vm.export [vm-name|id] [flags]
```

将 VM 导出为 OVF。

参数：
      -h, --help   help for vm.export

---

### vm.gpu.add

```
goct vm.gpu.add [flags]
```

为 VM 添加 GPU 设备。

参数：
      -h, --help   help for vm.gpu.add

---

### vm.gpu.ls

```
goct vm.gpu.ls [vm-name|id] [flags]
```

列出 VM GPU 设备。

参数：
      -h, --help   help for vm.gpu.ls

---

### vm.gpu.rm

```
goct vm.gpu.rm [flags]
```

移除 VM GPU 设备。

参数：
      -h, --help   help for vm.gpu.rm

---

### vm.info

```
goct vm.info [vm-name|id] [flags]
```

显示 VM 详情。

参数：
      -h, --help   help for vm.info

---

### vm.ip

```
goct vm.ip [vm-name|id] [flags]
```

显示 VM IP 地址（默认等待 VM tools 上报）。

参数：
  -a, --all            显示所有 IP（默认只显示第一个）
      --json           JSON 输出
  -h, --help           help for vm.ip
      --no-wait         不等待，立即返回
      --timeout int     等待超时秒数 (default 300)
      --v4              只显示 IPv4
      --v6              只显示 IPv6

---

### vm.ls

```
goct vm.ls [flags]
```

列出虚拟机。

参数：
      --id-only    只输出 ID
  -h, --help       help for vm.ls

---

### vm.migrate

```
goct vm.migrate [vm-name|id] [flags]
```

将 VM 迁移到另一台主机（省略 --host 让 CloudTower 选择）。

参数：
      --cluster string   目标集群 ID（可选）
  -h, --help            help for vm.migrate
      --host string      目标主机 ID（可选）

---

### vm.migrate.abort

```
goct vm.migrate.abort [vm-name|id] [flags]
```

取消跨集群迁移。

参数：
  -h, --help   help for vm.migrate.abort

---

### vm.migrate.across

```
goct vm.migrate.across [vm-name|id] [flags]
```

将 VM 迁移到另一个集群。

参数：
      --cluster string   目标集群 ID
      --host string      目标主机 ID（可选）
  -h, --help             help for vm.migrate.across

---

### vm.nic.add

```
goct vm.nic.add [flags]
```

为 VM 添加 NIC。

参数：
  -h, --help   help for vm.nic.add

---

### vm.nic.ls

```
goct vm.nic.ls [vm-name|id] [flags]
```

列出 VM NIC。

参数：
  -h, --help   help for vm.nic.ls

---

### vm.nic.rm

```
goct vm.nic.rm [flags]
```

按索引移除 VM NIC。

参数：
  -h, --help   help for vm.nic.rm

---

### vm.nic.update

```
goct vm.nic.update [flags]
```

更新 VM NIC 配置。

参数：
      --gateway string     网关
      --ip string          静态 IP
      --netmask string     子网掩码
  -h, --help               help for vm.nic.update
      --nic-id string      NIC ID
      --type string        类型：DHCP, IPV4
      --vlan string        VLAN ID 或名称
      --vm string          VM 名称或 ID

---

### vm.power.off

```
goct vm.power.off [vm-name|id] [flags]
```

关闭/断电 VM。

参数：
      --force    强制关闭
  -h, --help    help for vm.power.off

---

### vm.power.on

```
goct vm.power.on [vm-name|id] [flags]
```

开启 VM。

参数：
  -h, --help   help for vm.power.on

---

### vm.power.reset

```
goct vm.power.reset [vm-name|id] [flags]
```

重启/强制重置 VM。

参数：
      --force    强制重启
  -h, --help    help for vm.power.reset

---

### vm.power.resume

```
goct vm.power.resume [vm-name|id] [flags]
```

恢复挂起的 VM。

参数：
  -h, --help   help for vm.power.resume

---

### vm.power.suspend

```
goct vm.power.suspend [vm-name|id] [flags]
```

挂起 VM。

参数：
  -h, --help   help for vm.power.suspend

---

### vm.rebuild

```
goct vm.rebuild [vm-name|id] [flags]
```

从快照重建 VM。

参数：
  -h, --help   help for vm.rebuild

---

### vm.recover

```
goct vm.recover [vm-name|id] [flags]
```

从回收站恢复 VM。

参数：
  -h, --help   help for vm.recover

---

### vm.recycle

```
goct vm.recycle [vm-name|id] [flags]
```

将 VM 移至回收站。

参数：
  -h, --help   help for vm.recycle

---

### vm.reset-password

```
goct vm.reset-password [vm-name|id] [flags]
```

重置客户机操作系统密码。

参数：
      --password string   新密码
  -h, --help              help for vm.reset-password

---

### vm.shutdown

```
goct vm.shutdown [vm-name|id] [flags]
```

优雅关闭 VM（客户机操作系统关闭）。

参数：
      --force    强制关闭
  -h, --help     help for vm.shutdown

---

### vm.snapshot.create

```
goct vm.snapshot.create [vm-name|id] [flags]
```

创建快照。

参数：
      --name string    快照名称
  -h, --help          help for vm.snapshot.create

---

### vm.snapshot.ls

```
goct vm.snapshot.ls [vm-name|id] [flags]
```

列出 VM 快照。

参数：
  -h, --help   help for vm.snapshot.ls

---

### vm.snapshot.revert

```
goct vm.snapshot.revert [vm-name|id] [flags]
```

将 VM 回滚到快照。

参数：
  -h, --help   help for vm.snapshot.revert

---

### vm.snapshot.rm

```
goct vm.snapshot.rm [flags]
```

删除快照。

参数：
  -h, --help   help for vm.snapshot.rm

---

### vm.tools.install

```
goct vm.tools.install [vm-name|id] [flags]
```

在 VM 上安装 VMware Tools。

参数：
  -h, --help   help for vm.tools.install

---

### vm.update

```
goct vm.update [vm-name|id] [flags]
```

更新 VM 名称或描述。

参数：
      --description string   描述（空字符串清除）
  -h, --help                 help for vm.update
      --name string          新名称

---

### vm.vnc

```
goct vm.vnc [vm-name|id] [flags]
```

获取 VM VNC 连接信息。

参数：
  -h, --help   help for vm.vnc

---

## 主机

### host.disconnect

```
goct host.disconnect [host-name|id] [flags]
```

断开主机连接（SDK 不支持）。

参数：
  -h, --help   help for host.disconnect

---

### host.info

```
goct host.info [host-name|id] [flags]
```

显示主机详情。

参数：
  -h, --help   help for host.info

---

### host.ls

```
goct host.ls [flags]
```

列出主机。

参数：
      --id-only    只输出 ID
  -h, --help       help for host.ls

---

### host.maintenance.enter

```
goct host.maintenance.enter [host-name|id] [flags]
```

进入维护模式。

参数：
  -h, --help   help for host.maintenance.enter

---

### host.maintenance.exit

```
goct host.maintenance.exit [host-name|id] [flags]
```

退出维护模式。

参数：
  -h, --help   help for host.maintenance.exit

---

### host.reboot

```
goct host.reboot [host-name|id] [flags]
```

重启主机。

参数：
      --force    强制重启
  -h, --help     help for host.reboot

---

### host.reconnect

```
goct host.reconnect [host-name|id] [flags]
```

重新连接主机（SDK 不支持）。

参数：
  -h, --help   help for host.reconnect

---

### host.shutdown

```
goct host.shutdown [host-name|id] [flags]
```

关闭主机。

参数：
      --force    强制关闭
  -h, --help     help for host.shutdown

---

## 集群

### cluster.info

```
goct cluster.info [cluster-name|id] [flags]
```

显示集群详情。

参数：
  -h, --help   help for cluster.info

---

### cluster.ls

```
goct cluster.ls [flags]
```

列出集群。

参数：
      --id-only    只输出 ID
  -h, --help       help for cluster.ls

---

## 存储

### datastore.disk.ls

```
goct datastore.disk.ls [flags]
```

列出物理磁盘。

参数：
  -h, --help   help for datastore.disk.ls

---

### datastore.info

```
goct datastore.info [datastore-name|id] [flags]
```

显示存储详情。

参数：
  -h, --help   help for datastore.info

---

### datastore.ls

```
goct datastore.ls [flags]
```

列出存储。

参数：
      --id-only    只输出 ID
  -h, --help       help for datastore.ls

---

### storage.pool.ls

```
goct storage.pool.ls [flags]
```

列出超融合存储池（DiskPool）。

参数：
      --id-only    只输出 ID
  -h, --help       help for storage.pool.ls

---

## 网络

### network.info

```
goct network.info [network-name|id] [flags]
```

显示虚拟交换机详情。

参数：
  -h, --help   help for network.info

---

### network.ls

```
goct network.ls [flags]
```

列出虚拟交换机（VDS）。

参数：
      --id-only    只输出 ID
  -h, --help       help for network.ls

---

### vlan.create

```
goct vlan.create [flags]
```

创建 VLAN。

参数：
  -h, --help   help for vlan.create

---

### vlan.destroy

```
goct vlan.destroy [vlan-name|id] [flags]
```

删除 VLAN。

参数：
  -h, --help   help for vlan.destroy

---

### vlan.info

```
goct vlan.info [vlan-name|id] [flags]
```

显示 VLAN 详情。

参数：
  -h, --help   help for vlan.info

---

### vlan.ls

```
goct vlan.ls [flags]
```

列出 VLAN。

参数：
      --id-only    只输出 ID
  -h, --help       help for vlan.ls

---

## 任务

### events

```
goct events [flags]
```

显示最近事件/审计日志条目。

参数：
      --action string    按操作类型过滤（如 create_vm, delete_vm）
      --follow           持续跟踪新事件
      --json             JSON 输出
  -h, --help             help for events
      --limit int        最大结果数 (default 50)
      --user string      按用户过滤

---

### task.cancel

```
goct task.cancel [task-id] [flags]
```

取消任务（SDK 不支持）。

参数：
  -h, --help   help for task.cancel

---

### task.info

```
goct task.info [task-id] [flags]
```

显示任务详情。

参数：
  -h, --help   help for task.info

---

### task.ls

```
goct task.ls [flags]
```

列出任务。

参数：
      --id-only    只输出 ID
  -h, --help       help for task.ls

---

### task.wait

```
goct task.wait [task-id] [flags]
```

等待任务完成。

参数：
  -h, --help   help for task.wait

---

## 告警

### alert-rule.ls

```
goct alert-rule.ls [flags]
```

列出告警规则。

参数：
  -h, --help   help for alert-rule.ls

---

### alert.ack

```
goct alert.ack [alert-id] [flags]
```

确认（解决）告警。

参数：
  -h, --help   help for alert.ack

---

### alert.info

```
goct alert.info [alert-id] [flags]
```

显示告警详情。

参数：
  -h, --help   help for alert.info

---

### alert.ls

```
goct alert.ls [flags]
```

列出告警。

参数：
      --id-only    只输出 ID
  -h, --help       help for alert.ls

---

## 用户

### user.create

```
goct user.create [flags]
```

创建用户。

参数：
  -h, --help   help for user.create

---

### user.destroy

```
goct user.destroy [user-name|id] [flags]
```

删除用户。

参数：
  -h, --help   help for user.destroy

---

### user.info

```
goct user.info [user-name|id] [flags]
```

显示用户详情。

参数：
  -h, --help   help for user.info

---

### user.ls

```
goct user.ls [flags]
```

列出用户。

参数：
      --id-only    只输出 ID
  -h, --help       help for user.ls

---

## 指标

### cluster.metrics

```
goct cluster.metrics [cluster-name|id] [flags]
```

查询集群指标 (zbs_cluster_*)。

参数：
      --json       JSON 输出
  -h, --help       help for cluster.metrics
      --list        列出可用指标
      --range string   时间范围（如 5m, 1h, 7d）

---

### host.metrics

```
goct host.metrics [host-name|id] [flags]
```

查询主机指标。

参数：
      --json       JSON 输出
  -h, --help       help for host.metrics
      --list        列出可用指标
      --range string   时间范围（如 5m, 1h, 7d）

---

### sfs.metrics

```
goct sfs.metrics [flags]
```

查询 SFS 指标（TODO：未实现）。

参数：
  -h, --help   help for sfs.metrics

---

### vm.metrics

```
goct vm.metrics [vm-name|id] [flags]
```

查询 VM 指标 (elf_*)。

参数：
      --json       JSON 输出
  -h, --help       help for vm.metrics
      --list        列出可用指标
      --range string   时间范围（如 5m, 1h, 7d）

---

### vm.volume

```
goct vm.volume [vm-name|id] [flags]
```

查询 VM 卷指标 (elf_vm_disk_overall_*)。

参数：
      --json       JSON 输出
  -h, --help       help for vm.volume
      --list        列出可用指标
      --range string   时间范围（如 5m, 1h, 7d）

---

### volume.metrics

```
goct volume.metrics [flags]
```

查询卷指标。

参数：
      --json       JSON 输出
  -h, --help       help for volume.metrics
      --list        列出可用指标
      --range string   时间范围（如 5m, 1h, 7d）

---

## 其他命令

### completion

```
goct completion [shell] [flags]
```

生成指定 shell 的自动补全脚本。

参数：
  -h, --help   help for completion

---

### convert-to-vm

```
goct convert-to-vm [template-name|id] [flags]
```

将模板转换为 VM。

参数：
  -h, --help   help for convert-to-vm

---
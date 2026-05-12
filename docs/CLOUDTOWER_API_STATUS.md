# CloudTower API 实现状态

> 最后更新: 2026-05-12
> SDK: github.com/smartxworks/cloudtower-go-sdk/v2@v2.22.1

## 状态说明

- ✅ **已实现** - goct adapter 层已完成实现
- ⚠️ **部分实现** - 已有基础实现但功能不完整
- ❌ **未实现** - SDK 有但 goct 未实现
- ❓ **未知** - 需要进一步确认

---

## 完整 API 列表

### 计算资源

| SDK Client | 功能 | goct 状态 | 说明 |
|------------|------|-----------|------|
| **vm** | 虚拟机完整操作 | ⚠️ 部分 | 核心 VM 操作 |
| → vm.create | 创建 VM | ✅ | |
| → vm.clone | 克隆 VM | ⚠️ | 缺 linked clone (IsFullCopy) |
| → vm.destroy | 删除 VM | ✅ | |
| → vm.power | 电源操作 | ✅ | on/off/reset/suspend/resume |
| → vm.migrate | 迁移 | ✅ | 本集群迁移 |
| → vm.migrate_across_cluster | 跨集群迁移 | ✅ | |
| → vm.abort_migrate_across_cluster | 取消迁移 | ✅ | |
| → vm.reboot | 重启 | ✅ | |
| → vm.shutdown | 关闭 | ✅ | |
| → vm.rebuild | 从快照重建 | ✅ | |
| → vm.reset_password | 重置密码 | ✅ | |
| → vm.convert_to_vm | 模板转 VM | ✅ | |
| → vm.export | 导出 OVF | ✅ | |
| → vm.import | 导入 OVF | ❌ | **未实现** |
| → vm.rollback | 回滚快照 | ❌ | **未实现** |
| → vm.stop_in_cutover_migration | 停止迁移 | ❌ | **未实现** |
| → vm.add_pci_nic | 添加 PCI NIC | ❌ | **未实现** |
| → vm.remove_pci_nic | 移除 PCI NIC | ❌ | **未实现** |
| → vm.set_disk_resident_in_cache | 设置磁盘缓存 | ❌ | **未实现** |
| → vm.update_vm_advanced_options | 更新高级选项 | ❌ | **未实现** |
| → vm.update_vm_host_options | 更新主机选项 | ❌ | **未实现** |
| → vm.update_vm_io_policy | 更新 IO 策略 | ❌ | **未实现** |
| → vm.update_vm_nic_qos_option | 更新 NIC QoS | ❌ | **未实现** |
| → vm.update_vm_nic_vpc_info | 更新 VPC 信息 | ❌ | **未实现** |
| → vm.update_vm_owner | 更新所有者 | ❌ | **未实现** |
| → vm.add_vm_to_folder | 添加到文件夹 | ❌ | **未实现** |
| → vm.remove_vm_to_folder | 从文件夹移除 | ❌ | **未实现** |
| **host** | 主机操作 | ⚠️ 部分 | |
| → host.ls | 列表主机 | ✅ | |
| → host.info | 主机详情 | ✅ | |
| → host.register | 注册主机 | ✅ | CloudTower 专有 |
| → host.disconnect | 断开连接 | ✅ | |
| → host.reconnect | 重新连接 | ✅ | |
| → host.enter_maintenance | 进入维护 | ✅ | |
| → host.exit_maintenance | 退出维护 | ✅ | |
| → host.shutdown | 关机 | ✅ | |
| → host.reboot | 重启 | ✅ | |
| → host.add | 添加主机 | ❌ | **未实现** |
| → host.remove | 移除主机 | ❌ | **未实现** |
| → host.account.* | 主机账户管理 | ❌ | **未实现** |
| → host.cert.* | 证书管理 | ❌ | **未实现** |
| → host.esxcli | 执行 ESXCLI | ❌ | **未实现** |
| → host.firewall.* | 防火墙配置 | ❌ | **未实现** |
| → host.option.ls/set | 高级设置 | ❌ | **未实现** |
| → host.storage.* | 存储配置 | ❌ | **未实现** |
| → host.tpm.* | TPM 证明 | ❌ | **未实现** |
| → host.vnic.* | 虚拟 NIC | ❌ | **未实现** |
| → host.vswitch.* | vSwitch | ❌ | **未实现** |
| → host.service.* | 服务管理 | ❌ | **未实现** |
| → host.autostart.* | 自动启动 | ❌ | **未实现** |
| → host.date.* | 日期时间 | ❌ | **未实现** |
| → host.portgroup.* | 端口组 | ❌ | **未实现** |
| **cluster** | 集群操作 | ⚠️ 部分 | |
| → cluster.ls | 列表集群 | ✅ | |
| → cluster.info | 集群详情 | ✅ | |
| → cluster.register | 注册集群 | ✅ | CloudTower 专有 |
| → cluster.add | 添加主机到集群 | ❌ | **未实现** |
| → cluster.change | 修改集群 | ❌ | **未实现** |
| → cluster.module.* | 集群模块 | ❌ | **未实现** |
| → cluster.group.* | 集群组 | ❌ | **未实现** |
| → cluster.rule.* | HA/DRS 规则 | ❌ | **未实现** |
| → cluster.override.* | 覆盖设置 | ❌ | **未实现** |
| → cluster.stretch | 延伸集群 | ❌ | **未实现** |
| → cluster.vlcm.* | 生命周期管理 | ❌ | **未实现** |
| → cluster.upgrade_history | 升级历史 | ❌ | **未实现** |
| **vm_snapshot** | 快照操作 | ✅ | |
| → snapshot.create | 创建快照 | ✅ | |
| → snapshot.delete | 删除快照 | ✅ | |
| → snapshot.revert | 恢复到快照 | ✅ | |
| → snapshot.ls | 列表快照 | ✅ | |
| → snapshot.tree | 快照树视图 | ❌ | **未实现** |
| → snapshot.export | 导出快照 | ❌ | **未实现** |
| **vm_template** | VM 模板 | ❌ | **未实现** |
| **vm_folder** | VM 文件夹 | ❌ | **未实现** |
| **vm_placement_group** | VM 放置组 | ❌ | **未实现** |

### 存储资源

| SDK Client | 功能 | goct 状态 | 说明 |
|------------|------|-----------|------|
| **datastore** | 存储操作 | ⚠️ 部分 | |
| → datastore.ls | 列表存储 | ✅ | |
| → datastore.info | 存储详情 | ✅ | |
| → datastore.disk.ls | 磁盘列表 | ✅ | |
| → datastore.cluster.* | 存储集群 | ❌ | **未实现** |
| → datastore.vsan.* | VSAN 存储 | ❌ | **未实现** |
| → datastore.maintenance.* | 维护模式 | ❌ | **未实现** |
| → datastore.cp/mv | 文件复制/移动 | ❌ | **未实现** |
| → datastore.tail | 日志跟踪 | ❌ | **未实现** |
| → datastore.mkdir | 创建目录 | ❌ | **未实现** |
| **disk** | 磁盘操作 | ⚠️ 部分 | |
| → disk.ls | 列表磁盘 | ✅ | |
| → disk.attach | 热附加磁盘 | ❌ | **未实现** |
| → disk.detach | 热分离磁盘 | ❌ | **未实现** |
| → disk.register | 注册磁盘 | ❌ | **未实现** |
| → disk.snapshot.* | 磁盘快照 | ❌ | **未实现** |
| → disk.tags.* | 磁盘标签 | ❌ | **未实现** |
| → disk.metadata.* | 元数据 | ❌ | **未实现** |
| → disk.rdm.* | RDM 映射 | ❌ | **未实现** |
| **disk_pool** | 超融合存储池 | ✅ | |
| → storage.pool.ls | 列表存储池 | ✅ | |
| **elf_storage_policy** | 存储策略 | ❌ | **未实现** |
| **snapshot_plan** | 快照计划 | ❌ | **未实现** |
| **backup_plan** | 备份计划 | ❌ | **未实现** |
| **backup_restore** | 备份恢复 | ❌ | **未实现** |
| **iscsi_target** | iSCSI 目标 | ❌ | **未实现** |
| **iscsi_lun** | iSCSI LUN | ❌ | **未实现** |
| **nfs_export** | NFS 导出 | ❌ | **未实现** |
| **nvmf_namespace** | NVMe 命名空间 | ❌ | **未实现** |
| **nvmf_subsystem** | NVMe 子系统 | ❌ | **未实现** |

### 网络资源

| SDK Client | 功能 | goct 状态 | 说明 |
|------------|------|-----------|------|
| **network** | 网络操作 | ⚠️ 部分 | |
| → network.ls | 列表网络 | ✅ | |
| → network.info | 网络详情 | ✅ | |
| **vlan** | VLAN 管理 | ✅ | |
| → vlan.ls | 列表 VLAN | ✅ | |
| → vlan.create | 创建 VLAN | ✅ | |
| → vlan.delete | 删除 VLAN | ✅ | |
| → vlan.info | VLAN 详情 | ✅ | |
| **vds** | 分布式交换机 | ❌ | **未实现** |
| **security_group** | 安全组 | ❌ | **未实现** |
| **isolation_policy** | 隔离策略 | ❌ | **未实现** |
| **nic** | NIC 操作 | ❌ | **未实现** |
| **virtual_private_cloud*** | VPC 网络 | ❌ | **未实现** (15+ clients) |
| **network_policy_rule_service** | 网络策略规则 | ❌ | **未实现** |

### 虚拟机设备

| SDK Client | 功能 | goct 状态 | 说明 |
|------------|------|-----------|------|
| **vm_disk** | 磁盘管理 | ⚠️ 部分 | |
| → vm.disk.add | 添加磁盘 | ✅ | |
| → vm.disk.rm | 移除磁盘 | ✅ | |
| → vm.disk.expand | 扩展磁盘 | ✅ | |
| → vm.disk.update | 更新磁盘 | ✅ | |
| → vm.disk.ls | 列表磁盘 | ✅ | |
| **vm_nic** | NIC 管理 | ⚠️ 部分 | |
| → vm.nic.add | 添加 NIC | ✅ | |
| → vm.nic.rm | 移除 NIC | ✅ | |
| → vm.nic.update | 更新 NIC | ✅ | |
| → vm.nic.ls | 列表 NIC | ✅ | |
| **vm_cdrom** | CD-ROM 管理 | ⚠️ 部分 | |
| → vm.cdrom.add | 添加 CD-ROM | ✅ | |
| → vm.cdrom.rm | 移除 CD-ROM | ✅ | |
| → vm.cdrom.eject | 弹出 ISO | ✅ | |
| → vm.cdrom.toggle | 启用/禁用 | ✅ | |
| → vm.cdrom.ls | 列表 CD-ROM | ✅ | |
| **gpu_device** | GPU 设备 | ⚠️ 部分 | |
| → vm.gpu.add | 添加 GPU | ✅ | |
| → vm.gpu.rm | 移除 GPU | ✅ | |
| → vm.gpu.ls | 列表 GPU | ✅ | |
| **pci_device** | PCI 设备 | ❌ | **未实现** |
| **usb_device** | USB 设备 | ❌ | **未实现** |
| **pmem_dimm** | PMem 内存 | ❌ | **未实现** |

### 内容库

| SDK Client | 功能 | goct 状态 | 说明 |
|------------|------|-----------|------|
| **content_library_vm_template** | 内容库模板 | ⚠️ 部分 | |
| → template.ls | 列表模板 | ⚠️ | 仅 List，无 create/delete |
| → vm.create --from-template | 从模板创建 VM | ✅ | |
| **content_library_image** | 内容库镜像 | ❌ | **未实现** |
| **elf_image** | ELF 镜像 | ❌ | **未实现** |
| **svt_image** | SVT 镜像 | ❌ | **未实现** |
| **cluster_image** | 集群镜像 | ❌ | **未实现** |
| **upload_task** | 上传任务 | ❌ | **未实现** |
| **application** | 应用管理 | ❌ | **未实现** |
| **deploy** | 部署管理 | ❌ | **未实现** |

### 标签和分类

| SDK Client | 功能 | goct 状态 | 说明 |
|------------|------|-----------|------|
| **label** | 标签管理 | ❌ | **未实现** |
| → label.create | 创建标签 | ❌ | |
| → label.delete | 删除标签 | ❌ | |
| → label.update | 更新标签 | ❌ | |
| → label.get | 获取标签 | ❌ | |
| → label.add_to_resources | 添加到资源 | ❌ | |
| → label.remove_from_resources | 从资源移除 | ❌ | |
| **entity_filter** | 实体过滤器 | ❌ | **未实现** |
| **vm_entity_filter_result** | 过滤结果 | ❌ | **未实现** |

### 回收站

| SDK Client | 功能 | goct 状态 | 说明 |
|------------|------|-----------|------|
| **recycle_bin** | 回收站 | ⚠️ 部分 | |
| → vm.recycle | 移入回收站 | ✅ | |
| → vm.recover | 从回收站恢复 | ✅ | |
| → vm.ls (--recycle) | 回收站列表 | ❌ | **未实现** (需要过滤) |
| **global_settings** | 全局设置 | ❌ | **未实现** |
| → recycle_bin_setting | 回收站设置 | ❌ | **未实现** |
| → global_recycle_bin_setting | 全局回收站设置 | ❌ | **未实现** |

### 硬件拓扑

| SDK Client | 功能 | goct 状态 | 说明 |
|------------|------|-----------|------|
| **rack_topo** | 机架拓扑 | ❌ | **未实现** |
| **brick_topo** | Brick 拓扑 | ❌ | **未实现** |
| **node_topo** | 节点拓扑 | ❌ | **未实现** |
| **zone_topo** | 可用区拓扑 | ❌ | **未实现** |
| **cluster_topo** | 集群拓扑 | ❌ | **未实现** |
| **zone** | 可用区 | ❌ | **未实现** |
| **datacenter** | 数据中心 | ❌ | **未实现** |
| **business_host** | 业务主机 | ❌ | **未实现** |
| **business_host_group** | 业务主机组 | ❌ | **未实现** |
| **discovered_host** | 发现的主机 | ❌ | **未实现** |
| **vsphere_esxi_account** | vSphere ESXi 账户 | ❌ | **未实现** |

### 告警和监控

| SDK Client | 功能 | goct 状态 | 说明 |
|------------|------|-----------|------|
| **alert** | 告警 | ⚠️ 部分 | |
| → alert.ls | 列表告警 | ✅ | |
| → alert.info | 告警详情 | ✅ | |
| → alert.ack | 确认告警 | ✅ | |
| **alert_rule** | 告警规则 | ❌ | **未实现** |
| **alert_notifier** | 告警通知器 | ❌ | **未实现** |
| **global_alert_rule** | 全局告警规则 | ❌ | **未实现** |
| **metrics** | 指标查询 | ⚠️ 部分 | |
| → vm.metrics | VM 指标 | ✅ | |
| → host.metrics | 主机指标 | ✅ | |
| → cluster.metrics | 集群指标 | ✅ | |
| → volume.metrics | 卷指标 | ✅ | |
| → sfs.metrics | SFS 指标 | ⚠️ | SDK 未实现 |

### 用户和认证

| SDK Client | 功能 | goct 状态 | 说明 |
|------------|------|-----------|------|
| **user** | 用户管理 | ⚠️ 部分 | |
| → user.ls | 列表用户 | ✅ | |
| → user.info | 用户详情 | ✅ | |
| → user.create | 创建用户 | ✅ | |
| → user.delete | 删除用户 | ✅ | |
| **license** | 许可证 | ❌ | **未实现** |
| **organization** | 组织 | ❌ | **未实现** |
| **role** (user_role_next) | 角色管理 | ❌ | **未实现** |
| **user_role_next** | 用户角色 | ❌ | **未实现** |
| **login_client** | 登录客户端 | ❌ | **未实现** |

### 系统和配置

| SDK Client | 功能 | goct 状态 | 说明 |
|------------|------|-----------|------|
| **api_info** | API 信息 | ✅ | |
| → about | 服务器版本 | ✅ | |
| **task** | 任务 | ✅ | |
| → task.ls | 列表任务 | ✅ | |
| → task.info | 任务详情 | ✅ | |
| → task.wait | 等待任务 | ✅ | |
| → task.cancel | 取消任务 | ✅ | |
| **global_settings** | 全局设置 | ❌ | **未实现** |
| **cluster_settings** | 集群设置 | ❌ | **未实现** |
| **smtp_server** | SMTP 服务器 | ❌ | **未实现** |
| **snmp_transport** | SNMP 传输 | ❌ | **未实现** |
| **snmp_trap_receiver** | SNMP Trap 接收器 | ❌ | **未实现** |
| **ntp** | NTP 配置 | ❌ | **未实现** |
| **ipmi** | IPMI 配置 | ❌ | **未实现** |
| **registry_service** | 注册表服务 | ❌ | **未实现** |
| **vcenter_account** | vCenter 账户 | ❌ | **未实现** |

### 日志和审计

| SDK Client | 功能 | goct 状态 | 说明 |
|------------|------|-----------|------|
| **log_collection** | 日志收集 | ❌ | **未实现** |
| **log_service_config** | 日志服务配置 | ❌ | **未实现** |
| **system_audit_log** | 系统审计日志 | ❌ | **未实现** |
| **user_audit_log** | 用户审计日志 | ❌ | **未实现** |
| **table_reporter** | 表格报告 | ❌ | **未实现** |
| **report_task** | 报告任务 | ❌ | **未实现** |
| **report_template** | 报告模板 | ❌ | **未实现** |

### 虚拟化和平台

| SDK Client | 功能 | goct 状态 | 说明 |
|------------|------|-----------|------|
| **everoute_cluster** | Everoute 集群 | ❌ | **未实现** |
| **everoute_license** | Everoute 许可证 | ❌ | **未实现** |
| **everoute_package** | Everoute 包 | ❌ | **未实现** |
| **v2_everoute_license** | v2 许可证 | ❌ | **未实现** |
| **witness** | 见证 | ❌ | **未实现** |
| **witness_service** | 见证服务 | ❌ | **未实现** |
| **consistency_group** | 一致性组 | ❌ | **未实现** |
| **consistency_group_snapshot** | 一致性组快照 | ❌ | **未实现** |
| **snapshot_group** | 快照组 | ❌ | **未实现** |

### 复制和灾备

| SDK Client | 功能 | goct 状态 | 说明 |
|------------|------|-----------|------|
| **replication_plan** | 复制计划 | ❌ | **未实现** |
| **replica_vm** | 副本 VM | ❌ | **未实现** |
| **backup_service** | 备份服务 | ❌ | **未实现** |
| **backup_store_repository** | 备份存储库 | ❌ | **未实现** |
| **backup_target_execution** | 备份目标执行 | ❌ | **未实现** |
| **backup_plan_execution** | 备份计划执行 | ❌ | **未实现** |
| **backup_restore_execution** | 备份恢复执行 | ❌ | **未实现** |
| **backup_restore_point** | 备份恢复点 | ❌ | **未实现** |
| **ecp_license** | ECP 许可证 | ❌ | **未实现** |

### 其他

| SDK Client | 功能 | goct 状态 | 说明 |
|------------|------|-----------|------|
| **namespace_group** | 命名空间组 | ❌ | **未实现** |
| **resource_change** | 资源变更 | ❌ | **未实现** |
| **cluster_upgrade_history** | 集群升级历史 | ❌ | **未实现** |
| **graph** | 图形 | ❌ | **未实现** |
| **view** | 视图 | ❌ | **未实现** |
| **observability** | 可观测性 | ❌ | **未实现** |
| **cloud_tower_application** | CloudTower 应用 | ❌ | **未实现** |
| **cloud_tower_application_package** | 应用包 | ❌ | **未实现** |
| **upload_task** | 上传任务 | ❌ | **未实现** |

---

## 优先级建议

### 第一优先级（核心功能，SDK 已确认支持）

1. **回收站 VM 列表** - `vm.ls --recycle` 过滤 in_recycle_bin
2. **标签管理** - 完整的 label CRUD + 资源标签关联
3. **内容库模板管理** - template.ls, template.rm
4. **Linked Clone** - `vm.clone --linked` (IsFullCopy=false)
5. **Import OVF/VM** - 导入虚拟机

### 第二优先级（重要功能）

6. **Application 管理** - application, deploy, upload_task
7. **硬件拓扑** - rack_topo, brick_topo, node_topo, zone_topo
8. **快照计划** - snapshot_plan
9. **存储策略** - elf_storage_policy
10. **VM 文件夹** - vm_folder, vm_placement_group

### 第三优先级（增强功能）

11. **备份计划** - backup_plan, backup_restore
12. **PCI/USB 设备管理** - pci_device, usb_device
13. **角色管理** - user_role_next
14. **全局设置** - global_settings (含回收站设置)
15. **审计日志** - system_audit_log, user_audit_log

### 第四优先级（高级功能）

16. **VPC 网络** - virtual_private_cloud_* (需要额外部署)
17. **Everoute** - everoute_*
18. **复制/灾备** - replication_plan, replica_vm
19. **报告** - report_task, report_template
20. **NTP/SNMP/IPMI** - 系统配置

---

## 参考

- SDK 源码: `$GOPATH/pkg/mod/github.com/smartxworks/cloudtower-go-sdk/v2@v2.22.1/`
- goct adapter: `/Users/liguoqiang/project/goct/pkg/adapter/`
- goct 命令: `/Users/liguoqiang/project/goct/cmd/`

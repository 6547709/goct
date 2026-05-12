---
description: goct CLI tool for CloudTower management
---

# CloudTower Credentials

- **API URL**: https://10.0.50.210:443
- **Username**: root
- **Password**: VMware1!
- **Source**: LOCAL (认证源)

```bash
# 测试连接
GOCT_INSECURE=true ./goct metrics vm.metrics elf_cpu_usage --range 5m
```

## 注意事项

- demo02 VM 没有 metrics 数据（points: null），但 API 返回正常
- 需要用 `--source local` 指定认证源

# 外部参考资料目录

本目录用于存放 **不纳入主仓库** 的 CloudTower 官方文档与 SDK 资料。
被 `.gitignore` 排除，仅作开发期本地查阅。

## 目录约定

```
docs/references/
├── README.md                       # 本文件（入仓库）
└── cloudtower-skills/              # SmartX 官方 OpenAPI Skill 文档（不入仓库）
    └── skills/cloudtower-api/
        ├── SKILL.md                # Skill 入口
        └── references/
            ├── authentication.md   # 认证流程
            ├── operations/         # 573 个 API 操作详细文档
            ├── resources/          # 130 个资源类型文档
            └── schemas/            # ~1900 个 schema 定义
```

## 拉取/更新方式

```bash
cd docs/references
git clone --depth 1 https://github.com/smartxworks/cloudtower-skills.git
rm -rf cloudtower-skills/.git cloudtower-skills/.github
```

> **注**：该仓库存在大小写冲突路径（`schemas/ROLE` vs `schemas/Role`），在 macOS 默认文件系统（不区分大小写）下会被合并；Linux 会得到完整两份。仅本地查阅，不影响 goct 编译。

## 与 cloudtower-go-sdk 的关系

| 资料 | 用途 | 信任级别 |
| --- | --- | --- |
| `docs/plans/2026-04-29-sdk-cheatsheet.md` | adapter 层编码弹药（SDK 真实方法签名） | **★★★ 首要** |
| `docs/references/cloudtower-skills/.../operations/<Op>.md` | 接口语义、请求/响应字段细节、错误码 | **★★ 补充** |
| `docs/references/cloudtower-skills/.../authentication.md` | 登录流程权威说明（含 UserSource 枚举默认值） | **★★ 补充** |
| `~/go/pkg/mod/.../cloudtower-go-sdk/v2@v2.22.1/` | SDK 源码（生成式，与 OpenAPI 同源） | **★★★ 首要** |

**优先级原则**：
1. **adapter 编码** → 先看 cheatsheet → 再看 SDK 源码
2. **接口语义模糊** → 查 cloudtower-skills 的 operations/<Op>.md
3. **二者冲突** → 以 SDK 源码为准（编译期就能验证）

## 已确认的 OpenAPI 文档关键洞察

来自 cloudtower-skills 的几条 SDK README 未强调的重要信息：

1. **同步 vs 异步**：写操作返回 `WithTask_Xxx_`，若 `task_id == null` 表示该次调用是同步完成的，service 层不应再调 `task.Watch`。
2. **写操作响应数据语义**：成功响应不代表操作完成；返回 `data` 中除 `id` 外的字段应视为临时数据，业务侧需在 task 完成后重新 GET 拿最新状态。
3. **认证 source 默认值**：`LoginInput.source` 为枚举 `AUTHN/LDAP/LOCAL/SSO`，未指定时使用 `LOCAL`。
4. **同一资源多个变体方法**：例如 `restart-vm`（优雅）vs `force-restart-vm`（强制）；`shutdown-vm`（优雅）vs `poweroff-vm`（强制）；adapter 应区分暴露而非合并。

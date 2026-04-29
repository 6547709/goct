这份 SPEC.md 专为 Claude Code 编写，采用了高度结构化的中文描述，旨在清晰地传达架构设计意图、工程规范以及应对 SDK 升级的策略。
你可以直接将以下内容保存为项目根目录下的 SPEC.md，然后对 Claude Code 说：“请阅读 SPEC.md，并按照其中的架构设计和步骤开始开发 goct 工具。”
项目技术规格书：goct (CloudTower CLI)
1. 项目愿景
项目名称： goct (CloudTower CLI)
设计目标： 模仿 VMware 的 govc 工具，为 SmartX CloudTower 提供一个简洁、无状态、高可组合性的命令行界面。
核心原则： - 极简主义： 模仿 govc 的动宾结构（如 vm.ls, vm.info）。
• 可维护性： 采用“防腐层”设计，确保 CloudTower SDK 升级时，CLI 逻辑改动最小。
• 自动化友好： 默认支持表格输出，同时提供完整的 JSON 支持以便于脚本集成。
2. 技术栈要求
• 编程语言： Go 1.21+
• CLI 框架： spf13/cobra (核心框架)
• 配置管理： spf13/viper (处理环境变量和配置文件)
• SDK： github.com/smartxworks/cloudtower-go-sdk
• UI 组件： olekukonko/tablewriter (用于表格格式化输出)
3. 架构设计 (三层架构)
为了确保 SDK 更新时不破坏整体代码，必须严格遵守以下分层：
1. 命令层 (/cmd)： • 仅负责 Cobra 命令的定义、参数 (Flags) 解析。 • 不包含具体的业务逻辑或直接的 SDK 调用。
2. 服务层 (/pkg/service)： • 编排业务逻辑（例如：根据名称查找 VM ID，然后执行关机）。 • 定义稳定的接口，屏蔽 SDK 的复杂性。
3. 适配器层/防腐层 (/pkg/adapter)： • 核心关键点：所有对 cloudtower-go-sdk 的直接引用必须限制在此包内。 • 负责将 SDK 的模型转换为 CLI 的内部模型。 • 如果 SDK 升级导致函数签名变更，仅需在此处进行适配。
4. 关键特性实现指南
A. 认证与持久化
• 优先级：命令行 Flag > 环境变量 > 配置文件 (~/.goct.yaml)。
• 环境变量：GOCT_URL,GOCT_USERNAME,GOCT_PASSWORD,GOCT_INSECURE,GOCT_CLUSTER
• Session 缓存：参考 govc，将有效的 Session Token 缓存到本地临时文件，避免频繁登录。
B. Flag 注入模式 (Flag Embedding)
• 严禁在每个子命令中重复定义 --endpoint 等参数。
• 创建 pkg/flags 包，定义通用的结构体（如 ConnectionFlags, OutputFlags, SearchFlags）。
• 子命令通过嵌入这些结构体来实现参数复用。
C. 异步任务观察器 (Task Watcher)
• CloudTower 的操作多为异步。实现一个通用的 WaitTask(taskID string) 函数。
• 当执行变更操作（如创建、关机）时，自动轮询 Task 状态，并显示进度条或 Spinner，直到任务完成。
D. 统一输出引擎
• 支持 -format=table (默认) 和 -json。
• -json 模式下，直接输出 SDK 返回的原始 Data 结构，确保字段完整性。
5. 目录结构推荐
goct/
├── cmd/                # 命令定义
│   ├── root.go         # 根命令与全局配置初始化
│   ├── vm/             # 虚拟机相关子命令 (ls, info, power)
│   └── cluster/        # 集群相关子命令
├── pkg/
│   ├── adapter/        # SDK 防腐层 (唯一允许 import SDK 的地方)
│   ├── client/         # SDK 客户端初始化与认证逻辑
│   ├── output/         # 统一格式化输出逻辑
│   └── task/           # Task 状态轮询逻辑
├── main.go             # 程序入口
└── SPEC.md             # 本规范文档

6. 开发路径 (给 Claude Code 的指令序列)
第一阶段：基础设施
1. 初始化项目结构，安装 cobra 和 viper。
2. 实现 root.go，处理基础认证 Flags 和环境变量。
3. 实现 pkg/client，确保能成功连接到 CloudTower 并获取版本信息。
第二阶段：只读命令与输出
1. 实现 pkg/output，支持表格和 JSON 渲染。
2. 实现 goct vm ls 命令：调用 pkg/adapter 获取虚拟机列表。
3. 实现 goct vm info <name> 命令：支持通过名称或 ID 查询。
第三阶段：变更命令与任务追踪
1. 实现 pkg/task 中的轮询逻辑。
2. 实现 goct vm power.on/off 命令，并在命令结束后自动跟踪生成的 Task。
7. 编码规范
• 错误处理：必须使用 fmt.Errorf("context: %w", err) 包装错误，严禁丢弃错误。
• 命名风格：子命令使用小写短横线或点分隔，例如 vm.ls 或 vm-ls（推荐遵循 govc 的 . 风格）。
• 注释：每个非平凡的函数必须有中文注释说明其意图。

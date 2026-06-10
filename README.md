# Kratos Bootstrap

## 项目亮点

- **配置契约优先**：全部配置通过 Protobuf 定义，类型安全、强约束、自动生成

## 技术栈

| 层级 | 技术 | 说明 |
| --- | --- | --- |
| 配置定义 | Protobuf + buf.build | 接口契约优先，类型安全 |
| 多模块管理 | Go Workspace | 模块化设计，按需引入 |

## 项目结构

```
kratos-bootstrap/
├── api/                                # Protobuf API 定义与生成代码
│   ├── protos/conf/v1/                 # .proto 源文件（配置结构定义）
│   └── gen/go/conf/v1/                 # buf 生成的 Go 代码
├── bootstrap/                          # 应用引导核心
│   ├── app_info.go                     # 应用元信息管理
│   ├── bootstrap.go                    # 引导入口（配置加载/日志/注册/追踪）
│   ├── cli.go                          # CLI 命令行框架（Cobra）
│   ├── context.go                      # 引导上下文
├── config/                             # 配置加载模块
│   ├── config.go                       # 配置提供者（本地/远程）
│   ├── factory.go                      # 配置源工厂
├── server/                             # 服务模块
│   ├── middleware/                     # 中间件
│   ├── rest.go                         # REST/HTTP 服务端
│   ├── white_list.go                   # 白名单管理
├── example/                            # 示例模块
├── docs/                               # 文档
│   ├── go-work-guidelines.md           # Go Workspace 使用指南
│   └── modules.md                      # 多模块依赖管理
├── go.work                             # Go Workspace 配置
└── Makefile                            # 构建与代码生成命令
```

## 核心功能

### 应用引导

| 功能 | 说明 |
| --- | --- |
| 统一启动入口 | 通过 `RunApp` 函数封装配置加载、日志初始化、服务注册、链路追踪全流程 |
| CLI 框架 | 基于 Cobra 的命令行框架，支持子命令定制与 Flag 注入 |
| 守护进程 | 原生守护进程模式，支持后台运行与 PID 管理 |
| 应用元信息 | 统一管理应用名称、版本号、实例 ID、项目空间等元数据 |
| 优雅退出 | 内置信号捕获与优雅关停机制，确保服务安全下线 |

### 配置管理

| 功能 | 说明 |
| --- | --- |
| 配置提供者 | 支持本地文件配置，可扩展支持 Apollo/Consul/Etcd/Kubernetes/Nacos/Polaris |
| 工厂模式 | 配置源工厂注册机制，按需加载不同配置源 |
| 配置扫描 | 自动扫描并加载多个配置文件 |

### API 契约

| 功能 | 说明 |
| --- | --- |
| Protobuf 定义 | 应用元信息、服务器配置、TLS 配置、中间件配置、配置源定义 |
| 代码生成 | buf 自动生成 Go/HTTP/errors/TypeScript 代码 |
| 类型安全 | 配置字段强类型约束，编译期检查 |

### 服务通信

| 功能 | 说明 |
| --- | --- |
| REST (HTTP) 服务端 | HTTP/RESTful API 服务，支持 CORS、pprof 调试 |
| 中间件链 | Recovery / Tracing / Validate / RateLimit / Metadata |
| 限流 | BBR（Google BBR 算法）自适应限流 |
| 白名单 | 基于方法名的中间件白名单机制 |

## 设计原则

- **配置驱动**：全部组件通过 Protobuf 配置定义初始化，类型安全，避免硬编码
- **非侵入封装**：封装层仅做配置转译与实例创建，不侵入框架原生 API，保留框架原生使用方式
- **工厂模式**：注册中心、日志、追踪等组件均采用工厂注册模式，按需加载，松耦合
- **模块化设计**：每个功能模块独立 Go Module，按需引入，不产生冗余依赖
- **防御性编程**：全面的空指针检查与错误处理，确保组件缺失时不影响主流程

## 适用场景

- 微服务架构下的标准化应用引导与基础设施管理

## 许可证

项目基于 MIT 许可证开源，允许自由使用、修改和分发，需保留原版权信息。

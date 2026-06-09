# Go Work 项目准则

本文档描述本项目如何使用 Go Work 功能管理多模块项目。

## 项目结构

```
kratos-bootstrap/
├── go.mod              # 根模块（空模块，仅用于版本管理）
├── go.work             # Go Work 配置文件
├── go.work.sum         # Work 依赖校验和
├── api/                # API 模块（Protocol Buffers 定义）
│   ├── go.mod
│   └── gen/            # 生成的 Go 代码
└── bootstrap/          # Bootstrap 模块（服务引导框架）
    ├── go.mod
    ├── bootstrap.go
    ├── context.go
    ├── cli.go
    └── bootstrap_test.go
```

## Go Work 配置

### go.work 文件

```go
go 1.26.1

use (
    .
    ./api
    ./bootstrap
)
```

### 配置说明

| 配置项 | 说明 |
|--------|------|
| `go 1.26.1` | 项目使用的 Go 版本 |
| `use .` | 根目录模块 |
| `use ./api` | API 模块 |
| `use ./bootstrap` | Bootstrap 模块 |

## 使用准则

### 1. 模块依赖关系

```
bootstrap (主模块)
    └── api (依赖模块)
```

- `bootstrap` 模块依赖 `api` 模块
- `api` 模块在 `bootstrap/go.mod` 中通过 replace 或直接引用

### 2. 本地开发

使用 `go work` 命令进行本地开发：

```bash
# 在项目根目录执行
cd /home/murphy/workspace/golang/src/github.com/kalandramo/kratos-bootstrap

# 运行所有模块的测试
go test ./...

# 运行特定模块的测试
go test ./bootstrap/...

# 构建所有模块
go build ./...

# 构建特定模块
go build ./bootstrap
```

### 3. 添加新模块

当需要添加新模块时，编辑 `go.work` 文件：

```bash
# 1. 创建新模块目录
mkdir ./newmodule
go mod init github.com/kalandramo/kratos-bootstrap/newmodule

# 2. 添加到 go.work
go work use ./newmodule

# 或直接编辑 go.work 文件
```

### 4. 模块间依赖

在 `bootstrap/go.mod` 中引用 `api` 模块：

```go
require github.com/kalandramo/kratos-bootstrap/api v0.0.0

// 本地开发时自动通过 go.work 解析
// 发布时需要正确配置版本号
```

### 5. 同步依赖

```bash
# 同步所有模块的依赖
go work sync

# 更新 go.work.sum
go work use -r
```

## 命令参考

| 命令 | 说明 |
|------|------|
| `go work use <dir>` | 添加模块到 work 空间 |
| `go work use -r` | 递归查找并添加所有模块 |
| `go work sync` | 同步模块依赖 |
| `go work edit` | 编辑 go.work 文件 |
| `go work why` | 查看模块为何被包含 |

## 最佳实践

1. **统一 Go 版本**：所有模块使用相同的 Go 版本（1.26.1）

2. **模块职责清晰**：
   - `api`：仅包含 Protocol Buffers 定义和生成的代码
   - `bootstrap`：服务引导框架核心逻辑

3. **依赖管理**：
   - 避免循环依赖
   - 使用 `go work sync` 同步依赖版本

4. **CI/CD**：
   - 确保 CI 环境中安装 Go 1.26.1
   - 提交前运行 `go work sync`

## 常见问题

### Q: 为什么需要 go.work？

A: 当项目包含多个 Go 模块时，`go.work` 允许：
- 在单一工作空间中开发多个模块
- 自动解析模块间的本地依赖
- 避免手动配置 replace 指令

### Q: go.mod 和 go.work 的区别？

A:
- `go.mod`：定义单个模块的依赖
- `go.work`：定义多个模块的工作空间关系

### Q: 如何排除某个模块？

A: 从 `go.work` 中移除对应的 `use` 行：
```bash
go work use -r  # 重新扫描
# 或手动编辑 go.work
```

## 参考文档

- [Go Workspaces](https://go.dev/ref/mod#workspaces)
- [go work 命令文档](https://pkg.go.dev/cmd/go#hdr-Workspaces)

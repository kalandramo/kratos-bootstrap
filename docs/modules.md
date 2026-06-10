# 多模块依赖管理

## 项目结构

本项目采用 Go Workspace 模式管理多个模块：

```
kratos-bootstrap/
├── api/           # API 定义模块 (proto + 生成的 Go 代码)
├── bootstrap/     # 引导模块 (依赖 api)
├── config/        # 配置模块 (依赖 api)
├── server/        # 服务模块 (依赖 api)
├── example/       # 示例模块 (依赖 api, bootstrap, server)
└── go.work        # Go Workspace 配置文件
```

## 依赖关系

```
api (基础模块，无内部依赖)
  ├── bootstrap
  ├── config
  └── server
        └── example
```

## 工作原理

### Go Workspace

`go.work` 文件将多个本地模块组合成一个工作区，使得模块间可以使用本地版本而非远程版本：

```go
go 1.26.1

use (
    .
    ./api
    ./bootstrap
    ./config
    ./server
    ./example
)
```

### 依赖更新流程

当某个模块（如 `api`）更新后，其他模块需要通过以下步骤更新依赖：

1. **`go work use`** - 将模块添加到 workspace
2. **`go work sync`** - 同步所有模块的依赖版本
3. **`go mod tidy`** - 整理各模块的依赖

## 使用方式

### 更新所有模块依赖

当 `api` 等基础模块更新后，运行：

```bash
make update-all
```

这会：
1. 将所有模块添加到 workspace
2. 同步 workspace 依赖
3. 各模块自动使用 workspace 中的最新本地版本

### 整理所有依赖

```bash
make tidy
```

这会：
1. 运行 `go work sync` 同步 workspace
2. 对每个模块运行 `go mod tidy` 清理未使用的依赖

### 生成代码

```bash
make gen    # 或 make api
```

生成 protobuf 代码。

## 验证

更新后验证构建：

```bash
go build ./api/... ./bootstrap/... ./server/... ./config/... ./example/...
```

## 常见问题

### Q: 为什么不用 `go get`？

A: 在 Workspace 模式下，`go work sync` 会自动同步所有模块使用 workspace 中定义的本地版本，无需手动 `go get`。

### Q: 提交时需要提交什么？

A: 
- `go.work` - workspace 配置
- 各模块的 `go.mod` / `go.sum` - 依赖信息
- 生成的代码（如 `api/gen/go/`）

### Q: 如何添加新模块？

A: 
1. 创建新模块目录和 `go.mod`
2. 运行 `go work use ./新模块`
3. 运行 `make tidy`

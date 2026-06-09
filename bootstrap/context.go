package bootstrap

import (
	"context"
	"sync"

	kratosRegistry "github.com/go-kratos/kratos/v3/registry"
	conf "github.com/kalandramo/kratos-bootstrap/api/gen/go/conf/v1"
)

// Context 引导上下文
type Context struct {
	// config  *conf.Bootstrap // 引导配置
	appInfo *conf.AppInfo // 应用信息

	registrar kratosRegistry.Registrar // 服务注册器

	customConfig sync.Map // 自定义配置项
	values       sync.Map // 自定义值存储

	rootCtx context.Context    // 应用级根上下文（可用于优雅关闭）
	cancel  context.CancelFunc // 取消函数
}

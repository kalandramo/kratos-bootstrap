package bootstrap

import (
	"fmt"

	"github.com/go-kratos/kratos/v3"
	"github.com/go-kratos/kratos/v3/transport"
	"github.com/spf13/cobra"

	"github.com/kalandramo/kratos-bootstrap/config"
)

// NewApp 创建应用程序
func NewApp(ctx *Context, srv ...transport.Server) *kratos.App {
	var opts []kratos.Option

	if ctx.registrar != nil {
		opts = append(opts, kratos.Registrar(ctx.registrar))
	}

	if len(srv) > 0 {
		opts = append(opts, kratos.Server(srv...))
	}

	if ctx.appInfo.Metadata != nil {
		opts = append(opts, kratos.Metadata(ctx.appInfo.Metadata))
	}

	if ctx.appInfo.AppId != "" {
		registerName := ctx.appInfo.Project + "/" + ctx.appInfo.AppId
		opts = append(opts, kratos.Name(registerName))
	}

	if ctx.appInfo.Version != "" {
		opts = append(opts, kratos.Version(ctx.appInfo.Version))
	}

	if ctx.appInfo.InstanceId != "" {
		opts = append(opts, kratos.ID(ctx.appInfo.InstanceId))
	}

	return kratos.New(opts...)
}

// InitAppFunc 应用初始化函数类型
type InitAppFunc func(ctx *Context) (app *kratos.App, cleanup func(), err error)

// RunApp 运行应用程序并允许在执行前对 root 命令做定制。
// opts 可用于注册子命令、对 root 添加 flag 或其他修改。
func RunApp(ctx *Context, initApp InitAppFunc, opts ...func(root *cobra.Command)) error {
	if ctx == nil {
		return fmt.Errorf("bootstrap context is nil")
	}

	// 注入命令行参数
	root := NewRootCmd(flags, func(cmd *cobra.Command, args []string) error {
		return bootstrap(ctx, initApp)
	})

	// 允许调用方定制 root（如添加子命令、注册额外 flag 等）
	for _, opt := range opts {
		if opt != nil {
			opt(root)
		}
	}

	// 如果 flags 实现了 Register，就在 Execute 前注册到命令上，确保 cobra 能解析这些 flag
	if rb, ok := interface{}(flags).(interface{ Register(cmd *cobra.Command) }); ok {
		rb.Register(root)
	}

	if err := root.Execute(); err != nil {
		return err
	}

	return nil
}

// bootstrap 应用引导启动
func bootstrap(ctx *Context, initApp InitAppFunc) error {
	// 打印应用信息
	ctx.PrintAppInfo()

	var err error

	// load configs
	if err = config.LoadBootstrapConfig(flags.Conf); err != nil {
		return err
	}

	// get bootstrap config
	ctx.config = config.GetBootstrapConfig()
	if ctx.config == nil {
		return fmt.Errorf("bootstrap config is nil")
	}

	// init app
	app, cleanup, err := initApp(ctx)
	if err != nil {
		return err
	}
	defer cleanup()

	// run the app.
	if err = app.Run(); err != nil {
		return err
	}

	return nil
}

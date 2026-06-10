package main

import (
	"context"

	"github.com/go-kratos/kratos/v3"
	"github.com/go-kratos/kratos/v3/middleware"
	"github.com/go-kratos/kratos/v3/transport/http"

	confv1 "github.com/kalandramo/kratos-bootstrap/api/gen/go/conf/v1"
	"github.com/kalandramo/kratos-bootstrap/bootstrap"
	"github.com/kalandramo/kratos-bootstrap/server"
)

const (
	ProjectName = "example"
	AppId       = "base-test-service" // 后台服务
)

var version = "v1.0.0"

// NewDiscoveryName 构建服务发现名称
func NewDiscoveryName(serviceName string) string {
	return ProjectName + "/" + serviceName
}

func main() {
	if err := runApp(); err != nil {
		panic(err)
	}
}

func runApp() error {
	ctx := bootstrap.NewContext(
		context.Background(),
		&confv1.AppInfo{
			Project: ProjectName,
			AppId:   AppId,
			Version: version,
		},
	)

	return bootstrap.RunApp(ctx, initApp)
}

func initApp(ctx *bootstrap.Context) (*kratos.App, func(), error) {
	hs, err := NewRestServer(ctx, nil)
	if err != nil {
		return nil, nil, err
	}

	app := newApp(ctx, hs)
	return app, func() {}, nil
}

func newApp(
	ctx *bootstrap.Context,
	hs *http.Server,
) *kratos.App {
	return bootstrap.NewApp(ctx,
		hs,
	)
}

// NewRestServer new an REST server.
func NewRestServer(ctx *bootstrap.Context, middlewares []middleware.Middleware) (*http.Server, error) {
	cfg := ctx.GetConfig()

	if cfg == nil || cfg.Server == nil || cfg.Server.Rest == nil {
		return nil, nil
	}

	srv, err := server.CreateRestServer(cfg,
		middlewares...,
	)
	if err != nil {
		return nil, err
	}

	return srv, nil
}

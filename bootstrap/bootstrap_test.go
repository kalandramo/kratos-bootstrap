package bootstrap

import (
	"context"
	"testing"

	"github.com/go-kratos/kratos/v3"
	confv1 "github.com/kalandramo/kratos-bootstrap/api/gen/go/conf/v1"
	"github.com/stretchr/testify/assert"
)

func initApp(ctx *Context) (*kratos.App, func(), error) {
	app := NewApp(ctx)
	return app, func() {}, nil
}

func TestNewApp_WithConfig(t *testing.T) {
	tests := []struct {
		name     string
		appInfo  *confv1.AppInfo
		validate func(t *testing.T, app *kratos.App)
	}{
		{
			name: "with name and version",
			appInfo: &confv1.AppInfo{
				AppId:   "test-service",
				Version: "v1.0.0",
			},
			validate: func(t *testing.T, app *kratos.App) {
				assert.NotNil(t, app)
			},
		},
		{
			name: "with metadata",
			appInfo: &confv1.AppInfo{
				AppId:    "test-service",
				Version:  "v1.0.0",
				Metadata: map[string]string{"env": "test"},
			},
			validate: func(t *testing.T, app *kratos.App) {
				assert.NotNil(t, app)
			},
		},
		{
			name: "with custom instance id",
			appInfo: &confv1.AppInfo{
				AppId:      "test-service",
				Version:    "v1.0.0",
				InstanceId: "custom-instance-id",
			},
			validate: func(t *testing.T, app *kratos.App) {
				assert.NotNil(t, app)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := NewContext(context.Background(), tt.appInfo)
			app := NewApp(ctx)
			tt.validate(t, app)
		})
	}
}

func TestInitAppFunc(t *testing.T) {
	ctx := NewContext(context.Background(), &confv1.AppInfo{
		AppId:   "test",
		Version: "v1.0.0",
	})

	app, cleanup, err := initApp(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, app)
	assert.NotNil(t, cleanup)

	// cleanup 应该幂等
	cleanup()
}

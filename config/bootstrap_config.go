package config

import (
	"reflect"
	"sync"

	confv1 "github.com/kalandramo/kratos-bootstrap/api/gen/go/conf/v1"
	"google.golang.org/protobuf/proto"
)

var (
	muBC         sync.RWMutex
	initOnce     sync.Once
	configList   []proto.Message
	configSet    map[uintptr]struct{}
	commonConfig *confv1.Bootstrap
)

func GetBootstrapConfig() *confv1.Bootstrap {
	initBootstrapConfig()
	muBC.RLock()
	defer muBC.RUnlock()
	return commonConfig
}

// initBootstrapConfig 初始化引导配置（仅执行一次）
func initBootstrapConfig() {
	initOnce.Do(func() {
		muBC.Lock()
		defer muBC.Unlock()

		// 初始化集合与列表
		configList = make([]proto.Message, 0)
		configSet = make(map[uintptr]struct{})

		if commonConfig == nil {
			commonConfig = &confv1.Bootstrap{}
		}

		// 按需添加根与子配置，使用去重函数
		addConfigLocked(commonConfig)

		if commonConfig.Server == nil {
			commonConfig.Server = &confv1.Server{}
		}
		addConfigLocked(commonConfig.Server)
	})
}

// addConfigLocked 假定已持有 muBC 锁，添加时会去重并确保参数为指针
func addConfigLocked(c proto.Message) {
	if c == nil {
		return
	}

	v := reflect.ValueOf(c)
	if !v.IsValid() || v.Kind() != reflect.Ptr || v.IsNil() {
		// 只接受非 nil 的指针类型
		return
	}

	addr := v.Pointer()
	if _, exists := configSet[addr]; exists {
		return
	}

	configList = append(configList, c)
	configSet[addr] = struct{}{}
}

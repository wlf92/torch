package registry

import (
	"context"

	"github.com/wlf92/torch/pkg/known"
)

type IRegistry interface {
	Register(ctx context.Context, ins *ServiceInstance) error                     // 注册服务实例
	Deregister(ctx context.Context, ins *ServiceInstance) error                   // 解注册服务实例
	Watch(ctx context.Context, serviceName string) (IWatcher, error)              // 监听相同服务名的服务实例变化
	Services(ctx context.Context, serviceName string) ([]*ServiceInstance, error) // 获取服务实例列表
}

type IDiscovery interface {
	Watch(ctx context.Context, serviceName string) (IWatcher, error)              // 监听相同服务名的服务实例变化
	Services(ctx context.Context, serviceName string) ([]*ServiceInstance, error) // 获取服务实例列表
}

type IWatcher interface {
	Next() ([]*ServiceInstance, error) // 返回服务实例列表
	Stop() error                       // 停止监听
}

type ServiceInstance struct {
	ID       string      `json:"id"`       // 服务实体ID，每个服务实体ID唯一
	Name     string      `json:"name"`     // 服务实体名
	Kind     known.Kind  `json:"kind"`     // 服务实体类型
	Alias    string      `json:"alias"`    // 服务实体别名
	State    known.State `json:"state"`    // 服务实例状态
	Routes   []uint32    `json:"routes"`   // 服务路由ID
	Endpoint string      `json:"endpoint"` // 服务器实体暴露端口
}

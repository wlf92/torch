package consul

import (
	"context"
	"net"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/spf13/cast"

	"github.com/wlf92/torch/internal/launch"
	"github.com/wlf92/torch/pkg/known"
	"github.com/wlf92/torch/registry"
)

var _ registry.IRegistry = &Registry{}

type options struct {
	addr                           string          // 客户端连接地址, 内建客户端配置，默认为127.0.0.1:8500
	client                         *api.Client     // 外部客户端, 外部客户端配置，存在外部客户端时，优先使用外部客户端，默认为nil
	ctx                            context.Context // 上下文, 默认为context.Background
	enableHealthCheck              bool            // 是否启用健康检查, 默认为true
	healthCheckInterval            int             // 健康检查时间间隔（秒），仅在启用健康检查后生效, 默认10秒
	healthCheckTimeout             int             // 健康检查超时时间（秒），仅在启用健康检查后生效, 默认5秒
	enableHeartbeatCheck           bool            // 是否启用心跳检查, 默认为true
	heartbeatCheckInterval         int             // 心跳检查时间间隔（秒），仅在启用心跳检查后生效, 默认10秒
	deregisterCriticalServiceAfter int             // 健康检测失败后自动注销服务时间（秒）, 默认30秒
}

type Registry struct {
	err        error
	ctx        context.Context
	cancel     context.CancelFunc
	opts       *options
	watchers   sync.Map
	registrars sync.Map
}

func NewRegistry() *Registry {
	o := new(options)
	o.addr = launch.Config.Consul.Addr
	o.ctx = context.Background()
	o.enableHealthCheck = launch.Config.Consul.HealthCheck
	o.healthCheckInterval = launch.Config.Consul.HealthCheckInterval
	o.healthCheckTimeout = launch.Config.Consul.HealthCheckTimeout
	o.enableHeartbeatCheck = launch.Config.Consul.HeartbeatCheck
	o.heartbeatCheckInterval = launch.Config.Consul.HeartbeatCheckInterval
	o.deregisterCriticalServiceAfter = launch.Config.Consul.DeregisterCriticalServiceAfter

	r := &Registry{}
	r.opts = o
	r.ctx, r.cancel = context.WithCancel(o.ctx)

	if o.client == nil {
		config := api.DefaultConfig()
		if o.addr != "" {
			config.Address = o.addr
		}
		o.client, r.err = api.NewClient(config)
	}

	return r
}

// Register 注册服务实例
func (r *Registry) Register(ctx context.Context, ins *registry.ServiceInstance) error {
	if r.err != nil {
		return r.err
	}

	v, ok := r.registrars.Load(ins.ID)
	if ok {
		return v.(*registrar).register(ctx, ins)
	}

	reg := newRegistrar(r)

	if err := reg.register(ctx, ins); err != nil {
		return err
	}

	r.registrars.Store(ins.ID, reg)

	return nil
}

// Deregister 解注册服务实例
func (r *Registry) Deregister(ctx context.Context, ins *registry.ServiceInstance) error {
	v, ok := r.registrars.Load(ins.ID)
	if ok {
		return v.(*registrar).deregister(ctx, ins)
	}

	return r.opts.client.Agent().ServiceDeregister(ins.ID)
}

// Services 获取服务实例列表
func (r *Registry) Services(ctx context.Context, serviceName string) ([]*registry.ServiceInstance, error) {
	if r.err != nil {
		return nil, r.err
	}

	v, ok := r.watchers.Load(serviceName)
	if ok {
		return v.(*watcherMgr).services(), nil
	} else {
		services, _, err := r.services(ctx, serviceName, 0, true)
		return services, err
	}
}

// Watch 监听服务
func (r *Registry) Watch(ctx context.Context, serviceName string) (registry.IWatcher, error) {
	if r.err != nil {
		return nil, r.err
	}

	v, ok := r.watchers.Load(serviceName)
	if ok {
		return v.(*watcherMgr).fork(), nil
	}

	w, err := newWatcherMgr(r, ctx, serviceName)
	if err != nil {
		return nil, err
	}
	r.watchers.Store(serviceName, w)

	return w.fork(), nil
}

// 获取服务实体列表
func (r *Registry) services(ctx context.Context, serviceName string, waitIndex uint64, passingOnly bool) ([]*registry.ServiceInstance, uint64, error) {
	opts := &api.QueryOptions{
		WaitIndex: waitIndex,
		WaitTime:  60 * time.Second,
	}
	opts.WithContext(ctx)

	entries, meta, err := r.opts.client.Health().Service(serviceName, "", passingOnly, opts)
	if err != nil {
		return nil, 0, err
	}

	services := make([]*registry.ServiceInstance, 0, len(entries))
	for _, entry := range entries {
		ins := &registry.ServiceInstance{
			ID:     entry.Service.ID,
			Name:   entry.Service.Service,
			Routes: make([]uint32, 0, len(entry.Service.Meta)),
		}

		for scheme, addr := range entry.Service.TaggedAddresses {
			if scheme == "lan_ipv4" || scheme == "wan_ipv4" || scheme == "lan_ipv6" || scheme == "wan_ipv6" {
				continue
			}
			ins.Endpoint = (&url.URL{
				Scheme: scheme,
				Host:   net.JoinHostPort(addr.Address, strconv.Itoa(addr.Port)),
			}).String()
		}
		if ins.Endpoint == "" {
			continue
		}

		for k, v := range entry.Service.Meta {
			switch k {
			case metaFieldKind:
				ins.Kind = known.Kind(v)
			case metaFieldAlias:
				ins.Alias = v
			case metaFieldState:
				ins.State = known.State(v)
			case metaFieldRoutes:
				arr := strings.Split(v, ",")
				for _, v := range arr {
					ins.Routes = append(ins.Routes, cast.ToUint32(v))
				}
			}
		}

		services = append(services, ins)
	}

	return services, meta.LastIndex, nil
}

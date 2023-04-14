package router

import (
	"errors"
	"fmt"
	"sync"

	"github.com/wlf92/torch/endpoint"
	"github.com/wlf92/torch/pkg/log"
	"github.com/wlf92/torch/registry"
)

var (
	ErrNotFoundRoute    = errors.New("not found route")
	ErrNotFoundEndpoint = errors.New("not found endpoint")
)

type Router struct {
	strategy BalanceStrategy

	rw        sync.RWMutex
	routes    map[uint32]*Route             // 节点路由表
	endpoints map[string]*endpoint.Endpoint // 服务实例端点
}

func NewRouter(strategy BalanceStrategy) *Router {
	return &Router{
		routes:    make(map[uint32]*Route),
		endpoints: make(map[string]*endpoint.Endpoint),
		strategy:  strategy,
	}
}

// ReplaceServices 替换服务实例
func (r *Router) ReplaceServices(services ...*registry.ServiceInstance) {
	r.rw.Lock()
	defer r.rw.Unlock()

	r.routes = make(map[uint32]*Route, len(services))
	r.endpoints = make(map[string]*endpoint.Endpoint, len(services))

	for _, service := range services {
		if err := r.addService(service); err != nil {
			log.Errorw(fmt.Sprintf("service instance add failed, insID: %s kind: %s name: %s endpoint: %s err: %v", service.ID, service.Kind, service.Name, service.Endpoint, err))
		}
	}
}

// RemoveServices 移除服务实例
func (r *Router) RemoveServices(services ...*registry.ServiceInstance) {
	r.rw.Lock()
	defer r.rw.Unlock()

	for _, service := range services {
		delete(r.endpoints, service.ID)
		for _, id := range service.Routes {
			if route, ok := r.routes[id]; ok {
				route.removeEndpoint(service.ID)
			}
		}
	}
}

// AddServices 添加服务实例
func (r *Router) AddServices(services ...*registry.ServiceInstance) {
	r.rw.Lock()
	defer r.rw.Unlock()

	for _, service := range services {
		if err := r.addService(service); err != nil {
			log.Errorw(fmt.Sprintf("service instance add failed, insID: %s kind: %s name: %s endpoint: %s err: %v", service.ID, service.Kind, service.Name, service.Endpoint, err))
		}
	}
}

// AddService 添加服务实例
func (r *Router) AddService(service *registry.ServiceInstance) error {
	r.rw.Lock()
	defer r.rw.Unlock()

	return r.addService(service)
}

// 添加服务实例
func (r *Router) addService(service *registry.ServiceInstance) error {
	ep, err := endpoint.ParseEndpoint(service.Endpoint)
	if err != nil {
		return err
	}

	r.endpoints[service.ID] = ep
	for _, id := range service.Routes {
		route, ok := r.routes[id]
		if !ok {
			route = newRoute(r, id)
			r.routes[id] = route
		}
		route.addEndpoint(service.ID, ep)
	}

	return nil
}

// FindServiceRoute 查找节点路由
func (r *Router) FindServiceRoute(routeID uint32) (*Route, error) {
	r.rw.RLock()
	defer r.rw.RUnlock()

	route, ok := r.routes[routeID]
	if !ok {
		return nil, ErrNotFoundRoute
	}

	return route, nil
}

// FindServiceEndpoint 查找服务端口
func (r *Router) FindServiceEndpoint(insID string) (*endpoint.Endpoint, error) {
	r.rw.RLock()
	defer r.rw.RUnlock()

	ep, ok := r.endpoints[insID]
	if !ok {
		return nil, ErrNotFoundEndpoint
	}

	return ep, nil
}

// IterationServiceEndpoint 迭代服务端口
func (r *Router) IterationServiceEndpoint(fn func(insID string, ep *endpoint.Endpoint) bool) {
	r.rw.RLock()
	defer r.rw.RUnlock()

	for insID, ep := range r.endpoints {
		if fn(insID, ep) == false {
			break
		}
	}
}

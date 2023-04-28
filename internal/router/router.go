package router

import (
	"errors"
	"fmt"
	"sync"

	"github.com/wlf92/torch/internal/endpoint"
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
	msgRoutes map[uint32]*Route             // 消息路由表
	svcRoutes map[string]*Route             // 微服务路由表
	endpoints map[string]*endpoint.Endpoint // 服务实例端点
}

func NewRouter(strategy BalanceStrategy) *Router {
	return &Router{
		msgRoutes: make(map[uint32]*Route),
		svcRoutes: make(map[string]*Route),
		endpoints: make(map[string]*endpoint.Endpoint),
		strategy:  strategy,
	}
}

// ReplaceServices 替换服务实例
func (r *Router) ReplaceServices(services ...*registry.ServiceInstance) {
	r.rw.Lock()
	defer r.rw.Unlock()

	r.msgRoutes = make(map[uint32]*Route, len(services))
	r.svcRoutes = make(map[string]*Route, len(services))
	r.endpoints = make(map[string]*endpoint.Endpoint, len(services))

	for _, service := range services {
		if err := r.addService(service); err != nil {
			log.Errorw(fmt.Sprintf("service instance add failed, insID: %s kind: %s name: %s endpoint: %s err: %v", service.ID, service.Kind, service.Name, service.Endpoint, err))
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

// 添加服务实例
func (r *Router) addService(service *registry.ServiceInstance) error {
	ep, err := endpoint.ParseEndpoint(service.Endpoint)
	if err != nil {
		return err
	}

	r.endpoints[service.ID] = ep

	route, ok := r.svcRoutes[service.Alias]
	if !ok {
		route = newRoute(r)
		r.svcRoutes[service.Alias] = route
	}
	route.addEndpoint(service.ID, ep)

	for _, id := range service.Routes {
		route, ok := r.msgRoutes[id]
		if !ok {
			route = newRoute(r)
			r.msgRoutes[id] = route
		}
		route.addEndpoint(service.ID, ep)
	}

	return nil
}

// FindMsgRoute 查找节点路由
func (r *Router) FindMsgRoute(routeID uint32) (*Route, error) {
	r.rw.RLock()
	defer r.rw.RUnlock()

	route, ok := r.msgRoutes[routeID]
	if !ok {
		return nil, ErrNotFoundRoute
	}

	return route, nil
}

// FindSvcRoute 查找节点路由
func (r *Router) FindSvcRoute(alias string) (*Route, error) {
	r.rw.RLock()
	defer r.rw.RUnlock()

	route, ok := r.svcRoutes[alias]
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
		if !fn(insID, ep) {
			break
		}
	}
}

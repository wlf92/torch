package router

import (
	"sync"
	"sync/atomic"

	"github.com/wlf92/torch/internal/endpoint"
)

type BalanceStrategy string

const (
	Random           BalanceStrategy = "random" // 随机
	RoundRobin       BalanceStrategy = "rr"     // 轮询
	WeightRoundRobin BalanceStrategy = "wrr"    // 加权轮询
)

type endpointInfo struct {
	insID    string
	index    int
	endpoint *endpoint.Endpoint
}

type Route struct {
	router  *Router // 路由器
	counter int64   // 轮询计数器

	rw          sync.RWMutex
	endpointMap map[string]*endpointInfo // 服务端点字典
	endpointArr []*endpointInfo          // 服务端点表
}

func newRoute(router *Router) *Route {
	return &Route{
		router:      router,
		endpointMap: make(map[string]*endpointInfo),
		endpointArr: make([]*endpointInfo, 0),
	}
}

// FindEndpoint 查询路由服务端点
func (r *Route) FindEndpoint(insID ...string) (*endpoint.Endpoint, error) {
	if len(insID) == 0 || insID[0] == "" {
		switch r.router.strategy {
		case RoundRobin:
			return r.roundRobinDispatch()
		default:
			return r.randomDispatch()
		}
	}

	return r.fixedDispatch(insID[0])
}

// 添加服务端点
func (r *Route) addEndpoint(insID string, ep *endpoint.Endpoint) {
	r.rw.Lock()
	defer r.rw.Unlock()

	switch r.router.strategy {
	case RoundRobin, WeightRoundRobin:
		info, ok := r.endpointMap[insID]
		if ok {
			info.endpoint = ep
		} else {
			info = &endpointInfo{
				insID:    insID,
				index:    len(r.endpointArr),
				endpoint: ep,
			}
			r.endpointArr = append(r.endpointArr, info)
			r.endpointMap[insID] = info
		}
	default:
		r.endpointMap[insID] = &endpointInfo{endpoint: ep}
	}
}

// 固定分配
func (r *Route) fixedDispatch(insID string) (*endpoint.Endpoint, error) {
	r.rw.RLock()
	defer r.rw.RUnlock()

	info, ok := r.endpointMap[insID]
	if !ok {
		return nil, ErrNotFoundEndpoint
	}

	return info.endpoint, nil
}

// 随机分配
func (r *Route) randomDispatch() (*endpoint.Endpoint, error) {
	r.rw.RLock()
	defer r.rw.RUnlock()

	for _, info := range r.endpointMap {
		return info.endpoint, nil
	}

	return nil, ErrNotFoundEndpoint
}

// 轮询分配
func (r *Route) roundRobinDispatch() (*endpoint.Endpoint, error) {
	r.rw.RLock()
	defer r.rw.RUnlock()

	if len(r.endpointArr) == 0 {
		return nil, ErrNotFoundEndpoint
	}

	counter := atomic.AddInt64(&r.counter, 1)
	index := int(counter % int64(len(r.endpointArr)))

	return r.endpointArr[index].endpoint, nil
}

package torch

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/wlf92/torch/internal/launch"
	"github.com/wlf92/torch/pkg/known"
	"github.com/wlf92/torch/pkg/log"
	"github.com/wlf92/torch/registry"
	"github.com/wlf92/torch/transport"
	"github.com/wlf92/torch/utils/xnet"
	"google.golang.org/grpc"
)

var _ IComponent = (*Node)(nil)

type Node struct {
	ctx    context.Context
	cancel context.CancelFunc

	registry registry.IRegistry
	instance *registry.ServiceInstance

	rpc *grpc.Server

	routes map[uint32]transport.RouteHandler

	cf *launch.Node
}

func NewNode(name string) *Node {
	nd := new(Node)
	nd.ctx, nd.cancel = context.WithCancel(context.Background())
	nd.routes = make(map[uint32]transport.RouteHandler)
	nd.cf = launch.Config.GetNodeByName(name)
	return nd
}

func (nd *Node) Name() string {
	return "node"
}

func (nd *Node) Init() {
	if nd.registry == nil {
		log.Fatalw("registry can not be empty")
	}
	if nd.cf.RPCPort == 0 {
		log.Fatalw("rpc_port can not be empty")
	}
}

func (nd *Node) Start() {
	nd.startRPCServer()
	nd.registerServiceInstance()

	log.Infow("node server startup successful")
}

func (nd *Node) Destroy() {
	nd.deregisterServiceInstance()
	nd.stopRPCServer()
	nd.cancel()
}

func (nd *Node) SetRegistry(r registry.IRegistry) {
	nd.registry = r
}

// 注册服务实例
func (nd *Node) registerServiceInstance() {
	ip, _ := xnet.InternalIP()
	ip = fmt.Sprintf("//%s:%d", ip, nd.cf.RPCPort)

	nd.instance = &registry.ServiceInstance{
		ID:       nd.Name(),
		Name:     string(known.Node),
		Kind:     known.Node,
		Alias:    nd.Name(),
		State:    known.Work,
		Endpoint: ip,
	}

	for k := range nd.routes {
		nd.instance.Routes = append(nd.instance.Routes, k)
	}

	ctx, cancel := context.WithTimeout(nd.ctx, 10*time.Second)
	err := nd.registry.Register(ctx, nd.instance)
	cancel()

	if err != nil {
		log.Fatalw(fmt.Sprintf("register service instance failed: %v", err))
	}
}

// 解注册服务实例
func (nd *Node) deregisterServiceInstance() {
	ctx, cancel := context.WithTimeout(nd.ctx, 10*time.Second)
	err := nd.registry.Deregister(ctx, nd.instance)
	defer cancel()
	if err != nil {
		log.Errorw(fmt.Sprintf("deregister service instance failed: %v", err))
	}
}

func (nd *Node) startRPCServer() {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", nd.cf.RPCPort))
	if err != nil {
		return
	}

	ln, err := net.ListenTCP(addr.Network(), addr)
	if err != nil {
		return
	}

	nd.rpc = grpc.NewServer()
	nd.rpc.RegisterService(&transport.Inner_ServiceDesc, &transport.Server{Routes: nd.routes})

	go func() {
		if err := nd.rpc.Serve(ln); err != nil {
			log.Fatalw(fmt.Sprintf("failed to serve: %v", err))
		}
	}()
}

func (nd *Node) stopRPCServer() {
	nd.rpc.Stop()
}

func (nd *Node) AddRouteHandler(route uint32, handler transport.RouteHandler) {
	nd.routes[route] = handler
}

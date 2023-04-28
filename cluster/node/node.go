package node

import (
	"context"
	"fmt"
	"net"
	"reflect"
	"sync"
	"time"

	"github.com/wlf92/torch"
	"github.com/wlf92/torch/internal/launch"
	"github.com/wlf92/torch/internal/router"
	"github.com/wlf92/torch/pkg/known"
	"github.com/wlf92/torch/pkg/log"
	"github.com/wlf92/torch/registry"
	"github.com/wlf92/torch/transport"
	"github.com/wlf92/torch/utils/xnet"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
)

var _ torch.IComponent = (*Node)(nil)

type Node struct {
	ctx    context.Context
	cancel context.CancelFunc

	registry registry.IRegistry
	instance *registry.ServiceInstance

	rpc     *grpc.Server
	rpcDesc *grpc.ServiceDesc
	rpcObj  interface{}

	routes map[uint32]interface{}

	rpcRouter  *router.Router // rpc路由器
	gateRouter *router.Router // gate路由器

	cf        *launch.Node
	name      string
	mpClients sync.Map
}

func Create(name string) *Node {
	nd := new(Node)
	nd.ctx, nd.cancel = context.WithCancel(context.Background())
	nd.routes = make(map[uint32]interface{})
	nd.cf = launch.Config.GetNodeByName(name)
	nd.rpcRouter = router.NewRouter(router.Random)
	nd.gateRouter = router.NewRouter(router.Random)
	nd.name = name
	nd.mpClients = sync.Map{}
	return nd
}

func (nd *Node) Name() string {
	return nd.name
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
	nd.watchServiceInstance()

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
	ep := fmt.Sprintf("//%s:%d", ip, nd.cf.RPCPort)

	nd.instance = &registry.ServiceInstance{
		ID:       fmt.Sprintf("%s-%s:%d", nd.Name(), ip, nd.cf.RPCPort),
		Name:     string(known.Node),
		Kind:     known.Node,
		Alias:    nd.Name(),
		State:    known.Work,
		Endpoint: ep,
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

func (nd *Node) watchServiceInstance() {
	for _, v := range []known.Kind{known.Node, known.Gate} {
		rctx, rcancel := context.WithTimeout(nd.ctx, 10*time.Second)
		watcher, err := nd.registry.Watch(rctx, string(v))
		rcancel()

		if err != nil {
			log.Fatalw(fmt.Sprintf("the service instance watch failed: %v", err))
		}
		go func(k known.Kind) {
			defer watcher.Stop()
			for {
				select {
				case <-nd.ctx.Done():
					return
				default:
					// exec watch
				}

				services, err := watcher.Next()
				if err != nil {
					continue
				}

				if k == known.Node {
					nd.rpcRouter.ReplaceServices(services...)
				} else if k == known.Gate {
					nd.gateRouter.ReplaceServices(services...)
				}
			}
		}(v)
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

func (nd *Node) SetRpcService(sd *grpc.ServiceDesc, ss interface{}) {
	nd.rpcDesc = sd
	nd.rpcObj = ss
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
	if nd.rpcDesc != nil {
		nd.rpc.RegisterService(nd.rpcDesc, nd.rpcObj)
	}

	go func() {
		if err := nd.rpc.Serve(ln); err != nil {
			log.Fatalw(fmt.Sprintf("failed to serve: %v", err))
		}
	}()
}

func (nd *Node) stopRPCServer() {
	nd.rpc.Stop()
}

func (nd *Node) AddRouteHandler(route uint32, handler interface{}) {
	tp := reflect.TypeOf(handler)
	if tp.Kind() != reflect.Func {
		panic("AddRouteHandler: handler must be func")
	}
	if tp.NumIn() != 2 || tp.NumOut() != 1 {
		panic("AddRouteHandler: in/out count error")
	}
	if tp.In(0).Kind() != reflect.Int64 {
		panic("AddRouteHandler")
	}
	if _, ok := reflect.New(tp.In(1).Elem()).Interface().(proto.Message); !ok {
		panic("AddRouteHandler")
	}
	if _, ok := reflect.New(tp.Out(0).Elem()).Interface().(proto.Message); !ok {
		panic("AddRouteHandler")
	}
	nd.routes[route] = handler
}

func (nd *Node) GetServiceClient(alias string) *grpc.ClientConn {
	route, err := nd.rpcRouter.FindSvcRoute(alias)
	if err != nil {
		return nil
	}

	ep, err := route.FindEndpoint()
	if err != nil {
		return nil
	}

	client, ok := nd.mpClients.Load(ep.Address())
	if !ok {
		// 如果带宽够，一个连接够了，没必要多个连接，grpc内部会保证重连的问题，所以也不用处理
		for i := 0; i < 3; i++ {
			client, err = grpc.Dial(ep.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				if i == 2 {
					return nil
				}
				continue
			}
			nd.mpClients.Store(ep.Address(), client)
			break
		}
	}
	return client.(*grpc.ClientConn)
}

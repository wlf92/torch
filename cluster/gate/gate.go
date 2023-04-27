package gate

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/wlf92/torch"
	"github.com/wlf92/torch/internal/launch"
	"github.com/wlf92/torch/internal/router"
	"github.com/wlf92/torch/network"
	"github.com/wlf92/torch/packet"
	"github.com/wlf92/torch/pkg/known"
	"github.com/wlf92/torch/pkg/log"
	"github.com/wlf92/torch/registry"
	"github.com/wlf92/torch/session"
	"github.com/wlf92/torch/transport"
	"github.com/wlf92/torch/utils/xnet"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type FilterHandler func(msgId uint32, conn network.Conn) bool // 返回true则不继续往下走
type ErrHandler func(conn network.Conn, err error)

var _ torch.IComponent = (*Gateway)(nil)
var instance *Gateway

type Gateway struct {
	server network.IServer // 网关服务器

	ctx    context.Context
	cancel context.CancelFunc

	registry      registry.IRegistry
	instance      *registry.ServiceInstance
	errHandler    ErrHandler
	filterHandler FilterHandler

	rpc     *grpc.Server
	rpcDesc *grpc.ServiceDesc
	rpcObj  interface{}

	gateRouter *router.Router // 网关路由器

	mpClients  sync.Map
	mpSessions sync.Map
}

func Create() *Gateway {
	gw := new(Gateway)
	gw.ctx, gw.cancel = context.WithCancel(context.Background())
	gw.gateRouter = router.NewRouter(router.Random)
	gw.mpClients = sync.Map{}
	gw.mpSessions = sync.Map{}

	instance = gw
	return gw
}

func (gw *Gateway) Name() string {
	return "gateway"
}

func (gw *Gateway) Init() {
	if gw.server == nil {
		log.Fatalw("server can not be empty")
	}
	if gw.registry == nil {
		log.Fatalw("registry can not be empty")
	}

	if launch.Config.Gate.RPCPort == 0 {
		log.Fatalw("rpc_port can not be empty")
	}
}

func (gw *Gateway) Start() {
	gw.startNetworkServer()
	gw.startRPCServer()
	gw.registerServiceInstance()
	gw.watchServiceInstance()

	log.Infow("gate server startup successful")
}

func (gw *Gateway) Destroy() {
	gw.deregisterServiceInstance()
	gw.stopNetworkServer()
	gw.stopRPCServer()
	gw.cancel()
}

func (gw *Gateway) SetServer(s network.IServer) {
	gw.server = s
}

func (gw *Gateway) SetRegistry(r registry.IRegistry) {
	gw.registry = r
}

func (gw *Gateway) SetErrorHandler(handler ErrHandler) {
	gw.errHandler = handler
}

// 注册服务实例
func (gw *Gateway) registerServiceInstance() {
	ip, _ := xnet.InternalIP()
	ep := fmt.Sprintf("//%s:%d", ip, launch.Config.Gate.RPCPort)

	gw.instance = &registry.ServiceInstance{
		ID:       fmt.Sprintf("%s-%s:%d", gw.Name(), ip, launch.Config.Gate.RPCPort),
		Name:     string(known.Gate),
		Kind:     known.Gate,
		Alias:    gw.Name(),
		State:    known.Work,
		Endpoint: ep,
	}

	ctx, cancel := context.WithTimeout(gw.ctx, 10*time.Second)
	err := gw.registry.Register(ctx, gw.instance)
	cancel()

	if err != nil {
		log.Fatalw(fmt.Sprintf("register service instance failed: %v", err))
	}
}

// 解注册服务实例
func (gw *Gateway) deregisterServiceInstance() {
	ctx, cancel := context.WithTimeout(gw.ctx, 10*time.Second)
	err := gw.registry.Deregister(ctx, gw.instance)
	defer cancel()
	if err != nil {
		log.Errorw(fmt.Sprintf("deregister service instance failed: %v", err))
	}
}

func (gw *Gateway) SetRpcService(sd *grpc.ServiceDesc, ss interface{}) {
	gw.rpcDesc = sd
	gw.rpcObj = ss
}

func (gw *Gateway) startRPCServer() {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", launch.Config.Gate.RPCPort))
	if err != nil {
		return
	}

	ln, err := net.ListenTCP(addr.Network(), addr)
	if err != nil {
		return
	}

	gw.rpc = grpc.NewServer()
	if gw.rpcDesc != nil {
		gw.rpc.RegisterService(gw.rpcDesc, gw.rpcObj)
	}

	go func() {
		if err := gw.rpc.Serve(ln); err != nil {
			log.Fatalw(fmt.Sprintf("failed to serve: %v", err))
		}
	}()
}

func (gw *Gateway) stopRPCServer() {
	gw.rpc.Stop()
}

func (gw *Gateway) startNetworkServer() {
	gw.server.OnConnect(gw.handleConnect)
	gw.server.OnDisconnect(gw.handleDisconnect)
	gw.server.OnReceive(gw.handleReceive)

	if err := gw.server.Start(); err != nil {
		log.Fatalw("network server start failed: %v", err)
	}
}

func (gw *Gateway) stopNetworkServer() {
	if err := gw.server.Stop(); err != nil {
		log.Errorw(fmt.Sprintf("network server stop failed: %v", err))
	}
}

func (gw *Gateway) watchServiceInstance() {
	rctx, rcancel := context.WithTimeout(gw.ctx, 10*time.Second)
	watcher, err := gw.registry.Watch(rctx, string(known.Node))
	rcancel()

	if err != nil {
		log.Fatalw(fmt.Sprintf("the service instance watch failed: %v", err))
	}
	go func() {
		defer watcher.Stop()
		for {
			select {
			case <-gw.ctx.Done():
				return
			default:
				// exec watch
			}

			services, err := watcher.Next()
			if err != nil {
				continue
			}

			gw.gateRouter.ReplaceServices(services...)
		}
	}()
}

// 处理连接打开
func (gw *Gateway) handleConnect(conn network.Conn) {
	log.Infow("connnect one")

	s := session.NewSession()
	gw.mpSessions.Store(conn.ID(), s)
}

// 处理断开连接
func (gw *Gateway) handleDisconnect(conn network.Conn) {
	log.Infow("disconnnect one")
	gw.mpSessions.Delete(conn.ID())
}

func (gw *Gateway) SetFilterFunc(handler FilterHandler) {
	gw.filterHandler = handler
}

// 处理接收到的消息
func (gw *Gateway) handleReceive(conn network.Conn, data []byte, _ int) {
	message, err := packet.Unpack(data)
	if err != nil {
		log.Errorw(fmt.Sprintf("unpack data to struct failed: %v", err))
		return
	}

	if gw.filterHandler != nil && gw.filterHandler(message.Route, conn) {
		return
	}

	route, err := gw.gateRouter.FindMsgRoute(message.Route)
	if err != nil {
		if gw.errHandler != nil {
			gw.errHandler(conn, fmt.Errorf("gateway.no.route"))
		}
		return
	}

	ep, err := route.FindEndpoint()
	if err != nil {
		if gw.errHandler != nil {
			gw.errHandler(conn, fmt.Errorf("gateway.no.service"))
		}
		return
	}

	client, ok := gw.mpClients.Load(ep.Address())
	if !ok {
		// 如果带宽够，一个连接够了，没必要多个连接，grpc内部会保证重连的问题，所以也不用处理
		for i := 0; i < 3; i++ {
			ct, err := grpc.Dial(ep.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				if i == 2 {
					if gw.errHandler != nil {
						gw.errHandler(conn, fmt.Errorf("gateway.dial.fail"))
					}
					continue
				}
			}
			client = transport.NewInnerClient(ct)
			gw.mpClients.Store(ep.Address(), client)
			break
		}
	}

	rsp, err := client.(transport.InnerClient).RouteRpc(gw.ctx, &transport.RouteRpcReq{
		MsgId:  message.Route,
		Datas:  message.Buffer,
		UserId: conn.UID(),
	})
	if err != nil {
		if gw.errHandler != nil {
			gw.errHandler(conn, err)
		}
		return
	}

	if err == nil && rsp != nil && len(rsp.Datas) > 0 {
		conn.Send(rsp.Datas)
	}
}

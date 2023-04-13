package torch

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/wlf92/torch/network"
	"github.com/wlf92/torch/pkg/known"
	"github.com/wlf92/torch/pkg/log"
	"github.com/wlf92/torch/registry"
	"github.com/wlf92/torch/transport"
	"github.com/wlf92/torch/utils/xnet"
	"google.golang.org/grpc"
)

var _ IComponent = (*Gateway)(nil)

type Gateway struct {
	server network.IServer // 网关服务器

	ctx    context.Context
	cancel context.CancelFunc

	registry registry.IRegistry
	instance *registry.ServiceInstance

	rpc *grpc.Server
}

func NewGateway() *Gateway {
	gw := new(Gateway)
	gw.ctx, gw.cancel = context.WithCancel(context.Background())

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
}

func (gw *Gateway) Start() {
	gw.startNetworkServer()
	gw.startRPCServer()
	gw.registerServiceInstance()

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

// 注册服务实例
func (gw *Gateway) registerServiceInstance() {
	ip, _ := xnet.InternalIP()
	ip = fmt.Sprintf("//%s:9999", ip)

	gw.instance = &registry.ServiceInstance{
		ID:       gw.Name(),
		Name:     string(known.Gate),
		Kind:     known.Gate,
		Alias:    gw.Name(),
		State:    known.Work,
		Endpoint: ip,
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

func (gw *Gateway) startRPCServer() {
	addr, err := net.ResolveTCPAddr("tcp", ":9999")
	if err != nil {
		return
	}

	ln, err := net.ListenTCP(addr.Network(), addr)
	if err != nil {
		return
	}

	gw.rpc = grpc.NewServer()
	gw.rpc.RegisterService(&transport.Inner_ServiceDesc, &transport.Server{})

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

// 处理连接打开
func (gw *Gateway) handleConnect(conn network.Conn) {
	// s := g.sessions.Get().(*session.Session)
	// s.Init(conn)
	// g.group.AddSession(s)
}

// 处理断开连接
func (gw *Gateway) handleDisconnect(conn network.Conn) {
	// s, err := g.group.RemSession(session.Conn, conn.ID())
	// if err != nil {
	// 	log.Errorf("session remove failed, gid: %d, cid: %d, uid: %d, err: %v", g.opts.id, s.CID(), s.UID(), err)
	// 	return
	// }

	// if uid := conn.UID(); uid > 0 {
	// 	ctx, cancel := context.WithTimeout(g.ctx, g.opts.timeout)
	// 	err = g.proxy.unbindGate(ctx, conn.ID(), uid)
	// 	cancel()
	// 	if err != nil {
	// 		log.Errorf("user unbind failed, gid: %d, uid: %d, err: %v", g.opts.id, uid, err)
	// 	}
	// }

	// s.Reset()
	// g.sessions.Put(s)
}

// 处理接收到的消息
func (gw *Gateway) handleReceive(conn network.Conn, data []byte, _ int) {
	// message, err := packet.Unpack(data)
	// if err != nil {
	// 	log.Errorf("unpack data to struct failed: %v", err)
	// 	return
	// }

	// ctx, cancel := context.WithTimeout(g.ctx, g.opts.timeout)
	// err = g.proxy.deliver(ctx, conn.ID(), conn.UID(), message)
	// cancel()
	// if err != nil {
	// 	log.Warnf("deliver message failed: %v", err)
	// }
}

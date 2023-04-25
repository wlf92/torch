package ws

import (
	"fmt"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/wlf92/torch/internal/launch"
	"github.com/wlf92/torch/network"
	"github.com/wlf92/torch/pkg/log"
)

type UpgradeHandler func(w http.ResponseWriter, r *http.Request) (allowed bool)

var _ network.IServer = (*server)(nil)

type server struct {
	addr              string
	listener          net.Listener              // 监听器
	connMgr           *connMgr                  // 连接管理器
	startHandler      network.StartHandler      // 服务器启动hook函数
	stopHandler       network.CloseHandler      // 服务器关闭hook函数
	connectHandler    network.ConnectHandler    // 连接打开hook函数
	disconnectHandler network.DisconnectHandler // 连接关闭hook函数
	receiveHandler    network.ReceiveHandler    // 接收消息hook函数
	upgradeHandler    UpgradeHandler            // HTTP协议升级成WS协议hook函数
}

func NewServer() *server {
	s := &server{}
	s.connMgr = newConnMgr(s)
	return s
}

func (s *server) Addr() string {
	return s.addr
}

func (s *server) Protocol() string {
	return "websocket"
}

func (s *server) Start() error {
	if err := s.init(); err != nil {
		return err
	}

	if s.startHandler != nil {
		s.startHandler()
	}

	go s.serve()

	return nil
}

func (s *server) Stop() error {
	if err := s.listener.Close(); err != nil {
		return err
	}

	s.connMgr.close()

	return nil
}

func (s *server) init() error {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", launch.Config.Gate.WsPort))
	if err != nil {
		return err
	}

	ln, err := net.ListenTCP(addr.Network(), addr)
	if err != nil {
		return err
	}

	s.listener = ln
	return nil
}

func (s *server) serve() {
	upgrader := websocket.Upgrader{
		ReadBufferSize:    4096,
		WriteBufferSize:   4096,
		EnableCompression: true,
		CheckOrigin:       func(r *http.Request) bool { return true },
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if s.upgradeHandler != nil && !s.upgradeHandler(w, r) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Errorw("websocket upgrade error: %v", err)
			return
		}

		if err = s.connMgr.allocate(conn); err != nil {
			_ = conn.Close()
		}
	})
	http.Serve(s.listener, nil)
}

func (s *server) OnStart(handler network.StartHandler) {
	s.startHandler = handler
}

func (s *server) OnUpgrade(handler UpgradeHandler) {
	s.upgradeHandler = handler
}

func (s *server) OnStop(handler network.CloseHandler) {
	s.stopHandler = handler
}

func (s *server) OnConnect(handler network.ConnectHandler) {
	s.connectHandler = handler
}

func (s *server) OnReceive(handler network.ReceiveHandler) {
	s.receiveHandler = handler
}

func (s *server) OnDisconnect(handler network.DisconnectHandler) {
	s.disconnectHandler = handler
}

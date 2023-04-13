package network

type (
	StartHandler      func()
	CloseHandler      func()
	ConnectHandler    func(conn Conn)
	DisconnectHandler func(conn Conn)
	ReceiveHandler    func(conn Conn, msg []byte, wsMsgType int)
)

type IServer interface {
	Addr() string                           // 监听地址
	Start() error                           // 启动服务器
	Stop() error                            // 关闭服务器
	Protocol() string                       // 协议
	OnStart(handler StartHandler)           // 监听服务器启动
	OnStop(handler CloseHandler)            // 监听服务器关闭
	OnConnect(handler ConnectHandler)       // 监听连接打开
	OnReceive(handler ReceiveHandler)       // 监听接收消息
	OnDisconnect(handler DisconnectHandler) // 监听连接断开
}

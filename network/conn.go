package network

import (
	"errors"
	"net"
)

const (
	ConnOpened ConnState = iota + 1 // 连接打开
	ConnHanged                      // 连接挂起
	ConnClosed                      // 连接关闭
)

var (
	ErrConnectionHanged  = errors.New("connection is hanged")
	ErrConnectionClosed  = errors.New("connection is closed")
	ErrIllegalMsgType    = errors.New("illegal message type")
	ErrTooManyConnection = errors.New("too many connection")
)

type (
	ConnState int32

	Conn interface {
		ID() int64                     // 获取连接ID
		UID() int64                    // 获取用户ID
		Send(bts []byte) error         // 发送消息
		State() ConnState              // 获取连接状态
		Close(isForce ...bool) error   // 关闭连接
		LocalIP() (string, error)      // 获取本地IP
		LocalAddr() (net.Addr, error)  // 获取本地地址
		RemoteIP() (string, error)     // 获取远端IP
		RemoteAddr() (net.Addr, error) // 获取远端地址
	}
)

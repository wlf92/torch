package session

import (
	"sync"

	"github.com/wlf92/torch/network"
)

type Session struct {
	rw   sync.RWMutex // 读写锁
	conn network.Conn // 连接
}

func NewSession() *Session {
	return &Session{}
}

func (s *Session) Init(conn network.Conn) {
	s.rw.Lock()
	defer s.rw.Unlock()

	s.conn = conn
}

func (s *Session) CID() int64 {
	s.rw.RLock()
	defer s.rw.RUnlock()

	return s.conn.ID()
}

func (s *Session) UID() int64 {
	s.rw.RLock()
	defer s.rw.RUnlock()

	return s.conn.UID()
}

func (s *Session) Bind(uid int64) {
	s.rw.Lock()
	defer s.rw.Unlock()

	s.conn.Bind(uid)
}

func (s *Session) Send(msg []byte) error {
	s.rw.RLock()
	defer s.rw.RUnlock()

	return s.conn.Send(msg)
}

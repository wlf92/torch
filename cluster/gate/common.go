package gate

import "github.com/wlf92/torch/session"

func BroadCast(userIds []int64, msgId uint32, datas []byte) {

}

func BindUser(connId, userId int64) error {
	v, ok := instance.mpSessions.Load(connId)
	if !ok {
		return nil
	}

	s := v.(*session.Session)
	s.Bind(userId)
	return nil
}

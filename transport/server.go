package transport

import (
	"context"
)

type RouteHandler func(p []byte) []byte

type Server struct {
	UnimplementedInnerServer
	Routes map[uint32]RouteHandler
}

func (e *Server) MessageRoute(ctx context.Context, req *MessageRouteReq) (*MessageRouteRsp, error) {
	reply := &MessageRouteRsp{}

	for _, v := range req.Msgs {
		routeId := v.GetMsgId()
		if f, ok := e.Routes[routeId]; ok {
			bts := f(v.GetContent())
			if bts != nil {
				reply.Msgs = append(reply.Msgs, &SingleBack{
					UserId:  v.UserId,
					Content: bts,
				})
			}
		}
	}
	return reply, nil
}

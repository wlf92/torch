package transport

import context "context"

type RouteHandler func(channel, area int32, uid int64, bts []byte) ([]byte, error)

type Server struct {
	UnimplementedInnerServer
	Routes map[uint32]RouteHandler
}

func (e *Server) RouteRpc(ctx context.Context, req *RouteRpcReq) (*RouteRpcRsp, error) {
	reply := &RouteRpcRsp{}

	routeId := req.GetMsgId()
	if f, ok := e.Routes[routeId]; ok {
		bts, err := f(req.GetChannelId(), req.GetAreaId(), req.GetUserId(), req.GetDatas())
		if err != nil {
			return reply, err
		}

		if bts != nil {
			reply.UserId = req.UserId
			reply.Datas = bts
		}
	}

	return reply, nil
}

func (e *Server) HttpRpc(ctx context.Context, req *HttpRpcReq) (*HttpRpcRsp, error) {
	reply := &HttpRpcRsp{}

	// for _, v := range req.Msgs {
	// 	routeId := req.GetMsgId()
	// 	if f, ok := e.Routes[routeId]; ok {
	// 		bts := f(req.GetChannelId(), req.GetAreaId(), req.GetUserId(), req.GetDatas())
	// 		if bts != nil {
	// 			reply.Msgs = append(reply.Msgs, &SingleRouteBack{
	// 				UserId: req.UserId,
	// 				Datas:  bts,
	// 			})
	// 		}
	// 	}
	// }

	return reply, nil
}

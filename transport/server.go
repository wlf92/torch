package transport

import (
	context "context"
	"fmt"
	reflect "reflect"

	"github.com/wlf92/torch/packet"
	"google.golang.org/protobuf/proto"
)

type Server struct {
	UnimplementedInnerServer
	Routes map[uint32]interface{}
}

func (e *Server) RouteRpc(ctx context.Context, req *RouteRpcReq) (*RouteRpcRsp, error) {
	reply := &RouteRpcRsp{}

	routeId := req.GetMsgId()
	if handler, ok := e.Routes[routeId]; ok {
		obj, ok := reflect.New(reflect.TypeOf(handler).In(1).Elem()).Interface().(proto.Message)
		if !ok {
			return nil, fmt.Errorf("inner.route.req.create.fail")
		}

		err := proto.Unmarshal(req.GetDatas(), obj)
		if err != nil {
			return nil, fmt.Errorf("inner.route.req.parse.fail")
		}

		f := reflect.ValueOf(handler)
		callResults := f.Call([]reflect.Value{reflect.ValueOf(req.GetUserId()), reflect.ValueOf(obj)})

		if len(callResults) != 1 {
			return nil, fmt.Errorf("inner.route.rsp.nil")
		}

		// 如果没有返回，则默认这个请求是没有数据返回的，用于埋点之类的
		if callResults[0].IsNil() {
			return nil, nil
		}

		bts, err := proto.Marshal(callResults[0].Interface().(proto.Message))
		if err != nil {
			return nil, fmt.Errorf("inner.route.rsp.serializable.fail")
		}

		reply.Datas = packet.Pack(&packet.Message{Route: routeId + 1, Buffer: bts})
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

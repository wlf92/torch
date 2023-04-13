package transport

import (
	"context"
)

type Server struct {
	UnimplementedInnerServer
}

func (e *Server) MessageRoute(ctx context.Context, req *MessageRouteReq) (*MessageRouteRsp, error) {

	return nil, nil
}

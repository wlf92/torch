package transport

import (
	"context"

	"google.golang.org/grpc"
)

type client struct {
	ct InnerClient
}

func NewClient(cc *grpc.ClientConn) *client {
	return &client{ct: NewInnerClient(cc)}
}

func (c *client) RouteRpc(ctx context.Context, req *RouteRpcReq) (rsp *RouteRpcRsp, err error) {
	return c.ct.RouteRpc(ctx, req)
}

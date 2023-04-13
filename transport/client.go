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

func (c *client) MessageRoute(ctx context.Context, req *MessageRouteReq) (rsp *MessageRouteReq, err error) {

	return
}

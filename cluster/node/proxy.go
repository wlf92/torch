package node

import "google.golang.org/grpc"

type IProxy interface {
	GetServiceClient(alias string) *grpc.ClientConn
	GetGateClient(insId string) *grpc.ClientConn
}

package torch

import (
	"github.com/wlf92/torch/network"
	"github.com/wlf92/torch/pkg/log"
)

type Gateway struct {
	server network.IServer // 网关服务器
}

func (gw *Gateway) Name() string {
	return "gateway"
}

func (gw *Gateway) Init() {
	if gw.server == nil {
		log.Fatalw("server can not be empty")
	}
}

func (gw *Gateway) Start() {
	if err := gw.server.Start(); err != nil {
		log.Fatalw("network server start failed: %v", err)
	}

	log.Infow("gate server startup successful")
}

func (gw *Gateway) Destroy() {

}

func (gw *Gateway) SetServer(s network.IServer) {
	gw.server = s
}

func NewGateway() *Gateway {
	return new(Gateway)
}

package torch

type Gateway struct {
}

func (b *Gateway) Name() string {
	return "gateway"
}

func (b *Gateway) Init() {

}

func (b *Gateway) Start() {

}

func (b *Gateway) Restart() {

}

func (b *Gateway) Destroy() {

}

func NewGateway() IComponent {
	return new(Gateway)
}

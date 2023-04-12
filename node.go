package torch

type Node struct {
}

func (b *Node) Name() string {
	return "node"
}

func (b *Node) Init() {

}

func (b *Node) Start() {

}

func (b *Node) Restart() {

}

func (b *Node) Destroy() {

}

func NewNode() IComponent {
	return new(Node)
}

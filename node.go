package torch

var _ IComponent = (*Node)(nil)

type Node struct {
}

func (b *Node) Name() string {
	return "node"
}

func (b *Node) Init() {

}

func (b *Node) Start() {

}

func (b *Node) Destroy() {

}

func NewNode() *Node {
	return new(Node)
}

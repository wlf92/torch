package known

const (
	// XRequestIDKey 用来定义上下文中的键，代表请求的 uuid.
	XRequestIDKey = "X-Request-ID"

	// XUsernameKey 用来定义上下文的键，代表请求的所有者.
	XUsernameKey = "X-Username"
)

type Kind string

const (
	Gate Kind = "gate" // 网关服
	Node Kind = "node" // 节点服
)

type State string

const (
	Work State = "work" // 工作（节点正常工作，可以分配更多玩家到该节点）
	Busy State = "busy" // 繁忙（节点资源紧张，不建议分配更多玩家到该节点上）
	Hang State = "hang" // 挂起（节点即将关闭，正处于资源回收中）
	Shut State = "shut" // 关闭（节点已经关闭，无法正常访问该节点）
)

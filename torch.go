package torch

type IComponent interface {
	Name() string // Name 组件名称
	Init()        // Init 初始化组件
	Start()       // Start 启动组件
	Restart()     // Restart 重启组件
	Destroy()     // Destroy 销毁组件
}

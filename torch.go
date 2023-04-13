package torch

import (
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/wlf92/torch/pkg/log"
)

type IComponent interface {
	Name() string // Name 组件名称
	Init()        // Init 初始化组件
	Start()       // Start 启动组件
	Destroy()     // Destroy 销毁组件
}

type Container struct {
	sig       chan os.Signal
	component IComponent
}

func NewContainer(component IComponent) *Container {
	return &Container{sig: make(chan os.Signal), component: component}
}

func (c *Container) Serve() {
	c.component.Init()
	c.component.Start()

	switch runtime.GOOS {
	case `windows`:
		signal.Notify(c.sig, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	default:
		signal.Notify(c.sig, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGKILL, syscall.SIGTERM)
	}

	sig := <-c.sig
	log.Errorw("process got signal %v, container will close", sig)
	signal.Stop(c.sig)

	c.component.Destroy()
}

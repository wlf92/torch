package locate

import (
	"context"

	"github.com/wlf92/torch/pkg/known"
)

type Locator interface {
	Get(ctx context.Context, uid int64, insKind known.Kind) (string, error)     // Get 获取用户定位
	Set(ctx context.Context, uid int64, insKind known.Kind, insID string) error // Set 设置用户定位
	Rem(ctx context.Context, uid int64, insKind known.Kind, insID string) error // Rem 移除用户定位
	Watch(ctx context.Context, insKinds ...known.Kind) (Watcher, error)         // Watch 监听用户定位变化
}

type Watcher interface {
	Next() ([]*Event, error) // Next 返回用户位置列表
	Stop() error             // Stop 停止监听
}

type Event struct {
	UID     int64      `json:"uid"`      // 用户ID
	Type    EventType  `json:"type"`     // 事件类型
	InsID   string     `json:"ins_id"`   // 实例ID
	InsKind known.Kind `json:"ins_kind"` // 实例类型
}

type EventType int

const (
	SetLocation EventType = iota // 设置定位
	RemLocation                  // 移除定位
)

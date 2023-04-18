package eventbus

import (
	"context"
	"encoding/json"
	"fmt"

	"time"

	"github.com/wlf92/torch/internal/xvalue"
	"github.com/wlf92/torch/pkg/log"
	"github.com/wlf92/torch/utils/xconv"
	"github.com/wlf92/torch/utils/xtime"
	"github.com/wlf92/torch/utils/xuuid"
)

var globalEventbus Eventbus

type Eventbus interface {
	Close() error                                                              // Close 关闭事件总线
	Publish(ctx context.Context, topic string, message interface{}) error      // Publish 发布事件
	Subscribe(ctx context.Context, topic string, handler EventHandler) error   // Subscribe 订阅事件
	Unsubscribe(ctx context.Context, topic string, handler EventHandler) error // Unsubscribe 取消订阅
}

type EventHandler func(event *Event)

type Event struct {
	ID        string       // 事件ID
	Topic     string       // 事件主题
	Payload   xvalue.Value // 事件载荷
	Timestamp time.Time    // 事件时间
}

type data struct {
	ID        string `json:"id"`        // 事件ID
	Topic     string `json:"topic"`     // 事件主题
	Payload   string `json:"payload"`   // 事件载荷
	Timestamp int64  `json:"timestamp"` // 事件时间
}

// PackData 打包数据
func PackData(topic string, payload interface{}) ([]byte, error) {
	id, err := xuuid.UUID()
	if err != nil {
		return nil, err
	}

	return json.Marshal(&data{
		ID:        id,
		Topic:     topic,
		Payload:   xconv.String(payload),
		Timestamp: xtime.Now().UnixNano(),
	})
}

// UnpackData 解析
func UnpackData(v []byte) (*Event, error) {
	d := &data{}

	err := json.Unmarshal(v, d)
	if err != nil {
		return nil, err
	}

	return &Event{
		ID:        d.ID,
		Topic:     d.Topic,
		Payload:   xvalue.NewValue(d.Payload),
		Timestamp: xtime.UnixNano(d.Timestamp),
	}, nil
}

// SetEventbus 设置事件总线
func SetEventbus(eventbus Eventbus) {
	if globalEventbus != nil {
		if err := globalEventbus.Close(); err != nil {
			log.Errorw(fmt.Sprintf("the old eventbus close failed: %v", err))
		}
	}

	globalEventbus = eventbus
}

// GetEventbus 获取事件总线
func GetEventbus() Eventbus {
	return globalEventbus
}

// Publish 发布事件
func Publish(ctx context.Context, topic string, message interface{}) error {
	if globalEventbus == nil {
		log.Warnw("the eventbus component is not injected, and the publish operation will be ignored.")
		return nil
	}

	return globalEventbus.Publish(ctx, topic, message)
}

// Subscribe 订阅事件
func Subscribe(ctx context.Context, topic string, handler EventHandler) error {
	if globalEventbus == nil {
		log.Warnw("the eventbus component is not injected, and the subscribe operation will be ignored.")
		return nil
	}

	return globalEventbus.Subscribe(ctx, topic, handler)
}

// Unsubscribe 取消订阅
func Unsubscribe(ctx context.Context, topic string, handler EventHandler) error {
	if globalEventbus == nil {
		log.Warnw("the eventbus component is not injected, and the unsubscribe operation will be ignored.")
		return nil
	}

	return globalEventbus.Unsubscribe(ctx, topic, handler)
}

// Close 关闭事件总线
func Close() error {
	if globalEventbus == nil {
		return nil
	}

	return globalEventbus.Close()
}

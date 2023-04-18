package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/wlf92/torch/locate"
	"github.com/wlf92/torch/pkg/known"
	"github.com/wlf92/torch/pkg/log"
	"golang.org/x/sync/singleflight"
)

const (
	userLocationsKey = "%s:locate:user:%d:locations" // hash
	channelEventKey  = "%s:locate:channel:%v:event"  // channel
)

var _ locate.Locator = &Locator{}

type Locator struct {
	ctx      context.Context
	cancel   context.CancelFunc
	sfg      singleflight.Group // singleFlight
	watchers sync.Map

	client redis.UniversalClient
	prefix string
}

func NewLocator(ctx context.Context, client redis.UniversalClient) *Locator {
	l := &Locator{}
	l.ctx, l.cancel = context.WithCancel(ctx)
	l.client = client
	l.prefix = "locate"
	return l
}

// Get 获取用户定位
func (l *Locator) Get(ctx context.Context, uid int64, insKind known.Kind) (string, error) {
	key := fmt.Sprintf(userLocationsKey, l.prefix, uid)
	val, err, _ := l.sfg.Do(key+string(insKind), func() (interface{}, error) {
		val, err := l.client.HGet(ctx, key, string(insKind)).Result()
		if err != nil && err != redis.Nil {
			return "", err
		}
		return val, nil
	})

	if err != nil {
		return "", err
	}

	return val.(string), nil
}

// Set 设置用户定位
func (l *Locator) Set(ctx context.Context, uid int64, insKind known.Kind, insID string) error {
	key := fmt.Sprintf(userLocationsKey, l.prefix, uid)
	err := l.client.HSet(ctx, key, string(insKind), insID).Err()
	if err != nil {
		return err
	}

	err = l.publish(ctx, uid, insKind, insID, locate.SetLocation)
	if err != nil {
		log.Errorf("location event publish failed: %v", err)
	}

	return nil
}

// Rem 移除用户定位
func (l *Locator) Rem(ctx context.Context, uid int64, insKind known.Kind, insID string) error {
	oldInsID, err := l.Get(ctx, uid, insKind)
	if err != nil {
		return err
	}

	if oldInsID == "" || oldInsID != insID {
		return nil
	}

	key := fmt.Sprintf(userLocationsKey, l.prefix, uid)
	err = l.client.HDel(ctx, key, string(insKind)).Err()
	if err != nil {
		return err
	}

	err = l.publish(ctx, uid, insKind, insID, locate.RemLocation)
	if err != nil {
		log.Errorf("location event publish failed: %v", err)
	}

	return nil
}

func (l *Locator) publish(ctx context.Context, uid int64, insKind known.Kind, insID string, eventType locate.EventType) error {
	msg, err := marshal(&locate.Event{
		UID:     uid,
		Type:    eventType,
		InsID:   insID,
		InsKind: insKind,
	})
	if err != nil {
		return err
	}

	channel := fmt.Sprintf(channelEventKey, l.prefix, string(insKind))

	return l.client.Publish(ctx, channel, msg).Err()
}

func (l *Locator) toUniqueKey(insKinds ...known.Kind) string {
	sort.Slice(insKinds, func(i, j int) bool {
		return insKinds[i] < insKinds[j]
	})

	keys := make([]string, 0, len(insKinds))
	for _, insKind := range insKinds {
		keys = append(keys, string(insKind))
	}

	return strings.Join(keys, "&")
}

// Watch 监听用户定位变化
func (l *Locator) Watch(ctx context.Context, insKinds ...known.Kind) (locate.Watcher, error) {
	key := l.toUniqueKey(insKinds...)

	v, ok := l.watchers.Load(key)
	if ok {
		return v.(*watcherMgr).fork(), nil
	}

	w, err := newWatcherMgr(ctx, l, key, insKinds...)
	if err != nil {
		return nil, err
	}

	l.watchers.Store(key, w)

	return w.fork(), nil
}

func marshal(event *locate.Event) (string, error) {
	buf, err := json.Marshal(event)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func unmarshal(data []byte) (*locate.Event, error) {
	event := &locate.Event{}
	if err := json.Unmarshal(data, event); err != nil {
		return nil, err
	}
	return event, nil
}

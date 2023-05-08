package redis

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/wlf92/torch/mq"
)

type MessageQueue struct {
	ctx    context.Context
	cancel context.CancelFunc

	client   redis.UniversalClient
	status   int32
	mpTopics sync.Map
	mpGroups sync.Map
}

func Create(ctx context.Context, client redis.UniversalClient) mq.IMessageQueue {
	m := &MessageQueue{}
	m.ctx, m.cancel = context.WithCancel(ctx)
	m.client = client
	return m
}

func (m *MessageQueue) check(topic string, channel string) error {
	info, err := m.client.XInfoStream(m.ctx, topic).Result()
	if err != nil {
		err = m.Publish(topic, map[string]interface{}{"active": 1})
		if err != nil {
			return err
		}
		info, err = m.client.XInfoStream(m.ctx, topic).Result()
		if err != nil {
			return err
		}
	}

	m.mpTopics.Store(topic, struct{}{})

	groups, err := m.client.XInfoGroups(m.ctx, topic).Result()
	if err != nil {
		return err
	}

	isExist := false
	for _, v := range groups {
		if v.Name == channel {
			isExist = true
		}
		m.mpGroups.Store(fmt.Sprintf("%s-%s", topic, v.Name), struct{}{})
	}

	if !isExist {
		m.client.XGroupCreate(m.ctx, topic, channel, info.LastGeneratedID)
		m.mpGroups.Store(fmt.Sprintf("%s-%s", topic, channel), struct{}{})
	}
	return nil
}

func (m *MessageQueue) Publish(topic string, body map[string]interface{}) error {
	err := m.client.XAdd(m.ctx, &redis.XAddArgs{
		MaxLen: 5,
		Approx: false,
		Stream: topic,
		Values: body,
	}).Err()
	return err
}

func (m *MessageQueue) Subscribe(topic string, channel string, handler mq.HandlerFunc) error {
	if channel == "" {
		channel = "default"
	}

	go func() {
		for atomic.LoadInt32(&m.status) == 0 {
			time.Sleep(time.Second * 5)
		}

		for {
			_, ok1 := m.mpTopics.Load(topic)
			_, ok2 := m.mpGroups.Load(fmt.Sprintf("%s-%s", topic, channel))
			if !(ok1 && ok2) {
				err := m.check(topic, channel)
				if err != nil {
					time.Sleep(time.Second * 2)
					continue
				}
				break
			}
			break
		}

		for {
			entries, err := m.client.XReadGroup(m.ctx, &redis.XReadGroupArgs{
				Group:    channel,
				Consumer: "default",
				Streams:  []string{topic, ">"},
				Count:    5,
				Block:    0,
				NoAck:    false,
			}).Result()

			if err != nil {
				time.Sleep(time.Second * 3)
				continue
			}

			for i := 0; i < len(entries[0].Messages); i++ {
				messageID := entries[0].Messages[i].ID
				values := entries[0].Messages[i].Values

				if values["rpt"] != nil {
					bts := values["rpt"].(string)
					err := handler([]byte(bts))
					if err != nil {
						log.Fatal(err)
					}
				}
				m.client.XAck(m.ctx, topic, channel, messageID)
			}
		}
	}()
	return nil
}

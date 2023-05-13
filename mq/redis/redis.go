package redis

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/wlf92/torch/mq"
)

const ErrBusyGroup = "BUSYGROUP Consumer Group name already exists"

type MessageQueue struct {
	ctx    context.Context
	cancel context.CancelFunc

	client redis.UniversalClient
}

func Create(ctx context.Context, client redis.UniversalClient) mq.IMessageQueue {
	m := &MessageQueue{}
	m.ctx, m.cancel = context.WithCancel(ctx)
	m.client = client
	return m
}

func (m *MessageQueue) Publish(topic string, body []byte) error {
	err := m.client.XAdd(m.ctx, &redis.XAddArgs{
		MaxLen: 1000,
		Approx: false,
		Stream: topic,
		Values: map[string]interface{}{"body": body},
	}).Err()
	return err
}

func (m *MessageQueue) Subscribe(topic string, group string, handler mq.HandlerFunc) error {
	if group == "" {
		group = "default"
	}

	err := m.client.XGroupCreateMkStream(context.Background(), topic, group, "$").Err()
	if err != nil && err.Error() != ErrBusyGroup {
		return err
	}

	go func() {
		for {
			entries, err := m.client.XReadGroup(m.ctx, &redis.XReadGroupArgs{
				Group:    group,
				Consumer: "default",
				Streams:  []string{topic, ">"},
				Count:    1,
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
				err := handler([]byte(values["body"].(string)))
				if err != nil {
					log.Fatal(err)
				}

				m.client.XAck(m.ctx, topic, group, messageID)
			}
		}
	}()
	return nil
}

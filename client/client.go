package client

import (
	"context"
	"reflect"
	"strings"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/rlog"
	"github.com/sirupsen/logrus"
)

type Client struct {
	consumer    rocketmq.PushConsumer
	registerMap map[string]func(msg []byte) error
}

func NewClient(group string) (*Client, error) {
	rlog.SetLogLevel("warn")
	c, err := rocketmq.NewPushConsumer(
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{"127.0.0.1:9876"})),
		consumer.WithRetry(2),
		consumer.WithGroupName(group),
		consumer.WithConsumeFromWhere(consumer.ConsumeFromLastOffset),
	)
	if err != nil {
		logrus.Errorf("init consumer err %v", err)
		return nil, err
	}
	m := make(map[string]func(msg []byte) error)
	return &Client{consumer: c, registerMap: m}, nil
}
func (c *Client) Register(h Model, f func(msg []byte) error) {
	name := strings.ToLower(reflect.TypeOf(h).Name())
	c.registerMap[name] = f
}
func (c *Client) Run() error {
	for name, f := range c.registerMap {
		err := c.consumer.Subscribe(name, consumer.MessageSelector{}, func(ctx context.Context,
			msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
			for _, msg := range msgs {
				if err := f(msg.Body); err != nil {
					logrus.Error(err)
					c.consumer.Shutdown()
				}
			}
			return consumer.ConsumeSuccess, nil
		})
		if err != nil {
			logrus.Errorf("register map err %v", err)
			return err
		}
	}
	err := c.consumer.Start()
	if err != nil {
		logrus.Errorf("start consumer err %v", err)
		return err
	}
	return nil
}
func (c *Client) Stop() error {
	return c.consumer.Shutdown()
}

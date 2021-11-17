package client

import (
	"reflect"
	"strings"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/sirupsen/logrus"
)

type Client struct {
	clientID    string
	consumer    stan.Conn
	registerMap map[string]func(m []byte)
}
type Config struct {
	Host      string
	Port      string
	ClusterID string
	ClientID  string
}

func NewClient(c Config) (*Client, error) {
	nc, err := nats.Connect(c.Host + ":" + c.Port)
	if err != nil {
		logrus.Fatalf("[NewClient] nats connect err %+v", err)
	}
	sc, err := stan.Connect(c.ClusterID, c.ClientID, stan.NatsConn(nc))
	if err != nil {
		logrus.Fatalf("[NewClient] stan connect err %+v", err)
	}
	m := make(map[string]func(m []byte))
	return &Client{consumer: sc, registerMap: m, clientID: c.ClientID}, nil
}
func (c *Client) Register(h Model, f func(m []byte)) {
	name := strings.ToLower(reflect.TypeOf(h).Name())
	c.registerMap[name] = f
}
func (c *Client) Run() error {
	for name, f := range c.registerMap {
		_, err := c.consumer.Subscribe(name, func(msg *stan.Msg) {
			f(msg.Data)
		}, stan.DurableName(name+c.clientID))
		if err != nil {
			logrus.Errorf("[Run] sub %s err: %+v", name, err)
		}
	}
	return nil
}
func (c *Client) Close() error {
	err := c.consumer.Close()
	if err != nil {
		logrus.Errorf("[Close] close client error")
	}
	return nil
}

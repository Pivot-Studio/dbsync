package client

import (
	"reflect"
	"strings"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/sirupsen/logrus"
)

type Client struct {
	clientID string
	consumer stan.Conn
	cfg      Config
}
type Config struct {
	Host      string
	Port      string
	ClusterID string
	ClientID  string
	ProName   string
}

func NewClient(cfg Config) (*Client, error) {
	nc, err := nats.Connect(cfg.Host + ":" + cfg.Port)
	if err != nil {
		logrus.Fatalf("[NewClient] nats connect err %+v", err)
	}
	sc, err := stan.Connect(cfg.ClusterID, cfg.ClientID, stan.NatsConn(nc))
	if err != nil {
		logrus.Fatalf("[NewClient] stan connect err %+v", err)
	}
	return &Client{consumer: sc, clientID: cfg.ClientID, cfg: cfg}, nil
}
func (c *Client) Register(h Model, f func(m []byte)) {
	name := strings.ToLower(reflect.TypeOf(h).Name())
	_, err := c.consumer.Subscribe(c.cfg.ProName+name, func(msg *stan.Msg) {
		f(msg.Data)
	}, stan.DurableName(c.cfg.ProName+name+c.clientID))
	if err != nil {
		logrus.Errorf("[Run] sub %s err: %+v", name, err)
	}
}
func (c *Client) Close() error {
	err := c.consumer.Close()
	if err != nil {
		logrus.Errorf("[Close] close client error")
	}
	return nil
}

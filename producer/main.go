package main

import (
	"context"
	"log"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/sirupsen/logrus"
)

var (
	c *canal.Canal
	p rocketmq.Producer
)

type EventHandler struct {
	canal.DummyEventHandler
}

func (h *EventHandler) OnRow(e *canal.RowsEvent) error {
	// p.SendAsync(context.Background(), OnSendFail, primitive.NewMessage(e.Table.Name, e.RawData))
	logrus.Infof("[%s]radata:%+v\nrows:%+v", e.Table.Name, e.RawData, e.Rows)
	return nil
}

func OnSendFail(ctx context.Context, result *primitive.SendResult, err error) {

}

func main() {
	cfg := canal.NewDefaultConfig()
	cfg.Addr = "rm-2ze90979778u1xh5uvo.mysql.rds.aliyuncs.com:3306"
	cfg.User = "pivotstudio"
	// We only care table canal_test in test db
	cfg.Dump.TableDB = "husthole"
	cfg.Dump.Tables = []string{"replies", "holes"}
	var err error
	c, err = canal.NewCanal(cfg)
	if err != nil {
		log.Fatal(err)
	}
	// Register a handler to handle RowsEvent
	// c.SetEventHandler(&EventHandler{})
	// p, _ = roctmq.NewProducer(
	// 	producer.WikethNsResolver(primitive.NewPassthroughResolver([]string{"127.0.0.1:9876"})),
	// 	producer.WithRetry(2),
	// 	producer.WithQueueSelector(producer.NewManualQueueSelector()))

	// err = p.Start()
	// if err != nil {
	// 	fmt.Printf("start producer error: %s", err.Error())
	// 	os.Exit(1)
	// }
	// Start canal
	c.Run()
}

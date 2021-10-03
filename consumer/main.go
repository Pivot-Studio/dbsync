package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/go-mysql-org/go-mysql/replication"
	"github.com/sirupsen/logrus"
)

var (
	parser *replication.BinlogParser
)

func init() {
	parser = replication.NewBinlogParser()
}

func main() {
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithGroupName("testGroup"),
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{"127.0.0.1:9876"})),
	)
	err := c.Subscribe("test", consumer.MessageSelector{}, func(ctx context.Context,
		msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for i := range msgs {
			fmt.Printf("subscribe callback: %v \n", msgs[i])
			ev, err := parser.Parse(msgs[i].Body)
			if err != nil {
				logrus.Errorf("[]error:%+v", err)
				return consumer.ConsumeRetryLater, err
			}
			e := ev.Event.(*replication.RowsEvent)
			logrus.Infof("%s %+v", ev.Header.EventType, e.Rows)
		}

		return consumer.ConsumeSuccess, nil
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	// Note: start after subscribe
	err = c.Start()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	time.Sleep(time.Hour)
	err = c.Shutdown()
	if err != nil {
		fmt.Printf("shutdown Consumer error: %s", err.Error())
	}
}

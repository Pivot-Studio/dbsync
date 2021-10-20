package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/rlog"
	"github.com/sirupsen/logrus"
)

func main() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	rlog.SetLogLevel("warn")
	c, err := rocketmq.NewPushConsumer(
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{"127.0.0.1:9876"})),
		consumer.WithRetry(2),
	)
	if err != nil {
		logrus.Fatalf("init consumer err %v", err)
	}
	err = c.Subscribe("RocketTest", consumer.MessageSelector{}, func(ctx context.Context,
		msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		logrus.Infof("recive num: %v \n", len(msgs))
		for _, msg := range msgs {
			logrus.Infof("msg content: %s", msg.Body)
		}
		return consumer.ConsumeSuccess, nil
	})
	if err != nil {
		logrus.Fatal("init consumer err %v", err)
	}
	done := make(chan struct{}, 1)
	go func() {
		err = c.Start()
		if err != nil {
			logrus.Fatalf("start consumer err %v", err)
		}
		done <- struct{}{}
	}()
	n := <-sc
	logrus.Infof("receive signal %v, closing", n)
	err = c.Shutdown()
	if err != nil {
		logrus.Fatalf("shutdown err consumer err %v", err)
	}
	<-done

}

package mq

import (
	"context"
	"dbsync/model"
	"encoding/json"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/sirupsen/logrus"
)

type RockerTransfer struct {
	topic string
	p     rocketmq.Producer
}

func (r *RockerTransfer) DoBulk(rows []model.RowRequest) error {
	var ms []*primitive.Message
	for _, row := range rows {
		rowByte, err := json.Marshal(row)
		if err != nil {
			logrus.Errorf("json marshal err %v\nmessage row %v", err, row)
			return err
		}
		ms = append(ms, &primitive.Message{Topic: r.topic, Body: rowByte})
	}
	_, err := r.p.SendSync(context.Background(), ms...)
	if err != nil {
		logrus.Errorf("rocket send message err %v\nmessage rows %v", err, rows)
		return err
	}
	logrus.Info("send message success,length:%d", len(rows))
	return nil
}
func (r *RockerTransfer) InitRocket() error {
	p, err := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{"127.0.0.1:9876"})),
		producer.WithRetry(2),
	)
	if err != nil {
		logrus.Errorf("init rocket err %v", err)
		return err
	}
	logrus.Info("init rocket success")
	r.p = p
	return nil
}
func (r *RockerTransfer) Run() error {
	err := r.p.Start()
	if err != nil {
		logrus.Errorf("start rocket err %v", err)
		return err
	}
	logrus.Info("start rocket success")
	return nil
}

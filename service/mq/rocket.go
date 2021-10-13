package mq

import (
	"context"
	"dbsync/model"
	"encoding/json"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
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

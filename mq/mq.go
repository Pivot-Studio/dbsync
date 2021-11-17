package mq

import (
	"github.com/Pivot-Studio/dbsync/conf"
	"github.com/Pivot-Studio/dbsync/model"

	"github.com/sirupsen/logrus"
)

var MQ MessageQueue

type MessageQueue interface {
	DoBulk(rows []*model.RowRequest) error
	Initialize() error
	Close() error
}

func init() {
	switch conf.C.MQName {
	case "stan":
		{
			logrus.Infof("[init] use stan as storage")
			MQ = &StanMQ{}
		}
	default:
		{
			logrus.Fatalf("[init] unkonwn mq,try to use stan")
		}
	}
	err := MQ.Initialize()
	if err != nil {
		logrus.Fatalf("[init] mq init err: %+v", err)
	}
}

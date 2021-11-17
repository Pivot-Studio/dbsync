package mq

import (
	"dbsync/conf"
	"dbsync/model"

	json "github.com/json-iterator/go"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/sirupsen/logrus"
)

type StanMQ struct {
	sc stan.Conn
}

func (s *StanMQ) Initialize() error {
	nc, err := nats.Connect(conf.C.Stan.Host + ":" + conf.C.Stan.Port)
	if err != nil {
		logrus.Errorf("[Initialize] nats connect err %+v", err)
		return err
	}
	sc, err := stan.Connect(conf.C.Stan.ClusterID, conf.C.Stan.ClientID, stan.NatsConn(nc))
	if err != nil {
		logrus.Errorf("[Initialize] stan connect err %+v", err)
		return err
	}
	s.sc = sc
	return nil
}
func (s *StanMQ) Close() error {
	err := s.sc.Close()
	if err != nil {
		logrus.Errorf("[Close] stan close err: %+v", err)
		return err
	}
	return nil
}
func (s *StanMQ) DoBulk(rows []*model.RowRequest) error {
	for _, row := range rows {
		b, err := json.Marshal(row)
		if err != nil {
			logrus.Errorf("[DoBulk] marshal row err: %+v", err)
			return err
		}
		err = s.sc.Publish(row.Table, b)
		if err != nil {
			logrus.Errorf("[DoBulk] publish row err: %+v", err)
			return err
		}
	}
	return nil
}

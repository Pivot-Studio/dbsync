package river

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Pivot-Studio/dbsync/conf"
	"github.com/Pivot-Studio/dbsync/mq"
	"github.com/Pivot-Studio/dbsync/storage"

	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/sirupsen/logrus"
)

var (
	bulkSize = 128
	interval = 200 * time.Millisecond
)

type River struct {
	canal      *canal.Canal
	ctx        context.Context
	wg         sync.WaitGroup
	cancel     context.CancelFunc
	syncCh     chan interface{}
	storageDao storage.PositionStorager
	transfer   mq.MessageQueue
}

func NewRiver() (*River, error) {
	r := new(River)
	r.syncCh = make(chan interface{}, 4096)
	r.ctx, r.cancel = context.WithCancel(context.Background())
	r.transfer = mq.MQ
	r.storageDao = storage.StorageDao
	err := r.initCanal()
	if err != nil {
		logrus.Error("[NewRiver] init canal err %+v", err)
		return nil, err
	}
	return r, nil
}
func (r *River) Run() error {
	r.wg.Add(1)
	go r.syncLoop()
	logrus.Info("[Run] sleep 5 seconds......")
	time.Sleep(5 * time.Second)

	pos, err := r.storageDao.Get()
	if err != nil {
		logrus.Warnf("[Run] get pos in storage err %+v", err)
		pos, err = r.canal.GetMasterPos()
		if err != nil {
			logrus.Errorf("[Run] get master pos err %+v", err)
			return err
		}
		logrus.Infof("[Run] get master pos %+v", pos)
	}
	err = r.canal.RunFrom(pos)
	if err != nil {
		logrus.Errorf("[Run] start canal err %+v", err)
		return err
	}
	return nil
}
func (r *River) Close() {
	logrus.Infof("[Close] closing river")
	r.cancel()
	r.canal.Close()
	err := r.storageDao.Close()
	if err != nil {
		logrus.Errorf("[Close] close storage err: %+v", err)
	}
	err = r.transfer.Close()
	if err != nil {
		logrus.Errorf("[Close] close mq err: %+v", err)
	}
	r.wg.Wait()
}
func (r *River) Ctx() context.Context {
	return r.ctx
}
func (r *River) initCanal() error {
	cfg := canal.NewDefaultConfig()
	cfg.Addr = conf.C.Mysql.Host + ":" + conf.C.Mysql.Port
	cfg.User = conf.C.Mysql.Username
	cfg.Password = conf.C.Mysql.Password
	c, err := canal.NewCanal(cfg)
	if err != nil {
		logrus.Error(err)
		return fmt.Errorf("[initCanal] err: %+v", err)
	}
	c.SetEventHandler(&eventHandler{r})
	r.canal = c
	return nil
}

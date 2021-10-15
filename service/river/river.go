package river

import (
	"context"
	"fmt"
	"sync"
	"time"

	"dbsync/service/mq"
	"dbsync/storage"

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
	transfer   *mq.RockerTransfer
}

func NewRiver() (*River, error) {
	r := new(River)
	r.syncCh = make(chan interface{}, 4096)
	r.ctx, r.cancel = context.WithCancel(context.Background())

	r.transfer = new(mq.RockerTransfer)
	err := r.transfer.InitRocket()
	if err != nil {
		logrus.Error("init rocket err %v", err)
		return nil, err
	}

	err = r.initStorager(&storage.BoltPositionStorage{})
	if err != nil {
		logrus.Error("init storager err %v", err)
		return nil, err
	}

	err = r.initCanal()
	if err != nil {
		logrus.Error("init canal err %v", err)
		return nil, err
	}

	logrus.Info("init river success")
	return r, nil
}
func (r *River) Run() error {

	err := r.transfer.Run()
	if err != nil {
		logrus.Errorf("start transfer err %v", err)
		return err
	}

	r.wg.Add(1)
	go r.syncLoop()
	logrus.Info("sleep 5 seconds......")
	time.Sleep(5 * time.Second)

	pos, err := r.storageDao.Get()
	if err != nil {
		logrus.Warnf("get pos in storage err %v", err)
		pos, err = r.canal.GetMasterPos()
		if err != nil {
			logrus.Errorf("get master pos err %v", err)
			return err
		}
		logrus.Infof("get master pos %v", pos)
	}
	err = r.canal.RunFrom(pos)
	if err != nil {
		logrus.Errorf("start canal err %v", err)
		return err
	}
	return nil
}
func (r *River) Close() {
	logrus.Infof("closing river")
	r.cancel()
	r.canal.Close()
	r.storageDao.Close()
	r.wg.Wait()
}
func (r *River) Ctx() context.Context {
	return r.ctx
}
func (r *River) initCanal() error {
	cfg := canal.NewDefaultConfig()
	cfg.Addr = "127.0.0.1:3306"
	cfg.User = "canal"
	cfg.Password = "canal"
	c, err := canal.NewCanal(cfg)
	if err != nil {
		logrus.Error(err)
		return fmt.Errorf("init canal err %v", err)
	}
	c.SetEventHandler(&eventHandler{r})
	r.canal = c
	return nil
}
func (r *River) initStorager(s storage.PositionStorager) error {
	r.storageDao = s
	err := r.storageDao.Initialize()
	if err != nil {
		logrus.Errorf("init storager err %v", err)
		return err
	}
	return nil
}

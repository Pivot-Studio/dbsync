package storage

import (
	"dbsync/conf"

	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/sirupsen/logrus"
)

var (
	StorageDao  PositionStorager
	positionKey = conf.C.Storage.PosKey
)

type PositionStorager interface {
	Initialize() error
	Save(pos mysql.Position) error
	Get() (mysql.Position, error)
	Close() error
}

func init() {
	switch conf.C.Storage.DaoName {
	case "redis":
		{
			logrus.Infof("[init] use redis as storage")
			StorageDao = &redisPositionStorage{}
		}
	case "bolt":
		{
			logrus.Infof("[init] use bolt as storage")
			StorageDao = &boltPositionStorage{}
		}
	default:
		{
			logrus.Fatalf("[init] unkonwn stroage,try to use redis or bolt")
		}
	}
	err := StorageDao.Initialize()
	if err != nil {
		logrus.Fatalf("[init] storage init err: %+v", err)
	}
}

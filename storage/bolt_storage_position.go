package storage

import (
	"dbsync/conf"
	"encoding/json"
	"fmt"

	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/sirupsen/logrus"
	"go.etcd.io/bbolt"
)

var bucket = []byte(conf.C.Bolt.Bucket)

type boltPositionStorage struct {
	bolt *bbolt.DB
}

func (b *boltPositionStorage) Save(pos mysql.Position) error {
	logrus.Infof("[Save] position at: %+v", pos)
	return b.bolt.Update(func(tx *bbolt.Tx) error {
		bt := tx.Bucket(bucket)
		data, err := json.Marshal(pos)
		if err != nil {
			return err
		}
		return bt.Put([]byte(positionKey), data)
	})
}

func (b *boltPositionStorage) Get() (mysql.Position, error) {
	var entity mysql.Position
	err := b.bolt.View(func(tx *bbolt.Tx) error {
		bt := tx.Bucket(bucket)
		data := bt.Get([]byte(positionKey))
		if data == nil {
			return fmt.Errorf("[Get] position storage not found")
		}
		return json.Unmarshal(data, &entity)
	})

	return entity, err
}
func (b *boltPositionStorage) Initialize() error {
	db, err := bbolt.Open(conf.C.Bolt.File, 0600, nil)
	if err != nil {
		logrus.Errorf("[Initialize] bbolt err %+v", err)
		return err
	}
	err = db.Update(func(tx *bbolt.Tx) error {
		tx.CreateBucketIfNotExists(bucket)
		return nil
	})
	if err != nil {
		logrus.Errorf("[Initialize] bucket err %+v", err)
		return err
	}
	logrus.Info("[Initialize] bolt success")
	b.bolt = db
	return nil
}
func (b *boltPositionStorage) Close() error {
	err := b.bolt.Close()
	if err != nil {
		logrus.Errorf("[Close] bolt err %+v", err)
		return err
	}
	return nil
}

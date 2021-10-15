package storage

import (
	"encoding/json"
	"fmt"

	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/sirupsen/logrus"
	"go.etcd.io/bbolt"
)

type BoltPositionStorage struct {
	bolt *bbolt.DB
}

func (b *BoltPositionStorage) Save(pos mysql.Position) error {
	logrus.Infof("save position %v", pos)
	return b.bolt.Update(func(tx *bbolt.Tx) error {
		bt := tx.Bucket(positionBucket)
		data, err := json.Marshal(pos)
		if err != nil {
			return err
		}
		return bt.Put(positionKey, data)
	})
}

func (b *BoltPositionStorage) Get() (mysql.Position, error) {
	var entity mysql.Position
	err := b.bolt.View(func(tx *bbolt.Tx) error {
		bt := tx.Bucket(positionBucket)
		data := bt.Get(positionKey)
		if data == nil {
			return fmt.Errorf("position storage not found")
		}
		return json.Unmarshal(data, &entity)
	})

	return entity, err
}
func (b *BoltPositionStorage) Initialize() error {
	db, err := bbolt.Open("my.db", 0600, nil)
	if err != nil {
		logrus.Errorf("init bbolt err %v", err)
		return err
	}
	err = db.Update(func(tx *bbolt.Tx) error {
		tx.CreateBucketIfNotExists(positionBucket)
		return nil
	})
	if err != nil {
		logrus.Errorf("create bucket err %v", err)
		return err
	}
	logrus.Info("init bolt success")
	b.bolt = db
	return nil
}
func (b *BoltPositionStorage) Close() error {
	err := b.bolt.Close()
	if err != nil {
		logrus.Errorf("close bolt err %v", err)
		return err
	}
	return nil
}

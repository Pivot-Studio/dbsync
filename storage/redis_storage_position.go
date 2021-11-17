package storage

import (
	"context"
	"dbsync/conf"

	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-redis/redis/v8"
	json "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
)

var ctx = context.Background()

type redisPositionStorage struct {
	rdb *redis.Client
}

func (r *redisPositionStorage) Save(pos mysql.Position) error {
	logrus.Infof("[Save] position at: %+v", pos)
	data, err := json.Marshal(pos)
	if err != nil {
		logrus.Errorf("[Save] marshal json parse err: %+v", err)
		return err
	}
	err = r.rdb.Set(ctx, positionKey, data, 0).Err()
	if err != nil {
		logrus.Errorf("[Save] postion err: %+v", err)
		return err
	}
	return nil
}

func (r *redisPositionStorage) Get() (mysql.Position, error) {
	var entity mysql.Position
	data, err := r.rdb.Get(ctx, positionKey).Bytes()
	if err != nil {
		return entity, err
	}
	err = json.Unmarshal(data, &entity)
	if err != nil {
		logrus.Errorf("[Get] unmarshal json parse err: %+v", err)
		return entity, err
	}
	return entity, nil
}
func (r *redisPositionStorage) Initialize() error {
	r.rdb = redis.NewClient(&redis.Options{
		Addr:     conf.C.Redis.Host + ":" + conf.C.Redis.Port,
		Password: conf.C.Redis.Password,
	})
	_, err := r.rdb.Ping(ctx).Result()
	if err != nil {
		logrus.Errorf("[Initialize] ping redis err: %+v", err)
		return err
	}
	logrus.Infof("[Initialize] redis success")
	return nil
}
func (r *redisPositionStorage) Close() error {
	err := r.rdb.Close()
	if err != nil {
		logrus.Errorf("[Close] redis err: %+v", err)
		return err
	}
	return nil
}

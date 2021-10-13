package storage

import (
	"go.etcd.io/bbolt"
)

type boltPositionStorage struct {
	bolt *bbolt.DB
}

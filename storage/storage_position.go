package storage

import (
	"github.com/go-mysql-org/go-mysql/mysql"
)

var (
	positionBucket = []byte("PositionBucket")
	positionKey    = []byte("PositionKey")
)

type PositionStorager interface {
	Initialize() error
	Save(pos mysql.Position) error
	Get() (mysql.Position, error)
	Close() error
}

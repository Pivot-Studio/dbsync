package storage

import (
	"github.com/go-mysql-org/go-mysql/mysql"
)

type PositionStorager interface {
	Initialize() error
	Save(pos mysql.Position) error
	Get() (mysql.Position, error)
}

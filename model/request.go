package model

import (
	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/mysql"
)

type RowRequest struct {
	RowsEvent canal.RowsEvent
}

type PosRequest struct {
	Pos   mysql.Position
	Force bool
}

package model

import (
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/schema"
)

type RowRequest struct {
	Table      string
	Column     []schema.TableColumn
	Action     string
	BeforeData []interface{}
	AfterData  []interface{}
}
type PosRequest struct {
	Pos   mysql.Position
	Force bool
}

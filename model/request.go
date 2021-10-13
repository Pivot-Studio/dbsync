package model

import "github.com/go-mysql-org/go-mysql/mysql"

type RowRequest struct {
	// Action    string
	// Timestamp uint32
	// Old       []interface{}
	// Row       []interface{}
}

type PosRequest struct {
	Pos   mysql.Position
	Force bool
}

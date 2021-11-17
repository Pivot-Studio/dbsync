package river

import (
	"fmt"
	"github.com/Pivot-Studio/dbsync/model"

	"github.com/go-mysql-org/go-mysql/canal"
)

func makeInsertRequest(e *canal.RowsEvent) ([]*model.RowRequest, error) {
	reqs := make([]*model.RowRequest, 0, len(e.Rows))
	for _, row := range e.Rows {
		req := &model.RowRequest{Column: e.Table.Columns, Action: e.Action, Table: e.Table.Name}
		req.AfterData = row
		reqs = append(reqs, req)
	}
	return reqs, nil
}
func makeDeleteRequest(e *canal.RowsEvent) ([]*model.RowRequest, error) {
	reqs := make([]*model.RowRequest, 0, len(e.Rows))
	for _, row := range e.Rows {
		req := &model.RowRequest{Column: e.Table.Columns, Action: e.Action, Table: e.Table.Name}
		req.BeforeData = row
		reqs = append(reqs, req)
	}
	return reqs, nil
}
func makeUpdateRequest(e *canal.RowsEvent) ([]*model.RowRequest, error) {
	reqs := make([]*model.RowRequest, 0, len(e.Rows))
	for i := 0; i < len(e.Rows); i += 2 {
		req := &model.RowRequest{Column: e.Table.Columns, Action: e.Action, Table: e.Table.Name}
		req.BeforeData, req.AfterData = e.Rows[i], e.Rows[i+1]
		reqs = append(reqs, req)
	}
	return reqs, nil
}
func (r *River) makeRequest(e *canal.RowsEvent) (reqs []*model.RowRequest, err error) {
	switch e.Action {
	case canal.InsertAction:
		reqs, err = makeInsertRequest(e)
	case canal.DeleteAction:
		reqs, err = makeDeleteRequest(e)
	case canal.UpdateAction:
		reqs, err = makeUpdateRequest(e)
	default:
		err = fmt.Errorf("invalid rows action %s", e.Action)
	}
	return
}

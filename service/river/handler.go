package river

import (
	"dbsync/model"
	"fmt"

	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
)

type eventHandler struct {
	r *River
}

func (h *eventHandler) OnRotate(e *replication.RotateEvent) error {
	pos := mysql.Position{
		Name: string(e.NextLogName),
		Pos:  uint32(e.Position),
	}

	h.r.syncCh <- model.PosRequest{Pos: pos, Force: true}

	return h.r.ctx.Err()
}

func (h *eventHandler) OnTableChanged(schema, table string) error {
	return nil
}

func (h *eventHandler) OnDDL(nextPos mysql.Position, _ *replication.QueryEvent) error {
	h.r.syncCh <- model.PosRequest{Pos: nextPos, Force: true}
	return h.r.ctx.Err()
}

func (h *eventHandler) OnXID(nextPos mysql.Position) error {
	h.r.syncCh <- model.PosRequest{Pos: nextPos, Force: false}
	return h.r.ctx.Err()
}

func (h *eventHandler) OnRow(e *canal.RowsEvent) error {
	var err error
	req := model.RowRequest{RowsEvent: *e}
	reqs := []*model.RowRequest{&req}
	if err != nil {
		h.r.cancel()
		return fmt.Errorf("make %s ES request err %v, close sync", e.Action, err)
	}

	h.r.syncCh <- reqs
	return h.r.ctx.Err()
}

func (h *eventHandler) OnGTID(gtid mysql.GTIDSet) error {
	return nil
}

func (h *eventHandler) OnPosSynced(pos mysql.Position, set mysql.GTIDSet, force bool) error {
	return nil
}

func (h *eventHandler) String() string {
	return "ESRiverEventHandler"
}

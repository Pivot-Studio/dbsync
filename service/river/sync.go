package river

import (
	"dbsync/model"
	"time"

	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/sirupsen/logrus"
)

func (r *River) syncLoop() {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	defer r.wg.Done()

	lastSavedTime := time.Now()
	reqs := make([]model.RowRequest, 0, 1024)

	var pos mysql.Position

	for {
		needFlush := false
		needSavePos := false

		select {
		case v := <-r.syncCh:
			switch v := v.(type) {
			case model.PosRequest:
				now := time.Now()
				if v.Force || now.Sub(lastSavedTime) > 3*time.Second {
					lastSavedTime = now
					needFlush = true
					needSavePos = true
					pos = v.Pos
				}
			case []model.RowRequest:
				reqs = append(reqs, v...)
				needFlush = len(reqs) >= bulkSize
			}
		case <-ticker.C:
			needFlush = true
		case <-r.ctx.Done():
			return
		}

		if needFlush {
			if err := r.transfer.DoBulk(reqs); err != nil {
				logrus.Errorf("do mq bulk err %v, close sync", err)
				r.cancel()
				return
			}
			reqs = reqs[0:0]
		}

		if needSavePos {
			if err := r.storageDao.Save(pos); err != nil {
				logrus.Errorf("save sync position %v err %v, close sync", pos, err)
				r.cancel()
				return
			}
		}
	}
}

package river

import (
	"context"
	"sync"
	"time"

	"dbsync/service/mq"
	"dbsync/storage"

	"github.com/go-mysql-org/go-mysql/canal"
)

var (
	bulkSize = 128
	interval = 200 * time.Millisecond
)

type River struct {
	canal      *canal.Canal
	ctx        context.Context
	wg         sync.WaitGroup
	cancel     context.CancelFunc
	syncCh     chan interface{}
	storageDao storage.PositionStorager
	transfer   *mq.RockerTransfer
}

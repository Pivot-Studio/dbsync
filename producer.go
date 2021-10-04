package dbsync

import (
	"encoding/json"
	"io/ioutil"
	"os"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
)

// PersistenMeta persistent
type PersistentMeta struct {
	pos mysql.Position
}

type Producer struct {
	db       *badger.DB
	meta     PersistentMeta
	syncer   *replication.BinlogSyncer
	streamer *replication.BinlogStreamer
}

// NewProducer creates a new producer by config file specified by path
func NewProducer(path string) (p *Producer) {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	b, _ := ioutil.ReadAll(f)
	var config replication.BinlogSyncerConfig
	json.Unmarshal(b, &config)
	db, err := badger.Open(badger.DefaultOptions("badger"))
	if err != nil {
		panic(err)
	}
	return &Producer{
		db:     db,
		syncer: replication.NewBinlogSyncer(config),
	}
}

func (p *Producer) restorePos() {
	err := p.db.View(func(txn *badger.Txn) error {
		// TODO
		return nil
	})
	if err != nil {
	}
}

func (p *Producer) Start() error {
	var err error
	p.streamer, err = p.syncer.StartSync(mysql.Position{})
	if err != nil {
		return err
	}
	return nil
}

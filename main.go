package dbsync

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
)

func m_ain() {
	// Create a binlog syncer with a unique server id, the server id must be different from other MySQL's.
	// flavor is mysql or mariadb
	cfg := replication.BinlogSyncerConfig{
		ServerID: 100,
		Flavor:   "mysql",
		Host:     "rm-2ze90979778u1xh5uvo.mysql.rds.aliyuncs.com",
		Port:     3306,
		User:     "pivotstudio",
		Password: "PivotStudio@2020",
	}
	syncer := replication.NewBinlogSyncer(cfg)

	// Start sync with specified binlog file and position
	streamer, _ := syncer.StartSync(mysql.Position{})
	// or you can start a gtid replication like
	// streamer, _ := syncer.StartSyncGTID(gtidSet)
	// the mysql GTID set likes this "de278ad0-2106-11e4-9f8e-6edd0ca20947:1-2"
	// the mariadb GTID set likes this "0-1-100"
	for {
		ev, _ := streamer.GetEvent(context.Background())
		switch ev.Event.(type) {
		case *replication.RowsEvent:
			Dump(ev, os.Stdout)
		}
	}
}

func Dump(be *replication.BinlogEvent, w io.Writer) {
	// fmt.Fprintf(w, "TableID: %d\n", e.TableID)
	// fmt.Fprintf(w, "Flags: %d\n", e.Flags)
	// fmt.Fprintf(w, "Column count: %d\n", e.ColumnCount)
	e := be.Event.(*replication.RowsEvent)
	if string(e.Table.Table) != "holes" {
		return
	}
	var action string
	switch be.Header.EventType {
	case replication.WRITE_ROWS_EVENTv1, replication.WRITE_ROWS_EVENTv2:
		action = "insert"
	case replication.DELETE_ROWS_EVENTv1, replication.DELETE_ROWS_EVENTv2:
		action = "delete"
	case replication.UPDATE_ROWS_EVENTv1, replication.UPDATE_ROWS_EVENTv2:
		action = "update"
	}
	fmt.Fprintf(w, "%s table [%s]", action, string(e.Table.Table))
	// fmt.Fprintf(w, "cloumn name\t")

	for _, colname := range e.Table.ColumnName {
		fmt.Fprintf(w, "%s\t", string(colname))
	}
	fmt.Fprintf(w, "\n_______________________________________________________________\n")
	for _, rows := range e.Rows {
		for _, d := range rows {
			if _, ok := d.([]byte); ok {
				fmt.Fprintf(w, "%q\t", d)
			} else {
				fmt.Fprintf(w, "%#v\t", d)
			}
		}
		fmt.Fprintf(w, "\n^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^\n")
	}
	fmt.Fprintf(w, "\n_______________________________________________________________\n")
}

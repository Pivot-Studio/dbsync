package client

import (
	"dbsync/model"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Rican7/conjson"
	"github.com/Rican7/conjson/transform"
	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/schema"
	"github.com/sirupsen/logrus"
)

const (
	mysqlTimeFormat string = "2006-01-02 15:04:05"
	mysqlDateFormat string = "2006-01-02"
)

type Model interface{}

func Build(before, after Model, msg []byte) error {
	var m model.RowRequest
	err := json.Unmarshal(msg, &m)
	if err != nil {
		logrus.Errorf("[Build] pos json parse err %v", err)
	}
	switch m.Action {
	case canal.InsertAction:
		err = buildInsertMsg(after, m)
	case canal.DeleteAction:
		err = buildDeleteMsg(before, m)
	case canal.UpdateAction:
		err = buildUpdateMsg(before, after, m)
	}
	if err != nil {
		return fmt.Errorf("[Build] %s message err %v", m.Action, err)
	}
	return nil
}
func buildInsertMsg(dest interface{}, msg model.RowRequest) error {
	destMap := make(map[string]interface{})
	for k, c := range msg.Column {
		destMap[c.Name] = makeReqColumnData(&c, msg.AfterData[k])
	}
	return mapToDest(&destMap, dest)
}
func buildDeleteMsg(dest interface{}, msg model.RowRequest) error {
	destMap := make(map[string]interface{})
	for k, c := range msg.Column {
		destMap[c.Name] = makeReqColumnData(&c, msg.BeforeData[k])
	}
	return mapToDest(&destMap, dest)
}
func mapToDest(m *map[string]interface{}, dest interface{}) error {

	b, err := json.Marshal(m)
	if err != nil {
		logrus.Errorf("[mapToDest] map json parse err %v", err)
		return err
	}
	json.Unmarshal(
		b,
		conjson.NewUnmarshaler(dest, transform.ConventionalKeys()),
	)
	if err != nil {
		logrus.Errorf("[mapToDest] dest json parse err %v", err)
		return err
	}
	return nil
}
func buildUpdateMsg(before, after Model, msg model.RowRequest) error {
	beforeMap, afterMap := make(map[string]interface{}), make(map[string]interface{})
	for k, c := range msg.Column {
		beforeMap[c.Name] = makeReqColumnData(&c, msg.BeforeData[k])
		afterMap[c.Name] = makeReqColumnData(&c, msg.AfterData[k])
	}
	err := mapToDest(&beforeMap, before)
	if err != nil {
		logrus.Errorf("[buildUpdateMsg] set map to dest err %v", err)
	}
	err = mapToDest(&afterMap, after)
	if err != nil {
		logrus.Errorf("[buildUpdateMsg] set map to dest err %v", err)
	}
	return nil
}
func makeReqColumnData(col *schema.TableColumn, value interface{}) interface{} {
	switch col.Type {
	case schema.TYPE_ENUM:
		switch value := value.(type) {
		case int64:
			// for binlog, ENUM may be int64, but for dump, enum is string
			eNum := value - 1
			if eNum < 0 || eNum >= int64(len(col.EnumValues)) {
				// we insert invalid enum value before, so return empty
				logrus.Warnf("[makeReqColumnData] invalid binlog enum index %d, for enum %+v", eNum, col.EnumValues)
				return ""
			}

			return col.EnumValues[eNum]
		}
	case schema.TYPE_SET:
		switch value := value.(type) {
		case int64:
			// for binlog, SET may be int64, but for dump, SET is string
			bitmask := value
			sets := make([]string, 0, len(col.SetValues))
			for i, s := range col.SetValues {
				if bitmask&int64(1<<uint(i)) > 0 {
					sets = append(sets, s)
				}
			}
			return strings.Join(sets, ",")
		}
	case schema.TYPE_BIT:
		switch value := value.(type) {
		case string:
			// for binlog, BIT is int64, but for dump, BIT is string
			// for dump 0x01 is for 1, \0 is for 0
			if value == "\x01" {
				return int64(1)
			}

			return int64(0)
		}
	case schema.TYPE_STRING:
		switch value := value.(type) {
		case []byte:
			return string(value[:])
		}
	case schema.TYPE_JSON:
		var f interface{}
		var err error
		switch v := value.(type) {
		case string:
			err = json.Unmarshal([]byte(v), &f)
		case []byte:
			err = json.Unmarshal(v, &f)
		}
		if err == nil && f != nil {
			return f
		}
	case schema.TYPE_DATETIME, schema.TYPE_TIMESTAMP:
		switch v := value.(type) {
		case string:
			loc, err := time.LoadLocation("UTC")
			if err != nil {
				return err
			}
			vt, err := time.ParseInLocation(mysqlTimeFormat, string(v), loc)
			if err != nil || vt.IsZero() { // failed to parse date or zero date
				return nil
			}
			return vt
		}
	case schema.TYPE_DATE:
		switch v := value.(type) {
		case string:
			vt, err := time.Parse(mysqlDateFormat, string(v))
			if err != nil || vt.IsZero() { // failed to parse date or zero date
				return nil
			}
			return vt.Format(mysqlDateFormat)
		}
	}

	return value
}

package conf

import (
	"io/ioutil"
	"os"

	json "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
)

var C *Config

type Config struct {
	Version string
	Redis   struct {
		Host     string
		Port     string
		Password string
	}
	Mysql struct {
		Host     string
		Port     string
		Username string
		Password string
	}
	Bolt struct {
		File   string
		Bucket string
	}
	Stan struct {
		Host      string
		Port      string
		ClusterID string
		ClientID  string
	}
	Storage struct {
		DaoName string
		PosKey  string
	}
	MQName string
}

func init() {
	C = &Config{}
	f, err := os.Open("config.json")
	defer f.Close()
	if err != nil {
		logrus.Fatalf("[init] open config error:%+v", err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		logrus.Fatalf("[init] read config error:%+v", err)
	}
	err = json.Unmarshal(b, C)
	if err != nil {
		logrus.Fatalf("[init] unmarshal config error:%+v", err)
	}
}

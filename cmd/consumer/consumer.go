package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"dbsync/client"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Holes struct {
	ID                 uint
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          gorm.DeletedAt `gorm:"index"`
	HoleId             uint           `gorm:"primarykey"`
	OwnerEmail         string
	Content            string `gorm:"type:varchar(1037)"`
	ImageUrl           string
	CreatedTimestamp   int64
	CreatedIp          string
	LastReplyTimestamp int64
	ThumbupNum         int
	ReplyNum           int
	FollowNum          int
	PvNum              int
	IsDeleted          bool
	ForestId           int //树洞所属的小树林
}

func HoleTest(msg []byte) {
	var holeBefore Holes
	var holeAfter Holes
	client.Build(&holeBefore, &holeAfter, msg)
	logrus.Infof("before %v,after %v", holeBefore, holeAfter)
}
func main() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	c, err := client.NewClient(client.Config{
		ClusterID: "stan",
		ClientID:  "client1",
		Host:      "nats://nats.default.svc.cluster.local",
		Port:      "4222",
	})
	if err != nil {
		logrus.Fatalf("init client err %v", err)
	}
	c.Register(Holes{}, HoleTest)
	err = c.Run()
	if err != nil {
		logrus.Fatalf("start consumer err %v", err)
	}
	n := <-sc
	logrus.Infof("receive signal %v, closing", n)
	err = c.Close()
	if err != nil {
		logrus.Fatalf("shutdown err consumer err %v", err)
	}
}

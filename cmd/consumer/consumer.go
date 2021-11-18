package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Pivot-Studio/dbsync/client"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Replies struct {
	gorm.Model
	Reply_user string `gorm:"index"`
	Content    string
	ImageUrl   string
	//refer to the rank in certain hole
	LocalReplyId   uint
	PostAliasIndex uint
	HoleId         uint
	ThumbupNum     uint
	IsDeleted      bool
	//reply to an another reply
	ReplyTo int
}
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
	if holeBefore == (Holes{}) {
		logrus.Infof("[hole] new hole %v", holeAfter)
	}
}
func ReplyTest(msg []byte) {
	var replyBefore Replies
	var replyAfter Replies
	client.Build(&replyBefore, &replyAfter, msg)
	if replyBefore == (Replies{}) {
		logrus.Infof("[hole] new hole %v", replyAfter)
	}
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
	c.Register(Replies{}, ReplyTest)
	n := <-sc
	logrus.Infof("receive signal %v, closing", n)
	err = c.Close()
	if err != nil {
		logrus.Fatalf("shutdown err consumer err %v", err)
	}
}

package main

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Hole struct {
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

func main() {
	var h Hole
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		logrus.Error(err)
	}
	var hFind Hole
	h = Hole{
		ID:               3,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		HoleId:           1038,
		OwnerEmail:       "tan@qq.com",
		Content:          "今天是个好日子",
		CreatedIp:        "127.0.0.1",
		CreatedTimestamp: time.Now().In(loc).Unix(),
		FollowNum:        1000,
		PvNum:            199,
		ReplyNum:         1212,
		IsDeleted:        false,
		ForestId:         77,
	}
	connStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&charset=utf8mb4,utf8",
		"canal", "canal", "127.0.0.1", "canal_test")
	db, err := gorm.Open(mysql.Open(connStr), &gorm.Config{})
	if err != nil {
		logrus.Fatal(err)
	}
	err = db.Create(&h).Error
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Infof("creat h %v", h)
	db.Where("id=?", 3).Find(&hFind)
	logrus.Warnf("find %v", hFind)
}

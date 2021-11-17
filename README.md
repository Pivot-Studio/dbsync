## Consumer Usage
``` go
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
	ForestId           int 
}

func HoleTest(msg []byte) error {
	var holeBefore Holes
	var holeAfter Holes
	client.Build(&holeBefore, &holeAfter, msg)
	logrus.Infof("before %v,after %v", holeBefore, holeAfter)
	return nil
}
func main() {
    // 这里的consemer1是消费者的唯一标识ID，切勿重复
    // 一个消费者可以订阅多张表，默认广播模式
    c, err := client.NewClient("consemer1")
    if err != nil {
	    logrus.Fatalf("init client err %v", err)
    }
    // 在接收到消息后会调用回调函数HoleTest
    c.Register(Holes{}, HoleTest)
    err = c.Run()
    if err != nil {
        logrus.Fatalf("start consumer err %v", err)
    }
}
```
完整的例子在[cmd/consumer/consumer.go](./cmd/consumer/consumer.go)

**需要注意的点**
- c.Register里面传入的model和回调函数中反序列化的model必须是同一个类型，因为go没泛型，在这里做不了类型检查，使用错误会panic
- client.Build(&holeBefore, &holeAfter, msg),这里必须传引用，因为用的interface{}所以也做不了类型检查
- 数据插入时,holeBefore是空,holeAfter是插入的值
- 数据删除时,holeAfter是空,holeBefore是删除的值
- 数据更改时,holeBefore是更改前的值,hoelAfter是更改后的值
- 具有time.Time类型的字段,如CreatAt这种,会绑定UTC时间,和gorm一致,如果实际使用得加8h
- 为了保证消息消费的顺序,目前的逻辑是如果某条消息消费错误会返回状态SuspendCurrentQueueAMoment，而不是放入延时队列,并且同一个表的数据在一个队列当中。
## Producer Usage
```bash
go run cmd/producer/producer.go
```
完整例子在[cmd/producer/producer.go](./cmd/producer/producer.go)
**需要注意的点**
- 持久化逻辑：目前binlog的postion持久化发生在以下几种情况
    - 每隔3s把累积消息一次性投递到mq后
    - 在小于3s内消息累积超过1024条
    - binlog文件改变后（比如从mysql-bin.000010变成mysql-bin.000011）
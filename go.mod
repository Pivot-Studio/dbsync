module github.com/Pivot-Studio/dbsync

go 1.15

require (
	github.com/apache/rocketmq-client-go v1.2.4
	github.com/apache/rocketmq-client-go/v2 v2.1.0
	github.com/go-mysql-org/go-mysql v1.3.0
	github.com/sirupsen/logrus v1.4.1
)

replace github.com/go-mysql-org/go-mysql v1.3.0 => github.com/Pivot-Studio/go-mysql v1.3.1

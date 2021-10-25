package main

import (
	"os"
	"os/signal"
	"syscall"

	"dbsync/client"

	"github.com/sirupsen/logrus"
)

type Holes struct {
	Id       int
	Contents string
}

func HoleTest(msg []byte) error {
	var holeBefore Holes
	var holeAfter Holes
	client.Build(&holeBefore, &holeAfter, msg)
	logrus.Infof("before %v,after %v", holeBefore, holeAfter)
	return nil
}
func main() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	c, err := client.NewClient("consemer 3")
	if err != nil {
		logrus.Fatalf("init client err %v", err)
	}
	c.Register(Holes{}, HoleTest)
	done := make(chan struct{}, 1)
	go func() {
		err = c.Run()
		if err != nil {
			logrus.Fatalf("start consumer err %v", err)
		}
		done <- struct{}{}
	}()
	n := <-sc
	logrus.Infof("receive signal %v, closing", n)
	err = c.Stop()
	if err != nil {
		logrus.Fatalf("shutdown err consumer err %v", err)
	}
	<-done

}

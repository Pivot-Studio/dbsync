package main

import (
	"dbsync/river"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

func main() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	r, err := river.NewRiver()
	if err != nil {
		logrus.Fatal(err)
	}
	done := make(chan struct{}, 1)
	go func() {
		r.Run()
		done <- struct{}{}
	}()

	select {
	case n := <-sc:
		logrus.Infof("receive signal %v, closing", n)
	case <-r.Ctx().Done():
		logrus.Infof("context is done with %v, closing", r.Ctx().Err())
	}

	r.Close()
	<-done
}

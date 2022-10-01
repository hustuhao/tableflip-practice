package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/cloudflare/tableflip"
	"github.com/robfig/cron/v3"
)

var (
	x int64
)

func main() {
	var (
		//listenAddr = flag.String("listen", "localhost:8080", "`Address` to listen on")
		pidFile = flag.String("pid-file", "./app.pid", "`Path` to pid file")
	)
	log.SetPrefix(fmt.Sprintf("%d ", os.Getpid()))

	logFilePath := "../../log/cron/"
	os.Mkdir(logFilePath, os.ModePerm)

	logFile, err := os.OpenFile(logFilePath+time.Now().Format("2006-01-02")+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile)
	log.Printf("pidFile:%s", *pidFile)
	upg, err := tableflip.New(tableflip.Options{
		PIDFile: *pidFile,
	})
	if err != nil {
		panic(err)
	}
	defer upg.Stop()

	// Do an upgrade on SIGHUP
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGUSR2)
		for range sig {
			log.Printf("reciece sigusr2")
			err := upg.Upgrade()
			if err != nil {
				log.Println("Upgrade failed:", err)
			}
		}
	}()

	// cron 这个库的不是0秒开始，按进程启动时间固定间隔运行
	// 更新了v3库 需要额外支持"秒"
	c := cron.New(cron.WithSeconds())
	c.AddFunc("@every 1s", cronTask)

	c.Start()

	log.Printf("ready")
	if err := upg.Ready(); err != nil {
		panic(err)
	}

	log.Printf("exit1")
	<-upg.Exit()
	log.Printf("exit2")

	// Make sure to set a deadline on exiting the process
	// after upg.Exit() is closed. No new upgrades can be
	// performed if the parent doesn't exit.
	ctx := c.Stop()
	select {
	case <-ctx.Done():
		log.Println("handle cron stop.", ctx.Err())
	case <-time.After(60 * time.Second):
		log.Println("Graceful shutdown timed out")
		os.Exit(1)
	}

	log.Printf("cron end, pid: %d", os.Getpid())
}

func cronTask() {
	i := atomic.AddInt64(&x, 1)
	log.Printf("Task%d start------\n", i)
	for i := 0; i < 2; i++ {
		time.Sleep(time.Second)
		//log.Printf("I'm wokring hard.\n")
	}
	log.Printf("Task%d end------\n", i)
}

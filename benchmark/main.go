package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var reconnectCount = make(map[string]int)
var ch = make(chan string, 10)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
}

func main() {
	port := flag.Int("port", 8087, "")
	flag.Parse()

	gracefulShutdown()
	srv := gin.Default()

	srv.GET("/", func(ctx *gin.Context) {
		go count()
		remoteAddr := ctx.Request.RemoteAddr
		logrus.Info("remote address:" + remoteAddr)
		ch <- remoteAddr
		fmt.Println("query params:", ctx.Query("lang"))
		fmt.Println("headers:", ctx.GetHeader("Authorization"))
	})

	srv.Run(fmt.Sprintf(":%d", *port))
}

func gracefulShutdown() {
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		close(ch)
		time.Sleep(time.Second)
		fmt.Printf("http reconnect count:%v \n", len(reconnectCount)-1) // 連線重建次數
		os.Exit(0)
	}()
}

func count() {
	for s := range ch {
		reconnectCount[s]++
	}
}

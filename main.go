package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/agamotto-cloud/imok/pkg/config"
	"github.com/agamotto-cloud/imok/pkg/handler"
)

func main() {
	config.Initialize()
	config.C.Show()

	// 优雅关闭服务
	shutdownC := make(chan struct{})
	go func() {
		defer func() {
			shutdownC <- struct{}{}
		}()

		signalC := make(chan os.Signal, 1)
		signal.Notify(signalC, syscall.SIGINT, syscall.SIGTERM)

		<-signalC

		log.Println("Shutdown Server")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// 停止 API 服务
		if err := handler.Shutdown(ctx); err != nil {
			log.Println("Shutdown Server", err.Error())
		}
	}()

	log.Println("Good To Go")

	//启动服务
	handler.Run()

	//等待关闭所有资源
	<-shutdownC
}

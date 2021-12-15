package handler

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/agamotto-cloud/imok/pkg/config"
	"github.com/agamotto-cloud/imok/pkg/service"
	"github.com/gin-gonic/gin"
)

var (
	eng    *gin.Engine
	server *http.Server
	once   sync.Once
)

func init() {
	once.Do(func() {
		eng = gin.Default()
	})
}

//Regist 注册 API
func Regist() {

	v1 := eng.Group("/v1")

	v1.Handle("GET", "/healz", service.Healz)

	v1.Handle("GET", "/regist", service.Regist)
}

// Run 启动Http服务
func Run() {
	Regist()
	server = &http.Server{
		Addr:    config.C.Server.Listen,
		Handler: eng,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Println("Server Run", err.Error())
	}
}

// Shutdown Server
func Shutdown(ctx context.Context) error {
	return server.Shutdown(ctx)
}

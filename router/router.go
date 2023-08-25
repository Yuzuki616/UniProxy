package router

import (
	"UniProxy/handle"
	"UniProxy/middleware"
	"fmt"
	"github.com/gin-gonic/gin"
)

var engine *gin.Engine

func Init() {
	gin.SetMode(gin.ReleaseMode)
	engine = gin.New()
	engine.Use(middleware.Logger, gin.Recovery())
}

func loadRoute() {
	// status
	engine.POST("initParams", handle.InitParams)
	engine.GET("getStatus", handle.GetStatus)
	// servers
	engine.GET("getServers", handle.GetServers)
	// proxy
	engine.POST("startUniProxy", handle.StartUniProxy)
	engine.GET("stopUniProxy", handle.StopUniProxy)
	engine.GET("setSystemProxy", handle.SetSystemProxy)
	engine.GET("clearSystemProxy", handle.ClearSystemProxy)
}

func Start(host string, port int) error {
	loadRoute()
	err := engine.Run(fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return err
	}
	return nil
}

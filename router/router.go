package router

import (
	"UniProxy/geo"
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
	engine.GET("geosite.db", func(c *gin.Context) {
		c.Header("content-disposition", "attachment; filename=\"geosite.db\"")
		c.Data(200, "application/octet-stream", geo.Site)
	})
	engine.GET("geoip.db", func(c *gin.Context) {
		c.Header("content-disposition", "attachment; filename=\"geoip.db\"")
		c.Data(200, "application/octet-stream", geo.Ip)
	})
}

func Start(host string, port int) error {
	loadRoute()
	err := engine.Run(fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return err
	}
	return nil
}

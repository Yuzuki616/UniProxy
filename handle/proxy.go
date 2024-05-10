package handle

import (
	"UniProxy/proxy"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type StartUniProxyRequest struct {
	Tag        string `json:"tag"`
	Uuid       string `json:"uuid"`
	GlobalMode bool   `json:"global_mode"`
}

func StartUniProxy(c *gin.Context) {
	p := StartUniProxyRequest{}
	err := c.ShouldBindJSON(&p)
	if err != nil {
		c.JSON(400, Rsp{
			Success: false,
		})
		return
	}
	proxy.GlobalMode = p.GlobalMode
	err = proxy.StartProxy(p.Tag, p.Uuid, servers[p.Tag])
	if err != nil {
		log.WithField("err", err).Error("start proxy failed")
		c.JSON(400, Rsp{
			Success: false,
			Message: err.Error(),
		})
		return
	}
	c.JSON(200, Rsp{
		Success: true,
		Message: "ok",
		Data: StatusData{
			Inited:      inited,
			Running:     proxy.Running,
			GlobalMode:  proxy.GlobalMode,
			SystemProxy: proxy.SystemProxy,
		},
	})
	return
}

func StopUniProxy(c *gin.Context) {
	if proxy.Running {
		proxy.StopProxy()
	}
	c.JSON(200, Rsp{
		Success: true,
		Message: "ok",
	})
}

func SetSystemProxy(c *gin.Context) {
	c.JSON(200, Rsp{
		Success: true,
		Message: "ok",
	})
}

func ClearSystemProxy(c *gin.Context) {
	err := proxy.ClearSystemProxy()
	if err != nil {
		c.JSON(200, Rsp{
			Success: false,
			Message: err.Error(),
		})
	}
	c.JSON(200, Rsp{
		Success: true,
		Message: "ok",
	})
}

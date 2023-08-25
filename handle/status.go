package handle

import (
	"V2bProxy/common/encrypt"
	"V2bProxy/proxy"
	"V2bProxy/v2b"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
)

type initParamsRequest struct {
	MixedPort int    `json:"mixed_port"`
	AppName   string `json:"app_name"`
	Url       string `json:"url"`
	Token     string `json:"token"`
	License   string `json:"license"`
	UserPath  string `json:"user_path"`
}

var inited bool

func InitParams(c *gin.Context) {
	p := initParamsRequest{}
	err := c.ShouldBindJSON(&p)
	if err != nil {
		c.JSON(400, &Rsp{Success: false, Message: err.Error()})
		return
	}
	if encrypt.Sha([]byte(encrypt.Sha([]byte(p.Url))+"1145141919")) != p.License {
		c.JSON(400, &Rsp{Success: false})
		return
	}
	f, err := os.OpenFile(path.Join(p.UserPath, "uniproxy.log"), os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		c.JSON(400, &Rsp{Success: false, Message: err.Error()})
		return
	}
	log.SetOutput(f)
	v2b.Init(p.Url, p.Token)
	proxy.InPort = p.MixedPort
	proxy.DataPath = p.UserPath
	inited = true
	c.JSON(200, &Rsp{Success: true})
}

func GetStatus(c *gin.Context) {
	c.JSON(200, &Rsp{
		Success: true,
		Data: StatusData{
			Inited:      inited,
			Running:     proxy.Running,
			GlobalMode:  proxy.GlobalMode,
			SystemProxy: proxy.SystemProxy,
		},
	})
}

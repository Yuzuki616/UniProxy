package middleware

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"time"
)

func Logger(c *gin.Context) {
	path := c.Request.URL.Path
	raw := c.Request.URL.RawQuery
	// Start timer
	start := time.Now()
	log.Infof("[%s:%d] | %s %s",
		c.Request.Method, c.Writer.Status(), c.ClientIP()+" |", path)
	// Process request
	c.Next()
	latency := time.Now().Sub(start)
	if latency > time.Minute {
		latency = latency - latency%time.Second
	}
	if raw != "" {
		path = path + "?" + raw
	}
	log.Infof("[%s:%d] %s | %s %s",
		c.Request.Method, c.Writer.Status(), latency, c.ClientIP()+" |", path)
}

package handle

import (
	"UniProxy/common/balance"
	"UniProxy/conf"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

var urlBalance *balance.List[string]

func ReverseProxy(c *gin.Context) {
	if urlBalance == nil {
		urlBalance = balance.New[string](conf.C.Api.Balance, conf.C.Api.Baseurl)
	}
	if len(conf.C.Api.Balance) == 0 {
		return
	}
	u := urlBalance.Next()
	target, err := url.Parse(u)
	if err != nil {
		log.WithField("err", err).Error("parse url failed")
		c.JSON(400, Rsp{
			Success: false,
		})
	}
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Director = func(req *http.Request) {
		req.Host = target.Host
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		if len(req.Form) != 0 {
			r := strings.NewReader(req.Form.Encode())
			req.ContentLength = int64(r.Len())
			req.Body = io.NopCloser(r)
		}
	}
	proxy.ServeHTTP(c.Writer, c.Request)
}

package v2b

import (
	"UniProxy/common/balance"
	"encoding/json"
	"errors"
	"github.com/avast/retry-go/v4"
	"github.com/go-resty/resty/v2"
	"time"
)

var (
	clients *balance.List[*resty.Client]
	etag    string
)

func Init(b string, url []string, auth string) {
	cs := make([]*resty.Client, len(url))
	for i, u := range url {
		cs[i] = resty.New().
			SetTimeout(time.Second*40).
			SetQueryParam("auth_data", auth).
			SetBaseURL(u).
			SetRetryCount(3).
			SetRetryWaitTime(3 * time.Second)
	}
	clients = balance.New[*resty.Client](b, cs)
}

type ServerFetchRsp struct {
	Data []ServerInfo `json:"data"`
}

type ServerInfo struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Host        string `json:"host"`
	Port        int    `json:"port"`
	Network     string `json:"network"`
	Type        string `json:"type"`
	Cipher      string `json:"cipher"`
	Tls         int    `json:"tls"`
	Flow        string `json:"flow"`
	TlsSettings struct {
		ServerName    string `json:"serverName"`
		ServerPort    string `json:"server_port"`
		AllowInsecure string `json:"allowInsecure"`
		RealityDest   string `json:"server_name"`
		ShortId       string `json:"short_id"`
		PublicKey     string `json:"public_key"`
	} `json:"tls_settings"`
	NetworkSettings struct {
		Path       string      `json:"path"`
		Headers    interface{} `json:"headers"`
		ServerName string      `json:"serverName"`
	} `json:"networkSettings"`
	CreatedAt     interface{} `json:"created_at"`
	AllowInsecure int         `json:"insecure"`
	LastCheckAt   interface{} `json:"last_check_at"`
	Tags          interface{} `json:"tags"`
	UpMbps        int         `json:"up_mbps"`
	ServerName    string      `json:"server_name"`
	ServerKey     string      `json:"server_key"`
	DownMbps      int         `json:"down_mbps"`
}

func GetServers() ([]ServerInfo, error) {
	var r *resty.Response
	err := retry.Do(func() error {
		c := clients.Next()
		rsp, err := c.R().
			SetHeader("If-None-Match", etag).
			Get("api/v1/user/server/fetch")
		if err != nil {
			return err
		}
		if rsp.StatusCode() == 304 {
			return nil
		}
		etag = rsp.Header().Get("ETag")
		if rsp.StatusCode() != 200 {
			return nil
		}
		r = rsp
		return nil
	}, retry.Attempts(3))
	if err != nil {
		return nil, err
	}
	if r.StatusCode() == 304 {
		return nil, nil
	}
	rsp := &ServerFetchRsp{}
	err = json.Unmarshal(r.Body(), rsp)
	if err != nil {
		return nil, err
	}
	if len(rsp.Data) == 0 {
		return nil, errors.New("no servers")
	}
	return rsp.Data, nil
}

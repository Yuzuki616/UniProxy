package v2b

import (
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
	"time"
)

var (
	client *resty.Client
)

func Init(url, auth string) {
	client = resty.New().
		SetTimeout(time.Second*40).
		SetQueryParam("auth_data", auth).
		SetBaseURL(url).
		SetRetryCount(3).
		SetRetryWaitTime(3 * time.Second)
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
	CreatedAt     int         `json:"created_at"`
	AllowInsecure int         `json:"insecure"`
	LastCheckAt   string      `json:"last_check_at"`
	Tags          interface{} `json:"tags"`
	ServerName    string      `json:"server_name"`
	ServerKey     string      `json:"server_key"`
	UpMbps        int         `json:"up_mbps"`
	DownMbps      int         `json:"down_mbps"`
}

func GetServers() ([]ServerInfo, error) {
	r, err := client.R().
		Get("api/v1/user/server/fetch")
	if err != nil {
		return nil, err
	}
	if r.StatusCode() == 304 {
		return nil, nil
	}
	if r.StatusCode() != 200 {
		return nil, errors.New(r.String())
	}
	client.SetHeader("If-None-Match", r.Header().Get("ETag"))
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

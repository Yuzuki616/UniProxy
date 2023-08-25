package proxy

import (
	"V2bProxy/common/sysproxy"
	"V2bProxy/v2b"
	box "github.com/sagernet/sing-box"
)

var (
	Running     bool
	SystemProxy bool
	GlobalMode  bool
	TunMode     bool
	InPort      int
	DataPath    string
)

var client *box.Box

func StartProxy(tag string, uuid string, server *v2b.ServerInfo) error {
	if !Running {
		StopProxy()
	}
	SystemProxy = true
	c, err := GetSingBoxConfig(uuid, server)
	if err != nil {
		return err
	}
	client, err = box.New(box.Options{Options: c})
	if err != nil {
		return err
	}
	err = client.Start()
	if err != nil {
		return err
	}
	Running = true
	return nil
}

func StopProxy() {
	if Running {
		client.Close()
		Running = false
	}
}

func ClearSystemProxy() error {
	if Running {
		client.Close()
		Running = false
		return nil
	}
	sysproxy.ClearSystemProxy()
	return nil
}

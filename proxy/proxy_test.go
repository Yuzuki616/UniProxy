package proxy

import (
	"UniProxy/v2b"
	"testing"
)

func TestStartProxy(t *testing.T) {
	v2b.Init("http://127.0.0.1:1022", "xxxxxxxx")
	s, _ := v2b.GetServers()
	t.Log(s[0])
	InPort = 1151
	GlobalMode = true
	t.Log(StartProxy("test", "xxxxxx", &s[0]))
	select {}
}

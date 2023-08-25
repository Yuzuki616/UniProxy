package v2b

import (
	"testing"
)

func TestGetServers(t *testing.T) {
	Init("http://127.0.0.1", "xxxxxxxxx")
	t.Log(GetServers())

}

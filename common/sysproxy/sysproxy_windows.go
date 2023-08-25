package sysproxy

import "github.com/sagernet/sing/common/wininet"

func ClearSystemProxy() error {
	return wininet.ClearSystemProxy()
}

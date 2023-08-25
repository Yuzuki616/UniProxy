package sysproxy

import (
	"github.com/sagernet/sing/common/shell"
	"strings"
)

func ClearSystemProxy() error {
	o, err := shell.Exec("networksetup", "-listallnetworkservices").ReadOutput()
	if err != nil {
		return err
	}
	nets := strings.Split(o, "\n")
	for i := range nets {
		shell.Exec("networksetup", "-setsocksfirewallproxystate", nets[i], "off").Attach().Run()
		shell.Exec("networksetup", "-setwebproxystate", nets[i], "off").Attach().Run()
		shell.Exec("networksetup", "-setsecurewebproxystate", nets[i], "off").Attach().Run()
	}
	return nil
}

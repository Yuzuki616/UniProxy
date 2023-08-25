//go:build !(windows || darwin)

package sysproxy

import (
	"os"
)

func ClearSystemProxy() error {
	return os.ErrInvalid
}

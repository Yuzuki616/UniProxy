package file

import (
	"fmt"
	"io"
	"os"
)

func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}

func Copy(path1 string, path2 string) error {
	f, err := os.Open(path1)
	if err != nil {
		return fmt.Errorf("open file1 error: %s", err)
	}
	f2, err := os.OpenFile(path2, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0744)
	if err != nil {
		return fmt.Errorf("open file2 error: %s", err)
	}
	_, err = io.Copy(f2, f)
	if err != nil {
		return fmt.Errorf("stream copy error: %s", err)
	}
	return nil
}

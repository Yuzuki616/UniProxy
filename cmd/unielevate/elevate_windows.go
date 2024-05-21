package main

import (
	"os"
	"os/exec"
)

func ExecElevateCommand(cmd string) (*exec.Cmd, error) {
	c := exec.Command("elevate.exe", "-k", cmd)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err := c.Start()
	if err != nil {
		return nil, err
	}
	return c, nil
}

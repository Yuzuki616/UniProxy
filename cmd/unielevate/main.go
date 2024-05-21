package main

import (
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const defaultName = "uniproxy_base.exe"

func main() {
	var cmdArgs []string
	if len(os.Args) > 1 {
		if strings.HasPrefix(os.Args[1], "-") {
			cmdArgs = make([]string, 0, len(os.Args)+2)
			cmdArgs = append(cmdArgs, defaultName)
			cmdArgs = append(cmdArgs, os.Args[1:]...)
			cmdArgs = append(cmdArgs, "-tun")
		} else {
			cmdArgs = os.Args[1:]
			cmdArgs = append(cmdArgs, "-tun")
		}
	} else {
		cmdArgs = os.Args[1:]
		cmdArgs = append(cmdArgs, "-tun")
	}
	c, err := ExecElevateCommand(strings.Join(os.Args, " "))
	if err != nil {
		return
	}
	defer c.Wait()

	// waiting signal
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	<-s
}

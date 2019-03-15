package main

import (
	"etcdInLeft/lib"
)

func main() {
	cmd_runner := &lib.Commands{
		CmdServer,
	}

	cmd_runner.Run()
}

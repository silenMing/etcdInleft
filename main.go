package main

import (
	"github.com/silenMing/etcdInleft/lib"
)

func main() {
	cmd_runner := &lib.Commands{
		CmdServer,
	}

	cmd_runner.Run()
}

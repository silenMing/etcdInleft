package main

import (
	"github.com/silenMing/etcdInleft/lib"
)

func main() {
	cmdRunner := &lib.Commands{
		CmdServer,
	}

	cmdRunner.Run()
}

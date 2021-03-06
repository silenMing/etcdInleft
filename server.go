package main

import (
	"runtime"
	"log"
	"github.com/silenMing/etcdInleft/lib"
)

var CmdServer = &lib.Command{
	UsageLine: "server",
	Short:     "服务管理",
	Long: `start  启动服务`,
}


func init() {
	CmdServer.Run = server
}

func server(cmd *lib.Command, args []string)  int{
	runtime.GOMAXPROCS(runtime.NumCPU())
	if len(args) < 1 {
		log.Println("缺少参数")
		return 1
	}

	switch args[0] {
	case "start":
		lib.ConnectEtcd("leftGetId", Cfg.EtcdAddr, Cfg.MyAddr)

	}

	return  0
}

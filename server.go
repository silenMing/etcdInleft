package main

import (
	"runtime"
	"log"
	"etcdInLeft/lib"
)

var cmdServer = &lib.Command{
	UsageLine: "server",
	Short:     "服务管理",
	Long: `start  启动服务`,
}


func init() {
	cmdServer.Run = server
}

func server(cmd *lib.Command, args []string)  int{
	runtime.GOMAXPROCS(runtime.NumCPU())
	if len(args) < 1 {
		log.Println("缺少参数")
		return 1
	}




	return  0
}

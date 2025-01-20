package main

import (
	"flag"

	"github.com/golang/glog"

	"github.com/Project-HAMi/dcu-dcgm/pkg/dcgm"
)

func main() {
	flag.Parse()
	defer glog.Flush()
	glog.Info("go-dcgm start ...")
	//初始化dcgm服务
	dcgm.Init()
	defer dcgm.ShutDown()

	//所有的pid列表
	dcgm.PidList()
	//获取进程信息（格式化打印信息）
	dcgm.ShowPids()
	//根据pid获取进程的名称
	dcgm.ProcessName(1)
}

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

	//硬件拓扑信息(支持json打印信息)
	dcgm.ShowWeightTopology([]int{0, 1, 2}, true)
	dcgm.ShowWeightTopology([]int{0, 1, 2}, false)
	//基于跳数显示硬件拓扑信息(支持json打印信息)
	dcgm.ShowHopsTopology([]int{0, 1, 2}, false)
	dcgm.ShowHopsTopology([]int{0, 1, 2}, true)
	//基于链接类型的硬件拓扑信息(支持json打印信息)
	dcgm.ShowTypeTopology([]int{0, 1, 2}, true)
	dcgm.ShowTypeTopology([]int{0, 1, 2}, false)
	//numa节点HW拓扑信息
	dcgm.ShowNumaTopology([]int{0, 1, 2})
	//显示硬件拓扑信息,包括权重、跳数、链接类型以及NUMA节点信息
	dcgm.ShowHwTopology([]int{0, 1, 2})
}

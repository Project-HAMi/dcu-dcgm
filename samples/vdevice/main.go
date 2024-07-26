package main

import (
	"flag"

	"github.com/golang/glog"

	"g.sugon.com/das/dcgm-dcu/pkg/dcgm"
)

func main() {
	flag.Parse()
	defer glog.Flush()
	glog.Info("go-dcgm start ...")
	dcgm.Init()
	defer dcgm.ShutDown()
	//dcgm.AllDeviceInfos()
	//dcgm.VDeviceCount()
	//dcgm.DeviceRemainingInfo(0)
	//dcgm.DeviceRemainingInfo(1)

}

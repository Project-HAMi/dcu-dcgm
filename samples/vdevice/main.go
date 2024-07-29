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
	dcgm.VDeviceCount()
	//dcgm.DeviceRemainingInfo(0)
	//dcgm.DeviceRemainingInfo(1)

	dcgm.CreateVDevices(0, 2, []int{4, 4}, []int{1024, 2048})
	//dcgm.DestroyVDevice(1)
	//dcgm.DestroySingleVDevice(0)
	//dcgm.UpdateSingleVDevice(5, 20, 8589934592)
	//dcgm.StopVDevice(0)
	//dcgm.StartVDevice(0)
	dcgm.EncryptionVMStatus()
	//dcgm.SetEncryptionVMStatus(true)
}

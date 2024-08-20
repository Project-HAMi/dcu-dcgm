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
	//获取物理设备信息列表
	dcgm.DeviceInfos()
	//DCU设备数量
	dcgm.DeviceCount()
	//虚拟设备信息
	dcgm.VDeviceSingleInfo(0)
	//虚拟设备总数量
	dcgm.VDeviceCount()
	//获取所有物理设备及其虚拟设备的信息列表
	dcgm.AllDeviceInfos()
	//销毁指定虚拟设备
	dcgm.DestroySingleVDevice(1)
	//销毁指定物理设备上的所有虚拟设备
	dcgm.DestroyVDevice(1)
	//更新虚拟设备资源
	dcgm.UpdateSingleVDevice(2, 10, 2048)
	//获取物理设备剩余资源
	dcgm.DeviceRemainingInfo(1)
	dcgm.CreateVDevices(0, 2, []int{10, 10}, []int{1024, 1024})
	dcgm.GetDeviceInfo(0)
	dcgm.GetDeviceByDvInd(0)
	dcgm.GetDeviceByDvInd(1)
}

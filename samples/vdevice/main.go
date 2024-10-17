package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang/glog"

	"g.sugon.com/das/dcgm-dcu/pkg/dcgm"
)

func main() {
	flag.Parse()
	defer glog.Flush()
	glog.Info("go-dcgm start ...")
	dcgm.Init()
	defer dcgm.ShutDown()
	////获取物理设备信息列表
	//dcgm.DeviceInfos()
	////DCU设备数量
	//dcgm.DeviceCount()
	////虚拟设备信息
	//dcgm.VDeviceSingleInfo(0)
	////虚拟设备总数量
	//dcgm.VDeviceCount()
	////获取所有物理设备及其虚拟设备的信息列表
	//dcgm.AllDeviceInfos()
	////销毁指定虚拟设备
	//dcgm.DestroySingleVDevice(1)
	////销毁指定物理设备上的所有虚拟设备
	//dcgm.DestroyVDevice(0)
	////更新虚拟设备资源
	//dcgm.UpdateSingleVDevice(2, 10, 2048)
	////获取物理设备剩余资源
	//dcgm.DeviceRemainingInfo(1)
	//dcgm.DeviceRemainingInfo(0)
	//dcgm.CreateVDevices(0, 2, []int{10, 10}, []int{1024, 1024})
	//dcgm.GetDeviceInfo(0)
	//dcgm.GetDeviceByDvInd(0)
	//dcgm.GetDeviceByDvInd(1)
	////启动虚拟设备
	//dcgm.StartVDevice(0)
	////关闭虚拟设备
	//dcgm.StopVDevice(0)

	// 创建一个通道来监听系统中断信号
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	// 开始循环，每隔5秒调用一次 dcgm.DeviceRemainingInfo
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 调用设备1的剩余资源
			glog.Infof("==================start==============================")
			glog.Info("Calling DeviceRemainingInfo for device 1")
			dcgm.DeviceRemainingInfo(1)
			// 调用设备0的剩余资源
			glog.Info("Calling DeviceRemainingInfo for device 0")
			dcgm.DeviceRemainingInfo(0)
			glog.Info("-----------------------------------------")
			//虚拟设备数量
			glog.Info("Calling dmiGetVDeviceCount for device")
			dcgm.VDeviceCount()
			glog.Info("-----------------------------------------")
			//虚拟设备信息
			glog.Info("Calling VDeviceSingleInfo for vdevice 1")
			dcgm.VDeviceSingleInfo(1)
			glog.Info("-----------------------------------------")
			//物理设备百分比
			glog.Info("Calling DevBusyPercent for device 1")
			dcgm.DevBusyPercent(1)
			glog.Info("Calling DevBusyPercent for device 0")
			dcgm.DevBusyPercent(0)
			glog.Info("-----------------------------------------")
			//虚拟设备百分比
			glog.Info("Calling VDevBusyPercent for vdevice 1")
			dcgm.VDevBusyPercent(1)
			glog.Infof("==================end==============================")

		case <-stopChan:
			// 收到中断信号，停止程序
			glog.Info("Received stop signal, exiting...")
			return
		}
	}
	//dcgm.CreateVDevices(0, 2, []int{10, 10}, []int{1024, 1024})

}

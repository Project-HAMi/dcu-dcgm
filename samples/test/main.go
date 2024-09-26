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
	//K100 AI不支持
	//dcgm.DevPowerCapRange(0, 0)
	//dcgm.PerfDeterminismMode(0, 800)
	//dcgm.DevOdVoltInfoGet(0)
	//dcgm.DevOdVoltInfoGet(1)
	//dcgm.DevPciBandwidthSet(0, 102)
	//dcgm.DevOdVoltInfoSet(0, 0, 700, 690)
	//dcgm.DevGpuClkFreqSet(0, dcgm.RSMI_CLK_TYPE_SYS, 0b0100)
	//dcgm.DevGpuClkFreqSet(0, dcgm.RSMI_CLK_TYPE_SYS, 4)
	//dcgm.DevGpuClkFreqSet(0, dcgm.RSMI_CLK_TYPE_SYS,800)
	//dcgm.DevGpuClkFreqSet(0, dcgm.RSMI_CLK_TYPE_SYS, 1)
	//dcgm.DevGpuClkFreqSet(0, dcgm.RSMI_CLK_TYPE_DCEF, 64)
	//dcgm.DevGpuClkFreqSet(0, dcgm.RSMI_CLK_TYPE_SYS, 4096)
	//dcgm.DevOverdriveLevelSet(0, 5)
	//dcgm.DevOverdriveLevelGet(0)
	//dcgm.EccStatus(0, dcgm.RSMIGpuBlockDF)
	//dcgm.EccCount(0, dcgm.RSMIGpuBlockDF)
	//
	//dcgm.EccStatus(0, dcgm.RSMIGpuBlockPCIEBIF)
	//dcgm.EccCount(0, dcgm.RSMIGpuBlockPCIEBIF)
	//dcgm.CollectDeviceMetrics()
	dcgm.EccBlocksInfo(0)
	//dcgm.VDeviceSingleInfo(0)
	//dcgm.VDeviceSingleInfo(1)
	dcgm.AllDeviceInfos()
}

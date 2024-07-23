package main

import (
	"flag"

	"github.com/golang/glog"

	"dcgm-dcu/pkg/dcgm"
)

func main() {
	flag.Parse()
	defer glog.Flush()
	glog.Info("go-dcgm start ...")
	dcgm.Init()
	defer dcgm.ShutDown()
	dcgm.NumMonitorDevices()
	//dcgm.DevName(0)
	//dcgm.DevPciBandwidth(0)
	////dcgm.DevPerfLevelSet(0, dcgm.RSMI_DEV_PERF_LEVEL_LOW)
	//dcgm.MemoryPercent(0)
	//dcgm.DevName(1)
	//dcgm.DevSku(1)
	//dcgm.DevBrand(1)
	//dcgm.DevVendorName(1)
	//dcgm.DevVramVendor(1)
	//dcgm.DevPciBandwidth(1)
	//dcgm.MemoryPercent(1)
	//dcgm.DevGpuMetricsInfo(0)
	dcgm.CollectDeviceMetrics()
	////[vram|vis_vram|gtt
	//dcgm.MemInfo(0, "vram")
	//dcgm.MemInfo(0, "vis_vram")
	//dcgm.MemInfo(0, "gtt")
	//dcgm.ProcessName(1)
	//dcgm.ProcessName(10)
	//dcgm.ProcessName(100)
	//dcgm.PerfLevel(0)
	//dcgm.PerfLevel(1)
	//dcgm.PidByName("mm_percpu_wq")
	//dcgm.PidByName("ksoftirqd/14")
	dcgm.DeviceInfos()
}

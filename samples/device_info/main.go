package main

import (
	"log"

	"go-dcgm/pkg/dcgm"
)

func main() {
	log.Println("go-dcgm start ...")
	dcgm.Init()
	defer dcgm.ShutDown()
	dcgm.NumMonitorDevices()
	dcgm.DevName(0)
	dcgm.DevPciBandwidth(0)
	dcgm.MemoryTotal(0)
	dcgm.MemoryUsed(0)
	//dcgm.DevPerfLevelSet(0, dcgm.RSMI_DEV_PERF_LEVEL_LOW)
	dcgm.MemoryPercent(0)
	dcgm.DevName(1)
	dcgm.DevPciBandwidth(1)
	dcgm.MemoryTotal(1)
	dcgm.MemoryUsed(1)
	dcgm.MemoryPercent(1)
	dcgm.DevGpuMetricsInfo(0)

}

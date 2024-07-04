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

}

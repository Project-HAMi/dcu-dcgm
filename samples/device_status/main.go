package main

import (
	"log"

	"dcgm-dcu/pkg/dcgm"
)

func main() {
	log.Println("go-dcgm start ...")
	dcgm.Init()
	defer dcgm.ShutDown()
	dcgm.DevPerfLevelSet(0, dcgm.RSMI_DEV_PERF_LEVEL_AUTO)

	//// 示例参数
	//dvInd := int(0) // 设备索引
	//count := int(2) // 计数器数量
	//
	//// 示例利用率计数器数组
	//utilizationCounters := []dcgm.RSMIUtilizationCounter{
	//	{Type: dcgm.RSMI_COARSE_GRAIN_GFX_ACTIVITY},
	//	{Type: dcgm.RSMI_COARSE_GRAIN_MEM_ACTIVITY},
	//}
	//
	//// 调用 rsmiUtilizationCountGet 函数
	//timestamp, err := dcgm.UtilizationCount(dvInd, utilizationCounters, count)
	//if err != nil {
	//	fmt.Printf("Error: %v\n", err)
	//	return
	//}
	//// 打印结果
	//fmt.Println("Timestamp:", timestamp)
	//fmt.Println("Utilization Counters:")
	//for _, counter := range utilizationCounters {
	//	fmt.Printf("Type: %v, Value: %v\n", counter.Type, counter.Value)
	//}
}

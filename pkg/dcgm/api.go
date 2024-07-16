package dcgm

import (
	"fmt"
	"log"
	"strconv"
)

// 初始化rocm_smi
func Init() error {
	return rsmiInit()
}

// 关闭rocm_smi
func ShutDown() error {
	return rsmiShutdown()
}

// 获取GPU数量
func NumMonitorDevices() (int, error) {
	return rsmiNumMonitorDevices()
}

// 获取设备利用率计数器
func UtilizationCount(dvInd int, utilizationCounters []RSMIUtilizationCounter, count int) (timestamp int64, err error) {
	return rsmiUtilizationCountGet(dvInd, utilizationCounters, count)
}

// 获取设备名称
func DevName(dvInd int) (name string, err error) {
	return rsmiDevNameGet(dvInd)
}

// 获取可用的pcie带宽列表
func DevPciBandwidth(dvInd int) RSMIPcieBandwidth {
	return rsmiDevPciBandwidthGet(dvInd)

}

// 内存总量
func MemoryTotal(dvInd int) int64 {
	return rsmiDevMemoryTotalGet(dvInd, RSMI_MEM_TYPE_FIRST)

}

// 内存使用量
func MemoryUsed(dvInd int) int64 {
	return rsmiDevMemoryUsageGet(dvInd, RSMI_MEM_TYPE_FIRST)

}

// 内存使用百分比
func MemoryPercent(dvInd int) int {
	return rsmiDevMemoryBusyPercentGet(dvInd)
}

// 获取设备温度值
//func DevTemp(dvInd int) int64 {
//	return go_rsmi_dev_temp_metric_get(dvInd)
//}

// 获取设别性能级别
func DevPerfLevelGet(dvInd int) (perf RSMIDevPerfLevel, err error) {
	return rsmiDevPerfLevelGet(dvInd)
}

// 设置设备PowerPlay性能级别
func DevPerfLevelSet(dvInd int, level RSMIDevPerfLevel) error {
	return rsmiDevPerfLevelSet(dvInd, level)
}

// 获取gpu度量信息
func DevGpuMetricsInfo(dvInd int) (gpuMetrics RSMIGPUMetrics, err error) {
	return rsmiDevGpuMetricsInfoGet(dvInd)
}

// 获取设备监控中的指标
func CollectDeviceMetrics() (devices []MonitorInfo, err error) {
	numMonitorDevices, err := rsmiNumMonitorDevices()
	if err != nil {
		return nil, err
	}
	for i := 0; i < numMonitorDevices; i++ {
		bdfid := rsmiDevPciIdGet(i)
		// 解析BDFID
		domain := (bdfid >> 32) & 0xffffffff
		bus := (bdfid >> 8) & 0xff
		dev := (bdfid >> 3) & 0x1f
		function := bdfid & 0x7
		// 格式化PCI ID
		pciBusNumber := fmt.Sprintf("%04X:%02X:%02X.%X", domain, bus, dev, function)
		//设备序列号
		deviceId := rsmiDevSerialNumberGet(i)
		//设备id
		devId := rsmiDevIdGet(i)
		//型号名称
		subSystemName := type2name[fmt.Sprintf("%X", devId)]
		//设备温度
		temperature := rsmiDevTempMetricGet(i, 0, RSMI_TEMP_CURRENT)
		t, err := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(temperature)/1000.0), 64)
		if err != nil {
			return nil, err
		}
		//设备平均功耗
		powerUsage := rsmiDevPowerAveGet(i, 0)
		pu, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(powerUsage)/1000000.0), 64)
		fmt.Printf("🔋 DCU[%v] power cap : %v \n", i, pu)
		//获取设备功率上限
		powerCap := rsmiDevPowerCapGet(i, 0)
		pc, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(powerCap)/1000000.0), 64)
		fmt.Printf("\U0001FAAB DCU[%v] power usage : %v \n", i, pc)
		//获取设备内存总量
		memoryCap := rsmiDevMemoryTotalGet(i, RSMI_MEM_TYPE_FIRST)
		mc, _ := strconv.ParseFloat(fmt.Sprintf("%f", float64(memoryCap)/1.0), 64)
		fmt.Printf(" DCU[%v] memory cap : %v \n", i, mc)
		//获取设备内存使用量
		memoryUsed := rsmiDevMemoryUsageGet(i, RSMI_MEM_TYPE_FIRST)
		mu, _ := strconv.ParseFloat(fmt.Sprintf("%f", float64(memoryUsed)/1.0), 64)
		fmt.Printf(" DCU[%v] memory used : %v \n", i, mu)
		//获取设备设备忙碌时间百分比
		utilizationRate, _ := rsmiDevBusyPercentGet(i)
		ur, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(utilizationRate)/1.0), 64)
		fmt.Printf(" DCU[%v] utilization rate : %v \n", i, ur)
		//获取pcie流量信息
		sent, received, maxPktSz := rsmiDevPciThroughputGet(i)
		pcieBwMb, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(received+sent)*float64(maxPktSz)/1024.0/1024.0), 64)
		fmt.Printf(" DCU[%v] PCIE  bandwidth : %v \n", i, pcieBwMb)
		//获取设备系统时钟速度列表
		clk, _ := rsmiDevGpuClkFreqGet(i, RSMI_CLK_TYPE_SYS)
		sclk, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(clk.Frequency[clk.Current])/1000000.0), 64)
		fmt.Printf(" DCU[%v] SCLK : %v \n", i, sclk)
		deviceInfo := MonitorInfo{
			PicBusNumber:    pciBusNumber,
			DeviceId:        deviceId,
			SubSystemName:   subSystemName,
			Temperature:     t,
			PowerUsage:      pu,
			powerCap:        pc,
			MemoryCap:       mc,
			MemoryUsed:      mu,
			UtilizationRate: ur,
			PcieBwMb:        pcieBwMb,
			Clk:             sclk,
		}
		devices = append(devices, deviceInfo)
	}
	log.Println("devices: ", dataToJson(devices))
	return
}

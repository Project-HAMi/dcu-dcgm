package dcgm

// 初始化rocm_smi
func Init() error {
	return go_rsmi_init()
}

// 关闭rocm_smi
func ShutDown() error {
	return go_rsmi_shutdown()
}

// 获取GPU数量
func NumMonitorDevices() (int, error) {
	return go_rsmi_num_monitor_devices()
}

// 获取设备利用率计数器
func UtilizationCount(dvInd int, utilizationCounters []RSMIUtilizationCounter, count int) (timestamp int64, err error) {
	return go_rsmi_utilization_count_get(dvInd, utilizationCounters, count)
}

// 获取设备名称
func DevName(dvInd int) (name string, err error) {
	return go_rsmi_dev_name_get(dvInd)
}

// 获取可用的pcie带宽列表
func DevPciBandwidth(dvInd int) RSMIPcieBandwidth {
	return go_rsmi_dev_pci_bandwidth_get(dvInd)

}

// 内存总量
func MemoryTotal(dvInd int) int64 {
	return go_rsmi_dev_memory_total_get(dvInd, RSMI_MEM_TYPE_FIRST)

}

// 内存使用量
func MemoryUsed(dvInd int) int64 {
	return go_rsmi_dev_memory_usage_get(dvInd, RSMI_MEM_TYPE_FIRST)

}

// 内存使用百分比
func MemoryPercent(dvInd int) int {
	return go_rsmi_dev_memory_busy_percent_get(dvInd)
}

// 获取设备温度值
//func DevTemp(dvInd int) int64 {
//	return go_rsmi_dev_temp_metric_get(dvInd)
//}

// 获取设别性能级别
func DevPerfLevelGet(dvInd int) (perf RSMIDevPerfLevel, err error) {
	return go_rsmi_dev_perf_level_get(dvInd)
}

// 设置设备PowerPlay性能级别
func DevPerfLevelSet(dvInd int, level RSMIDevPerfLevel) error {
	return go_rsmi_dev_perf_level_set(dvInd, level)
}

// 获取gpu度量信息
func DevGpuMetricsInfo(dvInd int) (gpuMetrics RSMIGPUMetrics, err error) {
	return go_rsmi_dev_gpu_metrics_info_get(dvInd)

}

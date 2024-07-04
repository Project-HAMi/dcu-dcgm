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

func DevPciBandwidth(dvInd int) RSMIPcieBandwidth {
	return go_rsmi_dev_pci_bandwidth_get(dvInd)

}

func MemoryTotal(dvInd int) int64 {
	return go_rsmi_dev_memory_total_get(dvInd, RSMI_MEM_TYPE_FIRST)

}

func MemoryUsed(dvInd int) int64 {
	return go_rsmi_dev_memory_usage_get(dvInd, RSMI_MEM_TYPE_FIRST)

}

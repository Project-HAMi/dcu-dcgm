package dcgm

/*
#cgo CFLAGS: -Wall -I./include
#cgo LDFLAGS: -L./lib -lrocm_smi64 -Wl,--unresolved-symbols=ignore-in-object-files
#include <stdint.h>
#include <kfd_ioctl.h>
#include <rocm_smi64Config.h>
#include <rocm_smi.h>
*/
import "C"
import (
	"fmt"
	"unsafe"

	"github.com/golang/glog"
)

// rsmiDevTempMetricGet 获取设备的温度度量值 *
func rsmiDevTempMetricGet(dvInd int, sensorType int, metric RSMITemperatureMetric) (temp int64, err error) {
	var temperature C.int64_t
	ret := C.rsmi_dev_temp_metric_get(C.uint32_t(dvInd), C.uint32_t(sensorType), C.rsmi_temperature_metric_t(metric), &temperature)
	if err = errorString(ret); err != nil {
		return 0, fmt.Errorf("rsmiDevTempMetricGet:%s", err)
	}
	temp = int64(temperature)
	return
}

// rsmiDevVoltMetricGet 获取设备的电压度量值
func rsmiDevVoltMetricGet(dvInd int, voltageType RSMIVoltageType, metric RSMIVoltageMetric) int64 {
	var voltage C.int64_t
	C.rsmi_dev_volt_metric_get(C.uint32_t(dvInd), C.rsmi_voltage_type_t(voltageType), C.rsmi_voltage_metric_t(metric), &voltage)
	return int64(voltage)
}

// rsmiDevFanSpeedSet 设置设备风扇转速，以rpm为单位
func rsmiDevFanSpeedSet(dvInd, sensorInd int, speed int64) (err error) {
	ret := C.rsmi_dev_fan_speed_set(C.uint32_t(dvInd), C.uint32_t(sensorInd), C.uint64_t(speed))
	glog.Infof("rsmi_dev_fan_speed_set_ret:%v, retstr:%v", ret, errorString(ret))
	if err = errorString(ret); err != nil {
		return fmt.Errorf("Error rsmi_dev_fan_speed_set: %s", err)
	}
	return nil
}

// rsmiDevBusyPercentGet 获取设备设备忙碌时间百分比
func rsmiDevBusyPercentGet(dvInd int) (busyPercent int, err error) {
	var cbusyPercent C.uint32_t
	ret := C.rsmi_dev_busy_percent_get(C.uint32_t(dvInd), &cbusyPercent)
	//glog.Info("rsmi_dev_busy_percent_get:", ret)
	if err = errorString(ret); err != nil {
		return 0, fmt.Errorf("Error rsmi_dev_busy_percent_get:%s", err)
	}
	busyPercent = int(cbusyPercent)
	return busyPercent, nil
}

// rsmiUtilizationCountGet 获取设备利用率计数器
func rsmiUtilizationCountGet(dvInd int, utilizationCounters []RSMIUtilizationCounter, count int) (timestamp int64, err error) {
	// 转换 Go 结构体数组到 C 结构体数组
	cUtilizationCounters := make([]C.rsmi_utilization_counter_t, len(utilizationCounters))
	for i, uc := range utilizationCounters {
		cUtilizationCounters[i] = C.rsmi_utilization_counter_t{
			_type: C.RSMI_UTILIZATION_COUNTER_TYPE(uc.Type),
			value: C.uint64_t(uc.Value),
		}
	}

	var ctimestamp C.uint64_t
	// 调用 C 函数
	ret := C.rsmi_utilization_count_get(
		C.uint32_t(dvInd),
		&cUtilizationCounters[0],
		C.uint32_t(count),
		&ctimestamp,
	)
	//glog.Infof("rsmi_utilization_count_get ret:%v ,retstr:%v ", ret, errorString(ret))
	if err = errorString(ret); err != nil {
		return 0, fmt.Errorf("Error rsmi_utilization_count_get:%s", err)
	}
	// 更新 Go 结构体数组中的值
	for i := range utilizationCounters {
		utilizationCounters[i].Value = uint64(cUtilizationCounters[i].value)
	}
	//glog.Infof("utilizationCounters:%v,timestamp:%v", utilizationCounters, int64(ctimestamp))

	return int64(ctimestamp), nil
}

// rsmiDevPerfLevelGet 获取设备的性能级别
func rsmiDevPerfLevelGet(dvInd int) (perf RSMIDevPerfLevel, err error) {
	var cPerfLevel C.rsmi_dev_perf_level_t
	ret := C.rsmi_dev_perf_level_get(C.uint32_t(dvInd), &cPerfLevel)
	if err = errorString(ret); err != nil {
		return RSMIDevPerfLevel(cPerfLevel), fmt.Errorf("Error rsmi_dev_perf_level_get:%s", err)
	}
	perf = RSMIDevPerfLevel(cPerfLevel)
	glog.Info("dev_perf_level:", perf)
	return perf, nil
}

// rsmiPerfDeterminismModeSet 设置设备的性能确定性模式
func rsmiPerfDeterminismModeSet(dvInd int, clkValue int64) (err error) {
	ret := C.rsmi_perf_determinism_mode_set(C.uint32_t(dvInd), C.uint64_t(clkValue))
	glog.Infof("dev_perf_determinism_mode ret:%v, retstr:%v", ret, errorString(ret))
	if err = errorString(ret); err != nil {
		return fmt.Errorf("Error rsmi_perf_determinism_mode_set:%s", err)
	}
	return
}

// rsmiDevOverdriveLevelGet 获取设备的超速百分比
func rsmiDevOverdriveLevelGet(dvInd int) (od int, err error) {
	var cod C.uint32_t
	ret := C.rsmi_dev_overdrive_level_get(C.uint32_t(dvInd), &cod)
	if err = errorString(ret); err != nil {
		return int(cod), fmt.Errorf("Error rsmi_dev_overdrive_level_get:%s", err)
	}
	od = int(cod)
	return
}

// rsmiDevGpuClkFreqGet 获取设备系统时钟速度列表
func rsmiDevGpuClkFreqGet(dvInd int, clkType RSMIClkType) (frequencies RSMIFrequencies, err error) {
	var cfrequencies C.rsmi_frequencies_t
	ret := C.rsmi_dev_gpu_clk_freq_get(C.uint32_t(dvInd), C.rsmi_clk_type_t(clkType), &cfrequencies)
	if err = errorString(ret); err != nil {
		return frequencies, fmt.Errorf("Error rsmi_dev_gpu_clk_freq_get:%s", err)
	}
	frequencies = RSMIFrequencies{
		NumSupported: uint32(cfrequencies.num_supported),
		Current:      uint32(cfrequencies.current),
		Frequency:    *(*[32]uint64)(unsafe.Pointer(&cfrequencies.frequency)),
	}
	return
}

// rsmiDevOdVoltInfoGet 获取设备电压/频率曲线信息
func rsmiDevOdVoltInfoGet(dvInd int) (odv RSMIOdVoltFreqData, err error) {
	var codv C.rsmi_od_volt_freq_data_t
	ret := C.rsmi_dev_od_volt_info_get(C.uint32_t(dvInd), &codv)
	glog.Infof("rsmi_dev_od_volt_info_get ret:%v, retstr:%v", ret, errorString(ret))
	if err = errorString(ret); err != nil {
		return odv, fmt.Errorf("Error rsmi_dev_od_volt_info_get:%s", err)
	}
	odv = RSMIOdVoltFreqData{
		CurrMclkRange: RSMIRange{
			LowerBound: uint64(codv.curr_sclk_range.lower_bound),
			UpperBound: uint64(codv.curr_sclk_range.upper_bound),
		},
		CurrSclkRange: RSMIRange{
			LowerBound: uint64(codv.curr_mclk_range.lower_bound),
			UpperBound: uint64(codv.curr_mclk_range.upper_bound),
		},
		SclkFreqLimits: RSMIRange{
			LowerBound: uint64(codv.sclk_freq_limits.lower_bound),
			UpperBound: uint64(codv.sclk_freq_limits.upper_bound),
		},
		MclkFreqLimits: RSMIRange{
			LowerBound: uint64(codv.mclk_freq_limits.lower_bound),
			UpperBound: uint64(codv.mclk_freq_limits.upper_bound),
		},
		Curve:      RSMIOdVoltCurve{},
		NumRegions: uint32(codv.num_regions),
	}
	for i := 0; i < len(codv.curve.vc_points); i++ {
		odv.Curve.VcPoints[i] = RSMIOdVddcPoint{
			Frequency: uint64(codv.curve.vc_points[i].frequency),
			Voltage:   uint64(codv.curve.vc_points[i].voltage),
		}
	}
	return
}

// rsmiDevGpuMetricsInfoGet 获取gpu度量信息
func rsmiDevGpuMetricsInfoGet(dvInd int) (gpuMetrics RSMIGPUMetrics, err error) {
	var cgpuMetrics C.rsmi_gpu_metrics_t
	ret := C.rsmi_dev_gpu_metrics_info_get(C.uint32_t(dvInd), &cgpuMetrics)
	if err = errorString(ret); err != nil {
		return gpuMetrics, fmt.Errorf("Error rsmi_dev_gpu_metrics_info_get:%s", err)
	}
	gpuMetrics = RSMIGPUMetrics{
		CommonHeader: MetricsTableHeader{
			StructureSize:   uint16(cgpuMetrics.common_header.structure_size),
			FormatRevision:  uint8(cgpuMetrics.common_header.format_revision),
			ContentRevision: uint8(cgpuMetrics.common_header.content_revision),
		},
		TemperatureEdge:        uint16(cgpuMetrics.temperature_edge),
		TemperatureHotspot:     uint16(cgpuMetrics.temperature_hotspot),
		TemperatureMem:         uint16(cgpuMetrics.temperature_mem),
		TemperatureVRGfx:       uint16(cgpuMetrics.temperature_vrgfx),
		TemperatureVRSoc:       uint16(cgpuMetrics.temperature_vrsoc),
		TemperatureVRMem:       uint16(cgpuMetrics.temperature_vrmem),
		AverageGfxActivity:     uint16(cgpuMetrics.average_gfx_activity),
		AverageUmcActivity:     uint16(cgpuMetrics.average_umc_activity),
		AverageMmActivity:      uint16(cgpuMetrics.average_mm_activity),
		AverageSocketPower:     uint16(cgpuMetrics.average_socket_power),
		EnergyAccumulator:      uint64(cgpuMetrics.energy_accumulator),
		SystemClockCounter:     uint64(cgpuMetrics.system_clock_counter),
		AverageGfxclkFrequency: uint16(cgpuMetrics.average_gfxclk_frequency),
		AverageSocclkFrequency: uint16(cgpuMetrics.average_socclk_frequency),
		AverageUclkFrequency:   uint16(cgpuMetrics.average_uclk_frequency),
		AverageVclk0Frequency:  uint16(cgpuMetrics.average_vclk0_frequency),
		AverageDclk0Frequency:  uint16(cgpuMetrics.average_dclk0_frequency),
		AverageVclk1Frequency:  uint16(cgpuMetrics.average_vclk1_frequency),
		AverageDclk1Frequency:  uint16(cgpuMetrics.average_dclk1_frequency),
		CurrentGfxclk:          uint16(cgpuMetrics.current_gfxclk),
		CurrentSocclk:          uint16(cgpuMetrics.current_socclk),
		CurrentUclk:            uint16(cgpuMetrics.current_uclk),
		CurrentVclk0:           uint16(cgpuMetrics.current_vclk0),
		CurrentDclk0:           uint16(cgpuMetrics.current_dclk0),
		CurrentVclk1:           uint16(cgpuMetrics.current_vclk1),
		CurrentDclk1:           uint16(cgpuMetrics.current_dclk1),
		ThrottleStatus:         uint32(cgpuMetrics.throttle_status),
		CurrentFanSpeed:        uint16(cgpuMetrics.current_fan_speed),
		PcieLinkWidth:          uint16(cgpuMetrics.pcie_link_width),
		PcieLinkSpeed:          uint16(cgpuMetrics.pcie_link_speed),
		Padding:                uint16(cgpuMetrics.padding),
		GfxActivityAcc:         uint32(cgpuMetrics.gfx_activity_acc),
		MemActivityAcc:         uint32(cgpuMetrics.mem_actvity_acc),
		TempetureHBM:           *((*[4]uint16)(unsafe.Pointer(&cgpuMetrics.temperature_hbm))),
	}
	glog.Info("rsmi_dev_gpu_metrics_info_get:%s", dataToJson(gpuMetrics))
	return
}

// rsmiDevEccStatusGet 获取GPU块的ECC状态
func rsmiDevEccStatusGet(dvInd int, block RSMIGpuBlock) (state RSMIRasErrState, err error) {
	glog.Infof("rsmiDevEccStatusGet: %d,%d", dvInd, block)
	var sstate C.rsmi_ras_err_state_t
	ret := C.rsmi_dev_ecc_status_get(C.uint32_t(dvInd), C.rsmi_gpu_block_t(block), &sstate)
	glog.Infof("rsmi_dev_ecc_status_get ret:%v", ret)
	if err = errorString(ret); err != nil {
		return state, fmt.Errorf("Error rsmi_dev_ecc_status_get:%s", err)
	}
	state = RSMIRasErrState(sstate)
	glog.Infof("rsmiDevEccStatusGet:%v", sstate)
	return
}

// rsmiDevEccCountGet 获取GPU块的错误计数
func rsmiDevEccCountGet(dvInd int, gpuBlock RSMIGpuBlock) (errorCount RSMIErrorCount, err error) {
	var cerrorCount C.rsmi_error_count_t
	ret := C.rsmi_dev_ecc_count_get(C.uint32_t(dvInd), C.rsmi_gpu_block_t(gpuBlock), &cerrorCount)
	glog.Infof("rsmiDevEccCountGet:%v,ret retstr:%v", ret, errorString(ret))
	if err = errorString(ret); err != nil {
		return errorCount, fmt.Errorf("Error rsmi_dev_ecc_count_get:%s", err)
	}
	errorCount = RSMIErrorCount{
		CorrectableErr:   uint64(cerrorCount.correctable_err),
		UncorrectableErr: uint64(cerrorCount.uncorrectable_err),
	}
	glog.Infof("DCUBlockType:%v, DevEccCount:%v", gpuBlock, dataToJson(errorCount))
	return
}

// rsmiDevEccEnabledGet 获取已启用的ECC位掩码
func rsmiDevEccEnabledGet(dvInd int) (enabledBlocks int64, err error) {
	var cenabledBlocks C.uint64_t
	ret := C.rsmi_dev_ecc_enabled_get(C.uint32_t(dvInd), &cenabledBlocks)
	if err = errorString(ret); err != nil {
		return enabledBlocks, fmt.Errorf("Error rsmi_dev_ecc_enabled_get:%s", err)
	}
	enabledBlocks = int64(cenabledBlocks)
	glog.Infof("DCUBlockType:%v", enabledBlocks)
	return
}

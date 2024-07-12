package dcgm

/*
#cgo CFLAGS: -Wall -I/opt/dtk-24.04/rocm_smi/include/rocm_smi
#cgo LDFLAGS: -L/opt/dtk-24.04/rocm_smi/lib -lrocm_smi64 -Wl,--unresolved-symbols=ignore-in-object-files
#include <stdint.h>
#include <kfd_ioctl.h>
#include <rocm_smi64Config.h>
#include <rocm_smi.h>
*/
import "C"
import (
	"fmt"
	"log"
	"unsafe"
)

type RSMITemperatureMetric C.rsmi_temperature_metric_t

const (
	RSMI_TEMP_CURRENT        RSMITemperatureMetric = C.RSMI_TEMP_CURRENT
	RSMI_TEMP_FIRST          RSMITemperatureMetric = C.RSMI_TEMP_FIRST
	RSMI_TEMP_MAX            RSMITemperatureMetric = C.RSMI_TEMP_MAX
	RSMI_TEMP_MIN            RSMITemperatureMetric = C.RSMI_TEMP_MIN
	RSMI_TEMP_MAX_HYST       RSMITemperatureMetric = C.RSMI_TEMP_MAX_HYST
	RSMI_TEMP_MIN_HYST       RSMITemperatureMetric = C.RSMI_TEMP_MIN_HYST
	RSMI_TEMP_CRITICAL       RSMITemperatureMetric = C.RSMI_TEMP_CRITICAL
	RSMI_TEMP_CRITICAL_HYST  RSMITemperatureMetric = C.RSMI_TEMP_CRITICAL_HYST
	RSMI_TEMP_EMERGENCY      RSMITemperatureMetric = C.RSMI_TEMP_EMERGENCY
	RSMI_TEMP_EMERGENCY_HYST RSMITemperatureMetric = C.RSMI_TEMP_EMERGENCY_HYST
	RSMI_TEMP_CRIT_MIN       RSMITemperatureMetric = C.RSMI_TEMP_CRIT_MIN
	RSMI_TEMP_CRIT_MIN_HYST  RSMITemperatureMetric = C.RSMI_TEMP_CRIT_MIN_HYST
	RSMI_TEMP_OFFSET         RSMITemperatureMetric = C.RSMI_TEMP_OFFSET
	RSMI_TEMP_LOWEST         RSMITemperatureMetric = C.RSMI_TEMP_LOWEST
	RSMI_TEMP_HIGHEST        RSMITemperatureMetric = C.RSMI_TEMP_HIGHEST
	RSMI_TEMP_LAST           RSMITemperatureMetric = C.RSMI_TEMP_LAST
)

type RSMIVoltageType C.rsmi_voltage_type_t

const (
	RSMI_VOLT_TYPE_FIRST   RSMIVoltageType = C.RSMI_VOLT_TYPE_FIRST
	RSMI_VOLT_TYPE_VDDGFX  RSMIVoltageType = C.RSMI_VOLT_TYPE_VDDGFX
	RSMI_VOLT_TYPE_LAST    RSMIVoltageType = C.RSMI_VOLT_TYPE_LAST
	RSMI_VOLT_TYPE_INVALID RSMIVoltageType = C.RSMI_VOLT_TYPE_INVALID
)

type RSMIVoltageMetric C.rsmi_voltage_metric_t

const (
	RSMI_VOLT_CURRENT  RSMIVoltageMetric = C.RSMI_VOLT_CURRENT //!< Voltage current value.
	RSMI_VOLT_FIRST    RSMIVoltageMetric = C.RSMI_VOLT_FIRST
	RSMI_VOLT_MAX      RSMIVoltageMetric = C.RSMI_VOLT_MAX      //!< Voltage max value.
	RSMI_VOLT_MIN_CRIT RSMIVoltageMetric = C.RSMI_VOLT_MIN_CRIT //!< Voltage critical min value.
	RSMI_VOLT_MIN      RSMIVoltageMetric = C.RSMI_VOLT_MIN      //!< Voltage min value.
	RSMI_VOLT_MAX_CRIT RSMIVoltageMetric = C.RSMI_VOLT_MAX_CRIT //!< Voltage critical max value.
	RSMI_VOLT_AVERAGE  RSMIVoltageMetric = C.RSMI_VOLT_AVERAGE  //!< Average voltage.
	RSMI_VOLT_LOWEST   RSMIVoltageMetric = C.RSMI_VOLT_LOWEST   //!< Historical minimum voltage.
	RSMI_VOLT_HIGHEST  RSMIVoltageMetric = C.RSMI_VOLT_HIGHEST  //!< Historical maximum voltage.
	RSMI_VOLT_LAST                       = C.RSMI_VOLT_LAST
)

type RSMIUtilizationCounterType C.RSMI_UTILIZATION_COUNTER_TYPE

const (
	RSMI_UTILIZATION_COUNTER_FIRST RSMIUtilizationCounterType = C.RSMI_UTILIZATION_COUNTER_FIRST
	RSMI_COARSE_GRAIN_GFX_ACTIVITY RSMIUtilizationCounterType = C.RSMI_COARSE_GRAIN_GFX_ACTIVITY
	RSMI_COARSE_GRAIN_MEM_ACTIVITY RSMIUtilizationCounterType = C.RSMI_COARSE_GRAIN_MEM_ACTIVITY
	RSMI_UTILIZATION_COUNTER_LAST  RSMIUtilizationCounterType = C.RSMI_UTILIZATION_COUNTER_LAST
)

type RSMIUtilizationCounter struct {
	Type  RSMIUtilizationCounterType
	Value uint64
}

type RSMIClkType C.rsmi_clk_type_t

const (
	RSMI_CLK_TYPE_SYS  RSMIClkType = C.RSMI_CLK_TYPE_SYS
	RSMI_CLK_TYPE_DF   RSMIClkType = C.RSMI_CLK_TYPE_DF
	RSMI_CLK_TYPE_DCEF RSMIClkType = C.RSMI_CLK_TYPE_DCEF
	RSMI_CLK_TYPE_SOC  RSMIClkType = C.RSMI_CLK_TYPE_SOC
	RSMI_CLK_TYPE_MEM  RSMIClkType = C.RSMI_CLK_TYPE_MEM
	RSMI_CLK_TYPE_PCIE RSMIClkType = C.RSMI_CLK_TYPE_PCIE
	RSMI_CLK_INVALID   RSMIClkType = C.RSMI_CLK_INVALID
)

type RSMIOdVoltFreqData struct {
	CurrSclkRange  RSMIRange
	CurrMclkRange  RSMIRange
	SclkFreqLimits RSMIRange
	MclkFreqLimits RSMIRange
	Curve          RSMIOdVoltCurve
	NumRegions     uint32
}

type RSMIRange struct {
	LowerBound uint64
	UpperBound uint64
}

type RSMIOdVoltCurve struct {
	VcPoints [3]RSMIOdVddcPoint
}

type RSMIOdVddcPoint struct {
	Frequency uint64
	Voltage   uint64
}

type MetricsTableHeader struct {
	StructureSize   uint16
	FormatRevision  uint8
	ContentRevision uint8
}

type RSMIGPUMetrics struct {
	CommonHeader           MetricsTableHeader
	TemperatureEdge        uint16
	TemperatureHotspot     uint16
	TemperatureMem         uint16
	TemperatureVRGfx       uint16
	TemperatureVRSoc       uint16
	TemperatureVRMem       uint16
	AverageGfxActivity     uint16
	AverageUmcActivity     uint16
	AverageMmActivity      uint16
	AverageSocketPower     uint16
	EnergyAccumulator      uint64
	SystemClockCounter     uint64
	AverageGfxclkFrequency uint16
	AverageSocclkFrequency uint16
	AverageUclkFrequency   uint16
	AverageVclk0Frequency  uint16
	AverageDclk0Frequency  uint16
	AverageVclk1Frequency  uint16
	AverageDclk1Frequency  uint16
	CurrentGfxclk          uint16
	CurrentSocclk          uint16
	CurrentUclk            uint16
	CurrentVclk0           uint16
	CurrentDclk0           uint16
	CurrentVclk1           uint16
	CurrentDclk1           uint16
	ThrottleStatus         uint32
	CurrentFanSpeed        uint16
	PcieLinkWidth          uint16
	PcieLinkSpeed          uint16
	Padding                uint16
	GfxActivityAcc         uint32
	MemActivityAcc         uint32
	TemperatureHBM         [4]uint16
}

// go_rsmi_dev_temp_metric_get 获取设备的温度度量值 *
func go_rsmi_dev_temp_metric_get(dvInd int, sensorType int, metric RSMITemperatureMetric) int64 {
	var temperature C.int64_t
	C.rsmi_dev_temp_metric_get(C.uint32_t(dvInd), C.uint32_t(sensorType), C.rsmi_temperature_metric_t(metric), &temperature)
	return int64(temperature)
}

// go_rsmi_dev_volt_metric_get 获取设备的电压度量值
func go_rsmi_dev_volt_metric_get(dvInd int, voltageType RSMIVoltageType, metric RSMIVoltageMetric) int64 {
	var voltage C.int64_t
	C.rsmi_dev_volt_metric_get(C.uint32_t(dvInd), C.rsmi_voltage_type_t(voltageType), C.rsmi_voltage_metric_t(metric), &voltage)
	return int64(voltage)
}

// go_rsmi_dev_fan_reset 将风扇复位为自动驱动控制
func go_go_rsmi_dev_fan_reset(dvInd, sensorInd int) (err error) {
	ret := C.rsmi_dev_fan_reset(C.uint32_t(dvInd), C.uint32_t(sensorInd))
	log.Println("go_rsmi_dev_fan_reset_ret:", ret)
	if err = errorString(ret); err != nil {
		return fmt.Errorf("Error go_rsmi_dev_fan_reset: %s", err)
	}
	return nil
}

// go_rsmi_dev_fan_speed_set 设置设备风扇转速，以rpm为单位
func go_rsmi_dev_fan_speed_set(dvInd, sensorInd int, speed int64) (err error) {
	ret := C.rsmi_dev_fan_speed_set(C.uint32_t(dvInd), C.uint32_t(sensorInd), C.uint64_t(speed))
	log.Println("go_rsmi_dev_fan_speed_set_ret:", ret)
	if err = errorString(ret); err != nil {
		return fmt.Errorf("Error go_rsmi_dev_fan_speed_set: %s", err)
	}
	return nil
}

// go_rsmi_dev_busy_percent_get 获取设备设备忙碌时间百分比
func go_rsmi_dev_busy_percent_get(dvInd int) (busyPercent int, err error) {
	var cbusyPercent C.uint32_t
	ret := C.rsmi_dev_busy_percent_get(C.uint32_t(dvInd), &cbusyPercent)
	log.Println("rsmi_dev_busy_percent_get:", ret)
	if err = errorString(ret); err != nil {
		return 0, fmt.Errorf("Error go_rsmi_dev_busy_percent_get:%s", err)
	}
	busyPercent = int(cbusyPercent)
	return busyPercent, nil
}

// go_rsmi_utilization_count_get 获取设备利用率计数器
func go_rsmi_utilization_count_get(dvInd int, utilizationCounters []RSMIUtilizationCounter, count int) (timestamp int64, err error) {
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

	if err = errorString(ret); err != nil {
		return 0, fmt.Errorf("Error go_rsmi_dev_busy_percent_get:%s", err)
	}
	// 更新 Go 结构体数组中的值
	for i := range utilizationCounters {
		utilizationCounters[i].Value = uint64(cUtilizationCounters[i].value)
	}

	return int64(ctimestamp), nil
}

// go_rsmi_dev_perf_level_get 获取设备的性能级别
func go_rsmi_dev_perf_level_get(dvInd int) (perf RSMIDevPerfLevel, err error) {
	var cPerfLevel C.rsmi_dev_perf_level_t
	ret := C.rsmi_dev_perf_level_get(C.uint32_t(dvInd), &cPerfLevel)
	if err = errorString(ret); err != nil {
		return RSMIDevPerfLevel(cPerfLevel), fmt.Errorf("Error go_rsmi_dev_perf_level_get:%s", err)
	}
	perf = RSMIDevPerfLevel(cPerfLevel)
	log.Println("dev_perf_level:", perf)
	return perf, nil
}

// go_rsmi_perf_determinism_mode_set 设置设备的性能确定性模式
func go_rsmi_perf_determinism_mode_set(dvInd int, clkValue int64) (err error) {

	ret := C.rsmi_perf_determinism_mode_set(C.uint32_t(dvInd), C.uint64_t(clkValue))
	if err = errorString(ret); err != nil {
		return fmt.Errorf("Error go_rsmi_perf_determinism_mode_set:%s", err)
	}
	return
}

// go_rsmi_dev_overdrive_level_get 获取设备的超速百分比
func go_rsmi_dev_overdrive_level_get(dvInd int) (od int, err error) {
	var cod C.uint32_t
	ret := C.rsmi_dev_overdrive_level_get(C.uint32_t(dvInd), &cod)
	if err = errorString(ret); err != nil {
		return int(cod), fmt.Errorf("Error go_rsmi_dev_overdrive_level_get:%s", err)
	}
	od = int(cod)
	return
}

// go_rsmi_dev_gpu_clk_freq_get 获取设备系统时钟速度列表
func go_rsmi_dev_gpu_clk_freq_get(dvInd int, clkType RSMIClkType) (frequencies RSMIFrequencies, err error) {
	var cfrequencies C.rsmi_frequencies_t
	ret := C.rsmi_dev_gpu_clk_freq_get(C.uint32_t(dvInd), C.rsmi_clk_type_t(clkType), &cfrequencies)
	if err = errorString(ret); err != nil {
		return frequencies, fmt.Errorf("Error go_rsmi_dev_gpu_clk_freq_get:%s", err)
	}
	frequencies = RSMIFrequencies{
		NumSupported: uint32(cfrequencies.num_supported),
		Current:      uint32(cfrequencies.current),
		Frequency:    *(*[32]uint64)(unsafe.Pointer(&cfrequencies.frequency)),
	}
	return
}

// go_rsmi_dev_od_volt_info_get 获取设备电压/频率曲线信息
func go_rsmi_dev_od_volt_info_get(dvInd int) (odv RSMIOdVoltFreqData, err error) {
	var codv C.rsmi_od_volt_freq_data_t
	ret := C.rsmi_dev_od_volt_info_get(C.uint32_t(dvInd), &codv)
	if err = errorString(ret); err != nil {
		return odv, fmt.Errorf("Error go_rsmi_dev_od_volt_info_get:%s", err)
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

// go_rsmi_dev_gpu_metrics_info_get 获取gpu度量信息
func go_rsmi_dev_gpu_metrics_info_get(dvInd int) (gpuMetrics RSMIGPUMetrics, err error) {
	var cgpuMetrics C.rsmi_gpu_metrics_t
	ret := C.rsmi_dev_gpu_metrics_info_get(C.uint32_t(dvInd), &cgpuMetrics)
	if err = errorString(ret); err != nil {
		return gpuMetrics, fmt.Errorf("Error go_rsmi_dev_gpu_metrics_info_get:%s", err)
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
		TemperatureHBM:         *((*[4]uint16)(unsafe.Pointer(&cgpuMetrics.temperature_hbm))),
	}
	log.Printf("go_rsmi_dev_gpu_metrics_info_get:%s", dataToJson(gpuMetrics))
	return
}

// go_rsmi_dev_ecc_status_get 获取GPU块的ECC状态
func rsmi_dev_ecc_status_get(dvInd int, block RSMIGpuBlock) (state RSMIRasErrState, err error) {
	var sstate C.rsmi_ras_err_state_t
	ret := C.rsmi_dev_ecc_status_get(C.uint32_t(dvInd), C.rsmi_gpu_block_t(block), &sstate)
	if err = errorString(ret); err != nil {
		return state, fmt.Errorf("Error go_rsmi_dev_ecc_status_get:%s", err)
	}
	state = RSMIRasErrState(sstate)
	return
}

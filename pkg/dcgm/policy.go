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

	"github.com/golang/glog"
)

// rsmiDevPerfLevelSet 设置设备PowerPlay性能级别
func rsmiDevPerfLevelSet(dvInd int, devPerfLevel RSMIDevPerfLevel) (err error) {
	glog.Info("dev_perf_level_set:", devPerfLevel)
	ret := C.rsmi_dev_perf_level_set(C.int32_t(dvInd), C.rsmi_dev_perf_level_t(devPerfLevel))
	glog.Infof("dev_perf_level_set ret:%v,retstr:%v", ret, errorString(ret))
	if err = errorString(ret); err != nil {
		return fmt.Errorf("dev_perf_level_set:%s", err)
	}
	return
}

// rsmiDevClkRangeSet 设置设备时钟范围信息
func rsmiDevClkRangeSet(dvInd int, minClkValue, maxClkValue int64, clkType RSMIClkType) (err error) {
	ret := C.rsmi_dev_clk_range_set(C.uint32_t(dvInd), C.uint64_t(minClkValue), C.uint64_t(maxClkValue), C.rsmi_clk_type_t(clkType))
	glog.Infof("rsmi_dev_clk_range_set ret:%v, retstr:%v", ret, errorString(ret))
	if err = errorString(ret); err != nil {
		glog.Errorf("Error rsmi_dev_clk_range_set:%s", err)
		return fmt.Errorf("Error rsmi_dev_clk_range_set:%s", err)
	}
	return
}

// rsmiDevOdVoltInfoSet 设置设备电压曲线点
func rsmiDevOdVoltInfoSet(dvInd, vPoint, clkValue, voltValue int) (err error) {
	ret := C.rsmi_dev_od_volt_info_set(C.uint32_t(dvInd), C.uint32_t(vPoint), C.uint64_t(clkValue), C.uint64_t(voltValue))
	glog.Infof("rsmi_dev_od_volt_info_set ret:%v", ret)
	if err = errorString(ret); err != nil {
		return fmt.Errorf("Error rsmi_dev_od_volt_info_set:%s", err)
	}
	return
}

// rsmiDevOverdriveLevelSet 设置设备超速百分比
func rsmiDevOverdriveLevelSet(dvInd, od int) (err error) {
	ret := C.rsmi_dev_overdrive_level_set(C.int32_t(dvInd), C.uint32_t(od))
	glog.Infof("rsmi_dev_overdrive_level_set ret:%v, retStr:%v", ret, errorString(ret))
	if err = errorString(ret); err != nil {
		return fmt.Errorf("Error rsmi_dev_overdrive_level_set:%s", err)
	}
	return
}

// rsmiDevGpuClkFreqSet 设置可用于指定时钟的频率集
func rsmiDevGpuClkFreqSet(dvInd int, clkType RSMIClkType, freqBitmask int64) (err error) {
	ret := C.rsmi_dev_gpu_clk_freq_set(C.uint32_t(dvInd), C.rsmi_clk_type_t(clkType), C.uint64_t(freqBitmask))
	if err = errorString(ret); err != nil {
		return fmt.Errorf("Error rsmi_dev_gpu_clk_freq_set:%s", err)
	}
	return nil
}

// rsmiDevCounterGroupSupported 判断设备是否支持特定事件组
func rsmiDevCounterGroupSupported(dvInd int, group RSMIEventGroup) (err error) {
	ret := C.rsmi_dev_counter_group_supported(C.uint32_t(dvInd), C.rsmi_event_group_t(group))
	if err = errorString(ret); err != nil {
		return fmt.Errorf("Error rsmi_dev_counter_group_supported:%s", err)
	}
	return
}

// rsmiDevCounterCreate 创建性能计数器对象
func rsmiDevCounterCreate(dvInd int, eventType RSMIEventType) (eventHandle EventHandle, err error) {
	var ceventHandle C.rsmi_event_handle_t
	ret := C.rsmi_dev_counter_create(C.uint32_t(dvInd), C.rsmi_event_type_t(eventType), &ceventHandle)
	if err = errorString(ret); err != nil {
		return eventHandle, fmt.Errorf("Error rsmi_dev_counter_create:%s", err)
	}
	eventHandle = EventHandle(ceventHandle)
	return
}

// rsmiDevCounterDestroy 释放性能计数器对象
func rsmiDevCounterDestroy(handle EventHandle) (err error) {
	var chandle C.rsmi_event_handle_t
	ret := C.rsmi_dev_counter_destroy(C.rsmi_event_handle_t(chandle))
	if err = errorString(ret); err != nil {
		return fmt.Errorf("Error rsmi_dev_counter_destroy:%s", err)
	}
	return
}

// rsmiCounterControl 发布性能计数器控制命令
func rsmiCounterControl(evtHandle EventHandle, cmd RSMICounterCommand) (err error) {
	ret := C.rsmi_counter_control(C.rsmi_event_handle_t(evtHandle), C.rsmi_counter_command_t(cmd), nil)

	if err := errorString(ret); err != nil {
		return fmt.Errorf("Error in rsmi_counter_control: %s", err)
	}
	return
}

// rsmiCounterRead 读取性能计数器的当前值
func rsmiCounterRead(handle EventHandle) (counterValue RSMICounterValue, err error) {
	var ccounterValue C.rsmi_counter_value_t
	ret := C.rsmi_counter_read(C.rsmi_event_handle_t(handle), &ccounterValue)
	if err = errorString(ret); err != nil {
		return counterValue, fmt.Errorf("Error rsmiCounterRead:%s", err)
	}
	counterValue = RSMICounterValue{
		Value:       uint64(ccounterValue.value),
		TimeEnabled: uint64(ccounterValue.time_enabled),
		TimeRunning: uint64(ccounterValue.time_running),
	}
	return
}

func rsmiCounterAvailableCountersGet(dvInd int, group RSMIEventGroup) (availAble int, err error) {
	var cavailAble C.uint32_t
	ret := C.rsmi_counter_available_counters_get(C.uint32_t(dvInd), C.rsmi_event_group_t(group), &cavailAble)
	if err = errorString(ret); err != nil {
		return availAble, fmt.Errorf("Error rsmiCounterAvailableCountersGet:%s", err)
	}
	availAble = int(cavailAble)
	return
}

// rsmiDevFanReset 将风扇复位为自动驱动控制
func rsmiDevFanReset(dvInd, sensorInd int) (err error) {
	ret := C.rsmi_dev_fan_reset(C.uint32_t(dvInd), C.uint32_t(sensorInd))
	glog.Info("rsmi_dev_fan_reset_ret:", ret)
	if err = errorString(ret); err != nil {
		return fmt.Errorf("Error rsmi_dev_fan_reset: %s", err)
	}
	return nil
}

// rsmiDevPowerProfileSet 设置设备功率配置文件
func rsmiDevPowerProfileSet(dvInd int, reserved int, profile RSNIPowerProfilePresetMasks) (err error) {
	ret := C.rsmi_dev_power_profile_set(C.uint32_t(dvInd), C.uint32_t(reserved), C.rsmi_power_profile_preset_masks_t(profile))
	glog.Info("rsmi_dev_power_profile_set ret:%v, retstr:%v", ret, errorString(ret))
	if err = errorString(ret); err != nil {
		glog.Errorf("Error rsmi_dev_power_profile_set:%v", err)
		return fmt.Errorf("Error rsmi_dev_power_profile_set:%s", err)
	}
	return
}

// rsmiDevXgmiErrorReset 重置设备的XGMI错误状态
func rsmiDevXgmiErrorReset(dvInd int) (err error) {
	ret := C.rsmi_dev_xgmi_error_reset(C.uint32_t(dvInd))
	glog.Infof(" rsmi_dev_xgmi_error_reset ret:%v,retStr:%v", ret, errorString(ret))
	if err = errorString(ret); err != nil {
		return fmt.Errorf("Error rsmiDevXgmiErrorReset:%s", err)
	}
	return
}

// rsmiDevXGMIErrorStatus 获取设备的XGMI错误状态
func rsmiDevXGMIErrorStatus(dvInd int) (status RSMIXGMIStatus, err error) {
	var cStatus C.rsmi_xgmi_status_t
	ret := C.rsmi_dev_xgmi_error_status(C.uint32_t(dvInd), &cStatus)
	glog.Infof(" rsmi_dev_xgmi_error_status ret:%v,retstr:%v", ret, errorString(ret))
	if err := errorString(ret); err != nil {
		return status, fmt.Errorf("Error RSMIDevXGMIErrorStatus: %s", err)
	}
	status = RSMIXGMIStatus(cStatus)
	glog.Infof("RSMIDevXGMIErrorStatus:%v", status)
	return
}

// rsmiDevXgmiHiveIdGet 获取设备的XGMI hive id
func rsmiDevXgmiHiveIdGet(dvInd int) (hiveId int64, err error) {
	var chiveId C.uint64_t
	ret := C.rsmi_dev_xgmi_hive_id_get(C.uint32_t(dvInd), &chiveId)
	glog.Infof("rsmi_dev_xgmi_hive_id_get ret:%v", ret)
	if err = errorString(ret); err != nil {
		return hiveId, fmt.Errorf("Error rsmiDevXgmiHiveIdGet:%s", err)
	}
	hiveId = int64(chiveId)
	glog.Infof("rsmi_dev_xgmi_hive_id_get hiveId:%v", hiveId)
	return
}

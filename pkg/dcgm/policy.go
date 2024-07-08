package dcgm

import "C"

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
)

type RSMIDevPerfLevel C.rsmi_dev_perf_level_t

const (
	RSMI_DEV_PERF_LEVEL_AUTO            RSMIDevPerfLevel = C.RSMI_DEV_PERF_LEVEL_AUTO
	RSMI_DEV_PERF_LEVEL_FIRST           RSMIDevPerfLevel = C.RSMI_DEV_PERF_LEVEL_FIRST
	RSMI_DEV_PERF_LEVEL_LOW             RSMIDevPerfLevel = C.RSMI_DEV_PERF_LEVEL_LOW
	RSMI_DEV_PERF_LEVEL_HIGH            RSMIDevPerfLevel = C.RSMI_DEV_PERF_LEVEL_HIGH
	RSMI_DEV_PERF_LEVEL_MANUAL          RSMIDevPerfLevel = C.RSMI_DEV_PERF_LEVEL_MANUAL
	RSMI_DEV_PERF_LEVEL_STABLE_STD      RSMIDevPerfLevel = C.RSMI_DEV_PERF_LEVEL_STABLE_STD
	RSMI_DEV_PERF_LEVEL_STABLE_PEAK     RSMIDevPerfLevel = C.RSMI_DEV_PERF_LEVEL_STABLE_PEAK
	RSMI_DEV_PERF_LEVEL_STABLE_MIN_MCLK RSMIDevPerfLevel = C.RSMI_DEV_PERF_LEVEL_STABLE_MIN_MCLK
	RSMI_DEV_PERF_LEVEL_STABLE_MIN_SCLK RSMIDevPerfLevel = C.RSMI_DEV_PERF_LEVEL_STABLE_MIN_SCLK
	RSMI_DEV_PERF_LEVEL_DETERMINISM     RSMIDevPerfLevel = C.RSMI_DEV_PERF_LEVEL_DETERMINISM
	RSMI_DEV_PERF_LEVEL_LAST            RSMIDevPerfLevel = C.RSMI_DEV_PERF_LEVEL_LAST
	RSMI_DEV_PERF_LEVEL_UNKNOWN         RSMIDevPerfLevel = C.RSMI_DEV_PERF_LEVEL_UNKNOWN
)

// go_rsmi_dev_perf_level_set 设置设备PowerPlay性能级别
func go_rsmi_dev_perf_level_set(dvInd int, devPerfLevel RSMIDevPerfLevel) (err error) {
	log.Println("dev_perf_level_set:", devPerfLevel)
	ret := C.rsmi_dev_perf_level_set(C.int32_t(dvInd), C.rsmi_dev_perf_level_t(devPerfLevel))
	if err = errorString(ret); err != nil {
		return fmt.Errorf("dev_perf_level_set:%s", err)
	}
	return
}

// go_rsmi_dev_clk_range_set 设置设备时钟范围信息
func go_rsmi_dev_clk_range_set(dvInd, minClkValue, maxClkValue uint64, clkType RSMIClkType) (err error) {
	ret := C.rsmi_dev_clk_range_set(C.int32_t(dvInd), C.int64_t(minClkValue), C.int64_t(maxClkValue), C.rsmi_clk_type_t(clkType))
	if err = errorString(ret); err != nil {
		return fmt.Errorf("Error rsmi_dev_clk_range_set:%s", err)
	}
	return
}

// go_rsmi_dev_od_volt_info_set 设置设备电压曲线点
func rsmi_dev_od_volt_info_set(dvInd, vPoint, clkValue, voltValue int) (err error) {
	ret := C.rsmi_dev_od_volt_info_set(C.int32_t(dvInd), C.int32_t(vPoint), C.uint64_t(clkValue), C.uint64_t(voltValue))
	if err = errorString(ret); err != nil {
		return fmt.Errorf("Error rsmi_dev_od_volt_info_set:%s", err)
	}
	return
}

// go_rsmi_dev_overdrive_level_set 设置设备超速百分比
func go_rsmi_dev_overdrive_level_set(dvInd, od int) (err error) {
	ret := C.rsmi_dev_overdrive_level_set(C.int32_t(dvInd), C.uint32_t(od))
	if err = errorString(ret); err != nil {
		return fmt.Errorf("Error rsmi_dev_overdrive_level_set:%s", err)
	}
	return
}

// go_rsmi_dev_gpu_clk_freq_set 设置可用于指定时钟的频率集
func go_rsmi_dev_gpu_clk_freq_set(dvInd int, clkType RSMIClkType, freqBitmask int64) (err error) {
	ret := C.rsmi_dev_gpu_clk_freq_set(C.int32_t(dvInd), C.rsmi_clk_type_t(clkType), C.uint64(freqBitmask))
	if err = errorString(ret); err != nil {
		return fmt.Errorf("Error rsmi_dev_gpu_clk_freq_set:%s", err)
	}
	return nil
}

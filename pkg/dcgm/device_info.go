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
	"encoding/json"
	"fmt"
	"unsafe"

	"github.com/golang/glog"
)

// rsmiNumMonitorDevices 获取gpu数量 *
func rsmiNumMonitorDevices() (gpuNum int, err error) {
	var p C.uint
	ret := C.rsmi_num_monitor_devices(&p)
	glog.Info("go_rsmi_num_monitor_devices_ret:", ret)
	if err = errorString(ret); err != nil {
		return 0, fmt.Errorf("Error go_rsmi_num_monitor_devices_ret: %s", err)
	}
	gpuNum = int(p)
	glog.Info("go_rsmi_num_monitor_devices:", gpuNum)
	return gpuNum, nil
}

// rsmiDevSkuGet 获取设备sku
func rsmiDevSkuGet(dvInd int) (sku int, err error) {
	var csku C.uint16_t
	ret := C.rsmi_dev_sku_get(C.uint32_t(dvInd), &csku)
	if err = errorString(ret); err != nil {
		return sku, err
	}
	sku = int(csku)
	glog.Info("rsmiDevSkuGet:", sku)
	return
}

// rsmiDevVendorIdGet 获取设备供应商id
func rsmiDevVendorIdGet(dvInd int) uint {
	var vid C.uint16_t
	C.rsmi_dev_vendor_id_get(C.uint32_t(dvInd), &vid)
	return uint(vid)
}

// rsmiDevIdGet 获取设备类型标识id
func rsmiDevIdGet(dvInd int) (id int, err error) {
	var cid C.uint16_t
	ret := C.rsmi_dev_id_get(C.uint32_t(dvInd), &cid)
	if err = errorString(ret); err != nil {
		return 0, fmt.Errorf("Error rsmiDevIdGet:%s", err)
	}
	id = int(cid)
	glog.Info("rsmiDevIdGet:", id, fmt.Sprintf("%x", id))
	return
}

// rsmiDevNameGet 获取设备名称
func rsmiDevNameGet(dvInd int) (nameStr string, err error) {
	name := make([]C.char, uint32(256))
	ret := C.rsmi_dev_name_get(C.uint32_t(dvInd), &name[0], 256)
	if err = errorString(ret); err != nil {
		return nameStr, fmt.Errorf("Error go_rsmi_dev_name_get: %s", err)
	}
	nameStr = C.GoString(&name[0])
	glog.Info("rsmiDevNameGet:", nameStr)
	return
}

// rsmiDevBrandGet 获取设备品牌名称
func rsmiDevBrandGet(dvInd int) (brand string, err error) {
	brands := make([]C.char, uint32(256))
	C.rsmi_dev_brand_get(C.uint32_t(dvInd), &brands[0], 256)
	brand = C.GoString(&brands[0])
	glog.Info("rsmiDevBrandGet:", brand)
	return
}

// rsmiDevVendorNameGet 获取设备供应商名称
func rsmiDevVendorNameGet(dvInd int) string {
	bname := make([]C.char, uint32(256))
	C.rsmi_dev_vendor_name_get(C.uint32_t(dvInd), &bname[0], 80)
	result := C.GoString(&bname[0])
	fmt.Println("rsmiDevVendorNameGet:", result)
	return result
}

// rsmiDevVramVendorGet 获取设备显存供应商名称
func rsmiDevVramVendorGet(dvInd int) string {
	bname := make([]C.char, uint32(256))
	C.rsmi_dev_vram_vendor_get(C.uint32_t(dvInd), &bname[0], 80)
	result := C.GoString(&bname[0])
	fmt.Printf("rsmiDevVramVendorGet: %s\n", result)
	return result
}

// rsmiDevSerialNumberGet 获取设备序列号
func rsmiDevSerialNumberGet(dvInd int) string {
	serialNumber := make([]C.char, uint32(256))
	C.rsmi_dev_serial_number_get(C.uint32_t(dvInd), &serialNumber[0], 256)
	result := C.GoString(&serialNumber[0])
	return result
}

// rsmiDevSubsystemIdGet 获取设备子系统id
func rsmiDevSubsystemIdGet(dvInd int) int {
	var id C.uint16_t
	C.rsmi_dev_subsystem_id_get(C.uint32_t(dvInd), &id)
	return int(id)
}

// rsmiDevSubsystemNameGet 获取设备子系统名称
func rsmiDevSubsystemNameGet(dvInd int) string {
	subSystemName := make([]C.char, uint32(256))
	C.rsmi_dev_subsystem_name_get(C.uint32_t(dvInd), &subSystemName[0], 256)
	result := C.GoString(&subSystemName[0])
	fmt.Printf("rsmiDevSubsystemNameGet:%s\n", result)
	return result
}

// rsmiDevDrmRenderMinorGet 获取设备drm次编号
func rsmiDevDrmRenderMinorGet(dvInd int) int {
	var id C.uint32_t
	C.rsmi_dev_drm_render_minor_get(C.uint32_t(dvInd), &id)
	return int(id)
}

// rsmiDevUniqueIdGet 获取设备唯一id
func rsmiDevUniqueIdGet(dvInd int) int64 {
	var uniqueId C.uint64_t
	C.rsmi_dev_unique_id_get(C.uint32_t(dvInd), &uniqueId)
	return int64(uniqueId)
}

// rsmiDevSubsystemVendorIdGet 获取设备子系统供应商id
func rsmiDevSubsystemVendorIdGet(dvInd int) int {
	var id C.uint16_t
	C.rsmi_dev_subsystem_vendor_id_get(C.uint32_t(dvInd), &id)
	return int(id)
}

/****************************************** PCIe *********************************************/

// rsmiDevPciBandwidthGet 获取可用的pcie带宽列表
func rsmiDevPciBandwidthGet(dvInd int) RSMIPcieBandwidth {
	var bandwidth C.rsmi_pcie_bandwidth_t
	C.rsmi_dev_pci_bandwidth_get(C.uint32_t(dvInd), &bandwidth)
	rsmiPcieBandwidth := RSMIPcieBandwidth{
		TransferRate: RSMIFrequencies{
			NumSupported: uint32(bandwidth.transfer_rate.num_supported),
			Current:      uint32(bandwidth.transfer_rate.current),
			Frequency:    *(*[32]uint64)(unsafe.Pointer(&bandwidth.transfer_rate)),
		},
		lanes: *(*[32]uint32)(unsafe.Pointer(&bandwidth.lanes)),
	}
	glog.Info("RSMIPcieBandwidth:%s", rsmiPcieBandwidth)
	return rsmiPcieBandwidth
}

// rsmiDevPciIdGet 获取唯一pci设备标识符
func rsmiDevPciIdGet(dvInd int) (bdfid int64, err error) {
	var cbdfid C.uint64_t
	ret := C.rsmi_dev_pci_id_get(C.uint32_t(dvInd), &cbdfid)
	if err = errorString(ret); err != nil {
		return bdfid, err
	}
	bdfid = int64(cbdfid)
	return
}

// rsmiTopoNumaAffinityGet 获取与设备关联的numa节点
func rsmiTopoNumaAffinityGet(dvInd int) int {
	var namaNode C.uint32_t
	C.rsmi_topo_numa_affinity_get(C.uint32_t(dvInd), &namaNode)
	return int(namaNode)
}

// rsmiDevPciThroughputGet 获取pcie流量信息
func rsmiDevPciThroughputGet(dvInd int) (sent int64, received int64, maxPktSz int64) {
	var csent, creceived, cmaxpktsz C.uint64_t
	C.rsmi_dev_pci_throughput_get(C.uint32_t(dvInd), &csent, &creceived, &cmaxpktsz)
	sent = int64(cmaxpktsz)
	received = int64(csent)
	maxPktSz = int64(creceived)
	glog.Info("sent: %d, received: %d, maxPktSz: %d\n", sent, received, maxPktSz)
	return
}

// rsmiDevPciReplayCounterGet 获取pcie重放计数
func rsmiDevPciReplayCounterGet(dvInd int) uint64 {
	var counter C.uint64_t
	C.rsmi_dev_pci_replay_counter_get(C.uint32_t(dvInd), &counter)
	return uint64(counter)
}

// rsmiDevPciBandwidthSet 设置可使用的pcie带宽集
func rsmiDevPciBandwidthSet(dvInd int, bwBitmask int64) {
	C.rsmi_dev_pci_bandwidth_set(C.uint32_t(dvInd), C.uint64_t(bwBitmask))
}

/****************************************** Power *********************************************/

// rsmiDevPowerAveGet 获取设备平均功耗
func rsmiDevPowerAveGet(dvInd int, senserId int) int64 {
	var power C.uint64_t
	C.rsmi_dev_power_ave_get(C.uint32_t(dvInd), C.uint32_t(senserId), &power)
	return int64(power)
}

// rsmiDevEnergyCountGet 获取设备的能量累加计数
func rsmiDevEnergyCountGet() {

}

// rsmiDevPowerCapGet 获取设备功率上限
func rsmiDevPowerCapGet(dvInd int, senserId int) (power int64, err error) {
	var cpower C.uint64_t
	ret := C.rsmi_dev_power_cap_get(C.uint32_t(dvInd), C.uint32_t(senserId), &cpower)
	if err = errorString(ret); err != nil {
		return power, fmt.Errorf("Error rsmiDevPowerCapGet:%s", err)
	}
	power = int64(cpower)
	return
}

// rsmiDevPowerCapRangeGet 获取设备功率有效值范围
func rsmiDevPowerCapRangeGet(dvInd int, senserId int) (max, min int64) {
	var cmax, cmin C.uint64_t
	C.rsmi_dev_power_cap_range_get(C.uint32_t(dvInd), C.uint32_t(senserId), &cmax, &cmin)
	max, min = int64(cmax), int64(cmin)
	return
}

// rsmiDevPowerProfileSet 设置设备功率配置文件
func rsmiDevPowerProfileSet(dvInd int, reserved int, profile RSNIPowerProfilePresetMasks) {
	C.rsmi_dev_power_profile_set(C.uint32_t(dvInd), C.uint32_t(reserved), C.rsmi_power_profile_preset_masks_t(profile))
}

/****************************************** Memory *********************************************/

// rsmiDevMemoryTotalGet 获取设备内存总量 *
func rsmiDevMemoryTotalGet(dvInd int, memoryType RSMIMemoryType) (total int64, err error) {
	var ctotal C.uint64_t
	ret := C.rsmi_dev_memory_total_get(C.uint32_t(dvInd), C.rsmi_memory_type_t(memoryType), &ctotal)
	if err = errorString(ret); err != nil {
		return total, fmt.Errorf("Error rsmiDevMemoryTotalGet:%s", err)
	}
	total = int64(ctotal)
	glog.Info("memory_total:", total)
	return
}

// rsmiDevMemoryUsageGet 获取当前设备内存使用情况 *
func rsmiDevMemoryUsageGet(dvInd int, memoryType RSMIMemoryType) (used int64, err error) {
	var cused C.uint64_t
	ret := C.rsmi_dev_memory_usage_get(C.uint32_t(dvInd), C.rsmi_memory_type_t(memoryType), &cused)
	if err = errorString(ret); err != nil {
		return used, fmt.Errorf("Error rsmiDevMemoryUsageGet:%s", err)
	}
	used = int64(cused)
	glog.Info("memory_used:", used)
	return
}

// rsmiDevMemoryBusyPercentGet 获取设备内存使用的百分比
func rsmiDevMemoryBusyPercentGet(dvInd int) int {
	var busyPercent C.uint32_t
	C.rsmi_dev_memory_busy_percent_get(C.uint32_t(dvInd), &busyPercent)
	glog.Info("busy_percent:", busyPercent)
	return int(busyPercent)
}

// rsmiDevMemoryReservedPagesGet 获取有关保留的(“已退休”)内存页的信息
func rsmiDevMemoryReservedPagesGet(dvInd int) (numPages int, records []RSMIRetiredPageRecord, err error) {
	var cnumPages C.uint32_t
	ret := C.rsmi_dev_memory_reserved_pages_get(C.uint32_t(dvInd), &cnumPages, nil)
	if ret != 0 {
		return 0, nil, fmt.Errorf("failed to get the number of pages, error code: %d", ret)
	}
	glog.Info("cnumPages:", cnumPages)
	glog.Info("cnumPages:", int(cnumPages))
	numPages = int(cnumPages)
	if numPages == 0 {
		return 0, nil, nil // No pages to retrieve
	}
	cRecords := make([]C.rsmi_retired_page_record_t, numPages)
	ret = C.rsmi_dev_memory_reserved_pages_get(C.uint32_t(dvInd), &cnumPages, (*C.rsmi_retired_page_record_t)(unsafe.Pointer(&cRecords[0])))
	if ret != 0 {
		return 0, nil, fmt.Errorf("failed to get the page records, error code: %d", ret)
	}

	records = make([]RSMIRetiredPageRecord, numPages)
	for i, rec := range cRecords {
		records[i] = RSMIRetiredPageRecord{
			PageAddress: uint64(rec.page_address),
			PageSize:    uint64(rec.page_size),
			Status:      RSMIMemoryPageStatus(rec.status),
		}
	}
	indent, _ := json.MarshalIndent(records, "", "  ")
	glog.Info("records:", indent)
	return
}

// rsmiDevFanRpmsGet 获取设备的风扇速度，实际转速
func rsmiDevFanRpmsGet(dvInd, sensorInd int) int64 {
	var speed C.int64_t
	C.rsmi_dev_fan_rpms_get(C.uint32_t(dvInd), C.uint32_t(sensorInd), &speed)
	return int64(speed)
}

// rsmiDevFanSpeedGet 获取设备的风扇速度，相对速度值
func rsmiDevFanSpeedGet(dvInd, sensorInd int) (speed int64, err error) {
	var cspeed C.int64_t
	ret := C.rsmi_dev_fan_speed_get(C.uint32_t(dvInd), C.uint32_t(sensorInd), &cspeed)
	if err = errorString(ret); err != nil {
		return speed, fmt.Errorf("Error rsmiDevFanSpeedGet:%s", err)
	}
	speed = int64(cspeed)
	return
}

// rsmiDevFanSpeedMaxGet 获取设备的风扇速度，最大风速
func rsmiDevFanSpeedMaxGet(dvInd, sensorInd int) (maxSpeed int64, err error) {
	var cmaxSpeed C.uint64_t
	ret := C.rsmi_dev_fan_speed_max_get(C.uint32_t(dvInd), C.uint32_t(sensorInd), &cmaxSpeed)
	if err = errorString(ret); err != nil {
		return maxSpeed, fmt.Errorf("Error rsmiDevFanSpeedMaxGet:%s", err)
	}
	maxSpeed = int64(cmaxSpeed)
	return
}

// rsmiDevOdVoltCurveRegionsGet
func rsmiDevOdVoltCurveRegionsGet(dvInd int) (numRegions int, buffer RSMIFreqVoltRegion, err error) {
	var cnumRegions C.uint32_t
	var cbuffer C.rsmi_freq_volt_region_t
	ret := C.rsmi_dev_od_volt_curve_regions_get(C.uint32_t(dvInd), &cnumRegions, &cbuffer)
	if err = errorString(ret); err != nil {
		return 0, buffer, fmt.Errorf("Error dev_od_volt_curve_regions_get:%S", err)
	}
	numRegions = int(cnumRegions)
	buffer = RSMIFreqVoltRegion{
		FreqRange: RSMIRange{
			LowerBound: uint64(cbuffer.freq_range.lower_bound),
			UpperBound: uint64(cbuffer.freq_range.upper_bound),
		},
		VoltRange: RSMIRange{
			LowerBound: uint64(cbuffer.freq_range.lower_bound),
			UpperBound: uint64(cbuffer.freq_range.upper_bound),
		},
	}
	return
}

// rsmiDevPowerProfilePresetsGet 获取可用预设电源配置文件列表并指示当前活动的配置文件
func rsmiDevPowerProfilePresetsGet(dvInd, sensorInd int) (powerProfileStatus RSMPowerProfileStatus, err error) {
	var cpowerProfileStatus C.rsmi_power_profile_status_t
	ret := C.rsmi_dev_power_profile_presets_get(C.uint32_t(dvInd), C.uint32_t(sensorInd), &cpowerProfileStatus)
	if err = errorString(ret); err != nil {
		return powerProfileStatus, fmt.Errorf("Error dev_power_profile_presets_get:%s", err)
	}
	powerProfileStatus = RSMPowerProfileStatus{
		AvailableProfiles: RSMIBitField(cpowerProfileStatus.available_profiles),
		Current:           RSMIPowerProfilePresetMasks(cpowerProfileStatus.current),
		NumProfiles:       uint32(cpowerProfileStatus.num_profiles),
	}
	return
}

// rsmiVersionGet 获取当前运行的RSMI版本
func rsmiVersionGet() (version RSMIVersion, err error) {

	var cVersion C.rsmi_version_t
	ret := C.rsmi_version_get(&cVersion)
	if err = errorString(ret); err != nil {
		return version, fmt.Errorf("Error to get version: %s", err)
	}
	version = RSMIVersion{
		Major: uint32(cVersion.major),
		Minor: uint32(cVersion.minor),
		Patch: uint32(cVersion.patch),
		Build: C.GoString(cVersion.build),
	}
	return
}

// rsmiVersionStrGet 获取当前系统的驱动程序版本
func rsmiVersionStrGet(component RSMISwComponent, len int) (varStr string, err error) {
	var cvarStr C.char
	ret := C.rsmi_version_str_get(C.rsmi_sw_component_t(component), &cvarStr, C.uint32_t(len))
	if err = errorString(ret); err != nil {
		return "", fmt.Errorf("Error rsmi_version_str_get:%s", err)
	}
	varStr = string(cvarStr)
	return
}

// rsmiDevVbiosVersionGet 获取VBIOS版本
func rsmiDevVbiosVersionGet(dvInd, len int) (vbios string, err error) {
	var cvbios C.char
	ret := C.rsmi_dev_vbios_version_get(C.uint32_t(dvInd), &cvbios, C.uint32_t(len))
	if err = errorString(ret); err != nil {
		return vbios, fmt.Errorf("Error rsmi_dev_vbios_version_get:%s", err)
	}
	vbios = string(cvbios)
	return
}

// rsmiDevFirmwareVersionGet 获取设备的固件版本
func rsmiDevFirmwareVersionGet(dvInd int, fwBlock RSMIFwBlock) (fwVersion int64, err error) {
	var cfwBlock C.uint64_t
	ret := C.rsmi_dev_firmware_version_get(C.uint32_t(dvInd), C.rsmi_fw_block_t(fwBlock), &cfwBlock)
	if err = errorString(ret); err != nil {
		return fwVersion, fmt.Errorf("Error rsmi_dev_firmware_version_get:%s", err)
	}
	fwVersion = int64(cfwBlock)
	return
}

// rsmiDevEccCountGet 获取GPU块的错误计数
func rsmiDevEccCountGet(dvInd int, gpuBlock RSMIGpuBlock) (errorCount RSMIErrorCount, err error) {
	var cerrorCount C.rsmi_error_count_t
	ret := C.rsmi_dev_ecc_count_get(C.uint32_t(dvInd), C.rsmi_gpu_block_t(gpuBlock), &cerrorCount)
	if err = errorString(ret); err != nil {
		return errorCount, fmt.Errorf("Error rsmi_dev_ecc_count_get:%s", err)
	}
	errorCount = RSMIErrorCount{
		CorrectableErr:   uint64(cerrorCount.correctable_err),
		UncorrectableErr: uint64(cerrorCount.uncorrectable_err),
	}
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
	return
}

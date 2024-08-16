package dcgm

/*
#cgo CFLAGS: -Wall -I./include
#cgo LDFLAGS: -L./lib -lrocm_smi64 -lhydmi -Wl,--unresolved-symbols=ignore-in-object-files
#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <kfd_ioctl.h>
#include <rocm_smi64Config.h>
#include <rocm_smi.h>
#include <dmi_virtual.h>
#include <dmi_error.h>
#include <dmi.h>
#include <dmi_mig.h>
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
	//glog.Info("rsmiDevNameGet:", nameStr)
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
func rsmiDevVendorNameGet(dvInd int) (bname string, err error) {
	cbname := make([]C.char, uint32(256))
	ret := C.rsmi_dev_vendor_name_get(C.uint32_t(dvInd), &cbname[0], 80)
	if err = errorString(ret); err != nil {
		glog.Errorf("Error rsmi_dev_vendor_name_get:%v", err)
		return bname, fmt.Errorf("Error rsmi_dev_vendor_name_get:%v", err)
	}
	bname = C.GoString(&cbname[0])
	//glog.Infof("rsmiDevVendorNameGet:%v", bname)
	return
}

// rsmiDevVramVendorGet 获取设备显存供应商名称
func rsmiDevVramVendorGet(dvInd int) (result string, err error) {
	bname := make([]C.char, uint32(256))
	ret := C.rsmi_dev_vram_vendor_get(C.uint32_t(dvInd), &bname[0], 80)
	if err = errorString(ret); err != nil {
		return "", fmt.Errorf("Error rsmi_dev_vram_vendor_get:%s", err)
	}
	result = C.GoString(&bname[0])
	glog.Infof("rsmiDevVramVendorGet: %v", result)
	return
}

// rsmiDevSerialNumberGet 获取设备序列号
func rsmiDevSerialNumberGet(dvInd int) (serialNumber string, err error) {
	cserialNumber := make([]C.char, uint32(256))
	ret := C.rsmi_dev_serial_number_get(C.uint32_t(dvInd), &cserialNumber[0], 256)
	if err = errorString(ret); err != nil {
		glog.Errorf("Error rsmi_dev_serial_number_get:%v, errstr:%v", err, errorString(ret))
		return "", fmt.Errorf("Error rsmi_dev_serial_number_get:%s", err)
	}
	serialNumber = C.GoString(&cserialNumber[0])
	return
}

// rsmiDevSubsystemIdGet 获取设备子系统id
func rsmiDevSubsystemIdGet(dvInd int) int {
	var id C.uint16_t
	C.rsmi_dev_subsystem_id_get(C.uint32_t(dvInd), &id)
	return int(id)
}

// rsmiDevSubsystemNameGet 获取设备子系统名称
func rsmiDevSubsystemNameGet(dvInd int) (subSystemName string, err error) {
	csubSystemName := make([]C.char, uint32(256))
	ret := C.rsmi_dev_subsystem_name_get(C.uint32_t(dvInd), &csubSystemName[0], 256)
	if err = errorString(ret); err != nil {
		glog.Errorf("Error rsmi_dev_subsystem_name_get:%v", err)
		return subSystemName, fmt.Errorf("Error rsmi_dev_subsystem_name_get:%s", err)
	}
	subSystemName = C.GoString(&csubSystemName[0])
	//glog.Infof("rsmiDevSubsystemNameGet:%v", subSystemName)
	return
}

// rsmiDevDrmRenderMinorGet 获取设备drm次编号
func rsmiDevDrmRenderMinorGet(dvInd int) int {
	var id C.uint32_t
	C.rsmi_dev_drm_render_minor_get(C.uint32_t(dvInd), &id)
	return int(id)
}

// rsmiDevUniqueIdGet 获取设备唯一id
func rsmiDevUniqueIdGet(dvInd int) (uniqueId int64, err error) {
	var cuniqueId C.uint64_t
	ret := C.rsmi_dev_unique_id_get(C.uint32_t(dvInd), &cuniqueId)
	if err = errorString(ret); err != nil {
		glog.Errorf("Error rsmi_dev_unique_id_get:%v, retstr:%v", ret, errorString(ret))
		return uniqueId, fmt.Errorf("Error rsmi_dev_unique_id_get:%s", err)
	}
	uniqueId = int64(cuniqueId)
	return
}

// rsmiDevSubsystemVendorIdGet 获取设备子系统供应商id
func rsmiDevSubsystemVendorIdGet(dvInd int) int {
	var id C.uint16_t
	C.rsmi_dev_subsystem_vendor_id_get(C.uint32_t(dvInd), &id)
	return int(id)
}

/****************************************** PCIe *********************************************/

// rsmiDevPciBandwidthGet 获取可用的pcie带宽列表
func rsmiDevPciBandwidthGet(dvInd int) (rsmiPcieBandwidth RSMIPcieBandwidth, err error) {
	var bandwidth C.rsmi_pcie_bandwidth_t
	ret := C.rsmi_dev_pci_bandwidth_get(C.uint32_t(dvInd), &bandwidth)
	if err = errorString(ret); err != nil {
		return rsmiPcieBandwidth, fmt.Errorf("Error rsmi_dev_pci_bandwidth_get%s", err)
	}
	rsmiPcieBandwidth = RSMIPcieBandwidth{
		TransferRate: RSMIFrequencies{
			NumSupported: uint32(bandwidth.transfer_rate.num_supported),
			Current:      uint32(bandwidth.transfer_rate.current),
			Frequency:    *(*[32]uint64)(unsafe.Pointer(&bandwidth.transfer_rate)),
		},
		lanes: *(*[32]uint32)(unsafe.Pointer(&bandwidth.lanes)),
	}
	glog.Infof("RSMIPcieBandwidth:%v", dataToJson(rsmiPcieBandwidth))
	return
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
func rsmiTopoNumaAffinityGet(dvInd int) (namaNode int, err error) {
	var cnamaNode C.uint32_t
	ret := C.rsmi_topo_numa_affinity_get(C.uint32_t(dvInd), &cnamaNode)
	if err = errorString(ret); err != nil {
		glog.Errorf("Error rsmi_topo_numa_affinity_get ret:%v, retstr:%v", ret, errorString(ret))
		return namaNode, fmt.Errorf("Error rsmi_topo_numa_affinity_get:%s", err)
	}
	namaNode = int(cnamaNode)
	return
}

// rsmiDevPciThroughputGet 获取pcie流量信息
func rsmiDevPciThroughputGet(dvInd int) (sent int64, received int64, maxPktSz int64, err error) {
	var csent, creceived, cmaxpktsz C.uint64_t
	ret := C.rsmi_dev_pci_throughput_get(C.uint32_t(dvInd), &csent, &creceived, &cmaxpktsz)
	if err = errorString(ret); err != nil {
		return 0, 0, 0, fmt.Errorf("Error rsmi_dev_pci_throughput_get:%s", err)
	}
	sent = int64(cmaxpktsz)
	received = int64(csent)
	maxPktSz = int64(creceived)
	glog.Infof("sent: %v, received: %v, maxPktSz: %v", sent, received, maxPktSz)
	return
}

// rsmiDevPciReplayCounterGet 获取pcie重放计数
func rsmiDevPciReplayCounterGet(dvInd int) (counter int64, err error) {
	var ccounter C.uint64_t
	ret := C.rsmi_dev_pci_replay_counter_get(C.uint32_t(dvInd), &ccounter)
	if err = errorString(ret); err != nil {
		return counter, fmt.Errorf("Error rsmi_dev_pci_replay_counter_get:%s", err)
	}
	counter = int64(ccounter)
	glog.Infof("counter:%v", ccounter)
	return
}

// rsmiDevPciBandwidthSet 设置可使用的pcie带宽集
func rsmiDevPciBandwidthSet(dvInd int, bwBitmask int64) {
	C.rsmi_dev_pci_bandwidth_set(C.uint32_t(dvInd), C.uint64_t(bwBitmask))
}

/****************************************** Power *********************************************/

// rsmiDevPowerAveGet 获取设备平均功耗
func rsmiDevPowerAveGet(dvInd int, senserId int) (power int64, err error) {
	var cpower C.uint64_t
	ret := C.rsmi_dev_power_ave_get(C.uint32_t(dvInd), C.uint32_t(senserId), &cpower)
	if err = errorString(ret); err != nil {
		return power, fmt.Errorf("Error rsmiDevPowerAveGet:%v", err)
	}
	power = int64(cpower)
	return
}

// rsmiDevEnergyCountGet 获取设备的能量累加计数
func rsmiDevEnergyCountGet(dvInd int) (power uint64, counterResolution float32, timestamp uint64, err error) {
	var cPower C.uint64_t
	var cCounterResolution C.float
	var cTimestamp C.uint64_t
	ret := C.rsmi_dev_energy_count_get(C.uint32_t(dvInd), &cPower, &cCounterResolution, &cTimestamp)
	if ret != C.RSMI_STATUS_SUCCESS {
		return 0, 0, 0, fmt.Errorf("Error in rsmi_dev_energy_count_get: %s", errorString(ret))
	}
	return uint64(cPower), float32(cCounterResolution), uint64(cTimestamp), nil
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

/****************************************** Memory *********************************************/

// rsmiDevMemoryTotalGet 获取设备内存总量 *
func rsmiDevMemoryTotalGet(dvInd int, memoryType RSMIMemoryType) (total int64, err error) {
	var ctotal C.uint64_t
	ret := C.rsmi_dev_memory_total_get(C.uint32_t(dvInd), C.rsmi_memory_type_t(memoryType), &ctotal)
	if err = errorString(ret); err != nil {
		return total, fmt.Errorf("Error rsmiDevMemoryTotalGet:%s", err)
	}
	total = int64(ctotal)
	//glog.Info("memory_total:", total)
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
	//glog.Info("memory_used:", used)
	return
}

// rsmiDevMemoryBusyPercentGet 获取设备内存使用的百分比
func rsmiDevMemoryBusyPercentGet(dvInd int) (busyPercent int, err error) {
	var cbusyPercent C.uint32_t
	ret := C.rsmi_dev_memory_busy_percent_get(C.uint32_t(dvInd), &cbusyPercent)
	if err = errorString(ret); err != nil {
		return busyPercent, fmt.Errorf("Error rsmi_dev_memory_busy_percent_get:%s", err)
	}
	busyPercent = int(cbusyPercent)
	glog.Info("busy_percent:", busyPercent)
	return
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
func rsmiDevFanRpmsGet(dvInd, sensorInd int) (speed int64, err error) {
	var cspeed C.int64_t
	ret := C.rsmi_dev_fan_rpms_get(C.uint32_t(dvInd), C.uint32_t(sensorInd), &cspeed)
	if err = errorString(ret); err != nil {
		return speed, fmt.Errorf("Error rsmi_dev_fan_rpms_get:%s", err)
	}
	speed = int64(cspeed)
	glog.Infof("rsmi_dev_fan_rpms_get speed value: %v", speed)
	return
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
func rsmiDevOdVoltCurveRegionsGet(dvInd int) (numRegions int, regions []RSMIFreqVoltRegion, err error) {
	var cnumRegions C.uint32_t
	ret := C.rsmi_dev_od_volt_curve_regions_get(C.uint32_t(dvInd), &cnumRegions, nil)
	if err = errorString(ret); err != nil {
		return 0, nil, fmt.Errorf("Error dev_od_volt_curve_regions_get: %v", err)
	}

	cbuffer := make([]C.rsmi_freq_volt_region_t, cnumRegions)
	ret = C.rsmi_dev_od_volt_curve_regions_get(C.uint32_t(dvInd), &cnumRegions, &cbuffer[0])
	if err = errorString(ret); err != nil {
		return 0, nil, fmt.Errorf("Error dev_od_volt_curve_regions_get: %v", err)
	}

	regions = make([]RSMIFreqVoltRegion, cnumRegions)
	for i := 0; i < int(cnumRegions); i++ {
		regions[i] = RSMIFreqVoltRegion{
			FreqRange: RSMIRange{
				LowerBound: uint64(cbuffer[i].freq_range.lower_bound),
				UpperBound: uint64(cbuffer[i].freq_range.upper_bound),
			},
			VoltRange: RSMIRange{
				LowerBound: uint64(cbuffer[i].volt_range.lower_bound),
				UpperBound: uint64(cbuffer[i].volt_range.upper_bound),
			},
		}
	}
	numRegions = int(cnumRegions)
	return
}

// rsmiDevPowerProfilePresetsGet 获取可用预设电源配置文件列表并指示当前活动的配置文件
func rsmiDevPowerProfilePresetsGet(dvInd, sensorInd int) (powerProfileStatus RSMPowerProfileStatus, err error) {
	var cpowerProfileStatus C.rsmi_power_profile_status_t
	ret := C.rsmi_dev_power_profile_presets_get(C.uint32_t(dvInd), C.uint32_t(sensorInd), &cpowerProfileStatus)
	glog.Infof("rsmi_dev_power_profile_presets_get ret:%v, retstr:%v", ret, errorString(ret))
	if err = errorString(ret); err != nil {
		return powerProfileStatus, fmt.Errorf("Error dev_power_profile_presets_get:%s", err)
	}
	powerProfileStatus = RSMPowerProfileStatus{
		AvailableProfiles: RSMIBitField(cpowerProfileStatus.available_profiles),
		Current:           RSMIPowerProfilePresetMasks(cpowerProfileStatus.current),
		NumProfiles:       uint32(cpowerProfileStatus.num_profiles),
	}
	glog.Infof("powerProfileStatus value: %v", dataToJson(powerProfileStatus))
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
	cvarStr := make([]C.char, len)
	ret := C.rsmi_version_str_get(C.rsmi_sw_component_t(component), &cvarStr[0], C.uint32_t(len))
	if err = errorString(ret); err != nil {
		return "", fmt.Errorf("Error rsmi_version_str_get:%s", err)
	}
	varStr = C.GoString(&cvarStr[0])
	return
}

// rsmiDevVbiosVersionGet 获取VBIOS版本
func rsmiDevVbiosVersionGet(dvInd, len int) (vbios string, err error) {
	cvbios := make([]C.char, len)
	ret := C.rsmi_dev_vbios_version_get(C.uint32_t(dvInd), &cvbios[0], C.uint32_t(len))
	if err = errorString(ret); err != nil {
		return vbios, fmt.Errorf("Error rsmi_dev_vbios_version_get:%s", err)
	}
	vbios = C.GoString(&cvbios[0])
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

/*************************************VDCU******************************************/
// 设备数量
func dmiGetDeviceCount() (count int, err error) {
	var ccount C.int
	ret := C.dmiGetDeviceCount(&ccount)
	glog.Infof("dmiGetDeviceCount:%v", ret)
	if err = dmiErrorString(ret); err != nil {
		return 0, fmt.Errorf("Error vDeviceCount:%s", err)
	}
	count = int(ccount)
	glog.Infof("dmiDeviceCount:%v", count)
	return
}

// 设备信息
func dmiGetDeviceInfo(dvInd int) (deviceInfo DMIDeviceInfo, err error) {
	var cdeviceInfo C.dmiDeviceInfo
	ret := C.dmiGetDeviceInfo(C.int(dvInd), &cdeviceInfo)
	glog.Infof("dmiDeviceInfo ret:%v", ret)
	if err = dmiErrorString(ret); err != nil {
		return deviceInfo, fmt.Errorf("Error dmiGetDeviceInfo:%s", err)
	}
	// 创建一个新的变量来存储 name 字段
	var deviceName [DMI_NAME_SIZE]byte
	for i := 0; i < DMI_NAME_SIZE; i++ {
		deviceName[i] = byte(cdeviceInfo.name[i])
	}
	glog.Infof("deviceName:%v", deviceName)
	deviceInfo = DMIDeviceInfo{
		ComputeUnitCount: int(cdeviceInfo.compute_unit_count),
		GlobalMemSize:    uintptr(cdeviceInfo.global_mem_size),
		UsageMemSize:     uintptr(cdeviceInfo.usage_mem_size),
		DeviceID:         int(cdeviceInfo.device_id),
		Name:             ConvertASCIIToString(deviceName[:]),
	}

	glog.Infof("DeviceInfo: %v", dataToJson(deviceInfo))
	return
}

// 物理设备支持最大虚拟化设备数量
func dmiGetMaxVDeviceCount() (count int, err error) {
	var ccount C.int
	ret := C.dmiGetMaxVDeviceCount(&ccount)
	if err = dmiErrorString(ret); err != nil {
		return 0, fmt.Errorf("Error dmiGetMaxVDeviceCount:%s", err)
	}
	count = int(ccount)
	return
}

// 虚拟设备数量
func dmiGetVDeviceCount() (count int, err error) {
	var ccount C.int
	ret := C.dmiGetVDeviceCount(&ccount)
	if err = dmiErrorString(ret); err != nil {
		return 0, fmt.Errorf("Error dmiGetVDeviceCount:%s", err)
	}
	count = int(ccount)
	glog.Infof("dmiGetVDeviceCount:%v", count)
	return
}

// 虚拟设备信息
func dmiGetVDeviceInfo(vDvInd int) (vDeviceInfo DMIVDeviceInfo, err error) {
	var cvDeviceInfo C.dmiDeviceInfo
	ret := C.dmiGetVDeviceInfo(C.int(vDvInd), &cvDeviceInfo)
	glog.Infof("dmiGetVDeviceInfo ret:%v", ret)
	glog.Infof("cgo cvDeviceInfo:%v", dataToJson(cvDeviceInfo))
	if err = dmiErrorString(ret); err != nil {
		return vDeviceInfo, fmt.Errorf("Error dmiGetVDeviceInfo:%s", err)
	}
	// 创建一个新的变量来存储 name 字段
	var deviceName [DMI_NAME_SIZE]byte
	for i := 0; i < DMI_NAME_SIZE; i++ {
		deviceName[i] = byte(cvDeviceInfo.name[i])
	}
	vDeviceInfo = DMIVDeviceInfo{
		ComputeUnitCount: int(cvDeviceInfo.compute_unit_count),
		GlobalMemSize:    uintptr(cvDeviceInfo.global_mem_size),
		UsageMemSize:     uintptr(cvDeviceInfo.usage_mem_size),
		ContainerID:      uint64(cvDeviceInfo.container_id),
		DeviceID:         int(cvDeviceInfo.device_id),
		Name:             ConvertASCIIToString(deviceName[:]),
	}
	glog.Infof("vDeviceInfo: %v", dataToJson(vDeviceInfo))
	return
}

// 指定物理设备剩余的CU和内存
func dmiGetDeviceRemainingInfo(dvInd int) (cus, memories uint64, err error) {
	var ccus, cmemories C.size_t
	ret := C.dmiGetDeviceRemainingInfo(C.int(dvInd), &ccus, &cmemories)
	glog.Infof("dmiGetDeviceRemainingInfo ret:%v", ret)
	if err = dmiErrorString(ret); err != nil {
		return cus, memories, fmt.Errorf("Error dmiGetDeviceRemainingInfo:%s", err)
	}
	cus = uint64(ccus)
	memories = uint64(cmemories)
	glog.Infof("cus:%v,memories:%v", cus, memories)
	return
}

// 创建指定数量的虚拟设备
//
//	deviceID := 0
//	vdevCount := 2
//	vdevCUs := []int{4, 4}
//	vdevMemSize := []int{1024, 2048}
//
// 物理设备 ID: 0
//
//	├── 虚拟设备 1
//	│    ├── 计算单元: 4
//	│    └── 内存大小: 1024 字节
//	└── 虚拟设备 2
//	     ├── 计算单元: 4
//	     └── 内存大小: 2048 字节
func dmiCreateVDevices(dvInd int, vDevCount int, vDevCUs []int, vDevMemSize []int) (err error) {
	if len(vDevCUs) != vDevCount || len(vDevMemSize) != vDevCount {
		return fmt.Errorf("Invalid args")
	}

	fmt.Printf("deviceID: %d, vDevCount: %d, vDevCUs: %v, vDevMemSize: %v\n", dvInd, vDevCount, vDevCUs, vDevMemSize)
	// Allocate C arrays from Go slices
	cVdevCus := (*C.int)(C.malloc(C.size_t(len(vDevCUs)) * C.sizeof_int))
	cVdevMemSize := (*C.int)(C.malloc(C.size_t(len(vDevMemSize)) * C.sizeof_int))

	if cVdevCus == nil || cVdevMemSize == nil {
		return fmt.Errorf("Memory allocation failed")
	}
	defer C.free(unsafe.Pointer(cVdevCus))
	defer C.free(unsafe.Pointer(cVdevMemSize))

	// Copy values from Go slices to C arrays
	for i := 0; i < len(vDevCUs); i++ {
		*((*C.int)(unsafe.Pointer(uintptr(unsafe.Pointer(cVdevCus)) + uintptr(i)*unsafe.Sizeof(*cVdevCus)))) = C.int(vDevCUs[i])
		*((*C.int)(unsafe.Pointer(uintptr(unsafe.Pointer(cVdevMemSize)) + uintptr(i)*unsafe.Sizeof(*cVdevMemSize)))) = C.int(vDevMemSize[i])
	}
	// Print first elements for verification
	fmt.Printf("cVdevCus[0]: %d, cVdevCus[1]: %d, cVdevMemSize[0]: %d, cVdevMemSize[1]: %d\n",
		*((*C.int)(unsafe.Pointer(cVdevCus))),
		*((*C.int)(unsafe.Pointer(uintptr(unsafe.Pointer(cVdevCus)) + uintptr(1)*unsafe.Sizeof(*cVdevCus)))),
		*((*C.int)(unsafe.Pointer(cVdevMemSize))),
		*((*C.int)(unsafe.Pointer(uintptr(unsafe.Pointer(cVdevMemSize)) + uintptr(1)*unsafe.Sizeof(*cVdevMemSize)))),
	)

	ret := C.dmiCreateVDevices(C.int(dvInd), C.int(vDevCount),
		cVdevCus, cVdevMemSize)
	glog.Infof("dmiCreateVDevices ret:%v ,err:%v", ret, dmiErrorString(ret))
	if err = dmiErrorString(ret); err != nil {
		return fmt.Errorf("Error dmiCreateVDevices:%s", err)
	}
	return
}

// 销毁指定物理设备上的所有虚拟设备
func dmiDestroyVDevices(dvInd int) (err error) {
	ret := C.dmiDestroyVDevices(C.int(dvInd))
	glog.Infof("dmiDestroyVDevices ret:%v", ret)
	if err = dmiErrorString(ret); err != nil {
		return fmt.Errorf("Error dmiDestroyVDevices:%s", err)
	}
	return
}

// 销毁指定虚拟设备
func dmiDestroySingleVDevice(vDvInd int) (err error) {
	ret := C.dmiDestroySingleVDevice(C.int(vDvInd))
	glog.Infof("dmiDestroySingleVDevice ret:%v", ret)
	if err = dmiErrorString(ret); err != nil {
		return fmt.Errorf("Error dmiDestroySingleVDevice:%s", err)
	}
	return
}

// 更新指定设备资源大小，vDevCUs和vDevMemSize为-1是不更改
func dmiUpdateSingleVDevice(vDvInd int, vDevCUs int, vDevMemSize int) (err error) {
	ret := C.dmiUpdateSingleVDevice(C.int(vDvInd), C.int(vDevCUs), C.int(vDevMemSize))
	glog.Infof("dmiUpdateSingleVDevice ret:%v, retstr:%v", ret, dmiErrorString(ret))
	if err = dmiErrorString(ret); err != nil {
		return fmt.Errorf("Error dmiUpdateSingleVDevice:%s", err)
	}
	return
}

// 启动虚拟设备
func dmiStartVDevice(vDvInd int) (err error) {
	ret := C.dmiStartVDevice(C.int(vDvInd))
	glog.Infof("StartVDevice ret:%v", ret)
	if err = dmiErrorString(ret); err != nil {
		return fmt.Errorf("Error dmiStartVDevice:%s", err)
	}
	return
}

// 停止虚拟设备
func dmiStopVDevice(vDvInd int) (err error) {
	ret := C.dmiStopVDevice(C.int(vDvInd))
	glog.Infof("dmiStopVDevice ret:%v,retmessage:%v", ret, dmiErrorString(ret))
	if err = dmiErrorString(ret); err != nil {
		return fmt.Errorf("Error dmiStopVDevice:%s", err)
	}
	return
}

// 返回物理设备使用百分比
func dmiGetDevBusyPercent(dvInd int) (percent int, err error) {
	var cpercent C.int
	ret := C.dmiGetDevBusyPercent(C.int(dvInd), &cpercent)
	if err = dmiErrorString(ret); err != nil {
		return percent, fmt.Errorf("Error dmiGetDevBusyPercent:%s", err)
	}
	percent = int(cpercent)
	return
}

// 返回虚拟设备使用百分比
func dmiGetVDevBusyPercent(vDvInd int) (percent int, err error) {
	var cpercent C.int
	ret := C.dmiGetVDevBusyPercent(C.int(vDvInd), &cpercent)
	if err = dmiErrorString(ret); err != nil {
		return percent, fmt.Errorf("Error dmiGetVDevBusyPercent:%s", err)
	}
	percent = int(cpercent)
	return
}

// 设置虚拟机加密状态 status为true，则开启加密虚拟机，否则关闭
func dmiSetEncryptionVMStatus(status bool) (err error) {
	ret := C.dmiSetEncryptionVMStatus(C.bool(status))
	if err = dmiErrorString(ret); err != nil {
		return fmt.Errorf("Error dmiSetEncryptionVMStatus:%s", err)
	}
	return
}

// 获取加密虚拟机状态
func dmiGetEncryptionVMStatus() (status bool, err error) {
	var cstatus C.bool
	ret := C.dmiGetEncryptionVMStatus(&cstatus)
	if err = dmiErrorString(ret); err != nil {
		return false, fmt.Errorf("Error dmiGetEncryptionVMStatus:%s", err)
	}
	status = bool(cstatus)
	glog.Infof("DmiGetEncryptionVMStatus: %v", status)
	return
}

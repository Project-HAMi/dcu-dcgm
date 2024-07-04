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
	"encoding/json"
	"fmt"
	"log"
	"unsafe"
)

type RSMIPcieBandwidth struct {
	TransferRate RSMIFrequencies
	lanes        [32]uint32
}

type RSMIFrequencies struct {
	NumSupported uint32
	Current      uint32
	Frequency    [32]uint64
}

type RSNIPowerProfilePresetMasks C.rsmi_power_profile_preset_masks_t

const (
	RSMI_PWR_PROF_PRST_CUSTOM_MASK       RSNIPowerProfilePresetMasks = C.RSMI_PWR_PROF_PRST_CUSTOM_MASK       //!< Custom Power Profile
	RSMI_PWR_PROF_PRST_VIDEO_MASK        RSNIPowerProfilePresetMasks = C.RSMI_PWR_PROF_PRST_VIDEO_MASK        //!< Video Power Profile
	RSMI_PWR_PROF_PRST_POWER_SAVING_MASK RSNIPowerProfilePresetMasks = C.RSMI_PWR_PROF_PRST_POWER_SAVING_MASK //!< Power Saving Profile
	RSMI_PWR_PROF_PRST_COMPUTE_MASK      RSNIPowerProfilePresetMasks = C.RSMI_PWR_PROF_PRST_COMPUTE_MASK      //!< Compute Saving Profile
	RSMI_PWR_PROF_PRST_VR_MASK           RSNIPowerProfilePresetMasks = C.RSMI_PWR_PROF_PRST_VR_MASK           //!< VR Power Profile

	//!< 3D Full Screen Power Profile
	RSMI_PWR_PROF_PRST_3D_FULL_SCR_MASK RSNIPowerProfilePresetMasks = C.RSMI_PWR_PROF_PRST_3D_FULL_SCR_MASK
	RSMI_PWR_PROF_PRST_BOOTUP_DEFAULT   RSNIPowerProfilePresetMasks = C.RSMI_PWR_PROF_PRST_BOOTUP_DEFAULT //!< Default Boot Up Profile
	RSMI_PWR_PROF_PRST_LAST             RSNIPowerProfilePresetMasks = RSMI_PWR_PROF_PRST_BOOTUP_DEFAULT

	//!< Invalid power profile
	RSMI_PWR_PROF_PRST_INVALID RSNIPowerProfilePresetMasks = C.RSMI_PWR_PROF_PRST_INVALID
)

type RSMIMemoryType C.rsmi_memory_type_t

const (
	RSMI_MEM_TYPE_FIRST    RSMIMemoryType = C.RSMI_MEM_TYPE_FIRST
	RSMI_MEM_TYPE_VRAM     RSMIMemoryType = C.RSMI_MEM_TYPE_VRAM
	RSMI_MEM_TYPE_VIS_VRAM RSMIMemoryType = C.RSMI_MEM_TYPE_VIS_VRAM
	RSMI_MEM_TYPE_GTT      RSMIMemoryType = C.RSMI_MEM_TYPE_GTT
	RSMI_MEM_TYPE_LAST     RSMIMemoryType = C.RSMI_MEM_TYPE_LAST
)

type RSMIRetiredPageRecord struct {
	PageAddress uint64               //!< Start address of page
	PageSize    uint64               //!< Page size
	Status      RSMIMemoryPageStatus //!< Page "reserved" status
}

type RSMIMemoryPageStatus C.rsmi_memory_page_status_t

const (
	RSMI_MEM_PAGE_STATUS_RESERVED     RSMIMemoryPageStatus = C.RSMI_MEM_PAGE_STATUS_RESERVED
	RSMI_MEM_PAGE_STATUS_PENDING      RSMIMemoryPageStatus = C.RSMI_MEM_PAGE_STATUS_PENDING
	RSMI_MEM_PAGE_STATUS_UNRESERVABLE RSMIMemoryPageStatus = C.RSMI_MEM_PAGE_STATUS_UNRESERVABLE
)

// go_rsmi_num_monitor_devices 获取gpu数量 *
func go_rsmi_num_monitor_devices() (gpuNum int, err error) {
	var p C.uint
	ret := C.rsmi_num_monitor_devices(&p)
	log.Println("go_rsmi_num_monitor_devices_ret:", ret)
	if err = errorString(ret); err != nil {
		return 0, fmt.Errorf("Error go_rsmi_num_monitor_devices_ret: %s", err)
	}
	gpuNum = int(p)
	log.Println("go_rsmi_num_monitor_devices:", gpuNum)
	return gpuNum, nil
}

// go_rsmi_dev_sku_get 获取设备sku
func go_rsmi_dev_sku_get(dvInd int) string {
	var sku C.char
	C.rsmi_dev_sku_get(C.uint32_t(dvInd), &sku)
	return string(sku)
}

// go_rsmi_dev_vendor_id_get 获取设备供应商id
func go_rsmi_dev_vendor_id_get(dvInd int) uint {
	var vid C.uint16_t
	C.rsmi_dev_vendor_id_get(C.uint32_t(dvInd), &vid)
	return uint(vid)
}

// go_rsmi_dev_id_get 获取设备id
func go_rsmi_dev_id_get(dvInd int) uint {
	var id C.uint16_t
	C.rsmi_dev_id_get(C.uint32_t(dvInd), &id)
	return uint(id)
}

// go_rsmi_dev_name_get 获取设备名称
func go_rsmi_dev_name_get(dvInd int) (nameStr string, err error) {
	name := make([]C.char, uint32(256))
	ret := C.rsmi_dev_name_get(C.uint32_t(dvInd), &name[0], 256)
	if err = errorString(ret); err != nil {
		return nameStr, fmt.Errorf("Error go_rsmi_dev_name_get: %s", err)
	}
	nameStr = C.GoString(&name[0])
	log.Println("go_rsmi_dev_name_get:", nameStr)
	return
}

// go_rsmi_dev_brand_get 获取设备品牌名称
func go_rsmi_dev_brand_get(dvInd int) string {
	brand := make([]C.char, uint32(256))
	C.rsmi_dev_brand_get(C.uint32_t(dvInd), &brand[0], 256)
	result := C.GoString(&brand[0])
	return result
}

// go_rsmi_dev_vendor_name_get 获取设备供应商名称
func go_rsmi_dev_vendor_name_get(dvInd int) string {
	bname := make([]C.char, uint32(256))
	C.rsmi_dev_vendor_name_get(C.uint32_t(dvInd), &bname[0], 80)
	result := C.GoString(&bname[0])
	return result
}

// go_rsmi_dev_vram_vendor_get 获取设备显存供应商名称
func go_rsmi_dev_vram_vendor_get(dvInd int) string {
	bname := make([]C.char, uint32(256))
	C.rsmi_dev_vram_vendor_get(C.uint32_t(dvInd), &bname[0], 80)
	result := C.GoString(&bname[0])
	return result
}

// go_rsmi_dev_serial_number_get 获取设备序列号 *
func go_rsmi_dev_serial_number_get(dvInd int) string {
	serialNumber := make([]C.char, uint32(256))
	C.rsmi_dev_serial_number_get(C.uint32_t(dvInd), &serialNumber[0], 256)
	result := C.GoString(&serialNumber[0])
	return result
}

// go_rsmi_dev_subsystem_id_get 获取设备子系统id
func go_rsmi_dev_subsystem_id_get(dvInd int) int {
	var id C.uint16_t
	C.rsmi_dev_subsystem_id_get(C.uint32_t(dvInd), &id)
	return int(id)
}

// go_rsmi_dev_subsystem_name_get 获取设备子系统名称 *
func go_rsmi_dev_subsystem_name_get(dvInd int) string {
	subSystemName := make([]C.char, uint32(256))
	C.rsmi_dev_subsystem_name_get(C.uint32_t(dvInd), &subSystemName[0], 256)
	result := C.GoString(&subSystemName[0])
	return result
}

// go_rsmi_dev_drm_render_minor_get 获取设备drm次编号
func go_rsmi_dev_drm_render_minor_get(dvInd int) int {
	var id C.uint32_t
	C.rsmi_dev_drm_render_minor_get(C.uint32_t(dvInd), &id)
	return int(id)
}

// go_rsmi_dev_unique_id_get 获取设备唯一id
func go_rsmi_dev_unique_id_get(dvInd int) int64 {
	var uniqueId C.uint64_t
	C.rsmi_dev_unique_id_get(C.uint32_t(dvInd), &uniqueId)
	return int64(uniqueId)
}

// go_rsmi_dev_subsystem_vendor_id_get 获取设备子系统供应商id
func go_rsmi_dev_subsystem_vendor_id_get(dvInd int) int {
	var id C.uint16_t
	C.rsmi_dev_subsystem_vendor_id_get(C.uint32_t(dvInd), &id)
	return int(id)
}

/****************************************** PCIe *********************************************/

// go_rsmi_dev_pci_bandwidth_get 获取可用的pcie带宽列表
func go_rsmi_dev_pci_bandwidth_get(dvInd int) RSMIPcieBandwidth {
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
	log.Println("RSMIPcieBandwidth:%s", rsmiPcieBandwidth)
	return rsmiPcieBandwidth
}

// go_rsmi_dev_pci_id_get 获取唯一pci设备标识符
func go_rsmi_dev_pci_id_get(dvInd int) int64 {
	var bdfid C.uint64_t
	C.rsmi_dev_pci_id_get(C.uint32_t(dvInd), &bdfid)
	return int64(bdfid)
}

// go_rsmi_topo_numa_affinity_get 获取与设备关联的numa节点
func go_rsmi_topo_numa_affinity_get(dvInd int) int {
	var namaNode C.uint32_t
	C.rsmi_topo_numa_affinity_get(C.uint32_t(dvInd), &namaNode)
	return int(namaNode)
}

// go_rsmi_dev_pci_throughput_get 获取pcie流量信息
func go_rsmi_dev_pci_throughput_get(dvInd int) (sent int64, received int64, maxPktSz int64) {
	var csent, creceived, cmaxpktsz C.uint64_t
	C.rsmi_dev_pci_throughput_get(C.uint32_t(dvInd), &csent, &creceived, &cmaxpktsz)
	sent = int64(cmaxpktsz)
	received = int64(csent)
	maxPktSz = int64(creceived)
	log.Printf("sent: %d, received: %d, maxPktSz: %d\n", sent, received, maxPktSz)
	return
}

// go_rsmi_dev_pci_replay_counter_get 获取pcie重放计数
func go_rsmi_dev_pci_replay_counter_get(dvInd int) uint64 {
	var counter C.uint64_t
	C.rsmi_dev_pci_replay_counter_get(C.uint32_t(dvInd), &counter)
	return uint64(counter)
}

// go_rsmi_dev_pci_bandwidth_set 设置可使用的pcie带宽集
func go_rsmi_dev_pci_bandwidth_set(dvInd int, bwBitmask int64) {
	C.rsmi_dev_pci_bandwidth_set(C.uint32_t(dvInd), C.uint64_t(bwBitmask))
}

/****************************************** Power *********************************************/

// go_rsmi_dev_power_ave_get 获取设备平均功耗
func go_rsmi_dev_power_ave_get(dvInd int, senserId int) int64 {
	var power C.uint64_t
	C.rsmi_dev_power_ave_get(C.uint32_t(dvInd), C.uint32_t(senserId), &power)
	return int64(power)
}

// go_rsmi_dev_energy_count_get 获取设备的能量累加计数
func go_rsmi_dev_energy_count_get() {

}

// go_rsmi_dev_power_cap_get 获取设备功率上限
func go_rsmi_dev_power_cap_get(dvInd int, senserId int) int64 {
	var power C.uint64_t
	C.rsmi_dev_power_cap_get(C.uint32_t(dvInd), C.uint32_t(senserId), &power)
	return int64(power)
}

// go_rsmi_dev_power_cap_range_get 获取设备功率有效值范围
func go_rsmi_dev_power_cap_range_get(dvInd int, senserId int) (max, min int64) {
	var cmax, cmin C.uint64_t
	C.rsmi_dev_power_cap_range_get(C.uint32_t(dvInd), C.uint32_t(senserId), &cmax, &cmin)
	max, min = int64(cmax), int64(cmin)
	return
}

// go_rsmi_dev_power_profile_set 设置设备功率配置文件
func go_rsmi_dev_power_profile_set(dvInd int, reserved int, profile RSNIPowerProfilePresetMasks) {
	C.rsmi_dev_power_profile_set(C.uint32_t(dvInd), C.uint32_t(reserved), C.rsmi_power_profile_preset_masks_t(profile))
}

/****************************************** Memory *********************************************/

// go_rsmi_dev_memory_total_get 获取设备内存总量 *
func go_rsmi_dev_memory_total_get(dvInd int, memoryType RSMIMemoryType) (total int64) {
	var ctotal C.uint64_t
	C.rsmi_dev_memory_total_get(C.uint32_t(dvInd), C.rsmi_memory_type_t(memoryType), &ctotal)
	total = int64(ctotal)
	log.Println("memory_total:", total)
	return
}

// go_rsmi_dev_memory_usage_get 获取当前设备内存使用情况 *
func go_rsmi_dev_memory_usage_get(dvInd int, memoryType RSMIMemoryType) (used int64) {
	var cused C.uint64_t
	C.rsmi_dev_memory_usage_get(C.uint32_t(dvInd), C.rsmi_memory_type_t(memoryType), &cused)
	used = int64(cused)
	log.Println("memory_used:", used)
	return
}

// go_rsmi_dev_memory_busy_percent_get 获取设备内存使用的百分比
func go_rsmi_dev_memory_busy_percent_get(dvInd int) int {
	var busyPercent C.uint32_t
	C.rsmi_dev_memory_busy_percent_get(C.uint32_t(dvInd), &busyPercent)
	return int(busyPercent)
}

// go_rsmi_dev_memory_reserved_pages_get
func go_rsmi_dev_memory_reserved_pages_get(dvInd int) (numPages int, records []RSMIRetiredPageRecord, err error) {
	var cnumPages C.uint32_t
	ret := C.rsmi_dev_memory_reserved_pages_get(C.uint32_t(dvInd), &cnumPages, nil)
	if ret != 0 {
		return 0, nil, fmt.Errorf("failed to get the number of pages, error code: %d", ret)
	}
	log.Println("cnumPages:", cnumPages)
	log.Println("cnumPages:", int(cnumPages))
	numPages = int(cnumPages)
	if numPages == 0 {
		return 0, nil, nil // No pages to retrieve
	}
	cRecords := make([]C.rsmi_retired_page_record_t, numPages)
	ret = C.rsmi_dev_memory_reserved_pages_get(C.uint32_t(dvInd), &cnumPages, (*C.rsmi_retired_page_record_t)(unsafe.Pointer(&records[0])))
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
	log.Printf("records:", indent)
	return
}

// go_rsmi_dev_fan_rpms_get 获取设备的风扇速度，实际转速
func go_rsmi_dev_fan_rpms_get(dvInd, sensorInd int) int64 {
	var speed C.int64_t
	C.rsmi_dev_fan_rpms_get(C.uint32_t(dvInd), C.uint32_t(sensorInd), &speed)
	return int64(speed)
}

// go_rsmi_dev_fan_speed_get 获取设备的风扇速度，相对速度值
func go_rsmi_dev_fan_speed_get(dvInd, sensorInd int) int64 {
	var speed C.int64_t
	C.rsmi_dev_fan_speed_get(C.uint32_t(dvInd), C.uint32_t(sensorInd), &speed)
	return int64(speed)
}

// go_rsmi_dev_fan_speed_max_get 获取设备的风扇速度，最大风速
func go_rsmi_dev_fan_speed_max_get(dvInd, sensorInd int) int64 {
	var maxSpeed C.uint64_t
	C.rsmi_dev_fan_speed_max_get(C.uint32_t(dvInd), C.uint32_t(sensorInd), &maxSpeed)
	return int64(maxSpeed)
}

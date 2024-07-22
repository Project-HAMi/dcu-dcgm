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
)

// rsmiComputeProcessInfoGet 获取当前使用GPU的所有进程信息
func rsmiComputeProcessInfoGet() (processInfo []RSMIProcessInfo, numItems int, err error) {
	var cprocessInfo C.rsmi_process_info_t
	var cnumItems C.uint32_t
	ret := C.rsmi_compute_process_info_get(&cprocessInfo, &cnumItems)
	if err := errorString(ret); err != nil {
		return processInfo, numItems, fmt.Errorf("Error rsmiComputeProcessInfoGet: %s", err)
	}
	// 创建一个大小为numItems的切片
	processInfo = make([]RSMIProcessInfo, int(cnumItems))
	// 第二次调用以获取实际的数据
	ret = C.rsmi_compute_process_info_get((*C.rsmi_process_info_t)(unsafe.Pointer(&processInfo[0])), &cnumItems)
	if err := errorString(ret); err != nil {
		return processInfo, numItems, fmt.Errorf("Error in go_rsmi_compute_process_info_get: %s", err)
	}
	numItems = int(cnumItems)
	return
}

// rsmiComputeProcessInfoByPidGet 获取指定进程的进程信息
func rsmiComputeProcessInfoByPidGet(pid int) (proc RSMIProcessInfo, err error) {
	var cproc C.rsmi_process_info_t
	ret := C.rsmi_compute_process_info_by_pid_get(C.uint32_t(pid), &cproc)
	if err = errorString(ret); err != nil {
		return proc, fmt.Errorf("Error rsmiComputeProcessInfoByPidGet:%s", err)
	}
	proc = RSMIProcessInfo{
		ProcessID:   uint32(cproc.process_id),
		Pasid:       uint32(cproc.pasid),
		VramUsage:   uint64(cproc.vram_usage),
		SdmaUsage:   uint64(cproc.sdma_usage),
		CuOccupancy: uint32(cproc.cu_occupancy),
	}
	return
}

// rsmiComputeProcessGpusGet 获取进程当前正在使用的设备索引
func rsmiComputeProcessGpusGet(pid int) (dvIndices []int, err error) {
	var cnumDevices C.uint32_t
	// 第一次调用以获取numDevices的值
	ret := C.rsmi_compute_process_gpus_get(C.uint32_t(pid), nil, &cnumDevices)
	if err := errorString(ret); err != nil {
		return dvIndices, fmt.Errorf("Error in RSMIComputeProcessGPUsGet (initial call): %s", err)
	}

	// 创建一个大小为numDevices的切片
	dvIndicesC := make([]C.uint32_t, cnumDevices)
	// 第二次调用以获取实际的数据
	ret = C.rsmi_compute_process_gpus_get(C.uint32_t(pid), &dvIndicesC[0], &cnumDevices)
	if err := errorString(ret); err != nil {
		return nil, fmt.Errorf("Error in RSMIComputeProcessGPUsGet: %s", err)
	}
	// 将C数组转换为Go切片
	dvIndices = make([]int, cnumDevices)
	for i := 0; i < int(cnumDevices); i++ {
		dvIndices[i] = int(dvIndicesC[i])
	}
	return
}

// rsmiDevXGMIErrorStatus 获取设备的XGMI错误状态
func rsmiDevXGMIErrorStatus(dvInd int) (status RSMIXGMIStatus, err error) {
	var cStatus C.rsmi_xgmi_status_t

	ret := C.rsmi_dev_xgmi_error_status(C.uint32_t(dvInd), &cStatus)
	if err := errorString(ret); err != nil {
		return status, fmt.Errorf("Error RSMIDevXGMIErrorStatus: %s", err)
	}
	status = RSMIXGMIStatus(cStatus)
	return
}

// rsmiDevXgmiErrorReset 重置设备的XGMI错误状态
func rsmiDevXgmiErrorReset(dvInd int) (err error) {
	ret := C.rsmi_dev_xgmi_error_reset(C.uint32_t(dvInd))
	if err = errorString(ret); err != nil {
		return fmt.Errorf("Error rsmiDevXgmiErrorReset:%s", err)
	}
	return
}

// rsmiDevXgmiHiveIdGet 获取设备的XGMI hive id
func rsmiDevXgmiHiveIdGet(dvInd int) (hiveId int64, err error) {
	var chiveId C.uint64_t
	ret := C.rsmi_dev_xgmi_hive_id_get(C.uint32_t(dvInd), &chiveId)
	if err = errorString(ret); err != nil {
		return hiveId, fmt.Errorf("Error rsmiDevXgmiHiveIdGet:%s", err)
	}
	hiveId = int64(chiveId)
	return
}

// rsmiTopoGetNumaBodeBumber 获取设备的numa cpu节点号
func rsmiTopoGetNumaBodeBumber(dvInd int) (numaNode int, err error) {
	var cnumaNode C.uint32_t
	ret := C.rsmi_topo_get_numa_node_number(C.uint32_t(dvInd), &cnumaNode)
	if err = errorString(ret); err != nil {
		return numaNode, fmt.Errorf("Error rsmiTopoGetNumaBodeBumber:%s", err)
	}
	numaNode = int(cnumaNode)
	return
}

// rsmiTopoGetLinkWeight 获取2个gpu之间连接的权重
func rsmiTopoGetLinkWeight(dvIndSrc, dvIndDst int) (weight int64, err error) {
	var cweight C.uint64_t
	ret := C.rsmi_topo_get_link_weight(C.uint32_t(dvIndSrc), C.uint32_t(dvIndDst), &cweight)
	if err = errorString(ret); err != nil {
		return weight, fmt.Errorf("Error rsmiTopoGetLinkWeight:%S", err)
	}
	weight = int64(cweight)
	return
}

// rsmiTopoGetLinkType 获取2个gpu之间的hops和连接类型
func rsmiTopoGetLinkType(dvIndSrc, dvIndDst int) (hops int64, linkType RSMIIOLinkType, err error) {
	var chops C.uint64_t
	var clinkType C.RSMI_IO_LINK_TYPE
	ret := C.rsmi_topo_get_link_type(C.uint32_t(dvIndSrc), C.uint32_t(dvIndDst), &chops, &clinkType)
	if err = errorString(ret); err != nil {
		return hops, linkType, fmt.Errorf("Error rsmiTopoGetLinkType:%s", err)
	}
	hops = int64(chops)
	linkType = RSMIIOLinkType(clinkType)
	return
}

// rsmiDevSupportedFuncIteratorOpen 获取设备支持RSMI函数的函数名迭代器
func rsmiDevSupportedFuncIteratorOpen(dvInd int) (iterHandle RSMIFuncIDIterHandle, err error) {
	var handle C.rsmi_func_id_iter_handle_t
	ret := C.rsmi_dev_supported_func_iterator_open(C.uint32_t(dvInd), &handle)
	if err = errorString(ret); err != nil {
		return iterHandle, fmt.Errorf("Error rsmiDevSupportedFuncIteratorOpen: %s", err)
	}
	iterHandle = RSMIFuncIDIterHandle(handle)
	return
}

// rsmiDevSupportedVariantIteratorOpen 获取给定句柄的变体迭代器
func rsmiDevSupportedVariantIteratorOpen(iterHandle RSMIFuncIDIterHandle) (handle RSMIFuncIDIterHandle, err error) {
	var chandle C.rsmi_func_id_iter_handle_t
	ret := C.rsmi_dev_supported_variant_iterator_open(C.rsmi_func_id_iter_handle_t(iterHandle), &chandle)
	if err = errorString(ret); err != nil {
		return iterHandle, fmt.Errorf("Error rsmiDevSupportedVariantIteratorOpen: %s", err)
	}
	handle = RSMIFuncIDIterHandle(chandle)
	return
}

// rsmiFuncIterNext 推进函数标识符迭代器
func rsmiFuncIterNext(handle RSMIFuncIDIterHandle) (err error) {
	ret := C.rsmi_func_iter_next(C.rsmi_func_id_iter_handle_t(handle))
	if err = errorString(ret); err != nil {
		return fmt.Errorf("Error rsmiFuncIterNext:%s", err)
	}
	return
}

// rsmiDevSupportedFuncIteratorClose 关闭变量迭代器句柄
func rsmiDevSupportedFuncIteratorClose(handle RSMIFuncIDIterHandle) (err error) {
	cHandle := C.rsmi_func_id_iter_handle_t(handle)
	ret := C.rsmi_dev_supported_func_iterator_close(&cHandle)
	if err = errorString(ret); err != nil {
		return fmt.Errorf("Error rsmiDevSupportedFuncIteratorClose:%s", err)
	}
	return
}

// rsmiFuncIterValueGet 获取与函数/变量迭代器相关联的值
//func rsmiFuncIterValueGet(handle RSMIFuncIDIterHandle) (value RSMIFuncIDValue, err error) {
//	var cvalue C.rsmi_func_id_value_t
//	// 调用C函数
//	ret := C.rsmi_func_iter_value_get(C.rsmi_func_id_iter_handle_t(handle), &cvalue)
//	if err = errorString(ret); err != nil {
//		return value, fmt.Errorf("Error rsmiFuncIterValueGet:%s", err)
//	}
//	value.ID = uint64(cvalue.id)
//	value.Name = C.GoString((*C.char)(unsafe.Pointer(cvalue.name)))
//	value.MemoryType = RSMIMemoryType(*(*C.rsmi_memory_type_t)(unsafe.Pointer(&cvalue)))
//	value.TempMetric = RSMITemperatureMetric(*(*C.rsmi_temperature_metric_t)(unsafe.Pointer(&cvalue)))
//	value.EventType = RSMIEventType(*(*C.rsmi_event_type_t)(unsafe.Pointer(&cvalue)))
//	value.EventGroup = RSMIEventGroup(*(*C.rsmi_event_group_t)(unsafe.Pointer(&cvalue)))
//	value.ClkType = RSMIClkType(*(*C.rsmi_clk_type_t)(unsafe.Pointer(&cvalue)))
//	value.FwBlock = RSMIFwBlock(*(*C.rsmi_fw_block_t)(unsafe.Pointer(&cvalue)))
//	value.GpuBlock = RSMIGpuBlock(*(*C.rsmi_gpu_block_t)(unsafe.Pointer(&cvalue)))
//	return
//}

// rsmiEventNotificationInit 准备收集GPU事件通知
func rsmiEventNotificationInit(deInd int) (err error) {
	ret := C.rsmi_event_notification_init(C.uint32_t(deInd))
	if err = errorString(ret); err != nil {
		return fmt.Errorf("Rrror rsmiEventNotificationInit:%s", err)
	}
	return
}

// rsmiEventNotificationMaskSet 设置设备指定要收集的事件。
func rsmiEventNotificationMaskSet(dvInd int, mask int64) (err error) {
	ret := C.rsmi_event_notification_mask_set(C.uint32_t(dvInd), C.uint64_t(mask))
	if err = errorString(ret); err != nil {
		return fmt.Errorf("Rrror rsmiEventNotificationMaskSet:%s", err)
	}
	return
}

// rsmiEventNotificationGet 收集事件通知，等待指定时间
func rsmiEventNotificationGet(timeoutMs int) (numElem int, datas []RSMIEEvtNotificationData, err error) {
	var cnumElen C.uint32_t
	ret := C.rsmi_event_notification_get(C.int(timeoutMs), &cnumElen, nil)
	if err = errorString(ret); err != nil {
		return 0, nil, fmt.Errorf("Error rsmiEventNotificationGet,numElem:%s", err)
	}
	numElem = int(cnumElen)
	cdatas := make([]C.rsmi_evt_notification_data_t, numElem)
	ret = C.rsmi_event_notification_get(C.int(timeoutMs), &cnumElen, (*C.rsmi_evt_notification_data_t)(unsafe.Pointer(&cdatas[0])))
	if err = errorString(ret); err != nil {
		return numElem, nil, fmt.Errorf("Error rsmiEventNotificationGet,datas:%s", err)
	}
	datas = make([]RSMIEEvtNotificationData, numElem)
	for i, data := range cdatas {
		datas[i] = RSMIEEvtNotificationData{
			DvInd:   uint32(data.dv_ind),
			Event:   RSMIEvtNotificationType(data.event),
			Message: *(*[64]byte)(unsafe.Pointer(&data.message)),
		}
	}
	return
}

// rsmiEventNotificationStop 关闭任何文件句柄并释放由GPU事件通知使用的任何资源。
func rsmiEventNotificationStop(dvInd int) (err error) {
	ret := C.rsmi_event_notification_stop(C.uint32_t(dvInd))
	if err = errorString(ret); err != nil {
		return fmt.Errorf("Error rsmiEventNotificationStop:%s", err)
	}
	return
}

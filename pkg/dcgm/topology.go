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
import "fmt"

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

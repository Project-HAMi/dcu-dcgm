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
	"unsafe"
)

func errorString(result C.rsmi_status_t) error {
	if RSMIStatus(result) == RSMI_STATUS_SUCCESS {
		return nil
	}
	var cStatusString *C.char
	statusCode := C.rsmi_status_string(result, (**C.char)(unsafe.Pointer(&cStatusString)))
	if RSMIStatus(statusCode) != RSMI_STATUS_SUCCESS {
		return fmt.Errorf("error: %s", statusCode)
	}
	goStatusString := C.GoString(cStatusString)
	return fmt.Errorf("%s", goStatusString)
}
func dataToJson(data any) string {
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error serializing to JSON:", err)
	}
	return string(jsonData)
}

// 获取所提供的RSMI错误状态的描述
func go_rsmi_status_string(status RSMIStatus) (statusStr string, err error) {
	var cstatusStr *C.char
	ret := C.rsmi_status_string(C.rsmi_status_t(status), (**C.char)(unsafe.Pointer(&cstatusStr)))
	if err = errorString(ret); err != nil {
		return statusStr, fmt.Errorf("Error go_rsmi_status_string:%s", err)
	}
	statusStr = C.GoString(cstatusStr)
	return
}

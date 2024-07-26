package dcgm

/*
#cgo CFLAGS: -Wall -I./include
#cgo LDFLAGS: -L./lib -lrocm_smi64 -lhydmi -Wl,--unresolved-symbols=ignore-in-object-files
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

func dmiErrorString(result C.dmiStatus) error {
	if DMIStatus(result) == DMI_STATUS_SUCCESS {
		return nil
	}
	var cStatusString *C.char
	statusCode := C.dmiGetStatusString(result, (**C.char)(unsafe.Pointer(&cStatusString)))
	if DMIStatus(statusCode) != DMI_STATUS_SUCCESS {
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
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func indexOf(slice []string, item string) int {
	for i, s := range slice {
		if s == item {
			return i
		}
	}
	return -1
}

func perfLevelString(i int) string {
	switch i {
	case 0:
		return "AUTO"
	case 1:
		return "LOW"
	case 2:
		return "HIGH"
	case 3:
		return "MANUAL"
	case 4:
		return "STABLE_STD"
	case 5:
		return "STABLE_PEAK"
	case 6:
		return "STABLE_MIN_MCLK"
	case 7:
		return "STABLE_MIN_SCLK"
	default:
		return "UNKNOWN"
	}
}

func ConvertASCIIToString(asciiCodes []byte) string {
	var result []rune
	for _, code := range asciiCodes {
		// Stop at the first null character
		if code == 0 {
			break
		}
		// Filter out non-ASCII characters
		if code > 127 {
			continue
		}
		result = append(result, rune(code))
	}
	return string(result)
}

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
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
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

// parseConfig 解析配置文件内容为DMIVDeviceInfo结构体
func parseConfig(filePath string) (*DMIVDeviceInfo, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &DMIVDeviceInfo{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ": ", 2)
		if len(parts) < 2 {
			continue
		}
		key := parts[0]
		value := parts[1]

		switch key {
		case "cu_count":
			config.ComputeUnitCount, _ = strconv.Atoi(value)
		case "mem":
			// 解析内存大小，例如 "4096 MiB"
			memParts := strings.Fields(value)
			if len(memParts) == 2 {
				memSize, err := strconv.Atoi(memParts[0])
				if err == nil {
					// 转换为字节数（假设单位是 MiB）
					config.GlobalMemSize = uintptr(memSize * 1024 * 1024)
				}
			}
		case "device_id":
			config.DeviceID, _ = strconv.Atoi(value)
		case "vdev_id":
			config.VMinorNumber, _ = strconv.Atoi(value)
		case "PciBusId":
			config.PicBusNumber = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return config, nil
}

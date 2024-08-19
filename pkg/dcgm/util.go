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
	jsonData, err := json.MarshalIndent(data, "", "  ")
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
			config.PciBusNumber = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return config, nil
}

// 打印二维数组
func print2DArray(data [][]string) {
	for _, row := range data {
		fmt.Println(strings.Join(row, "\t"))
	}
}

// 打印超出规格运行的警告，并提示用户接受条款
func confirmOutOfSpecWarning(autoRespond bool) {
	warning := `
          ******WARNING******

          Operating your AMD GPU outside of official AMD specifications or outside of
          factory settings, including but not limited to the conducting of overclocking,
          over-volting or under-volting (including use of this interface software,
          even if such software has been directly or indirectly provided by AMD or otherwise
          affiliated in any way with AMD), may cause damage to your AMD GPU, system components
          and/or result in system failure, as well as cause other problems.
          DAMAGES CAUSED BY USE OF YOUR AMD GPU OUTSIDE OF OFFICIAL AMD SPECIFICATIONS OR
          OUTSIDE OF FACTORY SETTINGS ARE NOT COVERED UNDER ANY AMD PRODUCT WARRANTY AND
          MAY NOT BE COVERED BY YOUR BOARD OR SYSTEM MANUFACTURER'S WARRANTY.
          Please use this utility with caution.
          `

	fmt.Println(warning)

	var userInput string
	if !autoRespond {
		fmt.Print("Do you accept these terms? [y/N] ")
		fmt.Scanln(&userInput)
	} else {
		userInput = "y"
	}

	userInput = strings.ToLower(userInput)
	if userInput == "yes" || userInput == "y" {
		return
	} else {
		fmt.Println("Confirmation not given. Exiting without setting value")
		os.Exit(1)
	}
}
func profileString(profile interface{}) string {
	dictionary := map[int]string{
		1:  "CUSTOM",
		2:  "VIDEO",
		4:  "POWER SAVING",
		8:  "COMPUTE",
		16: "VR",
		32: "3D FULL SCREEN",
		64: "BOOTUP DEFAULT",
	}

	switch v := profile.(type) {
	case int:
		if name, ok := dictionary[v]; ok {
			return name
		}
	case string:
		if num, err := strconv.Atoi(v); err == nil {
			if name, ok := dictionary[num]; ok {
				return name
			}
		} else {
			for key, val := range dictionary {
				if val == v {
					return strconv.Itoa(key)
				}
			}
		}
	}
	return "UNKNOWN"
}

func profileEnum(profile string) RSNIPowerProfilePresetMasks {
	dictionary := map[string]RSNIPowerProfilePresetMasks{
		"CUSTOM":         RSMI_PWR_PROF_PRST_CUSTOM_MASK,
		"VIDEO":          RSMI_PWR_PROF_PRST_VIDEO_MASK,
		"POWER SAVING":   RSMI_PWR_PROF_PRST_POWER_SAVING_MASK,
		"COMPUTE":        RSMI_PWR_PROF_PRST_COMPUTE_MASK,
		"VR":             RSMI_PWR_PROF_PRST_VR_MASK,
		"3D FULL SCREEN": RSMI_PWR_PROF_PRST_3D_FULL_SCR_MASK,
		"BOOTUP DEFAULT": RSMI_PWR_PROF_PRST_BOOTUP_DEFAULT,
	}

	if val, ok := dictionary[profile]; ok {
		return val
	}
	return RSMI_PWR_PROF_PRST_INVALID
}

func printTableLog(headers []string, data [][]string, device int, title string) {
	fmt.Printf("Device: %d - %s\n", device, title)
	fmt.Println(headers)
	for _, row := range data {
		fmt.Println(row)
	}
	fmt.Println()
}

func formatMatrixToJSON(deviceList []int, matrix [][]int64, metricName string) {
	for rowIndx := 0; rowIndx < len(deviceList); rowIndx++ {
		for colInd := rowIndx + 1; colInd < len(deviceList); colInd++ {
			valueStr := matrix[deviceList[rowIndx]][deviceList[colInd]]
			fmt.Printf(metricName+"\n", deviceList[rowIndx], deviceList[colInd])
			fmt.Println(valueStr)
		}
	}
}

func formatMatrixToStrJSON(deviceList []int, matrix [][]string, metricName string) {
	for rowIndx := 0; rowIndx < len(deviceList); rowIndx++ {
		for colInd := rowIndx + 1; colInd < len(deviceList); colInd++ {
			valueStr := matrix[deviceList[rowIndx]][deviceList[colInd]]
			fmt.Printf(metricName+"\n", deviceList[rowIndx], deviceList[colInd])
			fmt.Println(valueStr)
		}
	}
}

func printTableRow(format string, displayString interface{}) {
	if format != "" {
		fmt.Printf(format, displayString)
	} else {
		fmt.Print(displayString)
	}
	fmt.Print(" ")
}

// 获取指定目录下的文件列表，如果目录不存在或为空，返回空切片
func getConfigFiles(dir string) ([]os.DirEntry, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			// 如果目录不存在，返回空切片
			return []os.DirEntry{}, nil
		}
		return nil, err
	}
	return files, nil
}

// 解析配置文件内容
func parseConfigFile(filePath string) (map[string]string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(content), "\n")
	config := make(map[string]string)
	for _, line := range lines {
		parts := strings.Split(line, ":")
		if len(parts) == 2 {
			config[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return config, nil
}

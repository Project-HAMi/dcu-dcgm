package router

import (
	"fmt"

	"g.sugon.com/das/dcgm-dcu/pkg/dcgm"
)

// 将字符串转换为 RSMIDevPerfLevel 类型
func ConvertToRSMIDevPerfLevel(level string) (dcgm.RSMIDevPerfLevel, error) {
	switch level {
	case "AUTO":
		return dcgm.RSMI_DEV_PERF_LEVEL_AUTO, nil
	case "FIRST":
		return dcgm.RSMI_DEV_PERF_LEVEL_FIRST, nil
	case "LOW":
		return dcgm.RSMI_DEV_PERF_LEVEL_LOW, nil
	case "HIGH":
		return dcgm.RSMI_DEV_PERF_LEVEL_HIGH, nil
	case "MANUAL":
		return dcgm.RSMI_DEV_PERF_LEVEL_MANUAL, nil
	case "STABLE_STD":
		return dcgm.RSMI_DEV_PERF_LEVEL_STABLE_STD, nil
	case "STABLE_PEAK":
		return dcgm.RSMI_DEV_PERF_LEVEL_STABLE_PEAK, nil
	case "STABLE_MIN_MCLK":
		return dcgm.RSMI_DEV_PERF_LEVEL_STABLE_MIN_MCLK, nil
	case "STABLE_MIN_SCLK":
		return dcgm.RSMI_DEV_PERF_LEVEL_STABLE_MIN_SCLK, nil
	case "DETERMINISM":
		return dcgm.RSMI_DEV_PERF_LEVEL_DETERMINISM, nil
	case "LAST":
		return dcgm.RSMI_DEV_PERF_LEVEL_LAST, nil
	case "UNKNOWN":
		return dcgm.RSMI_DEV_PERF_LEVEL_UNKNOWN, nil
	default:
		return dcgm.RSMI_DEV_PERF_LEVEL_UNKNOWN, fmt.Errorf("invalid level string: %s", level)
	}
}

// ConvertToRSMIGpuBlock 函数定义
func ConvertToRSMIGpuBlock(block string) (dcgm.RSMIGpuBlock, error) {
	switch block {
	case "INVALID":
		return dcgm.RSMIGpuBlockInvalid, nil
	case "FIRST":
		return dcgm.RSMIGpuBlockFirst, nil
	case "UMC":
		return dcgm.RSMIGpuBlockUMC, nil
	case "SDMA":
		return dcgm.RSMIGpuBlockSDMA, nil
	case "GFX":
		return dcgm.RSMIGpuBlockGFX, nil
	case "MMHUB":
		return dcgm.RSMIGpuBlockMMHUB, nil
	case "ATHUB":
		return dcgm.RSMIGpuBlockATHUB, nil
	case "PCIEBIF":
		return dcgm.RSMIGpuBlockPCIEBIF, nil
	case "HDP":
		return dcgm.RSMIGpuBlockHDP, nil
	case "XGMIWAFL":
		return dcgm.RSMIGpuBlockXGMIWAFL, nil
	case "DF":
		return dcgm.RSMIGpuBlockDF, nil
	case "SMN":
		return dcgm.RSMIGpuBlockSMN, nil
	case "SEM":
		return dcgm.RSMIGpuBlockSEM, nil
	case "MP0":
		return dcgm.RSMIGpuBlockMP0, nil
	case "MP1":
		return dcgm.RSMIGpuBlockMP1, nil
	case "FUSE":
		return dcgm.RSMIGpuBlockFuse, nil
	case "MCA":
		return dcgm.RSMIGpuBlockMCA, nil
	case "LAST":
		return dcgm.RSMIGpuBlockLast, nil
	case "RESERVED":
		return dcgm.RSMIGpuBlockReserved, nil
	default:
		return dcgm.RSMIGpuBlockInvalid, fmt.Errorf("invalid block string: %s", block)
	}
}

// ConvertToRSMISwComponent 函数定义
func ConvertToRSMISwComponent(component string) (dcgm.RSMISwComponent, error) {
	switch component {
	case "FIRST":
		return dcgm.RSMISwCompFirst, nil
	case "DRIVER":
		return dcgm.RSMISwCompDriver, nil
	case "LAST":
		return dcgm.RSMISwCompLast, nil
	default:
		return dcgm.RSMISwCompFirst, fmt.Errorf("invalid component string: %s", component)
	}
}

// Response represents a basic structure for API responses.
type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// SuccessResponse creates a successful response payload with custom data.
func SuccessResponse(data map[string]interface{}) Response {
	return Response{
		Message: "成功",
		Data:    data,
	}
}

// ErrorResponse creates an error response payload.
func ErrorResponse(data interface{}) Response {
	return Response{
		Message: "失败",
		Data:    data,
	}
}

// ConvertFrequencyToSclkClock 将sclk频率值转换为对应的十进制值
func ConvertFrequencyToSclkClock(freq string) (int64, error) {
	switch freq {
	case "600":
		return 1, nil
	case "700":
		return 2, nil
	case "750":
		return 4, nil
	case "800":
		return 8, nil
	case "900":
		return 16, nil
	case "1000":
		return 32, nil
	case "1106":
		return 64, nil
	case "1200":
		return 128, nil
	case "1270":
		return 256, nil
	case "1319":
		return 512, nil
	case "1400":
		return 1024, nil
	case "1500":
		return 2048, nil
	default:
		return 0, fmt.Errorf("invalid frequency value: %s", freq)
	}
}

// ConvertFrequencyToSocclkClock 将socclk频率值转换为对应的十进制值
func ConvertFrequencyToSocclkClock(freq string) (int64, error) {
	switch freq {
	case "309":
		return 1, nil
	case "523":
		return 2, nil
	case "566":
		return 4, nil
	case "618":
		return 8, nil
	case "680":
		return 16, nil
	case "755":
		return 32, nil
	case "850":
		return 64, nil
	case "971":
		return 128, nil
	default:
		return 0, fmt.Errorf("invalid frequency value: %s", freq)
	}
}

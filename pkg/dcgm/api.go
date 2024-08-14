package dcgm

import "C"
import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang/glog"
)

// @Summary 初始化 DCGM
// @Description 初始化 (DCGM) 库。
// @Produce json
// @Success 200 {object} string "成功初始化"
// @Failure 500 {object} error "初始化失败"
// @Router /Init [post]
func Init() error {
	return rsmiInit()
}

// @Summary 关闭 DCGM
// @Description 关闭 Data Center GPU Manager (DCGM) 库。
// @Produce json
// @Success 200 {object} string "成功关闭"
// @Failure 500 {object} error "关闭失败"
// @Router /ShutDown [post]
func ShutDown() error {
	return rsmiShutdown()
}

// @Summary 获取 GPU 数量
// @Description 获取监视的 GPU 数量。
// @Produce json
// @Success 200 {int} int "GPU 数量"
// @Failure 500 {object} error "获取 GPU 数量失败"
// @Router /NumMonitorDevices [get]
func NumMonitorDevices() (int, error) {
	return rsmiNumMonitorDevices()
}

// 获取设备利用率计数器
// @Summary 获取设备利用率计数器
// @Description 根据设备索引获取利用率计数器
// @Param dvInd query int true "设备索引"
// @Param utilizationCounters body []RSMIUtilizationCounter true "利用率计数器对象列表"
// @Param count query int true "计数器的数量"
// @Success 200 {object} int64 "返回的时间戳"
// @Failure 400 {object} error "请求失败"
// @Router /utilizationcount [post]
func UtilizationCount(dvInd int, utilizationCounters []RSMIUtilizationCounter, count int) (timestamp int64, err error) {
	return rsmiUtilizationCountGet(dvInd, utilizationCounters, count)
}

// @Summary 获取设备名称
// @Description 根据设备 ID 获取设备名称。
// @Produce json
// @Param dvInd path int true "设备 ID"
// @Success 200 {string} name "设备名称"
// @Failure 400 {object} error "请求错误"
// @Failure 404 {object} error "设备未找到"
// @Router /DevName [get]
func DevName(dvInd int) (name string, err error) {
	return rsmiDevNameGet(dvInd)
}

// 获取设备SKU
// @Summary 获取设备SKU
// @Description 根据设备索引获取SKU
// @Param dvInd query int true "设备索引"
// @Success 200 {int} sku "返回设备SKU"
// @Failure 400 {object} error "请求失败"
// @Router /DevSku [get]
func DevSku(dvInd int) (sku int, err error) {
	return rsmiDevSkuGet(dvInd)
}

// 获取设备品牌名称
// @Summary 获取设备品牌名称
// @Description 根据设备索引获取品牌名称
// @Param dvInd query int true "设备索引"
// @Success 200 {string} brand "返回设备品牌名称"
// @Failure 400 {object} error "请求失败"
// @Router /DevBrand [get]
func DevBrand(dvInd int) (brand string, err error) {
	return rsmiDevBrandGet(dvInd)
}

// 获取设备供应商名称
// @Summary 获取设备供应商名称
// @Description 根据设备索引获取供应商名称
// @Param dvInd query int true "设备索引"
// @Success 200 {string} bname "返回设备供应商名称"
// @Failure 400 {object} error "请求失败"
// @Router /DevVendorName [get]
func DevVendorName(dvInd int) (bname string, err error) {
	return rsmiDevVendorNameGet(dvInd)
}

// 获取设备显存供应商名称
// @Summary 获取设备显存供应商名称
// @Description 根据设备索引获取显存供应商名称
// @Param dvInd query int true "设备索引"
// @Success 200 {string} name "返回显存供应商名称"
// @Failure 400 {object} error "请求失败"
// @Router /DevVramVendor [get]
func DevVramVendor(dvInd int) (name string, err error) {
	return rsmiDevVramVendorGet(dvInd)
}

// @Summary 获取可用的 PCIe 带宽列表
// @Description 根据设备 ID 获取设备的可用 PCIe 带宽列表。
// @Produce json
// @Param dvInd path int true "设备 ID"
// @Success 200 {RSMIPcieBandwidth} rsmiPcieBandwidth "PCIe 带宽列表"
// @Failure 400 {object} error "请求错误"
// @Failure 404 {object} error "设备未找到"
// @Router /DevPciBandwidth [get]
func DevPciBandwidth(dvInd int) (rsmiPcieBandwidth RSMIPcieBandwidth, err error) {
	return rsmiDevPciBandwidthGet(dvInd)

}

// @Summary 获取内存使用百分比
// @Description 根据设备 ID 获取设备内存的使用百分比。
// @Produce json
// @Param dvInd path int true "设备 ID"
// @Success 200 {int} busyPercent "内存使用百分比"
// @Failure 400 {object} error "请求错误"
// @Failure 404 {object} error "设备未找到"
// @Router /MemoryPercent [get]
func MemoryPercent(dvInd int) (busyPercent int, err error) {
	return rsmiDevMemoryBusyPercentGet(dvInd)
}

// 获取设备温度值
//func DevTemp(dvInd int) int64 {
//	return go_rsmi_dev_temp_metric_get(dvInd)
//}

// @Summary 设置设备 PowerPlay 性能级别
// @Description 根据设备 ID 设置 PowerPlay 性能级别。
// @Produce json
// @Param dvInd path int true "设备 ID"
// @Param level query string true "要设置的性能级别"
// @Success 200 {string} string "操作成功"
// @Failure 400 {object} error "请求错误"
// @Failure 404 {object} error "设备未找到"
// @Router /DevPerfLevelSet [post]
func DevPerfLevelSet(dvInd int, level RSMIDevPerfLevel) error {
	return rsmiDevPerfLevelSet(dvInd, level)
}

// @Summary 获取 GPU 度量信息
// @Description 根据设备 ID 获取 GPU 的度量信息。
// @Produce json
// @Param dvInd query int true "设备 ID"
// @Success 200 {object} RSMIGPUMetrics "GPU 度量信息"
// @Failure 400 {object} error "请求错误"
// @Failure 404 {object} error "设备未找到"
// @Router /DevGpuMetricsInfo [get]
func DevGpuMetricsInfo(dvInd int) (gpuMetrics RSMIGPUMetrics, err error) {
	return rsmiDevGpuMetricsInfoGet(dvInd)
}

// @Summary 获取设备监控中的指标
// @Description 收集所有设备的监控指标信息。
// @Produce json
// @Success 200 {array} MonitorInfo "设备监控指标信息列表"
// @Failure 400 {object} error "请求错误"
// @Failure 404 {object} error "设备未找到"
// @Router /CollectDeviceMetrics [get]
func CollectDeviceMetrics() (monitorInfos []MonitorInfo, err error) {
	numMonitorDevices, err := rsmiNumMonitorDevices()
	if err != nil {
		return nil, err
	}
	for i := 0; i < numMonitorDevices; i++ {
		bdfid, err := rsmiDevPciIdGet(i)
		if err != nil {
			return nil, err
		}
		// 解析BDFID
		domain := (bdfid >> 32) & 0xffffffff
		bus := (bdfid >> 8) & 0xff
		dev := (bdfid >> 3) & 0x1f
		function := bdfid & 0x7
		// 格式化PCI ID
		pciBusNumber := fmt.Sprintf("%04x:%02x:%02x.%x", domain, bus, dev, function)
		//设备序列号
		deviceId, _ := rsmiDevSerialNumberGet(i)
		//获取设备类型标识id
		devTypeId, _ := rsmiDevIdGet(i)
		//型号名称
		devTypeName := type2name[fmt.Sprintf("%x", devTypeId)]
		//设备温度
		temperature, _ := rsmiDevTempMetricGet(i, 0, RSMI_TEMP_CURRENT)
		t, err := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(temperature)/1000.0), 64)
		if err != nil {
			return nil, err
		}
		//设备平均功耗
		powerUsage, _ := rsmiDevPowerAveGet(i, 0)
		pu, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(powerUsage)/1000000.0), 64)
		glog.Infof("\U0001FAAB DCU[%v] power usage : %.0f", i, pu)
		//获取设备功率上限
		powerCap, _ := rsmiDevPowerCapGet(i, 0)
		pc, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(powerCap)/1000000.0), 64)
		glog.Infof("🔋 DCU[%v] power cap : %.0f", i, pc)
		//获取设备内存总量
		memoryCap, _ := rsmiDevMemoryTotalGet(i, RSMI_MEM_TYPE_FIRST)
		mc, _ := strconv.ParseFloat(fmt.Sprintf("%f", float64(memoryCap)/1.0), 64)
		glog.Infof("DCU[%v] memory total: %.0f", i, mc)
		//获取设备内存使用量
		memoryUsed, _ := rsmiDevMemoryUsageGet(i, RSMI_MEM_TYPE_FIRST)
		mu, _ := strconv.ParseFloat(fmt.Sprintf("%f", float64(memoryUsed)/1.0), 64)
		glog.Infof(" DCU[%v] memory used : %.0f ", i, mu)
		//获取设备忙碌时间百分比
		utilizationRate, _ := rsmiDevBusyPercentGet(i)
		ur, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(utilizationRate)/1.0), 64)
		glog.Infof(" DCU[%v] utilization rate : %.0f", i, ur)
		//获取pcie流量信息
		sent, received, maxPktSz, _ := rsmiDevPciThroughputGet(i)
		pcieBwMb, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(received+sent)*float64(maxPktSz)/1024.0/1024.0), 64)
		glog.Infof(" DCU[%v] PCIE  bandwidth : %.0f", i, pcieBwMb)
		//获取设备系统时钟速度列表
		clk, _ := rsmiDevGpuClkFreqGet(i, RSMI_CLK_TYPE_SYS)
		sclk, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(clk.Frequency[clk.Current])/1000000.0), 64)
		glog.Infof(" DCU[%v] SCLK : %.0f", i, sclk)
		monitorInfo := MonitorInfo{
			MinorNumber:     i,
			PciBusNumber:    pciBusNumber,
			DeviceId:        deviceId,
			SubSystemName:   devTypeName,
			Temperature:     t,
			PowerUsage:      pu,
			PowerCap:        pc,
			MemoryCap:       mc,
			MemoryUsed:      mu,
			UtilizationRate: ur,
			PcieBwMb:        pcieBwMb,
			Clk:             sclk,
		}
		monitorInfos = append(monitorInfos, monitorInfo)
	}
	glog.Info("monitorInfos: ", dataToJson(monitorInfos))
	return
}

/*func CollectVDeviceMetrics() (devices []PhysicalDeviceInfo, err error) {


}*/

// AllDeviceInfos 获取所有物理设备及其虚拟设备的信息列表
// @Summary 获取所有物理设备及其虚拟设备的信息列表
// @Description 返回所有物理设备及其虚拟设备的信息
// @Produce json
// @Success 200 {array} PhysicalDeviceInfo "物理设备及其虚拟设备信息列表"
// @Failure 500 {object} error "服务器内部错误"
// @Router /AllDeviceInfos [get]
func AllDeviceInfos() ([]PhysicalDeviceInfo, error) {
	var allDevices []PhysicalDeviceInfo

	// 获取物理设备数量
	deviceCount, err := rsmiNumMonitorDevices()
	if err != nil {
		return nil, err
	}

	// 用于保存所有物理设备的信息
	deviceMap := make(map[int]*PhysicalDeviceInfo)

	// 获取所有物理设备信息
	for i := 0; i < deviceCount; i++ {
		//物理设备支持最大虚拟化设备数量
		maxVDeviceCount, _ := dmiGetMaxVDeviceCount()

		//物理设备使用百分比
		//devPercent, _ := dmiGetDevBusyPercent(i)
		//deviceInfo.Percent = devPercent

		bdfid, err := rsmiDevPciIdGet(i)
		if err != nil {
			return nil, err
		}
		// 解析BDFID
		domain := (bdfid >> 32) & 0xffffffff
		bus := (bdfid >> 8) & 0xff
		dev := (bdfid >> 3) & 0x1f
		function := bdfid & 0x7
		// 格式化PCI ID
		pciBusNumber := fmt.Sprintf("%04x:%02x:%02x.%x", domain, bus, dev, function)
		//设备序列号
		deviceId, _ := rsmiDevSerialNumberGet(i)
		//获取设备类型标识id
		devTypeId, _ := rsmiDevIdGet(i)
		//型号名称
		devTypeName := type2name[fmt.Sprintf("%x", devTypeId)]
		//设备温度
		temperature, _ := rsmiDevTempMetricGet(i, 0, RSMI_TEMP_CURRENT)
		t, err := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(temperature)/1000.0), 64)
		if err != nil {
			return nil, err
		}
		//设备平均功耗
		powerUsage, _ := rsmiDevPowerAveGet(i, 0)
		pu, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(powerUsage)/1000000.0), 64)
		//glog.Infof("\U0001FAAB DCU[%v] power usage : %.0f", i, pu)
		//获取设备功率上限
		powerCap, _ := rsmiDevPowerCapGet(i, 0)
		pc, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(powerCap)/1000000.0), 64)
		//glog.Infof("🔋 DCU[%v] power cap : %.0f", i, pc)
		//获取设备内存总量
		memoryCap, _ := rsmiDevMemoryTotalGet(i, RSMI_MEM_TYPE_FIRST)
		mc, _ := strconv.ParseFloat(fmt.Sprintf("%f", float64(memoryCap)/1.0), 64)
		//glog.Infof("DCU[%v] memory total: %.0f", i, mc)
		//获取设备内存使用量
		memoryUsed, _ := rsmiDevMemoryUsageGet(i, RSMI_MEM_TYPE_FIRST)
		mu, _ := strconv.ParseFloat(fmt.Sprintf("%f", float64(memoryUsed)/1.0), 64)
		//glog.Infof(" DCU[%v] memory used : %.0f ", i, mu)
		//获取设备设备忙碌时间百分比
		utilizationRate, _ := rsmiDevBusyPercentGet(i)
		ur, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(utilizationRate)/1.0), 64)
		//glog.Infof(" DCU[%v] utilization rate : %.0f", i, ur)
		//获取pcie流量信息
		sent, received, maxPktSz, _ := rsmiDevPciThroughputGet(i)
		pcieBwMb, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(received+sent)*float64(maxPktSz)/1024.0/1024.0), 64)
		//glog.Infof(" DCU[%v] PCIE  bandwidth : %.0f", i, pcieBwMb)
		//获取设备系统时钟速度列表
		clk, _ := rsmiDevGpuClkFreqGet(i, RSMI_CLK_TYPE_SYS)
		sclk, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(clk.Frequency[clk.Current])/1000000.0), 64)
		//glog.Infof(" DCU[%v] SCLK : %.0f", i, sclk)
		computeUnit := computeUnitType[devTypeName]
		device := Device{
			MinorNumber:      i,
			PciBusNumber:     pciBusNumber,
			DeviceId:         deviceId,
			SubSystemName:    devTypeName,
			Temperature:      t,
			PowerUsage:       pu,
			PowerCap:         pc,
			MemoryCap:        mc,
			MemoryUsed:       mu,
			UtilizationRate:  ur,
			PcieBwMb:         pcieBwMb,
			Clk:              sclk,
			ComputeUnitCount: computeUnit,
			MaxVDeviceCount:  maxVDeviceCount,
		} // 创建PhysicalDeviceInfo并存入map
		pdi := PhysicalDeviceInfo{
			Device:         device,
			VirtualDevices: []DMIVDeviceInfo{},
		}
		deviceMap[device.MinorNumber] = &pdi
	}

	// 获取虚拟设备数量
	//vDeviceCount, err := dmiGetVDeviceCount()
	//vDeviceCount := deviceCount * 4
	//if err != nil {
	//	return nil, err
	//}
	//// 获取所有虚拟设备信息并关联到对应的物理设备
	//for j := 0; j < vDeviceCount; j++ {
	//	vDeviceInfo, err := dmiGetVDeviceInfo(j)
	//	glog.Infof("vDeviceInfo error: %v", err)
	//	if err == nil {
	//		vDevPercent, _ := dmiGetVDevBusyPercent(j)
	//		vDeviceInfo.Percent = vDevPercent
	//		vDeviceInfo.VMinorNumber = j
	//		// 找到对应的物理设备并将虚拟设备添加到其VirtualDevices中
	//		if pdi, exists := deviceMap[vDeviceInfo.DeviceID]; exists {
	//			pdi.VirtualDevices = append(pdi.VirtualDevices, vDeviceInfo)
	//		}
	//	}
	//	if err != nil {
	//		return nil, fmt.Errorf("Error getting virtual device info for virtual device %d: %s", j, err)
	//	}
	//}

	dirPath := "/etc/vdev"
	// 读取目录中的文件列表
	files, err := os.ReadDir(dirPath)
	if err != nil {
		glog.Errorf("无法读取目录: %v", err)
	}
	// 打印文件数量
	//fmt.Printf("文件数量: %d\n", len(files))
	// 逐个读取并解析每个文件的内容
	for _, file := range files {
		// 确保是文件而不是子目录
		if !file.IsDir() {
			filePath := filepath.Join(dirPath, file.Name())
			config, err := parseConfig(filePath)
			if err != nil {
				glog.Errorf("无法解析文件 %s: %v", filePath, err)
				continue
			}
			//glog.Infof("文件: %s\n配置: %+v\n", filePath, config)
			// 找到对应的物理设备并将虚拟设备添加到其VirtualDevices中
			if pdi, exists := deviceMap[config.DeviceID]; exists {
				pdi.VirtualDevices = append(pdi.VirtualDevices, *config)
			}
		}
	}
	// 将map中的所有PhysicalDeviceInfo转为slice
	for _, pdi := range deviceMap {
		allDevices = append(allDevices, *pdi)
	}
	glog.Infof("allDevices:%v", dataToJson(allDevices))
	return allDevices, nil
}

// PicBusInfo 获取设备的总线信息
// @Summary 获取设备的总线信息
// @Description 根据设备索引返回对应的总线信息（BDF格式）
// @Produce json
// @Param dvInd query int true "设备索引"
// @Success 200 {string} string "返回设备的总线信息"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /PicBusInfo [get]
func PicBusInfo(dvInd int) (picID string, err error) {
	bdfid, err := rsmiDevPciIdGet(dvInd)
	if err != nil {
		return "", err
	}
	// Parse BDFID
	domain := (bdfid >> 32) & 0xffffffff
	bus := (bdfid >> 8) & 0xff
	devID := (bdfid >> 3) & 0x1f
	function := bdfid & 0x7
	// Format and return the bus identifier
	picID = fmt.Sprintf("%04X:%02X:%02X.%X", domain, bus, devID, function)
	return
}

// FanSpeedInfo 获取风扇转速信息
// @Summary 获取风扇转速信息
// @Description 根据设备索引返回当前风扇转速及其占最大转速的百分比
// @Produce json
// @Param dvInd query int true "设备索引"
// @Success 200 {int64} fanLevel "返回当前风扇转速"
// @Success 200 {float64} fanPercentage "返回风扇转速百分比"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /FanSpeedInfo [get]
func FanSpeedInfo(dvInd int) (fanLevel int64, fanPercentage float64, err error) {
	// 当前转速
	fanLevel, err = rsmiDevFanSpeedGet(dvInd, 0)
	if err != nil {
		return 0, 0, err
	}
	// 最大转速
	fanMax, err := rsmiDevFanSpeedMaxGet(dvInd, 0)
	if err != nil {
		return 0, 0, err
	}
	// Calculate fan speed percentage
	fanPercentage = (float64(fanLevel) / float64(fanMax)) * 100
	return
}

// GPUUse 当前GPU使用的百分比
// @Summary 获取当前GPU使用的百分比
// @Description 根据设备索引返回当前GPU的使用百分比
// @Produce json
// @Param dvInd query int true "设备索引"
// @Success 200 {int} percent "返回GPU使用的百分比"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /GPUUse [get]
func GPUUse(dvInd int) (percent int, err error) {
	percent, err = rsmiDevBusyPercentGet(dvInd)
	if err != nil {
		return 0, err
	}
	return
}

// DevID 设备ID的十六进制值
// @Summary 获取设备ID的十六进制值
// @Description 根据设备索引返回设备ID的十六进制值
// @Produce json
// @Param dvInd query int true "设备索引"
// @Success 200 {int} id "返回设备ID的十六进制值"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /DevID [get]
func DevID(dvInd int) (id int, err error) {
	id, err = rsmiDevIdGet(dvInd)
	if err != nil {
		return 0, err
	}
	return
}

// MaxPower 设备的最大功率
// @Summary 获取设备的最大功率
// @Description 根据设备索引返回设备的最大功率（以瓦特为单位）
// @Produce json
// @Param dvInd query int true "设备索引"
// @Success 200 {int64} power "返回设备的最大功率"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /MaxPower [get]
func MaxPower(dvInd int) (power int64, err error) {
	power, err = rsmiDevPowerCapGet(dvInd, 0)
	if err != nil {
		return 0, err
	}
	glog.Infof("Max power: %v", (power / 1000000))
	return (power / 1000000), nil
}

// MemInfo 获取设备的指定内存使用情况
// @Summary 获取设备的指定内存使用情况
// @Description 根据设备索引和内存类型返回内存的使用量和总量
// @Produce json
// @Param dvInd query int true "设备索引"
// @Param memType query string true "内存类型（可选值: vram, vis_vram, gtt）"
// @Success 200 {int64} memUsed "返回指定内存类型的使用量"
// @Success 200 {int64} memTotal "返回指定内存类型的总量"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /MemInfo [get]
func MemInfo(dvInd int, memType string) (memUsed int64, memTotal int64, err error) {
	memType = strings.ToUpper(memType)
	if !contains(memoryTypeL, memType) {
		fmt.Println(dvInd, fmt.Sprintf("Invalid memory type %s", memType))
		return 0, 0, fmt.Errorf("invalid memory type")
	}
	memTypeIndex := RSMIMemoryType(indexOf(memoryTypeL, memType))
	memUsed, err = rsmiDevMemoryUsageGet(dvInd, memTypeIndex)
	if err != nil {
		return memUsed, memTotal, err
	}
	fmt.Println(dvInd, fmt.Sprintf("memUsed: %d", memUsed))
	memTotal, err = rsmiDevMemoryTotalGet(dvInd, memTypeIndex)
	if err != nil {
		return memUsed, memTotal, err
	}
	fmt.Println(dvInd, fmt.Sprintf("memTotal: %d", memTotal))
	return
}

// DeviceInfos 获取设备信息列表
// @Summary 获取设备信息列表
// @Description 返回所有设备的详细信息列表
// @Produce json
// @Success 200 {array} DeviceInfo "返回设备信息列表"
// @Failure 500 {object} error "服务器内部错误"
// @Router /DeviceInfos [get]
func DeviceInfos() (deviceInfos []DeviceInfo, err error) {
	numDevices, err := rsmiNumMonitorDevices()
	if err != nil {
		return nil, err
	}
	for i := 0; i < numDevices; i++ {
		bdfid, err := rsmiDevPciIdGet(i)
		if err != nil {
			return nil, err
		}
		// 解析BDFID
		domain := (bdfid >> 32) & 0xffffffff
		bus := (bdfid >> 8) & 0xff
		dev := (bdfid >> 3) & 0x1f
		function := bdfid & 0x7
		// 格式化PCI ID
		pciBusNumber := fmt.Sprintf("%04X:%02X:%02X.%X", domain, bus, dev, function)
		//设备序列号
		deviceId, _ := rsmiDevSerialNumberGet(i)
		//获取设备类型标识id
		devTypeId, _ := rsmiDevIdGet(i)
		devType := fmt.Sprintf("%x", devTypeId)
		//型号名称
		devTypeName := type2name[devType]
		//获取设备内存总量
		memoryTotal, _ := rsmiDevMemoryTotalGet(i, RSMI_MEM_TYPE_FIRST)
		mt, _ := strconv.ParseFloat(fmt.Sprintf("%f", float64(memoryTotal)/1.0), 64)
		glog.Info(" DCU[%v] memory total memory total: %.0f", i, mt)
		//获取设备内存使用量
		memoryUsed, _ := rsmiDevMemoryUsageGet(i, RSMI_MEM_TYPE_FIRST)
		mu, _ := strconv.ParseFloat(fmt.Sprintf("%f", float64(memoryUsed)/1.0), 64)
		glog.Info(" DCU[%v] memory used :%.0f", i, mu)
		computeUnit := computeUnitType[devTypeName]
		glog.Info(" DCU[%v] computeUnit : %.0f", i, computeUnit)
		deviceInfo := DeviceInfo{
			DvInd:        i,
			DeviceId:     deviceId,
			DevType:      devType,
			DevTypeName:  devTypeName,
			PciBusNumber: pciBusNumber,
			MemoryTotal:  mt,
			MemoryUsed:   mu,
			ComputeUnit:  computeUnit,
		}
		deviceInfos = append(deviceInfos, deviceInfo)
	}
	glog.Info("deviceInfos: ", dataToJson(deviceInfos))
	return
}

// ProcessName 获取指定PID的进程名
// @Summary 获取指定PID的进程名
// @Description 根据进程ID（PID）返回对应的进程名称
// @Produce json
// @Param pid query int true "进程ID"
// @Success 200 {string} string "返回进程名称"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /ProcessName [get]
func ProcessName(pid int) string {
	if pid < 1 {
		glog.Info("PID must be greater than 0")
		return "UNKNOWN"
	}
	cmd := exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "comm=")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		glog.Info("Error executing command:", err)
		return "UNKNOWN"
	}
	pName := out.String()
	if pName == "" {
		return "UNKNOWN"
	}
	// Remove the substrings surrounding from process name (b' and \n')
	pName = strings.TrimPrefix(pName, "b'")
	pName = strings.TrimSuffix(pName, "\\n'")
	glog.Info("Process name: %s\n", pName)
	return strings.TrimSpace(pName)
}

// PerfLevel 获取设备的当前性能水平
// @Summary 获取设备的当前性能水平
// @Description 返回指定设备的当前性能等级
// @Produce json
// @Param dvInd query int true "设备索引"
// @Success 200 {string} string "返回当前性能水平"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /PerfLevel [get]
func PerfLevel(dvInd int) (perf string, err error) {
	level, err := rsmiDevPerfLevelGet(dvInd)
	if err != nil {
		return perf, err
	}
	perf = perfLevelString(int(level))
	glog.Infof("Perf level: %v", perf)
	return
}

// getPid 获取特定应用程序的进程 ID
func PidByName(name string) (pid string, err error) {
	glog.Info("pidName: %s\n", name)
	cmd := exec.Command("pidof", name)
	output, err := cmd.Output()
	glog.Info("output:", output)
	if err != nil {
		glog.Info("Error: %v\nOutput: %s", err, string(output))
	} else {
		glog.Info("Output: %s", string(output))
	}
	// 移除末尾的换行符并返回 PID
	pid = strings.TrimSpace(string(output))
	glog.Info("pid: %s\n", pid)
	return
}

// Power 获取设备的平均功耗
// @Summary 获取设备的平均功耗
// @Description 返回指定设备的平均功耗
// @Produce json
// @Param dvInd query int true "设备索引"
// @Success 200 {int64} int64 "返回平均功耗（瓦特）"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /Power [get]
func Power(dvInd int) (power int64, err error) {
	powerAve, err := rsmiDevPowerAveGet(dvInd, 0)
	power = powerAve / 1000000
	glog.Infof("Power: %v", power)
	if err != nil {
		return power, err
	}
	return
}

// EccStatus 获取GPU块的ECC状态
// @Summary 获取GPU块的ECC状态
// @Description 返回指定GPU块的ECC状态
// @Produce json
// @Param dvInd query int true "设备索引"
// @Param block query string true "GPU块"
// @Success 200 {string} string "返回ECC状态"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /EccStatus [get]
func EccStatus(dvInd int, block RSMIGpuBlock) (state string, err error) {
	eccStatus, err := rsmiDevEccStatusGet(dvInd, block)
	state = rasErrStaleMachine[eccStatus]
	return
}

// Temperature 获取设备温度
// @Summary 获取设备温度
// @Description 返回指定设备的当前温度
// @Produce json
// @Param dvInd query int true "设备索引"
// @Param sensorType query int true "传感器类型"
// @Success 200 {float64} float64 "返回温度（摄氏度）"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /Temperature [get]
func Temperature(dvInd int, sensorType int) (temp float64, err error) {
	deviceTemp, err := rsmiDevTempMetricGet(dvInd, sensorType, RSMI_TEMP_CURRENT)
	temp, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", float64(deviceTemp)/1000.0), 64)
	glog.Infof("device Temperature:%v", temp)
	return
}

// VbiosVersion 获取设备的VBIOS版本
// @Summary 获取设备的VBIOS版本
// @Description 返回指定设备的VBIOS版本
// @Produce json
// @Param dvInd query int true "设备索引"
// @Success 200 {string} string "返回VBIOS版本"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /VbiosVersion [get]
func VbiosVersion(dvInd int) (vbios string, err error) {
	vbios, err = rsmiDevVbiosVersionGet(dvInd, 256)
	glog.Infof("VbiosVersion:%v", vbios)
	return
}

// Version 获取当前系统的驱动程序版本
// @Summary 获取当前系统的驱动程序版本
// @Description 返回指定组件的驱动程序版本
// @Produce json
// @Param component query string true "驱动组件"
// @Success 200 {string} string "返回驱动程序版本"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /Version [get]
func Version(component RSMISwComponent) (varStr string, err error) {
	varStr, err = rsmiVersionStrGet(component, 256)
	glog.Infof("component; Version:%v,%v", component, varStr)
	return
}

// ResetClocks 将设备的时钟重置为默认值
// @Summary 重置设备时钟
// @Description 重置指定设备的时钟和性能等级为默认值
// @Produce json
// @Param dvIdList body []int true "设备ID列表"
// @Success 200 {array} FailedMessage "返回失败消息列表"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /ResetClocks [post]
func ResetClocks(dvIdList []int) (failedMessage []FailedMessage) {
	errorMap := make(map[int][]string)
	glog.Info(" Reset Clocks ")
	for _, device := range dvIdList {
		// Reset OverDrive
		err := rsmiDevOverdriveLevelSet(device, 0)
		if err != nil {
			errorMap[device] = append(errorMap[device], "Unable to reset OverDrive")
			glog.Errorf("Unable to reset OverDrive, device: %v, error: %v", device, err)
		}
		// Reset PerfLevel
		err = rsmiDevPerfLevelSet(device, RSMI_DEV_PERF_LEVEL_AUTO)
		if err != nil {
			errorMap[device] = append(errorMap[device], "Unable to reset clocks")
			glog.Errorf("Unable to reset clocks, device: %v, error: %v", device, err)
		}

		// Set performance level to auto
		err = rsmiDevPerfLevelSet(device, RSMI_DEV_PERF_LEVEL_AUTO)
		if err != nil {
			errorMap[device] = append(errorMap[device], "Unable to set performance level to auto")
			glog.Errorf("Unable to set performance level to auto, device: %v, error: %v", device, err)
		}
	}
	for id, msg := range errorMap {
		failedMessage = append(failedMessage, FailedMessage{ID: id, ErrorMsg: strings.Join(msg, "; ")})
	}
	return
}

// ResetFans 复位风扇驱动控制
// @Summary 复位风扇控制
// @Description 重置指定设备的风扇控制为默认值
// @Produce json
// @Param dvIdList body []int true "设备ID列表"
// @Success 200 {string} string "复位成功"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /ResetFans [post]
func ResetFans(dvIdList []int) (err error) {
	for _, id := range dvIdList {
		err := rsmiDevFanReset(id, 0)
		glog.Infof("Resetting fan :%v", id)
		if err != nil {
			glog.Errorf("Unable reset Fan dvId:%v ,err:%v", id, err)
		}
	}
	return
}

// ResetProfile 重置设备的配置文件
// @Summary 重置指定设备的电源配置文件和性能级别
// @Produce json
// @Param dvIdList body []int true "设备ID列表"
// @Success 200 {array} FailedMessage "返回失败的设备及其错误信息"
// @Router /ResetProfile [post]
func ResetProfile(dvIdList []int) (failedMessage []FailedMessage) {
	errorMap := make(map[int][]string)
	for _, id := range dvIdList {
		err := rsmiDevPowerProfileSet(id, 0, RSMI_PWR_PROF_PRST_BOOTUP_DEFAULT)
		if err != nil {
			errorMap[id] = append(errorMap[id], "Unable to reset OverDrive")
			glog.Errorf("Unable to reset OverDrive, device: %v, error: %v", id, err)
		}
		// Reset PerfLevel
		err = rsmiDevPerfLevelSet(id, RSMI_DEV_PERF_LEVEL_AUTO)
		if err != nil {
			errorMap[id] = append(errorMap[id], "Unable to reset PerfLevel")
			glog.Errorf("Unable to reset PerfLevel, device: %v, error: %v", id, err)
		}
	}
	for id, msg := range errorMap {
		failedMessage = append(failedMessage, FailedMessage{ID: id, ErrorMsg: strings.Join(msg, "; ")})
	}
	return
}

// ResetXGMIErr 重置设备的XGMI错误状态
// @Summary 重置指定设备的XGMI错误状态
// @Produce json
// @Param dvIdList body []int true "设备ID列表"
// @Success 200 {array} FailedMessage "返回失败的设备及其错误信息"
// @Router /ResetXGMIErr [post]
func ResetXGMIErr(dvIdList []int) (failedMessage []FailedMessage) {
	errorMap := make(map[int][]string)
	for _, id := range dvIdList {
		err := rsmiDevXgmiErrorReset(id)
		if err != nil {
			errorMap[id] = append(errorMap[id], "Unable to reset XGMI error")
			glog.Errorf("Unable to reset XGMI error, device: %v, error: %v", id, err)
		}
	}
	for id, msg := range errorMap {
		failedMessage = append(failedMessage, FailedMessage{ID: id, ErrorMsg: strings.Join(msg, "; ")})
	}
	return
}

// XGMIErrorStatus 获取XGMI错误状态
// @Summary 获取指定设备的XGMI错误状态
// @Produce json
// @Param dvInd query int true "设备索引"
// @Success 200 {object} RSMIXGMIStatus "返回设备的XGMI错误状态"
// @Router /XGMIErrorStatus [get]
func XGMIErrorStatus(dvInd int) (status RSMIXGMIStatus, err error) {
	return rsmiDevXGMIErrorStatus(dvInd)
}

// XGMIHiveIdGet 获取设备的XGMI hive id
// @Summary 获取指定设备的XGMI hive id
// @Produce json
// @Param dvInd query int true "设备索引"
// @Success 200 {int64} int64 "返回设备的XGMI hive id"
// @Router /XGMIHiveIdGet [get]
func XGMIHiveIdGet(dvInd int) (hiveId int64, err error) {
	return rsmiDevXgmiHiveIdGet(dvInd)
}

// ResetPerfDeterminism 重置Performance Determinism
// @Summary 重置指定设备的性能决定性设置
// @Produce json
// @Param dvIdList body []int true "设备ID列表"
// @Success 200 {array} FailedMessage "返回失败的设备及其错误信息"
// @Router /ResetPerfDeterminism [post]
func ResetPerfDeterminism(dvIdList []int) (failedMessage []FailedMessage) {
	errorMap := make(map[int][]string)
	for _, device := range dvIdList {
		// Set performance level to auto
		err := rsmiDevPerfLevelSet(device, RSMI_DEV_PERF_LEVEL_AUTO)
		if err != nil {
			errorMap[device] = append(errorMap[device], "Unable to diable performance determinism")
			glog.Errorf("Unable to diable performance determinism, device: %v, error: %v", device, err)
		}
	}
	for id, msg := range errorMap {
		failedMessage = append(failedMessage, FailedMessage{ID: id, ErrorMsg: strings.Join(msg, "; ")})
	}
	return
}

// 为设备选定的时钟类型设定相应的频率范围
func SetClockRange(dvIdList []int, clkType string, minvalue string, maxvalue string, autoRespond bool) {
	if clkType != "sclk" && clkType != "mclk" {
		glog.Infof("device :%v,Invalid range identifier %v", dvIdList, clkType)
		glog.Infof("Unsupported range type %s", clkType)
		return
	}
	minVal, errMin := strconv.ParseInt(minvalue, 10, 64)
	maxVal, errMax := strconv.ParseInt(maxvalue, 10, 64)
	if errMin != nil || errMax != nil {
		glog.Errorf("Unable to set %s range", clkType)
		glog.Infof("%s or %s is not an integer", minvalue, maxvalue)
		return
	}
	confirmOutOfSpecWarning(autoRespond)
	for _, device := range dvIdList {
		err := rsmiDevClkRangeSet(device, minVal, maxVal, rsmiClkNamesDict[clkType])
		if err == nil {
			glog.Errorf("device:%v Successfully set %v from %v(MHz) to %v(MHz)", clkType, minVal, maxVal)
		} else {
			glog.Errorf("device:%v Unable to set %v from %v(MHz) to %v(MHz)", device, clkType, minVal, maxVal)
		}
	}
}

//设置电压曲线

// SetPowerPlayTableLevel 设置 PowerPlay 级别
// @Summary Set PowerPlay table level for devices
// @Description This function sets the PowerPlay table level for a list of devices. It checks the validity of the input values and adjusts the voltage settings accordingly.
// @Tags Device
// @Param dvIdList body []int true "List of device IDs"
// @Param clkType query string true "Clock type (sclk or mclk)"
// @Param point query string true "Voltage point"
// @Param clk query string true "Clock value in MHz"
// @Param volt query string true "Voltage value in mV"
// @Param autoRespond query bool false "Automatically respond to out-of-spec warnings"
// @Success 200 {string} string "PowerPlay table level set successfully"
// @Failure 400 {string} string "Invalid input or unable to set PowerPlay table level"
// @Router /SetPowerPlayTableLevel [post]
func SetPowerPlayTableLevel(dvIdList []int, clkType string, point string, clk string, volt string, autoRespond bool) {
	value := fmt.Sprintf("%s %s %s", point, clk, volt)
	_, errPoint := strconv.Atoi(point)
	_, errClk := strconv.Atoi(clk)
	_, errVolt := strconv.Atoi(volt)
	if errPoint != nil || errClk != nil || errVolt != nil {
		glog.Infof("Unable to set PowerPlay table level")
		glog.Infof("Non-integer characters are present in %s", value)
		return
	}
	confirmOutOfSpecWarning(autoRespond)
	for _, device := range dvIdList {
		pointVal, _ := strconv.Atoi(point)
		clkVal, _ := strconv.Atoi(clk)
		voltVal, _ := strconv.Atoi(volt)
		if clkType == "sclk" || clkType == "mclk" {
			err := rsmiDevOdVoltInfoSet(device, pointVal, clkVal, voltVal)
			if err == nil {
				glog.Infof("device:%v Successfully set voltage point %v to %v(MHz) %v(mV)", point, clk, volt)
			} else {
				glog.Errorf("device:%v Unable to set voltage point %v to %v(MHz) %v(mV)", point, clk, volt)

			}
		} else {
			glog.Errorf("device:%v Unable to set %s range", clkType)
			glog.Infof("Unsupported range type %s", clkType)
		}
	}
}

// SetClockOverDrive 设置时钟速度为 OverDrive
// @Summary Set Clock OverDrive for devices
// @Description This function sets the Clock OverDrive level for a list of devices. It adjusts the clock speed and ensures the performance level is set to manual if needed.
// @Tags Device
// @Param dvIdList body []int true "List of device IDs"
// @Param clktype query string true "Clock type (sclk or mclk)"
// @Param value query string true "OverDrive value as a percentage (0-20%)"
// @Param autoRespond query bool false "Automatically respond to out-of-spec warnings"
// @Success 200 {string} string "Clock OverDrive set successfully"
// @Failure 400 {string} string "Invalid input or unable to set Clock OverDrive"
// @Router /SetClockOverDrive [post]
func SetClockOverDrive(dvIdList []int, clktype string, value string, autoRespond bool) {
	glog.Infof("Set Clock OverDrive Range: 0 to 20%")
	intValue, err := strconv.Atoi(value)
	if err != nil {
		glog.Infof("Unable to set OverDrive level")
		glog.Errorf("%s it is not an integer", value)
		return
	}

	confirmOutOfSpecWarning(autoRespond)

	for _, device := range dvIdList {
		if intValue < 0 {
			glog.Errorf("Unable to set OverDrive device: %v", device)
			glog.Infof("Overdrive cannot be less than 0%")
			return
		}
		if intValue > 20 {
			glog.Infof("device:%v,Setting OverDrive to 20%", device)
			glog.Infof("OverDrive cannot be set to a value greater than 20%")
			intValue = 20
		}
		perf, _ := PerfLevel(device)
		if perf != "MANUAL" {
			err := rsmiDevPerfLevelSet(device, RSMI_DEV_PERF_LEVEL_MANUAL)

			if err == nil {
				glog.Infof("device:%v Performance level set to manual", device)
			} else {
				glog.Errorf("device:%v Unable to set performance level to manual")
			}
		}
		if clktype == "mclk" {
			fsFile := fmt.Sprintf("/sys/class/drm/card%d/device/pp_mclk_od", device)
			if _, err := os.Stat(fsFile); os.IsNotExist(err) {
				glog.Infof("Unable to write to sysfs file")
				glog.Warning("does not exist ", fsFile)
				continue
			}
			f, err := os.OpenFile(fsFile, os.O_WRONLY, 0644)
			if err == nil {
				glog.Infof("Unable to write to sysfs file %v", fsFile)
				glog.Warning("IO or OS error")
				continue
			}
			defer f.Close()
			_, err = f.WriteString(fmt.Sprintf("%v", intValue))
			if err != nil {
				glog.Infof("Unable to write to sysfs file %v", fsFile)
				glog.Warning("IO or OS error")
				continue
			}
			glog.Infof("device%v Successfully set %s OverDrive to %d%%", device, clktype, intValue)
		} else if clktype == "sclk" {
			err := rsmiDevOverdriveLevelSet(device, intValue)

			if err == nil {
				glog.Infof("device:%v Successfully set %s OverDrive to %d%%", device, clktype, intValue)
			} else {
				glog.Errorf("device:%v Unable to set %s OverDrive to %d%%", device, clktype, intValue)
			}
		} else {
			glog.Errorf("device:%v Unable to set OverDrive", device)
			glog.Errorf("Unsupported clock type %v", clktype)
		}
	}
}

// 设置时钟频率级别以启用性能确定性
func SetPerfDeterminism(dvIdList []int, clkvalue string) (failedMessage []FailedMessage, err error) {
	// 验证 clkvalue 是否为有效的整数
	intValue, err := strconv.ParseInt(clkvalue, 10, 64)
	if err != nil {
		glog.Errorf("Unable to set Performance Determinism")
		glog.Errorf("clkvalue:%v is not an integer", clkvalue)
		return failedMessage, fmt.Errorf("clkvalue:%v is not an integer", clkvalue)
	}

	errorMap := make(map[int][]string)
	// 遍历每个设备并设置性能确定性模式
	for _, device := range dvIdList {
		err := rsmiPerfDeterminismModeSet(device, intValue)
		if err != nil {
			errorMap[device] = append(errorMap[device], "Unable to set performance determinism")
			glog.Errorf("Unable to set performance determinism and clock frequency to %v for device %v", clkvalue, device)
		}
	}
	for id, msg := range errorMap {
		failedMessage = append(failedMessage, FailedMessage{ID: id, ErrorMsg: strings.Join(msg, "; ")})
	}
	return
}

// 设置风扇转速 [0-255] Fan speed level
func SetFanSpeed(dvIdList []int, fan string) {
	for _, device := range dvIdList {
		var fanLevel int64
		var err error
		lastChar := fan[len(fan)-1:]
		if lastChar == "%" {
			percentValue, err := strconv.Atoi(fan[:len(fan)-1])
			if err != nil {
				glog.Errorf("Invalid fan speed percentage: %s", fan)
				continue
			}
			fanLevel = int64(percentValue * 255 / 100)
		} else {
			fanLevel, err = strconv.ParseInt(fan, 10, 64)
			if err != nil {
				glog.Errorf("Invalid fan speed value: %s", fan)
				continue
			}
		}
		glog.Infof("Setting fan speed fanLevel value to %v", fanLevel)
		err = rsmiDevFanSpeedSet(device, 0, fanLevel)
		if err != nil {
			log.Printf("Failed to set fan speed for device %d", device)
		}
	}
}

// 获取设备的风扇速度，实际转速
func DevFanRpms(dvInd int) (speed int64, err error) {
	return rsmiDevFanRpmsGet(dvInd, 0)
}

// 设置设备性能 level:auto、low、high、normal
func SetPerformanceLevel(deviceList []int, level string) (failedMessages []FailedMessage) {
	for _, device := range deviceList {
		devPerfLevel, valid := validLevels[level]
		if !valid {
			glog.Errorf("device :%v Unable to set Performance Level, Invalid Performance level: %v", device, level)
		} else {
			err := rsmiDevPerfLevelSet(device, devPerfLevel)
			if err != nil {
				glog.Errorf("device:%v Failed to set performance level to %v", device, level)
				failedMessages = append(failedMessages, FailedMessage{
					ID:       device,
					ErrorMsg: fmt.Sprintf("Failed to set performance level to %v", level),
				})
			}
		}
	}
	return
}

// 设置功率配置
func SetProfile(dvIdList []int, profile string) (failedMessages []FailedMessage) {

	for _, device := range dvIdList {
		// 获取先前的配置文件
		status, err := rsmiDevPowerProfilePresetsGet(device, 0)
		glog.Infof("status.Current: %v, int:%v", status.Current, int(status.Current))
		if err == nil {
			previousProfile := profileString(int(status.Current))

			// 确定期望的配置文件
			glog.Infof("previousProfile value: %v", previousProfile)
			glog.Infof("desiredProfile value: %v", profile)
			glog.Infof("previousProfile and desiredProfile:%v", profile == previousProfile)
			if profile == "UNKNOWN" {
				glog.Errorf("device:%v Unable to set profile to: %v (UNKNOWN profile)", device, profile)
				failedMessages = append(failedMessages, FailedMessage{ID: device, ErrorMsg: fmt.Sprintf("Unable to set profile to: %s (UNKNOWN profile)", profile)})
				continue
			}

			// 设置配置文件
			if previousProfile == profile {
				glog.Infof("device:%v Profile was already set to%v", device, previousProfile)
			} else {
				err := rsmiDevPowerProfileSet(device, 0, profileEnum(profile))
				if err == nil {
					// 获取当前配置文件
					profileStatus, err := rsmiDevPowerProfilePresetsGet(device, 0)
					if err == nil {
						currentProfile := profileString(int(profileStatus.Current))
						if currentProfile == profile {
							glog.Infof("device:%v Successfully set profile to:%v", device, profile)
						} else {
							glog.Errorf("device:%v Failed to set profile to: %v", device, profile)
							failedMessages = append(failedMessages, FailedMessage{ID: device, ErrorMsg: fmt.Sprintf("Failed to set profile to: %s", profile)})
						}
					}
				} else {
					failedMessages = append(failedMessages, FailedMessage{ID: device, ErrorMsg: fmt.Sprintf("Failed to set profile to: %s", profile)})
				}
			}
		}
	}

	return
}

func DevPowerProfileSet(dvInd int, reserved int, profile RSNIPowerProfilePresetMasks) (err error) {
	return rsmiDevPowerProfileSet(dvInd, reserved, profile)
}

func DevPowerProfilePresetsGet(dvInd, sensorInd int) (powerProfileStatus RSMPowerProfileStatus, err error) {
	return rsmiDevPowerProfilePresetsGet(dvInd, sensorInd)
}

// 获取设备总线信息
func GetBus(device int) (picId string, err error) {

	bdfid, err := rsmiDevPciIdGet(device)
	if err != nil {
		return picId, err
	}
	domain := (bdfid >> 32) & 0xffffffff
	bus := (bdfid >> 8) & 0xff
	dev := (bdfid >> 3) & 0x1f
	function := bdfid & 0x7
	picId = fmt.Sprintf("%04X:%02X:%02X.%X", domain, bus, dev, function)
	return
}

//显示设备的概要信息
//func ShowAllConcise()  {
//}

// 显示设备硬件信息
func ShowAllConciseHw(dvIdList []int) {
	header := []string{"GPU", "DID", "GFX RAS", "SDMA RAS", "UMC RAS", "VBIOS", "BUS"}
	headWidths := make([]int, len(header))
	for i, head := range header {
		headWidths[i] = len(head) + 2
	}

	values := make(map[string][]string)
	for _, device := range dvIdList {
		gpuid, _ := rsmiDevIdGet(device)
		gfxRas, _ := EccStatus(device, RSMIGpuBlockGFX)
		sdmaRas, _ := EccStatus(device, RSMIGpuBlockSDMA)
		umcRas, _ := EccStatus(device, RSMIGpuBlockUMC)
		vbios, _ := VbiosVersion(device)
		bus, _ := GetBus(device)
		values[fmt.Sprintf("card%d", device)] = []string{
			fmt.Sprintf("GPU%d", device), strconv.Itoa(gpuid), gfxRas, sdmaRas, umcRas, vbios, bus,
		}
	}

	valWidths := make(map[int][]int)
	for _, device := range dvIdList {
		valWidths[device] = make([]int, len(values[fmt.Sprintf("card%d", device)]))
		for i, val := range values[fmt.Sprintf("card%d", device)] {
			valWidths[device][i] = len(val) + 2
		}
	}
	maxWidths := headWidths
	for _, device := range dvIdList {
		for col := range valWidths[device] {
			if valWidths[device][col] > maxWidths[col] {
				maxWidths[col] = valWidths[device][col]
			}
		}
	}

	for i, head := range header {
		fmt.Printf("%-*s", maxWidths[i], head)
	}
	fmt.Println()

	for _, device := range dvIdList {
		for i, val := range values[fmt.Sprintf("card%d", device)] {
			fmt.Printf("%-*s", maxWidths[i], val)
		}
		fmt.Println()
	}

}

// 显示时钟信息
func ShowClocks(dvIdList []int) {

	for _, device := range dvIdList {
		for clkType, clkID := range rsmiClkNamesDict {

			freq, err := rsmiDevGpuClkFreqGet(device, clkID)
			if err == nil {
				glog.Infof("device:%v Supported %v frequencies on GPU%v", device, clkType, device)
				for x := 0; x < int(freq.NumSupported); x++ {
					fr := fmt.Sprintf("%dMhz", freq.Frequency[x]/1000000)
					if uint32(x) == freq.Current {
						glog.Infof("Device %d: %d %s *", device, x, fr)
					} else {
						glog.Infof("Device %d: %d %s", device, x, fr)
					}
				}
			} else {
				glog.Errorf("device:%v clkType:%v frequency is unsupported", device, clkType)

			}
		}
		bw, err := rsmiDevPciBandwidthGet(device)
		if err == nil {
			glog.Infof("Supported PCIe frequencies on GPU%d", device)
			for x := 0; x < int(bw.TransferRate.NumSupported); x++ {
				fr := fmt.Sprintf("%.1fGT/s x%d", float64(bw.TransferRate.Frequency[x])/1000000000, bw.lanes[x])
				if uint32(x) == bw.TransferRate.Current {
					glog.Infof("Device %d: %d %s *", device, x, fr)
				} else {
					glog.Infof("Device %d: %d %s", device, x, fr)
				}
			}
		}
	}
}

// 展示风扇转速和风扇级别
func ShowCurrentFans(dvIdList []int, printJSON bool) {
	glog.Info("--------- Current Fan Metric ---------")
	var sensorInd uint32 = 0

	for _, device := range dvIdList {
		fanLevel, fanSpeed, err := FanSpeedInfo(device)
		if err != nil {
			glog.Errorf("Unable to detect fan speed for GPU %v: %v", device, err)
			continue
		}

		fanSpeed = float64(int64(fanSpeed + 0.5)) // 四舍五入

		if fanLevel == 0 || fanSpeed == 0 {
			glog.Infof("Device %v: Unable to detect fan speed", device)
			glog.Infof("Current fan speed is: %v", fanSpeed)
			glog.Infof("Current fan level is: %v", fanLevel)
			glog.Infof("GPU might be cooled with a non-PWM fan")
			continue
		}
		if printJSON {
			glog.Infof("Device %v: Fan speed (level): %v", device, fanLevel)
			glog.Infof("Device %v: Fan speed (%%): %.0f", device, fanSpeed)
		} else {
			glog.Infof("Device %v: Fan Level: %d (%.0f%%)", device, fanLevel, fanSpeed)
		}

		rpmSpeed, err := rsmiDevFanRpmsGet(device, int(sensorInd))
		if err == nil {
			glog.Infof("Device %v: Fan RPM: %v", device, rpmSpeed)
		} else {
			glog.Errorf("Device %v: Error getting fan RPM: %v", device, err)
		}
	}
	glog.Info("--------------------------------------")
}

// 显示所有设备的所有可用温度传感器的温度
func ShowCurrentTemps(dvIdList []int) (temperatureInfos []TemperatureInfo, err error) {
	glog.Info("--------- Temperature ---------")
	for _, device := range dvIdList {
		sensorTemps := make(map[string]float64)
		for _, sensor := range tempTypeList {
			temp, err := Temperature(device, sensor.Type)
			if err != nil {
				glog.Errorf("Error getting temperature for device %d sensor %d: %v", device, sensor.Type, err)
			} else {
				glog.Infof("Device %d Temperature (Sensor %s): %.2f°C", device, sensor.Name, temp)
				sensorTemps[sensor.Name] = temp
			}
			deviceTempInfo := TemperatureInfo{
				DeviceID:    device,
				SensorTemps: sensorTemps,
			}
			temperatureInfos = append(temperatureInfos, deviceTempInfo)
		}
	}
	glog.Info("--------------------------------")
	glog.Infof("temperatureInfos:%v", dataToJson(temperatureInfos))
	return
}

// 显示给定设备列表中指定固件类型的固件版本信息
func ShowFwInfo(dvIdList []int, fwType []string) (fwInfos []FirmwareInfo, err error) {
	var firmwareBlocks []string
	if len(fwType) == 0 || contains(fwType, "all") {
		firmwareBlocks = fwBlockNames
	} else {
		for _, name := range fwType {
			if contains(fwBlockNames, strings.ToUpper(name)) {
				firmwareBlocks = append(firmwareBlocks, strings.ToUpper(name))
			}
		}
	}
	for _, device := range dvIdList {
		fwVerMap := make(map[string]string)
		for _, fwName := range firmwareBlocks {
			fwNameUpper := strings.ToUpper(fwName)
			fwVersion, err := rsmiDevFirmwareVersionGet(device, RSMIFwBlock(indexOf(fwBlockNames, fwNameUpper)))
			if err != nil {
				glog.Errorf("Error getting firmware version for device %v firmware block %v: %v", device, fwNameUpper, err)
				continue
			}

			var formattedFwVersion string
			if fwNameUpper == "VCN" || fwNameUpper == "VCE" || fwNameUpper == "UVD" || fwNameUpper == "SOS" {
				formattedFwVersion = fmt.Sprintf("0x%s", strings.ToUpper(fmt.Sprintf("%08x", fwVersion)))
			} else if fwNameUpper == "TA XGMI" || fwNameUpper == "TA RAS" || fwNameUpper == "SMC" {
				formattedFwVersion = fmt.Sprintf("%02d.%02d.%02d.%02d",
					(fwVersion>>24)&0xFF, (fwVersion>>16)&0xFF, (fwVersion>>8)&0xFF, fwVersion&0xFF)
			} else if fwNameUpper == "ME" || fwNameUpper == "MC" || fwNameUpper == "CE" {
				formattedFwVersion = fmt.Sprintf("\t\t%d", fwVersion)
			} else {
				formattedFwVersion = fmt.Sprintf("\t%d", fwVersion)
			}

			fwVerMap[fwNameUpper] = formattedFwVersion
			glog.Infof("Device %v %v firmware version: %v", device, fwNameUpper, formattedFwVersion)
		}
		fwInfos = append(fwInfos, FirmwareInfo{
			DeviceID:    device,
			FirmwareVer: fwVerMap,
		})
	}
	glog.Infof("fwInfos:%v", dataToJson(fwInfos))
	return
}

// 获取进程列表
func PidList() (pidList []string, err error) {
	processInfo, numItems, err := rsmiComputeProcessInfoGet()
	if err != nil {
		return nil, err
	}
	if numItems == 0 {
		return
	}
	for i := 0; i < numItems; i++ {
		pidList = append(pidList, fmt.Sprintf("%d", processInfo[i].ProcessID))
	}
	glog.Infof("pidList:%v", pidList)
	return
}

// 获取设备的粗粒度利用率
func GetCoarseGrainUtil(device int, typeName *string) (utilizationCounters []RSMIUtilizationCounter, err error) {
	var length int

	if typeName != nil {
		// 获取特定类型的利用率计数器
		var i RSMIUtilizationCounterType
		var found bool
		for index, name := range utilizationCounterName {
			if name == *typeName {
				i = RSMIUtilizationCounterType(index)
				found = true
				break
			}
		}
		if !found {
			glog.Infof("No such coarse grain counter type: %v", *typeName)
			return nil, fmt.Errorf("no such coarse grain counter type")
		}
		length = 1
		utilizationCounters = make([]RSMIUtilizationCounter, length)
		utilizationCounters[0].Type = i
	} else {
		// 获取所有类型的利用率计数器
		length = int(RSMI_UTILIZATION_COUNTER_LAST) + 1
		utilizationCounters = make([]RSMIUtilizationCounter, length)
		for i := 0; i < length; i++ {
			utilizationCounters[i].Type = RSMIUtilizationCounterType(i)
		}
	}
	_, err = rsmiUtilizationCountGet(device, utilizationCounters, length)
	if err != nil {
		return nil, err
	}
	return
}

// DCU使用率
func ShowGpuUse(dvIdList []int) (deviceUseInfos []DeviceUseInfo, err error) {
	fmt.Printf(" time GPU is busy\n ")

	for _, device := range dvIdList {
		deviceUseInfo := DeviceUseInfo{
			DeviceID:    device,
			Utilization: make(map[string]uint64),
		}

		// 获取 GPU 使用百分比
		percent, err := GPUUse(device)
		if err != nil {
			fmt.Printf("Device %d: GPU use Unsupported\n", device)
			deviceUseInfo.GPUUsePercent = -1

		} else {
			fmt.Printf("Device %d: GPU use (%%) %d\n", device, percent)
			deviceUseInfo.GPUUsePercent = percent
		}

		// 获取粗粒度利用率
		typeName := "GFX Activity"
		utilCounters, err := GetCoarseGrainUtil(device, &typeName)
		if err != nil {
			fmt.Printf("Device %d: Error getting coarse grain utilization: %v\n", device, err)
		} else {
			for _, counter := range utilCounters {
				fmt.Printf("Device %d: %s %d\n", device, utilizationCounterName[counter.Type], counter.Value)
				if int(counter.Type) < len(utilizationCounterName) {
					deviceUseInfo.Utilization[utilizationCounterName[counter.Type]] = counter.Value
				}
			}
		}
		deviceUseInfos = append(deviceUseInfos, deviceUseInfo)
	}
	glog.Infof("deviceUseInfos: %v", dataToJson(deviceUseInfos))
	return
}

// showEnergy 展示设备消耗的能量
func ShowEnergy(dvIdList []int) {
	for _, device := range dvIdList {
		power, counterResolution, _, err := rsmiDevEnergyCountGet(device)
		if err != nil {
			glog.Errorf("Error getting energy count for device %d: %v\n", device, err)
			continue
		}
		fmt.Printf("Device %d Energy counter: %d\n", device, power)
		fmt.Printf("Device %d Accumulated Energy (uJ): %.2f\n", device, float64(power)*float64(counterResolution))
	}
}

// 设备内存信息
func ShowMemInfo(dvIdList []int, memTypes []string) {
	var returnTypes []string

	if len(memTypes) == 1 && memTypes[0] == "all" {
		returnTypes = memoryTypeL
	} else {
		for _, memType := range memTypes {
			if contains(memoryTypeL, memType) {
				returnTypes = append(returnTypes, memType)
			} else {
				log.Printf("Invalid memory type: %s", memType)
				return
			}
		}
	}

	fmt.Println(" Memory Usage (Bytes) ")
	for _, device := range dvIdList {
		for _, mem := range returnTypes {
			memInfoUsed, memInfoTotal, err := MemInfo(device, mem)
			if err != nil {
				log.Printf("Error getting memory info for device %d: %v", device, err)
				continue
			}
			fmt.Println("device ", device, fmt.Sprintf("%s Total Memory (B)", mem), memInfoTotal)
			fmt.Println("device ", device, fmt.Sprintf("%s Total Used Memory (B)", mem), memInfoUsed)
		}
	}
	fmt.Println("End of Memory Usage")
}

// 设备内存使用情况
func ShowMemUse(dvIdList []int) {
	fmt.Println("Current Memory Use")
	for _, device := range dvIdList {
		busyPercent, err := rsmiDevMemoryBusyPercentGet(device)
		if err == nil {
			fmt.Println("device: ", device, "GPU memory use (%)", busyPercent)
		}
		typeName := "Memory Activity"
		utilCounters, err := GetCoarseGrainUtil(device, &typeName)
		if err == nil {
			for _, utCounter := range utilCounters {
				fmt.Println("device: ", device, utilizationCounterName[utCounter.Type], utCounter.Value)
			}
		} else {
			glog.Errorf("Device %d: Failed to get coarse grain util counters: %v", device, err)
		}
	}
}

// 显示设备供应商信息
func ShowMemVendor(dvIdList []int) (deviceMemVendorInfos []DeviceMemVendorInfo, err error) {
	for _, device := range dvIdList {
		vendor, err := rsmiDevVramVendorGet(device)
		if err == nil {
			glog.Infof("device:%v  GPU memory vendor:%v", device, vendor)
			deviceMemVendorInfos = append(deviceMemVendorInfos, DeviceMemVendorInfo{DeviceID: device, Vendor: vendor})
		} else {
			glog.Warning("GPU memory vendor missing or not supported")
			deviceMemVendorInfos = append(deviceMemVendorInfos, DeviceMemVendorInfo{DeviceID: device, Vendor: ""})
		}
	}
	glog.Infof("GPU memory vendor: %v", dataToJson(deviceMemVendorInfos))
	return
}

// PCIe带宽使用情况
func ShowPcieBw(dvIdList []int) (pcieBandwidthInfos []PcieBandwidthInfo, err error) {
	for _, device := range dvIdList {
		sent, received, maxPktSz, err := rsmiDevPciThroughputGet(device)
		if err == nil {
			// 计算带宽
			bw := ((float64(received) + float64(sent)) * float64(maxPktSz)) / 1024.0 / 1024.0
			bwstr := fmt.Sprintf("%.3f", bw)
			glog.Infof("device:%v Estimated maximum PCIe bandwidth over the last second (MB/s):%v", device, bwstr)
			pcieBandwidthInfos = append(pcieBandwidthInfos, PcieBandwidthInfo{DeviceID: device, Sent: sent, Received: received, MaxPktSz: maxPktSz, Bw: bw})
		} else {
			glog.Warning("GPU PCIe bandwidth usage not supported")
			pcieBandwidthInfos = append(pcieBandwidthInfos, PcieBandwidthInfo{DeviceID: device, Sent: 0, Received: 0, MaxPktSz: 0, Bw: 0})
		}
	}
	glog.Infof("pcieBandwidthInfos:%v", dataToJson(pcieBandwidthInfos))
	return
}

// 设备PCIe重放计数
func ShowPcieReplayCount(dvIdList []int) (pcieReplayCountInfos []PcieReplayCountInfo, err error) {
	for _, device := range dvIdList {
		count, err := rsmiDevPciReplayCounterGet(device)
		if err == nil {
			glog.Infof("device:%v PCIe Replay Count:%v", device, count)
			pcieReplayCountInfos = append(pcieReplayCountInfos, PcieReplayCountInfo{DeviceID: device, Count: count})
		} else {
			glog.Warning("GPU PCIe replay count not supported")
			pcieReplayCountInfos = append(pcieReplayCountInfos, PcieReplayCountInfo{DeviceID: device, Count: 0})
		}
	}
	glog.Infof("pcieReplayCountInfos:%v", dataToJson(pcieReplayCountInfos))
	return
}

// 获取进程信息
func ShowPids() {
	fmt.Printf("========== KFD Processes ==========\n")
	dataArray := [][]string{
		{"PID", "PROCESS NAME", "GPU(s)", "VRAM USED", "SDMA USED", "CU OCCUPANCY"},
	}

	pidList, err := PidList()
	if err != nil {
		fmt.Printf("Error getting PID list: %v\n", err)
		fmt.Printf("==========\n")
		return
	}

	if len(pidList) == 0 {
		fmt.Println("No KFD PIDs currently running")
		fmt.Printf("==========\n")
		return
	}

	for _, pidStr := range pidList {
		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			log.Println("Error converting PID:", err)
			continue
		}

		gpuNumber := "UNKNOWN"
		vramUsage := "UNKNOWN"
		sdmaUsage := "UNKNOWN"
		cuOccupancy := "UNKNOWN"

		dvIndices, err := rsmiComputeProcessGpusGet(pid)
		if err == nil {
			gpuNumber = fmt.Sprintf("%d", len(dvIndices))
		} else {
			glog.Warning("Unable to fetch GPU number by PID")
		}

		proc, err := rsmiComputeProcessInfoByPidGet(pid)
		if err == nil {
			vramUsage = fmt.Sprintf("%d", proc.VramUsage)
			sdmaUsage = fmt.Sprintf("%d", proc.SdmaUsage)
			cuOccupancy = fmt.Sprintf("%d", proc.CuOccupancy)
		} else {
			glog.Warning("Unable to fetch process info by PID")
		}

		dataArray = append(dataArray, []string{
			pidStr,
			getProcessName(pid),
			gpuNumber,
			vramUsage,
			sdmaUsage,
			cuOccupancy,
		})
	}

	fmt.Println("KFD process information:")
	print2DArray(dataArray)
	fmt.Printf("==========\n")
}

func getProcessName(pid int) string {
	if pid < 1 {
		log.Println("PID must be greater than 0")
		return "UNKNOWN"
	}

	cmd := exec.Command("sh", "-c", fmt.Sprintf("ps -p %d -o comm=", pid))
	output, err := cmd.Output()
	if err != nil {
		log.Println("Error executing command:", err)
		return "UNKNOWN"
	}

	pName := strings.TrimSpace(string(output))
	if pName == "" {
		pName = "UNKNOWN"
	}

	return pName
}

// 获取设备平均功率
func ShowPower(dvIdList []int) (devicePowerInfos []DevicePowerInfo, err error) {
	fmt.Println("========== Power Consumption ==========")

	for _, device := range dvIdList {
		power, err := Power(device)
		if err != nil {
			glog.Errorf("device:%v Unable to get Average Graphics Package Power Consumption", device)
			devicePowerInfos = append(devicePowerInfos, DevicePowerInfo{DeviceID: device, Power: -1})
			continue
		}
		if power != 0 {
			fmt.Println("device:", device, "Average Graphics Package Power (W)", fmt.Sprintf("%d", power))
			devicePowerInfos = append(devicePowerInfos, DevicePowerInfo{DeviceID: device, Power: power})
		} else {
			glog.Errorf("device:%v Unable to get Average Graphics Package Power Consumption", device)
			devicePowerInfos = append(devicePowerInfos, DevicePowerInfo{DeviceID: device, Power: -1})
		}
	}
	glog.Infof("devicePowerInfos:%v", dataToJson(devicePowerInfos))
	return
}

// 显示当前GPU内存时钟频率和电压
func ShowPowerPlayTable(dvIdList []int) (devicePowerPlayInfos []DevicePowerPlayInfo, err error) {
	fmt.Println("========== GPU Memory clock frequencies and voltages ==========")
	for _, device := range dvIdList {
		odv, err := rsmiDevOdVoltInfoGet(device)
		if err != nil {
			log.Printf("Error retrieving voltage info for device %d: %v\n", device, err)
			continue
		}

		od_sclk := []string{
			fmt.Sprintf("0: %dMhz", odv.CurrSclkRange.LowerBound/1000000),
			fmt.Sprintf("1: %dMhz", odv.CurrSclkRange.UpperBound/1000000),
		}

		od_mclk := fmt.Sprintf("1: %dMhz", odv.CurrMclkRange.UpperBound/1000000)

		od_vddc_curve := make([]string, 3)
		for position := 0; position < 3; position++ {
			od_vddc_curve[position] = fmt.Sprintf("%d: %dMhz %dmV", position,
				odv.Curve.VcPoints[position].Frequency/1000000,
				odv.Curve.VcPoints[position].Voltage)
		}

		od_range := []string{
			fmt.Sprintf("SCLK: %dMhz %dMhz", odv.SclkFreqLimits.LowerBound/1000000, odv.SclkFreqLimits.UpperBound/1000000),
			fmt.Sprintf("MCLK: %dMhz %dMhz", odv.MclkFreqLimits.LowerBound/1000000, odv.MclkFreqLimits.UpperBound/1000000),
		}

		for position := 0; position < 3; position++ {
			od_range = append(od_range, fmt.Sprintf("VDDC_CURVE_SCLK[%d]: %dMhz", position, odv.Curve.VcPoints[position].Frequency/1000000))
			od_range = append(od_range, fmt.Sprintf("VDDC_CURVE_VOLT[%d]: %dmV", position, odv.Curve.VcPoints[position].Voltage))
		}

		powerPlayInfo := DevicePowerPlayInfo{
			DeviceID:      device,
			OD_SCLK:       od_sclk,
			OD_MCLK:       od_mclk,
			OD_VDDC_CURVE: od_vddc_curve,
			OD_RANGE:      od_range,
		}

		devicePowerPlayInfos = append(devicePowerPlayInfos, powerPlayInfo)
	}

	fmt.Println("===============================================================")
	return
}

// 显示设备列表中所请求的产品名称
func ShowProductName(dvIdList []int) (deviceProductInfos []DeviceproductInfo, err error) {
	fmt.Println("========== Product Info ==========")
	for _, device := range dvIdList {
		deviceProductInfo := DeviceproductInfo{DeviceID: device}

		// Retrieve card vendor
		vendor, err := rsmiDevVendorNameGet(device)
		if err != nil {
			log.Printf("Incompatible device. GPU[%d]: Expected vendor name: Advanced Micro Devices, Inc. [AMD/ATI]\nGPU[%d]: Actual vendor name: %s\n", device, device, vendor)
			continue
		}
		deviceProductInfo.CardVendor = vendor

		// Retrieve the device series
		series, err := rsmiDevNameGet(device)
		if err == nil {
			deviceProductInfo.CardSeries = series
			fmt.Printf("GPU[%d] Card series: %s\n", device, series)
		}

		// Retrieve the device model
		model, err := rsmiDevSubsystemNameGet(device)
		if err == nil {
			deviceProductInfo.CardModel = model
			fmt.Printf("GPU[%d] Card model: %s\n", device, model)
		}

		fmt.Printf("GPU[%d] Card vendor: %s\n", device, vendor)

		// Retrieve the device SKU
		vbios, err := rsmiDevVbiosVersionGet(device, 256)
		if err == nil {
			deviceProductInfo.CardSKU = vbios
			fmt.Printf("GPU[%d] Card SKU: %s\n", device, deviceProductInfo.CardSKU)
		}

		deviceProductInfos = append(deviceProductInfos, deviceProductInfo)

	}

	fmt.Println("==================================")
	glog.Infof("deviceProductInfos:%v", dataToJson(deviceProductInfos))
	return
}

// 可用电源配置文件
func ShowProfile(dvIdList []int) (deviceProfiles []DeviceProfile, err error) {
	fmt.Println(" Show Power Profiles ")

	for _, device := range dvIdList {
		status, err := rsmiDevPowerProfilePresetsGet(device, 0)
		if err != nil {
			log.Printf("Error getting power profile presets: %v", err)
			continue
		}

		binaryMaskString := fmt.Sprintf("%07b", status.AvailableProfiles)
		bitMaskPosition := 0
		profileNumber := 0
		var profiles []string

		for bitMaskPosition < 7 {
			if binaryMaskString[6-bitMaskPosition] == '1' {
				profileNumber++
				var profileInfo string
				if 1<<bitMaskPosition == int(status.Current) {
					profileInfo = fmt.Sprintf("%d. Available power profile (#%d of 7): %s*", profileNumber, bitMaskPosition+1, profileString(1<<bitMaskPosition))
				} else {
					profileInfo = fmt.Sprintf("%d. Available power profile (#%d of 7): %s", profileNumber, bitMaskPosition+1, profileString(1<<bitMaskPosition))
				}
				profiles = append(profiles, profileInfo)
			}
			bitMaskPosition++
		}

		deviceProfiles = append(deviceProfiles, DeviceProfile{DeviceID: device, Profiles: profiles})
	}
	glog.Infof("deviceProfiles: %v", dataToJson(deviceProfiles))
	return
}

// 电流或电压范围
func ShowRange(dvIdList []int, rangeType string) {
	if rangeType != "sclk" && rangeType != "mclk" && rangeType != "voltage" {
		fmt.Println(0, fmt.Sprintf("Invalid range identifier %s", rangeType))
		return
	}

	fmt.Println(fmt.Sprintf(" Show Valid %s Range ", rangeType))

	for _, device := range dvIdList {
		odvf, err := rsmiDevOdVoltInfoGet(device)
		if err != nil {
			log.Printf("Error getting OD volt info: %v", err)
			fmt.Println(device, fmt.Sprintf("Unable to display %s range", rangeType))
			continue
		}
		switch rangeType {
		case "sclk":
			fmt.Println(device, fmt.Sprintf("Valid sclk range: %dMhz - %dMhz",
				odvf.CurrSclkRange.LowerBound/1000000, odvf.CurrSclkRange.UpperBound/1000000))
		case "mclk":
			fmt.Println(device, fmt.Sprintf("Valid mclk range: %dMhz - %dMhz",
				odvf.CurrMclkRange.LowerBound/1000000, odvf.CurrMclkRange.UpperBound/1000000))
		case "voltage":
			numRegions, regions, err := rsmiDevOdVoltCurveRegionsGet(device)
			if err != nil {
				log.Printf("Error getting OD volt curve regions: %v", err)
				fmt.Println(device, fmt.Sprintf("Unable to display %s range", rangeType))
				continue
			}
			for i := 0; i < numRegions; i++ {
				fmt.Println(device, fmt.Sprintf("Region %d: Valid voltage range: %dmV - %dmV",
					i, regions[i], regions[i].VoltRange.UpperBound))
			}
		}
	}

	fmt.Println(" End of Range Display ")
}

// 显示设备列表中指定类型的退役页
func ShowRetiredPages(dvIdList []int, retiredType string) {
	fmt.Println(" Pages Info ")
	if retiredType == "" {
		retiredType = "all"
	}

	for _, device := range dvIdList {
		_, records, err := rsmiDevMemoryReservedPagesGet(device)
		if err != nil {
			log.Printf("Unable to retrieve reserved page info for device %d: %v", device, err)
			continue
		}

		var data [][]string
		for _, rec := range records {
			status := MemoryPageStatus[rec.Status]
			if status == retiredType || retiredType == "all" {
				data = append(data, []string{
					fmt.Sprintf("0x%X", rec.PageAddress),
					fmt.Sprintf("0x%X", rec.PageSize),
					status,
				})
			}
		}

		if len(data) > 0 {
			printTableLog([]string{"Page address", "Page size", "Status"}, data, device, retiredType+" PAGES INFO")
		}
	}
	fmt.Println(" Pages Info ")
}

// 设备序列号
func ShowSerialNumber(dvIdList []int) (deviceSerialInfos []DeviceSerialInfo, err error) {
	fmt.Println("----- Serial Number -----")
	for _, device := range dvIdList {
		serialNumber, err := rsmiDevSerialNumberGet(device)
		deviceSerialInfo := DeviceSerialInfo{
			DeviceID: device,
		}
		if err == nil && serialNumber != "" {
			deviceSerialInfo.SerialNumber = serialNumber
		} else {
			deviceSerialInfo.SerialNumber = "N/A"
		}
		deviceSerialInfos = append(deviceSerialInfos, deviceSerialInfo)
		fmt.Printf("Device %d - Serial Number: %s\n", device, deviceSerialInfo.SerialNumber)
	}
	fmt.Println("------------------------")
	glog.Infof("deviceSerialInfos:%v", dataToJson(deviceSerialInfos))
	return
}

// 唯一设备ID
func ShowUId(dvIdList []int) (deviceUIdInfos []DeviceUIdInfo, err error) {
	fmt.Println("----- Unique ID -----")
	for _, device := range dvIdList {
		uniqueId, err := rsmiDevUniqueIdGet(device)
		deviceUIdInfo := DeviceUIdInfo{
			DeviceID: device,
		}
		if err == nil && uniqueId != 0 {
			deviceUIdInfo.UId = fmt.Sprintf("0x%x", uniqueId)
		} else {
			deviceUIdInfo.UId = "N/A"
		}
		deviceUIdInfos = append(deviceUIdInfos, deviceUIdInfo)
		fmt.Printf("Device %d - Unique ID: %s\n", device, deviceUIdInfo.UId)
	}
	fmt.Println("---------------------")
	return
}

// showVbiosVersion 打印并返回设备的VBIOS版本信息
func ShowVbiosVersion(dvIdList []int) (deviceVBIOSInfos []DeviceVBIOSInfo, err error) {
	fmt.Println("----- VBIOS -----")
	for _, device := range dvIdList {
		vbios, err := VbiosVersion(device)
		if err != nil {
			fmt.Printf("Error fetching VBIOS version for device %d: %v\n", device, err)
			deviceVBIOSInfos = append(deviceVBIOSInfos, DeviceVBIOSInfo{
				DeviceID: device,
				VBIOS:    "Error",
			})
		} else {
			fmt.Printf("Device %d VBIOS version: %s\n", device, vbios)
			deviceVBIOSInfos = append(deviceVBIOSInfos, DeviceVBIOSInfo{
				DeviceID: device,
				VBIOS:    vbios,
			})
		}
	}
	fmt.Println("---------------")
	glog.Infof("deviceVBIOSInfos:%v", dataToJson(deviceVBIOSInfos))
	return
}

// showEvents 显示设备的事件
func ShowEvents(dvIdList []int, eventTypes []string) {
	fmt.Println("----- Show Events -----")
	fmt.Println("Press 'q' or 'ctrl + c' to quit")

	var eventTypeList []string
	for _, event := range eventTypes { // 清理列表中的错误值
		cleanEvent := strings.ReplaceAll(event, ",", "")
		if contains(notificationTypeNames, strings.ToUpper(cleanEvent)) {
			eventTypeList = append(eventTypeList, strings.ToUpper(cleanEvent))
		} else {
			fmt.Printf("Ignoring unrecognized event type %s\n", cleanEvent)
		}
	}

	if len(eventTypeList) == 0 {
		eventTypeList = notificationTypeNames
	}

	var wg sync.WaitGroup
	for _, device := range dvIdList {
		wg.Add(1)
		go func(device int) {
			defer wg.Done()
			printEventList(device, 1000, eventTypeList)
		}(device)
		time.Sleep(250 * time.Millisecond)
	}

	go func() {
		var input string
		for {
			fmt.Scanln(&input)
			if input == "q" {
				for _, device := range dvIdList {
					if err := rsmiEventNotificationStop(device); err != nil {
						fmt.Printf("GPU[%d]: Unable to end event notifications: %v\n", device, err)
					}
				}
				break
			}
		}
	}()

	wg.Wait()
}

// 当前电压信息
func ShowVoltage(dvIdList []int) (deviceVoltageInfos []DeviceVoltageInfo, err error) {
	for _, device := range dvIdList {
		// 默认电压类型和度量标准
		vtype := RSMI_VOLT_TYPE_FIRST
		met := RSMI_VOLT_CURRENT //
		voltage := rsmiDevVoltMetricGet(device, vtype, met)
		if voltage != 0 {
			fmt.Printf("Device %d: Voltage (mV) = %d\n", device, voltage)
			deviceVoltageInfos = append(deviceVoltageInfos, DeviceVoltageInfo{
				DeviceID: device,
				Voltage:  voltage,
			})
		} else {
			log.Printf("GPU %d voltage not supported\n", device)
		}
	}
	glog.Infof("deviceVoltageInfos:%v", dataToJson(deviceVoltageInfos))
	return
}

// 电压曲线点
func ShowVoltageCurve(dvIdList []int) {
	fmt.Println("------------ Voltage Curve Points ------------")
	for _, device := range dvIdList {
		odv, err := rsmiDevOdVoltInfoGet(device)
		if err != nil {
			log.Printf("GPU %d: Voltage Curve is not supported: %v\n", device, err)
			continue
		}

		for position, point := range odv.Curve.VcPoints {
			fmt.Printf("Device %d: Voltage point %d: %d MHz %d mV\n", device, position, point.Frequency/1000000, point.Voltage)
		}
	}
	fmt.Println("----------------------------------------------")
}

// XGMI错误状态
func ShowXgmiErr(dvIdList []int, printJSON bool) {
	fmt.Println("------------ XGMI Error Status ------------")
	for _, device := range dvIdList {
		status, err := rsmiDevXGMIErrorStatus(device)
		if err != nil {
			log.Printf("Error retrieving XGMI status for device %d: %v\n", device, err)
			continue
		}

		var desc string
		switch status {
		case RSMIXGMIStatusNoErrors:
			desc = "No errors detected since last read"
		case RSMIXGMIStatusError:
			desc = "Single error detected since last read"
		case RSMIXGMIStatusMultipleErrors:
			desc = "Multiple errors detected since last read"
		default:
			log.Printf("Invalid return value from xgmi_error for device %d\n", device)
			continue
		}

		if printJSON {
			fmt.Printf("Device %d: XGMI Error count: %d\n", device, status)
		} else {
			fmt.Printf("Device %d: XGMI Error count: %d (%s)\n", device, status, desc)
		}
	}
	fmt.Println("-------------------------------------------")
}

// 硬件拓扑信息
func ShowWeightTopology(dvIdList []int, printJSON bool) {
	// 初始化矩阵存储设备间的权重
	gpuLinksWeight := make([][]int64, len(dvIdList))
	for i := range gpuLinksWeight {
		gpuLinksWeight[i] = make([]int64, len(dvIdList))
	}

	fmt.Println("------------ Weight between two GPUs ------------")
	for _, srcDevice := range dvIdList {
		for _, destDevice := range dvIdList {
			if srcDevice == destDevice {
				gpuLinksWeight[srcDevice][destDevice] = 0
			} else {
				weight, err := rsmiTopoGetLinkWeight(srcDevice, destDevice)
				if err != nil {
					log.Printf("Cannot read Link Weight between device %d and %d: %v\n", srcDevice, destDevice, err)
					continue
				}
				gpuLinksWeight[srcDevice][destDevice] = weight
			}
		}
	}

	if printJSON {
		formatMatrixToJSON(dvIdList, gpuLinksWeight, "(Topology) Weight between DRM devices %d and %d")
		return
	}

	// 打印矩阵表格
	printTableRow("", "      ")
	for _, row := range dvIdList {
		printTableRow("%-12s", fmt.Sprintf("GPU%d", row))
	}
	fmt.Println()
	for _, gpu1 := range dvIdList {
		printTableRow("%-6s", fmt.Sprintf("GPU%d", gpu1))
		for _, gpu2 := range dvIdList {
			if gpu1 == gpu2 {
				printTableRow("%-12s", "0")
			} else {
				printTableRow("%-12d", gpuLinksWeight[gpu1][gpu2])
			}
		}
		fmt.Println()
	}
	fmt.Println("-------------------------------------------------")
}

// 基于跳数显示硬件拓扑信息
func ShowHopsTopology(dvIdList []int, printJSON bool) {
	// 初始化矩阵存储设备间的跳数
	gpuLinksHops := make([][]int64, len(dvIdList))
	for i := range gpuLinksHops {
		gpuLinksHops[i] = make([]int64, len(dvIdList))
	}

	fmt.Println("------------ Hops between two GPUs ------------")
	for _, srcDevice := range dvIdList {
		for _, destDevice := range dvIdList {
			if srcDevice == destDevice {
				gpuLinksHops[srcDevice][destDevice] = 0
			} else {
				hops, _, err := rsmiTopoGetLinkType(srcDevice, destDevice)
				if err != nil {
					log.Printf("Cannot read Link Hops between device %d and %d: %v\n", srcDevice, destDevice, err)
					continue
				}
				gpuLinksHops[srcDevice][destDevice] = hops
			}
		}
	}

	if printJSON {
		formatMatrixToJSON(dvIdList, gpuLinksHops, "(Topology) Hops between DRM devices %d and %d")
		return
	}

	// 打印矩阵表格
	printTableRow("", "      ")
	for _, row := range dvIdList {
		printTableRow("%-12s", fmt.Sprintf("GPU%d", row))
	}
	fmt.Println()
	for _, gpu1 := range dvIdList {
		printTableRow("%-6s", fmt.Sprintf("GPU%d", gpu1))
		for _, gpu2 := range dvIdList {
			if gpu1 == gpu2 {
				printTableRow("%-12s", "0")
			} else {
				printTableRow("%-12d", gpuLinksHops[gpu1][gpu2])
			}
		}
		fmt.Println()
	}
	fmt.Println("-------------------------------------------------")
}

// 基于链接类型的硬件拓扑信息
func ShowTypeTopology(dvIdList []int, printJSON bool) {
	// 初始化矩阵存储设备间的链接类型
	gpuLinksType := make([][]string, len(dvIdList))
	for i := range gpuLinksType {
		gpuLinksType[i] = make([]string, len(dvIdList))
	}

	fmt.Println("------------ Link Type between two GPUs ------------")
	for _, srcDevice := range dvIdList {
		for _, destDevice := range dvIdList {
			if srcDevice == destDevice {
				gpuLinksType[srcDevice][destDevice] = "0"
			} else {
				_, linkType, err := rsmiTopoGetLinkType(srcDevice, destDevice)
				if err != nil {
					log.Printf("Cannot read Link Type between device %d and %d: %v\n", srcDevice, destDevice, err)
					continue
				}
				switch linkType {
				case 1:
					gpuLinksType[srcDevice][destDevice] = LinkTypePCIE
				case 2:
					gpuLinksType[srcDevice][destDevice] = LinkTypeXGMI
				default:
					gpuLinksType[srcDevice][destDevice] = LinkTypeUnknown
				}
			}
		}
	}

	if printJSON {
		formatMatrixToStrJSON(dvIdList, gpuLinksType, "(Topology) Link type between DRM devices %d and %d")
		return
	}

	// 打印矩阵表格
	printTableRow("", "      ")
	for _, row := range dvIdList {
		printTableRow("%-12s", fmt.Sprintf("GPU%d", row))
	}
	fmt.Println()
	for _, gpu1 := range dvIdList {
		printTableRow("%-6s", fmt.Sprintf("GPU%d", gpu1))
		for _, gpu2 := range dvIdList {
			if gpu1 == gpu2 {
				printTableRow("%-12s", "0")
			} else {
				printTableRow("%-12s", gpuLinksType[gpu1][gpu2])
			}
		}
		fmt.Println()
	}
	fmt.Println("----------------------------------------------------")
}

// numa节点HW拓扑信息
func ShowNumaTopology(dvIdList []int) {
	fmt.Println("---------- Numa Nodes ----------")

	for _, device := range dvIdList {
		// 获取 NUMA 节点编号
		numaNode, err := rsmiTopoGetNumaBodeBumber(device)
		if err == nil {
			fmt.Printf("Device %d: (Topology) Numa Node: %d\n", device, numaNode)
		} else {
			glog.Errorf("device:%v Cannot read Numa Node", device)
		}

		// 获取 NUMA 关联信息
		numaAffinity, err := rsmiTopoNumaAffinityGet(device)
		if err == nil {
			fmt.Println(device, "(Topology) Numa Affinity", numaAffinity)
		} else {
			glog.Errorf("device:%v Cannot read Numa Affinity", device)
		}
	}
}

// 显示硬件拓扑信息,包括权重、跳数、链接类型以及NUMA节点信息
func ShowHwTopology(dvIdList []int) {
	ShowWeightTopology(dvIdList, true)

	ShowHopsTopology(dvIdList, true)

	ShowTypeTopology(dvIdList, true)

	ShowNumaTopology(dvIdList)
}

/*************************************VDCU******************************************/
// 设备数量
func DeviceCount() (count int, err error) {
	return dmiGetDeviceCount()
}

// 虚拟设备信息
func VDeviceSingleInfo(dvInd int) (vDeviceInfo DMIVDeviceInfo, err error) {
	return dmiGetVDeviceInfo(dvInd)
}

// 虚拟设备数量
func VDeviceCount() (count int, err error) { return dmiGetVDeviceCount() }

// 指定物理设备剩余的CU和内存
//func DeviceRemainingInfo(dvInd int) (cus, memories uintptr, err error) {
//	return dmiGetDeviceRemainingInfo(dvInd)
//}

// CreateVDevices 创建指定数量的虚拟设备
// @Summary 创建虚拟设备
// @Description 在指定的物理设备上创建指定数量的虚拟设备。
// @Tags 虚拟设备
// @Param dvInd query int true "物理设备的索引"
// @Param vDevCount query int true "要创建的虚拟设备数量"
// @Param vDevCUs query []int true "每个虚拟设备的计算单元数量"
// @Param vDevMemSize query []int true "每个虚拟设备的内存大小"
// @Success 200 {string} string "虚拟设备创建成功"
// @Failure 400 {string} string "创建虚拟设备失败"
// @Router /CreateVDevices [post]
func CreateVDevices(dvInd int, vDevCount int, vDevCUs []int, vDevMemSize []int) (err error) {
	return dmiCreateVDevices(dvInd, vDevCount, vDevCUs, vDevMemSize)
}

// DestroyVDevice 销毁指定物理设备上的所有虚拟设备
// @Summary 销毁所有虚拟设备
// @Description 销毁指定物理设备上的所有虚拟设备。
// @Tags 虚拟设备
// @Param dvInd query int true "物理设备的索引"
// @Success 200 {string} string "虚拟设备销毁成功"
// @Failure 400 {string} string "虚拟设备销毁失败"
// @Router /DestroyVDevice [delete]
func DestroyVDevice(dvInd int) (err error) {
	return dmiDestroyVDevices(dvInd)
}

// DestroySingleVDevice 销毁指定虚拟设备
// @Summary 销毁单个虚拟设备
// @Description 销毁指定索引的虚拟设备。
// @Tags 虚拟设备
// @Param vDvInd query int true "虚拟设备的索引"
// @Success 200 {string} string "虚拟设备销毁成功"
// @Failure 400 {string} string "虚拟设备销毁失败"
// @Router /DestroySingleVDevice [delete]
func DestroySingleVDevice(vDvInd int) (err error) {
	return dmiDestroySingleVDevice(vDvInd)
}

// UpdateSingleVDevice 更新指定设备资源大小
// @Summary 更新虚拟设备资源
// @Description 更新指定虚拟设备的计算单元和内存大小。如果 vDevCUs 或 vDevMemSize 为 -1，则对应的资源不更改。
// @Tags 虚拟设备
// @Param vDvInd query int true "虚拟设备的索引"
// @Param vDevCUs query int true "更新后的计算单元数量"
// @Param vDevMemSize query int true "更新后的内存大小"
// @Success 200 {string} string "虚拟设备更新成功"
// @Failure 400 {string} string "虚拟设备更新失败"
// @Router /UpdateSingleVDevice [put]
func UpdateSingleVDevice(vDvInd int, vDevCUs int, vDevMemSize int) (err error) {
	return dmiUpdateSingleVDevice(vDvInd, vDevCUs, vDevMemSize)
}

// 启动虚拟设备
func StartVDevice(vDvInd int) (err error) {
	return dmiStartVDevice(vDvInd)
}

// 停止虚拟设备
func StopVDevice(vDvInd int) (err error) {
	return dmiStopVDevice(vDvInd)
}

// 设置虚拟机加密状态 status为true，则开启加密虚拟机，否则关闭
func SetEncryptionVMStatus(status bool) (err error) {
	return dmiSetEncryptionVMStatus(status)
}

// 获取加密虚拟机状态
func EncryptionVMStatus() (status bool, err error) {
	return dmiGetEncryptionVMStatus()
}

// 打印事件列表方法
func PrintEventList(device int, delay int, eventList []string) {
	printEventList(device, delay, eventList)
}

package dcgm

import "C"
import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
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
func Init() (err error) {
	devCount := listFilesInDevDri()
	glog.Infof("devCount:%v", devCount)
	maxRetries := 12                   // 最大重试次数
	retryCount := 0                    // 记录连续返回相同设备数量的次数
	lastNumDevices := -1               // 记录上一次获取的设备数量
	restartTimeout := 10 * time.Second // 每次重试等待10秒
	initFailCount := 0                 // rsmiInit 连续失败的计数
	maxInitFails := 6                  // 连续失败最大次数
	for {
		err = rsmiInit() // 初始化rsmi
		if err == nil {
			ShutDown()
			for retryCount < maxRetries {
				rsmiInit()
				numDevices, _ := NumMonitorDevices() // 获取GPU设备数量
				if numDevices == devCount {
					glog.Infof("DCU initialization is complete:%v", numDevices)
					return nil // 数量相等，初始化成功，结束函数
				} else {
					if numDevices == lastNumDevices {
						retryCount++ // 记录连续返回相同设备数量的次数
					} else {
						retryCount = 0 // 数量变化时重置计数
					}

					glog.Infof("retryCount:%v", retryCount)
					if retryCount >= maxRetries {
						glog.Infof("设备数量连续 %d 次相同但与 devCount 不相等，初始化失败", maxRetries)
						return
					}
					lastNumDevices = numDevices // 更新记录的设备数量
					ShutDown()                  // 数量不相等，执行关机操作
				}
				time.Sleep(restartTimeout) // 等待10秒
			}
		} else {
			initFailCount++ // 初始化失败，计数加一
			glog.Infof("初始化失败: %v. 10秒后重试...\n", err)

			if initFailCount >= maxInitFails {
				glog.Errorf("rsmiInit 连续 %d 次失败，终止初始化: %v", maxInitFails, err)
				return err // 连续6次失败，返回错误信息
			}
		}
		time.Sleep(restartTimeout) // 等待10秒后再次重试
	}
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
// @Success 200 {string} brand "设备品牌名称"
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
// @Success 200 {object} RSMIPcieBandwidth "PCIe 带宽列表"
// @Failure 400 {object} error "请求错误"
// @Failure 404 {object} error "设备未找到"
// @Router /DevPciBandwidth [get]
func DevPciBandwidth(dvInd int) (rsmiPcieBandwidth RSMIPcieBandwidth, err error) {
	return rsmiDevPciBandwidthGet(dvInd)
}

func DevPciBandwidthSet(dvInd int, bwBitmask int64) (err error) {
	return rsmiDevPciBandwidthSet(dvInd, bwBitmask)
}

// @Summary 获取内存使用百分比
// @Description 根据设备 ID 获取设备内存的CollectDeviceMetrics使用百分比。
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

// DevGpuMetricsInfo 获取 GPU 度量信息
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

func DevPowerCapRange(dvInd int, senserId int) (max, min int64, err error) {
	return rsmiDevPowerCapRangeGet(dvInd, senserId)
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
	var wg sync.WaitGroup
	monitorInfos = make([]MonitorInfo, numMonitorDevices)
	deviceResults := make(chan MonitorInfo, numMonitorDevices) // Create a channel to collect results

	for i := 0; i < numMonitorDevices; i++ {
		wg.Add(1)
		go func(deviceIndex int) {
			defer wg.Done()

			var wgDevice sync.WaitGroup
			var muDevice sync.Mutex
			monitorInfo := MonitorInfo{MinorNumber: deviceIndex}

			// Collect PCI ID
			wgDevice.Add(1)
			go func() {
				defer wgDevice.Done()
				bdfid, err := rsmiDevPciIdGet(deviceIndex)
				if err != nil {
					glog.Errorf("Failed to get PCI ID for device %d: %v", deviceIndex, err)
					return
				}
				domain := (bdfid >> 32) & 0xffffffff
				bus := (bdfid >> 8) & 0xff
				dev := (bdfid >> 3) & 0x1f
				function := bdfid & 0x7
				pciBusNumber := fmt.Sprintf("%04x:%02x:%02x.%x", domain, bus, dev, function)
				muDevice.Lock()
				monitorInfo.PciBusNumber = pciBusNumber
				muDevice.Unlock()
			}()

			// Collect Device Serial Number
			wgDevice.Add(1)
			go func() {
				defer wgDevice.Done()
				deviceId, _ := rsmiDevSerialNumberGet(deviceIndex)
				muDevice.Lock()
				monitorInfo.DeviceId = deviceId
				muDevice.Unlock()
			}()

			// Collect Device Type ID
			wgDevice.Add(1)
			go func() {
				defer wgDevice.Done()
				devTypeId, _ := rsmiDevIdGet(deviceIndex)
				devTypeName := type2name[fmt.Sprintf("%x", devTypeId)]
				muDevice.Lock()
				monitorInfo.SubSystemName = devTypeName
				muDevice.Unlock()
			}()

			// Collect Temperature
			wgDevice.Add(1)
			go func() {
				defer wgDevice.Done()
				temperature, _ := rsmiDevTempMetricGet(deviceIndex, 0, RSMI_TEMP_CURRENT)
				t, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(temperature)/1000.0), 64)
				muDevice.Lock()
				monitorInfo.Temperature = t
				muDevice.Unlock()
			}()

			// Collect Power Usage
			wgDevice.Add(1)
			go func() {
				defer wgDevice.Done()
				powerUsage, _ := rsmiDevPowerAveGet(deviceIndex, 0)
				pu, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(powerUsage)/1000000.0), 64)
				muDevice.Lock()
				monitorInfo.PowerUsage = pu
				muDevice.Unlock()
			}()

			// Collect Power Cap
			wgDevice.Add(1)
			go func() {
				defer wgDevice.Done()
				powerCap, _ := rsmiDevPowerCapGet(deviceIndex, 0)
				pc, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(powerCap)/1000000.0), 64)
				muDevice.Lock()
				monitorInfo.PowerCap = pc
				muDevice.Unlock()
			}()

			// Collect Memory Capacity
			wgDevice.Add(1)
			go func() {
				defer wgDevice.Done()
				memoryCap, _ := rsmiDevMemoryTotalGet(deviceIndex, RSMI_MEM_TYPE_FIRST)
				mc, _ := strconv.ParseFloat(fmt.Sprintf("%f", float64(memoryCap)/1.0), 64)
				muDevice.Lock()
				monitorInfo.MemoryCap = mc
				muDevice.Unlock()
			}()

			// Collect Memory Usage
			wgDevice.Add(1)
			go func() {
				defer wgDevice.Done()
				memoryUsed, _ := rsmiDevMemoryUsageGet(deviceIndex, RSMI_MEM_TYPE_FIRST)
				mu, _ := strconv.ParseFloat(fmt.Sprintf("%f", float64(memoryUsed)/1.0), 64)
				muDevice.Lock()
				monitorInfo.MemoryUsed = mu
				muDevice.Unlock()
			}()

			// Collect Utilization Rate
			wgDevice.Add(1)
			go func() {
				defer wgDevice.Done()
				utilizationRate, _ := rsmiDevBusyPercentGet(deviceIndex)
				ur, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(utilizationRate)/1.0), 64)
				muDevice.Lock()
				monitorInfo.UtilizationRate = ur
				muDevice.Unlock()
			}()

			// Collect PCIe Throughput
			wgDevice.Add(1)
			go func() {
				defer wgDevice.Done()
				sent, received, maxPktSz, _ := rsmiDevPciThroughputGet(deviceIndex)
				pcieBwMb, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(received+sent)*float64(maxPktSz)/1024.0/1024.0), 64)
				muDevice.Lock()
				monitorInfo.PcieBwMb = pcieBwMb
				muDevice.Unlock()
			}()

			// Collect GPU Clock Frequencies
			wgDevice.Add(1)
			go func() {
				defer wgDevice.Done()
				clk, _ := rsmiDevGpuClkFreqGet(deviceIndex, RSMI_CLK_TYPE_SYS)
				sclk, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(clk.Frequency[clk.Current])/1000000.0), 64)
				supported := clk.NumSupported
				var sclkFrequency []string
				for i := 0; i < int(supported); i++ {
					freq := fmt.Sprintf("%d", int(clk.Frequency[i]/1000000))
					sclkFrequency = append(sclkFrequency, freq)
				}
				muDevice.Lock()
				monitorInfo.Clk = sclk
				monitorInfo.SclkFrequency = sclkFrequency
				muDevice.Unlock()
			}()

			wgDevice.Add(1)
			go func() {
				defer wgDevice.Done()
				soc, _ := rsmiDevGpuClkFreqGet(deviceIndex, RSMI_CLK_TYPE_SOC)
				socclk, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(soc.Frequency[soc.Current])/1000000.0), 64)
				supported := soc.NumSupported
				var socclkFrequency []string
				for i := 0; i < int(supported); i++ {
					freq := fmt.Sprintf("%d", int(soc.Frequency[i]/1000000))
					socclkFrequency = append(socclkFrequency, freq)
				}

				muDevice.Lock()
				monitorInfo.Socclk = socclk
				monitorInfo.SocclkFrequency = socclkFrequency
				muDevice.Unlock()
			}()

			// Collect Performance Level
			wgDevice.Add(1)
			go func() {
				defer wgDevice.Done()
				perf, err := PerfLevel(deviceIndex)
				if err != nil {
					glog.Errorf("Failed to get performance level for device %d: %v", deviceIndex, err)
					return
				}
				muDevice.Lock()
				monitorInfo.PerfLevel = perf
				muDevice.Unlock()
			}()

			wgDevice.Wait()

			deviceResults <- monitorInfo // Send result to channel
		}(i)
	}

	// Close the channel once all Goroutines are done
	go func() {
		wg.Wait()
		close(deviceResults)
	}()

	// Collect results from channel
	for monitorInfo := range deviceResults {
		monitorInfos[monitorInfo.MinorNumber] = monitorInfo
	}

	glog.Info("monitorInfos: ", dataToJson(monitorInfos))
	return
}

/*func CollectVDeviceMetrics() (devices []PhysicalDeviceInfo, err error) {

}*/

func DevGpuClkFreqSet(dvInd int, clkType RSMIClkType, freqBitmask int64) (err error) {
	return rsmiDevGpuClkFreqSet(dvInd, clkType, freqBitmask)
}

// GetDeviceByDvInd 根据设备的 dvInd 获取物理设备信息
// @Summary 获取物理设备信息
// @Description 根据设备的 dvInd 获取物理设备信息
// @Tags Device
// @Param dvInd path int true "设备的 MinorNumber"
// @Success 200 {object} PhysicalDeviceInfo "返回物理设备信息"
// @Failure 404 {string} string "设备未找到"
// @Failure 500 {string} string "内部服务器错误"
// @Router /GetDeviceByDvInd [get]
func GetDeviceByDvInd(dvInd int) (physicalDeviceInfo PhysicalDeviceInfo, err error) {
	devices, err := AllDeviceInfos()
	if err != nil {
		return physicalDeviceInfo, err
	}
	for _, physicalDevice := range devices {
		if physicalDevice.Device.MinorNumber == dvInd {
			glog.Infof("physicalDevice:%v", dataToJson(physicalDevice))
			return physicalDevice, nil
		}
	}
	return physicalDeviceInfo, fmt.Errorf("device with MinorNumber %d not found", dvInd)
}

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
		devPercent, _ := dmiGetDevBusyPercent(i)

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
		blockInfos, err := EccBlocksInfo(i)
		cus, memories, _ := DeviceRemainingInfo(i)
		device := Device{
			MinorNumber:               i,
			PciBusNumber:              pciBusNumber,
			DeviceId:                  deviceId,
			SubSystemName:             devTypeName,
			Temperature:               t,
			PowerUsage:                pu,
			PowerCap:                  pc,
			MemoryCap:                 mc,
			MemoryUsed:                mu,
			UtilizationRate:           ur,
			PcieBwMb:                  pcieBwMb,
			Clk:                       sclk,
			ComputeUnitCount:          computeUnit,
			MaxVDeviceCount:           maxVDeviceCount,
			Percent:                   devPercent,
			VDeviceCount:              0,
			ComputeUnitRemainingCount: cus,
			MemoryRemaining:           memories,
			BlocksInfos:               blockInfos,
		} // 创建PhysicalDeviceInfo并存入map
		pdi := PhysicalDeviceInfo{
			Device:         device,
			VirtualDevices: []DMIVDeviceInfo{},
		}
		deviceMap[device.MinorNumber] = &pdi
	}

	// 获取虚拟设备数量
	//vDeviceCount, err := dmiGetVDeviceCount()
	vDeviceCount := deviceCount * 4
	if err != nil {
		return nil, err
	}
	// 获取所有虚拟设备信息并关联到对应的物理设备
	for j := 0; j < vDeviceCount; j++ {
		vDeviceInfo, err := dmiGetVDeviceInfo(j)
		glog.Infof("vDeviceInfo error: %v", err)
		if err == nil {
			vDevPercent, _ := dmiGetVDevBusyPercent(j)
			vDeviceInfo.Percent = vDevPercent
			vDeviceInfo.VMinorNumber = j
			// 找到对应的物理设备并将虚拟设备添加到其VirtualDevices中
			if pdi, exists := deviceMap[vDeviceInfo.DeviceID]; exists {
				// 更新虚拟设备的 PciBusNumber，使用物理设备的 pciBusNumber
				vDeviceInfo.PciBusNumber = pdi.Device.PciBusNumber
				// 将虚拟设备添加到物理设备的 VirtualDevices 列表中
				pdi.VirtualDevices = append(pdi.VirtualDevices, vDeviceInfo)
				// 更新物理设备的 VDeviceCount，等于当前虚拟设备的数量
				pdi.Device.VDeviceCount = len(pdi.VirtualDevices)
			}
		}
		if err != nil {
			glog.Errorf("Error getting virtual device info for virtual device %d: %s", j, err)
		}
	}

	//dirPath := "/etc/vdev"
	//// 读取目录中的文件列表
	//files, err := os.ReadDir(dirPath)
	//if err != nil {
	//	glog.Errorf("无法读取目录: %v", err)
	//}
	//
	//// 打印文件数量
	////fmt.Printf("文件数量: %d\n", len(files))
	//
	//// 逐个读取并解析每个文件的内容
	//for _, file := range files {
	//	//glog.Infof("/etc/vdev/file：%v", file)
	//	// 确保是文件而不是子目录
	//	if !file.IsDir() && strings.HasPrefix(file.Name(), "vdev") && strings.HasSuffix(file.Name(), ".conf") {
	//		filePath := filepath.Join(dirPath, file.Name())
	//		config, err := parseConfig(filePath)
	//		if err != nil {
	//			glog.Errorf("无法解析文件 %s: %v", filePath, err)
	//			continue
	//		}
	//		//glog.Infof("文件: %s\n配置: %+v\n", filePath, config)
	//		// 找到对应的物理设备并将虚拟设备添加到其VirtualDevices中
	//		if pdi, exists := deviceMap[config.DeviceID]; exists {
	//			pdi.VirtualDevices = append(pdi.VirtualDevices, *config)
	//			pdi.Device.VDeviceCount = len(pdi.VirtualDevices) // 更新 VDeviceCount
	//		}
	//	}
	//}

	// 将map中的所有PhysicalDeviceInfo转为slice
	for _, pdi := range deviceMap {
		allDevices = append(allDevices, *pdi)
	}
	//for i := range allDevices {
	//	device := &allDevices[i]
	//	var computeUnitCountTotal = 0
	//	var memoryTotal = 0
	//	for _, virtualDevice := range device.VirtualDevices {
	//		computeUnitCountTotal += virtualDevice.ComputeUnitCount
	//		memoryTotal += int(virtualDevice.GlobalMemSize)
	//	}
	//	//glog.Infof("VirtualDevice computeUnitCountTotal:%v  MemoryTotal:%v", computeUnitCountTotal, memoryTotal)
	//	//glog.Infof("VirtualDevice device.Device.ComputeUnitCount:%v", device.Device.ComputeUnitCount)
	//	//glog.Infof("VirtualDevice float64(computeUnitCountTotal):%v ", float64(computeUnitCountTotal))
	//	device.Device.ComputeUnitRemainingCount = uint64(device.Device.ComputeUnitCount - float64(computeUnitCountTotal))
	//	//glog.Infof("device.Device.ComputeUnitRemainingCount:%v", device.Device.ComputeUnitRemainingCount)
	//	device.Device.MemoryRemaining = uint64(device.Device.MemoryCap - float64(memoryTotal))
	//}
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

func EccCount(dvInd int, block RSMIGpuBlock) (errorCount RSMIErrorCount, err error) {
	errorCount, err = rsmiDevEccCountGet(dvInd, block)
	return
}

func EccBlocksInfo(dvInd int) (blocksInfos []BlocksInfo, err error) {
	// 定义所有的RSMIGpuBlock值
	blocks := []RSMIGpuBlock{
		RSMIGpuBlockATHUB,
		RSMIGpuBlockDF,
		RSMIGpuBlockFuse,
		RSMIGpuBlockGFX,
		RSMIGpuBlockHDP,
		RSMIGpuBlockMMHUB,
		RSMIGpuBlockMP0,
		RSMIGpuBlockMP1,
		RSMIGpuBlockPCIEBIF,
		RSMIGpuBlockSDMA,
		RSMIGpuBlockSEM,
		RSMIGpuBlockSMN,
		RSMIGpuBlockUMC,
		RSMIGpuBlockXGMIWAFL,
	}

	// 遍历所有的block，分别调用EccStatus和EccCount
	for _, block := range blocks {
		state, err := EccStatus(dvInd, block)
		if err != nil {
			glog.Errorf("EccStatus 调用错误: block: %v, 错误: %v\n", block, err)
			continue
		}
		//glog.Infof("EccStatus - block: %v, state: %v\n", block, state)

		// 当状态是“ENABLED”时，调用EccCount接口获取错误计数
		if state == "ENABLED" {
			errorCount, err := EccCount(dvInd, block)
			if err != nil {
				glog.Errorf("EccCount 调用错误: block: %v, 错误: %v\n", block, err)
				continue
			}
			//glog.Infof("EccCount - block: %v, CorrectableErr: %v, UncorrectableErr: %v\n", block, errorCount.CorrectableErr, errorCount.UncorrectableErr)
			// 将block信息添加到结果集中
			blocksInfos = append(blocksInfos, BlocksInfo{
				Block: ConvertFromRSMIGpuBlock(block),
				State: state,
				CE:    int64(errorCount.CorrectableErr),
				UE:    int64(errorCount.UncorrectableErr),
			})
		} else {
			// 状态不是ENABLED时，只添加状态信息，不获取错误计数
			blocksInfos = append(blocksInfos, BlocksInfo{
				Block: ConvertFromRSMIGpuBlock(block),
				State: state,
				CE:    0,
				UE:    0,
			})
		}
	}
	//glog.Infof("blocksInfos:%v", dataToJson(blocksInfos))
	return
}

func EccEnabled(dvInd int) (enabledBlocks int64, err error) {
	return rsmiDevEccEnabledGet(dvInd)
}

// 设置设备的性能确定性模式(K100 AI不支持)
func PerfDeterminismMode(dvInd int, clkValue int64) (err error) {
	return rsmiPerfDeterminismModeSet(dvInd, clkValue)
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

// 设置设备超速百分比
func DevOverdriveLevelSet(dvInd, od int) (err error) {
	return rsmiDevOverdriveLevelSet(dvInd, od)
}

// 获取设备的超速百分比
func DevOverdriveLevelGet(dvInd int) (od int, err error) {
	return rsmiDevOverdriveLevelGet(dvInd)
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
// @Summary 获取XGMI错误状态
// @Description 获取指定物理设备的XGMI（高速互连链路）错误状态。
// @Tags XGMI状态
// @Param dvInd query int true "物理设备的索引"
// @Success 200 {integer} int "返回XGMI错误状态码"
// @Failure 400 {string} string "获取XGMI错误状态失败"
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
func SetClockRange(dvIdList []int, clkType string, minvalue string, maxvalue string, autoRespond bool) (failedMessage []FailedMessage) {
	errorMap := make(map[int][]string)
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
			errorMap[device] = append(errorMap[device], err.Error())
			glog.Errorf("Unable to diable performance determinism, device: %v, error: %v", device, err)

		}
	}
	for id, msg := range errorMap {
		failedMessage = append(failedMessage, FailedMessage{ID: id, ErrorMsg: strings.Join(msg, "; ")})
	}
	glog.Infof("SetClockRange failedMessage:%v", failedMessage)
	return
}

// 设置电压曲线
func DevOdVoltInfoSet(dvInd, vPoint, clkValue, voltValue int) (err error) {
	return rsmiDevOdVoltInfoSet(dvInd, vPoint, clkValue, voltValue)
}

// SetPowerPlayTableLevel 设置 PowerPlay 级别
// @Summary 设置设备的 PowerPlay 表级别
// @Description 该函数为设备列表设置 PowerPlay 表级别。它会检查输入值的有效性并相应地调整电压设置。
// @Tags 设备
// @Param dvIdList body []int true "设备 ID 列表"
// @Param clkType query string true "时钟类型（sclk 或 mclk）"
// @Param point query string true "电压点"
// @Param clk query string true "时钟值（以 MHz 为单位）"
// @Param volt query string true "电压值（以 mV 为单位）"
// @Param autoRespond query bool false "自动响应超出规格的警告"
// @Success 200 {string} string "成功设置 PowerPlay 表级别"
// @Failure 400 {string} string "输入无效或无法设置 PowerPlay 表级别"
// @Router /SetPowerPlayTableLevel [post]
func SetPowerPlayTableLevel(dvIdList []int, clkType string, point string, clk string, volt string, autoRespond bool) (failedMessage []FailedMessage) {
	value := fmt.Sprintf("%s %s %s", point, clk, volt)
	_, errPoint := strconv.Atoi(point)
	_, errClk := strconv.Atoi(clk)
	_, errVolt := strconv.Atoi(volt)

	// 创建一个 errorMap 用来记录错误信息
	errorMap := make(map[int][]string)

	if errPoint != nil || errClk != nil || errVolt != nil {
		glog.Infof("Unable to set PowerPlay table level")
		glog.Infof("Non-integer characters are present in %s", value)
		// 这里可以返回错误信息
		failedMessage = append(failedMessage, FailedMessage{ID: -1, ErrorMsg: "Invalid non-integer characters in parameters"})
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
				glog.Infof("device:%v Successfully set voltage point %v to %v(MHz) %v(mV)", device, point, clk, volt)
			} else {
				errorMsg := fmt.Sprintf("Unable to set voltage point %v to %v(MHz) %v(mV)", point, clk, volt)
				glog.Errorf("device:%v %s", device, errorMsg)
				errorMap[device] = append(errorMap[device], errorMsg)
			}
		} else {
			errorMsg := fmt.Sprintf("Unsupported range type %s", clkType)
			glog.Errorf("device:%v Unable to set %s range", device, clkType)
			glog.Infof("Unsupported range type %s", clkType)
			errorMap[device] = append(errorMap[device], errorMsg)
		}
	}

	// 将 errorMap 转换为 failedMessage 列表
	for id, msg := range errorMap {
		failedMessage = append(failedMessage, FailedMessage{ID: id, ErrorMsg: strings.Join(msg, "; ")})
	}

	return
}

// SetClockOverDrive 设置时钟速度为 OverDrive
// @Summary 为设备设置时钟 OverDrive
// @Description 该函数为设备列表设置时钟 OverDrive 级别。它会调整时钟速度，并在需要时确保性能级别设置为手动模式。
// @Tags 设备
// @Param dvIdList body []int true "设备 ID 列表"
// @Param clktype query string true "时钟类型（sclk 或 mclk）"
// @Param value query string true "OverDrive 值，表示为百分比（0-20%）"
// @Param autoRespond query bool false "自动响应超出规格的警告"
// @Success 200 {string} string "成功设置时钟 OverDrive"
// @Failure 400 {string} string "输入无效或无法设置时钟 OverDrive"
// @Router /SetClockOverDrive [post]
func SetClockOverDrive(dvIdList []int, clktype string, value string, autoRespond bool) (failedMessage []FailedMessage) {
	glog.Infof("Set Clock OverDrive Range: 0 to 20%")
	intValue, err := strconv.Atoi(value)
	if err != nil {
		glog.Infof("Unable to set OverDrive level")
		glog.Errorf("%s it is not an integer", value)
		failedMessage = append(failedMessage, FailedMessage{ID: -1, ErrorMsg: "Invalid non-integer value for OverDrive"})
		return
	}

	confirmOutOfSpecWarning(autoRespond)

	for _, device := range dvIdList {
		if intValue < 0 {
			glog.Errorf("Unable to set OverDrive for device: %v", device)
			glog.Infof("Overdrive cannot be less than 0%")
			failedMessage = append(failedMessage, FailedMessage{ID: device, ErrorMsg: "OverDrive cannot be less than 0%"})
			continue
		}
		if intValue > 20 {
			glog.Infof("device:%v, Setting OverDrive to 20%%", device)
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
				failedMessage = append(failedMessage, FailedMessage{ID: device, ErrorMsg: err.Error()})
				continue
			}
		}
		if clktype == "mclk" {
			fsFile := fmt.Sprintf("/sys/class/drm/card%d/device/pp_mclk_od", device)
			if _, err := os.Stat(fsFile); os.IsNotExist(err) {
				glog.Infof("Unable to write to sysfs file")
				glog.Warning("File does not exist: ", fsFile)
				failedMessage = append(failedMessage, FailedMessage{ID: device, ErrorMsg: "Sysfs file does not exist for mclk OverDrive"})
				continue
			}
			f, err := os.OpenFile(fsFile, os.O_WRONLY, 0644)
			if err != nil {
				glog.Infof("Unable to open sysfs file %v", fsFile)
				glog.Warning("IO or OS error")
				failedMessage = append(failedMessage, FailedMessage{ID: device, ErrorMsg: "Unable to open sysfs file for mclk OverDrive"})
				continue
			}
			defer f.Close()
			_, err = f.WriteString(fmt.Sprintf("%v", intValue))
			if err != nil {
				glog.Infof("Unable to write to sysfs file %v", fsFile)
				glog.Warning("IO or OS error")
				failedMessage = append(failedMessage, FailedMessage{ID: device, ErrorMsg: "Unable to write to sysfs file for mclk OverDrive"})
				continue
			}
			glog.Infof("device%v Successfully set %s OverDrive to %d%%", device, clktype, intValue)
		} else if clktype == "sclk" {
			err := rsmiDevOverdriveLevelSet(device, intValue)
			if err == nil {
				glog.Infof("device:%v Successfully set %s OverDrive to %d%%", device, clktype, intValue)
			} else {
				glog.Errorf("device:%v Unable to set %s OverDrive to %d%%", device, clktype, intValue)
				failedMessage = append(failedMessage, FailedMessage{ID: device, ErrorMsg: err.Error()})
			}
		} else {
			glog.Errorf("device:%v Unable to set OverDrive", device)
			glog.Errorf("Unsupported clock type %v", clktype)
			failedMessage = append(failedMessage, FailedMessage{ID: device, ErrorMsg: "Unsupported clock type"})
		}
	}
	return
}

// SetPerfDeterminism 设置时钟频率级别以启用性能确定性
// @Summary 设置时钟频率级别以启用性能确定性
// @Description 根据设备ID列表和给定的时钟频率值，设置设备的性能确定性模式
// @Tags Device
// @Accept  json
// @Produce  json
// @Param dvIdList body []int true "设备 ID 列表"
// @Param clkvalue query string true "时钟频率值"
// @Success 200 {array} FailedMessage
// @Failure 400 {object} FailedMessage
// @Router /SetPerfDeterminism [post]
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
			errorMap[device] = append(errorMap[device], err.Error())
			glog.Errorf("Unable to set performance determinism and clock frequency to %v for device %v", clkvalue, device)
		}
	}
	for id, msg := range errorMap {
		failedMessage = append(failedMessage, FailedMessage{ID: id, ErrorMsg: strings.Join(msg, "; ")})
	}
	return
}

// SetFanSpeed 设置风扇转速 [0-255]
// @Summary 设置风扇转速
// @Description 根据设备ID列表和给定的风扇速度，设置设备的风扇速度
// @Tags Device
// @Accept  json
// @Produce  json
// @Param dvIdList body []int true "设备 ID 列表"
// @Param fan query string true "风扇速度值或百分比（如 50%）"
// @Success 200 {string} string "成功信息"
// @Failure 400 {string} string "失败信息"
// @Router /SetFanSpeed [post]
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

// DevFanRpms 获取设备的风扇速度
// @Summary 获取设备的风扇速度
// @Description 获取指定设备的风扇速度（RPM）
// @Tags Device
// @Accept  json
// @Produce  json
// @Param dvInd path int true "设备索引"
// @Success 200 {integer} int64 "风扇速度 (RPM)"
// @Failure 400 {string} string "失败信息"
// @Router /DevFanRpms/{dvInd} [get]
func DevFanRpms(dvInd int) (speed int64, err error) {
	return rsmiDevFanRpmsGet(dvInd, 0)
}

// SetPerformanceLevel 设置设备性能等级
// @Summary 设置设备性能等级
// @Description 根据设备ID列表和给定的性能等级，设置设备的性能等级
// @Tags Device
// @Accept  json
// @Produce  json
// @Param deviceList body []int true "设备 ID 列表"
// @Param level query string true "性能等级 (auto, low, high, normal)"
// @Success 200 {array} FailedMessage
// @Failure 400 {object} FailedMessage
// @Router /SetPerformanceLevel [post]
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

// SetProfile 设置功率配置
// @Summary 设置功率配置
// @Description 根据设备ID列表和给定的功率配置文件，设置设备的功率配置
// @Tags Power
// @Accept  json
// @Produce  json
// @Param dvIdList body []int true "设备 ID 列表"
// @Param profile query string true "功率配置文件名称"
// @Success 200 {array} FailedMessage "设置成功的消息列表"
// @Failure 400 {object} FailedMessage "失败的消息列表"
// @Router /SetProfile [post]
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
					glog.Errorf("device:%v Failed to set profile to: %v", device, err.Error())
					failedMessages = append(failedMessages, FailedMessage{ID: device, ErrorMsg: fmt.Sprintf("Failed to set profile to: %s", profile)})
				}
			}
		}
	}

	return
}

// DevPowerProfileSet 设置设备功率配置文件
// @Summary 设置设备功率配置文件
// @Description 设置指定设备的功率配置文件
// @Tags Power
// @Accept  json
// @Produce  json
// @Param dvInd path int true "设备索引"
// @Param reserved query int true "保留参数，通常为0"
// @Param profile query int true "功率配置文件的枚举值"
// @Success 200 {string} string "成功信息"
// @Failure 400 {string} string "失败信息"
// @Router /DevPowerProfileSet [post]
func DevPowerProfileSet(dvInd int, reserved int, profile RSNIPowerProfilePresetMasks) (err error) {
	return rsmiDevPowerProfileSet(dvInd, reserved, profile)
}

func DevPowerProfilePresetsGet(dvInd, sensorInd int) (powerProfileStatus RSMPowerProfileStatus, err error) {
	return rsmiDevPowerProfilePresetsGet(dvInd, sensorInd)
}

// GetBus 获取设备总线信息
// @Summary 获取设备总线信息
// @Description 获取指定设备的总线信息
// @Tags Device
// @Accept  json
// @Produce  json
// @Param device path int true "设备索引"
// @Success 200 {string} string "设备总线ID"
// @Failure 400 {string} string "失败信息"
// @Router /GetBus/{device} [get]
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

// ShowAllConciseHw 显示设备硬件信息
// @Summary 显示设备硬件信息
// @Description 显示指定设备列表的简要硬件信息
// @Tags Hardware
// @Accept  json
// @Produce  json
// @Param dvIdList body []int true "设备 ID 列表"
// @Success 200 {string} string "成功信息"
// @Failure 400 {string} string "失败信息"
// @Router /ShowAllConciseHw [post]
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

// ShowClocks 显示时钟信息
// @Summary 显示时钟信息
// @Description 显示指定设备的时钟信息
// @Tags Clock
// @Accept  json
// @Produce  json
// @Param dvIdList body []int true "设备 ID 列表"
// @Success 200 {string} string "成功信息"
// @Failure 400 {string} string "失败信息"
// @Router /ShowClocks [post]
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
				fr := fmt.Sprintf("%.1fGT/s x%d", float64(bw.TransferRate.Frequency[x])/1000000000, bw.Lanes[x])
				if uint32(x) == bw.TransferRate.Current {
					glog.Infof("Device %d: %d %s *", device, x, fr)
				} else {
					glog.Infof("Device %d: %d %s", device, x, fr)
				}
			}
		}
	}
}

// ShowCurrentFans 展示风扇转速和风扇级别
// @Summary 展示风扇转速和风扇级别
// @Description 显示指定设备的当前风扇转速和风扇级别
// @Tags Fan
// @Accept  json
// @Produce  json
// @Param dvIdList body []int true "设备 ID 列表"
// @Param printJSON query bool true "是否以 JSON 格式打印输出"
// @Success 200 {string} string "成功信息"
// @Failure 400 {string} string "失败信息"
// @Router /ShowCurrentFans [post]
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

// ShowCurrentTemps 显示所有设备的所有可用温度传感器的温度
// @Summary 显示设备温度传感器数据
// @Tags Temperature
// @Param dvIdList query []int true "设备 ID 列表"
// @Success 200 {object} TemperatureInfo "温度信息列表"
// @Failure 400 {object} error "错误信息"
// @Router /ShowCurrentTemps [get]
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

// ShowFwInfo 显示给定设备列表中指定固件类型的固件版本信息
// @Summary 显示设备固件版本信息
// @Tags Firmware
// @Param dvIdList query []int true "设备 ID 列表"
// @Param fwType query []string true "固件类型列表"
// @Success 200 {object} []FirmwareInfo "固件版本信息列表"
// @Failure 400 {object} error "错误信息"
// @Router /ShowFwInfo [get]
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

// PidList 获取进程列表
// @Summary 获取计算进程列表
// @Tags Process
// @Success 200 {array} string "进程 ID 列表"
// @Failure 400 {object} error "错误信息"
// @Router /PidList [get]
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

// GetCoarseGrainUtil 获取设备的粗粒度利用率
// @Summary 获取设备粗粒度利用率
// @Tags Utilization
// @Param device query int true "设备 ID"
// @Param typeName query string false "利用率计数器类型名称"
// @Success 200 {array} RSMIUtilizationCounter "利用率计数器列表"
// @Failure 400 {object} error "错误信息"
// @Router /GetCoarseGrainUtil [get]
func GetCoarseGrainUtil(device int, typeName string) (utilizationCounters []RSMIUtilizationCounter, err error) {
	var length int

	if typeName != "" {
		// 获取特定类型的利用率计数器
		var i RSMIUtilizationCounterType
		var found bool
		for index, name := range utilizationCounterName {
			if name == typeName {
				i = RSMIUtilizationCounterType(index)
				found = true
				break
			}
		}
		if !found {
			glog.Infof("No such coarse grain counter type: %v", typeName)
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

// ShowGpuUse DCU使用率
// @Summary 显示设备的 DCU 使用率
// @Tags DCU
// @Param dvIdList query []int true "设备 ID 列表"
// @Success 200 {object} []DeviceUseInfo "设备使用信息列表"
// @Failure 400 {object} error "错误信息"
// @Router /ShowGpuUse [get]
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
		utilCounters, err := GetCoarseGrainUtil(device, typeName)
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

// ShowEnergy 展示设备消耗的能量
// @Summary 展示设备的能量消耗
// @Description 获取并展示指定设备的能量消耗情况。
// @Tags 设备
// @Param dvIdList query []int true "设备ID列表"
// @Success 200 {string} string "成功返回设备的能量消耗信息"
// @Failure 400 {string} string "请求参数错误"
// @Router /showEnergy [get]
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

// ShowMemInfo 展示设备的内存信息
// @Summary 展示设备内存信息
// @Description 获取并展示指定设备的内存使用情况，包括不同类型的内存。
// @Tags 设备
// @Param dvIdList query []int true "设备ID列表"
// @Param memTypes query []string true "内存类型列表，如 'all' 或指定类型"
// @Success 200 {string} string "成功返回设备的内存信息"
// @Failure 400 {string} string "请求参数错误"
// @Router /showMemInfo [get]
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

// ShowMemUse 展示设备的内存使用情况
// @Summary 展示设备内存使用情况
// @Description 获取并展示指定设备的当前内存使用百分比和其他相关的利用率数据。
// @Tags 设备
// @Param dvIdList query []int true "设备ID列表"
// @Success 200 {string} string "成功返回设备的内存使用信息"
// @Failure 400 {string} string "请求参数错误"
// @Router /showMemUse [get]
func ShowMemUse(dvIdList []int) {
	fmt.Println("Current Memory Use")
	for _, device := range dvIdList {
		busyPercent, err := rsmiDevMemoryBusyPercentGet(device)
		if err == nil {
			fmt.Println("device: ", device, "GPU memory use (%)", busyPercent)
		}
		typeName := "Memory Activity"
		utilCounters, err := GetCoarseGrainUtil(device, typeName)
		if err == nil {
			for _, utCounter := range utilCounters {
				fmt.Println("device: ", device, utilizationCounterName[utCounter.Type], utCounter.Value)
			}
		} else {
			glog.Errorf("Device %d: Failed to get coarse grain util counters: %v", device, err)
		}
	}
}

// ShowMemVendor 展示设备供应商信息
// @Summary 展示设备的内存供应商信息
// @Description 获取并展示指定设备的内存供应商信息。
// @Tags 设备
// @Param dvIdList query []int true "设备ID列表"
// @Success 200 {object} []DeviceMemVendorInfo "成功返回设备的内存供应商信息"
// @Failure 400 {string} string "请求参数错误"
// @Router /showMemVendor [get]
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

// ShowPcieBw 展示设备的PCIe带宽使用情况
// @Summary 展示设备的PCIe带宽使用情况
// @Description 获取并展示指定设备的PCIe带宽使用情况，包括发送和接收的带宽。
// @Tags 设备
// @Param dvIdList query []int true "设备ID列表"
// @Success 200 {object} []PcieBandwidthInfo "成功返回设备的PCIe带宽使用信息"
// @Failure 400 {string} string "请求参数错误"
// @Router /showPcieBw [get]
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

// ShowPcieReplayCount 展示设备的PCIe重放计数
// @Summary 展示设备的PCIe重放计数
// @Description 获取并展示指定设备的PCIe重放计数。
// @Tags 设备
// @Param dvIdList query []int true "设备ID列表"
// @Success 200 {object} []PcieReplayCountInfo "设备的PCIe重放计数信息"
// @Failure 400 {string} string "请求参数错误"
// @Router /showPcieReplayCount [get]
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

// ShowPids 展示进程信息
// @Summary 展示系统中正在运行的KFD进程信息
// @Description 获取并展示当前系统中运行的KFD进程的详细信息。
// @Tags 系统
// @Success 200 {string} string "成功返回进程信息"
// @Failure 400 {string} string "请求错误"
// @Router /showPids [get]
func ShowPids() (err error) {
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
			GetProcessName(pid),
			gpuNumber,
			vramUsage,
			sdmaUsage,
			cuOccupancy,
		})
	}

	fmt.Println("KFD process information:")
	print2DArray(dataArray)
	fmt.Printf("==========\n")
	return
}

func GetProcessName(pid int) string {
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

// ShowPower 展示设备的平均功率
// @Summary 展示设备的平均功率消耗
// @Description 获取并展示指定设备的平均图形功率消耗。
// @Tags 设备
// @Param dvIdList query []int true "设备ID列表"
// @Success 200 {object} []DevicePowerInfo "设备的功率信息"
// @Failure 400 {string} string "请求参数错误"
// @Router /showPower [get]
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

// 获取设备电压/频率曲线信息(K100 AI不支持)
func DevOdVoltInfoGet(deInd int) (odv RSMIOdVoltFreqData, err error) {
	odv, err = rsmiDevOdVoltInfoGet(deInd)
	return
}

// ShowPowerPlayTable 展示设备的GPU内存时钟频率和电压
// @Summary 展示设备的GPU内存时钟频率和电压
// @Description 获取并展示指定设备的GPU内存时钟频率和电压表。
// @Tags 设备
// @Param dvIdList query []int true "设备ID列表"
// @Success 200 {object} []DevicePowerPlayInfo "设备的GPU时钟频率和电压信息"
// @Failure 400 {string} string "请求参数错误"
// @Router /showPowerPlayTable [get]
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
		glog.Infof("DevicePowerPlayInfo:%v", dataToJson(devicePowerPlayInfos))
	}

	fmt.Println("===============================================================")
	return
}

// ShowProductName 显示设备列表中所请求的产品名称
// @Summary 显示设备的产品名称
// @Description 获取并显示指定设备的产品名称、供应商、系列、型号和SKU信息。
// @Tags 设备
// @Param dvIdList query []int true "设备ID列表"
// @Success 200 {array} DeviceproductInfo "设备的产品信息列表"
// @Failure 400 {string} string "请求参数错误"
// @Router /showProductName [get]
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

// ShowProfile 可用电源配置文件
// @Summary 显示设备的电源配置文件
// @Description 获取并显示指定设备的电源配置文件，包括可用的电源配置选项。
// @Tags 设备
// @Param dvIdList query []int true "设备ID列表"
// @Success 200 {array} DeviceProfile "设备的电源配置文件信息列表"
// @Failure 400 {string} string "请求参数错误"
// @Router /showProfile [get]
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

// ShowRange 电流或电压范围
// @Summary 显示设备的电流或电压范围
// @Description 获取并显示指定设备的有效电流或电压范围。
// @Tags 设备
// @Param dvIdList query []int true "设备ID列表"
// @Param rangeType query string true "范围类型 (sclk, mclk, voltage)"
// @Success 200 {string} string "设备的电流或电压范围信息"
// @Failure 400 {string} string "请求参数错误"
// @Router /showRange [get]
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

// ShowRetiredPages 显示设备列表中指定类型的退役页
// @Summary 显示设备的退役页信息
// @Description 获取并显示指定设备的退役内存页信息。
// @Tags 设备
// @Param dvIdList query []int true "设备ID列表"
// @Param retiredType query string false "退役类型 (默认为'all')"
// @Success 200 {string} string "设备的退役页信息"
// @Failure 400 {string} string "请求参数错误"
// @Router /showRetiredPages [get]
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

// ShowSerialNumber 设备序列号
// @Summary 显示设备的序列号
// @Description 获取并显示指定设备的序列号信息。
// @Tags 设备
// @Param dvIdList query []int true "设备ID列表"
// @Success 200 {array} DeviceSerialInfo "设备的序列号信息列表"
// @Failure 400 {string} string "请求参数错误"
// @Router /showSerialNumber [get]
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

// ShowUId 唯一设备ID
// @Summary 显示设备的唯一ID
// @Description 获取并显示指定设备的唯一ID信息。
// @Tags 设备
// @Param dvIdList query []int true "设备ID列表"
// @Success 200 {array} DeviceUIdInfo "设备的唯一ID信息列表"
// @Failure 400 {string} string "请求参数错误"
// @Router /showUId [get]
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

// ShowVbiosVersion 打印并返回设备的VBIOS版本信息
// @Summary 显示设备的VBIOS版本
// @Description 获取并显示指定设备的VBIOS版本信息。
// @Tags 设备
// @Param dvIdList query []int true "设备ID列表"
// @Success 200 {array} DeviceVBIOSInfo "设备的VBIOS版本信息列表"
// @Failure 400 {string} string "请求参数错误"
// @Router /showVbiosVersion [get]
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

// ShowEvents 显示设备的事件
// @Summary 显示设备的事件
// @Description 获取并显示指定设备的事件信息。
// @Tags 设备
// @Param dvIdList query []int true "设备ID列表"
// @Param eventTypes query []string true "事件类型列表"
// @Success 200 {string} string "成功返回设备的事件信息"
// @Failure 400 {string} string "请求参数错误"
// @Router /showEvents [get]
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

// ShowVoltage 当前电压信息
// @Summary 显示设备的电压信息
// @Description 获取并显示指定设备的当前电压信息。
// @Tags 设备
// @Param dvIdList query []int true "设备ID列表"
// @Success 200 {array} DeviceVoltageInfo "设备的电压信息列表"
// @Failure 400 {string} string "请求参数错误"
// @Router /showVoltage [get]
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

// ShowVoltageCurve 电压曲线点
// @Summary 显示设备的电压曲线点
// @Description 获取并显示指定设备的电压曲线点信息。
// @Tags 设备
// @Param dvIdList query []int true "设备ID列表"
// @Success 200 {string} string "设备的电压曲线点信息"
// @Failure 400 {string} string "请求参数错误"
// @Router /showVoltageCurve [get]
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

// ShowXgmiErr 显示指定设备的 XGMI 错误状态。
//
// @Summary 显示 XGMI 错误状态
// @Description 显示一组 GPU 设备的 XGMI 错误状态。
// @Tags Topology
// @Param dvIdList query []int true "设备 ID 列表"
// @Param printJSON query bool false "是否以 JSON 格式输出"
// @Success 200 {string} string "XGMI 错误状态信息"
// @Router /showXgmiErr [get]
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

// ShowWeightTopology 显示 GPU 拓扑中两台设备之间的权重。
// @Summary 显示 GPU 拓扑权重
// @Description 显示 GPU 设备之间的权重信息。
// @Tags Topology
// @Param dvIdList query []int true "设备 ID 列表"
// @Param printJSON query bool false "是否以 JSON 格式输出"
// @Success 200 {string} string "GPU 拓扑权重信息"
// @Router /showWeightTopology [get]
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

// ShowHopsTopology 显示 GPU 拓扑中两台设备之间的跳数。
// @Summary 显示 GPU 拓扑跳数
// @Description 显示 GPU 设备之间的跳数信息。
// @Tags Topology
// @Param dvIdList query []int true "设备 ID 列表"
// @Param printJSON query bool false "是否以 JSON 格式输出"
// @Success 200 {string} string "GPU 拓扑跳数信息"
// @Router /showHopsTopology [get]

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

// ShowTypeTopology 显示 GPU 拓扑中两台设备之间的链接类型。
// @Summary 显示 GPU 拓扑链接类型
// @Description 显示 GPU 设备之间的链接类型信息。
// @Tags Topology
// @Param dvIdList query []int true "设备 ID 列表"
// @Param printJSON query bool false "是否以 JSON 格式输出"
// @Success 200 {string} string "GPU 拓扑链接类型信息"
// @Router /showTypeTopology [get]
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

// ShowNumaTopology 显示指定设备的 NUMA 节点信息。
// @Summary 显示 NUMA 节点信息
// @Description 显示一组 DCU 设备的 NUMA 节点和关联信息。
// @Tags Topology
// @Param dvIdList query []int true "设备 ID 列表"
// @Success 200 {string} string "NUMA 节点信息"
// @Router /showNumaTopology [get]
func ShowNumaTopology(dvIdList []int) (numaInfos []NumaInfo, err error) {
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
		// 将设备和 NUMA 信息存储在结构体中并添加到切片中
		numaInfo := NumaInfo{
			DeviceID:     device,
			NumaNode:     numaNode,
			NumaAffinity: numaAffinity,
		}
		numaInfos = append(numaInfos, numaInfo)

		glog.Infof("Device %d: Numa Node: %d, Numa Affinity: %d\n", device, numaNode, numaAffinity)
	}
	return
}

// ShowHwTopology 显示指定设备的完整硬件拓扑信息。
// @Summary 显示完整的硬件拓扑信息
// @Description 显示一组 GPU 设备的权重、跳数、链接类型和 NUMA 节点信息。
// @Tags Topology
// @Param dvIdList query []int true "设备 ID 列表"
// @Success 200 {string} string "完整的硬件拓扑信息"
// @Router /showHwTopology [get]
func ShowHwTopology(dvIdList []int) {
	ShowWeightTopology(dvIdList, true)

	ShowHopsTopology(dvIdList, true)

	ShowTypeTopology(dvIdList, true)

	ShowNumaTopology(dvIdList)
}

/*************************************VDCU******************************************/
// DeviceCount 返回设备的数量。
// @Summary 获取设备数量
// @Description 获取当前系统中的设备数量。
// @Tags Device
// @Success 200 {int} int "设备数量"
// @Failure 500 {object} string "内部服务器错误"
// @Router /deviceCount [get]
func DeviceCount() (count int, err error) {
	return dmiGetDeviceCount()
}

// VDeviceSingleInfo
// @Summary 获取单个虚拟设备的信息
// @Description 根据设备索引获取对应的虚拟设备信息
// @Tags VirtualDevice
// @Param vDvInd query int true "设备索引"
// @Success 200 {object} DMIVDeviceInfo "虚拟设备信息"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "内部服务器错误"
// @Router /VDeviceSingleInfo [get]
func VDeviceSingleInfo(vDvInd int) (vDeviceInfo DMIVDeviceInfo, err error) {
	glog.Infof("VDeviceSingleInfo vDvInd:%v", vDvInd)
	return dmiGetVDeviceInfo(vDvInd)
}

// VDeviceCount 返回虚拟设备的数量。
// @Summary 获取虚拟设备数量
// @Description 获取当前系统中的虚拟设备数量。
// @Tags Device
// @Success 200 {int} int "虚拟设备数量"
// @Failure 500 {object} string "内部服务器错误"
// @Router /vDeviceCount [get]
func VDeviceCount() (count int, err error) { return dmiGetVDeviceCount() }

// DeviceRemainingInfo 返回指定物理设备的剩余计算单元（CU）和内存信息。
// @Summary 获取设备剩余信息
// @Description 获取指定设备的剩余计算单元和内存信息。
// @Tags Device
// @Param dvInd path int true "物理设备索引"
// @Success 200 {string} uint64 "剩余的CU信息"
// @Success 200 {string} uint64 "剩余的内存信息"
// @Failure 400 {object} string "无效的设备索引"
// @Failure 500 {object} string "内部服务器错误"
// @Router /deviceRemainingInfo/{dvInd} [get]
func DeviceRemainingInfo(dvInd int) (cus, memories uint64, err error) {
	return dmiGetDeviceRemainingInfo(dvInd)
}

// CreateVDevices 创建指定数量的虚拟设备
// @Summary 创建虚拟设备
// @Description 在指定的物理设备上创建指定数量的虚拟设备，返回创建的虚拟设备ID集合。
// @Tags 虚拟设备
// @Param dvInd query int true "物理设备的索引"
// @Param vDevCount query int true "要创建的虚拟设备数量"
// @Param vDevCUs query []int true "每个虚拟设备的计算单元数量"
// @Param vDevMemSize query []int true "每个虚拟设备的内存大小"
// @Success 200 {array} int "虚拟设备创建成功，返回虚拟设备ID集合"
// @Failure 400 {string} string "创建虚拟设备失败"
// @Router /CreateVDevices [post]
func CreateVDevices(dvInd int, vDevCount int, vDevCUs []int, vDevMemSize []int) (vdevIDs []int, err error) {
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

// StartVDevice 启动虚拟设备
// @Summary 启动指定的虚拟设备
// @Description 启动虚拟设备，指定设备索引
// @Tags VirtualDevice
// @Param vDvInd path int true "虚拟设备索引"
// @Success 200 {string} string "操作成功"
// @Failure 400 {string} string "操作失败"
// @Router /StartVDevice/{vDvInd} [get]
func StartVDevice(vDvInd int) (err error) {
	return dmiStartVDevice(vDvInd)
}

func DevBusyPercent(dvInd int) (percent int, err error) {
	return dmiGetDevBusyPercent(dvInd)
}

func VDevBusyPercent(vDvInd int) (percent int, err error) {
	return dmiGetDevBusyPercent(vDvInd)
}

// StopVDevice 停止虚拟设备
// @Summary 停止指定的虚拟设备
// @Description 停止虚拟设备，指定设备索引
// @Tags VirtualDevice
// @Param vDvInd path int true "虚拟设备索引"
// @Success 200 {string} string "操作成功"
// @Failure 400 {string} string "操作失败"
// @Router /StopVDevice/{vDvInd} [get]
func StopVDevice(vDvInd int) (err error) {
	return dmiStopVDevice(vDvInd)
}

// SetEncryptionVMStatus 设置虚拟机加密状态
// @Summary 设置虚拟机加密状态
// @Description 根据提供的状态开启或关闭虚拟机加密
// @Tags VirtualDevice
// @Param status query bool true "加密状态"
// @Success 200 {string} string "操作成功"
// @Failure 400 {string} string "操作失败"
// @Router /SetEncryptionVMStatus [post]
func SetEncryptionVMStatus(status bool) (err error) {
	return dmiSetEncryptionVMStatus(status)
}

// EncryptionVMStatus 获取加密虚拟机状态
// @Summary 获取当前虚拟机的加密状态
// @Description 返回虚拟机是否处于加密状态
// @Tags VirtualDevice
// @Success 200 {boolean} boolean "加密状态"
// @Failure 400 {string} string "操作失败"
// @Router /EncryptionVMStatus [get]
func EncryptionVMStatus() (status bool, err error) {
	return dmiGetEncryptionVMStatus()
}

// PrintEventList 打印事件列表
// @Summary 打印设备的事件列表
// @Description 打印指定设备的事件列表，并设置延迟
// @Tags Event
// @Param device path int true "设备索引"
// @Param delay query int true "延迟时间（秒）"
// @Param eventList query []string true "事件列表"
// @Success 200 {string} string "操作成功"
// @Failure 400 {string} string "操作失败"
// @Router /PrintEventList/{device} [get]
func PrintEventList(device int, delay int, eventList []string) {
	printEventList(device, delay, eventList)
}

func GetDeviceInfo(dvInd int) (deviceInfo DMIDeviceInfo, err error) {
	return dmiGetDeviceInfo(dvInd)
}

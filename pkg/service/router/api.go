package router

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"g.sugon.com/das/dcgm-dcu/pkg/dcgm"
)

// @Summary 获取设备名称
// @Description 根据设备 ID 获取设备名称
// @Accept  json
// @Produce  json
// @Param   dvInd     path   int     true  "Device ID"
// @Success 200 {string} string "设备名称"
// @Failure 400 {object} error "Invalid device ID"
// @Failure 500 {object} error "Internal Server Error"
// @Router /devicename/{dvInd} [get]
func GetDevName(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid device ID"))
		return
	}

	name, err := dcgm.DevName(dvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	deviceName := map[string]interface{}{
		"deviceName": name,
	}

	c.JSON(http.StatusOK, SuccessResponse(deviceName))
}

// @Summary 获取 GPU 数量
// @Description 获取监视的 GPU 数量
// @Produce json
// @Success 200 {int} int "GPU 数量"
// @Failure 500 {object} error "获取 GPU 数量失败"
// @Router /NumMonitorDevices [get]
func GetNumMonitorDevices(c *gin.Context) {
	num, err := dcgm.NumMonitorDevices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse("获取 GPU 数量失败"))
		return
	}
	gpuCount := map[string]interface{}{
		"gpuCount": num,
	}
	c.JSON(http.StatusOK, SuccessResponse(gpuCount))
}

// @Summary 获取设备SKU
// @Description 根据设备索引获取SKU
// @Produce json
// @Param dvInd path int true "设备索引"
// @Success 200 {int} sku "返回设备SKU"
// @Failure 400 {object} error "获取设备SKU失败"
// @Router /DevSku/{dvInd} [get]
func GetDevSku(c *gin.Context) {
	dvIndStr := c.Param("dvInd")
	dvInd, err := strconv.Atoi(dvIndStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("设备索引无效"))
		return
	}
	sku, err := dcgm.DevSku(dvInd)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("获取设备SKU失败"))
		return
	}
	gpuSku := map[string]interface{}{
		"sku": sku,
	}
	c.JSON(http.StatusOK, SuccessResponse(gpuSku))
}

// 获取设备品牌名称
// @Summary 获取设备品牌名称
// @Description 根据设备索引获取品牌名称
// @Param dvInd path int true "设备索引"
// @Success 200 {string} brand "设备品牌名称"
// @Failure 400 {object} error "请求失败"
// @Router /DevBrand/{dvInd} [get]
func DevBrand(c *gin.Context) {
	dvInd, _ := strconv.Atoi(c.Param("dvInd"))
	brand, err := dcgm.DevBrand(dvInd)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}
	devBrand := map[string]interface{}{
		"brand": brand,
	}
	c.JSON(http.StatusOK, SuccessResponse(devBrand))
}

// 获取设备供应商名称
// @Summary 获取设备供应商名称
// @Description 根据设备索引获取供应商名称
// @Param dvInd path int true "设备索引"
// @Success 200 {string} bname "返回设备供应商名称"
// @Failure 400 {object} error "请求失败"
// @Router /DevVendorName/{dvInd} [get]
func DevVendorName(c *gin.Context) {
	dvInd, _ := strconv.Atoi(c.Param("dvInd"))
	bname, err := dcgm.DevVendorName(dvInd)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}
	devBrandName := map[string]interface{}{
		"bname": bname,
	}
	c.JSON(http.StatusOK, SuccessResponse(devBrandName))
}

// 获取设备显存供应商名称
// @Summary 获取设备显存供应商名称
// @Description 根据设备索引获取显存供应商名称
// @Param dvInd path int true "设备索引"
// @Success 200 {string} name "返回显存供应商名称"
// @Failure 400 {object} error "请求失败"
// @Router /DevVramVendor/{dvInd} [get]
func DevVramVendor(c *gin.Context) {
	dvInd, _ := strconv.Atoi(c.Param("dvInd"))
	name, err := dcgm.DevVramVendor(dvInd)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}
	devVramVendor := map[string]interface{}{
		"vendorName": name,
	}
	c.JSON(http.StatusOK, SuccessResponse(devVramVendor))
}

// 获取可用的 PCIe 带宽列表
// @Summary 获取可用的 PCIe 带宽列表
// @Description 根据设备 ID 获取设备的可用 PCIe 带宽列表。
// @Param dvInd path int true "设备 ID"
// @Success 200 {object} RSMIPcieBandwidth "PCIe 带宽列表"
// @Failure 400 {object} error "请求错误"
// @Failure 404 {object} error "设备未找到"
// @Router /DevPciBandwidth/{dvInd} [get]
func DevPciBandwidth(c *gin.Context) {
	dvInd, _ := strconv.Atoi(c.Param("dvInd"))
	rsmiPcieBandwidth, err := dcgm.DevPciBandwidth(dvInd)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("请求失败"))
		return
	}
	response := map[string]interface{}{
		"rsmiPcieBandwidth": rsmiPcieBandwidth,
	}
	c.JSON(http.StatusOK, SuccessResponse(response))
}

// DevPciBandwidthSet 设置设备允许的 PCIe 带宽
// @Summary 设置设备允许的 PCIe 带宽
// @Description 根据设备索引和带宽掩码限制设备允许的 PCIe 带宽
// @Produce json
// @Param dvInd query int true "设备索引"
// @Param bwBitmask query int64 true "带宽掩码"
// @Success 200 {string} string "操作成功"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "服务器内部错误"
// @Router /DevPciBandwidthSet [post]
func DevPciBandwidthSet(c *gin.Context) {
	var dvInd int
	var bwBitmask int64

	// 获取 query 中的 dvInd 和 bwBitmask 参数
	if err := c.ShouldBindQuery(&dvInd); err != nil {
		c.JSON(http.StatusBadRequest, "Invalid dvInd parameter")
		return
	}
	if err := c.ShouldBindQuery(&bwBitmask); err != nil {
		c.JSON(http.StatusBadRequest, "Invalid bwBitmask parameter")
		return
	}

	// 调用已有的 DevPciBandwidthSet 函数
	if err := dcgm.DevPciBandwidthSet(dvInd, bwBitmask); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(nil))
}

// 获取内存使用百分比
// @Summary 获取内存使用百分比
// @Description 根据设备 ID 获取设备内存的使用百分比。
// @Param dvInd path int true "设备 ID"
// @Success 200 {int} busyPercent "内存使用百分比"
// @Failure 400 {object} error "请求错误"
// @Failure 404 {object} error "设备未找到"
// @Router /MemoryPercent/{dvInd} [get]
func MemoryPercent(c *gin.Context) {
	dvInd, _ := strconv.Atoi(c.Param("dvInd"))
	busyPercent, err := dcgm.MemoryPercent(dvInd)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}
	memoryPercent := map[string]interface{}{
		"busyPercent": busyPercent,
	}
	c.JSON(http.StatusOK, SuccessResponse(memoryPercent))
}

// 设置设备 PowerPlay 性能级别
// @Summary 设置设备 PowerPlay 性能级别
// @Description 根据设备 ID 设置 PowerPlay 性能级别。
// @Param dvInd path int true "设备 ID"
// @Param level query string true "要设置的性能级别"
// @Success 200 {string} string "操作成功"
// @Failure 400 {object} error "请求错误"
// @Failure 404 {object} error "设备未找到"
// @Router /DevPerfLevelSet/{dvInd} [post]
func DevPerfLevelSet(c *gin.Context) {
	dvInd, _ := strconv.Atoi(c.Param("dvInd"))
	level := c.Query("level")

	// 将 level 字符串转换为 RSMIDevPerfLevel 类型
	levelConverted, err := ConvertToRSMIDevPerfLevel(level)
	if err != nil {
		// 如果转换失败，返回错误响应
		c.JSON(http.StatusBadRequest, ErrorResponse("无效的性能级别"))
		return
	}

	// 调用 dcgm.DevPerfLevelSet 并传入转换后的 level
	err = dcgm.DevPerfLevelSet(dvInd, levelConverted)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(nil))
}

// 获取 GPU 度量信息
// @Summary 获取 GPU 度量信息
// @Description 根据设备 ID 获取 GPU 的度量信息。
// @Param dvInd path int true "设备 ID"
// @Success 200 {object} RSMIGPUMetrics "GPU 度量信息"
// @Failure 400 {object} error "请求错误"
// @Failure 404 {object} error "设备未找到"
// @Router /DevGpuMetricsInfo/{dvInd} [get]
func DevGpuMetricsInfo(c *gin.Context) {
	dvInd, _ := strconv.Atoi(c.Param("dvInd"))
	gpuMetrics, err := dcgm.DevGpuMetricsInfo(dvInd)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, gpuMetrics)
}

// 获取设备监控中的指标
// @Summary 获取设备监控中的指标
// @Description 收集所有设备的监控指标信息。
// @Success 200 {array} MonitorInfo "设备监控指标信息列表"
// @Failure 400 {object} error "请求错误"
// @Failure 404 {object} error "设备未找到"
// @Router /CollectDeviceMetrics [get]
func CollectDeviceMetrics(c *gin.Context) {
	monitorInfos, err := dcgm.CollectDeviceMetrics()
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("请求失败"))
		return
	}
	c.JSON(http.StatusOK, monitorInfos)
}

// @Summary 获取设备信息
// @Description 根据设备 ID 获取物理设备的详细信息
// @Accept  json
// @Produce  json
// @Param   dvInd     path   int     true  "Device ID"
// @Success 200 {object} PhysicalDeviceInfo "设备信息"
// @Failure 400 {object} error "Invalid device ID"
// @Failure 500 {object} error "Internal Server Error"
// @Router /deviceinfo/{dvInd} [get]
func GetDeviceByDvInd(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid device ID"))
		return
	}
	deviceInfo, err := dcgm.GetDeviceByDvInd(dvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, deviceInfo)
}

// @Summary 获取所有物理设备信息
// @Description 该接口返回所有物理设备的详细信息。
// @Produce json
// @Success 200 {array} PhysicalDeviceInfo "所有设备的详细信息"
// @Failure 400 {object} error "无效的请求参数"
// @Failure 500 {object} error "服务器内部错误"
// @Router /AllDeviceInfos [get]
func AllDeviceInfos(c *gin.Context) {
	deviceInfos, err := dcgm.AllDeviceInfos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, deviceInfos)
}

// @Summary 获取总线信息
// @Description 获取设备的总线信息 (BDF格式)
// @Accept  json
// @Produce  json
// @Param   dvInd     path   int     true  "Device ID"
// @Success 200 {string} string "总线信息"
// @Failure 400 {object} error "Invalid device ID"
// @Failure 500 {object} error "Internal Server Error"
// @Router /picbusinfo/{dvInd} [get]
func PicBusInfo(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid device ID"))
		return
	}
	busInfo, err := dcgm.PicBusInfo(dvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	picBusInfo := map[string]interface{}{
		"busInfo": busInfo,
	}
	c.JSON(http.StatusOK, SuccessResponse(picBusInfo))
}

// @Summary 获取风扇转速
// @Description 获取指定设备的风扇转速及其占最大转速的百分比
// @Accept  json
// @Produce  json
// @Param   dvInd     path   int     true  "Device ID"
// @Success 200 {object} map[string]int "风扇转速信息"
// @Failure 400 {object} error "Invalid device ID"
// @Failure 500 {object} error "Internal Server Error"
// @Router /fanspeedinfo/{dvInd} [get]
func FanSpeedInfo(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid device ID"))
		return
	}

	fanLevel, fanPercentage, err := dcgm.FanSpeedInfo(dvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	fanSpeedInfo := map[string]interface{}{
		"fanLevel":      fanLevel,
		"fanPercentage": fanPercentage,
	}

	c.JSON(http.StatusOK, SuccessResponse(fanSpeedInfo))
}

// @Summary 获取DCU使用率
// @Description 获取指定设备的DCU当前使用百分比
// @Accept  json
// @Produce  json
// @Param   dvInd     path   int     true  "Device ID"
// @Success 200 {int} int "GPU使用百分比"
// @Failure 400 {object} error "Invalid device ID"
// @Failure 500 {object} error "Internal Server Error"
// @Router /gpuuse/{dvInd} [get]
func GPUUse(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid device ID"))
		return
	}

	gpuUsage, err := dcgm.GPUUse(dvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	gpuUse := map[string]interface{}{
		"gpuUsage": gpuUsage,
	}
	c.JSON(http.StatusOK, SuccessResponse(gpuUse))
}

// @Summary 获取设备ID的十六进制值
// @Description 根据设备索引返回设备ID的十六进制值
// @Produce json
// @Param dvInd query int true "设备索引"
// @Success 200 {int} id "返回设备ID的十六进制值"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /DevID [get]
func GetDevID(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Query("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid device ID"))
		return
	}

	id, err := dcgm.DevID(dvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	devID := map[string]interface{}{
		"id": id,
	}
	c.JSON(http.StatusOK, SuccessResponse(devID))
}

// @Summary 获取设备的最大功率
// @Description 根据设备索引返回设备的最大功率（以瓦特为单位）
// @Produce json
// @Param dvInd query int true "设备索引"
// @Success 200 {int64} power "返回设备的最大功率"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /MaxPower [get]
func GetMaxPower(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Query("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid device ID"))
		return
	}

	power, err := dcgm.MaxPower(dvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	devMaxPower := map[string]interface{}{
		"power": power,
	}
	c.JSON(http.StatusOK, SuccessResponse(devMaxPower))
}

// @Summary 获取设备的指定内存使用情况
// @Description 根据设备索引和内存类型返回内存的使用量和总量
// @Produce json
// @Param dvInd query int true "设备索引"
// @Param memType query string true "内存类型（可选值: vram, vis_vram, gtt）"
// @Success 200 {object} map[string]interface{} "返回指定内存类型的使用量和总量"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /MemInfo [get]
func GetMemInfo(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Query("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid device ID"))
		return
	}

	memType := c.Query("memType")
	memUsed, memTotal, err := dcgm.MemInfo(dvInd, memType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	memInfo := map[string]interface{}{
		"memUsed":  memUsed,
		"memTotal": memTotal,
	}

	c.JSON(http.StatusOK, SuccessResponse(memInfo))
}

// @Summary 获取设备信息列表
// @Description 返回所有设备的详细信息列表
// @Produce json
// @Success 200 {object} DeviceInfo "返回设备信息列表"
// @Failure 500 {object} error "服务器内部错误"
// @Router /DeviceInfos [get]
func GetDeviceInfos(c *gin.Context) {
	deviceInfos, err := dcgm.DeviceInfos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, deviceInfos)
}

// @Summary 获取指定PID的进程名
// @Description 根据进程ID（PID）返回对应的进程名称
// @Produce json
// @Param pid query int true "进程ID"
// @Success 200 {string} string "返回进程名称"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /ProcessName [get]
func GetProcessName(c *gin.Context) {
	pid, err := strconv.Atoi(c.Query("pid"))
	if err != nil || pid < 1 {
		c.JSON(http.StatusBadRequest, ErrorResponse("请求参数错误"))
		return
	}

	pName := dcgm.ProcessName(pid)
	processName := map[string]interface{}{
		"pName": pName,
	}
	c.JSON(http.StatusOK, SuccessResponse(processName))
}

// PerfLevel 获取设备的当前性能水平
// @Summary 获取设备的当前性能水平
// @Description 返回指定设备的当前性能等级
// @Produce json
// @Param dvInd path int true "设备索引"
// @Success 200 {string} string "返回当前性能水平"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /PerfLevel/{dvInd} [get]
func PerfLevel(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid device ID"))
		return
	}
	perf, err := dcgm.PerfLevel(dvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	perfLevel := map[string]interface{}{
		"perfLevel": perf,
	}
	c.JSON(http.StatusOK, SuccessResponse(perfLevel))
}

// Power 获取设备的平均功耗
// @Summary 获取设备的平均功耗
// @Description 返回指定设备的平均功耗
// @Produce json
// @Param dvInd path int true "设备索引"
// @Success 200 {int64} int64 "返回平均功耗（瓦特）"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /Power/{dvInd} [get]
func Power(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid device ID"))
		return
	}
	power, err := dcgm.Power(dvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	devPower := map[string]interface{}{
		"power": power,
	}
	c.JSON(http.StatusOK, SuccessResponse(devPower))
}

// EccStatus 获取GPU块的ECC状态
// @Summary 获取GPU块的ECC状态
// @Description 返回指定GPU块的ECC状态
// @Produce json
// @Param dvInd path int true "设备索引"
// @Param block query string true "GPU块"
// @Success 200 {string} string "返回ECC状态"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /EccStatus/{dvInd} [get]
func EccStatus(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid device ID"))
		return
	}
	block := c.Query("block")
	if block == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse("Missing block parameter"))
		return
	}

	// 将 block 字符串转换为 RSMIGpuBlock 类型
	blockConverted, err := ConvertToRSMIGpuBlock(block)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid block parameter"))
		return
	}

	eccStatus, err := dcgm.EccStatus(dvInd, blockConverted)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	status := map[string]interface{}{
		"eccStatus": eccStatus,
	}
	c.JSON(http.StatusOK, SuccessResponse(status))
}

// Temperature 获取设备温度
// @Summary 获取设备温度
// @Description 返回指定设备的当前温度
// @Produce json
// @Param dvInd path int true "设备索引"
// @Param sensorType query int true "传感器类型:0：Edge GPU temperature; 1:Junction/hotspot temperature;2:VRAM temperature"
// @Success 200 {float64} float64 "返回温度（摄氏度）"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /Temperature/{dvInd} [get]
func Temperature(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid device ID"))
		return
	}
	sensorType, err := strconv.Atoi(c.Query("sensorType"))
	if err != nil || sensorType < 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid sensor type"))
		return
	}
	temp, err := dcgm.Temperature(dvInd, sensorType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	devTemp := map[string]interface{}{
		"temp": temp,
	}
	c.JSON(http.StatusOK, SuccessResponse(devTemp))
}

// VbiosVersion 获取设备的VBIOS版本
// @Summary 获取设备的VBIOS版本
// @Description 返回指定设备的VBIOS版本
// @Produce json
// @Param dvInd path int true "设备索引"
// @Success 200 {string} string "返回VBIOS版本"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /VbiosVersion/{dvInd} [get]
func VbiosVersion(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid device ID"))
		return
	}
	vbios, err := dcgm.VbiosVersion(dvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	devVbios := map[string]interface{}{
		"vbios": vbios,
	}
	c.JSON(http.StatusOK, SuccessResponse(devVbios))
}

// Swagger 注解
// @Summary 设置 GPU 时钟频率
// @Description 设置 GPU 上指定时钟的允许频率。clkType 设置为默认值，无需传递。
// @Tags GPU
// @Accept json
// @Produce json
// @Param dvInd query int true "设备索引"
// @Param freqBitmask query int64 true "频率掩码"
// @Success 200 {object} map[string]interface{} "成功"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /DevGpuClkFreqSet [post]
func DevGpuClkFreqSet(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Query("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("无效的 dvInd 参数"))
		return
	}
	freqBitmask, err := strconv.ParseInt(c.Query("freqBitmask"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("无效的 bwBitmask 参数"))
		return
	}

	// 调用 DevGpuClkFreqSet 函数
	err = dcgm.DevGpuClkFreqSet(dvInd, dcgm.RSMI_CLK_TYPE_SYS, freqBitmask)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(nil))
}

// Version 获取当前系统的驱动程序版本
// @Summary 获取当前系统的驱动程序版本
// @Description 返回指定组件的驱动程序版本
// @Produce json
// @Param component query string true "驱动组件:FIRST、DRIVER、LAST"
// @Success 200 {string} string "返回驱动程序版本"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /Version [get]
func Version(c *gin.Context) {
	component := c.Query("component")
	if component == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse("Missing component parameter"))
		return
	}

	// 将 component 字符串转换为 RSMISwComponent 类型
	componentConverted, err := ConvertToRSMISwComponent(component)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid component parameter"))
		return
	}

	version, err := dcgm.Version(componentConverted)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	devVersion := map[string]interface{}{
		"version": version,
	}
	c.JSON(http.StatusOK, SuccessResponse(devVersion))
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
func ResetClocks(c *gin.Context) {
	var dvIdList []int
	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid JSON data"))
		return
	}
	failedMessages := dcgm.ResetClocks(dvIdList)
	response := map[string]interface{}{
		"failedMessages": failedMessages,
	}
	c.JSON(http.StatusOK, SuccessResponse(response))
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
func ResetFans(c *gin.Context) {
	var dvIdList []int
	if err := c.ShouldBindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}
	if err := dcgm.ResetFans(dvIdList); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(nil))
}

// ResetProfile 重置设备的配置文件
// @Summary 重置指定设备的电源配置文件和性能级别
// @Produce json
// @Param dvIdList body []int true "设备ID列表"
// @Success 200 {array} FailedMessage "返回失败的设备及其错误信息"
// @Router /ResetProfile [post]
func ResetProfile(c *gin.Context) {
	var dvIdList []int
	if err := c.ShouldBindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}
	failedMessage := dcgm.ResetProfile(dvIdList)
	response := map[string]interface{}{
		"failedMessages": failedMessage,
	}
	c.JSON(http.StatusOK, SuccessResponse(response))
}

// ResetXGMIErr 重置设备的XGMI错误状态
// @Summary 重置指定设备的XGMI错误状态
// @Produce json
// @Param dvIdList body []int true "设备ID列表"
// @Success 200 {array} FailedMessage "返回失败的设备及其错误信息"
// @Router /ResetXGMIErr [post]
func ResetXGMIErr(c *gin.Context) {
	var dvIdList []int
	if err := c.ShouldBindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}
	failedMessage := dcgm.ResetXGMIErr(dvIdList)
	response := map[string]interface{}{
		"failedMessages": failedMessage,
	}
	c.JSON(http.StatusOK, SuccessResponse(response))
}

// XGMIErrorStatus 获取XGMI错误状态
// @Summary 获取XGMI错误状态
// @Description 获取指定物理设备的XGMI（高速互连链路）错误状态。
// @Param dvInd query int true "物理设备的索引"
// @Success 200 {integer} int "返回XGMI错误状态码"
// @Failure 400 {string} string "获取XGMI错误状态失败"
// @Router /XGMIErrorStatus [get]
func XGMIErrorStatus(c *gin.Context) {
	dvIndStr := c.Query("dvInd")
	dvInd, err := strconv.Atoi(dvIndStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid device ID"))
		return
	}
	status, err := dcgm.XGMIErrorStatus(dvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	errorStatus := map[string]interface{}{
		"status": status,
	}
	c.JSON(http.StatusOK, SuccessResponse(errorStatus))
}

// XGMIHiveIdGet 获取设备的XGMI hive id
// @Summary 获取指定设备的XGMI hive id
// @Produce json
// @Param dvInd query int true "设备索引"
// @Success 200 {int64} int64 "返回设备的XGMI hive id"
// @Router /XGMIHiveIdGet [get]
func XGMIHiveIdGet(c *gin.Context) {
	dvIndStr := c.Query("dvInd")
	dvInd, err := strconv.Atoi(dvIndStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid device ID"))
		return
	}
	hiveId, err := dcgm.XGMIHiveIdGet(dvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	xgmihiveId := map[string]interface{}{
		"hiveId": hiveId,
	}
	c.JSON(http.StatusOK, SuccessResponse(xgmihiveId))
}

// ResetPerfDeterminism 处理重置Performance Determinism
// @Summary 重置Performance Determinism
// @Description 该接口用于重置指定设备的性能决定性设置。请求体中需要包含设备ID列表。
// @Accept json
// @Produce json
// @Param dvIdList body []int true "设备ID列表"
// @Success 200 {array} FailedMessage "返回失败的设备及其错误信息"
// @Failure 400 {object} error "无效的请求体"
// @Router /ResetPerfDeterminism [post]
func ResetPerfDeterminism(c *gin.Context) {
	var dvIdList []int
	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid request body"))
		return
	}

	failedMessages := dcgm.ResetPerfDeterminism(dvIdList)
	if len(failedMessages) > 0 {
		response := map[string]interface{}{
			"failedMessages": failedMessages,
		}
		c.JSON(http.StatusBadRequest, ErrorResponse(response))
	} else {
		c.JSON(http.StatusOK, SuccessResponse(nil))
	}
}

// SetClockRange 处理设置时钟频率范围
// @Summary 设置设备的时钟频率范围
// @Description 设置设备的时钟频率范围（sclk 或 mclk）
// @Accept json
// @Produce json
// @Param dvIdList body []int true "设备ID列表"
// @Param clkType query string true "时钟类型（sclk 或 mclk）"
// @Param minvalue query string true "最小值（MHz）"
// @Param maxvalue query string true "最大值（MHz）"
// @Param autoRespond query bool false "自动响应超出规格的警告"
// @Success 200 {string} string "时钟范围设置成功"
// @Failure 400 {object} error "无效的请求参数或无法设置时钟范围"
// @Router /SetClockRange [post]
func SetClockRange(c *gin.Context) {
	var dvIdList []int
	clkType := c.Query("clkType")
	minvalue := c.Query("minvalue")
	maxvalue := c.Query("maxvalue")
	autoRespond, _ := strconv.ParseBool(c.Query("autoRespond"))

	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid request body"))
		return
	}

	failedMessages := dcgm.SetClockRange(dvIdList, clkType, minvalue, maxvalue, autoRespond)
	if len(failedMessages) > 0 {
		response := map[string]interface{}{
			"failedMessages": failedMessages,
		}
		c.JSON(http.StatusBadRequest, ErrorResponse(response))
	} else {
		c.JSON(http.StatusOK, SuccessResponse(nil))
	}
}

// SetPowerPlayTableLevel 处理设置PowerPlay表级别
// @Summary 设置设备的PowerPlay表级别
// @Description 设置设备的PowerPlay表级别
// @Accept json
// @Produce json
// @Param dvIdList body []int true "设备ID列表"
// @Param clkType query string true "时钟类型（sclk 或 mclk）"
// @Param point query string true "电压点"
// @Param clk query string true "时钟值（MHz）"
// @Param volt query string true "电压值（mV）"
// @Param autoRespond query bool false "自动响应超出规格的警告"
// @Success 200 {object} map[string]interface{} "成功设置PowerPlay表级别"
// @Failure 400 {object} map[string]interface{} "无效的请求参数或无法设置PowerPlay表级别，返回失败消息列表"
// @Router /SetPowerPlayTableLevel [post]
func SetPowerPlayTableLevel(c *gin.Context) {
	var dvIdList []int
	clkType := c.Query("clkType")
	point := c.Query("point")
	clk := c.Query("clk")
	volt := c.Query("volt")
	autoRespond, _ := strconv.ParseBool(c.Query("autoRespond"))

	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid request body"))
		return
	}

	failedMessage := dcgm.SetPowerPlayTableLevel(dvIdList, clkType, point, clk, volt, autoRespond)
	if len(failedMessage) > 0 {
		response := map[string]interface{}{
			"failedMessages": failedMessage,
		}
		c.JSON(http.StatusBadRequest, ErrorResponse(response))
	} else {
		c.JSON(http.StatusOK, SuccessResponse(nil))
	}
}

// SetClockOverDrive 处理设置时钟OverDrive
// @Summary 设置设备的时钟OverDrive
// @Description 设置设备的时钟OverDrive
// @Accept json
// @Produce json
// @Param dvIdList body []int true "设备ID列表"
// @Param clktype query string true "时钟类型（sclk 或 mclk）"
// @Param value query string true "OverDrive值，表示为百分比（0-20%）"
// @Param autoRespond query bool false "自动响应超出规格的警告"
// @Success 200 {string} string "成功设置时钟OverDrive"
// @Failure 400 {object} string "无效的请求参数或无法设置时钟OverDrive"
// @Router /SetClockOverDrive [post]
func SetClockOverDrive(c *gin.Context) {
	var dvIdList []int
	clktype := c.Query("clktype")
	value := c.Query("value")
	autoRespond, _ := strconv.ParseBool(c.Query("autoRespond"))

	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid request body"))
		return
	}

	failedMessage := dcgm.SetClockOverDrive(dvIdList, clktype, value, autoRespond)
	if len(failedMessage) > 0 {
		response := map[string]interface{}{
			"failedMessages": failedMessage,
		}
		c.JSON(http.StatusBadRequest, ErrorResponse(response))
	} else {
		c.JSON(http.StatusOK, SuccessResponse(nil))
	}
}

// setPerfDeterminism 处理设置性能确定性
// @Summary 设置设备的性能确定性
// @Description 设置设备的性能确定性
// @Accept json
// @Produce json
// @Param dvIdList body []int true "设备ID列表"
// @Param clkvalue query string true "时钟频率值"
// @Success 200 {array} FailedMessage "返回失败的设备及其错误信息"
// @Failure 400 {object} error "无效的请求体或无法设置性能确定性"
// @Router /SetPerfDeterminism [post]
func SetPerfDeterminism(c *gin.Context) {
	var dvIdList []int
	clkvalue := c.Query("clkvalue")

	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid request body"))
		return
	}

	failedMessages, err := dcgm.SetPerfDeterminism(dvIdList, clkvalue)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}

	if len(failedMessages) > 0 {
		c.JSON(http.StatusOK, failedMessages)
	} else {
		c.JSON(http.StatusOK, SuccessResponse(nil))
	}
}

// SetFanSpeed 设置风扇转速
// @Summary 设置风扇转速
// @Description 根据设备ID列表和给定的风扇速度，设置设备的风扇速度
// @Accept  json
// @Produce  json
// @Param dvIdList body []int true "设备 ID 列表"
// @Param fan query string true "风扇速度值（0-255,单位:RPM）"
// @Success 200 {string} string "成功信息"
// @Failure 400 {string} string "失败信息"
// @Router /SetFanSpeed [post]
func SetFanSpeed(c *gin.Context) {
	var dvIdList []int
	fan := c.Query("fan")

	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid request body"))
		return
	}

	dcgm.SetFanSpeed(dvIdList, fan)
	c.JSON(http.StatusOK, SuccessResponse(nil))
}

// DevFanRpms 获取设备的风扇速度
// @Summary 获取设备的风扇速度
// @Description 获取指定设备的风扇速度（RPM）
// @Accept  json
// @Produce  json
// @Param dvInd path int true "设备索引"
// @Success 200 {integer} int64 "风扇速度 (RPM)"
// @Failure 400 {string} string "失败信息"
// @Router /DevFanRpms/{dvInd} [get]
func DevFanRpms(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid device ID"))
		return
	}

	speed, err := dcgm.DevFanRpms(dvInd)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}
	devSpeed := map[string]interface{}{
		"speed": speed,
	}
	c.JSON(http.StatusOK, SuccessResponse(devSpeed))
}

// SetPerformanceLevel 设置设备性能等级
// @Summary 设置设备性能等级
// @Description 根据设备ID列表和给定的性能等级，设置设备的性能等级
// @Accept  json
// @Produce  json
// @Param deviceList body []int true "设备 ID 列表"
// @Param level query string true "性能等级 (auto, low, high, normal)"
// @Success 200 {array} FailedMessage
// @Failure 400 {object} FailedMessage
// @Router /SetPerformanceLevel [post]
func SetPerformanceLevel(c *gin.Context) {
	var deviceList []int
	level := c.Query("level")

	if err := c.BindJSON(&deviceList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid request body"))
		return
	}

	failedMessages := dcgm.SetPerformanceLevel(deviceList, level)
	if len(failedMessages) > 0 {
		c.JSON(http.StatusOK, failedMessages)
	} else {
		c.JSON(http.StatusOK, SuccessResponse(nil))
	}
}

// SetProfile 设置功率配置
// @Summary 设置功率配置
// @Description 根据设备ID列表和给定的功率配置文件，设置设备的功率配置
// @Accept  json
// @Produce  json
// @Param dvIdList body []int true "设备 ID 列表"
// @Param profile query string true "功率配置文件名称:CUSTOM、VIDEO、POWER SAVING、COMPUTE、VR、3D FULL SCREEN、BOOTUP DEFAULT"
// @Success 200 {array} FailedMessage "设置成功的消息列表"
// @Failure 400 {object} FailedMessage "失败的消息列表"
// @Router /SetProfile [post]
func SetProfile(c *gin.Context) {
	var dvIdList []int
	profile := c.Query("profile")

	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid request body"))
		return
	}

	failedMessages := dcgm.SetProfile(dvIdList, profile)
	c.JSON(http.StatusOK, failedMessages)
}

// DevPowerProfileSet 设置设备功率配置文件
// @Summary 设置设备功率配置文件
// @Description 设置指定设备的功率配置文件
// @Accept  json
// @Produce  json
// @Param dvInd path int true "设备索引"
// @Param reserved query int true "保留参数，通常为0"
// @Param profile query int true "功率配置文件的枚举值"
// @Success 200 {string} string "成功信息"
// @Failure 400 {string} string "失败信息"
// @Router /DevPowerProfileSet [post]
func DevPowerProfileSet(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid device ID"))
		return
	}

	reserved, err := strconv.Atoi(c.Query("reserved"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid reserved value"))
		return
	}

	profileEnum, err := strconv.Atoi(c.Query("profile"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid profile value"))
		return
	}

	err = dcgm.DevPowerProfileSet(dvInd, reserved, dcgm.RSNIPowerProfilePresetMasks(profileEnum))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(nil))
}

// GetBus 获取设备总线信息
// @Summary 获取设备总线信息
// @Description 获取指定设备的总线信息
// @Accept  json
// @Produce  json
// @Param dvInd path int true "设备索引"
// @Success 200 {string} string "设备总线ID"
// @Failure 400 {string} string "失败信息"
// @Router /GetBus/{device} [get]
func GetBus(c *gin.Context) {
	device, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid device ID"))
		return
	}

	picId, err := dcgm.GetBus(device)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}
	devPicId := map[string]interface{}{
		"picId": picId,
	}

	c.JSON(http.StatusOK, SuccessResponse(devPicId))
}

// ShowAllConciseHw 显示设备硬件信息
// @Summary 显示设备硬件信息
// @Description 显示指定设备列表的简要硬件信息
// @Accept  json
// @Produce  json
// @Param dvIdList body []int true "设备 ID 列表"
// @Success 200 {string} string "成功信息"
// @Failure 400 {string} string "失败信息"
// @Router /ShowAllConciseHw [post]
func ShowAllConciseHw(c *gin.Context) {
	var dvIdList []int

	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid request body"))
		return
	}

	dcgm.ShowAllConciseHw(dvIdList)
	c.JSON(http.StatusOK, SuccessResponse(nil))
}

// ShowClocks 显示时钟信息
// @Summary 显示时钟信息
// @Description 显示指定设备的时钟信息
// @Accept  json
// @Produce  json
// @Param dvIdList body []int true "设备 ID 列表"
// @Success 200 {string} string "成功信息"
// @Failure 400 {string} string "失败信息"
// @Router /ShowClocks [post]
func ShowClocks(c *gin.Context) {
	var dvIdList []int

	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid request body"))
		return
	}

	dcgm.ShowClocks(dvIdList)
	c.JSON(http.StatusOK, SuccessResponse(nil))
}

// ShowCurrentFans 展示风扇转速和风扇级别
// @Summary 展示风扇转速和风扇级别
// @Description 显示指定设备的当前风扇转速和风扇级别
// @Accept  json
// @Produce  json
// @Param dvIdList body []int true "设备 ID 列表"
// @Success 200 {string} string "成功信息"
// @Failure 400 {string} string "失败信息"
// @Router /fans/current [post]
func ShowCurrentFans(c *gin.Context) {
	var dvIdList []int
	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid request body"))
		return
	}

	dcgm.ShowCurrentFans(dvIdList, true)
	c.JSON(http.StatusOK, SuccessResponse(nil))
}

// ShowCurrentTemps 显示设备温度传感器数据
// @Summary 显示设备温度传感器数据
// @Accept  json
// @Produce  json
// @Param dvIdList body []int true "设备 ID 列表"
// @Success 200 {object} TemperatureInfo "温度信息列表"
// @Failure 400 {object} error "错误信息"
// @Router /temps/current [post]
func ShowCurrentTemps(c *gin.Context) {
	var dvIdList []int
	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid request body"))
		return
	}

	temperatureInfo, err := dcgm.ShowCurrentTemps(dvIdList)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}
	temperatureInfos := map[string]interface{}{
		"temperatureInfos": temperatureInfo,
	}
	c.JSON(http.StatusOK, SuccessResponse(temperatureInfos))
}

// ShowFwInfo 显示设备固件版本信息
// @Summary 显示设备固件版本信息
// @Param dvIdList query []int true "设备 ID 列表"
// @Param fwType query []string true "固件类型列表"
// @Success 200 {object} []FirmwareInfo "固件版本信息列表"
// @Failure 400 {object} error "错误信息"
// @Router /firmware/info [get]
func ShowFwInfo(c *gin.Context) {
	var dvIdList []int
	var fwType []string

	if err := c.BindQuery(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid query parameters"))
		return
	}

	if err := c.BindQuery(&fwType); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid query parameters"))
		return
	}

	fwInfos, err := dcgm.ShowFwInfo(dvIdList, fwType)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}
	fwInfo := map[string]interface{}{
		"fwInfos": fwInfos,
	}
	c.JSON(http.StatusOK, SuccessResponse(fwInfo))
}

// PidList 获取计算进程列表
// @Summary 获取计算进程列表
// @Success 200 {array} string "进程 ID 列表"
// @Failure 400 {object} error "错误信息"
// @Router /process/list [get]
func PidList(c *gin.Context) {
	pidList, err := dcgm.PidList()
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}
	pidLists := map[string]interface{}{
		"pidList": pidList,
	}
	c.JSON(http.StatusOK, SuccessResponse(pidLists))
}

// GetCoarseGrainUtil 获取设备粗粒度利用率
// @Summary 获取设备粗粒度利用率
// @Param device body int true "设备 ID"
// @Param typeName body string false "利用率计数器类型名称"
// @Success 200 {array} RSMIUtilizationCounter "利用率计数器列表"
// @Failure 400 {object} error "错误信息"
// @Router /utilization/coarse [post]
func GetCoarseGrainUtil(c *gin.Context) {
	var request struct {
		Device   int    `json:"device"`
		TypeName string `json:"typeName"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid JSON body"))
		return
	}

	// 直接传递 TypeName 的值
	utilizationCounters, err := dcgm.GetCoarseGrainUtil(request.Device, request.TypeName)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}

	response := map[string]interface{}{
		"utilizationCounters": utilizationCounters,
	}
	c.JSON(http.StatusOK, SuccessResponse(response))
}

// ShowGpuUse 显示设备的 DCU 使用率
// @Summary 显示设备的 DCU 使用率
// @Param dvIdList body []int true "设备 ID 列表"
// @Success 200 {object} []DeviceUseInfo "设备使用信息列表"
// @Failure 400 {object} error "错误信息"
// @Router /gpu/use [post]
func ShowGpuUse(c *gin.Context) {
	var dvIdList []int
	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid JSON body"))
		return
	}
	deviceUseInfos, err := dcgm.ShowGpuUse(dvIdList)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}
	response := map[string]interface{}{
		"deviceUseInfos": deviceUseInfos,
	}
	c.JSON(http.StatusOK, SuccessResponse(response))
}

// ShowEnergy 展示设备消耗的能量
// @Summary 展示设备的能量消耗
// @Description 获取并展示指定设备的能量消耗情况。
// @Param dvIdList body []int true "设备 ID 列表"
// @Success 200 {string} string "成功返回设备的能量消耗信息"
// @Failure 400 {string} string "请求参数错误"
// @Router /energy [post]
func ShowEnergy(c *gin.Context) {
	var dvIdList []int
	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid JSON body"))
		return
	}

	dcgm.ShowEnergy(dvIdList)
	c.JSON(http.StatusOK, SuccessResponse(nil))
}

// ShowMemInfo 展示设备的内存信息
// @Summary 展示设备内存信息
// @Description 获取并展示指定设备的内存使用情况，包括不同类型的内存。
// @Param dvIdList body []int true "设备 ID 列表"
// @Param memTypes body []string true "内存类型列表，如 'all' 或指定类型"
// @Success 200 {string} string "成功返回设备的内存信息"
// @Failure 400 {string} string "请求参数错误"
// @Router /memory/info [post]
func ShowMemInfo(c *gin.Context) {
	var request struct {
		DvIdList []int    `json:"dvIdList"`
		MemTypes []string `json:"memTypes"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid JSON body"))
		return
	}

	dcgm.ShowMemInfo(request.DvIdList, request.MemTypes)
	c.JSON(http.StatusOK, SuccessResponse(nil))
}

// ShowMemUse 展示设备的内存使用情况
// @Summary 展示设备内存使用情况
// @Description 获取并展示指定设备的当前内存使用百分比和其他相关的利用率数据。
// @Param dvIdList body []int true "设备 ID 列表"
// @Success 200 {string} string "成功返回设备的内存使用信息"
// @Failure 400 {string} string "请求参数错误"
// @Router /memory/use [post]
func ShowMemUse(c *gin.Context) {
	var dvIdList []int
	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid JSON body"))
		return
	}

	dcgm.ShowMemUse(dvIdList)
	c.JSON(http.StatusOK, SuccessResponse(nil))
}

// ShowMemVendor 展示设备供应商信息
// @Summary 展示设备的内存供应商信息
// @Description 获取并展示指定设备的内存供应商信息。
// @Param dvIdList body []int true "设备 ID 列表"
// @Success 200 {object} []DeviceMemVendorInfo "成功返回设备的内存供应商信息"
// @Failure 400 {string} string "请求参数错误"
// @Router /memory/vendor [post]
func ShowMemVendor(c *gin.Context) {
	var dvIdList []int
	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid JSON body"))
		return
	}

	deviceMemVendorInfos, err := dcgm.ShowMemVendor(dvIdList)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}
	response := map[string]interface{}{
		"deviceMemVendorInfos": deviceMemVendorInfos,
	}
	c.JSON(http.StatusOK, SuccessResponse(response))
}

// ShowPcieBw 展示设备的PCIe带宽使用情况
// @Summary 展示设备的PCIe带宽使用情况
// @Description 获取并展示指定设备的PCIe带宽使用情况，包括发送和接收的带宽。
// @Param dvIdList body []int true "设备 ID 列表"
// @Success 200 {object} []PcieBandwidthInfo "成功返回设备的PCIe带宽使用信息"
// @Failure 400 {string} string "请求参数错误"
// @Router /pcie/bandwidth [post]
func ShowPcieBw(c *gin.Context) {
	var dvIdList []int
	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid JSON body"))
		return
	}

	pcieBandwidthInfos, err := dcgm.ShowPcieBw(dvIdList)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}
	response := map[string]interface{}{
		"pcieBandwidthInfos": pcieBandwidthInfos,
	}

	c.JSON(http.StatusOK, SuccessResponse(response))
}

// ShowPcieReplayCount 展示设备的PCIe重放计数
// @Summary 展示设备的PCIe重放计数
// @Description 获取并展示指定设备的PCIe重放计数。
// @Param dvIdList body []int true "设备 ID 列表"
// @Success 200 {object} []PcieReplayCountInfo "设备的PCIe重放计数信息"
// @Failure 400 {string} string "请求参数错误"
// @Router /pcie/replaycount [post]
func ShowPcieReplayCount(c *gin.Context) {
	var dvIdList []int
	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid JSON body"))
		return
	}

	pcieReplayCountInfos, err := dcgm.ShowPcieReplayCount(dvIdList)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}
	response := map[string]interface{}{
		"pcieReplayCountInfos": pcieReplayCountInfos,
	}

	c.JSON(http.StatusOK, SuccessResponse(response))
}

// ShowPids 展示进程信息
// @Summary 展示系统中正在运行的KFD进程信息
// @Description 获取并展示当前系统中运行的KFD进程的详细信息。
// @Success 200 {string} string "成功返回进程信息"
// @Failure 400 {string} string "请求错误"
// @Router /pids [get]
func ShowPids(c *gin.Context) {
	err := dcgm.ShowPids()
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(nil))
}

// @Summary 展示设备的平均功率消耗
// @Description 获取并展示指定设备的平均图形功率消耗
// @Param dvIdList body []int true "设备 ID 列表"
// @Success 200 {object} []DevicePowerInfo "设备的功率信息"
// @Failure 400 {string} string "请求参数错误"
// @Router /device/power [post]
func GetDevicePower(c *gin.Context) {
	var dvIdList []int
	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid JSON body"))
		return
	}
	devicePowerInfos, err := dcgm.ShowPower(dvIdList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	response := map[string]interface{}{
		"devicePowerInfos": devicePowerInfos,
	}
	c.JSON(http.StatusOK, SuccessResponse(response))
}

// @Summary 展示设备的GPU内存时钟频率和电压
// @Description 获取并展示指定设备的GPU内存时钟频率和电压表
// @Param dvIdList body []int true "设备 ID 列表"
// @Success 200 {object} []DevicePowerPlayInfo "设备的GPU时钟频率和电压信息"
// @Failure 400 {string} string "请求参数错误"
// @Router /device/powerplay [post]
func GetDevicePowerPlayTable(c *gin.Context) {
	var dvIdList []int
	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid JSON body"))
		return
	}
	devicePowerPlayInfos, err := dcgm.ShowPowerPlayTable(dvIdList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	response := map[string]interface{}{
		"devicePowerPlayInfos": devicePowerPlayInfos,
	}
	c.JSON(http.StatusOK, SuccessResponse(response))
}

// @Summary 显示设备的产品名称
// @Description 获取并显示指定设备的产品名称、供应商、系列、型号和SKU信息
// @Param dvIdList body []int true "设备 ID 列表"
// @Success 200 {array} DeviceproductInfo "设备的产品信息列表"
// @Failure 400 {string} string "请求参数错误"
// @Router /device/product [post]
func GetDeviceProductName(c *gin.Context) {
	var dvIdList []int
	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid JSON body"))
		return
	}
	deviceProductInfos, err := dcgm.ShowProductName(dvIdList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	response := map[string]interface{}{
		"deviceProductInfos": deviceProductInfos,
	}
	c.JSON(http.StatusOK, SuccessResponse(response))
}

// @Summary 显示设备的电源配置文件
// @Description 获取并显示指定设备的电源配置文件，包括可用的电源配置选项
// @Param dvIdList body []int true "设备 ID 列表"
// @Success 200 {array} DeviceProfile "设备的电源配置文件信息列表"
// @Failure 400 {string} string "请求参数错误"
// @Router /device/profile [post]
func GetDeviceProfile(c *gin.Context) {
	var dvIdList []int
	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid JSON body"))
		return
	}
	deviceProfiles, err := dcgm.ShowProfile(dvIdList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	response := map[string]interface{}{
		"deviceProfiles": deviceProfiles,
	}
	c.JSON(http.StatusOK, SuccessResponse(response))
}

// GetDeviceRange 显示设备的电流或电压范围
// @Summary 显示设备的电流或电压范围（K100_AI卡不支持该操作）
// @Description 获取并显示指定设备的有效电流或电压范围
// @Param dvIdList body []int true "设备ID列表"
// @Param rangeType body string true "范围类型 (sclk, mclk, voltage)"
// @Success 200 {string} string "设备的电流或电压范围信息"
// @Failure 400 {string} string "请求参数错误"
// @Router /device/range [post]
func GetDeviceRange(c *gin.Context) {
	var request struct {
		DvIdList  []int  `json:"dvIdList"`
		RangeType string `json:"rangeType"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid JSON body"))
		return
	}

	dcgm.ShowRange(request.DvIdList, request.RangeType)
	c.JSON(http.StatusOK, SuccessResponse(nil))
}

// @Summary 显示设备的退役页信息
// @Description 获取并显示指定设备的退役内存页信息
// @Param dvIdList body []int true "设备ID列表"
// @Param retiredType body string false "退役类型 (默认为'all')"
// @Success 200 {string} string "设备的退役页信息"
// @Failure 400 {string} string "请求参数错误"
// @Router /device/retiredpages [post]
func GetDeviceRetiredPages(c *gin.Context) {
	var request struct {
		DvIdList    []int  `json:"dvIdList"`
		RetiredType string `json:"retiredType"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid JSON body"))
		return
	}
	dcgm.ShowRetiredPages(request.DvIdList, request.RetiredType)
	c.JSON(http.StatusOK, SuccessResponse(nil))
}

// @Summary 显示设备的序列号
// @Description 获取并显示指定设备的序列号信息
// @Param dvIdList body []int true "设备ID列表"
// @Success 200 {array} DeviceSerialInfo "设备的序列号信息列表"
// @Failure 400 {string} string "请求参数错误"
// @Router /device/serialnumber [post]
func GetDeviceSerialNumber(c *gin.Context) {
	var dvIdList []int
	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid JSON body"))
		return
	}
	deviceSerialInfos, err := dcgm.ShowSerialNumber(dvIdList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	response := map[string]interface{}{
		"deviceSerialInfos": deviceSerialInfos,
	}
	c.JSON(http.StatusOK, SuccessResponse(response))
}

// ShowUId 显示设备的唯一ID
// @Summary 显示设备的唯一ID
// @Description 获取并显示指定设备的唯一ID信息。
// @Param dvIdList body []int true "设备ID列表"
// @Success 200 {array} DeviceUIdInfo "设备的唯一ID信息列表"
// @Failure 400 {string} string "请求参数错误"
// @Router /showUId [post]
func ShowUId(c *gin.Context) {
	var dvIdList []int
	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid JSON body"))
		return
	}

	deviceUIdInfos, err := dcgm.ShowUId(dvIdList)
	if err != nil {
		c.String(http.StatusInternalServerError, "服务器内部错误")
		return
	}
	response := map[string]interface{}{
		"deviceUIdInfos": deviceUIdInfos,
	}
	c.JSON(http.StatusOK, SuccessResponse(response))
}

// ShowVbiosVersion 显示设备的VBIOS版本
// @Summary 显示设备的VBIOS版本
// @Description 获取并显示指定设备的VBIOS版本信息。
// @Param dvIdList body []int true "设备ID列表"
// @Success 200 {array} DeviceVBIOSInfo "设备的VBIOS版本信息列表"
// @Failure 400 {string} string "请求参数错误"
// @Router /showVbiosVersion [post]
func ShowVbiosVersion(c *gin.Context) {
	var dvIdList []int
	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid JSON body"))
		return
	}

	deviceVBIOSInfos, err := dcgm.ShowVbiosVersion(dvIdList)
	if err != nil {
		c.String(http.StatusInternalServerError, "服务器内部错误")
		return
	}
	response := map[string]interface{}{
		"deviceVBIOSInfos": deviceVBIOSInfos,
	}
	c.JSON(http.StatusOK, SuccessResponse(response))
}

// ShowEvents 显示设备的事件
// @Summary 显示设备的事件
// @Description 获取并显示指定设备的事件信息。
// @Param dvIdList body []int true "设备ID列表"
// @Param eventTypes body []string true "事件类型列表"
// @Success 200 {string} string "成功返回设备的事件信息"
// @Failure 400 {string} string "请求参数错误"
// @Router /showEvents [post]
func ShowEvents(c *gin.Context) {
	var requestData struct {
		DvIdList   []int    `json:"dvIdList"`
		EventTypes []string `json:"eventTypes"`
	}

	if err := c.BindJSON(&requestData); err != nil {
		c.String(http.StatusBadRequest, "请求参数错误")
		return
	}

	dcgm.ShowEvents(requestData.DvIdList, requestData.EventTypes)
	c.String(http.StatusOK, "设备事件信息已显示")
}

// ShowVoltage 显示设备的电压信息
// @Summary 显示设备的电压信息
// @Description 获取并显示指定设备的当前电压信息。
// @Param dvIdList body []int true "设备ID列表"
// @Success 200 {array} DeviceVoltageInfo "设备的电压信息列表"
// @Failure 400 {string} string "请求参数错误"
// @Router /showVoltage [post]
func ShowVoltage(c *gin.Context) {
	var dvIdList []int
	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid JSON body"))
		return
	}

	deviceVoltageInfos, err := dcgm.ShowVoltage(dvIdList)
	if err != nil {
		c.String(http.StatusInternalServerError, "服务器内部错误")
		return
	}
	response := map[string]interface{}{
		"deviceVoltageInfos": deviceVoltageInfos,
	}
	c.JSON(http.StatusOK, SuccessResponse(response))
}

// ShowVoltageCurve 显示设备的电压曲线点
// @Summary 显示设备的电压曲线点
// @Description 获取并显示指定设备的电压曲线点信息。
// @Param dvIdList body []int true "设备ID列表"
// @Success 200 {string} string "设备的电压曲线点信息"
// @Failure 400 {string} string "请求参数错误"
// @Router /showVoltageCurve [post]
func ShowVoltageCurve(c *gin.Context) {
	var dvIdList []int
	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid JSON body"))
		return
	}

	dcgm.ShowVoltageCurve(dvIdList)
	c.String(http.StatusOK, "设备电压曲线点信息已显示")
}

// ShowXgmiErr 显示 XGMI 错误状态
// @Summary 显示 XGMI 错误状态
// @Description 显示一组 GPU 设备的 XGMI 错误状态。
// @Param dvIdList body []int true "设备 ID 列表"
// @Success 200 {string} string "XGMI 错误状态信息"
// @Router /showXgmiErr [post]
func ShowXgmiErr(c *gin.Context) {
	var dvIdList []int
	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid JSON body"))
		return
	}

	dcgm.ShowXgmiErr(dvIdList, true)
	c.String(http.StatusOK, "XGMI 错误状态信息已显示")
}

// ShowWeightTopology 显示 GPU 拓扑权重
// @Summary 显示 GPU 拓扑权重
// @Description 显示 GPU 设备之间的权重信息。
// @Param dvIdList body []int true "设备 ID 列表"
// @Success 200 {string} string "GPU 拓扑权重信息"
// @Router /showWeightTopology [post]
func ShowWeightTopology(c *gin.Context) {
	var dvIdList []int
	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid JSON body"))
		return
	}
	dcgm.ShowWeightTopology(dvIdList, true)
	c.String(http.StatusOK, "GPU 拓扑权重信息已显示")
}

// ShowHopsTopology 显示 GPU 拓扑跳数
// @Summary 显示 GPU 拓扑跳数
// @Description 显示 GPU 设备之间的跳数信息。
// @Param dvIdList body []int true "设备 ID 列表"
// @Success 200 {string} string "GPU 拓扑跳数信息"
// @Router /showHopsTopology [post]
func ShowHopsTopology(c *gin.Context) {
	var dvIdList []int
	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid JSON body"))
		return
	}

	dcgm.ShowHopsTopology(dvIdList, true)
	c.String(http.StatusOK, "GPU 拓扑跳数信息已显示")
}

// ShowTypeTopology 显示 GPU 拓扑中两台设备之间的链接类型。
// @Summary 显示 GPU 拓扑链接类型
// @Description 显示 GPU 设备之间的链接类型信息。
// @Param dvIdList body []int true "设备 ID 列表"
// @Success 200 {string} string "GPU 拓扑链接类型信息"
// @Router /showTypeTopology [post]
func ShowTypeTopology(c *gin.Context) {
	var dvIdList []int
	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid JSON body"))
		return
	}
	dcgm.ShowTypeTopology(dvIdList, true)
	c.String(http.StatusOK, " GPU拓扑中两台设备之间的链接类型已展示")
}

// ShowNumaTopology 显示指定设备的 NUMA 节点信息。
// @Summary 显示 NUMA 节点信息
// @Description 显示一组 GPU 设备的 NUMA 节点和关联信息。
// @Param dvIdList body []int true "设备 ID 列表"
// @Success 200 {string} string "NUMA 节点信息"
// @Router /showNumaTopology [post]
func ShowNumaTopology(c *gin.Context) {
	var dvIdList []int
	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid JSON body"))
		return
	}

	infos, _ := dcgm.ShowNumaTopology(dvIdList)
	response := map[string]interface{}{
		"infos": infos,
	}
	c.JSON(http.StatusOK, SuccessResponse(response))
}

// ShowHwTopology 显示指定设备的完整硬件拓扑信息。
// @Summary 显示完整的硬件拓扑信息
// @Description 显示一组 GPU 设备的权重、跳数、链接类型和 NUMA 节点信息。
// @Param dvIdList body []int true "设备 ID 列表"
// @Success 200 {string} string "完整的硬件拓扑信息"
// @Router /showHwTopology [post]
func ShowHwTopology(c *gin.Context) {
	var dvIdList []int
	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid JSON body"))
		return
	}

	dcgm.ShowHwTopology(dvIdList)
	c.String(http.StatusOK, "指定设备的完整硬件拓扑信息已显示")
}

// DeviceCount 返回设备的数量
// @Summary 获取设备数量
// @Description 获取当前系统中的设备数量
// @Success 200 {int} int "设备数量"
// @Failure 500 {object} string "内部服务器错误"
// @Router /deviceCount [get]
func DeviceCount(c *gin.Context) {
	count, err := dcgm.DeviceCount()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	response := map[string]interface{}{
		"count": count,
	}
	c.JSON(http.StatusOK, SuccessResponse(response))
}

// VDeviceSingleInfo 获取单个虚拟设备的信息
// @Summary 获取单个虚拟设备的信息
// @Description 根据设备索引获取对应的虚拟设备信息
// @Param vDvInd query int true "设备索引"
// @Success 200 {object} DMIVDeviceInfo "虚拟设备信息"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "内部服务器错误"
// @Router /VDeviceSingleInfo [get]
func VDeviceSingleInfo(c *gin.Context) {
	var vDvInd int
	if err := c.BindQuery(&vDvInd); err != nil {
		c.JSON(http.StatusBadRequest, "虚拟设备销毁失败")
		return
	}
	vDeviceInfo, err := dcgm.VDeviceSingleInfo(vDvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	response := map[string]interface{}{
		"vDeviceInfo": vDeviceInfo,
	}
	c.JSON(http.StatusOK, SuccessResponse(response))
}

// VDeviceCount 返回虚拟设备的数量
// @Summary 获取虚拟设备数量
// @Description 获取当前系统中的虚拟设备数量
// @Success 200 {int} int "虚拟设备数量"
// @Failure 500 {object} string "内部服务器错误"
// @Router /vDeviceCount [get]
func VDeviceCount(c *gin.Context) {
	vcount, err := dcgm.VDeviceCount()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	response := map[string]interface{}{
		"vcount": vcount,
	}
	c.JSON(http.StatusOK, SuccessResponse(response))
}

// DeviceRemainingInfo 返回指定物理设备的剩余计算单元（CU）和内存信息
// @Summary 获取设备剩余信息
// @Description 获取指定设备的剩余计算单元和内存信息
// @Param dvInd path int true "物理设备索引"
// @Success 200 {string} uint64 "剩余的CU信息"
// @Success 200 {string} uint64 "剩余的内存信息"
// @Failure 400 {object} string "无效的设备索引"
// @Failure 500 {object} string "内部服务器错误"
// @Router /deviceRemainingInfo/{dvInd} [get]
func DeviceRemainingInfo(c *gin.Context) {
	// 获取参数并转换为整数类型
	dvIndStr := c.Param("dvInd")
	dvInd, err := strconv.Atoi(dvIndStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, "无效的设备索引")
		return
	}

	cus, memories, err := dcgm.DeviceRemainingInfo(dvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	devRemainingInfo := map[string]interface{}{
		"cus":      cus,
		"memories": memories,
	}
	c.JSON(http.StatusOK, SuccessResponse(devRemainingInfo))
}

// CreateVDevices 创建指定数量的虚拟设备
// @Summary 创建虚拟设备
// @Description 在指定的物理设备上创建指定数量的虚拟设备，返回创建的虚拟设备ID集合
// @Param dvInd query int true "物理设备的索引"
// @Param vDevCount query int true "要创建的虚拟设备数量"
// @Param vDevCUs query []int true "每个虚拟设备的计算单元数量，多个值使用多个 vDevCUs 参数传递，例如：vDevCUs=10&vDevCUs=20"
// @Param vDevMemSize query []int true "每个虚拟设备的内存大小，多个值使用多个 vDevMemSize 参数传递，例如：vDevMemSize=1024&vDevMemSize=2048"
// @Success 200 {object} map[string]interface{} "虚拟设备创建成功，返回包含虚拟设备ID集合的对象"
// @Failure 400 {string} string "请求参数无效"
// @Failure 500 {string} string "创建虚拟设备失败"
// @Router /CreateVDevices [post]
func CreateVDevices(c *gin.Context) {
	var err error
	// 获取单个整数参数
	dvInd, err := strconv.Atoi(c.Query("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "无效的物理设备索引")
		return
	}

	vDevCount, err := strconv.Atoi(c.Query("vDevCount"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "无效的虚拟设备数量")
		return
	}

	// 获取切片参数
	vDevCUsStr := c.QueryArray("vDevCUs")
	vDevMemSizeStr := c.QueryArray("vDevMemSize")

	// 将字符串切片转换为整数切片
	vDevCUs := make([]int, len(vDevCUsStr))
	vDevMemSize := make([]int, len(vDevMemSizeStr))

	for i, v := range vDevCUsStr {
		vDevCUs[i], err = strconv.Atoi(v)
		if err != nil {
			c.JSON(http.StatusBadRequest, "无效的计算单元数量")
			return
		}
	}

	for i, v := range vDevMemSizeStr {
		vDevMemSize[i], err = strconv.Atoi(v)
		if err != nil {
			c.JSON(http.StatusBadRequest, "无效的内存大小")
			return
		}
	}

	// 调用业务逻辑层的函数
	vdevIDs, err := dcgm.CreateVDevices(dvInd, vDevCount, vDevCUs, vDevMemSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "创建虚拟设备失败")
		return
	}

	// 成功响应
	response := map[string]interface{}{
		"vdevIDs": vdevIDs,
	}
	c.JSON(http.StatusOK, SuccessResponse(response))
}

// DestroyVDevice 销毁指定物理设备上的所有虚拟设备
// @Summary 销毁所有虚拟设备
// @Description 销毁指定物理设备上的所有虚拟设备
// @Param dvInd query int true "物理设备的索引"
// @Success 200 {string} string "虚拟设备销毁成功"
// @Failure 400 {string} string "虚拟设备销毁失败"
// @Router /DestroyVDevice [delete]
func DestroyVDevice(c *gin.Context) {
	var dvInd int
	if err := c.BindQuery(&dvInd); err != nil {
		c.JSON(http.StatusBadRequest, "虚拟设备销毁失败")
		return
	}
	if err := dcgm.DestroyVDevice(dvInd); err != nil {
		c.JSON(http.StatusInternalServerError, "虚拟设备销毁失败:"+err.Error())
		return
	}
	c.JSON(http.StatusOK, "虚拟设备销毁成功")
}

// DestroySingleVDevice 销毁指定虚拟设备
// @Summary 销毁单个虚拟设备
// @Description 销毁指定索引的虚拟设备
// @Param vDvInd query int true "虚拟设备的索引"
// @Success 200 {string} string "虚拟设备销毁成功"
// @Failure 400 {string} string "虚拟设备销毁失败"
// @Router /DestroySingleVDevice [delete]
func DestroySingleVDevice(c *gin.Context) {
	var vDvInd int
	if err := c.BindQuery(&vDvInd); err != nil {
		c.JSON(http.StatusBadRequest, "虚拟设备销毁失败")
		return
	}
	if err := dcgm.DestroySingleVDevice(vDvInd); err != nil {
		c.JSON(http.StatusInternalServerError, "虚拟设备销毁失败:"+err.Error())
		return
	}
	c.JSON(http.StatusOK, "虚拟设备销毁成功")
}

// UpdateSingleVDevice 更新指定设备资源大小
// @Summary 更新虚拟设备资源
// @Description 更新指定虚拟设备的计算单元和内存大小
// @Param vDvInd query int true "虚拟设备的索引"
// @Param vDevCUs query int true "更新后的计算单元数量"
// @Param vDevMemSize query int true "更新后的内存大小"
// @Success 200 {string} string "虚拟设备更新成功"
// @Failure 400 {string} string "虚拟设备更新失败"
// @Router /UpdateSingleVDevice [put]
func UpdateSingleVDevice(c *gin.Context) {
	var vDvInd, vDevCUs, vDevMemSize int
	if err := c.ShouldBindQuery(&vDvInd); err != nil {
		c.JSON(http.StatusBadRequest, "虚拟设备更新失败")
		return
	}
	if err := c.ShouldBindQuery(&vDevCUs); err != nil {
		c.JSON(http.StatusBadRequest, "虚拟设备更新失败")
		return
	}
	if err := c.ShouldBindQuery(&vDevMemSize); err != nil {
		c.JSON(http.StatusBadRequest, "虚拟设备更新失败")
		return
	}
	if err := dcgm.UpdateSingleVDevice(vDvInd, vDevCUs, vDevMemSize); err != nil {
		c.JSON(http.StatusInternalServerError, "虚拟设备更新失败")
		return
	}
	c.JSON(http.StatusOK, "虚拟设备更新成功")
}

// 启动虚拟设备
// @Summary 启动指定的虚拟设备
// @Description 启动虚拟设备，指定设备索引
// @Param vDvInd path int true "虚拟设备索引"
// @Success 200 {string} string "操作成功"
// @Failure 400 {string} string "操作失败"
// @Router /StartVDevice/{vDvInd} [get]
func StartVDevice(c *gin.Context) {
	vDvInd, err := strconv.Atoi(c.Param("vDvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "请求参数错误")
		return
	}

	if err := dcgm.StartVDevice(vDvInd); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, "启动成功")
}

// 停止虚拟设备
// @Summary 停止指定的虚拟设备
// @Description 停止虚拟设备，指定设备索引
// @Param vDvInd path int true "虚拟设备索引"
// @Success 200 {string} string "操作成功"
// @Failure 400 {string} string "操作失败"
// @Router /StopVDevice/{vDvInd} [get]
func StopVDevice(c *gin.Context) {
	vDvInd, err := strconv.Atoi(c.Param("vDvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "请求参数错误")
		return
	}

	if err := dcgm.StopVDevice(vDvInd); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, "停止成功")
}

// 设置虚拟机加密状态
// @Summary 设置虚拟机加密状态
// @Description 根据提供的状态开启或关闭虚拟机加密
// @Param status query bool true "加密状态"
// @Success 200 {string} string "操作成功"
// @Failure 400 {string} string "操作失败"
// @Router /SetEncryptionVMStatus [post]
func SetEncryptionVMStatus(c *gin.Context) {
	var status bool
	if err := c.BindQuery(&status); err != nil {
		c.JSON(http.StatusBadRequest, "请求参数错误")
		return
	}

	if err := dcgm.SetEncryptionVMStatus(status); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, "设置成功")
}

// 获取加密虚拟机状态
// @Summary 获取当前虚拟机的加密状态
// @Description 返回虚拟机是否处于加密状态
// @Success 200 {boolean} boolean "加密状态"
// @Failure 400 {string} string "操作失败"
// @Router /EncryptionVMStatus [get]
func EncryptionVMStatus(c *gin.Context) {
	status, err := dcgm.EncryptionVMStatus()
	if err != nil {
		c.JSON(http.StatusBadRequest, "操作失败")
		return
	}
	response := map[string]interface{}{
		"status": status,
	}
	c.JSON(http.StatusOK, SuccessResponse(response))
}

// 打印事件列表
// @Summary 打印设备的事件列表
// @Description 打印指定设备的事件列表，并设置延迟
// @Param device path int true "设备索引"
// @Param delay query int true "延迟时间（秒）"
// @Param eventList query []string true "事件列表"
// @Success 200 {string} string "操作成功"
// @Failure 400 {string} string "操作失败"
// @Router /PrintEventList/{device} [get]
func PrintEventList(c *gin.Context) {
	device, err := strconv.Atoi(c.Param("device"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "请求参数错误")
		return
	}

	var delay int
	var eventList []string
	if err := c.BindQuery(&delay); err != nil || c.BindQuery(&eventList) != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	dcgm.PrintEventList(device, delay, eventList)
	c.JSON(http.StatusOK, "操作成功")
}

// GetDeviceInfo 获取设备信息
// @Summary 获取设备信息
// @Description 根据设备索引获取对应的设备信息
// @Param dvInd path int true "设备索引"
// @Success 200 {object} DMIDeviceInfo "设备信息"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "内部服务器错误"
// @Router /device/info/{dvInd} [get]
func GetDeviceInfo(c *gin.Context) {
	var dvInd int
	if err := c.BindUri(&dvInd); err != nil {
		c.JSON(http.StatusBadRequest, "请求参数错误")
		return
	}

	deviceInfo, err := dcgm.GetDeviceInfo(dvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "内部服务器错误")
		return
	}
	response := map[string]interface{}{
		"deviceInfo": deviceInfo,
	}
	c.JSON(http.StatusOK, SuccessResponse(response))
}

// DeviceControl 处理设备控制
// @Summary 控制设备的性能级别、时钟频率和风扇重置
// @Description 根据传入的设备控制信息，设置设备的性能级别、时钟频率，并可选择性重置风扇
// @Accept json
// @Produce json
// @Param deviceControl body DeviceControlInfo true "设备控制信息"
// @Success 200 {object} string "成功返回操作结果"
// @Failure 400 {object} string "无效的请求参数或操作失败"
// @Failure 500 {object} string "内部服务器错误"
// @Router /device/control [post]
func DeviceControl(c *gin.Context) {
	var deviceInfo DeviceControlInfo
	var validationErrors []string // 存储所有验证失败的字段和原因
	// 绑定 JSON 请求体到 deviceInfo 对象
	if err := c.ShouldBindJSON(&deviceInfo); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}

	// 验证 DvInd 参数
	dvInd := deviceInfo.DvInd
	if dvInd < 0 {
		validationErrors = append(validationErrors, "无效的 DvInd")
	}
	// 验证 Level 参数
	var levelConverted dcgm.RSMIDevPerfLevel
	if deviceInfo.PerfLevel != "" {
		var err error
		levelConverted, err = ConvertToRSMIDevPerfLevel(deviceInfo.PerfLevel)
		if err != nil {
			validationErrors = append(validationErrors, "无效的性能级别 Level")
		}
	}

	// 验证 SclkClock 参数
	var sclkClock int64
	if deviceInfo.SclkClock != "" {
		var err error
		sclkClock, err = ConvertFrequencyToSclkClock(deviceInfo.SclkClock)
		if err != nil {
			validationErrors = append(validationErrors, "无效的 SclkClock 参数："+err.Error())
		}
	}

	// 验证 SocclkClock 参数
	var socclkClock int64
	if deviceInfo.SocclkClock != "" {
		var err error
		socclkClock, err = ConvertFrequencyToSocclkClock(deviceInfo.SocclkClock)
		if err != nil {
			validationErrors = append(validationErrors, "无效的 SocclkClock 参数："+err.Error())
		}
	}

	// 如果有任何验证失败的参数，返回错误
	if len(validationErrors) > 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse("以下参数无效: "+strings.Join(validationErrors, ", ")))
		return
	}
	// 执行方法逻辑并收集错误
	var executionErrors []string

	// 调用 DevPerfLevelSet 函数，如果 Level 字段存在
	if deviceInfo.PerfLevel != "" {
		err := dcgm.DevPerfLevelSet(dvInd, levelConverted)
		if err != nil {
			executionErrors = append(executionErrors, "设置性能级别失败: "+err.Error())
		}
	}

	// 调用 DevGpuClkFreqSet 函数，如果 SclkClock 字段存在
	if deviceInfo.SclkClock != "" {
		err := dcgm.DevGpuClkFreqSet(dvInd, dcgm.RSMI_CLK_TYPE_SYS, sclkClock)
		if err != nil {
			executionErrors = append(executionErrors, "设置时钟频率失败: "+err.Error())
		}
	}

	// 调用 DevGpuClkFreqSet 函数，如果 SocclkClock 字段存在
	if deviceInfo.SocclkClock != "" {
		err := dcgm.DevGpuClkFreqSet(dvInd, dcgm.RSMI_CLK_TYPE_SOC, socclkClock)
		if err != nil {
			executionErrors = append(executionErrors, "设置时钟频率失败: "+err.Error())
		}
	}

	// 调用 ResetFans 函数，如果 ResetFan 为 true
	if deviceInfo.ResetFan {
		err := dcgm.ResetFans([]int{dvInd})
		if err != nil {
			executionErrors = append(executionErrors, "重置风扇失败: "+err.Error())
		}
	}

	// 如果有任何方法执行失败，返回失败的错误信息
	if len(executionErrors) > 0 {
		executionErrorInfo := map[string]interface{}{
			"executionErrors": executionErrors,
		}
		c.JSON(http.StatusBadRequest, ErrorResponse(executionErrorInfo))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(nil))
}

// EccBlocksInfo 处理设备控制
// @Summary 获取 ECC block 信息
// @Description 根据设备索引获取 ECC block 信息
// @Accept json
// @Produce json
// @Param dvInd query int true "设备索引"
// @Success 200 {array} BlocksInfo "ECC block 信息"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "内部服务器错误"
// @Router /EccBlocksInfo [get]
func EccBlocksInfo(c *gin.Context) {
	// 获取请求中的 dvInd 参数
	var dvInd int
	if err := c.BindQuery(&dvInd); err != nil {
		c.JSON(http.StatusBadRequest, "虚拟设备销毁失败")
		return
	}

	// 调用 EccBlocksInfo 函数
	blocksInfos, err := dcgm.EccBlocksInfo(dvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	// 返回成功响应
	response := map[string]interface{}{
		"blocksInfos": blocksInfos,
	}
	c.JSON(http.StatusOK, SuccessResponse(response))
}

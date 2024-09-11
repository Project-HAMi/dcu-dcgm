package router

import (
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	// 注册路由
	router.GET("/devicename/:dvInd", GetDevName)
	router.GET("/NumMonitorDevices", GetNumMonitorDevices)
	router.GET("/DevSku/:dvInd", GetDevSku)
	router.GET("/DevBrand/:dvInd", DevBrand)
	router.GET("/DevVendorName/:dvInd", DevVendorName)
	router.GET("/DevVramVendor/:dvInd", DevVramVendor)
	router.GET("/DevPciBandwidth/:dvInd", DevPciBandwidth)
	router.POST("/DevPciBandwidthSet", DevPciBandwidthSet)

	router.GET("/MemoryPercent/:dvInd", MemoryPercent)
	// 路由注册
	router.GET("/PerfLevel/:dvInd", PerfLevel)
	router.POST("/DevPerfLevelSet/:dvInd", DevPerfLevelSet)
	router.GET("/DevGpuMetricsInfo/:dvInd", DevGpuMetricsInfo)
	router.GET("/CollectDeviceMetrics", CollectDeviceMetrics)
	router.GET("/DeviceInfo/:dvInd", GetDeviceByDvInd)
	// 路由
	router.GET("/PicbusInfo/:dvInd", PicBusInfo)
	// 路由(K100 AI不支持)
	router.GET("/FanSpeedInfo/:dvInd", FanSpeedInfo)
	// 路由
	router.GET("/DCUUse/:dvInd", GPUUse)
	// 路由
	router.GET("/DevID", GetDevID)
	// 路由
	router.GET("/MaxPower", GetMaxPower)
	// 路由
	router.GET("/MemInfo", GetMemInfo)
	router.GET("/AllDeviceInfos", AllDeviceInfos)
	// 路由
	router.GET("/DeviceInfos", GetDeviceInfos)
	// 路由
	router.GET("/ProcessName", GetProcessName)
	// 路由注册
	router.GET("/Power/:dvInd", Power)

	//设置 GPU 时钟频率
	router.POST("/DevGpuClkFreqSet", DevGpuClkFreqSet)
	// 路由注册
	router.GET("/EccStatus/:dvInd", EccStatus)
	// 路由注册
	router.GET("/Temperature/:dvInd", Temperature)
	// 路由注册
	router.GET("/VbiosVersion/:dvInd", VbiosVersion)
	// 路由注册
	router.GET("/Version", Version)
	// 重置设备时钟(K100 AI不支持)
	router.POST("/ResetClocks", ResetClocks)
	router.POST("/ResetFans", ResetFans)

	router.POST("/ResetProfile", ResetProfile)
	//(K100 AI不支持)
	router.POST("/ResetXGMIErr", ResetXGMIErr)
	//(K100 AI不支持)
	router.GET("/XGMIErrorStatus", XGMIErrorStatus)
	router.GET("/XGMIHiveId", XGMIHiveIdGet)
	router.POST("/ResetPerfDeterminism", ResetPerfDeterminism)
	// 路由(K100 AI不支持)
	router.POST("/SetClockRange", SetClockRange)
	// 路由（K100_AI卡不支持）
	router.POST("/SetPowerPlayTableLevel", SetPowerPlayTableLevel)
	// 路由（sudo权限)
	router.POST("/SetClockOverDrive", SetClockOverDrive)
	// 路由（K100_AI卡不支持）
	router.POST("/SetPerfDeterminism", SetPerfDeterminism)
	// 设置风扇速度(K100 AI不支持)
	router.POST("/SetFanSpeed", SetFanSpeed)
	// 获取设备风扇转速(K100 AI不支持)
	router.GET("/DevFanRpms/:dvInd", DevFanRpms)
	// 设置设备性能等级
	router.POST("/SetPerformanceLevel", SetPerformanceLevel)
	// 设置功率配置文件（K100_AI卡不支持该操作,CUSTOM不支持，剩余几个类型超出安全范围）
	router.POST("/SetProfile", SetProfile)
	// 设置设备功率配置文件（K100_AI卡不支持该操作）
	router.POST("/DevPowerProfileSet/:dvInd", DevPowerProfileSet)
	// 获取设备总线信息
	router.GET("/GetBus/:dvInd", GetBus)
	// 显示设备硬件信息
	router.POST("/ShowAllConciseHw", ShowAllConciseHw)
	// 显示设备时钟信息
	router.POST("/ShowClocks", ShowClocks)
	router.POST("/fans/current", ShowCurrentFans)
	router.POST("/temps/current", ShowCurrentTemps)
	router.GET("/firmware/info", ShowFwInfo)
	router.GET("/process/list", PidList)
	router.POST("/utilization/coarse", GetCoarseGrainUtil)
	router.POST("/gpu/use", ShowGpuUse)
	router.POST("/energy", ShowEnergy)
	router.POST("/memory/info", ShowMemInfo)
	router.POST("/memory/use", ShowMemUse)
	router.POST("/memory/vendor", ShowMemVendor)
	router.POST("/pcie/bandwidth", ShowPcieBw)
	router.POST("/pcie/replaycount", ShowPcieReplayCount)
	router.GET("/process/name/:pid", GetProcessName)
	router.POST("/device/power", GetDevicePower)
	//（K100_AI卡不支持该操作）
	router.POST("/device/powerplay", GetDevicePowerPlayTable)
	router.POST("/device/product", GetDeviceProductName)
	router.POST("/device/profile", GetDeviceProfile)

	router.POST("/device/retiredpages", GetDeviceRetiredPages)
	router.POST("/device/serialnumber", GetDeviceSerialNumber)
	router.POST("/showUId", ShowUId)
	router.POST("/showVbiosVersion", ShowVbiosVersion)
	router.POST("/showVoltage", ShowVoltage)
	//（K100_AI卡不支持该操作）
	router.POST("/showVoltageCurve", ShowVoltageCurve)
	//（K100_AI卡不支持该操作）
	router.POST("/showXgmiErr", ShowXgmiErr)

	router.POST("/showWeightTopology", ShowWeightTopology)
	router.POST("/showHopsTopology", ShowHopsTopology)
	router.POST("/showTypeTopology", ShowTypeTopology)
	router.POST("/showNumaTopology", ShowNumaTopology)
	router.POST("/showHwTopology", ShowHwTopology)
	router.GET("/deviceCount", DeviceCount)
	router.GET("/VDeviceSingleInfo", VDeviceSingleInfo)
	router.GET("/vDeviceCount", VDeviceCount)
	router.GET("/deviceRemainingInfo/:dvInd", DeviceRemainingInfo)
	router.POST("/CreateVDevices", CreateVDevices)
	router.DELETE("/DestroyVDevice", DestroyVDevice)
	router.DELETE("/DestroySingleVDevice", DestroySingleVDevice)
	router.PUT("/UpdateSingleVDevice", UpdateSingleVDevice)
	// 启动指定的虚拟设备
	router.GET("/StartVDevice/:vDvInd", StartVDevice)
	// 停止指定的虚拟设备
	router.GET("/StopVDevice/:vDvInd", StopVDevice)
	// 设置虚拟机加密状态
	router.POST("/SetEncryptionVMStatus", SetEncryptionVMStatus)
	// 获取加密虚拟机状态
	router.GET("/EncryptionVMStatus", EncryptionVMStatus)
	// 打印设备的事件列表
	router.GET("/PrintEventList/:device", PrintEventList)
	router.GET("/device/info/:dvInd", GetDeviceInfo)
	// 路由
	router.POST("/device/control", DeviceControl)
	return router
}

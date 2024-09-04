package main

import (
	"flag"

	"github.com/golang/glog"

	"g.sugon.com/das/dcgm-dcu/pkg/dcgm"
)

// 添加注释以描述 server 信息
//	@title			Swagger Example API
//	@version		1.0
//	@description	This is a sample server celler server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/router/v1

// @securityDefinitions.basic	BasicAuth
func main() {
	flag.Parse()
	defer glog.Flush()
	glog.Info("go-dcgm start ...")
	//初始化dcgm服务
	dcgm.Init()
	defer dcgm.ShutDown()
	//DCU物理设备总数量
	dcgm.NumMonitorDevices()
	//DCU名称
	dcgm.DevName(0)
	//获取指定DCU设备可用的PCIE宽带列表
	dcgm.DevPciBandwidth(0)
	//设置DCU的PowerPlay性能级别
	dcgm.DevPerfLevelSet(0, dcgm.RSMI_DEV_PERF_LEVEL_LOW)
	//DCU设备内存使用百分比
	dcgm.MemoryPercent(0)
	//获取DCU设备名称
	dcgm.DevName(1)
	//获取设备的sku
	dcgm.DevSku(1)
	//获取设备品牌名称
	dcgm.DevBrand(1)
	//获取设备的供应商的名称
	dcgm.DevVendorName(1)
	//获取设备显存供应商名称
	dcgm.DevVramVendor(1)
	//获取指定DCU的度量信息
	dcgm.DevGpuMetricsInfo(0)
	//获取物理设备的监控指标
	dcgm.CollectDeviceMetrics()
	//DCU设备指定类型的内存使用情况 [vram|vis_vram|gtt
	dcgm.MemInfo(0, "vram")
	//dcgm.MemInfo(0, "vis_vram")
	//dcgm.MemInfo(0, "gtt")
	//获取设备当前的性能水平
	dcgm.PerfLevel(0)
	//获取设备的平均功率
	dcgm.Power(1)
	//设备的VBIOS版本信息
	dcgm.VbiosVersion(1)
	//获取系统的驱动程序版本
	dcgm.Version(dcgm.RSMISwCompFirst)
	dcgm.Version(dcgm.RSMISwCompDriver)
	dcgm.Version(dcgm.RSMISwCompLast)
	// 调用方法打印事件列表
	//获取设备的XGMI hive id
	dcgm.XGMIHiveIdGet(1)

	//批量展示显示设备硬件信息
	dcgm.ShowAllConciseHw([]int{0, 1, 2})
	//批量展示显示时钟信息
	dcgm.ShowClocks([]int{0, 1, 2})
	//展示风扇转速和风扇级别
	dcgm.ShowCurrentFans([]int{0, 1, 2}, false)
	//显示设备的所有可用温度传感器的温度
	dcgm.ShowCurrentTemps([]int{0, 1, 2})
	//显示设备中指定固件类型的固件版本信息
	dcgm.ShowFwInfo([]int{0, 1, 2}, []string{"all"})

	//获取DCU设备的的粗粒度利用率
	dcgm.GetCoarseGrainUtil(0, "all")
	//批量获取DCU的使用率
	dcgm.ShowGpuUse([]int{0, 1, 2})
	//批量获取设备消耗的能量
	dcgm.ShowEnergy([]int{0, 1, 2})
	//DCU设备的ID（十六进制表示）
	dcgm.DevID(0)
	//设备的最大功率值
	dcgm.MaxPower(0)
	//设备的不同类型的内存使用情况 memType:[vram|vis_vram|gtt]
	dcgm.MemInfo(0, "vram")
	//批量获取设备内存的信息
	dcgm.ShowMemInfo([]int{0, 1, 2}, []string{"VRAM", "VIS_VRAM"})
	//批量获取设备内存使用情况
	dcgm.ShowMemUse([]int{0, 1, 2})
	//批量获取备供应商信息
	dcgm.ShowMemVendor([]int{0, 1, 2})
	//批量获取设备的PCIe带宽使用情况
	dcgm.ShowPcieBw([]int{0, 1, 2})
	//批量获取设备PCIe重放计数
	dcgm.ShowPcieReplayCount([]int{0, 1, 2})

	//批量获取设备的平均功率
	dcgm.ShowPower([]int{0, 1, 2})

	//当前设备内存时钟频率和电压（K100_AI卡不支持该操作）
	dcgm.ShowPowerPlayTable([]int{0, 1, 2})
	//可用电源配置文件
	dcgm.ShowProfile([]int{0, 1, 2})

	//电流或电压范围（K100_AI卡不支持该操作）
	devices := []int{0, 1, 2}
	dcgm.ShowRange(devices, "sclk")
	dcgm.ShowRange(devices, "mclk")
	dcgm.ShowRange(devices, "voltage")

	//显示设备中指定类型的退役页
	dcgm.ShowRetiredPages([]int{0, 1, 2}, "all")
	//设备序列号
	dcgm.ShowSerialNumber([]int{0, 1, 2})
	//设备的唯一设备ID
	dcgm.ShowUId([]int{0, 1, 2})
	//设备的VBIOS版本信息（格式化打印并返回设备的VBIOS版本信息）
	dcgm.ShowVbiosVersion([]int{0, 1, 2})

	// 示例设备列表和事件类型
	deviceList := []int{0, 1, 2}
	eventTypes := []string{"VM_FAULT", "THERMAL_THROTTLE"}
	dcgm.ShowEvents(deviceList, eventTypes)
	//指定设备的当前电压信息
	dcgm.ShowVoltage([]int{0, 1, 2})

	//获取指定设备的电压曲线点（K100_AI卡不支持该操作）
	dcgm.ShowVoltageCurve([]int{0, 1, 2})
	//指定设备的XGMI错误状态（K100_AI卡不支持该操作）
	dcgm.ShowXgmiErr([]int{0, 1, 2}, true)

}

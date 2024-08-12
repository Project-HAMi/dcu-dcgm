package main

import (
	"flag"

	"github.com/golang/glog"

	"g.sugon.com/das/dcgm-dcu/pkg/dcgm"
)

func main() {
	flag.Parse()
	defer glog.Flush()
	glog.Info("go-dcgm start ...")
	//初始化dcgm服务
	dcgm.Init()
	defer dcgm.ShutDown()

	// 打印事件列表
	eventList := []string{"VMFAULT", "FIRST", "THERMAL_THROTTLE", "GPU_PRE_RESET", "GPU_POST_RESET", "LAST"} // 示例事件列表
	dcgm.PrintEventList(1, 100, eventList)

	//批量复位风扇驱动控制
	dcgm.ResetFans([]int{0, 1})
	//批量重置设备的配置文件
	dcgm.ResetProfile([]int{0, 1})

	// 重置设备的XGMI错误状态（K100_AI卡不支持该操作）
	dcgm.ResetXGMIErr([]int{0, 1})
	dcgm.XGMIErrorStatus(0)
	dcgm.XGMIErrorStatus(1)

	//为设备选定的时钟类型设定相应的频率范围（K100_AI卡不支持该操作）
	dcgm.SetClockRange([]int{0}, "sclk", "1", "100", true)
	//设置 PowerPlay 级别
	dcgm.SetPowerPlayTableLevel([]int{0}, "sclk", "1", "10", "100", true)

	//设置时钟频率级别以启用性能确定性（K100_AI卡不支持该操作）
	dcgm.SetPerfDeterminism([]int{0}, "900")
	//设置风扇转速
	dcgm.SetFanSpeed([]int{0}, "200")
	//获取设备风扇的实际转速
	dcgm.DevFanRpms(0)
	//批量设置设备性能 level:auto、low、high、normal
	dcgm.SetPerformanceLevel([]int{0}, "auto")
	//设置功率配置（K100_AI卡不支持该操作）
	dcgm.SetProfile([]int{0}, "BOOTUP DEFAULT")
	//设置设备功率配置文件（K100_AI卡不支持该操作）
	dcgm.DevPowerProfileSet(0, 0, dcgm.RSMI_PWR_PROF_PRST_BOOTUP_DEFAULT)
}

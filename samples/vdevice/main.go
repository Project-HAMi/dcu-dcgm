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
	dcgm.Init()
	defer dcgm.ShutDown()
	dcgm.VDeviceCount()
	dcgm.AllDeviceInfos()

	//ticker := time.NewTicker(10 * time.Second) // 创建一个每隔10秒触发的定时器
	//defer ticker.Stop()                        // 确保在函数结束时停止定时器
	//
	//quit := make(chan struct{}) // 创建一个用于停止循环的通道
	//
	//go func() {
	//	for {
	//		select {
	//		case <-ticker.C:
	//			dcgm.AllDeviceInfos() // 每隔10秒调用一次函数
	//		case <-quit:
	//			return
	//		}
	//	}
	//}()
	//
	//// 假设我们希望程序运行一段时间后停止，可以使用 time.Sleep 或者其他方式
	//time.Sleep(180 * time.Second) // 例如，运行60秒后停止
	//
	//close(quit) // 发送信号以停止循环
	//fmt.Println("程序结束")

	//dcgm.DestroySingleVDevice(1)
	//dcgm.DeviceSingleInfo(4)

	//dcgm.DeviceRemainingInfo(0)
	//dcgm.DeviceRemainingInfo(1)

	//dcgm.CreateVDevices(0, 2, []int{4, 4}, []int{1024, 2048})
	//dcgm.DestroyVDevice(1)

	//dcgm.UpdateSingleVDevice(5, 20, 8589934592)

	//dcgm.StartVDevice(2)
	//dcgm.StopVDevice(2)
	//
	//dcgm.EncryptionVMStatus()
	//dcgm.SetEncryptionVMStatus(true)

	// 定义目录路径
	//dirPath := "/etc/vdev"
	//
	//// 读取目录中的文件列表
	//files, err := os.ReadDir(dirPath)
	//if err != nil {
	//	log.Fatalf("无法读取目录: %v", err)
	//}
	//
	//// 打印文件数量
	//fmt.Printf("文件数量: %d\n", len(files))
	//
	//// 逐个读取并解析每个文件的内容
	//for _, file := range files {
	//	// 确保是文件而不是子目录
	//	if !file.IsDir() {
	//		filePath := filepath.Join(dirPath, file.Name())
	//		config, err := dcgm.ParseConfig(filePath)
	//		if err != nil {
	//			log.Printf("无法解析文件 %s: %v", filePath, err)
	//			continue
	//		}
	//		fmt.Printf("文件: %s\n配置: %+v\n", filePath, config)
	//	}
	//}
}

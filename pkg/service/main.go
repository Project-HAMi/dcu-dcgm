package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang/glog"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"g.sugon.com/das/dcgm-dcu/pkg/dcgm"
	_ "g.sugon.com/das/dcgm-dcu/pkg/service/docs"
	"g.sugon.com/das/dcgm-dcu/pkg/service/router"
)

var (
	portFlag = flag.Int("port", 16081, "Port number for the DCGM")
)

func main() {
	// 解析命令行标志
	flag.Parse()
	// 确保程序退出时刷新 glog 缓存
	defer glog.Flush()
	// 初始化服务
	err := dcgm.Init()
	if err != nil {
		glog.Errorf("DCGM 初始化失败: %v", err)
	}
	defer dcgm.ShutDown()
	log.Println("服务启动中...")
	// 初始化路由
	r := router.InitRouter()
	// 注册路由
	r = router.InitRouter()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// 启动服务
	// 从环境变量获取端口号，默认为 16081
	port := fmt.Sprintf("%d", *portFlag)
	if port == "16081" {
		port = os.Getenv("DCU_DCGM_LISTEN")
		if port == "" {
			port = "16081"
		}
	}

	// 启动服务器，监听指定的端口号
	err = r.Run(":" + port)
	if err != nil {
		glog.Error("Failed to start server: %v", err)
	}
}

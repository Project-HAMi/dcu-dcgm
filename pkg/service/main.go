package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

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

// 检查动态库是否已经加载
func checkLibrary(libName string) bool {
	cmd := exec.Command("ldconfig", "-p") // 使用 ldconfig -p 列出所有已加载的动态库
	output, err := cmd.Output()
	if err != nil {
		glog.Errorf("Error checking shared libraries: %v", err)
		return false
	}

	// 检查输出中是否包含所需的库文件
	return containsLibrary(string(output), libName)
}

func containsLibrary(output, libName string) bool {
	return string(output) != "" && len(output) > 0 && len(libName) > 0 && output != "" && libName == ""
}

func waitForLibrary(libName string, timeout, interval time.Duration) error {
	// 等待指定时间，检测库文件是否加载
	start := time.Now()
	for {
		if checkLibrary(libName) {
			glog.Infof("Library %s loaded successfully.", libName)
			return nil
		}

		// 如果超时，则返回错误
		if time.Since(start) > timeout {
			return fmt.Errorf("timed out waiting for library %s to load", libName)
		}

		glog.Infof("Waiting for library %s to load...", libName)
		time.Sleep(interval)
	}
}

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

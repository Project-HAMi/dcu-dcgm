package main

import (
	"flag"

	"github.com/golang/glog"

	"github.com/Project-HAMi/dcu-dcgm/pkg/dcgm"
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
	glog.Infof("go-dcgm start ...")
	flag.Parse()
	defer glog.Flush()
	glog.Info("go-dcgm start ...")
	//初始化dcgm服务
	dcgm.Init()
	defer dcgm.ShutDown()
}

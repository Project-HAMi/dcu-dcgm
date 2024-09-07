# DCU DCGM

## 组件信息

DCU DCGM 为 DCU 管理提供 Golang 绑定接口，是管理和监控DCU的工具。包括健康状态监控、功率、时钟频率调控，以及资源使用情况统计等。

## 组件使用前置条件

**需要在组件部署的主机安装DCU驱动。如果没安装驱动则需要添加动态链接库。**

##### 添加动态链接库具体步骤如下：

1. 将库文件libhydmi.so.1.4和librocm_smi64.so.2.8放在某一个路径下（例如：/usr/lib/路径下）。

2. libhydmi.so.1.4和librocm_smi64.so.2.8两个库文件创建软连接，详细配置如下：

   ```
   # ll /usr/lib | grep .so*
   lrwxrwxrwx  1 root root     22 Aug  6 08:25 libhydmi.so -> /usr/lib/libhydmi.so.1
   lrwxrwxrwx  1 root root     24 Aug  6 08:25 libhydmi.so.1 -> /usr/lib/libhydmi.so.1.4
   -rw-rw-r--  1 root root 834456 Aug  6 08:24 libhydmi.so.1.4
   lrwxrwxrwx  1 root root     27 Aug  6 08:25 librocm_smi64.so -> /usr/lib/librocm_smi64.so.2
   lrwxrwxrwx  1 root root     29 Aug  6 08:25 librocm_smi64.so.2 -> /usr/lib/librocm_smi64.so.2.8
   -rw-rw-r--  1 root root 789440 Aug  6 08:24 librocm_smi64.so.2.8
   ...
   ```

3. 库路径添加到环境变量中，执行命令

   ```
   export SH_DIR=/usr
   ```

4. 指定目录添加到 `LD_LIBRARY_PATH` 环境变量中

```
export LD_LIBRARY_PATH=LD_LIBRARY_PATH:$SH_DIR/lib
```

## 使用流程

目前DCGM-DCU支持两种方式使用，一种是作为依赖库的形式集成到个人项目中调用函数；另一种是把DCGM-DCU作为服务启动，通过HTTP请求的方式调用接口。具体使用步骤如下：

##### 作为依赖库的形式使用流程

*目前代码仅在内部gitlab中存放，其他项目调用流程如下：*

1. git clone本项目代码到本地，与调用者项目存放于同级目录中；

2. 调用者项目修改go.mod文件：

```
replace g.sugon.com/das/dcgm-dcu => /your/path/dcgm-dcu
```

3. 在本地项目中执行：

```
go mod tidy
```

4. 在golang文件中import相关依赖包之后使用即可，其中api.go为封装DCGM的API调用，提供与DCGM库交互的各种API接口，处理具体的功能调用。/pkg/samoles下是简单的test：

```go
import (
...
"g.sugon.com/das/dcgm-dcu/pkg/dcgm"
...)


func main(){
...
	dcgm.Init()
    defer dcgm.ShutDown()
...

}
```

##### 作为服务启动形式使用流程

1. 将提供的二进制包放在主机路径下。（例如：/home/dcgm）

2. 进入二进制包所在的路径下执行启动命令(默认端口：16081)

   ```
   //后台运行服务，并将日志打印到同目录下的dcgm.log中
   nohup ./dcgm -logtostderr > dcgm.log 2>&1 &
   
   //后台运行服务，并将日志打印到同目录下的dcgm.log中，包括info级别的日志
   nohup ./dcgm -logtostderr -v=2 > dcgm.log 2>&1 &
   
   //指定端口号启动服务
   export DCU_DCGM_LISTEN=12345
   nohup ./dcgm -logtostderr -v=2 > dcgm.log 2>&1 &
   或者
   nohup env DCU_DCGM_LISTEN=12345 ./dcgm -logtostderr -v=2 > dcgm.log 2>&1 &
   ```

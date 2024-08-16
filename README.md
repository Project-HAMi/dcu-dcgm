# DCU DCGM

## 组件信息

DCU DCGM 为 DCU 管理提供 Golang 绑定接口，是管理和监控DCU的工具。包括健康状态监控、功率、时钟频率调控，以及资源使用情况统计等。

## 组件使用前置条件

组件部署主机上安装DCU驱动，或在系统默认动态链接库加载路径下存在DCU动态链接库libhydmi.so和librocm_smi64.so，下面以/usr/lib动态链接库加载路径为例说明动态链接库配置详情。
```bash
# ll /usr/lib | grep .so*
lrwxrwxrwx  1 root root     22 Aug  6 08:25 libhydmi.so -> /usr/lib/libhydmi.so.1
lrwxrwxrwx  1 root root     24 Aug  6 08:25 libhydmi.so.1 -> /usr/lib/libhydmi.so.1.4
-rw-rw-r--  1 root root 834456 Aug  6 08:24 libhydmi.so.1.4
lrwxrwxrwx  1 root root     27 Aug  6 08:25 librocm_smi64.so -> /usr/lib/librocm_smi64.so.2
lrwxrwxrwx  1 root root     29 Aug  6 08:25 librocm_smi64.so.2 -> /usr/lib/librocm_smi64.so.2.8
-rw-rw-r--  1 root root 789440 Aug  6 08:24 librocm_smi64.so.2.8
...
```

## 使用流程

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

“# dcu-dcgm”

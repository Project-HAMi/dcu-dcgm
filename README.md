# DCU DCGM

## 组件信息

为 DCU 管理提供的 Golang 绑定接口。是管理和监控DCU的工具。包括健康状态监控、功率、时钟频率调控等，以及资源使用情况的统计。

## 前置条件条件

所在宿主机上安装DCU驱动

## 使用流程

*目前代码仅在内部gitlab中存放，开发调用流程如下：*

1.git clone本项目代码到本地，与调用者项目存放于相同目录中；

2.调用者项目修改go.mod文件：

```
replace g.sugon.com/das/dcgm-dcu => /your/path/dcgm-dcu
```

3.在本地项目中执行：

```
go mod tidy
```

4.在golang文件中import相关依赖包之后使用即可，其中api.go为封装DCGM的API调用，提供与DCGM库交互的各种API接口，处理具体的功能调用。/pkg/samoles下是简单的test：

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


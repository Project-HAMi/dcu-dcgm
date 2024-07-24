# DCU DCGM



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

4.在golang文件中import相关依赖包之后使用即可：

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


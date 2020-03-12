### 简介

一行代码接入需要监控的服务（或者说进程、程序）。持续性采集服务的各项内存指标（包含Go heap管理相关的各项指标），并以时间维度绘制成折现图。可在浏览器中查看。

### 效果图

TODO

Go runtime 内存相关的指标含义参见：https://pengrl.com/p/20031/

进程VMS和RSS的含义参见：https://pengrl.com/p/21292/

### 用法

#### 1. 被监控的进程中添加以下代码：

```golang
import "github.com/q191201771/pprofplus"

pprofplus.Start()
```

使用示例，以及更多个性化配置的方法见： [example/example.go](https://github.com/q191201771/pprofplus/blob/master/example/example.go)

支持的配置项见：[pprofplus.go](https://github.com/q191201771/pprofplus/blob/master/pprofplus.go#L5)

#### 2. 启动dashboard展示程序（与被监控的进程在一台机器）：

```shell
./dashboard
```

更多的定制化配置见： `./dashboard -h`

我在另一个repo中提交了编译好的二进制文件可供直接使用：https://github.com/q191201771/pprofplus.bin

#### 3. 浏览器访问网页查看图表：

`http://<hostname>:10001/pprofplus`

### 为什么要单独搞一个dashboard展示程序，而不直接放在被监控的进程中展示？

因为我不想因为展示部分使用到了内存，而影响到原有被监控进程的内存使用情况。

### 技术细节

说起来十分简单，接入服务的部分，会每间隔5秒（可配置），通过Go runtime的接口获取宿主进程的Go heap的情况，并且获取进程的虚拟内存和常驻内存。

获取到的数据会写入到本地文件中。

dashboard读取本地文件中的数据，将数据以http图表的方式返回给用户查看。

### 第三方库依赖：

#### pprofplus部分

- 获取进程虚拟内存和常驻内存： github.com/shirou/gopsutil/process

#### dashboard部分

- http图表： github.com/go-echarts/go-echarts/charts

### 编译dashboard注意事项

如果你不需要对代码进行二次开发，建议直接从 https://github.com/q191201771/pprofplus.bin 下载编译好的dashboard直接使用。

如果自行编译dashboard，并且想要发布到非开发机器上使用。需在编译前执行以下操作（一次即可）：

```shell
// Linux/MacOS
$ cd $GOPATH/src/github.com/chenjiandongx/go-echarts
$ ./doPackr.sh

// Windows
$ cd %GOPATH%\src\github.com\chenjiandongx\go-echarts
$ doPackr.bat
```

这是因为第三方库go-echarts说白了是对前端echarts库的封装，它在运行时需要用到echarts的一些静态文件。  
在非开发机器上由于没有这些静态文件，会导致运行出错。  
go-echarts的解决方案是通过上面的命令，将这些静态文件打包成了go文件，这样就直接编译进宿主程序中了。  
这其实也是我把dashboard独立出来的原因之一，不想让宿主程序依赖这个库。。

### TODO

- go mod 支持
- pprofplus和pprofplus.bin的版本管理

#### 采集端

- 增加打开、关闭接口，可以在服务运行期间修改是否采集
- 目前单个进程的dump文件没有切割

#### 展示端dashboard

- url中增加参数
    - 开始时间、结束时间，或者距离最近的时间段
- 目前只读取固定目录下serviceName匹配的最新的那个文件，可支持指定特定文件，特定pid
- 网页默认标题

#### 长远计划

结合实际使用场景，增加更多的指标，并且可能还会添加内存以外的指标

#### 其他

欢迎提issues讨论。

本项目只在macos和linux下做过测试，windows的表现未知。
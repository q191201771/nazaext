目录

- 简介
- 效果图
- 接入方法
- 实现细节
- 第三方库依赖
- TODO
- 其他

### 简介

pprofplus用于采集Go进程的各项内存指标（包含Go runtime内存管理相关的，以及RSS等），并按时间维度绘制成折线图，可以通过网页实时查看。配合原有Go pprof工具，可以快速监控和分析Go进程的内存使用情况。

### 效果图

![效果图](https://raw.githubusercontent.com/q191201771/pprofplus.bin/master/snapshot.jpg)

- 进程VMS和RSS的含义参见：[《[译] linux内存管理之RSS和VSZ的区别》](https://pengrl.com/p/21292/)
- Go runtime 内存相关的指标含义参见：[《Go pprof内存指标含义备忘录》](https://pengrl.com/p/20031/)

分析案例：

- [《Go进程的HeapReleased上升，但是RSS不下降造成内存泄漏？》](https://pengrl.com/p/20033/)

### 接入方法

#### 1. 被监控的进程中添加以下代码：

```golang
import "github.com/q191201771/pprofplus/pkg/pprofplus"

pprofplus.Start()
```

使用示例，以及更多个性化配置的方法见： [example/example.go](https://github.com/q191201771/pprofplus/blob/master/example/example.go)

支持的配置项见：[pprofplus.go](https://github.com/q191201771/pprofplus/blob/master/pkg/pprofplus/pprofplus.go#L13)

#### 2. 启动dashboard程序（web展示用，与被监控的进程在一台机器）：

```shell
./dashboard
```

更多的定制化配置见： `./dashboard -h`

我在另一个repo中提交了编译好的二进制文件可供直接使用：https://github.com/q191201771/pprofplus.bin

#### 3. 浏览器访问网页查看图表：

`http://<yourhostname>:10001/pprofplus`

#### 为什么要单独搞一个dashboard展示程序，而不直接放在被监控的进程中展示？

因为独立开可以保证，展示部分使用的内存，不会增加被监控进程的内存使用。

### 实现细节

说起来十分简单，接入宿主程序的部分，会周期性获取自身进程的各项指标。

获取到的数据会写入到本地文件中。

dashboard读取本地文件中的数据，将数据以图表的方式通过http页面返回给用户查看。

### 第三方库依赖：

#### pprofplus部分

- 获取进程虚拟内存和常驻内存： github.com/shirou/gopsutil/process

#### dashboard部分

- http图表： github.com/go-echarts/go-echarts/charts

#### 编译dashboard注意事项

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

#### 整体

- go mod 支持
- pprofplus和pprofplus.bin的版本管理，包括数据的版本

#### 采集端

- 增加打开、关闭接口，可以在服务运行期间修改是否采集
- 目前单个进程的dump文件没有切割

#### 展示端dashboard

- url中增加参数
    - 开始时间、结束时间，或者距离最近的时间段
- 目前只读取固定目录下serviceName匹配的最新的那个文件，可支持指定特定文件，特定pid
- dashboard中的一部分代码从app移入pkg中，业务方也可选择直接在宿主进程中直接接入dashboard

#### 长远计划

结合实际使用场景，增加更多的指标，并且可能还会添加内存以外的指标。绘制多张图表。但是会尽力小而美。

### 其他

欢迎提issues讨论。

本项目只在macos和linux下做过测试，windows的表现未知也不做保证。。

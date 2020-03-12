package pprofplus

import "sync"

type Option struct {
	CaptureIntervalSec int    // 采集间隔，单位秒
	DumpDir            string // 采集数据存放的目录
	ServiceName        string // 采集服务名，展示数据时可根据这个字段，展示相应服务的信息
}

// 配置默认值
var option = Option{
	CaptureIntervalSec: 5,
	DumpDir:            "/tmp/pprofplus/",
	ServiceName:        "pprofplus",
}

type Info struct {
	Timestamp int64

	Sys          uint64 `json:"sys"`
	HeapSys      uint64 `json:"heapsys"`
	HeapAlloc    uint64 `json:"heapalloc"`
	HeapInuse    uint64 `json:"heapinuse"`
	HeapReleased uint64 `json:"heapreleased"`
	HeapIdle     uint64 `json:"heapidle"`

	VMS uint64 `json:"vms"`
	RSS uint64 `json:"rss"`
}

type ModOption func(option *Option)

var once sync.Once

func Start(modOptions ...ModOption) {
	once.Do(func() {
		start(modOptions...)
	})
}

func start(modOptions ...ModOption) {
	for _, mo := range modOptions {
		mo(&option)
	}

	c := NewCapture(option.CaptureIntervalSec)
	d, _ := NewDump(option.DumpDir, option.ServiceName)

	go func() {
		infoC := c.doAsync()
		for {
			info := <-infoC
			d.do(info)
		}
	}()
}

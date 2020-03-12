// Copyright 2020, Chef.  All rights reserved.
// https://github.com/q191201771/pprofplus
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

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

func Start(modOptions ...ModOption) error {
	var err error
	once.Do(func() {
		err = start(modOptions...)
	})
	return err
}

func start(modOptions ...ModOption) error {
	for _, mo := range modOptions {
		mo(&option)
	}

	c := NewCapture(option.CaptureIntervalSec)
	d, err := NewDump(option.DumpDir, option.ServiceName)
	if err != nil {
		return err
	}

	go func() {
		infoC := c.doAsync()
		for {
			info := <-infoC
			d.do(info)
		}
	}()

	return nil
}

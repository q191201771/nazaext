// Copyright 2020, Chef.  All rights reserved.
// https://github.com/q191201771/pprofplus
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package main

import (
	"fmt"
	"github.com/q191201771/pprofplus"
	//"net/http"
	//_ "net/http/pprof"
	"os"
)

func main() {
	if err := pprofplus.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 通过这种方式设置个性化配置
	//pprofplus.Start(func(option *pprofplus.Option) {
	//	option.CaptureIntervalSec = 10
	//})

	fmt.Println("an example show how service interact with pprofplus.")

	// 被监控的程序依然可以使用http pprof，和pprofplus配合一块分析内存情况
	// 注意，这里的端口不要和dashboard的重了
	//http.ListenAndServe(":10002", nil)

	ch := make(chan struct{})
	<-ch
}

package main

import (
	"fmt"
	"github.com/q191201771/pprofplus"
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
	ch := make(chan struct{})
	<-ch
}

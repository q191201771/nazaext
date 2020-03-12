// Copyright 2020, Chef.  All rights reserved.
// https://github.com/q191201771/pprofplus
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/go-echarts/go-echarts/charts"
	"github.com/q191201771/pprofplus/pkg/pprofplus"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	dir         *string
	serviceName *string
	addr        *string
	uri         *string

	memUnit = MemUnit(MemUnitMByte)
)

type MemUnit int

const (
	MemUnitByte MemUnit = iota + 1
	MemUnitKByte
	MemUnitMByte
	MemUnitGByte
)

var memUnitBrief []string

func main() {
	dir = flag.String("dir", "/tmp/pprofplus/", "dir of pprofplus dump file")
	serviceName = flag.String("service", "pprofplus", "service name")
	addr = flag.String("addr", ":10001", "dashboard addr")
	uri = flag.String("uri", "/pprofplus", "dashboard uri")
	flag.Parse()
	fmt.Printf("dir=%s, service=%s, addr=%s, uri=%s\n", *dir, *serviceName, *addr, *uri)

	http.HandleFunc(*uri, requestHandler)
	err := http.ListenAndServe(*addr, nil)
	fmt.Println(err)
}

func getInfos() ([]pprofplus.Info, error) {
	var filename string
	var newestUnix int64

	err := filepath.Walk(*dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() && path != *dir {
			return filepath.SkipDir
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasPrefix(info.Name(), *serviceName) || !strings.HasSuffix(info.Name(), ".dump") {
			return nil
		}
		ss := strings.Split(info.Name(), ".")
		if len(ss) != 2 {
			return nil
		}
		ss = strings.Split(ss[0], "_")
		if len(ss) != 3 {
			return nil
		}
		t, err := time.Parse("20060102150405", ss[2])
		if err != nil {
			return nil
		}
		if t.Unix() > newestUnix {
			newestUnix = t.Unix()
			filename = path
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	fmt.Printf("filename=%s\n", filename)
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	lines := bytes.Split(content, []byte{'\n'})
	var ret []pprofplus.Info
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		var info pprofplus.Info
		if err := json.Unmarshal(line, &info); err != nil {
			return nil, err
		}
		ret = append(ret, info)
	}
	return ret, nil
}

func requestHandler(writer http.ResponseWriter, _ *http.Request) {
	infos, err := getInfos()
	if err != nil {
		_, _ = writer.Write([]byte(err.Error()))
		return
	}

	var x []string
	var Sys []float64
	var HeapSys []float64
	var HeapAlloc []float64
	var HeapInuse []float64
	var HeapReleased []float64
	var HeapIdle []float64
	//var HeapIdleMinusRleased []float64
	var VMS []float64
	var RSS []float64
	for _, info := range infos {
		x = append(x, time.Unix(info.Timestamp, 0).Format("01-02 15:04:05"))
		Sys = append(Sys, calcMemWithUnit(info.Sys, memUnit))
		HeapSys = append(HeapSys, calcMemWithUnit(info.HeapSys, memUnit))
		HeapAlloc = append(HeapAlloc, calcMemWithUnit(info.HeapAlloc, memUnit))
		HeapInuse = append(HeapInuse, calcMemWithUnit(info.HeapInuse, memUnit))
		HeapReleased = append(HeapReleased, calcMemWithUnit(info.HeapReleased, memUnit))
		HeapIdle = append(HeapIdle, calcMemWithUnit(info.HeapIdle, memUnit))
		//HeapIdleMinusRleased = append(HeapIdleMinusRleased, calcMemWithUnit(info.ms.HeapIdle-info.ms.HeapReleased, option.MemUint))
		VMS = append(VMS, calcMemWithUnit(info.VMS, memUnit))
		RSS = append(RSS, calcMemWithUnit(info.RSS, memUnit))
	}

	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.InitOpts{PageTitle: "pprofplus", Theme: charts.ThemeType.Infographic},
	)
	line.Title = "单位：" + memUnitBrief[memUnit]

	line.AddXAxis(x)
	opts := charts.LineOpts{Smooth: true}
	line.AddYAxis("Sys", Sys, opts)
	line.AddYAxis("HeapSys", HeapSys, opts)
	line.AddYAxis("HeapAlloc", HeapAlloc, opts)
	line.AddYAxis("HeapInuse", HeapInuse, opts)
	line.AddYAxis("HeapReleased", HeapReleased, opts)
	line.AddYAxis("HeapIdle", HeapIdle, opts)
	//line.AddYAxis("HeapIdleMinusRleased", HeapIdleMinusRleased, opts)
	line.AddYAxis("VMS", VMS, opts)
	line.AddYAxis("RSS", RSS, opts)
	line.Render(writer)
}

func calcMemWithUnit(nByte uint64, unit MemUnit) float64 {
	switch unit {
	case MemUnitByte:
		return float64(nByte)
	case MemUnitKByte:
		return float64(nByte) / 1024
	case MemUnitMByte:
		return float64(nByte) / 1024 / 1024
	case MemUnitGByte:
		return float64(nByte) / 1024 / 1024
	}
	panic("never reach here.")
}

func init() {
	memUnitBrief = []string{"wrong", "Byte", "KByte", "MByte", "GByte"}
}

package pprofplus

import (
	"github.com/shirou/gopsutil/process"
	"os"
	"runtime"
	"time"
)

type capture struct {
	captureIntervalSec int
}

func NewCapture(captureIntervalSec int) *capture {
	return &capture{
		captureIntervalSec: captureIntervalSec,
	}
}

func (c *capture) doAsync() chan Info {
	ret := make(chan Info)
	go func() {
		p := process.Process{
			Pid: int32(os.Getpid()),
		}

		c.do(p)

		t := time.Tick(time.Second * time.Duration(c.captureIntervalSec))
		for {
			select {
			case <-t:
				ret <- c.do(p)
			}
		}
	}()
	return ret
}

func (c *capture) do(p process.Process) Info {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	mis, _ := p.MemoryInfo()

	info := Info{
		Timestamp:    time.Now().Unix(),
		Sys:          ms.Sys,
		HeapSys:      ms.HeapSys,
		HeapAlloc:    ms.HeapAlloc,
		HeapInuse:    ms.HeapInuse,
		HeapReleased: ms.HeapReleased,
		HeapIdle:     ms.HeapIdle,
		VMS:          mis.VMS,
		RSS:          mis.RSS,
	}
	return info
}

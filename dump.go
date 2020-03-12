// Copyright 2020, Chef.  All rights reserved.
// https://github.com/q191201771/pprofplus
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package pprofplus

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"
)

type dump struct {
	fp *os.File
}

func NewDump(dir string, serviceName string) (*dump, error) {
	if err := os.MkdirAll(dir, 0666); err != nil {
		return nil, err
	}

	filename := fmt.Sprintf("%s_%d_%s.dump",
		serviceName, os.Getpid(), time.Now().Format("20060102150405"))
	filename = path.Join(dir, filename)

	fp, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	return &dump{
		fp: fp,
	}, nil
}

func (dump *dump) do(info Info) {
	raw, _ := json.Marshal(&info)
	_, _ = dump.fp.Write(raw)
	_, _ = dump.fp.Write([]byte{'\n'})
	_ = dump.fp.Sync()
}

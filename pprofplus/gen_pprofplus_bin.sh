#!/usr/bin/env bash

if [ ! -d "output/" ];then
    mkdir output
fi

export GOARCH=amd64

cd example
export GOOS=darwin
go build -o ../output/example.macos
export GOOS=linux
go build -o ../output/example.linux
cd -

cd app/dashboard
export GOOS=darwin
go build -o ../../output/dashboard.macos
export GOOS=linux
go build -o ../../output/dashboard.linux
cd -

if [ -d "../pprofplus.bin/" ];then
    cp output/* ../pprofplus.bin/
    GitCommitLog=`git log --pretty=oneline -n 1`
    GitCommitLog=${GitCommitLog//\'/\"}
    echo ${GitCommitLog} > ../pprofplus.bin/pprofplus.version
fi

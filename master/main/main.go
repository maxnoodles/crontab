package main

import (
	"demo/crontab/master"
	"flag"
	"fmt"
	"runtime"
)

var (
	confFile string
)

// 解析命令行参数
func initArgs() {
	// master -config ./master.json
	flag.StringVar(&confFile, "config", "./master.json", "传入指定的 master.json")
	flag.Parse()
}

func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {

	initArgs()

	initEnv()
	err := master.InitConfig(confFile)
	if err != nil {
		goto ERR
	}

	err = master.InitJobMgr()
	if err != nil {
		goto ERR
	}

	err = master.InitApiServer()
	if err != nil {
		goto ERR
	}

ERR:
	fmt.Println(err)
}

package main

import (
	"demo/crontab/master"
	"fmt"
	"runtime"
)

func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	initEnv()
	err := master.InitConfig("")
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

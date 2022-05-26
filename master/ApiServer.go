package master

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"
)

// ApiServer 任务的 http 接口
type ApiServer struct {
	httpServer *http.Server
}

var G_apiServer *ApiServer

func handleJobSave(w http.ResponseWriter, r *http.Request) {

}

//初始化服务
func InitApiServer() (err error) {
	mux := http.NewServeMux()
	mux.HandleFunc("/job/save", handleJobSave)

	// 启动 TCP 监听
	listen, err := net.Listen("tcp", ":"+strconv.Itoa(G_config.ApiPort))
	if err != nil {
		fmt.Println(err)
	}

	//创建一个 HTTP 服务
	httpServer := &http.Server{
		ReadHeaderTimeout: time.Duration(G_config.ApiReadTimeOut) * time.Millisecond,
		WriteTimeout:      time.Duration(G_config.ApiWriteTimeOut) * time.Millisecond,
		Handler:           mux,
	}

	G_apiServer = &ApiServer{
		httpServer: httpServer,
	}

	go httpServer.Serve(listen)
	return
}

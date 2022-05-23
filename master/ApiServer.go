package master

import (
	"fmt"
	"net"
	"net/http"
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
	listen, err := net.Listen("tcp", ":8070")
	if err != nil {
		fmt.Println(err)
	}

	//创建一个 HTTP 服务
	httpServer := &http.Server{
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
		Handler:           mux,
	}

	G_apiServer = &ApiServer{
		httpServer: httpServer,
	}

	go httpServer.Serve(listen)
	return
}

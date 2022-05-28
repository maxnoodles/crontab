package master

import (
	"demo/crontab/common"
	"encoding/json"
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

// 保存任务接口
// POSTjob = {"name": job, "command": "echo hello", "cronExpr": "* * * * * *"}
func handleJobSave(resp http.ResponseWriter, req *http.Request) {
	var (
		err     error
		postJob string
		job     common.Job
		oldJob  *common.Job
		bytes   []byte
	)

	err = req.ParseForm()
	if err != nil {
		goto ERR
	}

	postJob = req.PostForm.Get("job")
	if err = json.Unmarshal([]byte(postJob), &job); err != nil {
		goto ERR
	}
	// 保存到 etcd
	if oldJob, err = G_jobMgr.SaveJob(&job); err != nil {
		goto ERR
	}

	if bytes, err = common.BuildResponse(0, "success", oldJob); err == nil {
		_, err := resp.Write(bytes)
		if err != nil {
			goto ERR
		}
	}
	return

ERR:
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}

// 删除任务接口
// POST /job/delete name=job1
func handleJobDelete(resp http.ResponseWriter, req *http.Request) {
	var (
		err    error // interface
		name   string
		oldJob *common.Job
		bytes  []byte
	)

	// POST: a=1&b=2&c=3
	if err = req.ParseForm(); err != nil {
		goto ERR
	}

	// 删除的任务名
	name = req.PostForm.Get("name")

	// 去删除任务
	if oldJob, err = G_jobMgr.DeleteJob(name); err != nil {
		goto ERR
	}

	// 正常应答
	if bytes, err = common.BuildResponse(0, "success", oldJob); err == nil {
		resp.Write(bytes)
	}

	return

ERR:
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}

// 列举所有crontab任务
func handleJobList(resp http.ResponseWriter, req *http.Request) {
	fmt.Println(123)
	var (
		jobList []*common.Job
		err     error
		bytes   []byte
	)

	//fmt.Println("handleJobList")

	// 获取任务列表
	if jobList, err = G_jobMgr.ListJobs(); err != nil {
		goto ERR
	}

	// 返回正常应答
	if bytes, err = common.BuildResponse(0, "success", jobList); err == nil {
		resp.Write(bytes)
	}

	return

ERR:
	// 返回错误应答
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}

// 强制杀死某个任务
// POST /job/kill name=job1
func handleJobKill(resp http.ResponseWriter, req *http.Request) {
	var (
		err   error
		name  string
		bytes []byte
	)

	// 解析POST表单
	if err = req.ParseForm(); err != nil {
		goto ERR
	}

	// 要杀死的任务名
	name = req.PostForm.Get("name")

	if err = G_jobMgr.KillJob(name); err != nil {
		goto ERR
	}

	// 返回正常应答
	if bytes, err = common.BuildResponse(0, "success", nil); err == nil {
		resp.Write(bytes)
	}

	return

ERR:
	// 返回错误应答
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}

// InitApiServer 初始化服务
func InitApiServer() (err error) {
	var (
		mux           *http.ServeMux
		listener      net.Listener
		httpServer    *http.Server
		staticDir     http.Dir     // 静态文件根目录
		staticHandler http.Handler // 静态文件HTTP回调
	)
	// 配置路由 根据url 来调用相应的回调函数
	// TCP过来之后会来路由这里
	mux = http.NewServeMux()
	mux.HandleFunc("/job/save", handleJobSave)
	mux.HandleFunc("/job/delete", handleJobDelete)
	mux.HandleFunc("/job/list", handleJobList)
	mux.HandleFunc("/job/kill", handleJobKill)

	// /index.html 匹配最长的,既 最大匹配

	// 静态文件目录
	staticDir = http.Dir(G_config.WebRoot)
	staticHandler = http.FileServer(staticDir)
	mux.Handle("/", http.StripPrefix("/", staticHandler))

	// 启动TCP监听
	if listener, err = net.Listen("tcp", ":"+strconv.Itoa(G_config.ApiPort)); err != nil {
		return
	}

	// 创建一个HTTP服务
	httpServer = &http.Server{
		ReadTimeout:  time.Duration(G_config.ApiReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(G_config.ApiWriteTimeout) * time.Millisecond,
		Handler:      mux,
	}

	// 赋值单例
	G_apiServer = &ApiServer{
		httpServer: httpServer,
	}

	// 启动了服务端
	go httpServer.Serve(listener)

	return
}

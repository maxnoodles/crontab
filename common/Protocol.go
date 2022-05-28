package common

import "encoding/json"

//定时任务

type Job struct {
	Name     string `json:"name"`
	Commond  string `json:"commond"`
	CronExpr string `json:"cronExpr"`
}

type Response struct {
	Errno int         `json:"errno"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
}

func BuildResponse(errno int, msg string, data interface{}) (resp []byte, err error) {
	var response Response

	response.Errno = errno
	response.Msg = msg
	response.Data = data
	return json.Marshal(response)
}

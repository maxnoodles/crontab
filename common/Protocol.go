package common

//定时任务

type Job struct {
	Name     string `json:"name"`
	Commond  string `json:"commond"`
	CronExpr string `json:"cronExpr"`
}

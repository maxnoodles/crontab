package worker

import (
	"demo/crontab/common"
	"fmt"
	"time"
)

type Scheduler struct {
	jobEventChan chan *common.JobEvent
	jobPlanTable map[string]*common.JobSchedulePlan
}

func (s Scheduler) PushJobEvent(jobEvent *common.JobEvent) {
	fmt.Println("jobEvent", jobEvent)
	s.jobEventChan <- jobEvent
}

func (s Scheduler) handleJobEvent(jobEvent *common.JobEvent) {
	var (
		jobSchedulePlan *common.JobSchedulePlan
		//jobExecuteInfo *common.JobExecuteInfo
		//jobExecuting bool
		jobExisted bool
		err        error
	)
	switch jobEvent.EventType {
	case common.JOB_EVENT_SAVE:
		fmt.Println("to save")
		if jobSchedulePlan, err = common.BuildJobSchedulePlan(jobEvent.Job); err != nil {
			return
		}
		s.jobPlanTable[jobEvent.Job.Name] = jobSchedulePlan
	case common.JOB_EVENT_DELETE:
		if jobSchedulePlan, jobExisted = s.jobPlanTable[jobEvent.Job.Name]; jobExisted {
			delete(s.jobPlanTable, jobEvent.Job.Name)
		}
	case common.JOB_EVENT_KILL:
		// todo: 强杀任务
	}
}

// 重新计算任务调度状态
func (s Scheduler) TryScheduler() (scheduleAfter time.Duration) {
	var (
		jobPlan  *common.JobSchedulePlan
		now      time.Time
		nearTime *time.Time
	)
	// 如果任务列表为空，随便睡眠多久
	if len(s.jobPlanTable) == 0 {
		scheduleAfter = 1 * time.Second
		return
	}

	//1. 遍历所有任务
	now = time.Now()
	for _, jobPlan = range s.jobPlanTable {
		if jobPlan.NextTime.Before(now) || jobPlan.NextTime.Equal(now) {
			s.TryStartJob(jobPlan)
			jobPlan.NextTime = jobPlan.Expr.Next(now)
		}

		if nearTime == nil || jobPlan.NextTime.Before(*nearTime) {
			nearTime = &jobPlan.NextTime
		}

	}
	scheduleAfter = nearTime.Sub(now)
	return
}

func (s Scheduler) schedulerLoop() {
	fmt.Println("in loop")
	schedulerAfter := s.TryScheduler()
	scheduleTimer := time.NewTimer(schedulerAfter)
	for {
		select {
		case jobEvent := <-s.jobEventChan:
			fmt.Println(jobEvent)
			// 对内存中维护的任务队列做增删改查
			s.handleJobEvent(jobEvent)
		case <-scheduleTimer.C:
		}
		schedulerAfter := s.TryScheduler()
		scheduleTimer.Reset(schedulerAfter)
	}
}

func (s Scheduler) TryStartJob(plan *common.JobSchedulePlan) {
	fmt.Println("执行任务", plan.Job)
}

var (
	G_scheduler *Scheduler
)

func InitScheduler() (err error) {
	G_scheduler = &Scheduler{
		jobEventChan: make(chan *common.JobEvent, 1000),
		jobPlanTable: make(map[string]*common.JobSchedulePlan),
		//jobExcutingTable:make(map[string]*common.JobExecuteInfo),
		//jobResultChan:make(chan *common.JobExecuteResult,1000),
	}

	// 启动调度协程
	go G_scheduler.schedulerLoop()
	return
}

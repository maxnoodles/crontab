package master

import (
	"context"
	"demo/crontab/common"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

type JobMgr struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

var (
	G_jobMgr *JobMgr
)

func InitJobMgr() (err error) {
	config := clientv3.Config{
		Endpoints:   G_config.EtcdEndpoints,
		DialTimeout: time.Duration(G_config.EtcdDialTimeout) * time.Millisecond,
	}

	client, err := clientv3.New(config)
	if err != nil {
		return
	}

	kv := clientv3.KV(client)
	lease := clientv3.Lease(client)
	if err != nil {
		return err
	}

	G_jobMgr = &JobMgr{
		client: client,
		kv:     kv,
		lease:  lease,
	}

	return
}

func (JobMgr *JobMgr) SaveJob(job *common.Job) (oldJob *common.Job, err error) {

	//保存到 /cron/jobs/ 下
	jobKey := common.JOB_SAVE_DIR + job.Name

	// 更新到etcd
	jobValue, err := json.Marshal(job)
	if err != nil {
		return
	}

	putResp, err := JobMgr.kv.Put(context.TODO(), jobKey, string(jobValue), clientv3.WithPrevKV())
	if err != nil {
		return
	}

	if putResp.PrevKv != nil {
		oldJobObj := common.Job{}
		if err := json.Unmarshal(putResp.PrevKv.Value, &oldJobObj); err != nil {
			err = nil
		}
		oldJob = &oldJobObj
	}
	return
}

// 删除任务
func (jobMgr *JobMgr) DeleteJob(name string) (oldJob *common.Job, err error) {
	var (
		jobKey    string
		delResp   *clientv3.DeleteResponse
		oldJobObj common.Job
	)

	// etcd中保存任务的key
	jobKey = common.JOB_SAVE_DIR + name

	// 从etcd中删除
	if delResp, err = jobMgr.kv.Delete(context.TODO(), jobKey, clientv3.WithPrevKV()); err != nil {
		return
	}

	// 返回被删除的任务信息
	if len(delResp.PrevKvs) != 0 {
		// 旧值解析
		if err = json.Unmarshal(delResp.PrevKvs[0].Value, &oldJobObj); err != nil {
			// 如果删除成功的话 能不能得到一个旧值其实无所谓，所以这里把err强行置为nil
			err = nil
			return
		}
		oldJob = &oldJobObj
	}

	return
}

// 列举任务
func (jobMgr *JobMgr) ListJobs() (jobList []*common.Job, err error) {
	var (
		dirKey  string
		getResp *clientv3.GetResponse
		kvPair  *mvccpb.KeyValue
		job     *common.Job
	)

	// 任务保存目录
	dirKey = common.JOB_SAVE_DIR

	// 获取目录下所有任务信息
	if getResp, err = jobMgr.kv.Get(context.TODO(), dirKey, clientv3.WithPrefix()); err != nil {
		return
	}

	// 初始化数组空间
	jobList = make([]*common.Job, 0)
	// len(jobLlist) == 0

	//fmt.Println(len(getResp.Kvs))

	// 遍历所有任务 进行反序列化
	for _, kvPair = range getResp.Kvs {
		job = &common.Job{}
		if err = json.Unmarshal(kvPair.Value, job); err != nil {
			fmt.Println(err)
			err = nil
			continue
		}
		jobList = append(jobList, job)
	}

	return
}

// 杀死任务
func (jobMgr *JobMgr) KillJob(name string) (err error) {
	// 原理，更新一下key=/cron/killer/任务名、
	var (
		killerKey      string
		leaseGrantResp *clientv3.LeaseGrantResponse
		leaseId        clientv3.LeaseID
	)

	// 通知worker杀死对应任务
	killerKey = common.JOB_KILLER_DIR + name

	// 让worker监听到一次put操作，创建一个租约让其稍后自动过期
	if leaseGrantResp, err = jobMgr.lease.Grant(context.TODO(), 1); err != nil {
		return
	}

	// 租约ID
	leaseId = leaseGrantResp.ID

	// 设置killer标记
	if _, err = jobMgr.kv.Put(context.TODO(), killerKey, "", clientv3.WithLease(leaseId)); err != nil {
		return
	}

	return
}

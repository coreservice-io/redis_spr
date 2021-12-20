package USpr

import (
	"errors"
	"github.com/go-redis/redis/v8"
	"math/rand"
	"sync"
	"time"
)

type SprJobMgr struct {
	jobMap      sync.Map
	redisClient *redis.ClusterClient
}

type RedisConfig struct {
	Addr     string
	Port     int
	UserName string
	Password string
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func New(config RedisConfig) (*SprJobMgr, error) {
	rds, err := initRedisClient(config.Addr, config.Port, config.UserName, config.Password)
	if err != nil {
		return nil, errors.New("redis connect error")
	}
	sMgr := &SprJobMgr{
		redisClient: rds,
	}

	return sMgr, nil
}

func (smgr *SprJobMgr) AddSprJob(jobName string) error {
	_, exist := smgr.jobMap.Load(jobName)
	if exist {
		return errors.New("job already exist")
	}
	//new job
	job := newJob(jobName, smgr)
	smgr.jobMap.Store(jobName, job)
	//start loop
	job.startLoop()
	return nil
}

func (smgr *SprJobMgr) RemoveSprJob(jobName string) {
	job, exist := smgr.jobMap.Load(jobName)
	if !exist {
		return
	}
	//stop
	job.(*SprJob).stopLoop()
	//delete
	smgr.jobMap.Delete(jobName)
}

func (smgr *SprJobMgr) IsMaster(jobName string) bool {
	job, exist := smgr.jobMap.Load(jobName)
	if !exist {
		return false
	}
	return job.(*SprJob).IsMaster
}

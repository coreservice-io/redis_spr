package RedisSpr

import (
	"errors"
	"io"
	"math/rand"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/universe-30/ULog"
)

type SprJobMgr struct {
	jobMap      sync.Map
	redisClient *redis.ClusterClient

	logger ULog.Logger
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
		logger:      ULog.New(),
	}

	return sMgr, nil
}

func (smgr *SprJobMgr) SetOutPut(w io.Writer) {
	smgr.logger.SetOutput(w)
}

func (smgr *SprJobMgr) SetLevel(v ULog.LogLevel) {
	smgr.logger.SetLevel(v)
}

func (smgr *SprJobMgr) SetLogger(logger ULog.Logger) {
	smgr.logger = logger
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
		smgr.logger.Debugln(jobName, "is not exist")
		return false
	}
	return job.(*SprJob).IsMaster
}

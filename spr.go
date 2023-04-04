package redis_spr

import (
	"context"
	"errors"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/coreservice-io/log"
	"github.com/go-redis/redis/v8"
)

type SprJobMgr struct {
	jobMap      sync.Map
	redisClient *redis.ClusterClient
	prefix      string
	logger      log.Logger
}

type RedisConfig struct {
	Addr     string
	Port     int
	UserName string
	Password string
	Prefix   string
	UseTLS   bool
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func New(config RedisConfig) (*SprJobMgr, error) {
	rds, err := initRedisClient(config.Addr, config.Port, config.UserName, config.Password, config.UseTLS)
	if err != nil {
		return nil, errors.New("redis connect error")
	}

	prefix := config.Prefix
	if prefix != "" && !strings.HasSuffix(prefix, ":") {
		prefix = prefix + ":"
	}

	sMgr := &SprJobMgr{
		redisClient: rds,
		logger:      nil,
		prefix:      prefix,
	}

	return sMgr, nil
}

func (smgr *SprJobMgr) SetLogger(logger log.Logger) {
	smgr.logger = logger
}

func (smgr *SprJobMgr) GetLogger() log.Logger {
	return smgr.logger
}

func (smgr *SprJobMgr) AddSprJob(jobName string) error {
	_, exist := smgr.jobMap.Load(jobName)
	if exist {
		return errors.New("job already exist")
	}
	//new job
	job := newJob(context.Background(), jobName, smgr)
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
		if smgr.logger != nil {
			smgr.logger.Debugln("<USpr>", jobName, "is not exist")
		}
		return false
	}
	return job.(*SprJob).IsMaster
}

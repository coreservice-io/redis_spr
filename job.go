package USpr

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/universe-30/USafeGo"
	"math/rand"
	"time"
)

const LoopIntervalSec = 15
const MasterKeepTime = 90

type SprJob struct {
	JobName         string
	IsMaster        bool
	JobRand         string
	LoopIntervalSec int
	StopFlag        bool
	LastRuntime     int64
	sprJobMgr       *SprJobMgr
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func newJob(name string, sprMgr *SprJobMgr) *SprJob {
	s := &SprJob{
		JobName:         "spr:" + name,
		IsMaster:        false,
		JobRand:         fmt.Sprintf("%d", rand.Intn(100000000)+1),
		LoopIntervalSec: LoopIntervalSec,
		StopFlag:        false,
		LastRuntime:     0,
		sprJobMgr:       sprMgr,
	}
	return s
}

func (s *SprJob) startLoop() {
	USafeGo.GoInfiniteLoop(func() bool {
		if s.StopFlag {
			return false
		}
		s.run()
		return true
	}, nil, s.LoopIntervalSec, 15)
}

func (s *SprJob) stopLoop() {
	s.StopFlag = true
	s.IsMaster = false
}

func (s *SprJob) run() {
	//log.Println(s.JobName, "loop job run")
	if s.sprJobMgr.redisClient == nil {
		s.IsMaster = false
		return
	}

	//check jobname in redis
	value, err := s.sprJobMgr.redisClient.Get(context.Background(), s.JobName).Result()

	//get value
	if err == nil {
		//value error
		if value != s.JobRand {
			s.IsMaster = false
			return
		}

		//value==jobRand
		//keep master token
		s.IsMaster = true
		s.sprJobMgr.redisClient.Expire(context.Background(), s.JobName, time.Second*time.Duration(MasterKeepTime))

	} else if err == redis.Nil {
		//if no value
		success, err := s.sprJobMgr.redisClient.SetNX(context.Background(), s.JobName, s.JobRand, time.Second*time.Duration(MasterKeepTime)).Result()
		if err != nil || !success {
			s.IsMaster = false
			return
		}
		s.IsMaster = true
	} else {
		//other err
		s.IsMaster = false
		return
	}
}

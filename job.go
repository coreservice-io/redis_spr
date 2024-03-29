package redis_spr

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/go-redis/redis/v8"
)

const loopIntervalSec = 15
const redoDelaySecs = 15
const masterKeepTime = 90

type SprJob struct {
	JobName         string
	IsMaster        bool
	JobRand         string
	LoopIntervalSec int
	Ctx             context.Context
	LastRuntime     int64
	sprJobMgr       *SprJobMgr
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func newJob(ctx context.Context, name string, sprMgr *SprJobMgr) *SprJob {
	s := &SprJob{
		JobName:         sprMgr.prefix + "spr:" + name,
		IsMaster:        false,
		JobRand:         fmt.Sprintf("%d", rand.Intn(100000000)+1),
		LoopIntervalSec: loopIntervalSec,
		LastRuntime:     0,
		sprJobMgr:       sprMgr,
		Ctx:             ctx,
	}
	return s
}

func (s *SprJob) startLoop() {
	goInfiniteLoop(func() bool {
		select {
		case <-s.Ctx.Done():
			s.IsMaster = false
			s.sprJobMgr.jobMap.Delete(s.JobName)
			return false
		default:
			s.run()
			return true
		}
	}, nil, s.LoopIntervalSec, redoDelaySecs)
}

func (s *SprJob) run() {

	if s.sprJobMgr.redisClient == nil {
		s.IsMaster = false
		return
	}

	//check job name in redis
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
		s.sprJobMgr.redisClient.Expire(context.Background(), s.JobName, time.Second*time.Duration(masterKeepTime))

	} else if err == redis.Nil {
		//if no value
		success, err := s.sprJobMgr.redisClient.SetNX(context.Background(), s.JobName, s.JobRand, time.Second*time.Duration(masterKeepTime)).Result()
		if !success {
			s.IsMaster = false
			if err != nil && s.sprJobMgr.logger != nil {
				s.sprJobMgr.logger.Errorln("<USpr>", err)
			}
			return
		}
		s.IsMaster = true
	} else {
		//other err
		if s.sprJobMgr.logger != nil {
			s.sprJobMgr.logger.Errorln("<USpr>", err)
		}
		s.IsMaster = false
		return
	}
}

func goInfiniteLoop(function func() bool, onPanic func(err interface{}), interval int, redoDelaySec int) {
	runToken := make(chan struct{})
	stopSignal := make(chan struct{})
	go func() {
		for {
			select {
			case <-runToken:
				go func() {
					defer func() {
						if err := recover(); err != nil {
							if onPanic != nil {
								onPanic(err)
							}
							time.Sleep(time.Duration(redoDelaySec) * time.Second)
							runToken <- struct{}{}
						}
					}()
					for {
						isGoOn := function()
						if !isGoOn {
							stopSignal <- struct{}{}
							return
						}
						time.Sleep(time.Duration(interval) * time.Second)
					}
				}()
			case <-stopSignal:
				return
			}
		}
	}()
	runToken <- struct{}{}
}

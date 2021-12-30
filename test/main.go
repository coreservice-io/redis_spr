package main

import (
	"log"
	"time"

	"github.com/universe-30/ULog"
	"github.com/universe-30/USpr"
)

func main() {
	sprJobMgr, err := USpr.New(USpr.RedisConfig{
		Addr: "127.0.0.1",
		Port: 6379,
	})
	if err != nil {
		log.Fatalln(err)
	}

	sprJobMgr.SetLevel(ULog.DebugLevel)
	//USpr use package github.com/universe-30/ULog as default logger
	//You can log to other target by using function SetOutPut(w io.Writer)
	//or use function SetLogger(logger ULog.Logger) to use your own logger which implemented the ULog.Logger interface

	err = sprJobMgr.AddSprJob("testJob1")
	if err != nil {
		log.Println(err)
	}
	err = sprJobMgr.AddSprJob("testJob2")
	if err != nil {
		log.Println(err)
	}

	// use function IsMaster("jobName") to check whether the process get the master token or not
	// if return true means get the master token
	go func() {
		for {
			time.Sleep(time.Second)
			log.Println("testjob is master:", sprJobMgr.IsMaster("testJob1"))
			log.Println("testjob2 is master:", sprJobMgr.IsMaster("testJob2"))
		}
	}()

	// use function RemoveSprJob("jobName") to remove the job
	// removed job always return false when use IsMaster("jobName")
	time.AfterFunc(time.Second*25, func() {
		sprJobMgr.RemoveSprJob("testJob2")
	})

	time.Sleep(1 * time.Hour)
}

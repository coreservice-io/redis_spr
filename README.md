# redis_spr

Make sure specific job can only be run over 1 process among all machines and processes. If some process is killed or crashed, the job will be switched to other process on some machine.

### usage
```
go get "github.com/universe-30/redis_spr"
```

```go
 package main

import (
	"context"
	"log"
	"time"

	"github.com/coreservice-io/redis_spr"
)

func main() {
	sprJobMgr, err := redis_spr.New(redis_spr.RedisConfig{
		Addr:   "127.0.0.1",
		Port:   6379,
		Prefix: "sprExample",
	})
	if err != nil {
		log.Fatalln(err)
	}

	//or use function SetLogger(logger log.Logger) to use your own logger which implemented the log.Logger interface
	err = sprJobMgr.AddSprJob(context.Background(), "testJob1")
	if err != nil {
		log.Println(err)
	}
	err = sprJobMgr.AddSprJob(context.Background(), "testJob2")
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

	time.Sleep(1 * time.Hour)
}


```

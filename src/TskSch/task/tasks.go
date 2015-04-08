package task

import (
	"TskSch/command"
	"TskSch/execute"
	"TskSch/msgQ"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"gopkg.in/mgo.v2"
	"os"
	"sync"
	"time"
)

const ConcLimit int = 3

func Execute(file *os.File, session *mgo.Session, Conn redis.Conn) {
	for {
		var (
			Wg     sync.WaitGroup
			Idlist []string
		)
		Size := msgQ.Size(Conn)
		if Size <= ConcLimit {
			for i := 0; i < Size; i++ {
				Idlist = append(Idlist,msgQ.Pop(Conn))
			}
		} else {
			for i := 0; i < ConcLimit; i++ {
				Idlist = append(Idlist,msgQ.Pop(Conn))
			}
		}
		for _, Id := range Idlist {
			if Id != "" {
				Wg.Add(1)
				//SEARCHING THE COMMAND BASED ON ID IN FILE
				cmd := command.Search(Conn,Id)
				if cmd != "" {
					//EXECUTING THE COMMAND CONCURRENTLY
					args := Id + ":" + cmd 
					go execute.Exec(file, session, &Wg, args)
					Wg.Wait()
				} else {
					fmt.Println("COMMAND IS NOT ASSIGNED TO IT'S cmd_id ")
				}
			} else {
				fmt.Println("TASK_ID IS NOT ASSIGNED IN THE QUEUE")
			}
		}
		time.Sleep(time.Second * 5)
	}
}



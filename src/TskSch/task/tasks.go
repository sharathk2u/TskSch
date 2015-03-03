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
	fmt.Println(file)
	for {
		var (
			Wg     sync.WaitGroup
			Idlist [ConcLimit]string
		)
		Size := msgQ.Size(Conn)
		if Size <= ConcLimit {
			for i := 0; i < Size; i++ {
				Idlist[i] = msgQ.Pop(Conn)
			}
		} else {
			for i := 0; i < ConcLimit; i++ {
				Idlist[i] = msgQ.Pop(Conn)
			}
		}
		for _, Id := range Idlist {
			if Id != "" {
				Wg.Add(1)
				//SEARCHING THE COMMAND BASED ON ID IN FILE
				cmd := command.Search(Conn, Id)
				if cmd != "" {
					//EXECUTING THE COMMAND CONCURRENTLY
					go execute.Exec(file, session, &Wg, cmd)
				} else {
					fmt.Println("COMMAND IS NOT ASSIGNED TO IT'S cmd_id ")
				}
			} else {
				fmt.Println("TASK_ID IS NOT ASSIGNED IN THE QUEUE")
			}
		}
		Wg.Wait()
		time.Sleep(time.Second * 5)
	}
}

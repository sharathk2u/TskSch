package task

import (
	"TskSch/command"
	"TskSch/execute"
	"TskSch/msgQ"
	"TskSch/logger"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"gopkg.in/mgo.v2"
	"os"
	"sync"
	"code.google.com/p/goconf/conf"
	"time"
)

var ConcLimit int
var y time.Time
var managerPath string
var host string
var name string
func Execute(file *os.File, session *mgo.Session, Conn redis.Conn,logfile *os.File) {
	for {
		var (
			Wg     sync.WaitGroup
			Idlist []string
		)

		//Extracting conf file 
		Finfo, _ := os.Stat("../server.conf")
		v := Finfo.ModTime()
		if !v.Equal(y) {
			c, err := conf.ReadConfigFile("../server.conf")
			if err != nil {
				fmt.Println("CAN'T READ CONF FIILE",err)
			}
			p, _ := c.GetInt("taskagent","Conclimit")
			name , _ = c.GetString("taskagent","name")
			ConcLimit = p
			x, _ := c.GetString("manager","host")
			y, _ := c.GetString("manager","port")
			managerPath = x + ":" + y
			host ,_ = c.GetString("resultDB","host")
		}

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
				LogInfo := logger.Info(logfile)
				LogInfo.Println("ASKING THE COMMAND BASED ON ID :"+ Id,"TO MANAGER")
				//SEARCHING THE COMMAND BASED ON ID
				cmd := command.Search(Conn,Id,managerPath,host,name)
				LogInfo.Println("GOT COMMAND BASED ON ID :"+ Id,"FROM MANAGER")
				if cmd != "" {
					//EXECUTING THE COMMAND CONCURRENTLY
					LogInfo.Println("EXECUTING THE COMMAND CONCURRENTLY")
					args := Id + ":" + cmd
					go execute.Exec(file, session, &Wg, args,logfile)
					Wg.Wait()
				} else {
					fmt.Println("COMMAND IS NOT ASSIGNED TO IT'S cmd_id ")
					LogErr := logger.Failure(logfile)
					LogErr.Println("COMMAND IS NOT ASSIGNED TO IT'S cmd_id ")
				}
			} else {
				fmt.Println("TASK_ID IS NOT PRESENT IN THE QUEUE")
				LogErr := logger.Failure(logfile)
				LogErr.Println("TASK_ID IS NOT PRESENT IN THE QUEUE")
			}
		}
		y = v
	}
}



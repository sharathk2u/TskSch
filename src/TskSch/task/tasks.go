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
	"code.google.com/p/goconf/conf"
	"time"
)

var ConcLimit int
var y time.Time
var schedulerPath string
var username string
var password string
var host string

func Execute(file *os.File, session *mgo.Session, Conn redis.Conn) {
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
			ConcLimit = p
			x, _ := c.GetString("scheduler","host")
			y, _ := c.GetString("scheduler","port")
			schedulerPath = x + ":" + y
			username ,_ = c.GetString("resultDB","username")
			password ,_ = c.GetString("resultDB","password")
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
				//SEARCHING THE COMMAND BASED ON ID IN FILE
				cmd := command.Search(Conn,Id,schedulerPath,username,password,host)
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
		y = v
	}
}



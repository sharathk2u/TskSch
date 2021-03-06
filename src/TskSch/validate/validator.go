package main

import (
	"TskSch/msgQ"
	"TskSch/logger"
	"TskSch/resultDB"
	"TskSch/mailer"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"bytes"
	"os/exec"
	"strconv"
	"time"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"code.google.com/p/goconf/conf"
	"os"
	"sync"
	"strings"
)

type task struct {
	Task_id string
}

type taskp struct {
	Task_id string
	Pid     int
}

var y time.Time
var schedulerHost string
var ManagerHost string
var taskAgent string
var host1 string
var host2 string
var port string
func main() {
	for _ = range time.Tick(time.Second * 50){

		logfile := logger.LogValInit()
		LogInfo := logger.Info(logfile)
		LogErr := logger.Failure(logfile)
		
		var wg sync.WaitGroup
		wg.Add(4)

		//Extracting conf
		Finfo, _ := os.Stat("../server.conf")
		v := Finfo.ModTime()
		if !v.Equal(y) {
			c, err := conf.ReadConfigFile("../server.conf")
			if err != nil {
				fmt.Println("CAN'T READ CONF FIILE",err)
			}
			w, _ := c.GetString("scheduler", "host")
			x, _ := c.GetString("scheduler", "port")
			schedulerHost = w + ":" + x
			a, _ := c.GetString("manager", "host")
			b, _ := c.GetString("manager", "port")
			ManagerHost = a + ":" + b
			p, _ := c.GetString("taskagent","host")
			z, _ := c.GetString("taskagent","port")
			taskAgent = p + ":" + z
			host1 ,_ = c.GetString("resultDB","host")
			host2 ,_ = c.GetString("msgQ","host")
			port ,_ = c.GetString("msgQ","port")
		}

		//Checking liveliness of Scheduler
		schedulerPath := "http://"+schedulerHost+"/ping"
		go func(schedulerPath string,wg *sync.WaitGroup){
			res, err := http.Get(schedulerPath)
			if err == nil {
				body , _ := ioutil.ReadAll(res.Body)
				if string(body) != "" {
					var status interface{}
					err := json.Unmarshal([]byte(body), &status)
					if err == nil {
						if status.(map[string]interface{})["status"].(string) == "alive" {
							LogInfo.Println("Scheduer is alive")	
						}else {
							mailer.Mail("GOSERVE: Regarding Scheduler Status", "Scheduer is not alive")
							LogErr.Println("Scheduer is not alive")
						}
					}
				}
			}else{
				fmt.Println("Cannot connect to Scheduler",err)
				LogErr.Println("Cannot connect to Scheduler",err)
			}
			wg.Done()
		}(schedulerPath, &wg)

		//Checking liveliness of Manager
		ManagerPath := "http://"+ManagerHost+"/ping"
		go func(ManagerPath string,wg *sync.WaitGroup){
			res, err := http.Get(ManagerPath)
			if err == nil {
				body , _ := ioutil.ReadAll(res.Body)
				if string(body) != "" {
					var status interface{}
					err := json.Unmarshal([]byte(body), &status)
					if err == nil {
						if status.(map[string]interface{})["status"].(string) == "alive" {
							LogInfo.Println("Manager is alive")	
						}else {
							mailer.Mail("GOSERVE: Regarding Scheduler Status", "Scheduer is not alive")
							LogErr.Println("Scheduer is not alive")
						}
					}
				}
			}else{
				fmt.Println("Cannot connect to Scheduler",err)
				LogErr.Println("Cannot connect to Scheduler",err)
			}
			wg.Done()
		}(ManagerPath, &wg)

		//Checking liveliness of Task Agents
		go func(taskAgent string,wg *sync.WaitGroup){
			TaskHost := strings.Split( strings.Split(taskAgent ,":")[0] , ",")
			TaskPort :=	strings.Split( strings.Split(taskAgent ,":")[1] , ",")
			if( len(TaskHost) == len(TaskPort) ) {
				for i , _ := range TaskHost {
					taskagentPath := "http://"+TaskHost[i]+":"+TaskPort[i]+"/ping"
					res, err := http.Get(taskagentPath)
					if err == nil{
						body , _ := ioutil.ReadAll(res.Body)
						if string(body) != "" {
							var status interface{}
							err := json.Unmarshal([]byte(body), &status)
							if err == nil {
								if status.(map[string]interface{})["status"].(string) == "alive" {
									fmt.Println("Task Agent is alive")
									
								}else {
									mailer.Mail("GOSERVE: Regarding Task Aegent : "+ taskagentPath + " Status", taskagentPath +" : is not alive")
									LogErr.Println("Scheduer is not alive")
								}
							}
						}
					}else{
						fmt.Println("Cannot connect to task agent",err)
						LogErr.Println("Cannot connect to Task Aegent : " + taskagentPath ,err)
					}
				}
			}else{
				fmt.Println("ERROR IN CONFIG FILE")
				LogErr.Println("ERROR IN CONFIG FILE")
			}
			wg.Done()
		}(taskAgent , &wg)

		go func(wg *sync.WaitGroup,host1 string,host2 string,port string) {
			var ids []string
			res1 := []task{}
			res2 := []taskp{}

			//Connecting to mongodb
			session := resultDB.ResultdbInit(host1)
			session.SetMode(mgo.Monotonic, true)
			col := session.DB("TskSch").C("Result")

			//Connecting to msgQ
			Conn := msgQ.RedisInit(host2 ,port)

			//Checking liveness of msgQ to get the list of taskids in msgQ
			err := msgQ.Ping(Conn)
			if err !=nil{
				Ids, Err := redis.Values(Conn.Do("LRANGE", "task", "0", "-1"))
				if Err != nil {
					fmt.Println("Could not able to connect to msgQ",Err)
					mailer.Mail("GOSERVE: Regarding 'msgQ' Status", "Could not able to connect to 'msgQ'")
					LogErr.Println("Could not able to connect to msgQ",Err)
				}else {
					for _, val := range Ids {
						ids = append(ids, string(val.([]byte)))
					}
				}
			}else{
				LogErr.Println(err)
			}

			//Checking liveliness of mongodb
			err = resultDB.Ping(session)
			if err !=nil{ 
				//collecting all the taskids from resultDB which are not executed and not in execution state to check whether they are in msgQ 
				Err := col.Find(bson.M{"executed": false, "exec_stat": false}).Select(bson.M{"task_id":1}).All(&res1)
				if Err != nil {
					fmt.Println("Could not able to connect to mongodb",Err)
					mailer.Mail("GOSERVE: Regarding 'mongodb' Status", "Could not able to connect to 'mongodb'")
					LogErr.Println("Could not able to connect to mongodb",Err)
				}
				for _, val := range res1{
					flag := In(val.Task_id , ids )
					if flag != true {
						x, err := Conn.Do("RPUSH", "task", val.Task_id)
						if err != nil {
							fmt.Println(x,err)
							LogErr.Println(x,err)
						}
						LogInfo.Println("PUSHED" + val.Task_id + "TASK TO msgQ")
					}
				}

				//collecting all the taskids from resultDB which are not executed , not in execution state and which were poped from executer but not executed
				Err = col.Find(bson.M{"executed": false, "exec_stat": false ,"pid":bson.M{"$gt":0}}).Select(bson.M{"task_id":1,"pid":1}).All(&res2)
				if Err != nil {
					fmt.Println("Could not able to connect to mongodb",Err)
					mailer.Mail("GOSERVE: Regarding 'mongodb' Status", "Could not able to connect to 'mongodb'")
					LogErr.Println("Could not able to connect to mongodb",Err)
				}
				for _ , val := range res2{
					flag := Isalive(val.Pid)
					if flag != true {
						flag1 := In(val.Task_id , ids )
						if flag1 != true {
							x, err := Conn.Do("RPUSH", "task", val.Task_id)
							if err != nil {
								fmt.Println(x,err)
								LogErr.Println(x,err)
							}
							fmt.Println("PUSHED" + val.Task_id + "TASK TO msgQ" )
							LogInfo.Println("PUSHED" + val.Task_id + "TASK TO msgQ")
						}
					}
				}
			}else{
				resultDB.Restart(session)
				LogInfo.Println("RESULT DB RESTARTED")
			}
			wg.Done()
			fmt.Println("X!")
		}(&wg,host1,host2,port)
		wg.Wait()
		y = v
	}
}

// Function to check aliveness of the executer
func Isalive( pid int ) bool {
	cmd := "kill -0 " + strconv.Itoa(pid)
	var Errout bytes.Buffer
	cmds := exec.Command("sh", "-c",cmd )
	cmds.Stderr = &Errout
	cmds.Run()
	if Errout.String() != ""{
		return false
	}
	return true
}

//helper function
func In(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

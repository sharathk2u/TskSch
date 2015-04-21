package main

import (
	"TskSch/msgQ"
	"TskSch/resultDB"
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
)

type task struct {
	Task_id string
}

type taskp struct {
	Task_id string
	Pid     int
}

func main() {
	for{

		var ids []string
		res1 := []task{}
		res2 := []taskp{}

		//Connecting to mongodb
		session := resultDB.ResultdbInit()
		session.SetMode(mgo.Monotonic, true)
		col := session.DB("TskSch").C("Result")

		//Connecting to msgQ
		Conn := msgQ.RedisInit()

		//Checking liveliness of Scheduler 
		res, err := http.Get("http://127.0.0.1:8001/ping")
		if err == nil && res.Status == "200" {
			body , _ := ioutil.ReadAll(res.Body)
			if string(body) != "" {
				var status interface{}
				err := json.Unmarshal([]byte(body), &status)
				if err == nil {
					if status.(map[string]interface{})["status"].(string) == "alive" {
						fmt.Println("Scheduer is alive")
					}else {
						fmt.Println("Scheduer is not alive")
					}
				}
			}
		}else{
			fmt.Println("Cannot connect to Scheduler")
		}

		//Checking liveliness of Scheduler 
		ress, err1 := http.Get("http://127.0.0.1:8000/ping")
		if err1 == nil && ress.Status == "200" {
			body1 , _ := ioutil.ReadAll(ress.Body)
			if string(body1) != "" {
				var status interface{}
				err := json.Unmarshal([]byte(body1), &status)
				if err == nil {
					if status.(map[string]interface{})["status"].(string) == "alive" {
						fmt.Println("Scheduer is alive")
					}else {
						fmt.Println("Scheduer is not alive")
					}
				}
			}
		}else{
			fmt.Println("Cannot connect to Scheduler")
		}

		//Checking liveness of msgQ to get the list of taskids in msgQ
		err = msgQ.Ping(Conn)
		if err !=nil{
			Ids, Err := redis.Values(Conn.Do("LRANGE", "task", "0", "-1"))
			if Err != nil {
				fmt.Println("Could not able to connect to msgQ",Err)
			}
			for _, val := range Ids {
				ids = append(ids, string(val.([]byte)))
			}
		}

		//Checking liveliness of mongodb
		err = resultDB.Ping(session)
		if err !=nil{ 

			//collecting all the taskids from resultDB which are not executed and not in execution state to check whether they are in msgQ 
			Err := col.Find(bson.M{"executed": false, "exec_stat": false}).Select(bson.M{"task_id":1}).All(&res1)
			if Err != nil {
				fmt.Println("Could not able to connect to mongodb",Err)
			}
			for _, val := range res1{
				flag := In(val.Task_id , ids )
				if flag != true {
					x, err := Conn.Do("RPUSH", "task", val.Task_id)
					if err != nil {
						fmt.Println(x,err)
					}
					fmt.Println("PUSHED" + val.Task_id + "TASK TO msgQ" )
				}
			}

			//collecting all the taskids from resultDB which are not executed , not in execution state and which were poped from executer but not executed
			Err = col.Find(bson.M{"executed": false, "exec_stat": false ,"pid":bson.M{"$gt":0}}).Select(bson.M{"task_id":1,"pid":1}).All(&res2)
			if Err != nil {
				fmt.Println("Could not able to connect to mongodb",Err)
			}
			for _ , val := range res2{
				flag := Isalive(val.Pid)
				if flag != true {
					flag1 := In(val.Task_id , ids )
					if flag1 != true {
						x, err := Conn.Do("RPUSH", "task", val.Task_id)
						if err != nil {
							fmt.Println(x,err)
						}
						fmt.Println("PUSHED" + val.Task_id + "TASK TO msgQ" )
					}
				}
			}

		}else{
			resultDB.Restart(session)
		}
	time.Sleep(time.Second * 100)
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


//type Result struct {
//	Task_id   string //command ID
//	Executed  bool   //Executed staus
//	TOE       string //Time Of Execution
//	TTE       string //Time Taken to Execute
//	Pid       int    // process id of client
//	Exec_Stat bool   //Execution Status
//	output    string
//	err       string
//}

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
)

type task struct {
	Task_id string
}

type taskp struct {
	Task_id string
	Pid     int
}

func main() {
	fmt.Println("OOKK")
	var ids []string
	res1 := []task{}
	res2 := []taskp{}
	session := resultDB.ResultdbInit()
	Conn := msgQ.RedisInit()
	Ids, Err := redis.Values(Conn.Do("LRANGE", "task", "0", "-1"))
	if Err != nil {
		fmt.Println(Err)
	}
	for _, val := range Ids {
		ids = append(ids, string(val.([]byte)))
	}
	session.SetMode(mgo.Monotonic, true)
	col := session.DB("TskSch").C("Result")
	Err = col.Find(bson.M{"executed": false, "exec_stat": false}).Select(bson.M{"task_id":1}).All(&res1)
	if Err != nil {
		fmt.Println(Err)
	}
	Err = col.Find(bson.M{"executed": false, "exec_stat": false ,"pid":bson.M{"$gt":0}}).Select(bson.M{"task_id":1,"pid":1}).All(&res2)
	if Err != nil {
		fmt.Println(Err)
	}
	for _, val := range res1{
		flag := In(val.Task_id , ids )
		if flag != true {
			x, err := Conn.Do("RPUSH", "task", val.Task_id)
			if err != nil {
				fmt.Println(x,err)
			}
		}
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
			}
		}
	}
}
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
func In(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

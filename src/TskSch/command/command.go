package command

import (
	"TskSch/msgQ"
	"github.com/garyburd/redigo/redis"
	"net/http"
	"fmt"
	"io/ioutil"
)

//SEARCHING FOR COMMAND BASED ON THE ID POPED FROM MSG QUEUE
func Search(c redis.Conn, cmd_id string,schedulerPath string) string {
	s := "http://"+schedulerPath+"/askCommand?cmdId=" + cmd_id 
	res, err := http.Get(s)
	if err!=nil{
		fmt.Println("CAN'T CONNECT TO SCHEDULER TO GET THE TASK_CMD OF GIVEN TASK ID ")
		return ""
	}

	body , _ := ioutil.ReadAll(res.Body)
	if string(body) == "" {
		msgQ.Push(c, cmd_id)
		return ""
	}else{
		return string(body)
	}
}

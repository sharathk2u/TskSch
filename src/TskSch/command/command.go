package command

import (
	"TskSch/msgQ"
	"github.com/garyburd/redigo/redis"
	"net/http"
	"fmt"
	"io/ioutil"
)

//SEARCHING FOR COMMAND BASED ON THE ID POPED FROM MSG QUEUE
func Search(c redis.Conn, cmd_id string) string {

	res, err := http.Get("http://127.0.0.1:8001/ask?id=1")
	if err!=nil{
		fmt.Println("CAN'T CONNECT TO SCHEDULER TO GET THE TASK_CMD OF GIVEN TASK ID ",err)
	}

	body , _ := ioutil.ReadAll(res.Body) 

	if string(body) == "" {
		msgQ.Push(c, cmd_id)
		return ""
	}else{
		return string(body)
	}

}

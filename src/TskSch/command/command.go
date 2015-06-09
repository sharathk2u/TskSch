package command

import (
    "TskSch/msgQ"
    "github.com/garyburd/redigo/redis"
    "net/http"
    "fmt"
    "io/ioutil"
	"runtime/debug"
	"TskSch/mailer"
)

//SEARCHING FOR COMMAND BASED ON THE ID POPED FROM MSG QUEUE
func Search(c redis.Conn, cmd_id string,managerPath string, host string,name string) string {
    s := "http://"+managerPath+"/askCommand?cmdId=" + cmd_id+":"+host+"&agentName="+name
    res, err := http.Get(s)
    if err!=nil{
        fmt.Println("CAN'T CONNECT TO SCHEDULER TO GET THE TASK_CMD OF GIVEN TASK ID")
		mailer.Mail("GOSERVE: Unable to connect to the MANAGER", "Unable to establish connection with the Scheduer to get the command \n\n"+ err.Error()+"\n\nStack Trace: --------------------\n\n\n"+string(debug.Stack()))
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


package main

import (
	"TskSch/execute"
	"TskSch/logger"
	"TskSch/msgQ"
	"TskSch/resultDB"
	"TskSch/task"
	"github.com/zenazn/goji"
	"net/http"
	"strconv"
	_ "fmt"
)

func main() {
	//INITIALIZING THE LOG FILE
	file := logger.LogInit()

	//INITIALIZING THE MONGODB
	session := resultDB.ResultdbInit()

	//INITIALIZING THE REDIS DB
	Conn := msgQ.RedisInit()

	//CALLING THE TASK MODULE
	go task.Execute(file, session, Conn)

	//CLOSING ALL THE CONNECTION
	defer func() {
		file.Close()
		session.Close()
		Conn.Close()
	}()
	go Listen_Serve()

	select {}

}

func Listen_Serve() {
	goji.Get("/ping", ping)
	goji.Get("/tasks", gettask)
	goji.Serve()
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("{\"status\":\"alive\"}"))
}

func gettask(w http.ResponseWriter, r *http.Request) {
	var taskIds string = ""
	taskInfo := execute.Get()
	for i, val := range taskInfo{
	//	if(val == true){
		taskIds = taskIds +"{\"" + i + "\":\"" + strconv.FormatBool(val) + "\"},"
	//	}
	}
	w.Write([]byte("{\"taskIds\"" + ":" + taskIds + "}"))
}
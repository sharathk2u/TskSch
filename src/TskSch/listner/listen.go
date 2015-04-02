package main

import (
	"TskSch/execute"
	"TskSch/logger"
	"TskSch/msgQ"
	"TskSch/resultDB"
	"TskSch/task"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"fmt"
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

	m := mux.NewRouter()

	//PING
	m.HandleFunc("/ping",func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("{\"status\":\"alive\"}"))
	}).Methods("GET")

	//CURRENT RUNNING TASKS
	m.HandleFunc("/tasks", func (w http.ResponseWriter, req *http.Request) {
		var taskIds string = ""
		taskInfo := execute.Get()
		for i, val := range taskInfo{
			if(val == true){
			taskIds = taskIds +"{\"" + i + "\":\"" + strconv.FormatBool(val) + "\"},"
			}
		}
		w.Write([]byte("{\"taskIds\"" + ":" + taskIds + "}"))
	}).Methods("GET")

	//RUNNING THE SERVER AT PORT 8000
	err := http.ListenAndServe(":8000", m)
	if err != nil {
		fmt.Println("Error starting server on port.")
		fmt.Println(err)
	}
}



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
	"code.google.com/p/goconf/conf"
)

func main() {
	c, err := conf.ReadConfigFile("../server.conf")
	if err != nil {
		fmt.Println("CAN'T READ CONF FIILE",err)
	}
	username ,_ := c.GetString("resultDB","username")
	password ,_ := c.GetString("resultDB","password")
	host1 ,_ := c.GetString("resultDB","host")
	host2 ,_ := c.GetString("msgQ","host")
	port ,_ := c.GetString("msgQ","port")
	//INITIALIZING THE LOG FILE
	file := logger.LogInit()

	//INITIALIZING THE MONGODB
	session := resultDB.ResultdbInit(username,password,host1)

	//INITIALIZING THE REDIS DB
	Conn := msgQ.RedisInit(host2,port)

	//CALLING THE TASK MODULE
	go task.Execute(file, session, Conn)

	//CLOSING ALL THE CONNECTION
	defer func() {
		file.Close()
		session.Close()
		Conn.Close()
	}()
	
	//TO EXPOSE API's
	go listenServe()

	select {}

}

func listenServe() {

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



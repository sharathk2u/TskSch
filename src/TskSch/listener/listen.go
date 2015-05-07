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
        "io/ioutil"
)

func main() {

        c, err := conf.ReadConfigFile("../server.conf")

        if err != nil {
                fmt.Println("CAN'T READ CONF FIILE",err)
        }

        host1 ,_ := c.GetString("resultDB","host")
        host2 ,_ := c.GetString("msgQ","host")
        port ,_ := c.GetString("msgQ","port")
        agentname ,_ := c.GetString("taskagent","name")
        agenthost ,_ := c.GetString("taskagent","host")
        agentport ,_ := c.GetString("taskagent","port")
        managerhost ,_ := c.GetString("manager","host")
        managerport,_ := c.GetString("manager","port")

        agent := agenthost+":"+agentport+":"+agentname
        s := "http://" + managerhost + ":" + managerport + "/register?agent=" + agent
        res, err := http.Get(s)
        if err != nil{
                fmt.Println("CAN'T CONNECT TO MANAGER")
        }
        body , _ := ioutil.ReadAll(res.Body)
        if string(body) == "ok" {

                //INITIALIZING THE LOG FILE
                file := logger.LogInit()
                //INITIALIZING THE MONGODB
                session := resultDB.ResultdbInit(host1)

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
                go listenServe(agentport)

                select {}

        }else{
                fmt.Println(" NOT VALID TASK AGENT")
        }
}

func listenServe(agentport string){

        m := mux.NewRouter()

        //PING
        m.HandleFunc("/ping",func(w http.ResponseWriter, req *http.Request) {
                w.WriteHeader(200)
                w.Write([]byte("alive"))
        }).Methods("GET")

        //CURRENT RUNNING TASKS
        m.HandleFunc("/tasks", func (w http.ResponseWriter, req *http.Request) {
                var taskIds string = ""
                taskInfo := execute.Get()
                for k, val := range taskInfo{
                        if(val.Value == true){
								taskIds += "{"+
												"Task Id : "+ "\"" + k + "\","+
												"Task Name : "+ "\"" + val.Name +"\","+
												"Value : "+ "\"" + strconv.FormatBool(val.Value) + "\""+"}"                      	
                        }
                }
                w.Write([]byte(taskIds))
        }).Methods("GET")

        //RUNNING THE SERVER AT PORT 8000
        err := http.ListenAndServe(":"+agentport, m)
        if err != nil {
                fmt.Println("Error starting server on port.")
                fmt.Println(err)
        }
}

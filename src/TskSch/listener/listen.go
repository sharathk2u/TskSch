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
    "os"
    "code.google.com/p/goconf/conf"
    "io/ioutil"
	"runtime/debug"
	"TskSch/mailer"
	"archive/zip"
	"io"
	"strings"
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
        	mailer.Mail("GOSERVE: Unable to connect to the MANAGER", "Unable to establish connection with the Manager \n\n"+ err.Error()+"\n\nStack Trace: --------------------\n\n\n"+string(debug.Stack()))
        	return
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
				taskIds += "{"+"Task Id : "+ "\"" + k + "\","+
				"Task Name : "+ "\"" + val.Name +"\","+
				"Value : "+ "\"" + strconv.FormatBool(val.Value) + "\""+"}"                      	
           	}
        }
        w.Write([]byte(taskIds))
    }).Methods("GET")

	//Uploading the file
        m.HandleFunc("/upload",func(w http.ResponseWriter, r *http.Request) {
           fmt.Println("Processing Upload request")
            r.ParseMultipartForm(32 << 20)
            file, handler, err := r.FormFile("uploadfile")
            if err != nil {
                fmt.Println("could not open uploadfile",err)
                return
            }
            defer file.Close()
            tq := strings.Split(handler.Filename,"/")
            filename := tq[len(tq)-1]
            f, err := os.OpenFile("/home/solution/tmp/"+filename, os.O_WRONLY|os.O_CREATE, 0666)
            if err != nil {
                fmt.Println("cant create the zip file inside tmp folder",err)
                return
            }
            defer f.Close()
            _, err = io.Copy(f, file)
            if err != nil {
                fmt.Println("File did not uploaded")
                return
            }
            fmt.Println("File Uploaded to agent and ready for unzip")
            flag := unzip("/home/solution/tmp/"+filename)
            if flag != nil {
                fmt.Println("File did not unziped")
                return
            }else{
            	os.Remove("/home/solution/tmp/"+filename)
            }
        }).Methods("POST")

    
    //RUNNING THE SERVER AT PORT 8000
    err := http.ListenAndServe(":"+agentport, m)
    if err != nil {
        fmt.Println("Error starting server on port.")
        fmt.Println(err)
    }
}
func unzip(filename string) error {
	r, err := zip.OpenReader(filename)
    if err != nil {
        fmt.Println(err)
        return err
    }
    defer r.Close()
    
    err = os.Mkdir(strings.Split(filename,".")[0],0777)
	if err != nil {
		fmt.Println("Unable to create the directory for writing. Check your write access privilege",err)
		return err
	}
    for _, f := range r.File {
        rc, err := f.Open()
        if err != nil {
            fmt.Println(err)
        	return err
        }
        z, err := os.Create(strings.Split(filename,".")[0]+"/"+f.Name)
		if err != nil {
			fmt.Println(err)
			return err
		}
		io.Copy(z, rc)
        defer func() {
			rc.Close()
			z.Close()
        }()
   }
   return nil
}

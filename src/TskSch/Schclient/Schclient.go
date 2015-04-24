package main
import(
	"TskSch/scheduler"
	"TskSch/resultDB"
	"encoding/json"
	"net/http"
	"fmt"
	"github.com/gorilla/mux"
	"time"
	"strconv"
	"strings"
	"code.google.com/p/goconf/conf"
//	"runtime/pprof"
//	"os"
)
func main(){

	c, err := conf.ReadConfigFile("../server.conf")
	if err != nil {
		fmt.Println("CAN'T READ CONF FIILE",err)
	}
	username ,_ := c.GetString("resultDB","username")
	password ,_ := c.GetString("resultDB","password")
	host1 ,_ := c.GetString("resultDB","host")
	host2 ,_ := c.GetString("msgQ","host")
	port ,_ := c.GetString("msgQ","port")

	//INITIALIZING THE MONGODB
	Session := resultDB.ResultdbInit(username , password ,host1)

	//CLOSING ALL THE CONNECTION
	defer func(){
		Session.Close()
	}()
	
	go scheduler.Schedule(Session,username , password ,host1,host2,port)

	go listenServe(username , password ,host1)
	
//	go func (){
//		time.Sleep(100 * time.Second)
//		pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
//	}()
	select{}
}

func listenServe(username string, password string,host1 string) {

	m := mux.NewRouter()

	//INITIALIZING THE MONGODB
	Session := resultDB.ResultdbInit(username , password ,host1)

	//CLOSING ALL THE CONNECTION
	defer func(){
		Session.Close()
	}()

	//PING
	m.HandleFunc("/ping",func(w http.ResponseWriter, req *http.Request) {

		w.WriteHeader(200)
		w.Write([]byte("{\"status\":\"alive\"}"))

	}).Methods("GET")

	//ADDING THE TASKS TO MONGODB
	m.HandleFunc("/addTask",func(w http.ResponseWriter, req *http.Request) {
	
	taskData := req.FormValue("taskData")
	if taskData != "" {
		var taskJs interface{}
		err := json.Unmarshal([]byte(taskData), &taskJs)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Unable to unmarshall taskData", http.StatusBadRequest)
			return
		}
		resultDB.InsertSchedule(Session,taskJs)
	}else{
		http.Error(w, "taskData cannot be empty", http.StatusBadRequest)
	}
	
	}).Methods("GET")

	//UPDATE THE TASK
	m.HandleFunc("/updateTask", func(w http.ResponseWriter, req *http.Request){
		
		taskData := req.FormValue("taskData")
		if taskData != ""{
			var taskJs interface{}
			err := json.Unmarshal([]byte(taskData), &taskJs)
			if err != nil {
				fmt.Println(err)
				http.Error(w, "Unable to unmarshall taskData", http.StatusBadRequest)
				return
			}
			resultDB.Update(Session,taskJs,time.Now())
		}else{
			http.Error(w, "taskData cannot be empty", http.StatusBadRequest)
		}
		
	}).Methods("GET")

	//SEND TASK_CMD when Task Agent asks the Scheduler
	m.HandleFunc("/askCommand", func(w http.ResponseWriter, req *http.Request){
		
		result := req.FormValue("cmdId")
		cmd_id := strings.Split(result,":")[0]
		if cmd_id != ""{
			val , _ := strconv.Atoi(cmd_id)
			cmd := resultDB.Find(val,strings.Split(result,":")[1],strings.Split(result,":")[2],strings.Split(result,":")[3])
			w.Write([]byte(cmd))
		}else{
			http.Error(w, "cmd_id cannot be empty", http.StatusBadRequest)
		}
	
	}).Methods("GET")

	//RUNNING THE SERVER AT PORT 8001
	err := http.ListenAndServe(":8001", m)
	if err != nil {
		fmt.Println("Error starting server on port.")
		fmt.Println(err)
	}
}

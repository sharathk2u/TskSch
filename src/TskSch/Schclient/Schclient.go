package main
import(
//	"TskSch/scheduler"
	"TskSch/resultDB"
	"encoding/json"
	"net/http"
	"fmt"
	"github.com/gorilla/mux"
	"time"
)
func main(){

//	//INITIALIZING THE MONGODB
//	Session := resultDB.ResultdbInit()

//	//CLOSING ALL THE CONNECTION
//	defer func(){
//		Session.Close()
//	}()
//	
//	go scheduler.Schedule(Session)

	go Listen_Serve()
	
	select{}
}

func Listen_Serve() {

	m := mux.NewRouter()

	//INITIALIZING THE MONGODB
	Session := resultDB.ResultdbInit()

	//CLOSING ALL THE CONNECTION
	defer func(){
		Session.Close()
	}()

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
	m.HandleFunc("/ask", func(w http.ResponseWriter, req *http.Request){
		cmd_id := req.FormValue("cmd_id")
		cmd := resultDB.Find(Session , cmd_id) 
		w.Write([]byte(cmd))
	}).Methods("GET")

	//RUNNING THE SERVER AT PORT 8000
	err := http.ListenAndServe(":8000", m)
	if err != nil {
		fmt.Println("Error starting server on port.")
		fmt.Println(err)
	}
}

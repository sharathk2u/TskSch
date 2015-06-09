package main
import(
    "TskSch/scheduler"
    "TskSch/resultDB"
    "TskSch/logger"
	"TskSch/mailer"
	"encoding/json"
    "net/http"
    "fmt"
    "github.com/gorilla/mux"
    "time"
    "strconv"
    "code.google.com/p/goconf/conf"
	"os"
	"path/filepath"
	"bufio"
)

func main(){

        c, err := conf.ReadConfigFile("../server.conf")
        if err != nil {
                fmt.Println("CAN'T READ CONF FIILE",err)
        }
        host1 ,_ := c.GetString("resultDB","host")
        host2 ,_ := c.GetString("msgQ","host")
        port ,_ := c.GetString("msgQ","port")
        
        //INITIALIZING THE MONGODB
        Session := resultDB.ResultdbInit(host1)

        //INITIALIZING THE LOG FILE
        logfile := logger.LogSchInit()

        //CLOSING ALL THE CONNECTION
        defer func(){
                Session.Close()
                logfile.Close()
        }()
        
        go scheduler.Schedule(Session,host1,host2,port,logfile)

        go listenServe(host1,logfile)

        select{}
}
func listenServe(host1 string, logfile *os.File) {

    m := mux.NewRouter()

    //INITIALIZING THE MONGODB
    Session := resultDB.ResultdbInit(host1)

    //CLOSING ALL THE CONNECTION
    defer func(){
            Session.Close()
    }()

    //PING
    m.HandleFunc("/ping",func(w http.ResponseWriter, req *http.Request) {

        w.WriteHeader(200)
        w.Write([]byte("{\"status\":\"alive\"}"))
        LogInfo := logger.Info(logfile)
        LogInfo.Println("SCHEDULAR GOT PINGED")

    }).Methods("GET")
	
	//ADDING THE TASKS TO MONGODB
    m.HandleFunc("/addTask",func(w http.ResponseWriter, req *http.Request) {
		type TaskInfo struct{
            Name string
            Cmd string
            Hour int
            Minute int
            Second int
            Day int
            Week int
            R int
        }
		
		var taskJs TaskInfo
		req.ParseMultipartForm(2000000)	
                    
		if len(req.MultipartForm.Value) != 0 && len(req.MultipartForm.File)!=0 {
			hour, _ := strconv.Atoi(req.FormValue("hour"))
            minute,_ := strconv.Atoi(req.FormValue("minute"))
            second,_ := strconv.Atoi(req.FormValue("second"))
            day,_ := strconv.Atoi(req.FormValue("day"))
            week,_ := strconv.Atoi(req.FormValue("week"))
            r,_ := strconv.Atoi(req.FormValue("r"))
			
			fileDir := req.Form["name"][0]
			err := os.Mkdir("." + string(filepath.Separator) + fileDir,0777)
			if err != nil {
    			fmt.Println("Unable to create the directory for writing. Check your write access privilege",err)
				LogErr := logger.Failure(logfile)
        		LogErr.Println("Unable to create the directory for writing. Check your write access privilege",err)
			}
			
			file := req.MultipartForm.File
			
			LogInfo := logger.Info(logfile)
			LogInfo.Println("STARTED TO ADD THE TASK FILE")
			
			for _ , v := range file {
			 	
			 	o, er := os.Create("." + string(filepath.Separator) + fileDir+string(filepath.Separator)+v[0].Filename)
				if er != nil {
	    			fmt.Println("Unable to create the file for writing. Check your write access privilege",er,o)
	    			LogErr := logger.Failure(logfile)
	    			LogErr.Println("Unable to create the file for writing. Check your write access privilege",er)
	 			}
	 			// write the content from POST to the file
	 			fd , e := v[0].Open()
	 			if e != nil {
	    			 fmt.Println(e)
	 			}	
			 	r1 := bufio.NewReader(fd)
				s, e := Readln(r1)
				str := ""
				for e == nil {
					if s!=""{
						str += s+"\n"
					}
					s,e = Readln(r1)
				}
				o.Write([]byte(str))
				fmt.Println("	FILE : "+ v[0].Filename + " UPLAODED SUCCESSFULLY : ")
			}

			LogInfo.Println("ADDING OF THE TASK FILE COMPLETED")
			
			taskJs = TaskInfo{
		        Name : req.Form["name"][0],
		        Cmd : req.Form["cmd"][0],
		        Hour : hour,
		        Minute : minute,
		        Second : second,
		        Day : day,
		        Week : week,
		        R : r,
			}
			js, _ := json.Marshal(taskJs)
			out := resultDB.InsertSchedule(Session,js)
			LogInfo.Println("TASK ADDED TO THE DATABASE successfully")
			LogInfo.Println("THE ID OF THE TASK ADDED IS : "+ strconv.Itoa(out))
			w.Write([]byte( "{" + "\"status\"" + " : \"Inserted,\"" +"\"Id\""+" : "+ "\""+strconv.Itoa(out)+"\""+"}"))
			mailer.Mail("GOSERVE: Regarding Task addition", taskJs.Name + " ADDED ")
        }else{
            http.Error(w, "taskData cannot be empty", http.StatusBadRequest)
            mailer.Mail("GOSERVE: Regarding Task addition", "Unable to ADD " + taskJs.Name + " Please check the format of addition")
            LogErr := logger.Failure(logfile)
            LogErr.Println("taskData CANNOT BE EMPTY , BAD REQUEST")
        }
    }).Methods("POST")


	//UPDATE THE TASK
    m.HandleFunc("/updateTask", func(w http.ResponseWriter, req *http.Request){

		type TaskInfo struct{
			Id int
            Name string
            Cmd string
            Hour int
            Minute int
            Second int
            Day int
            Week int
            R int
        }
        var taskJs TaskInfo
		req.ParseMultipartForm(2000000)	
                    
		if len(req.MultipartForm.Value) != 0 {
			id , _ := strconv.Atoi(req.FormValue("id"))
			hour, _ := strconv.Atoi(req.FormValue("hour"))
            minute,_ := strconv.Atoi(req.FormValue("minute"))
            second,_ := strconv.Atoi(req.FormValue("second"))
            day,_ := strconv.Atoi(req.FormValue("day"))
            week,_ := strconv.Atoi(req.FormValue("week"))
            r,_ := strconv.Atoi(req.FormValue("r"))
			
			if len(req.MultipartForm.File)!=0 {
				fileDir := req.Form["name"][0]
				err := os.Mkdir("." + string(filepath.Separator) + fileDir,0777)
				if err != nil {
					fmt.Println("Unable to create the directory for writing. Check your write access privilege",err)
					LogErr := logger.Failure(logfile)
					LogErr.Println("Unable to create the directory for writing. Check your write access privilege",err)
				}
			
				file := req.MultipartForm.File

				LogInfo := logger.Info(logfile)
				LogInfo.Println("STARTED UPDATING THE TASK FILE")

				for _ , v := range file {
				 	
				 	o, er := os.Create("." + string(filepath.Separator) + fileDir+string(filepath.Separator)+v[0].Filename)
					if er != nil {
						fmt.Println("Unable to create the file for writing. Check your write access privilege",er,o)
		 				LogErr := logger.Failure(logfile)
		 				LogErr.Println("Unable to create the file for writing. Check your write access privilege",er)
		 			}
		 			// write the content from POST to the file
					fd , e := v[0].Open()
		 			if e != nil {
						 fmt.Println(e)
		 			}	
				 	r1 := bufio.NewReader(fd)
					s, e := Readln(r1)
					str := ""
					for e == nil {
						if s!=""{
							str += s+"\n"
						}
						s,e = Readln(r1)
					}
					o.Write([]byte(str))
					fmt.Println("	FILE : "+ v[0].Filename + " UPLAODED SUCCESSFULLY : ")
				}
				LogInfo.Println("TASK FILE UPDATED")
			}


			taskJs = TaskInfo{
				Id : id,
				Name : req.Form["name"][0],
				Cmd : req.Form["cmd"][0],
				Hour : hour,
				Minute : minute,
				Second : second,
				Day : day,
				Week : week,
				R : r,
			}
			js, _ := json.Marshal(taskJs)
			resultDB.Update(Session,js,time.Now())
			LogInfo := logger.Info(logfile)
			LogInfo.Println("UPDATED TASK ADDED TO DATABASE successfully")
			w.Write([]byte( "{" + "\"status\"" + " : \"updated\""+"}"))
			mailer.Mail("GOSERVE: Regarding Task updation", taskJs.Name + " UPDATED ")
		}else{
				http.Error(w, "taskData cannot be empty", http.StatusBadRequest)
				mailer.Mail("GOSERVE: Regarding Task updation", "Unable to UPDATE " + taskJs.Name + " Please check the format of addition ")
		}

    }).Methods("GET")

    //RUNNING THE SERVER AT PORT 8001
    err := http.ListenAndServe(":8001", m)
    if err != nil {
        fmt.Println("Error starting server on port.",err)
        LogErr := logger.Failure(logfile)
        LogErr.Println("Error starting server on port.",err)
    }
}
func Readln(r *bufio.Reader) (string, error) {
        var (
                isRead bool = true
                Err error = nil
                line, ln []byte
        )
        for isRead && Err == nil {
                line, isRead , Err = r.ReadLine()
                ln = append(ln, line...)
	}
        return string(ln),Err
}
